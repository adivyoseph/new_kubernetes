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

// Code generated by conversion-gen. DO NOT EDIT.

package v1alpha1

import (
	unsafe "unsafe"

	v1alpha1 "k8s.io/api/coordination/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	coordination "k8s.io/kubernetes/pkg/apis/coordination"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1alpha1.LeaseCandidate)(nil), (*coordination.LeaseCandidate)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_LeaseCandidate_To_coordination_LeaseCandidate(a.(*v1alpha1.LeaseCandidate), b.(*coordination.LeaseCandidate), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*coordination.LeaseCandidate)(nil), (*v1alpha1.LeaseCandidate)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_coordination_LeaseCandidate_To_v1alpha1_LeaseCandidate(a.(*coordination.LeaseCandidate), b.(*v1alpha1.LeaseCandidate), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.LeaseCandidateList)(nil), (*coordination.LeaseCandidateList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_LeaseCandidateList_To_coordination_LeaseCandidateList(a.(*v1alpha1.LeaseCandidateList), b.(*coordination.LeaseCandidateList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*coordination.LeaseCandidateList)(nil), (*v1alpha1.LeaseCandidateList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_coordination_LeaseCandidateList_To_v1alpha1_LeaseCandidateList(a.(*coordination.LeaseCandidateList), b.(*v1alpha1.LeaseCandidateList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.LeaseCandidateSpec)(nil), (*coordination.LeaseCandidateSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_LeaseCandidateSpec_To_coordination_LeaseCandidateSpec(a.(*v1alpha1.LeaseCandidateSpec), b.(*coordination.LeaseCandidateSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*coordination.LeaseCandidateSpec)(nil), (*v1alpha1.LeaseCandidateSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_coordination_LeaseCandidateSpec_To_v1alpha1_LeaseCandidateSpec(a.(*coordination.LeaseCandidateSpec), b.(*v1alpha1.LeaseCandidateSpec), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_LeaseCandidate_To_coordination_LeaseCandidate(in *v1alpha1.LeaseCandidate, out *coordination.LeaseCandidate, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_LeaseCandidateSpec_To_coordination_LeaseCandidateSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_LeaseCandidate_To_coordination_LeaseCandidate is an autogenerated conversion function.
func Convert_v1alpha1_LeaseCandidate_To_coordination_LeaseCandidate(in *v1alpha1.LeaseCandidate, out *coordination.LeaseCandidate, s conversion.Scope) error {
	return autoConvert_v1alpha1_LeaseCandidate_To_coordination_LeaseCandidate(in, out, s)
}

func autoConvert_coordination_LeaseCandidate_To_v1alpha1_LeaseCandidate(in *coordination.LeaseCandidate, out *v1alpha1.LeaseCandidate, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_coordination_LeaseCandidateSpec_To_v1alpha1_LeaseCandidateSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

// Convert_coordination_LeaseCandidate_To_v1alpha1_LeaseCandidate is an autogenerated conversion function.
func Convert_coordination_LeaseCandidate_To_v1alpha1_LeaseCandidate(in *coordination.LeaseCandidate, out *v1alpha1.LeaseCandidate, s conversion.Scope) error {
	return autoConvert_coordination_LeaseCandidate_To_v1alpha1_LeaseCandidate(in, out, s)
}

func autoConvert_v1alpha1_LeaseCandidateList_To_coordination_LeaseCandidateList(in *v1alpha1.LeaseCandidateList, out *coordination.LeaseCandidateList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]coordination.LeaseCandidate)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1alpha1_LeaseCandidateList_To_coordination_LeaseCandidateList is an autogenerated conversion function.
func Convert_v1alpha1_LeaseCandidateList_To_coordination_LeaseCandidateList(in *v1alpha1.LeaseCandidateList, out *coordination.LeaseCandidateList, s conversion.Scope) error {
	return autoConvert_v1alpha1_LeaseCandidateList_To_coordination_LeaseCandidateList(in, out, s)
}

func autoConvert_coordination_LeaseCandidateList_To_v1alpha1_LeaseCandidateList(in *coordination.LeaseCandidateList, out *v1alpha1.LeaseCandidateList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]v1alpha1.LeaseCandidate)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_coordination_LeaseCandidateList_To_v1alpha1_LeaseCandidateList is an autogenerated conversion function.
func Convert_coordination_LeaseCandidateList_To_v1alpha1_LeaseCandidateList(in *coordination.LeaseCandidateList, out *v1alpha1.LeaseCandidateList, s conversion.Scope) error {
	return autoConvert_coordination_LeaseCandidateList_To_v1alpha1_LeaseCandidateList(in, out, s)
}

func autoConvert_v1alpha1_LeaseCandidateSpec_To_coordination_LeaseCandidateSpec(in *v1alpha1.LeaseCandidateSpec, out *coordination.LeaseCandidateSpec, s conversion.Scope) error {
	out.BinaryVersion = in.BinaryVersion
	out.CompatibilityVersion = in.CompatibilityVersion
	out.TargetLease = in.TargetLease
	out.LeaseDurationSeconds = (*int32)(unsafe.Pointer(in.LeaseDurationSeconds))
	out.RenewTime = (*v1.MicroTime)(unsafe.Pointer(in.RenewTime))
	return nil
}

// Convert_v1alpha1_LeaseCandidateSpec_To_coordination_LeaseCandidateSpec is an autogenerated conversion function.
func Convert_v1alpha1_LeaseCandidateSpec_To_coordination_LeaseCandidateSpec(in *v1alpha1.LeaseCandidateSpec, out *coordination.LeaseCandidateSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_LeaseCandidateSpec_To_coordination_LeaseCandidateSpec(in, out, s)
}

func autoConvert_coordination_LeaseCandidateSpec_To_v1alpha1_LeaseCandidateSpec(in *coordination.LeaseCandidateSpec, out *v1alpha1.LeaseCandidateSpec, s conversion.Scope) error {
	out.BinaryVersion = in.BinaryVersion
	out.CompatibilityVersion = in.CompatibilityVersion
	out.TargetLease = in.TargetLease
	out.LeaseDurationSeconds = (*int32)(unsafe.Pointer(in.LeaseDurationSeconds))
	out.RenewTime = (*v1.MicroTime)(unsafe.Pointer(in.RenewTime))
	return nil
}

// Convert_coordination_LeaseCandidateSpec_To_v1alpha1_LeaseCandidateSpec is an autogenerated conversion function.
func Convert_coordination_LeaseCandidateSpec_To_v1alpha1_LeaseCandidateSpec(in *coordination.LeaseCandidateSpec, out *v1alpha1.LeaseCandidateSpec, s conversion.Scope) error {
	return autoConvert_coordination_LeaseCandidateSpec_To_v1alpha1_LeaseCandidateSpec(in, out, s)
}
