/*
Copyright 2015 The Kubernetes Authors.

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

package pod

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	podtest "k8s.io/kubernetes/pkg/kubelet/pod/testing"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
)

// Stub out mirror client for testing purpose.
func newTestManager() (*basicManager, *podtest.FakeMirrorClient) {
	fakeMirrorClient := podtest.NewFakeMirrorClient()
	manager := NewBasicPodManager().(*basicManager)
	return manager, fakeMirrorClient
}

// Tests that pods/maps are properly set after the pod update, and the basic
// methods work correctly.
func TestGetSetPods(t *testing.T) {
	var (
		mirrorPod = &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				UID:       types.UID("mirror-pod-uid"),
				Name:      "mirror-static-pod-name",
				Namespace: metav1.NamespaceDefault,
				Annotations: map[string]string{
					kubetypes.ConfigSourceAnnotationKey: "api",
					kubetypes.ConfigMirrorAnnotationKey: "mirror",
				},
			},
		}
		staticPod = &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				UID:         types.UID("static-pod-uid"),
				Name:        "mirror-static-pod-name",
				Namespace:   metav1.NamespaceDefault,
				Annotations: map[string]string{kubetypes.ConfigSourceAnnotationKey: "file"},
			},
		}
		normalPod = &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				UID:       types.UID("normal-pod-uid"),
				Name:      "normal-pod-name",
				Namespace: metav1.NamespaceDefault,
			},
		}
	)

	testCase := []struct {
		name         string
		podList      []*v1.Pod
		wantPod      *v1.Pod
		expectPods   []*v1.Pod
		expectGetPod *v1.Pod
		expectUID    types.UID
	}{
		{
			name:         "Get normal pod",
			podList:      []*v1.Pod{mirrorPod, staticPod, normalPod},
			wantPod:      normalPod,
			expectPods:   []*v1.Pod{staticPod, normalPod},
			expectGetPod: normalPod,
			expectUID:    types.UID("static-pod-uid"),
		},
		{
			name:         "Get static pod",
			podList:      []*v1.Pod{mirrorPod, staticPod, normalPod},
			wantPod:      staticPod,
			expectPods:   []*v1.Pod{staticPod, normalPod},
			expectGetPod: staticPod,
			expectUID:    types.UID("static-pod-uid"),
		},
	}
	for _, test := range testCase {
		t.Run(test.name, func(t *testing.T) {
			podManager, _ := newTestManager()
			podManager.SetPods(test.podList)
			actualPods := podManager.GetPods()
			sort.Slice(actualPods, func(i, j int) bool {
				return actualPods[i].Name < actualPods[j].Name
			})
			if !reflect.DeepEqual(actualPods, test.expectPods) {
				t.Errorf("Test case %s failed: actualPods and test.expectPods are not equal", test.name)
			}

			// Tests UID translation works as expected. Convert static pod UID for comparison only.
			if uid := podManager.TranslatePodUID(mirrorPod.UID); uid != kubetypes.ResolvedPodUID(test.expectUID) {
				t.Errorf("unable to translate UID %q to the static POD's UID %q; %#v",
					mirrorPod.UID, staticPod.UID, podManager.mirrorPodByUID)
			}
			fullName := fmt.Sprintf("%s_%s", test.wantPod.Name, test.wantPod.Namespace)
			// Test the basic Get methods.
			actualPod, ok := podManager.GetPodByFullName(fullName)
			if !ok || !reflect.DeepEqual(actualPod, test.expectGetPod) {
				t.Errorf("unable to get pod by full name; expected: %#v, got: %#v", test.expectGetPod, actualPod)
			}
			actualPod, ok = podManager.GetPodByName(test.wantPod.Namespace, test.wantPod.Name)
			if !ok || !reflect.DeepEqual(actualPod, test.expectGetPod) {
				t.Errorf("unable to get pod by name; expected: %#v, got: %#v", test.expectGetPod, actualPod)
			}
		})
	}
}

func TestRemovePods(t *testing.T) {
	var (
		mirrorPod = &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				UID:       types.UID("mirror-pod-uid"),
				Name:      "mirror-static-pod-name",
				Namespace: metav1.NamespaceDefault,
				Annotations: map[string]string{
					kubetypes.ConfigSourceAnnotationKey: "api",
					kubetypes.ConfigMirrorAnnotationKey: "mirror",
				},
			},
		}
		staticPod = &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				UID:         types.UID("static-pod-uid"),
				Name:        "mirror-static-pod-name",
				Namespace:   metav1.NamespaceDefault,
				Annotations: map[string]string{kubetypes.ConfigSourceAnnotationKey: "file"},
			},
		}
		normalPod = &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				UID:       types.UID("normal-pod-uid"),
				Name:      "normal-pod-name",
				Namespace: metav1.NamespaceDefault,
			},
		}
	)

	testCase := []struct {
		name             string
		podList          []*v1.Pod
		needToRemovePod  *v1.Pod
		expectPods       []*v1.Pod
		expectMirrorPods []*v1.Pod
	}{
		{
			name:             "Remove mirror pod",
			podList:          []*v1.Pod{mirrorPod, staticPod, normalPod},
			needToRemovePod:  mirrorPod,
			expectPods:       []*v1.Pod{normalPod, staticPod},
			expectMirrorPods: []*v1.Pod{},
		},
		{
			name:             "Remove static pod",
			podList:          []*v1.Pod{mirrorPod, staticPod, normalPod},
			needToRemovePod:  staticPod,
			expectPods:       []*v1.Pod{normalPod},
			expectMirrorPods: []*v1.Pod{mirrorPod},
		},
		{
			name:             "Remove normal pod",
			podList:          []*v1.Pod{mirrorPod, staticPod, normalPod},
			needToRemovePod:  normalPod,
			expectPods:       []*v1.Pod{staticPod},
			expectMirrorPods: []*v1.Pod{mirrorPod},
		},
	}

	for _, test := range testCase {
		t.Run(test.name, func(t *testing.T) {
			podManager, _ := newTestManager()
			podManager.SetPods(test.podList)
			podManager.RemovePod(test.needToRemovePod)
			actualPods1 := podManager.GetPods()
			actualPods2, actualMirrorPods, _ := podManager.GetPodsAndMirrorPods()
			// Sort the pods to make the test deterministic.
			sort.Slice(actualPods1, func(i, j int) bool {
				return actualPods1[i].Name < actualPods1[j].Name
			})
			sort.Slice(actualPods2, func(i, j int) bool {
				return actualPods2[i].Name < actualPods2[j].Name
			})
			sort.Slice(test.expectPods, func(i, j int) bool {
				return test.expectPods[i].Name < test.expectPods[j].Name
			})

			// Check if the actual pods and mirror pods match the expected ones.
			if !reflect.DeepEqual(actualPods1, actualPods2) {
				t.Errorf("Test case %s failed: actualPods1 and actualPods2 are not equal", test.name)
			}
			if !reflect.DeepEqual(actualPods1, test.expectPods) {
				t.Errorf("Test case %s failed: actualPods1 and expectPods are not equal", test.name)
			}
			if !reflect.DeepEqual(actualMirrorPods, test.expectMirrorPods) {
				t.Errorf("Test case %s failed: actualMirrorPods and expectMirrorPods are not equal", test.name)
			}
		})
	}
}
