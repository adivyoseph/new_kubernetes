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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1beta1 "k8s.io/api/networking/v1beta1"
	networkingv1beta1 "k8s.io/client-go/applyconfigurations/networking/v1beta1"
	gentype "k8s.io/client-go/gentype"
	clientset "k8s.io/client-go/kubernetes/typed/networking/v1beta1"
)

// fakeIngresses implements IngressInterface
type fakeIngresses struct {
	*gentype.FakeClientWithListAndApply[*v1beta1.Ingress, *v1beta1.IngressList, *networkingv1beta1.IngressApplyConfiguration]
	Fake *FakeNetworkingV1beta1
}

func newFakeIngresses(fake *FakeNetworkingV1beta1, namespace string) clientset.IngressInterface {
	return &fakeIngresses{
		gentype.NewFakeClientWithListAndApply[*v1beta1.Ingress, *v1beta1.IngressList, *networkingv1beta1.IngressApplyConfiguration](
			fake.Fake,
			namespace,
			v1beta1.SchemeGroupVersion.WithResource("ingresses"),
			v1beta1.SchemeGroupVersion.WithKind("Ingress"),
			func() *v1beta1.Ingress { return &v1beta1.Ingress{} },
			func() *v1beta1.IngressList { return &v1beta1.IngressList{} },
			func(dst, src *v1beta1.IngressList) { dst.ListMeta = src.ListMeta },
			func(list *v1beta1.IngressList) []*v1beta1.Ingress { return gentype.ToPointerSlice(list.Items) },
			func(list *v1beta1.IngressList, items []*v1beta1.Ingress) {
				list.Items = gentype.FromPointerSlice(items)
			},
		),
		fake,
	}
}
