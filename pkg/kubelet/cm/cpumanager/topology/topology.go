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

package topology

import (
	"fmt"

	cadvisorapi "github.com/google/cadvisor/info/v1"
	"k8s.io/klog/v2"
	"k8s.io/utils/cpuset"
)

// NUMANodeInfo is a map from NUMANode ID to a list of CPU IDs associated with
// that NUMANode.
type NUMANodeInfo map[int]cpuset.CPUSet

// CPUInfo contains the NUMA, socket, and core IDs associated with a CPU.
type CPUInfo struct {
	SocketID   int //server global and unique
	NUMANodeID int //server global and unique
	L3GroupID  int //server global and unique
	CoreID     int
}

// CPUDetails is a map from CPU ID to Core ID, Socket ID, and NUMA ID.
type CPUDetails map[int]CPUInfo

// CPUTopology contains details of node cpu, where :
// CPU  - logical CPU, cadvisor - thread
// Core - physical CPU, cadvisor - Core
// Socket - socket, cadvisor - Socket
// NUMA Node - NUMA cell, cadvisor - Node
type CPUTopology struct {
	NumCPUs      int
	NumCores     int
	NumSockets   int
	NumNUMANodes int
	CPUDetails   CPUDetails
}

type CPUResources struct {
	cpuNumber int // OS cpu number
	owner     int // free, reserved, shared, strict
}

type CoreResources struct {
	Id   int // HW Id
	Cpus []CPUResources
}

type L3GroupResources struct {
	Id int // system wide unique
	// ordered
	Cores []CoreResources
}

type NumaNodeResources struct {
	Id       int // system wide unique
	L3Groups []L3GroupResources
}

type SocketResources struct {
	Id int // system wide unique Id
	//
	NumaNodes []NumaNodeResources
}

// Server top down graph of CPU resources
type ServerTopology struct {
	Sockets []SocketResources
}

func (srv *ServerTopology) findSocket(socketId int) int {
	for index, socket := range srv.Sockets {
		if socket.Id == socketId {
			return index
		}
	}
	// new socket
	newSocket := SocketResources{}
	newSocket.Id = socketId
	srv.Sockets = append(srv.Sockets, newSocket)
	return len(srv.Sockets) - 1
}

func (srv *ServerTopology) findNuma(sockIndex int, numaId int) int {
	for nodeIndex, node := range srv.Sockets[sockIndex].NumaNodes {
		if node.Id == numaId {
			return nodeIndex
		}
	}
	// new node
	newNode := NumaNodeResources{}
	newNode.Id = numaId
	srv.Sockets[sockIndex].NumaNodes = append(srv.Sockets[sockIndex].NumaNodes, newNode)
	return len(srv.Sockets[sockIndex].NumaNodes) - 1
}

func (srv *ServerTopology) findL3Group(sockIndex int, numaIndex int, l3Group int) int {
	for l3GroupIndex, L3Group := range srv.Sockets[sockIndex].NumaNodes[numaIndex].L3Groups {
		if L3Group.Id == l3Group {
			return l3GroupIndex
		}
	}
	// new l3Group
	newL3Group := L3GroupResources{}
	newL3Group.Id = l3Group
	srv.Sockets[sockIndex].NumaNodes[numaIndex].L3Groups = append(srv.Sockets[sockIndex].NumaNodes[numaIndex].L3Groups, newL3Group)
	return len(srv.Sockets[sockIndex].NumaNodes[numaIndex].L3Groups) - 1
}

func (srv *ServerTopology) findCore(sockIndex int, numaIndex int, l3GroupIndex int, core int) int {
	for coreIndex, Core := range srv.Sockets[sockIndex].NumaNodes[numaIndex].L3Groups[l3GroupIndex].Cores {
		if Core.Id == core {
			return coreIndex
		}
	}
	// new core
	newCore := CoreResources{}
	newCore.Id = core
	srv.Sockets[sockIndex].NumaNodes[numaIndex].L3Groups[l3GroupIndex].Cores = append(srv.Sockets[sockIndex].NumaNodes[numaIndex].L3Groups[l3GroupIndex].Cores, newCore)
	return len(srv.Sockets[sockIndex].NumaNodes[numaIndex].L3Groups[l3GroupIndex].Cores) - 1
}

func (topo *CPUTopology) NewServerTopology() *ServerTopology {
	srv := ServerTopology{}

	for cpu, cpuInfo := range topo.CPUDetails {
		SocketIndex := srv.findSocket(cpuInfo.SocketID)
		NumaNodeIndex := srv.findNuma(SocketIndex, cpuInfo.NUMANodeID)
		L3GroupIndex := srv.findL3Group(SocketIndex, NumaNodeIndex, cpuInfo.L3GroupID)
		CoreIndex := srv.findCore(SocketIndex, NumaNodeIndex, L3GroupIndex, cpuInfo.CoreID)
		found := 0
		for _, Cpu := range srv.Sockets[SocketIndex].NumaNodes[NumaNodeIndex].L3Groups[L3GroupIndex].Cores[CoreIndex].Cpus {
			if cpu == Cpu {
				found = 1
			}
		}
		if found == 0 {
			newCpu := cpu
			srv.Sockets[SocketIndex].NumaNodes[NumaNodeIndex].L3Groups[L3GroupIndex].Cores[CoreIndex].Cpus =
				append(srv.Sockets[SocketIndex].NumaNodes[NumaNodeIndex].L3Groups[L3GroupIndex].Cores[CoreIndex].Cpus, newCpu)
		}
	}

	//sort cpus within a core
	if len(srv.Sockets[0].NumaNodes[0].L3Groups[0].Cores[0].Cpus) > 1 {
		for _, socket := range srv.Sockets {
			for _, node := range socket.NumaNodes {
				for _, l3Group := range node.L3Groups {
					for _, core := range l3Group.Cores {
						cpu0 := core.Cpus[0]
						if cpu0 > core.Cpus[1] {
							core.Cpus[0] = core.Cpus[1]
							core.Cpus[1] = cpu0
						}
					}
				}
			}
		}
	}
	//sort Cores within an L3Group
	if len(srv.Sockets[0].NumaNodes[0].L3Groups[0].Cores) > 1 {
		for _, socket := range srv.Sockets {
			for _, node := range socket.NumaNodes {
				for _, l3Group := range node.L3Groups {
					for cores := 0; cores < len(l3Group.Cores); cores++ {
						for coreIndex := 1; coreIndex < len(l3Group.Cores); coreIndex++ {
							if l3Group.Cores[coreIndex-1].Cpus[0] > l3Group.Cores[coreIndex].Cpus[0] {
								tempCore := l3Group.Cores[coreIndex-1].Cpus
								l3Group.Cores[coreIndex-1].Cpus = l3Group.Cores[coreIndex].Cpus
								l3Group.Cores[coreIndex].Cpus = tempCore
							}
						}
					}

				}
			}
		}
	}
	//sort L3Group with a numa node
	if len(srv.Sockets[0].NumaNodes[0].L3Groups) > 1 {
		for _, socket := range srv.Sockets {
			for _, node := range socket.NumaNodes {
				for l3Groups := 0; l3Groups < len(node.L3Groups); l3Groups++ {
					for l3GroupIndex := 1; l3GroupIndex < len(node.L3Groups); l3GroupIndex++ {
						if node.L3Groups[l3GroupIndex-1].Cores[0].Cpus[0] > node.L3Groups[l3GroupIndex].Cores[0].Cpus[0] {
							tempL3Group := node.L3Groups[l3GroupIndex-1]
							node.L3Groups[l3GroupIndex-1] = node.L3Groups[l3GroupIndex]
							node.L3Groups[l3GroupIndex] = tempL3Group

						}
					}
				}
			}
		}
	}

	//debug
	klog.InfoS("NewServerTopology()")
	for sockIndex := 0; sockIndex < len(srv.Sockets); sockIndex++ {
		klog.InfoS("socket", "index", sockIndex, "Id", srv.Sockets[sockIndex].Id)
		for numaIndex := 0; numaIndex < len(srv.Sockets[sockIndex].NumaNodes); numaIndex++ {
			klog.InfoS("node  ", "index", numaIndex, "Id", srv.Sockets[sockIndex].NumaNodes[numaIndex].Id)
			for l3Index := 0; l3Index < len(srv.Sockets[sockIndex].NumaNodes[numaIndex].L3Groups); l3Index++ {
				klog.InfoS("   l3group", "index", l3Index, "Id", srv.Sockets[sockIndex].NumaNodes[numaIndex].L3Groups[l3Index].Id)
				for coreIndex := 0; coreIndex < len(srv.Sockets[sockIndex].NumaNodes[numaIndex].L3Groups[l3Index].Cores); coreIndex++ {
					if len(srv.Sockets[sockIndex].NumaNodes[numaIndex].L3Groups[l3Index].Cores[coreIndex].Cpus) < 2 {
						klog.InfoS("        Core", "index", coreIndex, "cpu", srv.Sockets[sockIndex].NumaNodes[numaIndex].L3Groups[l3Index].Cores[coreIndex].Cpus[0])
					} else {
						klog.InfoS("        Core", "index", coreIndex, "cpu0", srv.Sockets[sockIndex].NumaNodes[numaIndex].L3Groups[l3Index].Cores[coreIndex].Cpus[0],
							"cpu1", srv.Sockets[sockIndex].NumaNodes[numaIndex].L3Groups[l3Index].Cores[coreIndex].Cpus[1])
					}
				}
			}
		}
	}

	return &srv
}

// mark reserved cpus in topology tree
// allow per cpu
func (srv *ServerTopology) setReservedCPUs(cpus cpuset.CPUSet) {

}

// CPUsPerCore returns the number of logical CPUs are associated with
// each core.
func (topo *CPUTopology) CPUsPerCore() int {
	if topo.NumCores == 0 {
		return 0
	}
	return topo.NumCPUs / topo.NumCores
}

// CPUsPerSocket returns the number of logical CPUs are associated with
// each socket.
func (topo *CPUTopology) CPUsPerSocket() int {
	if topo.NumSockets == 0 {
		return 0
	}
	return topo.NumCPUs / topo.NumSockets
}

// CPUCoreID returns the physical core ID which the given logical CPU
// belongs to.
func (topo *CPUTopology) CPUCoreID(cpu int) (int, error) {
	info, ok := topo.CPUDetails[cpu]
	if !ok {
		return -1, fmt.Errorf("unknown CPU ID: %d", cpu)
	}
	return info.CoreID, nil
}

// CPUCoreID returns the socket ID which the given logical CPU belongs to.
func (topo *CPUTopology) CPUSocketID(cpu int) (int, error) {
	info, ok := topo.CPUDetails[cpu]
	if !ok {
		return -1, fmt.Errorf("unknown CPU ID: %d", cpu)
	}
	return info.SocketID, nil
}

// CPUCoreID returns the NUMA node ID which the given logical CPU belongs to.
func (topo *CPUTopology) CPUNUMANodeID(cpu int) (int, error) {
	info, ok := topo.CPUDetails[cpu]
	if !ok {
		return -1, fmt.Errorf("unknown CPU ID: %d", cpu)
	}
	return info.NUMANodeID, nil
}

// KeepOnly returns a new CPUDetails object with only the supplied cpus.
func (d CPUDetails) KeepOnly(cpus cpuset.CPUSet) CPUDetails {
	result := CPUDetails{}
	for cpu, info := range d {
		if cpus.Contains(cpu) {
			result[cpu] = info
		}
	}
	return result
}

// NUMANodes returns all of the NUMANode IDs associated with the CPUs in this
// CPUDetails.
func (d CPUDetails) NUMANodes() cpuset.CPUSet {
	var numaNodeIDs []int
	for _, info := range d {
		numaNodeIDs = append(numaNodeIDs, info.NUMANodeID)
	}
	return cpuset.New(numaNodeIDs...)
}

// NUMANodesInSockets returns all of the logical NUMANode IDs associated with
// the given socket IDs in this CPUDetails.
func (d CPUDetails) NUMANodesInSockets(ids ...int) cpuset.CPUSet {
	var numaNodeIDs []int
	for _, id := range ids {
		for _, info := range d {
			if info.SocketID == id {
				numaNodeIDs = append(numaNodeIDs, info.NUMANodeID)
			}
		}
	}
	return cpuset.New(numaNodeIDs...)
}

// Sockets returns all of the socket IDs associated with the CPUs in this
// CPUDetails.
func (d CPUDetails) Sockets() cpuset.CPUSet {
	var socketIDs []int
	for _, info := range d {
		socketIDs = append(socketIDs, info.SocketID)
	}
	return cpuset.New(socketIDs...)
}

// CPUsInSockets returns all of the logical CPU IDs associated with the given
// socket IDs in this CPUDetails.
func (d CPUDetails) CPUsInSockets(ids ...int) cpuset.CPUSet {
	var cpuIDs []int
	for _, id := range ids {
		for cpu, info := range d {
			if info.SocketID == id {
				cpuIDs = append(cpuIDs, cpu)
			}
		}
	}
	return cpuset.New(cpuIDs...)
}

// SocketsInNUMANodes returns all of the logical Socket IDs associated with the
// given NUMANode IDs in this CPUDetails.
func (d CPUDetails) SocketsInNUMANodes(ids ...int) cpuset.CPUSet {
	var socketIDs []int
	for _, id := range ids {
		for _, info := range d {
			if info.NUMANodeID == id {
				socketIDs = append(socketIDs, info.SocketID)
			}
		}
	}
	return cpuset.New(socketIDs...)
}

// Cores returns all of the core IDs associated with the CPUs in this
// CPUDetails.
func (d CPUDetails) Cores() cpuset.CPUSet {
	var coreIDs []int
	for _, info := range d {
		coreIDs = append(coreIDs, info.CoreID)
	}
	return cpuset.New(coreIDs...)
}

// CoresInNUMANodes returns all of the core IDs associated with the given
// NUMANode IDs in this CPUDetails.
func (d CPUDetails) CoresInNUMANodes(ids ...int) cpuset.CPUSet {
	var coreIDs []int
	for _, id := range ids {
		for _, info := range d {
			if info.NUMANodeID == id {
				coreIDs = append(coreIDs, info.CoreID)
			}
		}
	}
	return cpuset.New(coreIDs...)
}

// CoresInSockets returns all of the core IDs associated with the given socket
// IDs in this CPUDetails.
func (d CPUDetails) CoresInSockets(ids ...int) cpuset.CPUSet {
	var coreIDs []int
	for _, id := range ids {
		for _, info := range d {
			if info.SocketID == id {
				coreIDs = append(coreIDs, info.CoreID)
			}
		}
	}
	return cpuset.New(coreIDs...)
}

// CPUs returns all of the logical CPU IDs in this CPUDetails.
func (d CPUDetails) CPUs() cpuset.CPUSet {
	var cpuIDs []int
	for cpuID := range d {
		cpuIDs = append(cpuIDs, cpuID)
	}
	return cpuset.New(cpuIDs...)
}

// CPUsInNUMANodes returns all of the logical CPU IDs associated with the given
// NUMANode IDs in this CPUDetails.
func (d CPUDetails) CPUsInNUMANodes(ids ...int) cpuset.CPUSet {
	var cpuIDs []int
	for _, id := range ids {
		for cpu, info := range d {
			if info.NUMANodeID == id {
				cpuIDs = append(cpuIDs, cpu)
			}
		}
	}
	return cpuset.New(cpuIDs...)
}

// CPUsInCores returns all of the logical CPU IDs associated with the given
// core IDs in this CPUDetails.
func (d CPUDetails) CPUsInCores(ids ...int) cpuset.CPUSet {
	var cpuIDs []int
	for _, id := range ids {
		for cpu, info := range d {
			if info.CoreID == id {
				cpuIDs = append(cpuIDs, cpu)
			}
		}
	}
	return cpuset.New(cpuIDs...)
}

// Discover returns CPUTopology based on cadvisor node info
func Discover(machineInfo *cadvisorapi.MachineInfo) (*CPUTopology, error) {
	if machineInfo.NumCores == 0 {
		return nil, fmt.Errorf("could not detect number of cpus")
	}

	CPUDetails := CPUDetails{}
	numPhysicalCores := 0

	for _, node := range machineInfo.Topology {
		numPhysicalCores += len(node.Cores)
		for _, core := range node.Cores {
			if coreID, err := getUniqueCoreID(core.Threads); err == nil {
				for _, cpu := range core.Threads {
					//the last level cache details may not be initialized on some platforms
					if len(core.UncoreCaches) == 0 {
						newCache := cadvisorapi.Cache{}
						newCache.Id = node.Id //the cache ID is global and must be unique, reuse node.Id
						core.UncoreCaches = append(core.UncoreCaches, newCache)
					}
					CPUDetails[cpu] = CPUInfo{
						CoreID:     coreID,
						SocketID:   core.SocketID,
						NUMANodeID: node.Id,
						L3GroupID:  core.UncoreCaches[0].Id,
					}
				}
			} else {
				klog.ErrorS(nil, "Could not get unique coreID for socket", "socket", core.SocketID, "core", core.Id, "threads", core.Threads)
				return nil, err
			}
		}
	}

	return &CPUTopology{
		NumCPUs:      machineInfo.NumCores,
		NumSockets:   machineInfo.NumSockets,
		NumCores:     numPhysicalCores,
		NumNUMANodes: CPUDetails.NUMANodes().Size(),
		CPUDetails:   CPUDetails,
	}, nil
}

// getUniqueCoreID computes coreId as the lowest cpuID
// for a given Threads []int slice. This will assure that coreID's are
// platform unique (opposite to what cAdvisor reports)
func getUniqueCoreID(threads []int) (coreID int, err error) {
	if len(threads) == 0 {
		return 0, fmt.Errorf("no cpus provided")
	}

	if len(threads) != cpuset.New(threads...).Size() {
		return 0, fmt.Errorf("cpus provided are not unique")
	}

	min := threads[0]
	for _, thread := range threads[1:] {
		if thread < min {
			min = thread
		}
	}

	return min, nil
}
