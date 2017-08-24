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

package cpumanager

import (
	"fmt"
	"reflect"
	"testing"

	"k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/kubelet/cm/cpuset"
	"k8s.io/kubernetes/pkg/kubelet/cpumanager/topology"
)

type staticPolicyTest struct {
	description     string
	topo            *topology.CPUTopology
	containerID     string
	stAssignments   map[string]cpuset.CPUSet
	stDefaultCPUSet cpuset.CPUSet
	pod             *v1.Pod
	expErr          error
	expCPUAlloc     bool
	expCSet         cpuset.CPUSet
}

func TestStaticPolicyName(t *testing.T) {
	policy := NewStaticPolicy(topoSingleSocketHT)

	policyName := policy.Name()
	if policyName != "static" {
		t.Errorf("StaticPolicy Name() error. expected: static, returned: %v",
			policyName)
	}
}

func TestStaticPolicyStart(t *testing.T) {
	policy := NewStaticPolicy(topoSingleSocketHT).(*staticPolicy)

	st := &mockState{
		assignments:   map[string]cpuset.CPUSet{},
		defaultCPUSet: cpuset.NewCPUSet(),
	}

	policy.Start(st)
	for cpuid := 1; cpuid < policy.topology.NumCPUs; cpuid++ {
		if !st.defaultCPUSet.Contains(cpuid) {
			t.Errorf("StaticPolicy Start() error. expected cpuid %d to be present in defaultCPUSet", cpuid)
		}
	}
}

func TestStaticPolicyRegister(t *testing.T) {
	testCases := []staticPolicyTest{
		{
			description:     "GuPodSingleCore, SingleSocketHT, ExpectAllocOneCPU",
			topo:            topoSingleSocketHT,
			containerID:     "fakeID2",
			stAssignments:   map[string]cpuset.CPUSet{},
			stDefaultCPUSet: cpuset.NewCPUSet(1, 2, 3, 4, 5, 6, 7),
			pod:             makePod("1000m", "1000m"),
			expErr:          nil,
			expCPUAlloc:     true,
			expCSet:         cpuset.NewCPUSet(4), // expect sibling of partial core
		},
		{
			description: "GuPodMultipleCores, SingleSocketHT, ExpectAllocOneCore",
			topo:        topoSingleSocketHT,
			containerID: "fakeID3",
			stAssignments: map[string]cpuset.CPUSet{
				"fakeID100": cpuset.NewCPUSet(2, 3, 6, 7),
			},
			stDefaultCPUSet: cpuset.NewCPUSet(1, 4, 5),
			pod:             makePod("2000m", "2000m"),
			expErr:          nil,
			expCPUAlloc:     true,
			expCSet:         cpuset.NewCPUSet(1, 5),
		},
		{
			description: "GuPodMultipleCores, DualSocketHT, ExpectAllocOneSocket",
			topo:        topoDualSocketHT,
			containerID: "fakeID3",
			stAssignments: map[string]cpuset.CPUSet{
				"fakeID100": cpuset.NewCPUSet(2),
			},
			stDefaultCPUSet: cpuset.NewCPUSet(1, 3, 4, 5, 6, 7, 8, 9, 10, 11),
			pod:             makePod("6000m", "6000m"),
			expErr:          nil,
			expCPUAlloc:     true,
			expCSet:         cpuset.NewCPUSet(1, 3, 5, 7, 9, 11),
		},
		{
			description: "GuPodMultipleCores, DualSocketHT, ExpectAllocThreeCores",
			topo:        topoDualSocketHT,
			containerID: "fakeID3",
			stAssignments: map[string]cpuset.CPUSet{
				"fakeID100": cpuset.NewCPUSet(1, 5),
			},
			stDefaultCPUSet: cpuset.NewCPUSet(2, 3, 4, 6, 7, 8, 9, 10, 11),
			pod:             makePod("6000m", "6000m"),
			expErr:          nil,
			expCPUAlloc:     true,
			expCSet:         cpuset.NewCPUSet(2, 3, 4, 8, 9, 10),
		},
		{
			description: "GuPodMultipleCores, DualSocketNoHT, ExpectAllocOneSocket",
			topo:        topoDualSocketNoHT,
			containerID: "fakeID1",
			stAssignments: map[string]cpuset.CPUSet{
				"fakeID100": cpuset.NewCPUSet(),
			},
			stDefaultCPUSet: cpuset.NewCPUSet(1, 3, 4, 5, 6, 7),
			pod:             makePod("4000m", "4000m"),
			expErr:          nil,
			expCPUAlloc:     true,
			expCSet:         cpuset.NewCPUSet(4, 5, 6, 7),
		},
		{
			description: "GuPodMultipleCores, DualSocketNoHT, ExpectAllocFourCores",
			topo:        topoDualSocketNoHT,
			containerID: "fakeID1",
			stAssignments: map[string]cpuset.CPUSet{
				"fakeID100": cpuset.NewCPUSet(4, 5),
			},
			stDefaultCPUSet: cpuset.NewCPUSet(1, 3, 6, 7),
			pod:             makePod("4000m", "4000m"),
			expErr:          nil,
			expCPUAlloc:     true,
			expCSet:         cpuset.NewCPUSet(1, 3, 6, 7),
		},
		{
			description: "GuPodMultipleCores, DualSocketHT, ExpectAllocOneSocketOneCore",
			topo:        topoDualSocketHT,
			containerID: "fakeID3",
			stAssignments: map[string]cpuset.CPUSet{
				"fakeID100": cpuset.NewCPUSet(2),
			},
			stDefaultCPUSet: cpuset.NewCPUSet(1, 3, 4, 5, 6, 7, 8, 9, 10, 11),
			pod:             makePod("8000m", "8000m"),
			expErr:          nil,
			expCPUAlloc:     true,
			expCSet:         cpuset.NewCPUSet(1, 3, 4, 5, 7, 9, 10, 11),
		},
		{
			description:     "NonGuPod, SingleSocketHT, NoAlloc",
			topo:            topoSingleSocketHT,
			containerID:     "fakeID1",
			stAssignments:   map[string]cpuset.CPUSet{},
			stDefaultCPUSet: cpuset.NewCPUSet(1, 2, 3, 4, 5, 6, 7),
			pod:             makePod("1000m", "2000m"),
			expErr:          nil,
			expCPUAlloc:     false,
			expCSet:         cpuset.NewCPUSet(),
		},
		{
			description:     "GuPodNonIntegerCore, SingleSocketHT, NoAlloc",
			topo:            topoSingleSocketHT,
			containerID:     "fakeID4",
			stAssignments:   map[string]cpuset.CPUSet{},
			stDefaultCPUSet: cpuset.NewCPUSet(1, 2, 3, 4, 5, 6, 7),
			pod:             makePod("977m", "977m"),
			expErr:          nil,
			expCPUAlloc:     false,
			expCSet:         cpuset.NewCPUSet(),
		},
		{
			description: "GuPodMultipleCores, SingleSocketHT, NoAllocExpectError",
			topo:        topoSingleSocketHT,
			containerID: "fakeID5",
			stAssignments: map[string]cpuset.CPUSet{
				"fakeID100": cpuset.NewCPUSet(1, 2, 3, 4, 5, 6),
			},
			stDefaultCPUSet: cpuset.NewCPUSet(7),
			pod:             makePod("2000m", "2000m"),
			expErr:          fmt.Errorf("not enough cpus available to satisfy request"),
			expCPUAlloc:     false,
			expCSet:         cpuset.NewCPUSet(),
		},
		{
			description: "GuPodMultipleCores, DualSocketHT, NoAllocExpectError",
			topo:        topoDualSocketHT,
			containerID: "fakeID5",
			stAssignments: map[string]cpuset.CPUSet{
				"fakeID100": cpuset.NewCPUSet(1, 2, 3),
			},
			stDefaultCPUSet: cpuset.NewCPUSet(4, 5, 6, 7, 8, 9, 10, 11),
			pod:             makePod("10000m", "10000m"),
			expErr:          fmt.Errorf("not enough cpus available to satisfy request"),
			expCPUAlloc:     false,
			expCSet:         cpuset.NewCPUSet(),
		},
	}

	for _, testCase := range testCases {
		policy := NewStaticPolicy(testCase.topo)

		st := &mockState{
			assignments:   testCase.stAssignments,
			defaultCPUSet: testCase.stDefaultCPUSet,
		}

		container := &testCase.pod.Spec.Containers[0]
		err := policy.RegisterContainer(st, testCase.pod, container, testCase.containerID)
		if !reflect.DeepEqual(err, testCase.expErr) {
			t.Errorf("StaticPolicy Register() error (%v). expected register error: %v but got: %v",
				testCase.description, testCase.expErr, err)
		}

		if testCase.expCPUAlloc {
			cset, found := st.assignments[testCase.containerID]
			if !found {
				t.Errorf("StaticPolicy Register() error (%v). expected container id %v to be present in assignments %v",
					testCase.description, testCase.containerID, st.assignments)
			}

			if !reflect.DeepEqual(cset, testCase.expCSet) {
				t.Errorf("StaticPolicy Register() error (%v). expected cpuset %v but got %v",
					testCase.description, testCase.expCSet, cset)
			}

			if !cset.Intersection(st.defaultCPUSet).IsEmpty() {
				t.Errorf("StaticPolicy Register() error (%v). expected cpuset %v to be disoint from the shared cpuset %v",
					testCase.description, cset, st.defaultCPUSet)
			}
		}

		if !testCase.expCPUAlloc {
			_, found := st.assignments[testCase.containerID]
			if found {
				t.Errorf("StaticPolicy Register() error (%v). Did not expect container id %v to be present in assignments %v",
					testCase.description, testCase.containerID, st.assignments)
			}
		}

		if st.defaultCPUSet.IsEmpty() && !st.GetPressure() {
			t.Error("StaticPolicy CPU pressure error: Default CPU set is empty but policy reporting no CPU pressure")
		}
	}
}

func TestStaticPolicyUnregister(t *testing.T) {
	testCases := []staticPolicyTest{
		{
			description: "SingleSocketHT, DeAllocOneContainer",
			topo:        topoSingleSocketHT,
			containerID: "fakeID1",
			stAssignments: map[string]cpuset.CPUSet{
				"fakeID1": cpuset.NewCPUSet(1, 2, 3),
			},
			stDefaultCPUSet: cpuset.NewCPUSet(4, 5, 6, 7),
			expCSet:         cpuset.NewCPUSet(1, 2, 3, 4, 5, 6, 7),
		},
		{
			description: "SingleSocketHT, DeAllocOneContainer, BeginEmpty",
			topo:        topoSingleSocketHT,
			containerID: "fakeID1",
			stAssignments: map[string]cpuset.CPUSet{
				"fakeID1": cpuset.NewCPUSet(1, 2, 3),
				"fakeID2": cpuset.NewCPUSet(4, 5, 6, 7),
			},
			stDefaultCPUSet: cpuset.NewCPUSet(),
			expCSet:         cpuset.NewCPUSet(1, 2, 3),
		},
		{
			description: "SingleSocketHT, DeAllocTwoContainer",
			topo:        topoSingleSocketHT,
			containerID: "fakeID1",
			stAssignments: map[string]cpuset.CPUSet{
				"fakeID1": cpuset.NewCPUSet(1, 3, 5),
				"fakeID2": cpuset.NewCPUSet(2, 4),
			},
			stDefaultCPUSet: cpuset.NewCPUSet(6, 7),
			expCSet:         cpuset.NewCPUSet(1, 3, 5, 6, 7),
		},
		{
			description: "SingleSocketHT, NoDeAlloc",
			topo:        topoSingleSocketHT,
			containerID: "fakeID2",
			stAssignments: map[string]cpuset.CPUSet{
				"fakeID1": cpuset.NewCPUSet(1, 3, 5),
			},
			stDefaultCPUSet: cpuset.NewCPUSet(2, 4, 6, 7),
			expCSet:         cpuset.NewCPUSet(2, 4, 6, 7),
		},
	}

	for _, testCase := range testCases {
		policy := NewStaticPolicy(testCase.topo)

		st := &mockState{
			assignments:   testCase.stAssignments,
			defaultCPUSet: testCase.stDefaultCPUSet,
		}

		policy.UnregisterContainer(st, testCase.containerID)

		if !reflect.DeepEqual(st.defaultCPUSet, testCase.expCSet) {
			t.Errorf("StaticPolicy UnRegister() error (%v). expected default cpuset %v but got %v",
				testCase.description, testCase.expCSet, st.defaultCPUSet)
		}

		if _, found := st.assignments[testCase.containerID]; found {
			t.Errorf("StaticPolicy UnRegister() error (%v). expected containerID %v not be in assignments %v",
				testCase.description, testCase.containerID, st.assignments)
		}
	}
}
