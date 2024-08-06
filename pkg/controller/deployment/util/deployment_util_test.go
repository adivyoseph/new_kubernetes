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

package util

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apiserver/pkg/storage/names"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
	featuregatetesting "k8s.io/component-base/featuregate/testing"
	"k8s.io/klog/v2/ktesting"
	"k8s.io/kubernetes/pkg/controller"
	"k8s.io/kubernetes/pkg/features"
	"k8s.io/utils/ptr"
)

func newDControllerRef(d *apps.Deployment) *metav1.OwnerReference {
	isController := true
	return &metav1.OwnerReference{
		APIVersion: "apps/v1",
		Kind:       "Deployment",
		Name:       d.GetName(),
		UID:        d.GetUID(),
		Controller: &isController,
	}
}

// generateRS creates a replica set, with the input deployment's template as its template
func generateRS(deployment apps.Deployment) apps.ReplicaSet {
	return generateRSWithReplicas(deployment, 0, apps.ReplicaSetStatus{})
}

// generateRSWithReplicas creates a replica set, with the input deployment's template as its template
// and includes replicas and the status in the replica set
func generateRSWithReplicas(deployment apps.Deployment, replicas int32, status apps.ReplicaSetStatus) apps.ReplicaSet {
	template := deployment.Spec.Template.DeepCopy()
	return apps.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			UID:             randomUID(),
			Name:            names.SimpleNameGenerator.GenerateName("replicaset"),
			Labels:          template.Labels,
			OwnerReferences: []metav1.OwnerReference{*newDControllerRef(&deployment)},
		},
		Spec: apps.ReplicaSetSpec{
			Replicas: ptr.To(replicas),
			Template: *template,
			Selector: &metav1.LabelSelector{MatchLabels: template.Labels},
		},
		Status: status,
	}
}

func randomUID() types.UID {
	return types.UID(strconv.FormatInt(rand.Int63(), 10))
}

func generateDeployment(image string) apps.Deployment {
	return generateDeploymentWithReplicas(image, 1)
}

// generateDeployment creates a deployment, with the input image as its template
func generateDeploymentWithReplicas(image string, replicas int32) apps.Deployment {
	podLabels := map[string]string{"name": image}
	terminationSec := int64(30)
	enableServiceLinks := v1.DefaultEnableServiceLinks
	return apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        image,
			Annotations: make(map[string]string),
		},
		Spec: apps.DeploymentSpec{
			Replicas: ptr.To(replicas),
			Selector: &metav1.LabelSelector{MatchLabels: podLabels},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: podLabels,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:                   image,
							Image:                  image,
							ImagePullPolicy:        v1.PullAlways,
							TerminationMessagePath: v1.TerminationMessagePathDefault,
						},
					},
					DNSPolicy:                     v1.DNSClusterFirst,
					TerminationGracePeriodSeconds: &terminationSec,
					RestartPolicy:                 v1.RestartPolicyAlways,
					SecurityContext:               &v1.PodSecurityContext{},
					EnableServiceLinks:            &enableServiceLinks,
				},
			},
		},
	}
}

func generatePodTemplateSpec(name, nodeName string, annotations, labels map[string]string) v1.PodTemplateSpec {
	return v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Annotations: annotations,
			Labels:      labels,
		},
		Spec: v1.PodSpec{
			NodeName: nodeName,
		},
	}
}

func TestEqualIgnoreHash(t *testing.T) {
	tests := []struct {
		Name           string
		former, latter v1.PodTemplateSpec
		expected       bool
	}{
		{
			"Same spec, same labels",
			generatePodTemplateSpec("foo", "foo-node", map[string]string{}, map[string]string{apps.DefaultDeploymentUniqueLabelKey: "value-1", "something": "else"}),
			generatePodTemplateSpec("foo", "foo-node", map[string]string{}, map[string]string{apps.DefaultDeploymentUniqueLabelKey: "value-1", "something": "else"}),
			true,
		},
		{
			"Same spec, only pod-template-hash label value is different",
			generatePodTemplateSpec("foo", "foo-node", map[string]string{}, map[string]string{apps.DefaultDeploymentUniqueLabelKey: "value-1", "something": "else"}),
			generatePodTemplateSpec("foo", "foo-node", map[string]string{}, map[string]string{apps.DefaultDeploymentUniqueLabelKey: "value-2", "something": "else"}),
			true,
		},
		{
			"Same spec, the former doesn't have pod-template-hash label",
			generatePodTemplateSpec("foo", "foo-node", map[string]string{}, map[string]string{"something": "else"}),
			generatePodTemplateSpec("foo", "foo-node", map[string]string{}, map[string]string{apps.DefaultDeploymentUniqueLabelKey: "value-2", "something": "else"}),
			true,
		},
		{
			"Same spec, the label is different, the former doesn't have pod-template-hash label, same number of labels",
			generatePodTemplateSpec("foo", "foo-node", map[string]string{}, map[string]string{"something": "else"}),
			generatePodTemplateSpec("foo", "foo-node", map[string]string{}, map[string]string{apps.DefaultDeploymentUniqueLabelKey: "value-2"}),
			false,
		},
		{
			"Same spec, the label is different, the latter doesn't have pod-template-hash label, same number of labels",
			generatePodTemplateSpec("foo", "foo-node", map[string]string{}, map[string]string{apps.DefaultDeploymentUniqueLabelKey: "value-1"}),
			generatePodTemplateSpec("foo", "foo-node", map[string]string{}, map[string]string{"something": "else"}),
			false,
		},
		{
			"Same spec, the label is different, and the pod-template-hash label value is the same",
			generatePodTemplateSpec("foo", "foo-node", map[string]string{}, map[string]string{apps.DefaultDeploymentUniqueLabelKey: "value-1"}),
			generatePodTemplateSpec("foo", "foo-node", map[string]string{}, map[string]string{apps.DefaultDeploymentUniqueLabelKey: "value-1", "something": "else"}),
			false,
		},
		{
			"Different spec, same labels",
			generatePodTemplateSpec("foo", "foo-node", map[string]string{"former": "value"}, map[string]string{apps.DefaultDeploymentUniqueLabelKey: "value-1", "something": "else"}),
			generatePodTemplateSpec("foo", "foo-node", map[string]string{"latter": "value"}, map[string]string{apps.DefaultDeploymentUniqueLabelKey: "value-1", "something": "else"}),
			false,
		},
		{
			"Different spec, different pod-template-hash label value",
			generatePodTemplateSpec("foo-1", "foo-node", map[string]string{}, map[string]string{apps.DefaultDeploymentUniqueLabelKey: "value-1", "something": "else"}),
			generatePodTemplateSpec("foo-2", "foo-node", map[string]string{}, map[string]string{apps.DefaultDeploymentUniqueLabelKey: "value-2", "something": "else"}),
			false,
		},
		{
			"Different spec, the former doesn't have pod-template-hash label",
			generatePodTemplateSpec("foo-1", "foo-node-1", map[string]string{}, map[string]string{"something": "else"}),
			generatePodTemplateSpec("foo-2", "foo-node-2", map[string]string{}, map[string]string{apps.DefaultDeploymentUniqueLabelKey: "value-2", "something": "else"}),
			false,
		},
		{
			"Different spec, different labels",
			generatePodTemplateSpec("foo", "foo-node-1", map[string]string{}, map[string]string{"something": "else"}),
			generatePodTemplateSpec("foo", "foo-node-2", map[string]string{}, map[string]string{"nothing": "else"}),
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runTest := func(t1, t2 *v1.PodTemplateSpec, reversed bool) {
				reverseString := ""
				if reversed {
					reverseString = " (reverse order)"
				}
				// Run
				equal := EqualIgnoreHash(t1, t2)
				if equal != test.expected {
					t.Errorf("%q%s: expected %v", test.Name, reverseString, test.expected)
					return
				}
				if t1.Labels == nil || t2.Labels == nil {
					t.Errorf("%q%s: unexpected labels becomes nil", test.Name, reverseString)
				}
			}

			runTest(&test.former, &test.latter, false)
			// Test the same case in reverse order
			runTest(&test.latter, &test.former, true)
		})
	}
}

func TestFindNewReplicaSet(t *testing.T) {
	now := metav1.Now()
	later := metav1.Time{Time: now.Add(time.Minute)}

	deployment := generateDeployment("nginx")
	newRS := generateRS(deployment)
	newRS.Labels[apps.DefaultDeploymentUniqueLabelKey] = "hash"
	newRS.CreationTimestamp = later

	newRSDup := generateRS(deployment)
	newRSDup.Labels[apps.DefaultDeploymentUniqueLabelKey] = "different-hash"
	newRSDup.CreationTimestamp = now

	oldDeployment := generateDeployment("nginx")
	oldDeployment.Spec.Template.Spec.Containers[0].Name = "nginx-old-1"
	oldRS := generateRS(oldDeployment)

	tests := []struct {
		Name       string
		deployment apps.Deployment
		rsList     []*apps.ReplicaSet
		expected   *apps.ReplicaSet
	}{
		{
			Name:       "Get new ReplicaSet with the same template as Deployment spec but different pod-template-hash value",
			deployment: deployment,
			rsList:     []*apps.ReplicaSet{&newRS, &oldRS},
			expected:   &newRS,
		},
		{
			Name:       "Get the oldest new ReplicaSet when there are more than one ReplicaSet with the same template",
			deployment: deployment,
			rsList:     []*apps.ReplicaSet{&newRS, &oldRS, &newRSDup},
			expected:   &newRSDup,
		},
		{
			Name:       "Get nil new ReplicaSet",
			deployment: deployment,
			rsList:     []*apps.ReplicaSet{&oldRS},
			expected:   nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			if rs := FindNewReplicaSet(&test.deployment, test.rsList); !reflect.DeepEqual(rs, test.expected) {
				t.Errorf("In test case %q, expected %#v, got %#v", test.Name, test.expected, rs)
			}
		})
	}
}

func TestFindOldReplicaSets(t *testing.T) {
	now := metav1.Now()
	later := metav1.Time{Time: now.Add(time.Minute)}
	before := metav1.Time{Time: now.Add(-time.Minute)}

	deployment := generateDeployment("nginx")
	newRS := generateRSWithReplicas(deployment, 1, apps.ReplicaSetStatus{})
	newRS.Labels[apps.DefaultDeploymentUniqueLabelKey] = "hash"
	newRS.CreationTimestamp = later

	newRSDup := generateRS(deployment)
	newRSDup.Labels[apps.DefaultDeploymentUniqueLabelKey] = "different-hash"
	newRSDup.CreationTimestamp = now

	oldDeployment := generateDeployment("nginx")
	oldDeployment.Spec.Template.Spec.Containers[0].Name = "nginx-old-1"
	oldRS := generateRS(oldDeployment)
	oldRS.CreationTimestamp = before

	tests := []struct {
		Name            string
		deployment      apps.Deployment
		rsList          []*apps.ReplicaSet
		expected        []*apps.ReplicaSet
		expectedRequire []*apps.ReplicaSet
	}{
		{
			Name:            "Get old ReplicaSets",
			deployment:      deployment,
			rsList:          []*apps.ReplicaSet{&newRS, &oldRS},
			expected:        []*apps.ReplicaSet{&oldRS},
			expectedRequire: nil,
		},
		{
			Name:            "Get old ReplicaSets with no new ReplicaSet",
			deployment:      deployment,
			rsList:          []*apps.ReplicaSet{&oldRS},
			expected:        []*apps.ReplicaSet{&oldRS},
			expectedRequire: nil,
		},
		{
			Name:            "Get old ReplicaSets with two new ReplicaSets, only the oldest new ReplicaSet is seen as new ReplicaSet",
			deployment:      deployment,
			rsList:          []*apps.ReplicaSet{&oldRS, &newRS, &newRSDup},
			expected:        []*apps.ReplicaSet{&oldRS, &newRS},
			expectedRequire: []*apps.ReplicaSet{&newRS},
		},
		{
			Name:            "Get empty old ReplicaSets",
			deployment:      deployment,
			rsList:          []*apps.ReplicaSet{&newRS},
			expected:        nil,
			expectedRequire: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			requireRS, allRS := FindOldReplicaSets(&test.deployment, test.rsList)
			sort.Sort(controller.ReplicaSetsByCreationTimestamp(allRS))
			sort.Sort(controller.ReplicaSetsByCreationTimestamp(test.expected))
			if !reflect.DeepEqual(allRS, test.expected) {
				t.Errorf("In test case %q, expected %#v, got %#v", test.Name, test.expected, allRS)
			}
			// RSs are getting filtered correctly by rs.spec.replicas
			if !reflect.DeepEqual(requireRS, test.expectedRequire) {
				t.Errorf("In test case %q, expected %#v, got %#v", test.Name, test.expectedRequire, requireRS)
			}
		})
	}
}

func TestGetReplicaCountForReplicaSets(t *testing.T) {
	rs1 := generateRSWithReplicas(generateDeployment("foo"), 1, apps.ReplicaSetStatus{Replicas: 2, TerminatingReplicas: 3})
	rs2 := generateRSWithReplicas(generateDeployment("bar"), 2, apps.ReplicaSetStatus{Replicas: 3, TerminatingReplicas: 1})
	rs3 := generateRSWithReplicas(generateDeployment("unsynced"), 3, apps.ReplicaSetStatus{})

	tests := []struct {
		name                     string
		sets                     []*apps.ReplicaSet
		expectedCount            int32
		expectedActual           int32
		expectedTerminating      int32
		expectSurgeCapacityCount int32
	}{
		{
			name:                     "scaling down rs1",
			sets:                     []*apps.ReplicaSet{&rs1},
			expectedCount:            1,
			expectedActual:           2,
			expectedTerminating:      3,
			expectSurgeCapacityCount: 2,
		},
		{
			name:                     "scaling down rs1 and rs2",
			sets:                     []*apps.ReplicaSet{&rs1, &rs2},
			expectedCount:            3,
			expectedActual:           5,
			expectedTerminating:      4,
			expectSurgeCapacityCount: 5,
		},
		{
			name:                     "scaling up rs3",
			sets:                     []*apps.ReplicaSet{&rs3},
			expectedCount:            3,
			expectedActual:           0,
			expectedTerminating:      0,
			expectSurgeCapacityCount: 3,
		},
		{
			name:                     "scaling down rs1 and rs2 and scaling up rs3",
			sets:                     []*apps.ReplicaSet{&rs1, &rs2, &rs3},
			expectedCount:            6,
			expectedActual:           5,
			expectedTerminating:      4,
			expectSurgeCapacityCount: 8,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			count := GetReplicaCountForReplicaSets(test.sets)
			if count != test.expectedCount {
				t.Errorf("expectedCount %+v, got %+v", test.expectedCount, count)
			}
			count = GetActualReplicaCountForReplicaSets(test.sets)
			if count != test.expectedActual {
				t.Errorf("expectedActual %+v, got %+v", test.expectedActual, count)
			}
			count = GetTerminatingReplicaCountForReplicaSets(test.sets)
			if count != test.expectedTerminating {
				t.Errorf("expectedTerminating %+v, got %+v", test.expectedTerminating, count)
			}
			count = GetReplicaSurgeCapacityCountForReplicaSets(test.sets)
			if count != test.expectSurgeCapacityCount {
				t.Errorf("expectSurgeCapacityCount %+v, got %+v", test.expectSurgeCapacityCount, count)
			}
		})
	}
}

func TestResolveFenceposts(t *testing.T) {
	tests := []struct {
		maxSurge          *string
		maxUnavailable    *string
		desired           int32
		expectSurge       int32
		expectUnavailable int32
		expectError       bool
	}{
		{
			maxSurge:          ptr.To("0%"),
			maxUnavailable:    ptr.To("0%"),
			desired:           0,
			expectSurge:       0,
			expectUnavailable: 1,
			expectError:       false,
		},
		{
			maxSurge:          ptr.To("39%"),
			maxUnavailable:    ptr.To("39%"),
			desired:           10,
			expectSurge:       4,
			expectUnavailable: 3,
			expectError:       false,
		},
		{
			maxSurge:          ptr.To("oops"),
			maxUnavailable:    ptr.To("39%"),
			desired:           10,
			expectSurge:       0,
			expectUnavailable: 0,
			expectError:       true,
		},
		{
			maxSurge:          ptr.To("55%"),
			maxUnavailable:    ptr.To("urg"),
			desired:           10,
			expectSurge:       0,
			expectUnavailable: 0,
			expectError:       true,
		},
		{
			maxSurge:          nil,
			maxUnavailable:    ptr.To("39%"),
			desired:           10,
			expectSurge:       0,
			expectUnavailable: 3,
			expectError:       false,
		},
		{
			maxSurge:          ptr.To("39%"),
			maxUnavailable:    nil,
			desired:           10,
			expectSurge:       4,
			expectUnavailable: 0,
			expectError:       false,
		},
		{
			maxSurge:          nil,
			maxUnavailable:    nil,
			desired:           10,
			expectSurge:       0,
			expectUnavailable: 1,
			expectError:       false,
		},
	}

	for num, test := range tests {
		t.Run(fmt.Sprintf("%d", num), func(t *testing.T) {
			var maxSurge, maxUnavail *intstr.IntOrString
			if test.maxSurge != nil {
				maxSurge = ptr.To(intstr.FromString(*test.maxSurge))
			}
			if test.maxUnavailable != nil {
				maxUnavail = ptr.To(intstr.FromString(*test.maxUnavailable))
			}
			surge, unavail, err := ResolveFenceposts(maxSurge, maxUnavail, test.desired)
			if err != nil && !test.expectError {
				t.Errorf("unexpected error %v", err)
			}
			if err == nil && test.expectError {
				t.Error("expected error")
			}
			if surge != test.expectSurge || unavail != test.expectUnavailable {
				t.Errorf("#%v got %v:%v, want %v:%v", num, surge, unavail, test.expectSurge, test.expectUnavailable)
			}
		})
	}
}

func TestNewRSNewReplicas(t *testing.T) {
	tests := []struct {
		name        string
		replicas    int32
		strategy    apps.DeploymentStrategy
		replicaSets []*apps.ReplicaSet
		expected    int32
	}{
		{
			name:     "empty deployment does not scale",
			replicas: 0,
			replicaSets: []*apps.ReplicaSet{
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 5, apps.ReplicaSetStatus{})),
			},
			expected: 0,
		},
		{
			name:     "rolling update can not scale up as there are too many pods in the old replica set",
			replicas: 9,
			strategy: apps.DeploymentStrategy{
				Type: apps.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &apps.RollingUpdateDeployment{
					MaxSurge: ptr.To(intstr.FromInt32(1)),
				},
			},
			replicaSets: []*apps.ReplicaSet{
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 5, apps.ReplicaSetStatus{})),
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 5, apps.ReplicaSetStatus{})),
			},
			expected: 5,
		},
		{
			name:     "rolling update can partially scale up as there is some additional capacity in the maxSurge",
			replicas: 6,
			strategy: apps.DeploymentStrategy{
				Type: apps.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &apps.RollingUpdateDeployment{
					MaxSurge: ptr.To(intstr.FromInt32(3)),
				},
			},
			replicaSets: []*apps.ReplicaSet{
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 2, apps.ReplicaSetStatus{})),
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 5, apps.ReplicaSetStatus{})),
			},
			expected: 4,
		},
		{
			name:     "rolling update can scale up as there is enough additional capacity in the maxSurge",
			replicas: 6,
			strategy: apps.DeploymentStrategy{
				Type: apps.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &apps.RollingUpdateDeployment{
					MaxSurge: ptr.To(intstr.FromInt32(10)),
				},
			},
			replicaSets: []*apps.ReplicaSet{
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 2, apps.ReplicaSetStatus{})),
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 5, apps.ReplicaSetStatus{})),
			},
			expected: 6,
		},
		{
			name:     "recreate should fully scale as it is always scaled after the pods from previous replicaset have terminated",
			replicas: 3,
			strategy: apps.DeploymentStrategy{
				Type: apps.RecreateDeploymentStrategyType,
			},
			replicaSets: []*apps.ReplicaSet{
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 1, apps.ReplicaSetStatus{})),
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 5, apps.ReplicaSetStatus{})),
			},
			expected: 3,
		},
	}

	testVariants := []struct {
		name                                 string
		enableDeploymentPodReplacementPolicy bool
		podReplacementPolicy                 *apps.DeploymentPodReplacementPolicy
	}{
		{
			name:                                 " with DeploymentPodReplacementPolicy=false",
			enableDeploymentPodReplacementPolicy: false,
		},

		{
			name:                                 " with DeploymentPodReplacementPolicy=true",
			enableDeploymentPodReplacementPolicy: true,
		},

		{
			name:                                 " with DeploymentPodReplacementPolicy=true and podReplacementPolicy=TerminationStarted",
			enableDeploymentPodReplacementPolicy: true,
			podReplacementPolicy:                 ptr.To(apps.TerminationStarted),
		},
	}

	for _, test := range tests {
		for _, testVariant := range testVariants {
			t.Run(test.name+testVariant.name, func(t *testing.T) {
				featuregatetesting.SetFeatureGateDuringTest(t, utilfeature.DefaultFeatureGate, features.DeploymentPodReplacementPolicy, testVariant.enableDeploymentPodReplacementPolicy)
				newDeployment := generateDeploymentWithReplicas("nginx", test.replicas)
				newDeployment.Spec.Strategy = test.strategy

				if testVariant.enableDeploymentPodReplacementPolicy {
					// TerminationStarted policy should not affect default behavior
					newDeployment.Spec.PodReplacementPolicy = testVariant.podReplacementPolicy
					for _, rs := range test.replicaSets {
						// replicas and terminating replicas should not affect default behavior
						rs.Status.Replicas = rand.Int31n(15)
						rs.Status.TerminatingReplicas = rand.Int31n(15)
					}
				}

				rs, err := NewRSNewReplicas(&newDeployment, test.replicaSets, test.replicaSets[0])
				if err != nil {
					t.Errorf("In test case %s, got unexpected error %v", test.name, err)
				}
				if rs != test.expected {
					t.Errorf("In test case %s, expected %+v, got %+v", test.name, test.expected, rs)
				}
			})
		}
	}
}

func TestNewRSNewReplicasWithTerminationCompleteDeploymentPodReplacementPolicy(t *testing.T) {
	tests := []struct {
		name                           string
		replicas                       int32
		strategy                       apps.DeploymentStrategy
		deploymentPodReplacementPolicy *apps.DeploymentPodReplacementPolicy
		replicaSets                    []*apps.ReplicaSet
		expected                       int32
	}{
		{
			name:     "rolling update with TerminationComplete policy should not scale when old replicas are still present",
			replicas: 5,
			strategy: apps.DeploymentStrategy{
				Type: apps.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &apps.RollingUpdateDeployment{
					MaxSurge: ptr.To(intstr.FromInt32(2)),
				},
			},
			replicaSets: []*apps.ReplicaSet{
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 1, apps.ReplicaSetStatus{Replicas: 1})),
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 3, apps.ReplicaSetStatus{Replicas: 2})),
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 0, apps.ReplicaSetStatus{Replicas: 3})),
			},
			expected: 1,
		},
		{
			name:     "rolling update with TerminationComplete policy should not scale when old and terminating replicas are still present",
			replicas: 5,
			strategy: apps.DeploymentStrategy{
				Type: apps.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &apps.RollingUpdateDeployment{
					MaxSurge: ptr.To(intstr.FromInt32(2)),
				},
			},
			replicaSets: []*apps.ReplicaSet{
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 1, apps.ReplicaSetStatus{TerminatingReplicas: 1})),
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 1, apps.ReplicaSetStatus{TerminatingReplicas: 1})),
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 0, apps.ReplicaSetStatus{Replicas: 1, TerminatingReplicas: 2})),
			},
			expected: 1,
		},
		{
			name:     "rolling update with TerminationComplete policy should not scale when terminating replicas are still present",
			replicas: 5,
			strategy: apps.DeploymentStrategy{
				Type: apps.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &apps.RollingUpdateDeployment{
					MaxSurge: ptr.To(intstr.FromInt32(2)),
				},
			},
			replicaSets: []*apps.ReplicaSet{
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 1, apps.ReplicaSetStatus{TerminatingReplicas: 2})),
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 0, apps.ReplicaSetStatus{TerminatingReplicas: 2})),
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 0, apps.ReplicaSetStatus{TerminatingReplicas: 3})),
			},
			expected: 1,
		},
		{
			name:     "rolling update with TerminationComplete policy should partially scale when terminating pods are present",
			replicas: 6,
			strategy: apps.DeploymentStrategy{
				Type: apps.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &apps.RollingUpdateDeployment{
					MaxSurge: ptr.To(intstr.FromInt32(3)),
				},
			},
			replicaSets: []*apps.ReplicaSet{
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 1, apps.ReplicaSetStatus{Replicas: 1, TerminatingReplicas: 1})),
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 0, apps.ReplicaSetStatus{Replicas: 1, TerminatingReplicas: 1})),
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 0, apps.ReplicaSetStatus{TerminatingReplicas: 1})),
			},
			expected: 5,
		},

		{
			name:     "rolling update with TerminationComplete policy should fully scale when most of the pods are terminated",
			replicas: 6,
			strategy: apps.DeploymentStrategy{
				Type: apps.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &apps.RollingUpdateDeployment{
					MaxSurge: ptr.To(intstr.FromInt32(3)),
				},
			},
			replicaSets: []*apps.ReplicaSet{
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 1, apps.ReplicaSetStatus{Replicas: 1, TerminatingReplicas: 1})),
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 0, apps.ReplicaSetStatus{TerminatingReplicas: 1})),
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 0, apps.ReplicaSetStatus{TerminatingReplicas: 1})),
			},
			expected: 6,
		},
		{
			name:     "recreate with TerminationComplete policy should not scale when old replicas are still present",
			replicas: 7,
			strategy: apps.DeploymentStrategy{
				Type: apps.RecreateDeploymentStrategyType,
			},
			replicaSets: []*apps.ReplicaSet{
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 1, apps.ReplicaSetStatus{Replicas: 1})),
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 3, apps.ReplicaSetStatus{Replicas: 2})),
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 0, apps.ReplicaSetStatus{Replicas: 3})),
			},
			expected: 1,
		},
		{
			name:     "recreate with TerminationComplete policy should not scale when old and terminating replicas are still present",
			replicas: 7,
			strategy: apps.DeploymentStrategy{
				Type: apps.RecreateDeploymentStrategyType,
			},
			replicaSets: []*apps.ReplicaSet{
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 1, apps.ReplicaSetStatus{TerminatingReplicas: 1})),
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 1, apps.ReplicaSetStatus{TerminatingReplicas: 1})),
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 0, apps.ReplicaSetStatus{Replicas: 1, TerminatingReplicas: 2})),
			},
			expected: 1,
		},
		{
			name:     "recreate with TerminationComplete policy should not scale when terminating replicas are still present",
			replicas: 7,
			strategy: apps.DeploymentStrategy{
				Type: apps.RecreateDeploymentStrategyType,
			},
			replicaSets: []*apps.ReplicaSet{
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 1, apps.ReplicaSetStatus{TerminatingReplicas: 2})),
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 0, apps.ReplicaSetStatus{TerminatingReplicas: 2})),
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 0, apps.ReplicaSetStatus{TerminatingReplicas: 3})),
			},
			expected: 1,
		},
		{
			name:     "recreate with TerminationComplete policy should partially scale when terminating pods are present",
			replicas: 3,
			strategy: apps.DeploymentStrategy{
				Type: apps.RecreateDeploymentStrategyType,
			},
			replicaSets: []*apps.ReplicaSet{
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 1, apps.ReplicaSetStatus{Replicas: 1})),
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 0, apps.ReplicaSetStatus{TerminatingReplicas: 1})),
			},
			expected: 2,
		},
		{
			name:     "recreate with TerminationComplete policy should fully scale when all pods terminated",
			replicas: 3,
			strategy: apps.DeploymentStrategy{
				Type: apps.RecreateDeploymentStrategyType,
			},
			replicaSets: []*apps.ReplicaSet{
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 1, apps.ReplicaSetStatus{Replicas: 1})),
				ptr.To(generateRSWithReplicas(generateDeployment("nginx"), 0, apps.ReplicaSetStatus{})),
			},
			expected: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			featuregatetesting.SetFeatureGateDuringTest(t, utilfeature.DefaultFeatureGate, features.DeploymentPodReplacementPolicy, true)
			newDeployment := generateDeploymentWithReplicas("nginx", test.replicas)
			newDeployment.Spec.Strategy = test.strategy
			newDeployment.Spec.PodReplacementPolicy = ptr.To(apps.TerminationComplete)

			rs, err := NewRSNewReplicas(&newDeployment, test.replicaSets, test.replicaSets[0])
			if err != nil {
				t.Errorf("In test case %s, got unexpected error %v", test.name, err)
			}
			if rs != test.expected {
				t.Errorf("In test case %s, expected %+v, got %+v", test.name, test.expected, rs)
			}
		})
	}
}

var (
	condProgressing = func() apps.DeploymentCondition {
		return apps.DeploymentCondition{
			Type:   apps.DeploymentProgressing,
			Status: v1.ConditionFalse,
			Reason: "ForSomeReason",
		}
	}

	condProgressing2 = func() apps.DeploymentCondition {
		return apps.DeploymentCondition{
			Type:   apps.DeploymentProgressing,
			Status: v1.ConditionTrue,
			Reason: "BecauseItIs",
		}
	}

	condAvailable = func() apps.DeploymentCondition {
		return apps.DeploymentCondition{
			Type:   apps.DeploymentAvailable,
			Status: v1.ConditionTrue,
			Reason: "AwesomeController",
		}
	}

	status = func() *apps.DeploymentStatus {
		return &apps.DeploymentStatus{
			Conditions: []apps.DeploymentCondition{condProgressing(), condAvailable()},
		}
	}
)

func TestGetCondition(t *testing.T) {
	exampleStatus := status()

	tests := []struct {
		name string

		status   apps.DeploymentStatus
		condType apps.DeploymentConditionType

		expected bool
	}{
		{
			name: "condition exists",

			status:   *exampleStatus,
			condType: apps.DeploymentAvailable,

			expected: true,
		},
		{
			name: "condition does not exist",

			status:   *exampleStatus,
			condType: apps.DeploymentReplicaFailure,

			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cond := GetDeploymentCondition(test.status, test.condType)
			exists := cond != nil
			if exists != test.expected {
				t.Errorf("%s: expected condition to exist: %t, got: %t", test.name, test.expected, exists)
			}
		})
	}
}

func TestSetCondition(t *testing.T) {
	tests := []struct {
		name string

		status *apps.DeploymentStatus
		cond   apps.DeploymentCondition

		expectedStatus *apps.DeploymentStatus
	}{
		{
			name: "set for the first time",

			status: &apps.DeploymentStatus{},
			cond:   condAvailable(),

			expectedStatus: &apps.DeploymentStatus{Conditions: []apps.DeploymentCondition{condAvailable()}},
		},
		{
			name: "simple set",

			status: &apps.DeploymentStatus{Conditions: []apps.DeploymentCondition{condProgressing()}},
			cond:   condAvailable(),

			expectedStatus: status(),
		},
		{
			name: "overwrite",

			status: &apps.DeploymentStatus{Conditions: []apps.DeploymentCondition{condProgressing()}},
			cond:   condProgressing2(),

			expectedStatus: &apps.DeploymentStatus{Conditions: []apps.DeploymentCondition{condProgressing2()}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			SetDeploymentCondition(test.status, test.cond)
			if !reflect.DeepEqual(test.status, test.expectedStatus) {
				t.Errorf("%s: expected status: %v, got: %v", test.name, test.expectedStatus, test.status)
			}
		})
	}
}

func TestRemoveCondition(t *testing.T) {
	tests := []struct {
		name string

		status   *apps.DeploymentStatus
		condType apps.DeploymentConditionType

		expectedStatus *apps.DeploymentStatus
	}{
		{
			name: "remove from empty status",

			status:   &apps.DeploymentStatus{},
			condType: apps.DeploymentProgressing,

			expectedStatus: &apps.DeploymentStatus{},
		},
		{
			name: "simple remove",

			status:   &apps.DeploymentStatus{Conditions: []apps.DeploymentCondition{condProgressing()}},
			condType: apps.DeploymentProgressing,

			expectedStatus: &apps.DeploymentStatus{},
		},
		{
			name: "doesn't remove anything",

			status:   status(),
			condType: apps.DeploymentReplicaFailure,

			expectedStatus: status(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RemoveDeploymentCondition(test.status, test.condType)
			if !reflect.DeepEqual(test.status, test.expectedStatus) {
				t.Errorf("%s: expected status: %v, got: %v", test.name, test.expectedStatus, test.status)
			}
		})
	}
}

func TestDeploymentComplete(t *testing.T) {
	deployment := func(desired, current, updated, available, maxUnavailable, maxSurge int32) *apps.Deployment {
		return &apps.Deployment{
			Spec: apps.DeploymentSpec{
				Replicas: &desired,
				Strategy: apps.DeploymentStrategy{
					RollingUpdate: &apps.RollingUpdateDeployment{
						MaxUnavailable: ptr.To(intstr.FromInt32(maxUnavailable)),
						MaxSurge:       ptr.To(intstr.FromInt32(maxSurge)),
					},
					Type: apps.RollingUpdateDeploymentStrategyType,
				},
			},
			Status: apps.DeploymentStatus{
				Replicas:          current,
				UpdatedReplicas:   updated,
				AvailableReplicas: available,
			},
		}
	}

	tests := []struct {
		name string

		d *apps.Deployment

		expected bool
	}{
		{
			name: "not complete: min but not all pods become available",

			d:        deployment(5, 5, 5, 4, 1, 0),
			expected: false,
		},
		{
			name: "not complete: min availability is not honored",

			d:        deployment(5, 5, 5, 3, 1, 0),
			expected: false,
		},
		{
			name: "complete",

			d:        deployment(5, 5, 5, 5, 0, 0),
			expected: true,
		},
		{
			name: "not complete: all pods are available but not updated",

			d:        deployment(5, 5, 4, 5, 0, 0),
			expected: false,
		},
		{
			name: "not complete: still running old pods",

			// old replica set: spec.replicas=1, status.replicas=1, status.availableReplicas=1
			// new replica set: spec.replicas=1, status.replicas=1, status.availableReplicas=0
			d:        deployment(1, 2, 1, 1, 0, 1),
			expected: false,
		},
		{
			name: "not complete: one replica deployment never comes up",

			d:        deployment(1, 1, 1, 0, 1, 1),
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, exp := DeploymentComplete(test.d, &test.d.Status), test.expected; got != exp {
				t.Errorf("expected complete: %t, got: %t", exp, got)
			}
		})
	}
}

func TestDeploymentProgressing(t *testing.T) {
	deployment := func(current, updated, ready, available int32) *apps.Deployment {
		return &apps.Deployment{
			Status: apps.DeploymentStatus{
				Replicas:          current,
				UpdatedReplicas:   updated,
				ReadyReplicas:     ready,
				AvailableReplicas: available,
			},
		}
	}
	newStatus := func(current, updated, ready, available int32) apps.DeploymentStatus {
		return apps.DeploymentStatus{
			Replicas:          current,
			UpdatedReplicas:   updated,
			ReadyReplicas:     ready,
			AvailableReplicas: available,
		}
	}

	tests := []struct {
		name string

		d         *apps.Deployment
		newStatus apps.DeploymentStatus

		expected bool
	}{
		{
			name: "progressing: updated pods",

			d:         deployment(10, 4, 4, 4),
			newStatus: newStatus(10, 6, 4, 4),

			expected: true,
		},
		{
			name: "not progressing",

			d:         deployment(10, 4, 4, 4),
			newStatus: newStatus(10, 4, 4, 4),

			expected: false,
		},
		{
			name: "progressing: old pods removed",

			d:         deployment(10, 4, 6, 6),
			newStatus: newStatus(8, 4, 6, 6),

			expected: true,
		},
		{
			name: "not progressing: less new pods",

			d:         deployment(10, 7, 3, 3),
			newStatus: newStatus(10, 6, 3, 3),

			expected: false,
		},
		{
			name: "progressing: less overall but more new pods",

			d:         deployment(10, 4, 7, 7),
			newStatus: newStatus(8, 8, 5, 5),

			expected: true,
		},
		{
			name: "progressing: more ready pods",

			d:         deployment(10, 10, 9, 8),
			newStatus: newStatus(10, 10, 10, 8),

			expected: true,
		},
		{
			name: "progressing: more available pods",

			d:         deployment(10, 10, 10, 9),
			newStatus: newStatus(10, 10, 10, 10),

			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, exp := DeploymentProgressing(test.d, &test.newStatus), test.expected; got != exp {
				t.Errorf("expected progressing: %t, got: %t", exp, got)
			}
		})
	}
}

func TestDeploymentTimedOut(t *testing.T) {
	var (
		null     *int32
		ten      = int32(10)
		infinite = int32(math.MaxInt32)
	)

	timeFn := func(min, sec int) time.Time {
		return time.Date(2016, 1, 1, 0, min, sec, 0, time.UTC)
	}
	deployment := func(condType apps.DeploymentConditionType, status v1.ConditionStatus, reason string, pds *int32, from time.Time) apps.Deployment {
		return apps.Deployment{
			Spec: apps.DeploymentSpec{
				ProgressDeadlineSeconds: pds,
			},
			Status: apps.DeploymentStatus{
				Conditions: []apps.DeploymentCondition{
					{
						Type:           condType,
						Status:         status,
						Reason:         reason,
						LastUpdateTime: metav1.Time{Time: from},
					},
				},
			},
		}
	}

	tests := []struct {
		name string

		d     apps.Deployment
		nowFn func() time.Time

		expected bool
	}{
		{
			name: "nil progressDeadlineSeconds specified - no timeout",

			d:        deployment(apps.DeploymentProgressing, v1.ConditionTrue, "", null, timeFn(1, 9)),
			nowFn:    func() time.Time { return timeFn(1, 20) },
			expected: false,
		},
		{
			name: "infinite progressDeadlineSeconds specified - no timeout",

			d:        deployment(apps.DeploymentProgressing, v1.ConditionTrue, "", &infinite, timeFn(1, 9)),
			nowFn:    func() time.Time { return timeFn(1, 20) },
			expected: false,
		},
		{
			name: "progressDeadlineSeconds: 10s, now - started => 00:01:20 - 00:01:09 => 11s",

			d:        deployment(apps.DeploymentProgressing, v1.ConditionTrue, "", &ten, timeFn(1, 9)),
			nowFn:    func() time.Time { return timeFn(1, 20) },
			expected: true,
		},
		{
			name: "progressDeadlineSeconds: 10s, now - started => 00:01:20 - 00:01:11 => 9s",

			d:        deployment(apps.DeploymentProgressing, v1.ConditionTrue, "", &ten, timeFn(1, 11)),
			nowFn:    func() time.Time { return timeFn(1, 20) },
			expected: false,
		},
		{
			name: "previous status was a complete deployment",

			d:        deployment(apps.DeploymentProgressing, v1.ConditionTrue, NewRSAvailableReason, nil, time.Time{}),
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nowFn = test.nowFn
			_, ctx := ktesting.NewTestContext(t)
			if got, exp := DeploymentTimedOut(ctx, &test.d, &test.d.Status), test.expected; got != exp {
				t.Errorf("expected timeout: %t, got: %t", exp, got)
			}
		})
	}
}

func TestMaxUnavailable(t *testing.T) {
	deployment := func(replicas int32, maxUnavailable intstr.IntOrString) apps.Deployment {
		return apps.Deployment{
			Spec: apps.DeploymentSpec{
				Replicas: func(i int32) *int32 { return &i }(replicas),
				Strategy: apps.DeploymentStrategy{
					RollingUpdate: &apps.RollingUpdateDeployment{
						MaxSurge:       ptr.To(intstr.FromInt32(1)),
						MaxUnavailable: &maxUnavailable,
					},
					Type: apps.RollingUpdateDeploymentStrategyType,
				},
			},
		}
	}
	tests := []struct {
		name       string
		deployment apps.Deployment
		expected   int32
	}{
		{
			name:       "maxUnavailable less than replicas",
			deployment: deployment(10, intstr.FromInt32(5)),
			expected:   int32(5),
		},
		{
			name:       "maxUnavailable equal replicas",
			deployment: deployment(10, intstr.FromInt32(10)),
			expected:   int32(10),
		},
		{
			name:       "maxUnavailable greater than replicas",
			deployment: deployment(5, intstr.FromInt32(10)),
			expected:   int32(5),
		},
		{
			name:       "maxUnavailable with replicas is 0",
			deployment: deployment(0, intstr.FromInt32(10)),
			expected:   int32(0),
		},
		{
			name: "maxUnavailable with Recreate deployment strategy",
			deployment: apps.Deployment{
				Spec: apps.DeploymentSpec{
					Strategy: apps.DeploymentStrategy{
						Type: apps.RecreateDeploymentStrategyType,
					},
				},
			},
			expected: int32(0),
		},
		{
			name:       "maxUnavailable less than replicas with percents",
			deployment: deployment(10, intstr.FromString("50%")),
			expected:   int32(5),
		},
		{
			name:       "maxUnavailable equal replicas with percents",
			deployment: deployment(10, intstr.FromString("100%")),
			expected:   int32(10),
		},
		{
			name:       "maxUnavailable greater than replicas with percents",
			deployment: deployment(5, intstr.FromString("100%")),
			expected:   int32(5),
		},
	}

	for _, test := range tests {
		t.Log(test.name)
		t.Run(test.name, func(t *testing.T) {
			maxUnavailable := MaxUnavailable(test.deployment)
			if test.expected != maxUnavailable {
				t.Fatalf("expected:%v, got:%v", test.expected, maxUnavailable)
			}
		})
	}
}

func TestGetNonNegativeIntFromAnnotation(t *testing.T) {
	tests := []struct {
		name          string
		annotations   map[string]string
		expectedValue int32
		expectedValid bool
		expectErr     bool
	}{
		{
			name: "invalid empty",
		},
		{
			name:        "invalid negative ",
			annotations: map[string]string{"test": "-1", "foo": "2"},
			expectErr:   true,
		},
		{
			name:        "invalid",
			annotations: map[string]string{"test": "invalid", "foo": "2"},
			expectErr:   true,
		},
		{
			name:          "valid",
			annotations:   map[string]string{"test": "13", "foo": "2"},
			expectedValue: 13,
			expectedValid: true,
		},
		{
			name:          "valid zero",
			annotations:   map[string]string{"test": "0", "foo": "2"},
			expectedValue: 0,
			expectedValid: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tDeployment := generateDeployment("nginx")
			tRS := generateRS(tDeployment)
			tRS.Annotations = test.annotations
			value, valid, err := getNonNegativeIntFromAnnotation(&tRS, "test")
			if test.expectedValue != value {
				t.Fatalf("expected value:%v, got:%v", test.expectedValue, value)
			}
			if test.expectedValid != valid {
				t.Fatalf("expected valid:%v, got:%v", test.expectedValid, valid)
			}
			if test.expectErr != (err != nil) {
				t.Fatalf("expected err:%v, got:%v", test.expectErr, err)
			}
		})
	}
}

// Set of simple tests for annotation related util functions
func TestAnnotationUtils(t *testing.T) {

	//Setup
	tDeployment := generateDeploymentWithReplicas("nginx", 10)
	tDeployment.Spec.Strategy = apps.DeploymentStrategy{
		RollingUpdate: &apps.RollingUpdateDeployment{
			MaxSurge:       ptr.To(intstr.FromInt32(1)),
			MaxUnavailable: ptr.To(intstr.FromInt32(2)),
		},
		Type: apps.RollingUpdateDeploymentStrategyType,
	}
	tRS := generateRS(tDeployment)
	tDeployment.Annotations[RevisionAnnotation] = "1"

	//Test Case 1: Check if anotations are copied properly from deployment to RS
	t.Run("SetNewReplicaSetAnnotations", func(t *testing.T) {
		_, ctx := ktesting.NewTestContext(t)

		//Try to set the increment revision from 11 through 20
		for i := 10; i < 20; i++ {

			nextRevision := fmt.Sprintf("%d", i+1)
			SetNewReplicaSetAnnotations(ctx, &tDeployment, &tRS, nextRevision, true, 5)
			//Now the ReplicaSets Revision Annotation should be i+1

			if i >= 12 {
				expectedHistoryAnnotation := fmt.Sprintf("%d,%d", i-1, i)
				if tRS.Annotations[RevisionHistoryAnnotation] != expectedHistoryAnnotation {
					t.Errorf("Revision History Expected=%s Obtained=%s", expectedHistoryAnnotation, tRS.Annotations[RevisionHistoryAnnotation])
				}
			}
			if tRS.Annotations[RevisionAnnotation] != nextRevision {
				t.Errorf("Revision Expected=%s Obtained=%s", nextRevision, tRS.Annotations[RevisionAnnotation])
			}
		}
	})

	//Test Case 2:  Check if annotations are set properly
	t.Run("SetReplicaSetScaleAnnotations", func(t *testing.T) {
		annotationUpdate, needUpdate := ComputeReplicaSetScaleAnnotations(&tRS, &tDeployment, false)
		if !needUpdate {
			t.Errorf("ComputeReplicaSetScaleAnnotations() failed")
		}
		SetReplicaSetScaleAnnotations(&tRS, annotationUpdate)
		value, ok := tRS.Annotations[DesiredReplicasAnnotation]
		if !ok {
			t.Errorf("ComputeReplicaSetScaleAnnotations / SetReplicaSetScaleAnnotations did not set DesiredReplicasAnnotation")
		}
		if value != "10" {
			t.Errorf("ComputeReplicaSetScaleAnnotations / SetReplicaSetScaleAnnotations did not set DesiredReplicasAnnotation correctly value=%s", value)
		}
		if value, ok = tRS.Annotations[MaxReplicasAnnotation]; !ok {
			t.Errorf("ComputeReplicaSetScaleAnnotations / SetReplicaSetScaleAnnotations did not set DesiredReplicasAnnotation")
		}
		if value != "11" {
			t.Errorf("ComputeReplicaSetScaleAnnotations / SetReplicaSetScaleAnnotations did not set MaxReplicasAnnotation correctly value=%s", value)
		}
	})

	//Test Case 3:  Check if annotations reflect deployments state
	tRS.Annotations[DesiredReplicasAnnotation] = "10"
	tRS.Status.AvailableReplicas = 10
	tRS.Spec.Replicas = new(int32)
	*tRS.Spec.Replicas = 10

	t.Run("IsSaturated", func(t *testing.T) {
		saturated := IsSaturated(&tDeployment, &tRS)
		if !saturated {
			t.Errorf("SetReplicasAnnotations Expected=true Obtained=false")
		}
	})
	//Tear Down
}

func TestComputeAndSetReplicaSetScaleAnnotations(t *testing.T) {
	desiredReplicas := fmt.Sprintf("%d", int32(13))
	maxReplicas := fmt.Sprintf("%d", int32(24))

	tests := []struct {
		name                           string
		replicaSet                     *apps.ReplicaSet
		partialScaling                 bool
		requiresPodReplacementPolicyFG bool
		expectedNeedUpdate             bool
		expectedAnnotations            map[string]string
	}{
		{
			name: "test Annotations nil",
			replicaSet: &apps.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{Name: "hello", Namespace: "test"},
				Spec: apps.ReplicaSetSpec{
					Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"foo": "bar"}},
				},
			},
			expectedAnnotations: map[string]string{DesiredReplicasAnnotation: desiredReplicas, MaxReplicasAnnotation: maxReplicas},
			expectedNeedUpdate:  true,
		},
		{
			name: "test desiredReplicas update",
			replicaSet: &apps.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "hello",
					Namespace:   "test",
					Annotations: map[string]string{DesiredReplicasAnnotation: "8", MaxReplicasAnnotation: maxReplicas, "foo": "bar"},
				},
				Spec: apps.ReplicaSetSpec{
					Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"foo": "bar"}},
				},
			},
			expectedNeedUpdate:  true,
			expectedAnnotations: map[string]string{DesiredReplicasAnnotation: desiredReplicas, MaxReplicasAnnotation: maxReplicas, "foo": "bar"},
		},
		{
			name: "test maxReplicas update",
			replicaSet: &apps.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "hello",
					Namespace:   "test",
					Annotations: map[string]string{DesiredReplicasAnnotation: desiredReplicas, MaxReplicasAnnotation: "16"},
				},
				Spec: apps.ReplicaSetSpec{
					Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"foo": "bar"}},
				},
			},
			expectedAnnotations: map[string]string{DesiredReplicasAnnotation: desiredReplicas, MaxReplicasAnnotation: maxReplicas},
			expectedNeedUpdate:  true,
		},
		{
			name: "test needn't update",
			replicaSet: &apps.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "hello",
					Namespace:   "test",
					Annotations: map[string]string{DesiredReplicasAnnotation: desiredReplicas, MaxReplicasAnnotation: maxReplicas},
				},
				Spec: apps.ReplicaSetSpec{
					Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"foo": "bar"}},
				},
			},
			expectedAnnotations: map[string]string{DesiredReplicasAnnotation: desiredReplicas, MaxReplicasAnnotation: maxReplicas},
			expectedNeedUpdate:  false,
		},
		{
			name: "test removes replicasBeforeScale when not partial scaling",
			replicaSet: &apps.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "hello",
					Namespace:   "test",
					Annotations: map[string]string{DesiredReplicasAnnotation: desiredReplicas, MaxReplicasAnnotation: maxReplicas, ReplicaSetReplicasBeforeScaleAnnotation: "5"},
				},
				Spec: apps.ReplicaSetSpec{
					Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"foo": "bar"}},
				},
			},
			partialScaling:      false,
			expectedAnnotations: map[string]string{DesiredReplicasAnnotation: desiredReplicas, MaxReplicasAnnotation: maxReplicas},
			expectedNeedUpdate:  true,
		},
		{
			name: "test starts partial scaling and keeps old desiredReplicas and maxReplicas",
			replicaSet: &apps.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "hello",
					Namespace:   "test",
					Annotations: map[string]string{DesiredReplicasAnnotation: "7", MaxReplicasAnnotation: "18"},
				},
				Spec: apps.ReplicaSetSpec{
					Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"foo": "bar"}},
					Replicas: ptr.To(int32(5)),
				},
			},
			partialScaling:                 true,
			requiresPodReplacementPolicyFG: true,
			expectedAnnotations:            map[string]string{DesiredReplicasAnnotation: "7", MaxReplicasAnnotation: "18", ReplicaSetReplicasBeforeScaleAnnotation: "5"},
			expectedNeedUpdate:             true,
		},
		{
			name: "test starts partial scaling with invalid ReplicaSetReplicasBeforeScale value, keeps old desiredReplicas and maxReplicas",
			replicaSet: &apps.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "hello",
					Namespace:   "test",
					Annotations: map[string]string{DesiredReplicasAnnotation: "7", MaxReplicasAnnotation: "18", ReplicaSetReplicasBeforeScaleAnnotation: "invalid"},
				},
				Spec: apps.ReplicaSetSpec{
					Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"foo": "bar"}},
					Replicas: ptr.To(int32(5)),
				},
			},
			partialScaling:                 true,
			requiresPodReplacementPolicyFG: true,
			expectedAnnotations:            map[string]string{DesiredReplicasAnnotation: "7", MaxReplicasAnnotation: "18", ReplicaSetReplicasBeforeScaleAnnotation: "5"},
			expectedNeedUpdate:             true,
		},
		{
			name: "test continues partial scaling and keeps old desiredReplicas and maxReplicas and replicasBeforeScale",
			replicaSet: &apps.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "hello",
					Namespace:   "test",
					Annotations: map[string]string{DesiredReplicasAnnotation: "7", MaxReplicasAnnotation: "18", ReplicaSetReplicasBeforeScaleAnnotation: "5"},
				},
				Spec: apps.ReplicaSetSpec{
					Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"foo": "bar"}},
					Replicas: ptr.To(int32(6)),
				},
			},
			partialScaling:                 true,
			requiresPodReplacementPolicyFG: true,
			expectedAnnotations:            map[string]string{DesiredReplicasAnnotation: "7", MaxReplicasAnnotation: "18", ReplicaSetReplicasBeforeScaleAnnotation: "5"},
			expectedNeedUpdate:             false,
		},
		{
			name: "test removes updates desiredReplicas and maxReplicas and replicasBeforeScale when finishing partial scaling",
			replicaSet: &apps.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "hello",
					Namespace:   "test",
					Annotations: map[string]string{DesiredReplicasAnnotation: "7", MaxReplicasAnnotation: "18", ReplicaSetReplicasBeforeScaleAnnotation: "5"},
				},
				Spec: apps.ReplicaSetSpec{
					Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"foo": "bar"}},
					Replicas: ptr.To(int32(10)),
				},
			},
			partialScaling:      false,
			expectedAnnotations: map[string]string{DesiredReplicasAnnotation: desiredReplicas, MaxReplicasAnnotation: maxReplicas},
			expectedNeedUpdate:  true,
		},
	}

	for _, podReplacementPolicyFG := range []bool{true, false} {
		for _, test := range tests {
			if test.requiresPodReplacementPolicyFG && !podReplacementPolicyFG {
				continue
			}
			t.Run(test.name+fmt.Sprintf("with podReplacementPolicyFG=%v", podReplacementPolicyFG), func(t *testing.T) {
				featuregatetesting.SetFeatureGateDuringTest(t, utilfeature.DefaultFeatureGate, features.DeploymentPodReplacementPolicy, podReplacementPolicyFG)
				tDeployment := generateDeploymentWithReplicas("hello", 13)
				tDeployment.Spec.Strategy = apps.DeploymentStrategy{
					Type: apps.RollingUpdateDeploymentStrategyType,
					RollingUpdate: &apps.RollingUpdateDeployment{
						MaxSurge:       ptr.To(intstr.FromInt32(11)),
						MaxUnavailable: ptr.To(intstr.FromInt32(2)),
					},
				}
				annotationUpdate, needUpdate := ComputeReplicaSetScaleAnnotations(test.replicaSet, &tDeployment, test.partialScaling)
				tReplicaSet := test.replicaSet.DeepCopy()
				SetReplicaSetScaleAnnotations(tReplicaSet, annotationUpdate)
				if needUpdate != test.expectedNeedUpdate {
					t.Errorf("Expected needUpdate %v, Got: %v", test.expectedNeedUpdate, needUpdate)
				}
				if !reflect.DeepEqual(test.expectedAnnotations, tReplicaSet.Annotations) {
					t.Errorf("invalid annotations: %v", cmp.Diff(test.expectedAnnotations, tReplicaSet.Annotations))
				}
			})
		}
	}

}

func TestGetDeploymentsForReplicaSet(t *testing.T) {
	fakeInformerFactory := informers.NewSharedInformerFactory(&fake.Clientset{}, 0*time.Second)
	var deployments []*apps.Deployment
	for i := 0; i < 3; i++ {
		deployment := &apps.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("deployment-%d", i),
				Namespace: "test",
			},
			Spec: apps.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": fmt.Sprintf("test-%d", i),
					},
				},
			},
		}
		deployments = append(deployments, deployment)
		fakeInformerFactory.Apps().V1().Deployments().Informer().GetStore().Add(deployment)
	}
	var rss []*apps.ReplicaSet
	for i := 0; i < 5; i++ {
		rs := &apps.ReplicaSet{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "test",
				Name:      fmt.Sprintf("test-replicaSet-%d", i),
				Labels: map[string]string{
					"app":   fmt.Sprintf("test-%d", i),
					"label": fmt.Sprintf("label-%d", i),
				},
			},
		}
		rss = append(rss, rs)
	}
	tests := []struct {
		name   string
		rs     *apps.ReplicaSet
		err    error
		expect []*apps.Deployment
	}{
		{
			name:   "GetDeploymentsForReplicaSet for rs-0",
			rs:     rss[0],
			expect: []*apps.Deployment{deployments[0]},
		},
		{
			name:   "GetDeploymentsForReplicaSet for rs-1",
			rs:     rss[1],
			expect: []*apps.Deployment{deployments[1]},
		},
		{
			name:   "GetDeploymentsForReplicaSet for rs-2",
			rs:     rss[2],
			expect: []*apps.Deployment{deployments[2]},
		},
		{
			name: "GetDeploymentsForReplicaSet for rs-3",
			rs:   rss[3],
			err:  fmt.Errorf("could not find deployments set for ReplicaSet %s in namespace %s with labels: %v", rss[3].Name, rss[3].Namespace, rss[3].Labels),
		},
		{
			name: "GetDeploymentsForReplicaSet for rs-4",
			rs:   rss[4],
			err:  fmt.Errorf("could not find deployments set for ReplicaSet %s in namespace %s with labels: %v", rss[4].Name, rss[4].Namespace, rss[4].Labels),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			get, err := GetDeploymentsForReplicaSet(fakeInformerFactory.Apps().V1().Deployments().Lister(), test.rs)
			if err != nil {
				if err.Error() != test.err.Error() {
					t.Errorf("Error from GetDeploymentsForReplicaSet: %v", err)
				}
			} else if !reflect.DeepEqual(get, test.expect) {
				t.Errorf("Expect deployments %v, but got %v", test.expect, get)
			}
		})
	}

}

func TestMinAvailable(t *testing.T) {
	maxSurge := ptr.To(intstr.FromInt32(1))
	deployment := func(replicas int32, maxUnavailable intstr.IntOrString) *apps.Deployment {
		return &apps.Deployment{
			Spec: apps.DeploymentSpec{
				Replicas: ptr.To(replicas),
				Strategy: apps.DeploymentStrategy{
					RollingUpdate: &apps.RollingUpdateDeployment{
						MaxSurge:       maxSurge,
						MaxUnavailable: &maxUnavailable,
					},
					Type: apps.RollingUpdateDeploymentStrategyType,
				},
			},
		}
	}
	tests := []struct {
		name       string
		deployment *apps.Deployment
		expected   int32
	}{
		{
			name:       "replicas greater than maxUnavailable",
			deployment: deployment(10, intstr.FromInt32(5)),
			expected:   5,
		},
		{
			name:       "replicas equal maxUnavailable",
			deployment: deployment(10, intstr.FromInt32(10)),
			expected:   0,
		},
		{
			name:       "replicas less than maxUnavailable",
			deployment: deployment(5, intstr.FromInt32(10)),
			expected:   0,
		},
		{
			name:       "replicas is 0",
			deployment: deployment(0, intstr.FromInt32(10)),
			expected:   0,
		},
		{
			name: "minAvailable with Recreate deployment strategy",
			deployment: &apps.Deployment{
				Spec: apps.DeploymentSpec{
					Replicas: ptr.To[int32](10),
					Strategy: apps.DeploymentStrategy{
						Type: apps.RecreateDeploymentStrategyType,
					},
				},
			},
			expected: 0,
		},
		{
			name:       "replicas greater than maxUnavailable with percents",
			deployment: deployment(10, intstr.FromString("60%")),
			expected:   4,
		},
		{
			name:       "replicas equal maxUnavailable with percents",
			deployment: deployment(10, intstr.FromString("100%")),
			expected:   int32(0),
		},
		{
			name:       "replicas less than maxUnavailable with percents",
			deployment: deployment(5, intstr.FromString("100%")),
			expected:   0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MinAvailable(tt.deployment); got != tt.expected {
				t.Errorf("MinAvailable() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetReplicaSetFraction(t *testing.T) {
	tests := []struct {
		name                                 string
		enableDeploymentPodReplacementPolicy bool
		deploymentReplicas                   int32
		deploymentStatusReplicas             int32
		deploymentMaxSurge                   int32
		rsReplicas                           int32
		rsAnnotations                        map[string]string
		expectedFraction                     int32
	}{
		{
			name:               "empty deployment always scales to 0",
			deploymentReplicas: 0,
			rsReplicas:         10,
			expectedFraction:   -10,
		},
		{
			name:               "unsynced deployment does not scale when max-replicas annotation is missing",
			deploymentReplicas: 10,
			rsReplicas:         5,
			expectedFraction:   0,
		},
		{
			name:                     "scale up by 1/5 should increase RS replicas by 1/5 when max-replicas annotation is missing",
			deploymentReplicas:       120,
			deploymentStatusReplicas: 100,
			rsReplicas:               50,
			expectedFraction:         10,
		},
		{
			name:               "scale up by 1/5 should increase RS replicas by 1/5",
			deploymentReplicas: 120,
			rsReplicas:         50,
			rsAnnotations: map[string]string{
				MaxReplicasAnnotation: "100",
			},
			expectedFraction: 10,
		},
		{
			name:               "scale up with maxSurge by 1/5 should increase RS replicas approximately by 1/5",
			deploymentReplicas: 120,
			deploymentMaxSurge: 10,
			rsReplicas:         50,
			rsAnnotations: map[string]string{
				MaxReplicasAnnotation: "110",
			},
			expectedFraction: 9, // expectedFraction is not the whole 1/5 (10) since maxSurge pods have to be taken into account
			// and replica sets with these surge pods should proportionally scale as well during a rollout
		},
		{
			name:               "scale down by 1/6 should decrease RS replicas by 1/6",
			deploymentReplicas: 10,
			rsReplicas:         6,
			rsAnnotations: map[string]string{
				MaxReplicasAnnotation: "12",
			},
			expectedFraction: -1,
		},
		{
			name:               "scale down with maxSurge by 1/6 should decrease RS replicas approximately by 1/6",
			deploymentReplicas: 100,
			deploymentMaxSurge: 10,
			rsReplicas:         50,
			rsAnnotations: map[string]string{
				MaxReplicasAnnotation: "130",
			},
			expectedFraction: -8,
		},
		{
			name:                                 "scale up by 1/5 should increase RS replicas by 1/5 even when replicaset-replicas-before-scale is present and DeploymentPodReplacementPolicy is disabled",
			enableDeploymentPodReplacementPolicy: false,
			deploymentReplicas:                   120,
			rsReplicas:                           55,
			rsAnnotations: map[string]string{
				MaxReplicasAnnotation:                   "100",
				ReplicaSetReplicasBeforeScaleAnnotation: "50",
			},
			expectedFraction: 11,
		},
		{
			name:                                 "partial scale up by 1/5 should increase RS replicas by 1/5 in the end when DeploymentPodReplacementPolicy is enabled",
			enableDeploymentPodReplacementPolicy: true,
			deploymentReplicas:                   120,
			rsReplicas:                           55,
			rsAnnotations: map[string]string{
				MaxReplicasAnnotation:                   "100",
				ReplicaSetReplicasBeforeScaleAnnotation: "50",
			},
			expectedFraction: 5,
		},
		{
			name:                                 "partial scale up with maxSurge by 1/5 should increase RS replicas approximately by 1/5 in the end when DeploymentPodReplacementPolicy is enabled",
			enableDeploymentPodReplacementPolicy: true,
			deploymentReplicas:                   120,
			deploymentMaxSurge:                   10,
			rsReplicas:                           54,
			rsAnnotations: map[string]string{
				MaxReplicasAnnotation:                   "110",
				ReplicaSetReplicasBeforeScaleAnnotation: "50",
			},
			expectedFraction: 5,
		},
		{
			name:                                 "partial scale down with maxSurge by 1/6 should decrease RS replicas approximately by 1/6 in the end when DeploymentPodReplacementPolicy is enabled",
			enableDeploymentPodReplacementPolicy: true,
			deploymentReplicas:                   100,
			deploymentMaxSurge:                   10,
			rsReplicas:                           47,
			rsAnnotations: map[string]string{
				MaxReplicasAnnotation:                   "130",
				ReplicaSetReplicasBeforeScaleAnnotation: "50",
			},
			expectedFraction: -5,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			featuregatetesting.SetFeatureGateDuringTest(t, utilfeature.DefaultFeatureGate, features.DeploymentPodReplacementPolicy, test.enableDeploymentPodReplacementPolicy)
			logger, _ := ktesting.NewTestContext(t)

			tDeployment := generateDeploymentWithReplicas("nginx", test.deploymentReplicas)
			tDeployment.Status.Replicas = test.deploymentStatusReplicas
			tDeployment.Spec.Strategy = apps.DeploymentStrategy{
				Type: apps.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &apps.RollingUpdateDeployment{
					MaxSurge:       ptr.To(intstr.FromInt32(test.deploymentMaxSurge)),
					MaxUnavailable: ptr.To(intstr.FromInt32(1)),
				},
			}

			tRS := generateRSWithReplicas(tDeployment, test.rsReplicas, apps.ReplicaSetStatus{})
			tRS.Annotations = test.rsAnnotations

			fraction := getReplicaSetFraction(logger, tRS, tDeployment)
			if test.expectedFraction != fraction {
				t.Fatalf("expected fraction: %v, got:%v", test.expectedFraction, fraction)
			}
		})
	}
}

func TestGetReplicaSetProportion(t *testing.T) {
	tests := []struct {
		name               string
		deploymentReplicas int32
		rsReplicas         int32
		rsAnnotations      map[string]string
		replicasToAdd      int32
		replicasAdded      int32
		expectedProportion int32
	}{
		{
			name:               "empty RS always scales to 0",
			deploymentReplicas: 10,
			rsReplicas:         0,
			rsAnnotations: map[string]string{
				MaxReplicasAnnotation: "5",
			},
			expectedProportion: 0,
		},
		{
			name:               "replicasToAdd=0 always scales to 0",
			deploymentReplicas: 10,
			rsReplicas:         5,
			rsAnnotations: map[string]string{
				MaxReplicasAnnotation: "5",
			},
			expectedProportion: 0,
		},
		{
			name:               "replicasToAdd=replicasAdded always scales to 0",
			deploymentReplicas: 10,
			rsReplicas:         5,
			rsAnnotations: map[string]string{
				MaxReplicasAnnotation: "5",
			},
			replicasToAdd:      5,
			replicasAdded:      5,
			expectedProportion: 0,
		},
		{
			name:               "scale up by 1/5 should fully increase RS replicas when there is enough capacity in replicasAdded",
			deploymentReplicas: 120,
			rsReplicas:         50,
			rsAnnotations: map[string]string{
				MaxReplicasAnnotation: "100",
			},
			replicasToAdd:      20,
			replicasAdded:      5,
			expectedProportion: 10,
		},

		{
			name:               "scale up by 1/5 should partially increase RS replicas when there is not enough capacity in replicasAdded",
			deploymentReplicas: 120,
			rsReplicas:         50,
			rsAnnotations: map[string]string{
				MaxReplicasAnnotation: "100",
			},
			replicasToAdd:      20,
			replicasAdded:      13,
			expectedProportion: 7,
		},
		{
			name:               "scale down by 1/6 should fully decrease RS replicas when there is enough capacity in replicasAdded",
			deploymentReplicas: 100,
			rsReplicas:         60,
			rsAnnotations: map[string]string{
				MaxReplicasAnnotation: "120",
			},
			replicasToAdd:      -20,
			replicasAdded:      -10,
			expectedProportion: -10,
		},
		{
			name:               "scale down by 1/6 should partially decrease RS replicas when there is not enough capacity in replicasAdded",
			deploymentReplicas: 100,
			rsReplicas:         60,
			rsAnnotations: map[string]string{
				MaxReplicasAnnotation: "120",
			},
			replicasToAdd:      -20,
			replicasAdded:      -19,
			expectedProportion: -1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logger, _ := ktesting.NewTestContext(t)

			tDeployment := generateDeploymentWithReplicas("nginx", test.deploymentReplicas)
			tRS := generateRSWithReplicas(tDeployment, test.rsReplicas, apps.ReplicaSetStatus{})
			tRS.Annotations = test.rsAnnotations

			proportion := getReplicaSetProportion(logger, &tRS, tDeployment, test.replicasToAdd, test.replicasAdded)
			if test.expectedProportion != proportion {
				t.Fatalf("expected proportion: %v, got:%v", test.expectedProportion, proportion)
			}
		})
	}
}

func TestPrepareProportionalScalePlan(t *testing.T) {
	tests := []struct {
		name                                 string
		enableDeploymentPodReplacementPolicy bool
		deploymentReplicas                   int32
		deploymentMaxSurge                   int32
		deploymentScaleUpdateReplicas        int32
		replicaSetsReplicas                  map[string]int32
		replicasetReplicasBeforeScale        map[string]int32
		terminatingReplica                   int32
		expectedScalePlan                    map[string]int32
	}{
		{
			name:                          "should scale up proportionally",
			deploymentReplicas:            100,
			deploymentScaleUpdateReplicas: 120,
			replicaSetsReplicas: map[string]int32{
				"rs-3": 50,
				"rs-2": 30,
				"rs-1": 20,
			},
			expectedScalePlan: map[string]int32{
				"rs-3": 60,
				"rs-2": 36,
				"rs-1": 24,
			},
		},
		{
			name:                          "should scale down proportionally",
			deploymentReplicas:            120,
			deploymentScaleUpdateReplicas: 100,
			replicaSetsReplicas: map[string]int32{
				"rs-3": 60,
				"rs-2": 36,
				"rs-1": 24,
			},
			expectedScalePlan: map[string]int32{
				"rs-3": 50,
				"rs-2": 30,
				"rs-1": 20,
			},
		},
		{
			name:                          "the largest replicaset rs-2 should be preferred with a leftover when scaling up proportionally",
			deploymentReplicas:            100,
			deploymentScaleUpdateReplicas: 111,
			replicaSetsReplicas: map[string]int32{
				"rs-3": 30,
				"rs-2": 40,
				"rs-1": 30,
			},
			expectedScalePlan: map[string]int32{
				"rs-3": 33,
				"rs-2": 45,
				"rs-1": 33,
			},
		},
		{
			name:                          "the newest and largest replicaset rs-3 should be preferred with a leftover when scaling up proportionally",
			deploymentReplicas:            100,
			deploymentScaleUpdateReplicas: 111,
			replicaSetsReplicas: map[string]int32{
				"rs-3": 40,
				"rs-2": 40,
				"rs-1": 20,
			},
			expectedScalePlan: map[string]int32{
				"rs-3": 45,
				"rs-2": 44,
				"rs-1": 22,
			},
		},
		{
			name:                          "the oldest and largest replicaset rs-3 should be preferred with a leftover when scaling down proportionally",
			deploymentReplicas:            111,
			deploymentScaleUpdateReplicas: 100,
			replicaSetsReplicas: map[string]int32{
				"rs-3": 44,
				"rs-2": 44,
				"rs-1": 23,
			},
			expectedScalePlan: map[string]int32{
				"rs-3": 40,
				"rs-2": 39,
				"rs-1": 21,
			},
		},
		{
			name:                          "should scale up proportionally with maxSurge",
			deploymentReplicas:            100,
			deploymentMaxSurge:            10,
			deploymentScaleUpdateReplicas: 120,
			replicaSetsReplicas: map[string]int32{
				"rs-3": 60, // latest RS is surging
				"rs-2": 30,
				"rs-1": 20,
			},
			expectedScalePlan: map[string]int32{
				"rs-3": 71, // latest RS is surging
				"rs-2": 35,
				"rs-1": 24,
			},
		},
		{
			name:                          "should scale up proportionally with maxSurge: second iteration",
			deploymentReplicas:            120,
			deploymentMaxSurge:            10,
			deploymentScaleUpdateReplicas: 130,
			replicaSetsReplicas: map[string]int32{
				"rs-3": 71, // latest RS is surging
				"rs-2": 35,
				"rs-1": 24,
			},
			expectedScalePlan: map[string]int32{
				"rs-3": 76, // latest RS is surging
				"rs-2": 38,
				"rs-1": 26,
			},
		},
		{
			name:                          "should scale down proportionally with maxSurge",
			deploymentReplicas:            120,
			deploymentMaxSurge:            10,
			deploymentScaleUpdateReplicas: 100,
			replicaSetsReplicas: map[string]int32{
				"rs-3": 71, // latest RS is surging
				"rs-2": 35,
				"rs-1": 24,
			},
			expectedScalePlan: map[string]int32{
				"rs-3": 60, // latest RS is surging
				"rs-2": 30,
				"rs-1": 20,
			},
		},
		{
			name:                                 "should scale up partially and proportionally with maxSurge: first iteration",
			enableDeploymentPodReplacementPolicy: true,
			deploymentReplicas:                   100,
			deploymentMaxSurge:                   10,
			deploymentScaleUpdateReplicas:        120,
			replicaSetsReplicas: map[string]int32{
				"rs-3": 50, // latest RS is surging
				"rs-2": 30,
				"rs-1": 20,
			},
			terminatingReplica: 15,
			expectedScalePlan: map[string]int32{
				"rs-3": 59, // latest RS is surging
				"rs-2": 35,
				"rs-1": 21,
			},
		},
		{
			name:                                 "should scale up partially and proportionally with maxSurge: second iteration",
			enableDeploymentPodReplacementPolicy: true,
			deploymentReplicas:                   100,
			deploymentMaxSurge:                   10,
			deploymentScaleUpdateReplicas:        120,
			replicaSetsReplicas: map[string]int32{
				"rs-3": 59, // latest RS is surging
				"rs-2": 35,
				"rs-1": 21,
			},
			replicasetReplicasBeforeScale: map[string]int32{
				"rs-3": 50, // latest RS is surging
				"rs-2": 30,
				"rs-1": 20,
			},
			terminatingReplica: 5,
			expectedScalePlan: map[string]int32{
				"rs-3": 66, // latest RS is surging
				"rs-2": 35,
				"rs-1": 24,
			},
		},
		{
			name:                                 "should scale up partially and proportionally with maxSurge: third iteration",
			enableDeploymentPodReplacementPolicy: true,
			deploymentReplicas:                   100,
			deploymentMaxSurge:                   10,
			deploymentScaleUpdateReplicas:        120,
			replicaSetsReplicas: map[string]int32{
				"rs-3": 66, // latest RS is surging
				"rs-2": 35,
				"rs-1": 24,
			},
			replicasetReplicasBeforeScale: map[string]int32{
				"rs-3": 50, // latest RS is surging
				"rs-2": 30,
				"rs-1": 20,
			},
			terminatingReplica: 0,
			expectedScalePlan: map[string]int32{
				"rs-3": 71, // latest RS is surging
				"rs-2": 35,
				"rs-1": 24,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			featuregatetesting.SetFeatureGateDuringTest(t, utilfeature.DefaultFeatureGate, features.DeploymentPodReplacementPolicy, test.enableDeploymentPodReplacementPolicy)
			logger, _ := ktesting.NewTestContext(t)

			tDeployment := generateDeploymentWithReplicas("nginx", test.deploymentReplicas)
			tDeployment.Spec.Strategy = apps.DeploymentStrategy{
				Type: apps.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &apps.RollingUpdateDeployment{
					MaxSurge:       ptr.To(intstr.FromInt32(test.deploymentMaxSurge)),
					MaxUnavailable: ptr.To(intstr.FromInt32(1)),
				},
			}
			maxReplicasAnnotationValue := fmt.Sprintf("%v", *tDeployment.Spec.Replicas+test.deploymentMaxSurge)

			allActiveRSs := []*apps.ReplicaSet{}
			allReplicas := int32(0)
			for rsName, rsReplicas := range test.replicaSetsReplicas {
				tRS := generateRSWithReplicas(tDeployment, rsReplicas, apps.ReplicaSetStatus{})
				tRS.Name = rsName
				tRS.Annotations = map[string]string{
					MaxReplicasAnnotation: maxReplicasAnnotationValue,
				}
				if replicasBeforeScale := test.replicasetReplicasBeforeScale[rsName]; replicasBeforeScale > 0 {
					tRS.Annotations[ReplicaSetReplicasBeforeScaleAnnotation] = fmt.Sprintf("%v", replicasBeforeScale)
				}
				minutesCreatedAfterDeployment, err := strconv.Atoi(rsName[3:])
				if err != nil {
					t.Errorf("unexpected error %v", err)
				}
				tRS.CreationTimestamp = metav1.Time{Time: time.Now().Add(time.Duration(minutesCreatedAfterDeployment) * time.Minute)}
				allActiveRSs = append(allActiveRSs, &tRS)
				allReplicas += rsReplicas
			}

			// simulate scale
			tDeployment.Spec.Replicas = &test.deploymentScaleUpdateReplicas

			deploymentReplicasToAdd := *tDeployment.Spec.Replicas + test.deploymentMaxSurge - allReplicas
			deploymentReplicasToAdd -= test.terminatingReplica

			switch {
			case deploymentReplicasToAdd >= 0:
				sort.Sort(controller.ReplicaSetsBySizeNewer(allActiveRSs))
			case deploymentReplicasToAdd < 0:
				sort.Sort(controller.ReplicaSetsBySizeOlder(allActiveRSs))
			}

			scalePlan := PrepareProportionalScalePlan(logger, allActiveRSs, &tDeployment, deploymentReplicasToAdd)

			if !reflect.DeepEqual(test.expectedScalePlan, scalePlan) {
				t.Fatalf("unexpected scale plan: %s", cmp.Diff(test.expectedScalePlan, scalePlan))
			}
		})
	}
}
