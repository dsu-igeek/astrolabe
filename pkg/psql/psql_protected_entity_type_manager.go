package psql

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/vmware-tanzu/astrolabe/pkg/astrolabe"
	v1 "github.com/zalando/postgres-operator/pkg/apis/acid.zalan.do/v1"
	"github.com/zalando/postgres-operator/pkg/util/k8sutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

type PSQLProtectedEntityTypeManager struct {
	KubeClient       k8sutil.KubernetesClient
	watchedNamespace string
	snapshotsDir     string
}

func NewPSQLProtectedEntityTypeManager() (PSQLProtectedEntityTypeManager, error) {
	restConfig, err := clientcmd.BuildConfigFromFlags("", "/home/dsmithuchida/.kube/config")
	if err != nil {
		return PSQLProtectedEntityTypeManager{}, err
	}
	kubeClient, err := k8sutil.NewFromConfig(restConfig)

	if err != nil {
		return PSQLProtectedEntityTypeManager{}, err
	}

	snapshotsDir := "/tmp/psql"

	snapshotsDirInfo, err := os.Stat(snapshotsDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(snapshotsDir, 0700)
			if err != nil {
				return PSQLProtectedEntityTypeManager{}, errors.Wrapf(err, "could not create snapshot directory %s", snapshotsDir)
			}
			snapshotsDirInfo, err = os.Stat(snapshotsDir)
		}
		if err != nil {
			return PSQLProtectedEntityTypeManager{}, errors.Wrapf(err, "stat on snapshotsdir %s failed", snapshotsDir)
		}
	}
	if !snapshotsDirInfo.IsDir() {
		return PSQLProtectedEntityTypeManager{}, errors.New(fmt.Sprintf("configured snapshots dir %s is not a directory", snapshotsDir))
	}

	returnPETM := PSQLProtectedEntityTypeManager{
		KubeClient:       kubeClient,
		watchedNamespace: "postgres-test",
		snapshotsDir:     snapshotsDir,
	}

	return returnPETM, nil
}

func (this PSQLProtectedEntityTypeManager) GetTypeName() string {
	return "psql"
}

func (this PSQLProtectedEntityTypeManager) GetProtectedEntity(ctx context.Context, id astrolabe.ProtectedEntityID) (astrolabe.ProtectedEntity, error) {
	_, err := this.getPostgresqlForPEID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve instance info")
	}
	return NewPSQLProtectedEntity(id, &this), nil

}

func (this PSQLProtectedEntityTypeManager) getPostgresqlForPEID(crx context.Context, id astrolabe.ProtectedEntityID) (v1.Postgresql, error) {
	var options metav1.ListOptions
	list, err := this.KubeClient.AcidV1ClientSet.AcidV1().Postgresqls(this.watchedNamespace).List(options)
	if err != nil {
		return v1.Postgresql{}, errors.Wrap(err, "could not retrieve postgres instances")
	}
	for _, curPSQL := range list.Items {
		if string(curPSQL.UID) == id.GetID() {
			return curPSQL, nil
		}
	}
	return v1.Postgresql{}, errors.New("Not found")
}

func (this PSQLProtectedEntityTypeManager) GetProtectedEntities(ctx context.Context) ([]astrolabe.ProtectedEntityID, error) {
	var options metav1.ListOptions
	list, err := this.KubeClient.AcidV1ClientSet.AcidV1().Postgresqls(this.watchedNamespace).List(options)
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve postgres instances")
	}

	var returnIDs []astrolabe.ProtectedEntityID
	for _, curPSQL := range list.Items {
		fmt.Printf("%v\n", curPSQL)
		id := astrolabe.NewProtectedEntityID(this.GetTypeName(), string(curPSQL.UID))
		returnIDs = append(returnIDs, id)
	}
	return returnIDs, nil
}

func (this PSQLProtectedEntityTypeManager) Copy(ctx context.Context, pe astrolabe.ProtectedEntity, options astrolabe.CopyCreateOptions) (astrolabe.ProtectedEntity, error) {
	panic("implement me")
}

func (this PSQLProtectedEntityTypeManager) CopyFromInfo(ctx context.Context, info astrolabe.ProtectedEntityInfo, options astrolabe.CopyCreateOptions) (astrolabe.ProtectedEntity, error) {
	panic("implement me")
}

func (this PSQLProtectedEntityTypeManager) Delete(ctx context.Context, id astrolabe.ProtectedEntityID) error {
	panic("implement me")
}
