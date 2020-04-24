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
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/astrolabe/pkg/astrolabe"
	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	"github.com/vmware-tanzu/velero/pkg/backup"
	"github.com/vmware-tanzu/velero/pkg/builder"
	"github.com/vmware-tanzu/velero/pkg/client"
	"github.com/vmware-tanzu/velero/pkg/discovery"
	"github.com/vmware-tanzu/velero/pkg/podexec"
	"io"
)

type KubernetesNamespaceProtectedEntity struct {
	petm      * KubernetesNamespaceProtectedEntityTypeManager
	id        astrolabe.ProtectedEntityID
	name      string
	logger    logrus.FieldLogger
}

func (this *KubernetesNamespaceProtectedEntity) GetDataReader(context.Context) (io.ReadCloser, error) {
	if !this.id.HasSnapshot() {
		vc := client.VeleroConfig{}
		f := client.NewFactory("astrolabe", vc)

		veleroClient, err := f.Client()
		if err != nil {
			return nil, err
		}
		discoveryClient := veleroClient.Discovery()

		dynamicClient, err := f.DynamicClient()
		if err != nil {
			return nil, err
		}

		discoveryHelper, err := discovery.NewHelper(discoveryClient, this.logger)
		if err != nil {
			return nil, err
		}
		dynamicFactory := client.NewDynamicFactory(dynamicClient)

		kubeClient, err := f.KubeClient()
		if err != nil {
			return nil, err
		}

		kubeClientConfig, err := f.ClientConfig()
		if err != nil {
			return nil, err
		}
		podCommandExecutor := podexec.NewPodCommandExecutor(kubeClientConfig, kubeClient.CoreV1().RESTClient())

		k8sBackupper, err := backup.NewKubernetesBackupper(discoveryHelper,
			dynamicFactory,
			podCommandExecutor,
			nil,
			0)
		if err != nil {
			return nil, err
		}

		snapshotUUID, err := uuid.NewRandom()
		if err != nil {
			return nil, err
		}

		reader, writer := io.Pipe()
		backupParams := 	builder.ForBackup(velerov1.DefaultNamespace, "astrolabe-" + snapshotUUID.String()).
			IncludedNamespaces(this.name).Result()

		request := backup.Request{
			Backup:                    backupParams,
		}

		go this.runBackup(k8sBackupper, request, writer)

		return reader, nil
	}
	return this.petm.internalRepo.GetDataReaderForSnapshot(this.id)
}

func (this * KubernetesNamespaceProtectedEntity)runBackup(k8sBackupper backup.Backupper, request backup.Request, writer io.WriteCloser) {
	defer writer.Close()
	k8sBackupper.Backup(this.logger, &request, writer, nil, nil)
}

func (this *KubernetesNamespaceProtectedEntity) GetMetadataReader(context.Context) (io.ReadCloser, error) {
	return nil, nil
}

func NewKubernetesNamespaceProtectedEntity(petm *KubernetesNamespaceProtectedEntityTypeManager, nsPEID astrolabe.ProtectedEntityID,
	name string) (*KubernetesNamespaceProtectedEntity, error) {
	returnPE := KubernetesNamespaceProtectedEntity{
		petm:      petm,
		id:        nsPEID,
		logger:    petm.logger,
		name: name,
	}
	return &returnPE, nil
}

func (this *KubernetesNamespaceProtectedEntity) GetInfo(ctx context.Context) (astrolabe.ProtectedEntityInfo, error) {

	data := []astrolabe.DataTransport{
		astrolabe.NewS3DataTransportForPEID(this.id, this.petm.s3Config),
	}

	md := []astrolabe.DataTransport{
//		astrolabe.NewS3MDTransportForPEID(this.id, this.petm.s3Config),
	}

	combined := []astrolabe.DataTransport{
		astrolabe.NewS3CombinedTransportForPEID(this.id, this.petm.s3Config),
	}
	retVal := astrolabe.NewProtectedEntityInfo(
		this.id,
		this.name,
		-1,
		data,
		md,
		combined,
		[]astrolabe.ProtectedEntityID{})
	return retVal, nil
}

func (this *KubernetesNamespaceProtectedEntity) GetCombinedInfo(ctx context.Context) ([]astrolabe.ProtectedEntityInfo, error) {
	return nil, nil

}

func (this *KubernetesNamespaceProtectedEntity) Snapshot(ctx context.Context) (astrolabe.ProtectedEntitySnapshotID, error) {
	if this.id.HasSnapshot() {
		return astrolabe.ProtectedEntitySnapshotID{}, errors.New(fmt.Sprintf("pe %s is a snapshot, cannot snapshot again", this.id.String()))
	}
	snapshotUUID, err := uuid.NewRandom()
	if err != nil {
		return astrolabe.ProtectedEntitySnapshotID{}, errors.Wrap(err, "Failed to create new UUID")
	}
	snapshotID := astrolabe.NewProtectedEntitySnapshotID(snapshotUUID.String())


	err = this.petm.internalRepo.WriteProtectedEntity(ctx, this, snapshotID)
	if err != nil {
		return astrolabe.ProtectedEntitySnapshotID{}, errors.Wrap(err, "Failed to create new snapshot")
	}
	return snapshotID, nil
}

func (this *KubernetesNamespaceProtectedEntity) ListSnapshots(ctx context.Context) ([]astrolabe.ProtectedEntitySnapshotID, error) {
	return this.petm.internalRepo.ListSnapshotsForPEID(this.id)

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
