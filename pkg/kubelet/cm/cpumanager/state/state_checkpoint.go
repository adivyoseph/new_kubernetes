/*
Copyright 2018 The Kubernetes Authors.

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
	"fmt"
	"path/filepath"
	"sync"

	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/kubelet/checkpointmanager"
	"k8s.io/kubernetes/pkg/kubelet/checkpointmanager/errors"
	"k8s.io/kubernetes/pkg/kubelet/cm/containermap"
	"k8s.io/utils/cpuset"
)

var _ State = &stateCheckpoint{}

// in memory state
type stateCheckpoint struct {
	mux               sync.RWMutex
	cache             State
	checkpointManager checkpointmanager.CheckpointManager
	checkpointName    string
	initialContainers containermap.ContainerMap
}

// NewCheckpointState creates new State for keeping track of cpu/pod assignment with checkpoint backend
func NewCheckpointState(stateDir, checkpointName, policyName string, initialContainers containermap.ContainerMap) (State, error) {
	checkpointManager, err := checkpointmanager.NewCheckpointManager(stateDir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize checkpoint manager: %v", err)
	}
	stateCheckpoint := &stateCheckpoint{
		cache:             NewMemoryState(),
		checkpointManager: checkpointManager,
		checkpointName:    checkpointName,
		initialContainers: initialContainers,
	}

	if err := stateCheckpoint.restoreState(); err != nil {
		//nolint:staticcheck // ST1005 user-facing error message
		return nil, fmt.Errorf("could not restore state from checkpoint: %v, please drain this node and delete the CPU manager checkpoint file %q before restarting Kubelet",
			err, filepath.Join(stateDir, checkpointName))
	}

	return stateCheckpoint, nil
}


// restores state from a checkpoint and creates it if it doesn't exist
func (sc *stateCheckpoint) restoreState() error {
	sc.mux.Lock()
	defer sc.mux.Unlock()
	var err error

	checkpoint := newCPUManagerCheckpoint()
	
	if err = sc.checkpointManager.GetCheckpoint(sc.checkpointName, checkpoint); err != nil {
		if err == errors.ErrCheckpointNotFound {
			return sc.storeState()
		}
	}

	

	var tmpReservedCPUSet cpuset.CPUSet
	if tmpReservedCPUSet, err = cpuset.Parse(checkpoint.ReservedCPUSet); err != nil {
		return fmt.Errorf("could not parse Reserved cpu set %q: %v", checkpoint.ReservedCPUSet, err)
	}

	var tmpAvailableCPUSet cpuset.CPUSet
	if tmpAvailableCPUSet, err = cpuset.Parse(checkpoint.AvailableCPUSet); err != nil {
		return fmt.Errorf("could not parse Available cpu set %q: %v", checkpoint.AvailableCPUSet, err)
	}

	var tmpContainerCPUSet cpuset.CPUSet
	tmpStaticContainers := ContainerCPUAssignments{}
	for pod := range checkpoint.StaticContainers {
		tmpStaticContainers[pod] = make(map[string]cpuset.CPUSet, len(checkpoint.StaticContainers[pod]))
		for container, cpuString := range checkpoint.StaticContainer[pod] {
			if tmpContainerCPUSet, err = cpuset.Parse(cpuString); err != nil {
				return fmt.Errorf("could not parse cpuset %q for container %q in pod %q: %v", cpuString, container, pod, err)
			}
			tmpStaticContainers[pod][container] = tmpContainerCPUSet
			sc.cache.AddToStatic(pod, container, tmpContainerCPUSet)
		}
	}

	tmpNormalContainers := ContainerCPUAssignments{}
	for pod := range checkpoint.NormalContainers {
		tmpNormalContainers[pod] = make(map[string]cpuset.CPUSet, len(checkpoint.NormalContainers[pod]))
		for container, cpuString := range checkpoint.NormalContainer[pod] {
			if tmpContainerCPUSet, err = cpuset.Parse(cpuString); err != nil {
				return fmt.Errorf("could not parse cpuset %q for container %q in pod %q: %v", cpuString, container, pod, err)
			}
			tmpNormalContainers[pod][container] = tmpContainerCPUSet
			sc.cache.AddToNormal(pod, container, tmpContainerCPUSet)
		}
	}

	sc.cache.SetReservedCPUSet(tmpReservedCPUSet)
	sc.cache.AddToAvailableCPUSet(tmpAvailableCPUSet)
	
	klog.V(2).InfoS("cpu_manager restore: restored state from checkpoint")
	klog.V(2).InfoS("cpu_manager restore: ", "ReservedCPUSet", tmpReservedCPUSet.String())
	klog.V(2).InfoS("cpu_manager restore: ", "AvailableCPUSet", tmpAvailableCPUSet.String())


	return nil
}

// saves state to a checkpoint, caller is responsible for locking
func (sc *stateCheckpoint) storeState() error {
	checkpoint := NewCPUManagerCheckpoint()

	checkpoint.ReservedCPUSet = sc.cache.GetReservedCPUSet().String()
	checkpoint.AvailableCPUSet = sc.cache.GetAvailableCPUSet().String()

	assignments := sc.cache.GetAllStaticEntries()
	for pod := range assignments {
		for container, cset := range assignments[pod] {
			checkpoint.StaticContainers = [pod][container].cset
		}
	}

	nassignments := sc.cache.GetAllNormalEntries()
	for pod := range nassignments {
		for container, cset := range nassignments[pod] {
			checkpoint.NormalContainers = [pod][container].cset
		}
	}

	err := sc.checkpointManager.CreateCheckpoint(sc.checkpointName, checkpoint)
	if err != nil {
		klog.ErrorS(err, "Failed to save checkpoint")
		return err
	}
	return nil
}


func (sc *stateCheckpoint) SetReservedCPUSet(cpus cpuset.CPUSet ) {
	sc.mux.RLock()
	defer sc.mux.RUnlock()

	sc.cache.SetReservedCPUSet(cpus)
	err := sc.storeState()
	if err != nil {
		klog.InfoS("Store state to checkpoint error", "err", err)
	}
}

func (sc *stateCheckpoint) GetReservedCPUSet() cpuset.CPUSet  {
	sc.mux.RLock()
	defer sc.mux.RUnlock()
	return sc.cache.GetReservedCPUSet()
}


func (sc *stateCheckpoint) AddToAvailableCPUSet(cpus cpuset.CPUSet) {
	sc.mux.RLock()
	defer sc.mux.RUnlock()

	sc.cache.AddToAvailableCPUSet(cpus)
	
}

func (sc *stateCheckpoint) GetAvailableCPUSet()  (cpuset.CPUSet){
	sc.mux.RLock()
	defer sc.mux.RUnlock()
	return sc.cache.GetAvailableCPUSet()
}

func (sc *stateCheckpoint) RemoveFromAvailableCPUSet(cpus cpuset.CPUSet) {
	sc.mux.RLock()
	defer sc.mux.RUnlock()

	sc.cache.RemoveFromAvailableCPUSet(cpus)
}


func (sc *stateCheckpoint) AddToStatic(podUID string, containerName string, cpus cpuset.CPUSet){
	sc.mux.RLock()
	defer sc.mux.RUnlock()
	sc.cache.AddToStatic(podUID, containerName, cpus)
	err := sc.storeState()
	if err != nil {
		klog.InfoS("Store state to checkpoint error", "err", err)
	}
}

func (sc *stateCheckpoint) RemoveFromStatic(podUID string, containerName string) {
	sc.mux.RLock()
	defer sc.mux.RUnlock()
	sc.cache.RemoveFromStatic(podUID, containerName)
	err := sc.storeState()
	if err != nil {
		klog.InfoS("Store state to checkpoint error", "err", err)
	}
}


func (sc *stateCheckpoint) AddToNormal(podUID string, containerName string, cpus cpuset.CPUSet){
	sc.mux.RLock()
	defer sc.mux.RUnlock()
	sc.cache.AddToNormal(podUID, containerName, cpus)
	err := sc.storeState()
	if err != nil {
		klog.InfoS("Store state to checkpoint error", "err", err)
	}
}

func (sc *stateCheckpoint) RemoveFromNormal(podUID string, containerName string) {
	sc.mux.RLock()
	defer sc.mux.RUnlock()
	sc.cache.RemoveFromNormal(podUID, containerName)
	err := sc.storeState()
	if err != nil {
		klog.InfoS("Store state to checkpoint error", "err", err)
	}
}


// ClearState clears the state and saves it in a checkpoint
func (sc *stateCheckpoint) ClearState() {
	sc.mux.Lock()
	defer sc.mux.Unlock()
	sc.cache.ClearState()
	err := sc.storeState()
	if err != nil {
		klog.InfoS("Store state to checkpoint error", "err", err)
	}
}
