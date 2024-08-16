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
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strings"

	"k8s.io/apimachinery/pkg/util/dump"
	"k8s.io/kubernetes/pkg/kubelet/checkpointmanager"
	"k8s.io/kubernetes/pkg/kubelet/checkpointmanager/checksum"
	"k8s.io/kubernetes/pkg/kubelet/checkpointmanager/errors"
)


var _ checkpointmanager.Checkpoint = &CPUManagerCheckpoint{}

type CPUManagerCheckpoint struct {
	ReservedCPUSet    string                   `json:"reservedCpuSet"`
	AvaiableCPUSet    string                   `json:"availableCpuSet"`
	StaticContainers  map[string]map[string]string `json:"staticContainers,omitempty"`
	NormalContainers  map[string]map[string]string `json:"normalContainers,omitempty"`
	Checksum      checksum.Checksum            `json:"checksum"`
}


// NewCPUManagerCheckpoint returns an instance of Checkpoint
func NewCPUManagerCheckpoint() *CPUManagerCheckpoint {
	return &CPUManagerCheckpoint{
		StaticContainers: make(map[string]map[string]string),
		NormalContainers: make(map[string]map[string]string),
	}
}


// MarshalCheckpoint returns marshalled checkpoint 
func (cp *CPUManagerCheckpoint) MarshalCheckpoint() ([]byte, error) {
	// make sure checksum wasn't set before so it doesn't affect output checksum
	cp.Checksum = 0
	cp.Checksum = checksum.New(cp)
	return json.Marshal(*cp)
}


// UnmarshalCheckpoint tries to unmarshal passed bytes to checkpoint 
func (cp *CPUManagerCheckpoint) UnmarshalCheckpoint(blob []byte) error {
	return json.Unmarshal(blob, cp)
}

// VerifyChecksum verifies that current checksum of checkpoint is valid 
func (cp *CPUManagerCheckpoint) VerifyChecksum() error {
	if cp.Checksum == 0 {
		// accept empty checksum for compatibility with old file backend
		return nil
	}
	ck := cp.Checksum
	cp.Checksum = 0
	err := ck.Verify(cp)
	cp.Checksum = ck
	return err
}
