/*
Copyright 2014 The Kubernetes Authors All rights reserved.

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

package rest_test

import (
	"testing"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/api/rest"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/util/fielderrors"
)

type fakeStrategy struct {
	runtime.ObjectTyper
	api.NameGenerator
}

var strategy = fakeStrategy{api.Scheme, api.SimpleNameGenerator}

func (fakeStrategy) NamespaceScoped() bool {
	return true
}

func (fakeStrategy) PrepareForCreate(obj runtime.Object) {
}

func (fakeStrategy) Validate(ctx api.Context, obj runtime.Object) fielderrors.ValidationErrorList {
	return nil
}

func TestCheckGeneratedNameError(t *testing.T) {
	expect := errors.NewNotFound("foo", "bar")
	if err := rest.CheckGeneratedNameError(strategy, expect, &api.Pod{}); err != expect {
		t.Errorf("NotFoundError should be ignored: %v", err)
	}

	expect = errors.NewAlreadyExists("foo", "bar")
	if err := rest.CheckGeneratedNameError(strategy, expect, &api.Pod{}); err != expect {
		t.Errorf("AlreadyExists should be returned when no GenerateName field: %v", err)
	}

	expect = errors.NewAlreadyExists("foo", "bar")
	if err := rest.CheckGeneratedNameError(strategy, expect, &api.Pod{ObjectMeta: api.ObjectMeta{GenerateName: "foo"}}); err == nil || !errors.IsServerTimeout(err) {
		t.Errorf("expected try again later error: %v", err)
	}
}
