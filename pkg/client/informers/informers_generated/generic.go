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

// This file was automatically generated by informer-gen with arguments: --input-dirs=[k8s.io/kubernetes/pkg/api,k8s.io/kubernetes/pkg/api/v1,k8s.io/kubernetes/pkg/apis/abac,k8s.io/kubernetes/pkg/apis/abac/v0,k8s.io/kubernetes/pkg/apis/abac/v1beta1,k8s.io/kubernetes/pkg/apis/apps,k8s.io/kubernetes/pkg/apis/apps/v1beta1,k8s.io/kubernetes/pkg/apis/authentication,k8s.io/kubernetes/pkg/apis/authentication/v1beta1,k8s.io/kubernetes/pkg/apis/authorization,k8s.io/kubernetes/pkg/apis/authorization/v1beta1,k8s.io/kubernetes/pkg/apis/autoscaling,k8s.io/kubernetes/pkg/apis/autoscaling/v1,k8s.io/kubernetes/pkg/apis/batch,k8s.io/kubernetes/pkg/apis/batch/v1,k8s.io/kubernetes/pkg/apis/batch/v2alpha1,k8s.io/kubernetes/pkg/apis/certificates,k8s.io/kubernetes/pkg/apis/certificates/v1alpha1,k8s.io/kubernetes/pkg/apis/componentconfig,k8s.io/kubernetes/pkg/apis/componentconfig/v1alpha1,k8s.io/kubernetes/pkg/apis/extensions,k8s.io/kubernetes/pkg/apis/extensions/v1beta1,k8s.io/kubernetes/pkg/apis/imagepolicy,k8s.io/kubernetes/pkg/apis/imagepolicy/v1alpha1,k8s.io/kubernetes/pkg/apis/meta/v1,k8s.io/kubernetes/pkg/apis/policy,k8s.io/kubernetes/pkg/apis/policy/v1beta1,k8s.io/kubernetes/pkg/apis/rbac,k8s.io/kubernetes/pkg/apis/rbac/v1alpha1,k8s.io/kubernetes/pkg/apis/storage,k8s.io/kubernetes/pkg/apis/storage/v1beta1] --internal-clientset-package=k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset --listers-package=k8s.io/kubernetes/pkg/client/listers --versioned-clientset-package=k8s.io/kubernetes/pkg/client/clientset_generated/clientset

package informers_generated

import (
	"fmt"
	api "k8s.io/kubernetes/pkg/api"
	api_v1 "k8s.io/kubernetes/pkg/api/v1"
	apps "k8s.io/kubernetes/pkg/apis/apps"
	v1beta1 "k8s.io/kubernetes/pkg/apis/apps/v1beta1"
	autoscaling "k8s.io/kubernetes/pkg/apis/autoscaling"
	v1 "k8s.io/kubernetes/pkg/apis/autoscaling/v1"
	batch "k8s.io/kubernetes/pkg/apis/batch"
	batch_v1 "k8s.io/kubernetes/pkg/apis/batch/v1"
	v2alpha1 "k8s.io/kubernetes/pkg/apis/batch/v2alpha1"
	certificates "k8s.io/kubernetes/pkg/apis/certificates"
	v1alpha1 "k8s.io/kubernetes/pkg/apis/certificates/v1alpha1"
	extensions "k8s.io/kubernetes/pkg/apis/extensions"
	extensions_v1beta1 "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	policy "k8s.io/kubernetes/pkg/apis/policy"
	policy_v1beta1 "k8s.io/kubernetes/pkg/apis/policy/v1beta1"
	rbac "k8s.io/kubernetes/pkg/apis/rbac"
	rbac_v1alpha1 "k8s.io/kubernetes/pkg/apis/rbac/v1alpha1"
	storage "k8s.io/kubernetes/pkg/apis/storage"
	storage_v1beta1 "k8s.io/kubernetes/pkg/apis/storage/v1beta1"
	cache "k8s.io/kubernetes/pkg/client/cache"
	schema "k8s.io/kubernetes/pkg/runtime/schema"
)

// GenericInformer is type of SharedIndexInformer which will locate and delegate to other
// sharedInformers based on type
type GenericInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() cache.GenericLister
}

type genericInformer struct {
	informer cache.SharedIndexInformer
	resource schema.GroupResource
}

// Informer returns the SharedIndexInformer.
func (f *genericInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

// Lister returns the GenericLister.
func (f *genericInformer) Lister() cache.GenericLister {
	return cache.NewGenericLister(f.Informer().GetIndexer(), f.resource)
}

// ForResource gives generic access to a shared informer of the matching type
// TODO extend this to unknown resources with a client pool
func (f *sharedInformerFactory) ForResource(resource schema.GroupVersionResource) (GenericInformer, error) {
	switch resource {
	// Group=Apis, Version=Apps
	case apps.SchemeGroupVersion.WithResource("statefulsets"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Apis().Apps().StatefulSets().Informer()}, nil

		// Group=Apis, Version=Autoscaling
	case autoscaling.SchemeGroupVersion.WithResource("horizontalpodautoscalers"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Apis().Autoscaling().HorizontalPodAutoscalers().Informer()}, nil

		// Group=Apis, Version=Batch
	case batch.SchemeGroupVersion.WithResource("cronjobs"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Apis().Batch().CronJobs().Informer()}, nil
	case batch.SchemeGroupVersion.WithResource("jobs"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Apis().Batch().Jobs().Informer()}, nil

		// Group=Apis, Version=Certificates
	case certificates.SchemeGroupVersion.WithResource("certificatesigningrequests"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Apis().Certificates().CertificateSigningRequests().Informer()}, nil

		// Group=Apis, Version=Extensions
	case extensions.SchemeGroupVersion.WithResource("daemonsets"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Apis().Extensions().DaemonSets().Informer()}, nil
	case extensions.SchemeGroupVersion.WithResource("deployments"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Apis().Extensions().Deployments().Informer()}, nil
	case extensions.SchemeGroupVersion.WithResource("ingresses"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Apis().Extensions().Ingresses().Informer()}, nil
	case extensions.SchemeGroupVersion.WithResource("networkpolicies"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Apis().Extensions().NetworkPolicies().Informer()}, nil
	case extensions.SchemeGroupVersion.WithResource("podsecuritypolicies"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Apis().Extensions().PodSecurityPolicies().Informer()}, nil
	case extensions.SchemeGroupVersion.WithResource("replicasets"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Apis().Extensions().ReplicaSets().Informer()}, nil
	case extensions.SchemeGroupVersion.WithResource("thirdpartyresources"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Apis().Extensions().ThirdPartyResources().Informer()}, nil

		// Group=Apis, Version=Policy
	case policy.SchemeGroupVersion.WithResource("poddisruptionbudgets"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Apis().Policy().PodDisruptionBudgets().Informer()}, nil

		// Group=Apis, Version=Rbac
	case rbac.SchemeGroupVersion.WithResource("clusterroles"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Apis().Rbac().ClusterRoles().Informer()}, nil
	case rbac.SchemeGroupVersion.WithResource("clusterrolebindings"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Apis().Rbac().ClusterRoleBindings().Informer()}, nil
	case rbac.SchemeGroupVersion.WithResource("roles"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Apis().Rbac().Roles().Informer()}, nil
	case rbac.SchemeGroupVersion.WithResource("rolebindings"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Apis().Rbac().RoleBindings().Informer()}, nil

		// Group=Apis, Version=Storage
	case storage.SchemeGroupVersion.WithResource("storageclasses"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Apis().Storage().StorageClasses().Informer()}, nil

		// Group=Apps, Version=V1beta1
	case v1beta1.SchemeGroupVersion.WithResource("statefulsets"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Apps().V1beta1().StatefulSets().Informer()}, nil

		// Group=Autoscaling, Version=V1
	case v1.SchemeGroupVersion.WithResource("horizontalpodautoscalers"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Autoscaling().V1().HorizontalPodAutoscalers().Informer()}, nil

		// Group=Batch, Version=V1
	case batch_v1.SchemeGroupVersion.WithResource("jobs"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Batch().V1().Jobs().Informer()}, nil

		// Group=Batch, Version=V2alpha1
	case v2alpha1.SchemeGroupVersion.WithResource("cronjobs"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Batch().V2alpha1().CronJobs().Informer()}, nil
	case v2alpha1.SchemeGroupVersion.WithResource("jobs"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Batch().V2alpha1().Jobs().Informer()}, nil

		// Group=Certificates, Version=V1alpha1
	case v1alpha1.SchemeGroupVersion.WithResource("certificatesigningrequests"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Certificates().V1alpha1().CertificateSigningRequests().Informer()}, nil

		// Group=Core, Version=V1
	case api_v1.SchemeGroupVersion.WithResource("componentstatuses"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Core().V1().ComponentStatuses().Informer()}, nil
	case api_v1.SchemeGroupVersion.WithResource("configmaps"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Core().V1().ConfigMaps().Informer()}, nil
	case api_v1.SchemeGroupVersion.WithResource("endpoints"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Core().V1().Endpoints().Informer()}, nil
	case api_v1.SchemeGroupVersion.WithResource("events"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Core().V1().Events().Informer()}, nil
	case api_v1.SchemeGroupVersion.WithResource("limitranges"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Core().V1().LimitRanges().Informer()}, nil
	case api_v1.SchemeGroupVersion.WithResource("namespaces"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Core().V1().Namespaces().Informer()}, nil
	case api_v1.SchemeGroupVersion.WithResource("nodes"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Core().V1().Nodes().Informer()}, nil
	case api_v1.SchemeGroupVersion.WithResource("persistentvolumes"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Core().V1().PersistentVolumes().Informer()}, nil
	case api_v1.SchemeGroupVersion.WithResource("persistentvolumeclaims"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Core().V1().PersistentVolumeClaims().Informer()}, nil
	case api_v1.SchemeGroupVersion.WithResource("pods"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Core().V1().Pods().Informer()}, nil
	case api_v1.SchemeGroupVersion.WithResource("podtemplates"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Core().V1().PodTemplates().Informer()}, nil
	case api_v1.SchemeGroupVersion.WithResource("replicationcontrollers"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Core().V1().ReplicationControllers().Informer()}, nil
	case api_v1.SchemeGroupVersion.WithResource("resourcequotas"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Core().V1().ResourceQuotas().Informer()}, nil
	case api_v1.SchemeGroupVersion.WithResource("secrets"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Core().V1().Secrets().Informer()}, nil
	case api_v1.SchemeGroupVersion.WithResource("services"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Core().V1().Services().Informer()}, nil
	case api_v1.SchemeGroupVersion.WithResource("serviceaccounts"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Core().V1().ServiceAccounts().Informer()}, nil

		// Group=Extensions, Version=V1beta1
	case extensions_v1beta1.SchemeGroupVersion.WithResource("daemonsets"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Extensions().V1beta1().DaemonSets().Informer()}, nil
	case extensions_v1beta1.SchemeGroupVersion.WithResource("deployments"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Extensions().V1beta1().Deployments().Informer()}, nil
	case extensions_v1beta1.SchemeGroupVersion.WithResource("ingresses"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Extensions().V1beta1().Ingresses().Informer()}, nil
	case extensions_v1beta1.SchemeGroupVersion.WithResource("podsecuritypolicies"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Extensions().V1beta1().PodSecurityPolicies().Informer()}, nil
	case extensions_v1beta1.SchemeGroupVersion.WithResource("replicasets"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Extensions().V1beta1().ReplicaSets().Informer()}, nil
	case extensions_v1beta1.SchemeGroupVersion.WithResource("thirdpartyresources"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Extensions().V1beta1().ThirdPartyResources().Informer()}, nil

		// Group=Pkg, Version=Api
	case api.SchemeGroupVersion.WithResource("componentstatuses"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Pkg().Api().ComponentStatuses().Informer()}, nil
	case api.SchemeGroupVersion.WithResource("configmaps"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Pkg().Api().ConfigMaps().Informer()}, nil
	case api.SchemeGroupVersion.WithResource("endpoints"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Pkg().Api().Endpoints().Informer()}, nil
	case api.SchemeGroupVersion.WithResource("events"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Pkg().Api().Events().Informer()}, nil
	case api.SchemeGroupVersion.WithResource("limitranges"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Pkg().Api().LimitRanges().Informer()}, nil
	case api.SchemeGroupVersion.WithResource("namespaces"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Pkg().Api().Namespaces().Informer()}, nil
	case api.SchemeGroupVersion.WithResource("nodes"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Pkg().Api().Nodes().Informer()}, nil
	case api.SchemeGroupVersion.WithResource("persistentvolumes"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Pkg().Api().PersistentVolumes().Informer()}, nil
	case api.SchemeGroupVersion.WithResource("persistentvolumeclaims"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Pkg().Api().PersistentVolumeClaims().Informer()}, nil
	case api.SchemeGroupVersion.WithResource("pods"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Pkg().Api().Pods().Informer()}, nil
	case api.SchemeGroupVersion.WithResource("podtemplates"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Pkg().Api().PodTemplates().Informer()}, nil
	case api.SchemeGroupVersion.WithResource("replicationcontrollers"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Pkg().Api().ReplicationControllers().Informer()}, nil
	case api.SchemeGroupVersion.WithResource("resourcequotas"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Pkg().Api().ResourceQuotas().Informer()}, nil
	case api.SchemeGroupVersion.WithResource("secrets"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Pkg().Api().Secrets().Informer()}, nil
	case api.SchemeGroupVersion.WithResource("services"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Pkg().Api().Services().Informer()}, nil
	case api.SchemeGroupVersion.WithResource("serviceaccounts"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Pkg().Api().ServiceAccounts().Informer()}, nil

		// Group=Policy, Version=V1beta1
	case policy_v1beta1.SchemeGroupVersion.WithResource("poddisruptionbudgets"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Policy().V1beta1().PodDisruptionBudgets().Informer()}, nil

		// Group=Rbac, Version=V1alpha1
	case rbac_v1alpha1.SchemeGroupVersion.WithResource("clusterroles"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Rbac().V1alpha1().ClusterRoles().Informer()}, nil
	case rbac_v1alpha1.SchemeGroupVersion.WithResource("clusterrolebindings"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Rbac().V1alpha1().ClusterRoleBindings().Informer()}, nil
	case rbac_v1alpha1.SchemeGroupVersion.WithResource("roles"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Rbac().V1alpha1().Roles().Informer()}, nil
	case rbac_v1alpha1.SchemeGroupVersion.WithResource("rolebindings"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Rbac().V1alpha1().RoleBindings().Informer()}, nil

		// Group=Storage, Version=V1beta1
	case storage_v1beta1.SchemeGroupVersion.WithResource("storageclasses"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Storage().V1beta1().StorageClasses().Informer()}, nil

	}

	return nil, fmt.Errorf("no informer found for %v", resource)
}
