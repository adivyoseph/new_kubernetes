/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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

// Package workflow provides Registry interface and it's RESTStorage
// implementation for storing Workflow api objects.
package etcd

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/registry/cachesize"
	"k8s.io/kubernetes/pkg/registry/generic"
	etcdgeneric "k8s.io/kubernetes/pkg/registry/generic/etcd"
	"k8s.io/kubernetes/pkg/registry/workflow"
	"k8s.io/kubernetes/pkg/runtime"
)

// rest implements a RESTStorage for workflows against etcd
type REST struct {
	*etcdgeneric.Etcd
}

// NewREST returns a RESTStorage object that will work against Workflows.
func NewREST(opts generic.RESTOptions) (*REST, *StatusREST) {
	prefix := "/workflows"

	newListFunc := func() runtime.Object { return &extensions.WorkflowList{} }
	storageInterface := opts.Decorator(
		opts.Storage, cachesize.GetWatchCacheSizeByResource(cachesize.Workflows), &extensions.Workflow{}, prefix, workflow.Strategy, newListFunc)

	store := &etcdgeneric.Etcd{
		NewFunc: func() runtime.Object { return &extensions.Workflow{} },

		// NewListFunc returns an object capable of storing results of an etcd list.
		NewListFunc: newListFunc,
		// Produces a path that etcd understands, to the root of the resource
		// by combining the namespace in the context with the given prefix
		KeyRootFunc: func(ctx api.Context) string {
			return etcdgeneric.NamespaceKeyRootFunc(ctx, prefix)
		},
		// Produces a path that etcd understands, to the resource by combining
		// the namespace in the context with the given prefix
		KeyFunc: func(ctx api.Context, name string) (string, error) {
			return etcdgeneric.NamespaceKeyFunc(ctx, prefix, name)
		},
		// Retrieve the name field of a workflow
		ObjectNameFunc: func(obj runtime.Object) (string, error) {
			return obj.(*extensions.Workflow).Name, nil
		},
		// Used to match objects based on labels/fields for list and watch
		PredicateFunc: func(label labels.Selector, field fields.Selector) generic.Matcher {
			return workflow.MatchWorkflow(label, field)
		},
		QualifiedResource: extensions.Resource("workflows"),

		// Used to validate workflow creation
		CreateStrategy: workflow.Strategy,

		// Used to validate workflow updates
		UpdateStrategy: workflow.Strategy,

		DeleteStrategy: workflow.Strategy,

		Storage: storageInterface,
	}

	statusStore := *store
	statusStore.UpdateStrategy = workflow.StatusStrategy

	return &REST{store}, &StatusREST{store: &statusStore}
}

// StatusREST implements the REST endpoint for changing the status of a resourcequota.
type StatusREST struct {
	store *etcdgeneric.Etcd
}

func (r *StatusREST) New() runtime.Object {
	return &extensions.Workflow{}
}

// Update alters the status subset of an object.
func (r *StatusREST) Update(ctx api.Context, obj runtime.Object) (runtime.Object, bool, error) {
	return r.store.Update(ctx, obj)
}
