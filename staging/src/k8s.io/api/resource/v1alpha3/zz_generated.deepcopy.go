//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha3

import (
	v1 "k8s.io/api/core/v1"
	resource "k8s.io/apimachinery/pkg/api/resource"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AllocationResult) DeepCopyInto(out *AllocationResult) {
	*out = *in
	in.Devices.DeepCopyInto(&out.Devices)
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = new(v1.NodeSelector)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AllocationResult.
func (in *AllocationResult) DeepCopy() *AllocationResult {
	if in == nil {
		return nil
	}
	out := new(AllocationResult)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CELDeviceSelector) DeepCopyInto(out *CELDeviceSelector) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CELDeviceSelector.
func (in *CELDeviceSelector) DeepCopy() *CELDeviceSelector {
	if in == nil {
		return nil
	}
	out := new(CELDeviceSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClassConfiguration) DeepCopyInto(out *ClassConfiguration) {
	*out = *in
	in.DeviceConfiguration.DeepCopyInto(&out.DeviceConfiguration)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClassConfiguration.
func (in *ClassConfiguration) DeepCopy() *ClassConfiguration {
	if in == nil {
		return nil
	}
	out := new(ClassConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Device) DeepCopyInto(out *Device) {
	*out = *in
	if in.Attributes != nil {
		in, out := &in.Attributes, &out.Attributes
		*out = make(map[QualifiedName]DeviceAttribute, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	if in.Capacity != nil {
		in, out := &in.Capacity, &out.Capacity
		*out = make(map[QualifiedName]resource.Quantity, len(*in))
		for key, val := range *in {
			(*out)[key] = val.DeepCopy()
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Device.
func (in *Device) DeepCopy() *Device {
	if in == nil {
		return nil
	}
	out := new(Device)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeviceAllocationConfiguration) DeepCopyInto(out *DeviceAllocationConfiguration) {
	*out = *in
	if in.Requests != nil {
		in, out := &in.Requests, &out.Requests
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.DeviceConfiguration.DeepCopyInto(&out.DeviceConfiguration)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeviceAllocationConfiguration.
func (in *DeviceAllocationConfiguration) DeepCopy() *DeviceAllocationConfiguration {
	if in == nil {
		return nil
	}
	out := new(DeviceAllocationConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeviceAllocationResult) DeepCopyInto(out *DeviceAllocationResult) {
	*out = *in
	if in.Results != nil {
		in, out := &in.Results, &out.Results
		*out = make([]DeviceRequestAllocationResult, len(*in))
		copy(*out, *in)
	}
	if in.Config != nil {
		in, out := &in.Config, &out.Config
		*out = make([]DeviceAllocationConfiguration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeviceAllocationResult.
func (in *DeviceAllocationResult) DeepCopy() *DeviceAllocationResult {
	if in == nil {
		return nil
	}
	out := new(DeviceAllocationResult)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeviceAttribute) DeepCopyInto(out *DeviceAttribute) {
	*out = *in
	if in.IntValue != nil {
		in, out := &in.IntValue, &out.IntValue
		*out = new(int64)
		**out = **in
	}
	if in.BoolValue != nil {
		in, out := &in.BoolValue, &out.BoolValue
		*out = new(bool)
		**out = **in
	}
	if in.StringValue != nil {
		in, out := &in.StringValue, &out.StringValue
		*out = new(string)
		**out = **in
	}
	if in.VersionValue != nil {
		in, out := &in.VersionValue, &out.VersionValue
		*out = new(string)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeviceAttribute.
func (in *DeviceAttribute) DeepCopy() *DeviceAttribute {
	if in == nil {
		return nil
	}
	out := new(DeviceAttribute)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeviceClaim) DeepCopyInto(out *DeviceClaim) {
	*out = *in
	if in.Requests != nil {
		in, out := &in.Requests, &out.Requests
		*out = make([]DeviceRequest, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Constraints != nil {
		in, out := &in.Constraints, &out.Constraints
		*out = make([]DeviceConstraint, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Config != nil {
		in, out := &in.Config, &out.Config
		*out = make([]DeviceClaimConfiguration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeviceClaim.
func (in *DeviceClaim) DeepCopy() *DeviceClaim {
	if in == nil {
		return nil
	}
	out := new(DeviceClaim)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeviceClaimConfiguration) DeepCopyInto(out *DeviceClaimConfiguration) {
	*out = *in
	if in.Requests != nil {
		in, out := &in.Requests, &out.Requests
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.DeviceConfiguration.DeepCopyInto(&out.DeviceConfiguration)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeviceClaimConfiguration.
func (in *DeviceClaimConfiguration) DeepCopy() *DeviceClaimConfiguration {
	if in == nil {
		return nil
	}
	out := new(DeviceClaimConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeviceClass) DeepCopyInto(out *DeviceClass) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeviceClass.
func (in *DeviceClass) DeepCopy() *DeviceClass {
	if in == nil {
		return nil
	}
	out := new(DeviceClass)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DeviceClass) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeviceClassList) DeepCopyInto(out *DeviceClassList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DeviceClass, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeviceClassList.
func (in *DeviceClassList) DeepCopy() *DeviceClassList {
	if in == nil {
		return nil
	}
	out := new(DeviceClassList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DeviceClassList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeviceClassSpec) DeepCopyInto(out *DeviceClassSpec) {
	*out = *in
	if in.Selectors != nil {
		in, out := &in.Selectors, &out.Selectors
		*out = make([]DeviceSelector, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Config != nil {
		in, out := &in.Config, &out.Config
		*out = make([]ClassConfiguration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.SuitableNodes != nil {
		in, out := &in.SuitableNodes, &out.SuitableNodes
		*out = new(v1.NodeSelector)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeviceClassSpec.
func (in *DeviceClassSpec) DeepCopy() *DeviceClassSpec {
	if in == nil {
		return nil
	}
	out := new(DeviceClassSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeviceConfiguration) DeepCopyInto(out *DeviceConfiguration) {
	*out = *in
	if in.Opaque != nil {
		in, out := &in.Opaque, &out.Opaque
		*out = new(OpaqueDeviceConfiguration)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeviceConfiguration.
func (in *DeviceConfiguration) DeepCopy() *DeviceConfiguration {
	if in == nil {
		return nil
	}
	out := new(DeviceConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeviceConstraint) DeepCopyInto(out *DeviceConstraint) {
	*out = *in
	if in.Requests != nil {
		in, out := &in.Requests, &out.Requests
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.MatchAttribute != nil {
		in, out := &in.MatchAttribute, &out.MatchAttribute
		*out = new(FullyQualifiedName)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeviceConstraint.
func (in *DeviceConstraint) DeepCopy() *DeviceConstraint {
	if in == nil {
		return nil
	}
	out := new(DeviceConstraint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeviceRequest) DeepCopyInto(out *DeviceRequest) {
	*out = *in
	if in.DeviceRequestDetails != nil {
		in, out := &in.DeviceRequestDetails, &out.DeviceRequestDetails
		*out = new(DeviceRequestDetails)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeviceRequest.
func (in *DeviceRequest) DeepCopy() *DeviceRequest {
	if in == nil {
		return nil
	}
	out := new(DeviceRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeviceRequestAllocationResult) DeepCopyInto(out *DeviceRequestAllocationResult) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeviceRequestAllocationResult.
func (in *DeviceRequestAllocationResult) DeepCopy() *DeviceRequestAllocationResult {
	if in == nil {
		return nil
	}
	out := new(DeviceRequestAllocationResult)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeviceRequestDetails) DeepCopyInto(out *DeviceRequestDetails) {
	*out = *in
	if in.Selectors != nil {
		in, out := &in.Selectors, &out.Selectors
		*out = make([]DeviceSelector, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeviceRequestDetails.
func (in *DeviceRequestDetails) DeepCopy() *DeviceRequestDetails {
	if in == nil {
		return nil
	}
	out := new(DeviceRequestDetails)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeviceSelector) DeepCopyInto(out *DeviceSelector) {
	*out = *in
	if in.CEL != nil {
		in, out := &in.CEL, &out.CEL
		*out = new(CELDeviceSelector)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeviceSelector.
func (in *DeviceSelector) DeepCopy() *DeviceSelector {
	if in == nil {
		return nil
	}
	out := new(DeviceSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OpaqueDeviceConfiguration) DeepCopyInto(out *OpaqueDeviceConfiguration) {
	*out = *in
	in.Parameters.DeepCopyInto(&out.Parameters)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OpaqueDeviceConfiguration.
func (in *OpaqueDeviceConfiguration) DeepCopy() *OpaqueDeviceConfiguration {
	if in == nil {
		return nil
	}
	out := new(OpaqueDeviceConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodSchedulingContext) DeepCopyInto(out *PodSchedulingContext) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodSchedulingContext.
func (in *PodSchedulingContext) DeepCopy() *PodSchedulingContext {
	if in == nil {
		return nil
	}
	out := new(PodSchedulingContext)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PodSchedulingContext) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodSchedulingContextList) DeepCopyInto(out *PodSchedulingContextList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PodSchedulingContext, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodSchedulingContextList.
func (in *PodSchedulingContextList) DeepCopy() *PodSchedulingContextList {
	if in == nil {
		return nil
	}
	out := new(PodSchedulingContextList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PodSchedulingContextList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodSchedulingContextSpec) DeepCopyInto(out *PodSchedulingContextSpec) {
	*out = *in
	if in.PotentialNodes != nil {
		in, out := &in.PotentialNodes, &out.PotentialNodes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodSchedulingContextSpec.
func (in *PodSchedulingContextSpec) DeepCopy() *PodSchedulingContextSpec {
	if in == nil {
		return nil
	}
	out := new(PodSchedulingContextSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodSchedulingContextStatus) DeepCopyInto(out *PodSchedulingContextStatus) {
	*out = *in
	if in.ResourceClaims != nil {
		in, out := &in.ResourceClaims, &out.ResourceClaims
		*out = make([]ResourceClaimSchedulingStatus, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodSchedulingContextStatus.
func (in *PodSchedulingContextStatus) DeepCopy() *PodSchedulingContextStatus {
	if in == nil {
		return nil
	}
	out := new(PodSchedulingContextStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceClaim) DeepCopyInto(out *ResourceClaim) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceClaim.
func (in *ResourceClaim) DeepCopy() *ResourceClaim {
	if in == nil {
		return nil
	}
	out := new(ResourceClaim)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceClaim) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceClaimConsumerReference) DeepCopyInto(out *ResourceClaimConsumerReference) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceClaimConsumerReference.
func (in *ResourceClaimConsumerReference) DeepCopy() *ResourceClaimConsumerReference {
	if in == nil {
		return nil
	}
	out := new(ResourceClaimConsumerReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceClaimList) DeepCopyInto(out *ResourceClaimList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ResourceClaim, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceClaimList.
func (in *ResourceClaimList) DeepCopy() *ResourceClaimList {
	if in == nil {
		return nil
	}
	out := new(ResourceClaimList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceClaimList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceClaimSchedulingStatus) DeepCopyInto(out *ResourceClaimSchedulingStatus) {
	*out = *in
	if in.UnsuitableNodes != nil {
		in, out := &in.UnsuitableNodes, &out.UnsuitableNodes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceClaimSchedulingStatus.
func (in *ResourceClaimSchedulingStatus) DeepCopy() *ResourceClaimSchedulingStatus {
	if in == nil {
		return nil
	}
	out := new(ResourceClaimSchedulingStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceClaimSpec) DeepCopyInto(out *ResourceClaimSpec) {
	*out = *in
	in.Devices.DeepCopyInto(&out.Devices)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceClaimSpec.
func (in *ResourceClaimSpec) DeepCopy() *ResourceClaimSpec {
	if in == nil {
		return nil
	}
	out := new(ResourceClaimSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceClaimStatus) DeepCopyInto(out *ResourceClaimStatus) {
	*out = *in
	if in.Allocation != nil {
		in, out := &in.Allocation, &out.Allocation
		*out = new(AllocationResult)
		(*in).DeepCopyInto(*out)
	}
	if in.ReservedFor != nil {
		in, out := &in.ReservedFor, &out.ReservedFor
		*out = make([]ResourceClaimConsumerReference, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceClaimStatus.
func (in *ResourceClaimStatus) DeepCopy() *ResourceClaimStatus {
	if in == nil {
		return nil
	}
	out := new(ResourceClaimStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceClaimTemplate) DeepCopyInto(out *ResourceClaimTemplate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceClaimTemplate.
func (in *ResourceClaimTemplate) DeepCopy() *ResourceClaimTemplate {
	if in == nil {
		return nil
	}
	out := new(ResourceClaimTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceClaimTemplate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceClaimTemplateList) DeepCopyInto(out *ResourceClaimTemplateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ResourceClaimTemplate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceClaimTemplateList.
func (in *ResourceClaimTemplateList) DeepCopy() *ResourceClaimTemplateList {
	if in == nil {
		return nil
	}
	out := new(ResourceClaimTemplateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceClaimTemplateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceClaimTemplateSpec) DeepCopyInto(out *ResourceClaimTemplateSpec) {
	*out = *in
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceClaimTemplateSpec.
func (in *ResourceClaimTemplateSpec) DeepCopy() *ResourceClaimTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(ResourceClaimTemplateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourcePool) DeepCopyInto(out *ResourcePool) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourcePool.
func (in *ResourcePool) DeepCopy() *ResourcePool {
	if in == nil {
		return nil
	}
	out := new(ResourcePool)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceSlice) DeepCopyInto(out *ResourceSlice) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceSlice.
func (in *ResourceSlice) DeepCopy() *ResourceSlice {
	if in == nil {
		return nil
	}
	out := new(ResourceSlice)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceSlice) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceSliceList) DeepCopyInto(out *ResourceSliceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ResourceSlice, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceSliceList.
func (in *ResourceSliceList) DeepCopy() *ResourceSliceList {
	if in == nil {
		return nil
	}
	out := new(ResourceSliceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceSliceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceSliceSpec) DeepCopyInto(out *ResourceSliceSpec) {
	*out = *in
	out.Pool = in.Pool
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = new(v1.NodeSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.Devices != nil {
		in, out := &in.Devices, &out.Devices
		*out = make([]Device, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceSliceSpec.
func (in *ResourceSliceSpec) DeepCopy() *ResourceSliceSpec {
	if in == nil {
		return nil
	}
	out := new(ResourceSliceSpec)
	in.DeepCopyInto(out)
	return out
}
