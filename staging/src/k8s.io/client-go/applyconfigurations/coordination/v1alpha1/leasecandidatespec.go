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

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LeaseCandidateSpecApplyConfiguration represents a declarative configuration of the LeaseCandidateSpec type for use
// with apply.
type LeaseCandidateSpecApplyConfiguration struct {
	BinaryVersion        *string       `json:"binaryVersion,omitempty"`
	CompatibilityVersion *string       `json:"compatibilityVersion,omitempty"`
	TargetLease          *string       `json:"targetLease,omitempty"`
	LeaseDurationSeconds *int32        `json:"leaseDurationSeconds,omitempty"`
	RenewTime            *v1.MicroTime `json:"renewTime,omitempty"`
}

// LeaseCandidateSpecApplyConfiguration constructs a declarative configuration of the LeaseCandidateSpec type for use with
// apply.
func LeaseCandidateSpec() *LeaseCandidateSpecApplyConfiguration {
	return &LeaseCandidateSpecApplyConfiguration{}
}

// WithBinaryVersion sets the BinaryVersion field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the BinaryVersion field is set to the value of the last call.
func (b *LeaseCandidateSpecApplyConfiguration) WithBinaryVersion(value string) *LeaseCandidateSpecApplyConfiguration {
	b.BinaryVersion = &value
	return b
}

// WithCompatibilityVersion sets the CompatibilityVersion field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the CompatibilityVersion field is set to the value of the last call.
func (b *LeaseCandidateSpecApplyConfiguration) WithCompatibilityVersion(value string) *LeaseCandidateSpecApplyConfiguration {
	b.CompatibilityVersion = &value
	return b
}

// WithTargetLease sets the TargetLease field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the TargetLease field is set to the value of the last call.
func (b *LeaseCandidateSpecApplyConfiguration) WithTargetLease(value string) *LeaseCandidateSpecApplyConfiguration {
	b.TargetLease = &value
	return b
}

// WithLeaseDurationSeconds sets the LeaseDurationSeconds field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the LeaseDurationSeconds field is set to the value of the last call.
func (b *LeaseCandidateSpecApplyConfiguration) WithLeaseDurationSeconds(value int32) *LeaseCandidateSpecApplyConfiguration {
	b.LeaseDurationSeconds = &value
	return b
}

// WithRenewTime sets the RenewTime field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the RenewTime field is set to the value of the last call.
func (b *LeaseCandidateSpecApplyConfiguration) WithRenewTime(value v1.MicroTime) *LeaseCandidateSpecApplyConfiguration {
	b.RenewTime = &value
	return b
}
