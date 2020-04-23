package psql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/vmware-tanzu/astrolabe/pkg/astrolabe"
	"io"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type PSQLProtectedEntity struct {
	id        astrolabe.ProtectedEntityID
	petm      *PSQLProtectedEntityTypeManager
	namespace string
}

func NewPSQLProtectedEntity(id astrolabe.ProtectedEntityID, petm *PSQLProtectedEntityTypeManager) PSQLProtectedEntity {
	return PSQLProtectedEntity{
		id:   id,
		petm: petm,
	}
}

func (this PSQLProtectedEntity) GetInfo(ctx context.Context) (astrolabe.ProtectedEntityInfo, error) {
	psql, err := this.petm.getPostgresqlForPEID(ctx, this.id)

	if err != nil {
		return nil, errors.Wrapf(err, "getPostgresqlForPEID failed for %s", this.id.String())
	}

	dataS3Transport, err := astrolabe.NewS3DataTransportForPEID(this.id, this.petm.s3Config)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create S3 data transport")
	}

	data := []astrolabe.DataTransport{
		dataS3Transport,
	}

	mdS3Transport, err := astrolabe.NewS3MDTransportForPEID(this.id, this.petm.s3Config)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create S3 md transport")
	}

	md := []astrolabe.DataTransport{
		mdS3Transport,
	}

	combinedS3Transport, err := astrolabe.NewS3CombinedTransportForPEID(this.id, this.petm.s3Config)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create S3 combined transport")
	}

	combined := []astrolabe.DataTransport{
		combinedS3Transport,
	}

	retVal := astrolabe.NewProtectedEntityInfo(
		this.id,
		psql.Name,
		-1,
		data,
		md,
		combined,
		[]astrolabe.ProtectedEntityID{})
	return retVal, nil
}

func (this PSQLProtectedEntity) GetCombinedInfo(ctx context.Context) ([]astrolabe.ProtectedEntityInfo, error) {
	panic("implement me")
}

func (this PSQLProtectedEntity) Snapshot(ctx context.Context) (astrolabe.ProtectedEntitySnapshotID, error) {
	if this.id.HasSnapshot() {
		return astrolabe.ProtectedEntitySnapshotID{}, errors.New(fmt.Sprintf("pe %s is a snapshot, cannot snapshot again", this.id.String()))
	}
	snapshotUUID, err := uuid.NewRandom()
	if err != nil {
		return astrolabe.ProtectedEntitySnapshotID{}, errors.Wrap(err, "Failed to create new UUID")
	}
	snapshotID := astrolabe.NewProtectedEntitySnapshotID(snapshotUUID.String())

	/*
		pod, err := this.petm.KubeClient.Pods(this.namespace).Create(&v1.Pod{
			TypeMeta:   meta_v1.TypeMeta{},
			ObjectMeta: meta_v1.ObjectMeta{
				Name: "snapshot-pg-"+snapshotID.String(),
			},
			Spec:       v1.PodSpec{},
			Status:     v1.PodStatus{},
		})
		pod.
		panic("implement me")
	*/
	peSnapshotDir := filepath.Join(this.petm.snapshotsDir, this.id.String())
	peSnapshotDirInfo, err := os.Stat(peSnapshotDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(peSnapshotDir, 0700)
			if err != nil {
				return astrolabe.ProtectedEntitySnapshotID{}, errors.Wrapf(err, "could not create snapshot directory %s", peSnapshotDir)
			}
			peSnapshotDirInfo, err = os.Stat(peSnapshotDir)
		}
		if err != nil {
			return astrolabe.ProtectedEntitySnapshotID{}, errors.Wrapf(err, "stat on snapshotsdir %s failed", peSnapshotDir)
		}
	}
	if !peSnapshotDirInfo.IsDir() {
		return astrolabe.ProtectedEntitySnapshotID{}, errors.New(fmt.Sprintf("snapshot directory %s for pe %s is not a directory", peSnapshotDir, this.id.String()))
	}
	snapshotPEID := this.id.IDWithSnapshot(snapshotID)

	snapshotMDFilename := mdFilenameForSnapshot(snapshotPEID)
	snapshotMDPath := filepath.Join(peSnapshotDir, snapshotMDFilename)
	snapshotMDFile, err := os.Create(snapshotMDPath)
	if err != nil {
		return astrolabe.ProtectedEntitySnapshotID{}, errors.Wrapf(err, "Failed to create md file %s", snapshotMDPath)
	}

	mdReader, err := this.GetMetadataReader(ctx)
	_, err = io.Copy(snapshotMDFile, mdReader)

	snapshotDataFilename := dataFilenameForSnapshot(snapshotPEID)
	snapshotDataPath := filepath.Join(peSnapshotDir, snapshotDataFilename)

	snapshotDataFile, err := os.Create(snapshotDataPath)
	if err != nil {
		return astrolabe.ProtectedEntitySnapshotID{}, errors.Wrapf(err, "Failed to create data file %s", snapshotDataPath)
	}

	dataReader, err := this.GetDataReader(ctx)
	if err != nil {
		return astrolabe.ProtectedEntitySnapshotID{}, errors.Wrap(err, "Failed to get cmd's stdout")
	}

	_, err = io.Copy(snapshotDataFile, dataReader)

	if err != nil {
		return astrolabe.ProtectedEntitySnapshotID{}, errors.Wrap(err, "Command failed")
	}

	return snapshotID, nil
}

type copyResults struct {
	written int64
	err     error
}

func copyAndNotify(dst io.WriteCloser, src io.ReadCloser, finished chan copyResults) {
	written, err := io.Copy(dst, src)
	dst.Close()
	src.Close()

	finished <- copyResults{
		written: written,
		err:     err,
	}
}

const snapshotDataExtension = ".psql-snap-data"

func dataFilenameForSnapshot(snapshotPEID astrolabe.ProtectedEntityID) string {
	return snapshotPEID.String() + snapshotDataExtension
}

const snapshotMDExtension = ".psql-snap-md"

func mdFilenameForSnapshot(snapshotPEID astrolabe.ProtectedEntityID) string {
	return snapshotPEID.String() + snapshotMDExtension
}

func snapshotPEIDForFilename(filename string) (astrolabe.ProtectedEntityID, error) {
	if !strings.HasSuffix(filename, snapshotDataExtension) {
		return astrolabe.ProtectedEntityID{}, errors.New(fmt.Sprintf("%s does not end with %s", filename, snapshotDataExtension))
	}
	peidStr := strings.TrimSuffix(filename, snapshotDataExtension)
	return astrolabe.NewProtectedEntityIDFromString(peidStr)
}

func (this PSQLProtectedEntity) ListSnapshots(ctx context.Context) ([]astrolabe.ProtectedEntitySnapshotID, error) {
	if this.id.HasSnapshot() {
		return nil, errors.New(fmt.Sprintf("pe %s is a snapshot, cannot list snapshots", this.id.String()))
	}
	peSnapshotDir := filepath.Join(this.petm.snapshotsDir, this.id.String())
	peSnapshotDirInfo, err := os.Stat(peSnapshotDir)
	var returnIDs = make([]astrolabe.ProtectedEntitySnapshotID, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return []astrolabe.ProtectedEntitySnapshotID{}, nil
		}
		return returnIDs, errors.Wrapf(err, "could not snap snapshot dir %s", peSnapshotDir)
	}
	if !peSnapshotDirInfo.IsDir() {
		return []astrolabe.ProtectedEntitySnapshotID{}, errors.New(fmt.Sprintf("snapshot directory %s for pe %s is not a directory", peSnapshotDir, this.id.String()))
	}
	fileInfo, err := ioutil.ReadDir(peSnapshotDir)
	if err != nil {
		return returnIDs, errors.Wrapf(err, "could not read snapshot dir %s", peSnapshotDir)
	}

	for _, curFileInfo := range fileInfo {
		if strings.HasSuffix(curFileInfo.Name(), snapshotDataExtension) {
			peid, err := snapshotPEIDForFilename(curFileInfo.Name())
			if err != nil {

			} else {
				returnIDs = append(returnIDs, peid.GetSnapshotID())
			}
		}
	}
	return returnIDs, nil
}

func (this PSQLProtectedEntity) DeleteSnapshot(ctx context.Context, snapshotToDelete astrolabe.ProtectedEntitySnapshotID) (bool, error) {
	panic("implement me")
}

func (this PSQLProtectedEntity) GetInfoForSnapshot(ctx context.Context, snapshotID astrolabe.ProtectedEntitySnapshotID) (*astrolabe.ProtectedEntityInfo, error) {
	panic("implement me")
}

func (this PSQLProtectedEntity) GetComponents(ctx context.Context) ([]astrolabe.ProtectedEntity, error) {
	return []astrolabe.ProtectedEntity{}, nil
}

func (this PSQLProtectedEntity) GetID() astrolabe.ProtectedEntityID {
	return this.id
}

func (this PSQLProtectedEntity) GetDataReader(ctx context.Context) (io.ReadCloser, error) {
	if !this.id.HasSnapshot() {
		psql, err := this.petm.getPostgresqlForPEID(ctx, this.id)
		if err != nil {
			return nil, errors.Wrapf(err, "could not retrive psql resource for %s", this.id.String())
		}
		namespace := psql.Namespace
		pghost := psql.ObjectMeta.Name
		pgsecret, err := this.petm.KubeClient.Secrets(namespace).Get("postgres."+pghost+".credentials", metav1.GetOptions{})
		if err != nil {
			return nil, errors.Wrap(err, "could not retrieve secret")
		}
		pguser := string(pgsecret.Data["username"])
		pgpassword := string(pgsecret.Data["password"])
		fmt.Printf("pguser = %s, pgpassword = %s\n", pguser, pgpassword)
		dumpUUID, err := uuid.NewRandom()
		if err != nil {
			return nil, errors.Wrap(err, "could not create UUID")
		}
		podName := "snapshot-pg-" + dumpUUID.String()
		cmd := exec.Command("/usr/bin/kubectl", "run", "-n", namespace, podName, "--image=dpcpinternal/pg-dump:0.0.5",
			"--env", "PGPASSWORD=" + pgpassword, "--env",
			"PGHOST=" + pghost, "--env", "PGUSER=" +pguser, "-it", "--restart=Never", "--rm")

		cmdStdout, err := cmd.StdoutPipe()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get cmd's stdout")
		}
		err = cmd.Start()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to start command")
		}

		return cmdStdout, nil
	}
	peSnapshotDir := filepath.Join(this.petm.snapshotsDir, this.id.GetBaseID().String())
	peSnapshotDirInfo, err := os.Stat(peSnapshotDir)
	if err != nil {
		return nil, errors.Wrapf(err, "stat on snapshotsdir %s failed", peSnapshotDir)
	}
	if !peSnapshotDirInfo.IsDir() {
		return nil, errors.New(fmt.Sprintf("snapshot directory %s for pe %s is not a directory", peSnapshotDir, this.id.String()))
	}
	snapshotFilename := dataFilenameForSnapshot(this.id)
	snapshotPath := filepath.Join(peSnapshotDir, snapshotFilename)
	reader, err := os.Open(snapshotPath)
	return reader, err
}

func (this PSQLProtectedEntity) GetMetadataReader(ctx context.Context) (io.ReadCloser, error) {
	if !this.id.HasSnapshot() {
		psql, err := this.petm.getPostgresqlForPEID(ctx, this.id)
		if err != nil {
			return nil, errors.Wrapf(err, "could not retrive psql resource for %s", this.id.String())
		}
		psqlBytes, err := json.Marshal(psql)
		return ioutil.NopCloser(bytes.NewReader(psqlBytes)), nil
	}
	peSnapshotDir := filepath.Join(this.petm.snapshotsDir, this.id.GetBaseID().String())
	peSnapshotDirInfo, err := os.Stat(peSnapshotDir)
	if err != nil {
		return nil, errors.Wrapf(err, "stat on snapshotsdir %s failed", peSnapshotDir)
	}
	if !peSnapshotDirInfo.IsDir() {
		return nil, errors.New(fmt.Sprintf("snapshot directory %s for pe %s is not a directory", peSnapshotDir, this.id.String()))
	}
	snapshotFilename := mdFilenameForSnapshot(this.id)
	snapshotPath := filepath.Join(peSnapshotDir, snapshotFilename)
	reader, err := os.Open(snapshotPath)
	return reader, err
}
