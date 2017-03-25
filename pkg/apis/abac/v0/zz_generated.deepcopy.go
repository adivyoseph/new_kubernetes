// +build !ignore_autogenerated

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

// This file was autogenerated by deepcopy-gen. Do not edit it manually!

package v0

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto will perform a deep copy of the receiver, writing to out. in must be non-nil.
func (in *Policy) DeepCopyInto(out *Policy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	return
}

// DeepCopy will perform a deep copy of the receiver, creating a new Policy.
func (x *Policy) DeepCopy() *Policy {
	if x == nil {
		return nil
	}
	out := new(Policy)
	x.DeepCopyInto(out)
	return out
}

// DeepCopyObject will perform a deep copy of the receiver, creating a new object.
func (x *Policy) DeepCopyObject() runtime.Object {
	if c := x.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}
