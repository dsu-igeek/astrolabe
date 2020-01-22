/*
 * Copyright 2019 the Astrolabe contributors
 * SPDX-License-Identifier: Apache-2.0
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package kubernetes

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/astrolabe/pkg/astrolabe"
	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	"github.com/vmware-tanzu/velero/pkg/backup"
	"github.com/vmware-tanzu/velero/pkg/builder"
	"github.com/vmware-tanzu/velero/pkg/client"
	"github.com/vmware-tanzu/velero/pkg/discovery"
	"github.com/vmware-tanzu/velero/pkg/podexec"
	"io"
	v1 "k8s.io/api/core/v1"
	"os"
)

type KubernetesNamespaceProtectedEntity struct {
	knpetm    * KubernetesNamespaceProtectedEntityTypeManager
	id        astrolabe.ProtectedEntityID
	namespace *v1.Namespace
	logger    logrus.FieldLogger
}

func (this *KubernetesNamespaceProtectedEntity) GetDataReader(context.Context) (io.ReadCloser, error) {
	return nil, nil
}

func (this *KubernetesNamespaceProtectedEntity) GetMetadataReader(context.Context) (io.Reader, error) {
	return nil, nil
}

func NewKubernetesNamespaceProtectedEntity(knpetm *KubernetesNamespaceProtectedEntityTypeManager,
	namespace *v1.Namespace) (*KubernetesNamespaceProtectedEntity, error) {
	nsPEID := astrolabe.NewProtectedEntityID("k8sns", namespace.Name)
	returnPE := KubernetesNamespaceProtectedEntity{
		knpetm:    knpetm,
		id:        nsPEID,
		namespace: namespace,
		logger:    knpetm.logger,
	}
	return &returnPE, nil
}

func (this *KubernetesNamespaceProtectedEntity) GetInfo(ctx context.Context) (astrolabe.ProtectedEntityInfo, error) {
	return nil, nil
}
func (this *KubernetesNamespaceProtectedEntity) GetCombinedInfo(ctx context.Context) ([]astrolabe.ProtectedEntityInfo, error) {
	return nil, nil

}

func (this *KubernetesNamespaceProtectedEntity) Snapshot(ctx context.Context) (astrolabe.ProtectedEntitySnapshotID, error) {
	vc := client.VeleroConfig{}
	f := client.NewFactory("astrolabe", vc)

	veleroClient, err := f.Client()
	if err != nil {
		return astrolabe.ProtectedEntitySnapshotID{}, err
	}
	discoveryClient := veleroClient.Discovery()

	dynamicClient, err := f.DynamicClient()
	if err != nil {
		return astrolabe.ProtectedEntitySnapshotID{}, err
	}

	discoveryHelper, err := discovery.NewHelper(discoveryClient, this.logger)
	if err != nil {
		return astrolabe.ProtectedEntitySnapshotID{}, err
	}
	dynamicFactory := client.NewDynamicFactory(dynamicClient)

	kubeClient, err := f.KubeClient()
	if err != nil {
		return astrolabe.ProtectedEntitySnapshotID{}, err
	}

	kubeClientConfig, err := f.ClientConfig()
	if err != nil {
		return astrolabe.ProtectedEntitySnapshotID{}, err
	}
	podCommandExecutor := podexec.NewPodCommandExecutor(kubeClientConfig, kubeClient.CoreV1().RESTClient())

	k8sBackupper, err := backup.NewKubernetesBackupper(discoveryHelper,
		dynamicFactory,
		podCommandExecutor,
		nil,
		0)
	if err != nil {
		return astrolabe.ProtectedEntitySnapshotID{}, err
	}

	snapshotUUID, err := uuid.NewRandom()
	if err != nil {
		return astrolabe.ProtectedEntitySnapshotID{}, err
	}

	backupParams := 	builder.ForBackup(velerov1.DefaultNamespace, "astrolabe-" + snapshotUUID.String()).Result()

	request := backup.Request{
		Backup:                    backupParams,
	}


	snapshotFileName := "/tmp/" + snapshotUUID.String()
	backupFile, err := os.Create(snapshotFileName)
	if err != nil {
		return astrolabe.ProtectedEntitySnapshotID{}, err
	}

	k8sBackupper.Backup(this.logger, &request, backupFile, nil, nil)
	//	Backup(logger logrus.FieldLogger, backup *Request, backupFile io.Writer, actions []velero.BackupItemAction, volumeSnapshotterGetter VolumeSnapshotterGetter) error
	return astrolabe.ProtectedEntitySnapshotID{}, nil
}
func (this *KubernetesNamespaceProtectedEntity) ListSnapshots(ctx context.Context) ([]astrolabe.ProtectedEntitySnapshotID, error) {
	return nil, nil

}
func (this *KubernetesNamespaceProtectedEntity) DeleteSnapshot(ctx context.Context,
	snapshotToDelete astrolabe.ProtectedEntitySnapshotID) (bool, error) {
	return false, nil

}
func (this *KubernetesNamespaceProtectedEntity) GetInfoForSnapshot(ctx context.Context,
	snapshotID astrolabe.ProtectedEntitySnapshotID) (*astrolabe.ProtectedEntityInfo, error) {
	return nil, nil

}

func (this *KubernetesNamespaceProtectedEntity) GetComponents(ctx context.Context) ([]astrolabe.ProtectedEntity, error) {
	return nil, nil

}

func (this *KubernetesNamespaceProtectedEntity) GetID() astrolabe.ProtectedEntityID {
	return this.id
}
