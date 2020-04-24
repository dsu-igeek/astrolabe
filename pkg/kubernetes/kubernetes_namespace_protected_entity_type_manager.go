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
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/astrolabe/pkg/astrolabe"
	"github.com/vmware-tanzu/astrolabe/pkg/localsnap"
	v1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesNamespaceProtectedEntityTypeManager struct {
	clientset  *kubernetes.Clientset
	logger     logrus.FieldLogger
	s3Config   astrolabe.S3Config
	internalRepo localsnap.LocalSnapshotRepo
}
const 	SnapshotsDirKey = "snapshotsDir"
const typename = "k8sns"

func NewKubernetesNamespaceProtectedEntityTypeManagerFromConfig(params map[string]interface{}, s3Config astrolabe.S3Config,
	logger logrus.FieldLogger) (*KubernetesNamespaceProtectedEntityTypeManager, error) {
	masterURLObj := params["masterURL"]
	masterURL := ""
	if masterURLObj != nil {
		masterURL = masterURLObj.(string)
	}

	snapshotsDir, hasSnapshotsDir := params[SnapshotsDirKey].(string)
	if !hasSnapshotsDir {
		return nil, errors.New("no " + SnapshotsDirKey + " param found")
	}

	localSnapshotRepo, err := localsnap.NewLocalSnapshotRepo(typename, snapshotsDir)
	if err != nil {
		return nil, err
	}

	kubeconfgPathObj := params["kubeconfig"]
	kubeconfigPath := ""
	if kubeconfgPathObj != nil {
		kubeconfigPath = kubeconfgPathObj.(string)
	}
	config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	returnTypeManager := KubernetesNamespaceProtectedEntityTypeManager{
		clientset: clientset,
		logger:    logger,
		s3Config:  s3Config,
		internalRepo: localSnapshotRepo,
	}
	return &returnTypeManager, nil
}

func (this *KubernetesNamespaceProtectedEntityTypeManager) GetTypeName() string {
	return typename
}

func (this *KubernetesNamespaceProtectedEntityTypeManager) GetProtectedEntity(ctx context.Context, id astrolabe.ProtectedEntityID) (
	astrolabe.ProtectedEntity, error) {
	if (id.HasSnapshot()) {
		peinfo, err := this.internalRepo.GetPEInfoForID(ctx, id)
		if err != nil {
			return nil, errors.Wrap(err, "could not get peinfo")
		}
		return NewKubernetesNamespaceProtectedEntity(this, id, peinfo.GetName())
	} else {
		namespace, err := this.getNamespaceForPEID(ctx, id)
		if err != nil {
			return nil, errors.Wrap(err, "could not get namespace for id")
		}
		return NewKubernetesNamespaceProtectedEntity(this, id, namespace.Name)
	}
}

func (this *KubernetesNamespaceProtectedEntityTypeManager) getNamespaceForPEID(ctx context.Context, id astrolabe.ProtectedEntityID) (*v1.Namespace, error){
	namespaces, err := this.clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "Could not retrieve namespaces")
	}
	for _, curNamespace := range namespaces.Items {
		if string(curNamespace.UID) == id.GetID() {
			return &curNamespace, nil
		}
	}
	return nil, errors.New("Not found")
}
func (this *KubernetesNamespaceProtectedEntityTypeManager) GetProtectedEntities(ctx context.Context) ([]astrolabe.ProtectedEntityID, error) {
	namespaceList, err := this.clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return []astrolabe.ProtectedEntityID{}, err
	}
	var returnList []astrolabe.ProtectedEntityID
	for _, namespace := range namespaceList.Items {
		returnList = append(returnList, astrolabe.NewProtectedEntityID(this.GetTypeName(),
			string(namespace.UID)))
	}

	return returnList, nil
}

func (this *KubernetesNamespaceProtectedEntityTypeManager) Copy(ctx context.Context, pe astrolabe.ProtectedEntity, options astrolabe.CopyCreateOptions) (astrolabe.ProtectedEntity, error) {
	return nil, nil
}

func (this *KubernetesNamespaceProtectedEntityTypeManager) CopyFromInfo(ctx context.Context, pe astrolabe.ProtectedEntityInfo, options astrolabe.CopyCreateOptions) (astrolabe.ProtectedEntity, error) {
	return nil, nil
}

func (this *KubernetesNamespaceProtectedEntityTypeManager) Delete(ctx context.Context, id astrolabe.ProtectedEntityID) error {
	return errors.New("Not implemented")
}