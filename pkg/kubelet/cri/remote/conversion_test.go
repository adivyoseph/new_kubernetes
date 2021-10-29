/*
Copyright 2021 The Kubernetes Authors.

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

package remote

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
	"k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

func fillFields(s interface{}) {
	fillFieldsOffset(s, 0)
}

func fillFieldsOffset(s interface{}, offset int) {
	reflectType := reflect.TypeOf(s).Elem()
	reflectValue := reflect.ValueOf(s).Elem()

	for i := 0; i < reflectType.NumField(); i++ {
		field := reflectValue.Field(i)
		typeName := reflectType.Field(i).Name

		// Skipping protobuf internal values
		if strings.HasPrefix(typeName, "XXX_") {
			continue
		}

		fillField(field, i+offset)
	}
}

func fillField(field reflect.Value, v int) {
	switch field.Kind() {
	case reflect.Bool:
		field.SetBool(true)

	case reflect.Float32, reflect.Float64:
		field.SetFloat(float64(v))

	case reflect.String:
		field.SetString(fmt.Sprint(v))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		field.SetInt(int64(v))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		field.SetUint(uint64(v))

	case reflect.Map:
		field.Set(reflect.MakeMap(field.Type()))

	case reflect.Array, reflect.Slice:
		slice := reflect.MakeSlice(field.Type(), 1, 1)
		field.Set(slice)
		first := slice.Index(0)

		if first.Type().Kind() == reflect.Ptr {
			first.Set(reflect.New(first.Type().Elem()))
			fillFieldsOffset(first.Interface(), v)
		} else {
			fillField(first, v)
		}

	case reflect.Ptr:
		val := reflect.New(field.Type().Elem())
		field.Set(val)
		fillFieldsOffset(field.Interface(), v)

	case reflect.Struct:
		fillFieldsOffset(field.Addr().Interface(), v)
	}

}

func assertEqual(t *testing.T, a, b proto.Message) {
	aBytes, err := proto.Marshal(a)
	assert.Nil(t, err)

	bBytes, err := proto.Marshal(b)
	assert.Nil(t, err)

	assert.Equal(t, aBytes, bBytes)
}

func TestFromV1alpha2VersionResponse(t *testing.T) {
	from := &v1alpha2.VersionResponse{}
	fillFields(from)
	to := fromV1alpha2VersionResponse(from)
	assertEqual(t, from, to)
}

func TestFromV1alpha2PodSandboxStatus(t *testing.T) {
	from := &v1alpha2.PodSandboxStatus{}
	fillFields(from)
	to := fromV1alpha2PodSandboxStatus(from)
	assertEqual(t, from, to)
}

func TestFromV1alpha2ListPodSandboxResponse(t *testing.T) {
	from := &v1alpha2.ListPodSandboxResponse{}
	fillFields(from)
	to := fromV1alpha2ListPodSandboxResponse(from)
	assertEqual(t, from, to)
}

func TestFromV1alpha2ListContainersResponse(t *testing.T) {
	from := &v1alpha2.ListContainersResponse{}
	fillFields(from)
	to := fromV1alpha2ListContainersResponse(from)
	assertEqual(t, from, to)
}

func TestFromV1alpha2ContainerStatus(t *testing.T) {
	from := &v1alpha2.ContainerStatus{}
	fillFields(from)
	to := fromV1alpha2ContainerStatus(from)
	assertEqual(t, from, to)
}

func TestFromV1alpha2ExecResponse(t *testing.T) {
	from := &v1alpha2.ExecResponse{}
	fillFields(from)
	to := fromV1alpha2ExecResponse(from)
	assertEqual(t, from, to)
}

func TestFromV1alpha2AttachResponse(t *testing.T) {
	from := &v1alpha2.AttachResponse{}
	fillFields(from)
	to := fromV1alpha2AttachResponse(from)
	assertEqual(t, from, to)
}

func TestFromV1alpha2PortForwardResponse(t *testing.T) {
	from := &v1alpha2.PortForwardResponse{}
	fillFields(from)
	to := fromV1alpha2PortForwardResponse(from)
	assertEqual(t, from, to)
}

func TestFromV1alpha2RuntimeStatus(t *testing.T) {
	from := &v1alpha2.RuntimeStatus{}
	fillFields(from)
	to := fromV1alpha2RuntimeStatus(from)
	assertEqual(t, from, to)
}

func TestFromV1alpha2ContainerStats(t *testing.T) {
	from := &v1alpha2.ContainerStats{}
	fillFields(from)
	to := fromV1alpha2ContainerStats(from)
	assertEqual(t, from, to)
}

func TestFromV1alpha2ImageFsInfoResponse(t *testing.T) {
	from := &v1alpha2.ImageFsInfoResponse{}
	fillFields(from)
	to := fromV1alpha2ImageFsInfoResponse(from)
	assertEqual(t, from, to)
}

func TestFromV1alpha2ListContainerStatsResponse(t *testing.T) {
	from := &v1alpha2.ListContainerStatsResponse{}
	fillFields(from)
	to := fromV1alpha2ListContainerStatsResponse(from)
	assertEqual(t, from, to)
}

func TestFromV1alpha2PodSandboxStats(t *testing.T) {
	from := &v1alpha2.PodSandboxStats{}
	fillFields(from)
	to := fromV1alpha2PodSandboxStats(from)
	assertEqual(t, from, to)
}

func TestFromV1alpha2ListPodSandboxStatsResponse(t *testing.T) {
	from := &v1alpha2.ListPodSandboxStatsResponse{}
	fillFields(from)
	to := fromV1alpha2ListPodSandboxStatsResponse(from)
	assertEqual(t, from, to)
}

func TestFromV1alpha2Image(t *testing.T) {
	from := &v1alpha2.Image{}
	fillFields(from)
	to := fromV1alpha2Image(from)
	fmt.Printf(":%+v\n", to)
	assertEqual(t, from, to)
}

func TestFromV1alpha2ListImagesResponse(t *testing.T) {
	from := &v1alpha2.ListImagesResponse{}
	fillFields(from)
	to := fromV1alpha2ListImagesResponse(from)
	assertEqual(t, from, to)
}

func TestV1alpha2PodSandboxConfig(t *testing.T) {
	from := &runtimeapi.PodSandboxConfig{}
	fillFields(from)
	to := v1alpha2PodSandboxConfig(from)
	assertEqual(t, from, to)
}

func TestV1alpha2PodSandboxFilter(t *testing.T) {
	from := &runtimeapi.PodSandboxFilter{}
	fillFields(from)
	to := v1alpha2PodSandboxFilter(from)
	assertEqual(t, from, to)
}

func TestV1alpha2ContainerConfig(t *testing.T) {
	from := &runtimeapi.ContainerConfig{}
	fillFields(from)
	to := v1alpha2ContainerConfig(from)
	assertEqual(t, from, to)
}

func TestV1alpha2ContainerFilter(t *testing.T) {
	from := &runtimeapi.ContainerFilter{}
	fillFields(from)
	to := v1alpha2ContainerFilter(from)
	assertEqual(t, from, to)
}

func TestV1alpha2LinuxContainerResources(t *testing.T) {
	from := &runtimeapi.LinuxContainerResources{}
	fillFields(from)
	to := v1alpha2LinuxContainerResources(from)
	assertEqual(t, from, to)
}

func TestV1alpha2ExecRequest(t *testing.T) {
	from := &runtimeapi.ExecRequest{}
	fillFields(from)
	to := v1alpha2ExecRequest(from)
	assertEqual(t, from, to)
}

func TestV1alpha2AttachRequest(t *testing.T) {
	from := &runtimeapi.AttachRequest{}
	fillFields(from)
	to := v1alpha2AttachRequest(from)
	assertEqual(t, from, to)
}

func TestV1alpha2PortForwardRequest(t *testing.T) {
	from := &runtimeapi.PortForwardRequest{}
	fillFields(from)
	to := v1alpha2PortForwardRequest(from)
	assertEqual(t, from, to)
}

func TestV1alpha2RuntimeConfig(t *testing.T) {
	from := &runtimeapi.RuntimeConfig{}
	fillFields(from)
	to := v1alpha2RuntimeConfig(from)
	assertEqual(t, from, to)
}

func TestV1alpha2ContainerStatsFilter(t *testing.T) {
	from := &runtimeapi.ContainerStatsFilter{}
	fillFields(from)
	to := v1alpha2ContainerStatsFilter(from)
	assertEqual(t, from, to)
}

func TestV1alpha2PodSandboxStatsFilter(t *testing.T) {
	from := &runtimeapi.PodSandboxStatsFilter{}
	fillFields(from)
	to := v1alpha2PodSandboxStatsFilter(from)
	assertEqual(t, from, to)
}

func TestV1alpha2ImageFilter(t *testing.T) {
	from := &runtimeapi.ImageFilter{}
	fillFields(from)
	to := v1alpha2ImageFilter(from)
	assertEqual(t, from, to)
}

func TestV1alpha2ImageSpec(t *testing.T) {
	from := &runtimeapi.ImageSpec{}
	fillFields(from)
	to := v1alpha2ImageSpec(from)
	assertEqual(t, from, to)
}

func TestV1alpha2AuthConfig(t *testing.T) {
	from := &runtimeapi.AuthConfig{}
	fillFields(from)
	to := v1alpha2AuthConfig(from)
	assertEqual(t, from, to)
}
