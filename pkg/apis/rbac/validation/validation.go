/*
Copyright 2016 The Kubernetes Authors.

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

package validation

import (
	"k8s.io/kubernetes/pkg/api/validation"
	"k8s.io/kubernetes/pkg/api/validation/path"
	"k8s.io/kubernetes/pkg/apis/rbac"
	"k8s.io/kubernetes/pkg/util/validation/field"
)

// Minimal validation of names for roles and bindings. Identical to the validation for Openshift. See:
// * https://github.com/kubernetes/kubernetes/blob/60db50/pkg/api/validation/name.go
// * https://github.com/openshift/origin/blob/388478/pkg/api/helpers.go
func minimalNameRequirements(name string, prefix bool) []string {
	return path.IsValidPathSegmentName(name)
}

func ValidateAddRoleRequest(obj *rbac.AddRoleRequest) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(obj.Spec.Subjects) == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "subjects"), "must provide at least one subject"))
	}
	subjectsPath := field.NewPath("spec", "subjects")
	for i, subject := range obj.Spec.Subjects {
		allErrs = append(allErrs, validateRoleBindingSubject(subject, true, subjectsPath.Index(i))...)
	}

	allErrs = append(allErrs, validateRoleBindingRoleRef(obj.Spec.RoleRef, field.NewPath("spec", "roleRef"))...)

	return allErrs
}

func validateRoleBindingRoleRef(roleRef rbac.RoleRef, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// TODO allow multiple API groups.  For now, restrict to one, but I can envision other experimental roles in other groups taking
	// advantage of the binding infrastructure
	if roleRef.APIGroup != rbac.GroupName {
		allErrs = append(allErrs, field.NotSupported(field.NewPath("roleRef", "apiGroup"), roleRef.APIGroup, []string{rbac.GroupName}))
	}

	switch roleRef.Kind {
	case "Role", "ClusterRole":
	default:
		allErrs = append(allErrs, field.NotSupported(field.NewPath("roleRef", "kind"), roleRef.Kind, []string{"Role", "ClusterRole"}))

	}

	if len(roleRef.Name) == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("roleRef", "name"), ""))
	} else {
		for _, msg := range minimalNameRequirements(roleRef.Name, false) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("roleRef", "name"), roleRef.Name, msg))
		}
	}

	return allErrs
}

func ValidateRole(role *rbac.Role) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validation.ValidateObjectMeta(&role.ObjectMeta, true, minimalNameRequirements, field.NewPath("metadata"))...)

	for i, rule := range role.Rules {
		if err := validatePolicyRule(rule, true, field.NewPath("rules").Index(i)); err != nil {
			allErrs = append(allErrs, err...)
		}
	}
	if len(allErrs) != 0 {
		return allErrs
	}
	return nil
}

func ValidateRoleUpdate(role *rbac.Role, oldRole *rbac.Role) field.ErrorList {
	allErrs := ValidateRole(role)
	allErrs = append(allErrs, validation.ValidateObjectMetaUpdate(&role.ObjectMeta, &oldRole.ObjectMeta, field.NewPath("metadata"))...)

	return allErrs
}

func ValidateClusterRole(role *rbac.ClusterRole) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validation.ValidateObjectMeta(&role.ObjectMeta, false, minimalNameRequirements, field.NewPath("metadata"))...)

	for i, rule := range role.Rules {
		if err := validatePolicyRule(rule, false, field.NewPath("rules").Index(i)); err != nil {
			allErrs = append(allErrs, err...)
		}
	}
	if len(allErrs) != 0 {
		return allErrs
	}
	return nil
}

func ValidateClusterRoleUpdate(role *rbac.ClusterRole, oldRole *rbac.ClusterRole) field.ErrorList {
	allErrs := ValidateClusterRole(role)
	allErrs = append(allErrs, validation.ValidateObjectMetaUpdate(&role.ObjectMeta, &oldRole.ObjectMeta, field.NewPath("metadata"))...)

	return allErrs
}

func validatePolicyRule(rule rbac.PolicyRule, isNamespaced bool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(rule.Verbs) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("verbs"), "verbs must contain at least one value"))
	}

	if len(rule.NonResourceURLs) > 0 {
		if isNamespaced {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("nonResourceURLs"), rule.NonResourceURLs, "namespaced rules cannot apply to non-resource URLs"))
		}
		if len(rule.APIGroups) > 0 || len(rule.Resources) > 0 || len(rule.ResourceNames) > 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("nonResourceURLs"), rule.NonResourceURLs, "rules cannot apply to both regular resources and non-resource URLs"))
		}
		return allErrs
	}

	if len(rule.APIGroups) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("apiGroups"), "resource rules must supply at least one api group"))
	}
	if len(rule.Resources) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("resources"), "resource rules must supply at least one resource"))
	}
	return allErrs
}

func ValidateRoleBinding(roleBinding *rbac.RoleBinding) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validation.ValidateObjectMeta(&roleBinding.ObjectMeta, true, minimalNameRequirements, field.NewPath("metadata"))...)
	allErrs = append(allErrs, validateRoleBindingRoleRef(roleBinding.RoleRef, field.NewPath("roleRef"))...)

	subjectsPath := field.NewPath("subjects")
	for i, subject := range roleBinding.Subjects {
		allErrs = append(allErrs, validateRoleBindingSubject(subject, true, subjectsPath.Index(i))...)
	}

	return allErrs
}

func ValidateRoleBindingUpdate(roleBinding *rbac.RoleBinding, oldRoleBinding *rbac.RoleBinding) field.ErrorList {
	allErrs := ValidateRoleBinding(roleBinding)
	allErrs = append(allErrs, validation.ValidateObjectMetaUpdate(&roleBinding.ObjectMeta, &oldRoleBinding.ObjectMeta, field.NewPath("metadata"))...)

	if oldRoleBinding.RoleRef != roleBinding.RoleRef {
		allErrs = append(allErrs, field.Invalid(field.NewPath("roleRef"), roleBinding.RoleRef, "cannot change roleRef"))
	}

	return allErrs
}

func ValidateClusterRoleBinding(roleBinding *rbac.ClusterRoleBinding) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validation.ValidateObjectMeta(&roleBinding.ObjectMeta, false, minimalNameRequirements, field.NewPath("metadata"))...)

	allErrs = append(allErrs, validateRoleBindingRoleRef(roleBinding.RoleRef, field.NewPath("roleRef"))...)

	// further restrict the roleref
	switch roleBinding.RoleRef.Kind {
	case "ClusterRole":
	default:
		allErrs = append(allErrs, field.NotSupported(field.NewPath("roleRef", "kind"), roleBinding.RoleRef.Kind, []string{"ClusterRole"}))

	}

	subjectsPath := field.NewPath("subjects")
	for i, subject := range roleBinding.Subjects {
		allErrs = append(allErrs, validateRoleBindingSubject(subject, false, subjectsPath.Index(i))...)
	}

	return allErrs
}

func ValidateClusterRoleBindingUpdate(roleBinding *rbac.ClusterRoleBinding, oldRoleBinding *rbac.ClusterRoleBinding) field.ErrorList {
	allErrs := ValidateClusterRoleBinding(roleBinding)
	allErrs = append(allErrs, validation.ValidateObjectMetaUpdate(&roleBinding.ObjectMeta, &oldRoleBinding.ObjectMeta, field.NewPath("metadata"))...)

	if oldRoleBinding.RoleRef != roleBinding.RoleRef {
		allErrs = append(allErrs, field.Invalid(field.NewPath("roleRef"), roleBinding.RoleRef, "cannot change roleRef"))
	}

	return allErrs
}

func validateRoleBindingSubject(subject rbac.Subject, isNamespaced bool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(subject.Name) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("name"), ""))
	}

	switch subject.Kind {
	case rbac.ServiceAccountKind:
		if len(subject.Name) > 0 {
			for _, msg := range validation.ValidateServiceAccountName(subject.Name, false) {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("name"), subject.Name, msg))
			}
		}
		if !isNamespaced && len(subject.Namespace) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("namespace"), ""))
		}

	case rbac.UserKind:
		// TODO(ericchiang): What other restrictions on user name are there?
		if len(subject.Name) == 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("name"), subject.Name, "user name cannot be empty"))
		}

	case rbac.GroupKind:
		// TODO(ericchiang): What other restrictions on group name are there?
		if len(subject.Name) == 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("name"), subject.Name, "group name cannot be empty"))
		}

	default:
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("kind"), subject.Kind, []string{rbac.ServiceAccountKind, rbac.UserKind, rbac.GroupKind}))
	}

	return allErrs
}
