package psql

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/vmware-tanzu/astrolabe/pkg/astrolabe"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type PSQLProtectedEntity struct {
	id astrolabe.ProtectedEntityID
	petm * PSQLProtectedEntityTypeManager
	namespace string
}

func NewPSQLProtectedEntity(id astrolabe.ProtectedEntityID, petm * PSQLProtectedEntityTypeManager) PSQLProtectedEntity {
	return PSQLProtectedEntity{
		id:   id,
		petm: petm,
	}
}

func (this PSQLProtectedEntity) GetInfo(ctx context.Context) (astrolabe.ProtectedEntityInfo, error) {
	panic("implement me")
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
	snapshotFilename := filenameForSnapshot(snapshotPEID)
	snapshotPath := filepath.Join(peSnapshotDir, snapshotFilename)

	snapshotFile, err := os.Create(snapshotPath)

	podName := "snapshot-pg-"+snapshotID.String()
	cmd := exec.Command("/usr/bin/kubectl", "run", "-n", "postgres-test", podName, "--image=dpcpinternal/pg-dump:0.0.5",
		"--env", "PGPASSWORD=2sw82j65LuzV2oC0ZoXhvtCt53daed3A6QrLyTSa42nXp91Nsyep1KlGw9woJDSs", "--env",
		"PGHOST=acid-minimal-cluster", "-it")

	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		return astrolabe.ProtectedEntitySnapshotID{}, errors.Wrap(err, "Failed to get cmd's stdout")
	}

	finished := make(chan copyResults)
	go copyAndNotify(snapshotFile, cmdStdout, finished)
	err = cmd.Start()
	if err != nil {
		return astrolabe.ProtectedEntitySnapshotID{}, errors.Wrap(err, "Failed to start command")
	}
	err = cmd.Wait()
	if err != nil {
		return astrolabe.ProtectedEntitySnapshotID{}, errors.Wrap(err, "Command failed")
	}
	results := <- finished
	if results.err != nil {
		return astrolabe.ProtectedEntitySnapshotID{}, errors.Wrapf(err, "Copy failed, copied = %d", results.written)
	}
	return snapshotID, nil
}

type copyResults struct {
	written int64
	err error
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
const snapshotExtension = ".psql-snap"

func filenameForSnapshot(snapshotPEID astrolabe.ProtectedEntityID) string {
		return snapshotPEID.String() + snapshotExtension
}

func snapshotPEIDForFilename(filename string) (astrolabe.ProtectedEntityID, error) {
	if !strings.HasSuffix(filename, snapshotExtension) {
		return astrolabe.ProtectedEntityID{}, errors.New(fmt.Sprintf("%s does not end with %s", filename, snapshotExtension))
	}
	peidStr := strings.TrimSuffix(filename, snapshotExtension)
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
		if strings.HasSuffix(curFileInfo.Name(), snapshotExtension) {
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
	panic("implement me")
}

func (this PSQLProtectedEntity) GetID() astrolabe.ProtectedEntityID {
	return this.id
}

func (this PSQLProtectedEntity) GetDataReader(ctx context.Context) (io.ReadCloser, error) {
	if !this.id.HasSnapshot() {
		// TODO - go ahead an do a pg dump here directly
		return nil, errors.New("Cannot read from non-snapshot")
	}
	peSnapshotDir := filepath.Join(this.petm.snapshotsDir, this.id.GetBaseID().String())
	peSnapshotDirInfo, err := os.Stat(peSnapshotDir)
	if err != nil {
		return nil, errors.Wrapf(err, "stat on snapshotsdir %s failed", peSnapshotDir)
	}
	if !peSnapshotDirInfo.IsDir() {
		return nil, errors.New(fmt.Sprintf("snapshot directory %s for pe %s is not a directory", peSnapshotDir, this.id.String()))
	}
	snapshotFilename := filenameForSnapshot(this.id)
	snapshotPath := filepath.Join(peSnapshotDir, snapshotFilename)
	reader, err := os.Open(snapshotPath)
	return reader, err
}

func (this PSQLProtectedEntity) GetMetadataReader(ctx context.Context) (io.ReadCloser, error) {
	panic("implement me")
}
