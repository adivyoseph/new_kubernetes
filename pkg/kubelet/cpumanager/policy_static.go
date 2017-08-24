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
	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	v1qos "k8s.io/kubernetes/pkg/api/v1/helper/qos"
	"k8s.io/kubernetes/pkg/kubelet/cm/cpuset"
	"k8s.io/kubernetes/pkg/kubelet/cpumanager/state"
	"k8s.io/kubernetes/pkg/kubelet/cpumanager/topology"
)

// PolicyStatic name of static policy
const PolicyStatic policyName = "static"

var _ Policy = &staticPolicy{}

type staticPolicy struct {
	topology *topology.CPUTopology
}

// Ensure staticPolicy implements Policy interface
var _ Policy = &staticPolicy{}

// NewStaticPolicy returns a cpuset manager policy that does not change
// CPU assignments for exclusively pinned guaranteed containers after
// the main container process starts.
func NewStaticPolicy(topology *topology.CPUTopology) Policy {
	return &staticPolicy{
		topology: topology,
	}
}

func (p *staticPolicy) Name() string {
	return string(PolicyStatic)
}

func (p *staticPolicy) Start(s state.State) {
	allCPUs := p.topology.CPUDetails.CPUs()
	// takeByTopology allocates CPUs associated with low-numbered cores from
	// allCPUs.
	//
	// For example: Given a system with 8 CPUs available and HT enabled,
	// if NumReservedCores=2, then reserved={0,4}
	reserved, _ := takeByTopology(p.topology, allCPUs, p.topology.NumReservedCores)
	s.SetDefaultCPUSet(allCPUs.Difference(reserved))
}

func (p *staticPolicy) RegisterContainer(s state.State, pod *v1.Pod, container *v1.Container, containerID string) error {
	glog.Infof("[cpumanager] static policy: RegisterContainer (pod: %s, container: %s, container id: %s)", pod.Name, container.Name, containerID)
	if numCPUs := guaranteedCPUs(pod, container); numCPUs != 0 {
		// container belongs in an exclusively allocated pool
		cpuset, err := p.allocateCPUs(s, numCPUs)
		if err != nil {
			glog.Errorf("[cpumanager] unable to allocate %d CPUs (container id: %s, error: %v)", numCPUs, containerID, err)
			return err
		}
		s.SetCPUSet(containerID, cpuset)
	}
	// container belongs in the shared pool (nothing to do; use default cpuset)
	return nil
}

func (p *staticPolicy) UnregisterContainer(s state.State, containerID string) error {
	glog.Infof("[cpumanager] static policy: UnregisterContainer (container id: %s)", containerID)
	if toRelease, ok := s.GetCPUSet(containerID); ok {
		s.Delete(containerID)
		p.releaseCPUs(s, toRelease)
	}
	return nil
}

func (p *staticPolicy) allocateCPUs(s state.State, numCPUs int) (cpuset.CPUSet, error) {
	glog.Infof("[cpumanager] allocateCpus: (numCPUs: %d)", numCPUs)
	result, err := takeByTopology(p.topology, s.GetDefaultCPUSet(), numCPUs)
	if err != nil {
		return cpuset.NewCPUSet(), err
	}
	// Remove allocated CPUs from the shared CPUSet.
	s.SetDefaultCPUSet(s.GetDefaultCPUSet().Difference(result))
	// Set CPU pressure condition.
	s.SetPressure(s.GetDefaultCPUSet().IsEmpty())

	glog.Infof("[cpumanager] allocateCPUs: returning \"%v\"", result)
	return result, nil
}

func (p *staticPolicy) releaseCPUs(s state.State, release cpuset.CPUSet) {
	// Mutate the shared pool, adding supplied cpus.
	s.SetDefaultCPUSet(s.GetDefaultCPUSet().Union(release))
	// Set CPU pressure condition.
	s.SetPressure(s.GetDefaultCPUSet().IsEmpty())
}

func guaranteedCPUs(pod *v1.Pod, container *v1.Container) int {
	if v1qos.GetPodQOS(pod) != v1.PodQOSGuaranteed {
		return 0
	}
	cpuQuantity := container.Resources.Requests[v1.ResourceCPU]
	glog.Infof("[cpumanager] guaranteedCpus (container: %s, cpu request: %v)", container.Name, cpuQuantity.MilliValue())
	if cpuQuantity.Value()*1000 != cpuQuantity.MilliValue() {
		return 0
	}
	// Safe downcast to do for all systems with < 2.1 billion CPUs.
	// Per the language spec, `int` is guaranteed to be at least 32 bits wide.
	// https://golang.org/ref/spec#Numeric_types
	return int(cpuQuantity.Value())
}
