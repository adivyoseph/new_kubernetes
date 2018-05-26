/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package internalversion

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	storage "k8s.io/kubernetes/pkg/apis/storage"
	scheme "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/scheme"
)

// VolumeSnapshotsGetter has a method to return a VolumeSnapshotInterface.
// A group's client should implement this interface.
type VolumeSnapshotsGetter interface {
	VolumeSnapshots(namespace string) VolumeSnapshotInterface
}

// VolumeSnapshotInterface has methods to work with VolumeSnapshot resources.
type VolumeSnapshotInterface interface {
	Create(*storage.VolumeSnapshot) (*storage.VolumeSnapshot, error)
	Update(*storage.VolumeSnapshot) (*storage.VolumeSnapshot, error)
	UpdateStatus(*storage.VolumeSnapshot) (*storage.VolumeSnapshot, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*storage.VolumeSnapshot, error)
	List(opts v1.ListOptions) (*storage.VolumeSnapshotList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *storage.VolumeSnapshot, err error)
	VolumeSnapshotExpansion
}

// volumeSnapshots implements VolumeSnapshotInterface
type volumeSnapshots struct {
	client rest.Interface
	ns     string
}

// newVolumeSnapshots returns a VolumeSnapshots
func newVolumeSnapshots(c *StorageClient, namespace string) *volumeSnapshots {
	return &volumeSnapshots{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the volumeSnapshot, and returns the corresponding volumeSnapshot object, and an error if there is any.
func (c *volumeSnapshots) Get(name string, options v1.GetOptions) (result *storage.VolumeSnapshot, err error) {
	result = &storage.VolumeSnapshot{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("volumesnapshots").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of VolumeSnapshots that match those selectors.
func (c *volumeSnapshots) List(opts v1.ListOptions) (result *storage.VolumeSnapshotList, err error) {
	result = &storage.VolumeSnapshotList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("volumesnapshots").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested volumeSnapshots.
func (c *volumeSnapshots) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("volumesnapshots").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a volumeSnapshot and creates it.  Returns the server's representation of the volumeSnapshot, and an error, if there is any.
func (c *volumeSnapshots) Create(volumeSnapshot *storage.VolumeSnapshot) (result *storage.VolumeSnapshot, err error) {
	result = &storage.VolumeSnapshot{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("volumesnapshots").
		Body(volumeSnapshot).
		Do().
		Into(result)
	return
}

// Update takes the representation of a volumeSnapshot and updates it. Returns the server's representation of the volumeSnapshot, and an error, if there is any.
func (c *volumeSnapshots) Update(volumeSnapshot *storage.VolumeSnapshot) (result *storage.VolumeSnapshot, err error) {
	result = &storage.VolumeSnapshot{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("volumesnapshots").
		Name(volumeSnapshot.Name).
		Body(volumeSnapshot).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *volumeSnapshots) UpdateStatus(volumeSnapshot *storage.VolumeSnapshot) (result *storage.VolumeSnapshot, err error) {
	result = &storage.VolumeSnapshot{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("volumesnapshots").
		Name(volumeSnapshot.Name).
		SubResource("status").
		Body(volumeSnapshot).
		Do().
		Into(result)
	return
}

// Delete takes name of the volumeSnapshot and deletes it. Returns an error if one occurs.
func (c *volumeSnapshots) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("volumesnapshots").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *volumeSnapshots) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("volumesnapshots").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched volumeSnapshot.
func (c *volumeSnapshots) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *storage.VolumeSnapshot, err error) {
	result = &storage.VolumeSnapshot{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("volumesnapshots").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
