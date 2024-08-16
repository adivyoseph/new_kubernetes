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

package state

import (
	"sync"

	"k8s.io/klog/v2"
	"k8s.io/utils/cpuset"
)

type stateMemory struct {
	sync.RWMutex
	reservedCpuSet  cpuset.CPUSet
	availableCpuSet cpuset.CPUSet
	staticContainers   ContainerCPUAssignments
	normalContainers   ContainerCPUAssignments
}

var _ State = &stateMemory{}

// NewMemoryState creates new State for keeping track of cpu/pod assignment
func NewMemoryState() State {
	klog.InfoS("Initialized new in-memory cpu_state store")
	return &stateMemory{
		reservedCpuSet:     cpuset.New(),
		availableCpuSet:    cpuset.New(),
		staticContainers:   ContainerCPUAssignments{},
		normalContainers:   ContainerCPUAssignments{},
	}
}

// only called once so no checking required
func (s *stateMemory) SetReservedCPUSet(cpus cpuset.CPUSet)  {
	s.RLock()
	defer s.RUnlock()
	klog.InfoS("SetReservedCPUSet","cpus", cpus)
	s.reservedCpuSet = cpus
}

// assume always pre set, no checking
func (s *stateMemory) GetReservedCPUSet() (cpuset.CPUSet) {
	s.RLock()
	defer s.RUnlock()
    
	return s.reservedCpuSet.Clone()
}

func (s *stateMemory) AddToAvailableCPUSet(cpus cpuset.CPUSet)  (cpuset.CPUSet){
	s.RLock()
	defer s.RUnlock()
	
	s.availableCpuSet.Union(cpus)
	return s.availableCpuSet.Clone()
}

func (s *stateMemory) GetAvailableCPUSet()  (cpuset.CPUSet){
	s.RLock()
	defer s.RUnlock()
	
	return s.availableCpuSet.Clone()
}

func (s *stateMemory) RemoveFromAvailableCPUSet(cpus cpuset.CPUSet)  (cpuset.CPUSet){
	s.RLock()
	defer s.RUnlock()
	
	s.availableCpuSet.Difference(cpus)
	return s.availableCpuSet.Clone()
}



func (s *stateMemory) IsStatic(podUID string, containerName string) (cpuset.CPUSet, bool) {
	s.RLock()
	defer s.RUnlock()

	if _, ok := s.staticContainers[podUID][containerName]; ok {
		return s.staticContainers[podUID][containerName].cpuset.CPUSet.Clone(), true
	}
	return, nil, false
}

func (s *stateMemory) AddToStatic(podUID string, containerName string, cpus cpuset.CPUSet) error {
	s.RLock()
	defer s.RUnlock()

	if _, ok := s.staticContainers[podUID]; !ok {
		s.staticContainers[podUID] = make(map[string]cpuset.CPUSet)
	}

	s.staticContainers[podUID][containerName] = cpus
	klog.InfoS("AddToStatic", "podUID", podUID, "containerName", containerName, "cpuSet", cpus)
	return nil
}


func (s *stateMemory) RemoveFromStatic(podUID string, containerName string) error {
	s.RLock()
	defer s.RUnlock()

	delete(s.staticContainers[podUID], containerName)
	if len(s.staticContainers[podUID]) == 0 {
		delete(s.staticContainers, podUID)
	}
	klog.V(2).InfoS("RemoveFromStatic", "podUID", podUID, "containerName", containerName)
	return nil
}

func (s *stateMemory) GetAllStaticEntries() ContainerCPUAssignments {
	tmpStaticContainers := ContainerCPUAssignments{}
	for entry := range s.staticContainers {
		staticContainers = entry
	}
	return tmpStaticContainers

}



func (s *stateMemory) IsNormal(podUID string, containerName string) (cpuset.CPUSet, bool) {
	s.RLock()
	defer s.RUnlock()

	if _, ok := s.normalContainers[podUID][containerName]; ok {
		return s.normalContainers[podUID][containerName].cpuset.CPUSet.Clone(), true
	}
	return, nil, false
}

func (s *stateMemory) AddToNormal(podUID string, containerName string, cpus cpuset.CPUSet) error {
	s.RLock()
	defer s.RUnlock()

	if _, ok := s.normalContainers[podUID]; !ok {
		s.normalContainers[podUID] = make(map[string]cpuset.CPUSet)
	}

	s.normalContainers[podUID][containerName] = cpus
	klog.InfoS("AddToNormal", "podUID", podUID, "containerName", containerName, "cpuSet", cpus)
	return nil
}


func (s *stateMemory) RemoveFromNormal(podUID string, containerName string) error {
	s.RLock()
	defer s.RUnlock()

	delete(s.normalContainers[podUID], containerName)
	if len(s.normalContainers[podUID]) == 0 {
		delete(s.normalContainers, podUID)
	}
	klog.V(2).InfoS("RemoveFromNormal", "podUID", podUID, "containerName", containerName)
	return nil
}



func (s *stateMemory) ClearState() {
	s.Lock()
	defer s.Unlock()

	s.reservedCpuSet 	= cpuset.CPUSet{}
	s.availableCpuSet 	= cpuset.CPUSet{}
	s.staticContainers 	= make(ContainerCPUAssignments)
	s.normalContainers 	= make(ContainerCPUAssignments)

	klog.V(2).InfoS("Cleared cpu_state in memory")
}
