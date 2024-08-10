package state

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/kubernetes/pkg/kubelet/checkpointmanager"
)

func Test_stateCheckpoint_storeState(t *testing.T) {
	cache := NewStateMemory()
	checkpointManager, _ := checkpointmanager.NewCheckpointManager(".")
	checkpointName := "pod_state_checkpoint"
	sc := &stateCheckpoint{
		cache:             cache,
		checkpointManager: checkpointManager,
		checkpointName:    checkpointName,
	}
	var err error
	type args struct {
		podResourceAllocation PodResourceAllocation
	}

	tests := []struct {
		name string
		args args
	}{}
	suffix := []string{"Ki", "Mi", "Gi", "Ti", "Pi", "Ei", "n", "u", "m", "k", "M", "G", "T", "P", "E", ""}
	factor := []string{"1", "0.1", "0.03", "10", "100", "512", "1000", "1024", "700", "10000"}
	for j, fact := range factor {
		for i, suf := range suffix {
			if (suf == "E" || suf == "Ei") && (fact == "1000" || fact == "10000") {
				// when fact is 1000 or 10000, suffix "E" or "Ei", the quantity value is overflow
				// see detail https://github.com/kubernetes/apimachinery/blob/95b78024e3feada7739b40426690b4f287933fd8/pkg/api/resource/quantity.go#L301
				continue
			}
			tests = append(tests, struct {
				name string
				args args
			}{
				name: fmt.Sprintf("format %d - cpu %s%s", j*len(suffix)+i, fact, suf),
				args: args{
					podResourceAllocation: PodResourceAllocation{
						"pod1": {
							"container1": {
								v1.ResourceCPU:    resource.MustParse(fmt.Sprintf("%s%s", fact, suf)),
								v1.ResourceMemory: resource.MustParse(fmt.Sprintf("%s%s", fact, suf)),
							},
						},
					},
				},
			})
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// reset environment
			err = sc.cache.ClearState()
			require.NoError(t, err, "failed to clear state")
			// remove file ./pod_state_checkpoint
			err = checkpointManager.RemoveCheckpoint(checkpointName)
			require.NoError(t, err, "failed to remove checkpoint")
			err = sc.cache.SetPodResourceAllocation(tt.args.podResourceAllocation)
			require.NoError(t, err, "failed to set pod resource allocation")
			// store state
			err = sc.storeState()
			if err != nil {
				t.Errorf("failed to store state: %v", err)
			}
			// deep copy cache
			originCache := NewStateMemory()
			podAllocation := sc.cache.GetPodResourceAllocation()
			originCache.SetPodResourceAllocation(podAllocation)
			podResizeStatus := sc.cache.GetResizeStatus()
			originCache.SetResizeStatus(podResizeStatus)

			// reset cache
			err = sc.cache.ClearState()
			require.NoError(t, err, "failed to clear state")
			// restore state
			err = sc.restoreState()
			require.NoError(t, err, "failed to restore state")
			// compare restored state with original state
			restoredCache := sc.cache
			require.Equal(t, len(originCache.GetPodResourceAllocation()), len(restoredCache.GetPodResourceAllocation()), "restored pod resource allocation is not equal to original pod resource allocation")
			for podUID, containerResourceList := range originCache.GetPodResourceAllocation() {
				require.Equal(t, len(containerResourceList), len(restoredCache.GetPodResourceAllocation()[podUID]), "restored pod resource allocation is not equal to original pod resource allocation")
				for containerName, resourceList := range containerResourceList {
					for name, quantity := range resourceList {
						if !quantity.Equal(restoredCache.GetPodResourceAllocation()[podUID][containerName][name]) {
							t.Errorf("restored pod resource allocation is not equal to original pod resource allocation")
						}
					}
				}
			}
			// remove file ./pod_state_checkpoint
			err = checkpointManager.RemoveCheckpoint(checkpointName)
			require.NoError(t, err, "failed to remove checkpoint")
		})
	}
}
