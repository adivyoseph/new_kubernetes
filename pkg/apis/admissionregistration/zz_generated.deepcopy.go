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

package admissionregistration

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AuditAnnotation) DeepCopyInto(out *AuditAnnotation) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AuditAnnotation.
func (in *AuditAnnotation) DeepCopy() *AuditAnnotation {
	if in == nil {
		return nil
	}
	out := new(AuditAnnotation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExpressionWarning) DeepCopyInto(out *ExpressionWarning) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExpressionWarning.
func (in *ExpressionWarning) DeepCopy() *ExpressionWarning {
	if in == nil {
		return nil
	}
	out := new(ExpressionWarning)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MatchCondition) DeepCopyInto(out *MatchCondition) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MatchCondition.
func (in *MatchCondition) DeepCopy() *MatchCondition {
	if in == nil {
		return nil
	}
	out := new(MatchCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MatchResources) DeepCopyInto(out *MatchResources) {
	*out = *in
	if in.NamespaceSelector != nil {
		in, out := &in.NamespaceSelector, &out.NamespaceSelector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.ObjectSelector != nil {
		in, out := &in.ObjectSelector, &out.ObjectSelector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.ResourceRules != nil {
		in, out := &in.ResourceRules, &out.ResourceRules
		*out = make([]NamedRuleWithOperations, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ExcludeResourceRules != nil {
		in, out := &in.ExcludeResourceRules, &out.ExcludeResourceRules
		*out = make([]NamedRuleWithOperations, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.MatchPolicy != nil {
		in, out := &in.MatchPolicy, &out.MatchPolicy
		*out = new(MatchPolicyType)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MatchResources.
func (in *MatchResources) DeepCopy() *MatchResources {
	if in == nil {
		return nil
	}
	out := new(MatchResources)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MutatingWebhook) DeepCopyInto(out *MutatingWebhook) {
	*out = *in
	in.ClientConfig.DeepCopyInto(&out.ClientConfig)
	if in.Rules != nil {
		in, out := &in.Rules, &out.Rules
		*out = make([]RuleWithOperations, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.FailurePolicy != nil {
		in, out := &in.FailurePolicy, &out.FailurePolicy
		*out = new(FailurePolicyType)
		**out = **in
	}
	if in.MatchPolicy != nil {
		in, out := &in.MatchPolicy, &out.MatchPolicy
		*out = new(MatchPolicyType)
		**out = **in
	}
	if in.NamespaceSelector != nil {
		in, out := &in.NamespaceSelector, &out.NamespaceSelector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.ObjectSelector != nil {
		in, out := &in.ObjectSelector, &out.ObjectSelector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.SideEffects != nil {
		in, out := &in.SideEffects, &out.SideEffects
		*out = new(SideEffectClass)
		**out = **in
	}
	if in.TimeoutSeconds != nil {
		in, out := &in.TimeoutSeconds, &out.TimeoutSeconds
		*out = new(int32)
		**out = **in
	}
	if in.AdmissionReviewVersions != nil {
		in, out := &in.AdmissionReviewVersions, &out.AdmissionReviewVersions
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ReinvocationPolicy != nil {
		in, out := &in.ReinvocationPolicy, &out.ReinvocationPolicy
		*out = new(ReinvocationPolicyType)
		**out = **in
	}
	if in.MatchConditions != nil {
		in, out := &in.MatchConditions, &out.MatchConditions
		*out = make([]MatchCondition, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MutatingWebhook.
func (in *MutatingWebhook) DeepCopy() *MutatingWebhook {
	if in == nil {
		return nil
	}
	out := new(MutatingWebhook)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MutatingWebhookConfiguration) DeepCopyInto(out *MutatingWebhookConfiguration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Webhooks != nil {
		in, out := &in.Webhooks, &out.Webhooks
		*out = make([]MutatingWebhook, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MutatingWebhookConfiguration.
func (in *MutatingWebhookConfiguration) DeepCopy() *MutatingWebhookConfiguration {
	if in == nil {
		return nil
	}
	out := new(MutatingWebhookConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MutatingWebhookConfiguration) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MutatingWebhookConfigurationList) DeepCopyInto(out *MutatingWebhookConfigurationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MutatingWebhookConfiguration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MutatingWebhookConfigurationList.
func (in *MutatingWebhookConfigurationList) DeepCopy() *MutatingWebhookConfigurationList {
	if in == nil {
		return nil
	}
	out := new(MutatingWebhookConfigurationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MutatingWebhookConfigurationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamedRuleWithOperations) DeepCopyInto(out *NamedRuleWithOperations) {
	*out = *in
	if in.ResourceNames != nil {
		in, out := &in.ResourceNames, &out.ResourceNames
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.RuleWithOperations.DeepCopyInto(&out.RuleWithOperations)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamedRuleWithOperations.
func (in *NamedRuleWithOperations) DeepCopy() *NamedRuleWithOperations {
	if in == nil {
		return nil
	}
	out := new(NamedRuleWithOperations)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ParamKind) DeepCopyInto(out *ParamKind) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ParamKind.
func (in *ParamKind) DeepCopy() *ParamKind {
	if in == nil {
		return nil
	}
	out := new(ParamKind)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ParamRef) DeepCopyInto(out *ParamRef) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ParamRef.
func (in *ParamRef) DeepCopy() *ParamRef {
	if in == nil {
		return nil
	}
	out := new(ParamRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Rule) DeepCopyInto(out *Rule) {
	*out = *in
	if in.APIGroups != nil {
		in, out := &in.APIGroups, &out.APIGroups
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.APIVersions != nil {
		in, out := &in.APIVersions, &out.APIVersions
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Scope != nil {
		in, out := &in.Scope, &out.Scope
		*out = new(ScopeType)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Rule.
func (in *Rule) DeepCopy() *Rule {
	if in == nil {
		return nil
	}
	out := new(Rule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RuleWithOperations) DeepCopyInto(out *RuleWithOperations) {
	*out = *in
	if in.Operations != nil {
		in, out := &in.Operations, &out.Operations
		*out = make([]OperationType, len(*in))
		copy(*out, *in)
	}
	in.Rule.DeepCopyInto(&out.Rule)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RuleWithOperations.
func (in *RuleWithOperations) DeepCopy() *RuleWithOperations {
	if in == nil {
		return nil
	}
	out := new(RuleWithOperations)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceReference) DeepCopyInto(out *ServiceReference) {
	*out = *in
	if in.Path != nil {
		in, out := &in.Path, &out.Path
		*out = new(string)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceReference.
func (in *ServiceReference) DeepCopy() *ServiceReference {
	if in == nil {
		return nil
	}
	out := new(ServiceReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TypeChecking) DeepCopyInto(out *TypeChecking) {
	*out = *in
	if in.ExpressionWarnings != nil {
		in, out := &in.ExpressionWarnings, &out.ExpressionWarnings
		*out = make([]ExpressionWarning, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TypeChecking.
func (in *TypeChecking) DeepCopy() *TypeChecking {
	if in == nil {
		return nil
	}
	out := new(TypeChecking)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValidatingAdmissionPolicy) DeepCopyInto(out *ValidatingAdmissionPolicy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ValidatingAdmissionPolicy.
func (in *ValidatingAdmissionPolicy) DeepCopy() *ValidatingAdmissionPolicy {
	if in == nil {
		return nil
	}
	out := new(ValidatingAdmissionPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ValidatingAdmissionPolicy) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValidatingAdmissionPolicyBinding) DeepCopyInto(out *ValidatingAdmissionPolicyBinding) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ValidatingAdmissionPolicyBinding.
func (in *ValidatingAdmissionPolicyBinding) DeepCopy() *ValidatingAdmissionPolicyBinding {
	if in == nil {
		return nil
	}
	out := new(ValidatingAdmissionPolicyBinding)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ValidatingAdmissionPolicyBinding) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValidatingAdmissionPolicyBindingList) DeepCopyInto(out *ValidatingAdmissionPolicyBindingList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ValidatingAdmissionPolicyBinding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ValidatingAdmissionPolicyBindingList.
func (in *ValidatingAdmissionPolicyBindingList) DeepCopy() *ValidatingAdmissionPolicyBindingList {
	if in == nil {
		return nil
	}
	out := new(ValidatingAdmissionPolicyBindingList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ValidatingAdmissionPolicyBindingList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValidatingAdmissionPolicyBindingSpec) DeepCopyInto(out *ValidatingAdmissionPolicyBindingSpec) {
	*out = *in
	if in.ParamRef != nil {
		in, out := &in.ParamRef, &out.ParamRef
		*out = new(ParamRef)
		**out = **in
	}
	if in.MatchResources != nil {
		in, out := &in.MatchResources, &out.MatchResources
		*out = new(MatchResources)
		(*in).DeepCopyInto(*out)
	}
	if in.ValidationActions != nil {
		in, out := &in.ValidationActions, &out.ValidationActions
		*out = make([]ValidationAction, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ValidatingAdmissionPolicyBindingSpec.
func (in *ValidatingAdmissionPolicyBindingSpec) DeepCopy() *ValidatingAdmissionPolicyBindingSpec {
	if in == nil {
		return nil
	}
	out := new(ValidatingAdmissionPolicyBindingSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValidatingAdmissionPolicyList) DeepCopyInto(out *ValidatingAdmissionPolicyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ValidatingAdmissionPolicy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ValidatingAdmissionPolicyList.
func (in *ValidatingAdmissionPolicyList) DeepCopy() *ValidatingAdmissionPolicyList {
	if in == nil {
		return nil
	}
	out := new(ValidatingAdmissionPolicyList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ValidatingAdmissionPolicyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValidatingAdmissionPolicySpec) DeepCopyInto(out *ValidatingAdmissionPolicySpec) {
	*out = *in
	if in.ParamKind != nil {
		in, out := &in.ParamKind, &out.ParamKind
		*out = new(ParamKind)
		**out = **in
	}
	if in.MatchConstraints != nil {
		in, out := &in.MatchConstraints, &out.MatchConstraints
		*out = new(MatchResources)
		(*in).DeepCopyInto(*out)
	}
	if in.Validations != nil {
		in, out := &in.Validations, &out.Validations
		*out = make([]Validation, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.FailurePolicy != nil {
		in, out := &in.FailurePolicy, &out.FailurePolicy
		*out = new(FailurePolicyType)
		**out = **in
	}
	if in.AuditAnnotations != nil {
		in, out := &in.AuditAnnotations, &out.AuditAnnotations
		*out = make([]AuditAnnotation, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ValidatingAdmissionPolicySpec.
func (in *ValidatingAdmissionPolicySpec) DeepCopy() *ValidatingAdmissionPolicySpec {
	if in == nil {
		return nil
	}
	out := new(ValidatingAdmissionPolicySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValidatingAdmissionPolicyStatus) DeepCopyInto(out *ValidatingAdmissionPolicyStatus) {
	*out = *in
	if in.TypeChecking != nil {
		in, out := &in.TypeChecking, &out.TypeChecking
		*out = new(TypeChecking)
		(*in).DeepCopyInto(*out)
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ValidatingAdmissionPolicyStatus.
func (in *ValidatingAdmissionPolicyStatus) DeepCopy() *ValidatingAdmissionPolicyStatus {
	if in == nil {
		return nil
	}
	out := new(ValidatingAdmissionPolicyStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValidatingWebhook) DeepCopyInto(out *ValidatingWebhook) {
	*out = *in
	in.ClientConfig.DeepCopyInto(&out.ClientConfig)
	if in.Rules != nil {
		in, out := &in.Rules, &out.Rules
		*out = make([]RuleWithOperations, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.FailurePolicy != nil {
		in, out := &in.FailurePolicy, &out.FailurePolicy
		*out = new(FailurePolicyType)
		**out = **in
	}
	if in.MatchPolicy != nil {
		in, out := &in.MatchPolicy, &out.MatchPolicy
		*out = new(MatchPolicyType)
		**out = **in
	}
	if in.NamespaceSelector != nil {
		in, out := &in.NamespaceSelector, &out.NamespaceSelector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.ObjectSelector != nil {
		in, out := &in.ObjectSelector, &out.ObjectSelector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.SideEffects != nil {
		in, out := &in.SideEffects, &out.SideEffects
		*out = new(SideEffectClass)
		**out = **in
	}
	if in.TimeoutSeconds != nil {
		in, out := &in.TimeoutSeconds, &out.TimeoutSeconds
		*out = new(int32)
		**out = **in
	}
	if in.AdmissionReviewVersions != nil {
		in, out := &in.AdmissionReviewVersions, &out.AdmissionReviewVersions
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.MatchConditions != nil {
		in, out := &in.MatchConditions, &out.MatchConditions
		*out = make([]MatchCondition, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ValidatingWebhook.
func (in *ValidatingWebhook) DeepCopy() *ValidatingWebhook {
	if in == nil {
		return nil
	}
	out := new(ValidatingWebhook)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValidatingWebhookConfiguration) DeepCopyInto(out *ValidatingWebhookConfiguration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Webhooks != nil {
		in, out := &in.Webhooks, &out.Webhooks
		*out = make([]ValidatingWebhook, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ValidatingWebhookConfiguration.
func (in *ValidatingWebhookConfiguration) DeepCopy() *ValidatingWebhookConfiguration {
	if in == nil {
		return nil
	}
	out := new(ValidatingWebhookConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ValidatingWebhookConfiguration) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValidatingWebhookConfigurationList) DeepCopyInto(out *ValidatingWebhookConfigurationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ValidatingWebhookConfiguration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ValidatingWebhookConfigurationList.
func (in *ValidatingWebhookConfigurationList) DeepCopy() *ValidatingWebhookConfigurationList {
	if in == nil {
		return nil
	}
	out := new(ValidatingWebhookConfigurationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ValidatingWebhookConfigurationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Validation) DeepCopyInto(out *Validation) {
	*out = *in
	if in.Reason != nil {
		in, out := &in.Reason, &out.Reason
		*out = new(v1.StatusReason)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Validation.
func (in *Validation) DeepCopy() *Validation {
	if in == nil {
		return nil
	}
	out := new(Validation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WebhookClientConfig) DeepCopyInto(out *WebhookClientConfig) {
	*out = *in
	if in.URL != nil {
		in, out := &in.URL, &out.URL
		*out = new(string)
		**out = **in
	}
	if in.Service != nil {
		in, out := &in.Service, &out.Service
		*out = new(ServiceReference)
		(*in).DeepCopyInto(*out)
	}
	if in.CABundle != nil {
		in, out := &in.CABundle, &out.CABundle
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WebhookClientConfig.
func (in *WebhookClientConfig) DeepCopy() *WebhookClientConfig {
	if in == nil {
		return nil
	}
	out := new(WebhookClientConfig)
	in.DeepCopyInto(out)
	return out
}
