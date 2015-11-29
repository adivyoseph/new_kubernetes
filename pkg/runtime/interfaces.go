/*
Copyright 2014 The Kubernetes Authors All rights reserved.

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

package runtime

import (
	"io"
	"net/url"

	"k8s.io/kubernetes/pkg/api/unversioned"
)

const (
	APIVersionInternal    = "__internal"
	APIVersionUnversioned = "__unversioned"
)

// Typer retrieves information about an object's group, version, and kind.
type Typer interface {
	// ObjectVersionAndKind returns the version and kind of the provided object, or an
	// error if the object is not recognized (IsNotRegisteredError will return true).
	ObjectVersionAndKind(Object) (*unversioned.GroupVersionKind, error)
}

type Encoder interface {
	// EncodeToStream writes an object to a stream.
	EncodeToStream(obj Object, stream io.Writer) error
}

type Decoder interface {
	// Decode attempts to deserialize the provided data using either the innate typing of the scheme or the
	// default kind, group, and version provided. It returns a decoded object as well as the actual kind, group, and
	// version decoded, or an error.
	Decode(data []byte, gvk *unversioned.GroupVersionKind) (Object, *unversioned.GroupVersionKind, error)
}

// Serializer is the core interface for transforming objects into a serialized format and back.
// Implementations may choose to perform conversion of the object, but no assumptions should be made.
type Serializer interface {
	Encoder
	Decoder
}

// Codec is a Serializer that deals with the details of versioning objects. It offers the same
// interface as Serializer, so this is a marker to consumers that care about the version of the objects
// they receive.
type Codec Serializer

// ParameterCodec defines methods for serializing and deserializing API objects to url.Values
type ParameterSerializer interface {
	DecodeParameters(parameters url.Values, gvk unversioned.GroupVersionKind) (Object, error)
	EncodeParameters(Object) (url.Values, error)
}

type ParameterCodec ParameterSerializer

// ObjectDecoder is a convenience interface for identifying serialized versions of objects
// and transforming them into Objects. It intentionally overlaps with ObjectTyper and
// Decoder for use in decode only paths.
// TODO: Consider removing this interface?
type ObjectDecoder interface {
	// Identical to Serializer#Decode
	Decode(data []byte, gvk *unversioned.GroupVersionKind) (Object, *unversioned.GroupVersionKind, error)
	// Recognizes returns true if the scheme is able to handle the provided version and kind
	// of an object.
	Recognizes(version, kind string) bool
}

///////////////////////////////////////////////////////////////////////////////
// Non-codec interfaces

type ObjectVersioner interface {
	ConvertToVersion(in Object, outVersion string) (out Object, err error)
}

// ObjectConvertor converts an object to a different version.
type ObjectConvertor interface {
	Convert(in, out interface{}) error
	ConvertToVersion(in Object, outVersion string) (out Object, err error)
	ConvertFieldLabel(version, kind, label, value string) (string, string, error)
}

// ObjectTyper contains methods for extracting the APIVersion and Kind
// of objects.
type ObjectTyper interface {
	// ObjectVersionAndKind returns the version and kind of the provided object, or an
	// error if the object is not recognized (IsNotRegisteredError will return true).
	ObjectVersionAndKind(Object) (version, kind string, err error)
	// Recognizes returns true if the scheme is able to handle the provided version and kind,
	// or more precisely that the provided version is a possible conversion or decoding
	// target.
	Recognizes(version, kind string) bool
}

// ObjectCreater contains methods for instantiating an object by kind and version.
type ObjectCreater interface {
	New(version, kind string) (out Object, err error)
}

// ObjectCopier duplicates an object.
type ObjectCopier interface {
	// Copy returns an exact copy of the provided Object, or an error if the
	// copy could not be completed.
	Copy(Object) (Object, error)
}

// ResourceVersioner provides methods for setting and retrieving
// the resource version from an API object.
type ResourceVersioner interface {
	SetResourceVersion(obj Object, version string) error
	ResourceVersion(obj Object) (string, error)
}

// SelfLinker provides methods for setting and retrieving the SelfLink field of an API object.
type SelfLinker interface {
	SetSelfLink(obj Object, selfLink string) error
	SelfLink(obj Object) (string, error)

	// Knowing Name is sometimes necessary to use a SelfLinker.
	Name(obj Object) (string, error)
	// Knowing Namespace is sometimes necessary to use a SelfLinker
	Namespace(obj Object) (string, error)
}

// All api types must support the Object interface. It's deliberately tiny so that this is not an onerous
// burden. Implement it with a pointer receiver; this will allow us to use the go compiler to check the
// one thing about our objects that it's capable of checking for us.
type Object interface {
	// This function is used only to enforce membership. It's never called.
	// TODO: Consider mass rename in the future to make it do something useful.
	IsAnAPIObject()
}
