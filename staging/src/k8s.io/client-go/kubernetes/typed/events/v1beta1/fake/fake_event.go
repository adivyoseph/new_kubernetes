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
	v1beta1 "k8s.io/api/events/v1beta1"
	eventsv1beta1 "k8s.io/client-go/applyconfigurations/events/v1beta1"
	gentype "k8s.io/client-go/gentype"
	clientset "k8s.io/client-go/kubernetes/typed/events/v1beta1"
)

// fakeEvents implements EventInterface
type fakeEvents struct {
	*gentype.FakeClientWithListAndApply[*v1beta1.Event, *v1beta1.EventList, *eventsv1beta1.EventApplyConfiguration]
	Fake *FakeEventsV1beta1
}

func newFakeEvents(fake *FakeEventsV1beta1, namespace string) clientset.EventInterface {
	return &fakeEvents{
		gentype.NewFakeClientWithListAndApply[*v1beta1.Event, *v1beta1.EventList, *eventsv1beta1.EventApplyConfiguration](
			fake.Fake,
			namespace,
			v1beta1.SchemeGroupVersion.WithResource("events"),
			v1beta1.SchemeGroupVersion.WithKind("Event"),
			func() *v1beta1.Event { return &v1beta1.Event{} },
			func() *v1beta1.EventList { return &v1beta1.EventList{} },
			func(dst, src *v1beta1.EventList) { dst.ListMeta = src.ListMeta },
			func(list *v1beta1.EventList) []*v1beta1.Event { return gentype.ToPointerSlice(list.Items) },
			func(list *v1beta1.EventList, items []*v1beta1.Event) { list.Items = gentype.FromPointerSlice(items) },
		),
		fake,
	}
}
