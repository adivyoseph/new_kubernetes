/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fixture

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pkgruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	federationclientset "k8s.io/kubernetes/federation/client/clientset_generated/federation_clientset"
	"k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
)

type ResourceFixture interface {
	Setup(t *testing.T, testClient, controllerClient federationclientset.Interface)
	Teardown(t *testing.T)

	Kind() string
	GetAdapter() ResourceAdapter
	New(namespace string) pkgruntime.Object
	Mutate(obj pkgruntime.Object) pkgruntime.Object
}

// TODO reuse resource adapters defined for use with a generic controller as per
// https://github.com/kubernetes/kubernetes/pull/41050
type ResourceAdapter interface {
	Kind() string
	Equivalent(obj1, obj2 pkgruntime.Object) bool
	GetNamespacedName(obj pkgruntime.Object) types.NamespacedName

	FedCreate(obj pkgruntime.Object) (pkgruntime.Object, error)
	FedGet(nsName types.NamespacedName) (pkgruntime.Object, error)
	FedUpdate(obj pkgruntime.Object) (pkgruntime.Object, error)
	FedDelete(nsName types.NamespacedName, options *metav1.DeleteOptions) error

	Get(client clientset.Interface, nsName types.NamespacedName) (pkgruntime.Object, error)
}
