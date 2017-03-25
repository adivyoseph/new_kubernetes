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

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto will perform a deep copy of the receiver, writing to out. in must be non-nil.
func (in *API) DeepCopyInto(out *API) {
	*out = *in
	return
}

// DeepCopy will perform a deep copy of the receiver, creating a new API.
func (x *API) DeepCopy() *API {
	if x == nil {
		return nil
	}
	out := new(API)
	x.DeepCopyInto(out)
	return out
}

// DeepCopyInto will perform a deep copy of the receiver, writing to out. in must be non-nil.
func (in *Etcd) DeepCopyInto(out *Etcd) {
	*out = *in
	if in.Endpoints != nil {
		in, out := &in.Endpoints, &out.Endpoints
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy will perform a deep copy of the receiver, creating a new Etcd.
func (x *Etcd) DeepCopy() *Etcd {
	if x == nil {
		return nil
	}
	out := new(Etcd)
	x.DeepCopyInto(out)
	return out
}

// DeepCopyInto will perform a deep copy of the receiver, writing to out. in must be non-nil.
func (in *MasterConfiguration) DeepCopyInto(out *MasterConfiguration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.API = in.API
	in.Etcd.DeepCopyInto(&out.Etcd)
	out.Networking = in.Networking
	if in.APIServerExtraArgs != nil {
		in, out := &in.APIServerExtraArgs, &out.APIServerExtraArgs
		*out = make(map[string]string)
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.ControllerManagerExtraArgs != nil {
		in, out := &in.ControllerManagerExtraArgs, &out.ControllerManagerExtraArgs
		*out = make(map[string]string)
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.SchedulerExtraArgs != nil {
		in, out := &in.SchedulerExtraArgs, &out.SchedulerExtraArgs
		*out = make(map[string]string)
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.APIServerCertSANs != nil {
		in, out := &in.APIServerCertSANs, &out.APIServerCertSANs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy will perform a deep copy of the receiver, creating a new MasterConfiguration.
func (x *MasterConfiguration) DeepCopy() *MasterConfiguration {
	if x == nil {
		return nil
	}
	out := new(MasterConfiguration)
	x.DeepCopyInto(out)
	return out
}

// DeepCopyObject will perform a deep copy of the receiver, creating a new object.
func (x *MasterConfiguration) DeepCopyObject() runtime.Object {
	if c := x.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto will perform a deep copy of the receiver, writing to out. in must be non-nil.
func (in *Networking) DeepCopyInto(out *Networking) {
	*out = *in
	return
}

// DeepCopy will perform a deep copy of the receiver, creating a new Networking.
func (x *Networking) DeepCopy() *Networking {
	if x == nil {
		return nil
	}
	out := new(Networking)
	x.DeepCopyInto(out)
	return out
}

// DeepCopyInto will perform a deep copy of the receiver, writing to out. in must be non-nil.
func (in *NodeConfiguration) DeepCopyInto(out *NodeConfiguration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	if in.DiscoveryTokenAPIServers != nil {
		in, out := &in.DiscoveryTokenAPIServers, &out.DiscoveryTokenAPIServers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy will perform a deep copy of the receiver, creating a new NodeConfiguration.
func (x *NodeConfiguration) DeepCopy() *NodeConfiguration {
	if x == nil {
		return nil
	}
	out := new(NodeConfiguration)
	x.DeepCopyInto(out)
	return out
}

// DeepCopyObject will perform a deep copy of the receiver, creating a new object.
func (x *NodeConfiguration) DeepCopyObject() runtime.Object {
	if c := x.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto will perform a deep copy of the receiver, writing to out. in must be non-nil.
func (in *TokenDiscovery) DeepCopyInto(out *TokenDiscovery) {
	*out = *in
	if in.Addresses != nil {
		in, out := &in.Addresses, &out.Addresses
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy will perform a deep copy of the receiver, creating a new TokenDiscovery.
func (x *TokenDiscovery) DeepCopy() *TokenDiscovery {
	if x == nil {
		return nil
	}
	out := new(TokenDiscovery)
	x.DeepCopyInto(out)
	return out
}
