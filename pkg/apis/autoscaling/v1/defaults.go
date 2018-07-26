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

package v1

import (
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_HorizontalPodAutoscaler(obj *autoscalingv1.HorizontalPodAutoscaler) {
	if obj.Spec.MinReplicas == nil {
		minReplicas := int32(1)
		obj.Spec.MinReplicas = &minReplicas
	}

	// NB: we apply a default for CPU utilization in conversion because
	// we need access to the annotations to properly apply the default.
	// Set autoscalingv1.HorizontalPodAutoscalerSpec.TargetCPUUtilizationPercentage
	// to 100, which has the same meaning as unset.
	if obj.Spec.TargetCPUUtilizationPercentage == nil {
		targetCPUUtilizationPercentage := int32(100)
		obj.Spec.TargetCPUUtilizationPercentage = &targetCPUUtilizationPercentage
	}
}
