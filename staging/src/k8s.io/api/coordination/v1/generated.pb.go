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

// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: k8s.io/api/coordination/v1/generated.proto

package v1

import (
	fmt "fmt"

	io "io"

	proto "github.com/gogo/protobuf/proto"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

func (m *Lease) Reset()      { *m = Lease{} }
func (*Lease) ProtoMessage() {}
func (*Lease) Descriptor() ([]byte, []int) {
	return fileDescriptor_239d5a4df3139dce, []int{0}
}
func (m *Lease) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Lease) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *Lease) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Lease.Merge(m, src)
}
func (m *Lease) XXX_Size() int {
	return m.Size()
}
func (m *Lease) XXX_DiscardUnknown() {
	xxx_messageInfo_Lease.DiscardUnknown(m)
}

var xxx_messageInfo_Lease proto.InternalMessageInfo

func (m *LeaseList) Reset()      { *m = LeaseList{} }
func (*LeaseList) ProtoMessage() {}
func (*LeaseList) Descriptor() ([]byte, []int) {
	return fileDescriptor_239d5a4df3139dce, []int{1}
}
func (m *LeaseList) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LeaseList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *LeaseList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LeaseList.Merge(m, src)
}
func (m *LeaseList) XXX_Size() int {
	return m.Size()
}
func (m *LeaseList) XXX_DiscardUnknown() {
	xxx_messageInfo_LeaseList.DiscardUnknown(m)
}

var xxx_messageInfo_LeaseList proto.InternalMessageInfo

func (m *LeaseSpec) Reset()      { *m = LeaseSpec{} }
func (*LeaseSpec) ProtoMessage() {}
func (*LeaseSpec) Descriptor() ([]byte, []int) {
	return fileDescriptor_239d5a4df3139dce, []int{2}
}
func (m *LeaseSpec) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LeaseSpec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *LeaseSpec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LeaseSpec.Merge(m, src)
}
func (m *LeaseSpec) XXX_Size() int {
	return m.Size()
}
func (m *LeaseSpec) XXX_DiscardUnknown() {
	xxx_messageInfo_LeaseSpec.DiscardUnknown(m)
}

var xxx_messageInfo_LeaseSpec proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Lease)(nil), "k8s.io.api.coordination.v1.Lease")
	proto.RegisterType((*LeaseList)(nil), "k8s.io.api.coordination.v1.LeaseList")
	proto.RegisterType((*LeaseSpec)(nil), "k8s.io.api.coordination.v1.LeaseSpec")
}

func init() {
	proto.RegisterFile("k8s.io/api/coordination/v1/generated.proto", fileDescriptor_239d5a4df3139dce)
}

var fileDescriptor_239d5a4df3139dce = []byte{
	// 558 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x94, 0x4f, 0x6f, 0xd3, 0x30,
	0x18, 0xc6, 0x9b, 0xb5, 0x45, 0xad, 0xcb, 0x46, 0x65, 0x7a, 0x88, 0x7a, 0x48, 0x46, 0x11, 0xd2,
	0x84, 0x84, 0x43, 0x27, 0x84, 0x10, 0x17, 0x46, 0x40, 0xc0, 0xa4, 0x4e, 0x48, 0xe9, 0x4e, 0x68,
	0x07, 0xdc, 0xc4, 0xa4, 0xa6, 0x4b, 0x1c, 0x62, 0xb7, 0xa8, 0x37, 0x3e, 0x02, 0x5f, 0x05, 0x3e,
	0x45, 0x8f, 0x3b, 0x4e, 0x1c, 0x22, 0x6a, 0xbe, 0x05, 0x27, 0x64, 0xf7, 0x2f, 0xfd, 0xa3, 0x4d,
	0xbb, 0xc5, 0xaf, 0xdf, 0xe7, 0xf7, 0x3e, 0x7e, 0x9c, 0x04, 0x3c, 0xec, 0x3d, 0xe3, 0x88, 0x32,
	0x07, 0x27, 0xd4, 0xf1, 0x19, 0x4b, 0x03, 0x1a, 0x63, 0x41, 0x59, 0xec, 0x0c, 0x9a, 0x4e, 0x48,
	0x62, 0x92, 0x62, 0x41, 0x02, 0x94, 0xa4, 0x4c, 0x30, 0x58, 0x9f, 0xf4, 0x22, 0x9c, 0x50, 0xb4,
	0xdc, 0x8b, 0x06, 0xcd, 0xfa, 0xa3, 0x90, 0x8a, 0x6e, 0xbf, 0x83, 0x7c, 0x16, 0x39, 0x21, 0x0b,
	0x99, 0xa3, 0x25, 0x9d, 0xfe, 0x27, 0xbd, 0xd2, 0x0b, 0xfd, 0x34, 0x41, 0xd5, 0x9f, 0x2c, 0xc6,
	0x46, 0xd8, 0xef, 0xd2, 0x98, 0xa4, 0x43, 0x27, 0xe9, 0x85, 0xaa, 0xc0, 0x9d, 0x88, 0x08, 0xbc,
	0xc1, 0x40, 0xdd, 0xd9, 0xa6, 0x4a, 0xfb, 0xb1, 0xa0, 0x11, 0x59, 0x13, 0x3c, 0xbd, 0x4a, 0xc0,
	0xfd, 0x2e, 0x89, 0xf0, 0xaa, 0xae, 0xf1, 0xd3, 0x00, 0xc5, 0x16, 0xc1, 0x9c, 0xc0, 0x8f, 0xa0,
	0xa4, 0xdc, 0x04, 0x58, 0x60, 0xd3, 0xd8, 0x37, 0x0e, 0x2a, 0x87, 0x8f, 0xd1, 0x22, 0x86, 0x39,
	0x14, 0x25, 0xbd, 0x50, 0x15, 0x38, 0x52, 0xdd, 0x68, 0xd0, 0x44, 0xef, 0x3b, 0x9f, 0x89, 0x2f,
	0x4e, 0x88, 0xc0, 0x2e, 0x1c, 0x65, 0x76, 0x4e, 0x66, 0x36, 0x58, 0xd4, 0xbc, 0x39, 0x15, 0xbe,
	0x05, 0x05, 0x9e, 0x10, 0xdf, 0xdc, 0xd1, 0xf4, 0x07, 0x68, 0x7b, 0xc8, 0x48, 0x5b, 0x6a, 0x27,
	0xc4, 0x77, 0x6f, 0x4f, 0x91, 0x05, 0xb5, 0xf2, 0x34, 0xa0, 0xf1, 0xc3, 0x00, 0x65, 0xdd, 0xd1,
	0xa2, 0x5c, 0xc0, 0xb3, 0x35, 0xe3, 0xe8, 0x7a, 0xc6, 0x95, 0x5a, 0xdb, 0xae, 0x4e, 0x67, 0x94,
	0x66, 0x95, 0x25, 0xd3, 0x6f, 0x40, 0x91, 0x0a, 0x12, 0x71, 0x73, 0x67, 0x3f, 0x7f, 0x50, 0x39,
	0xbc, 0x77, 0xa5, 0x6b, 0x77, 0x77, 0x4a, 0x2b, 0x1e, 0x2b, 0x9d, 0x37, 0x91, 0x37, 0x7e, 0xe5,
	0xa7, 0x9e, 0xd5, 0x39, 0xe0, 0x73, 0xb0, 0xd7, 0x65, 0xe7, 0x01, 0x49, 0x8f, 0x03, 0x12, 0x0b,
	0x2a, 0x86, 0xda, 0x79, 0xd9, 0x85, 0x32, 0xb3, 0xf7, 0xde, 0xfd, 0xb7, 0xe3, 0xad, 0x74, 0xc2,
	0x16, 0xa8, 0x9d, 0x2b, 0xd0, 0xeb, 0x7e, 0xaa, 0x27, 0xb7, 0x89, 0xcf, 0xe2, 0x80, 0xeb, 0x58,
	0x8b, 0xae, 0x29, 0x33, 0xbb, 0xd6, 0xda, 0xb0, 0xef, 0x6d, 0x54, 0xc1, 0x0e, 0xa8, 0x60, 0xff,
	0x4b, 0x9f, 0xa6, 0xe4, 0x94, 0x46, 0xc4, 0xcc, 0xeb, 0x00, 0x9d, 0xeb, 0x05, 0x78, 0x42, 0xfd,
	0x94, 0x29, 0x99, 0x7b, 0x47, 0x66, 0x76, 0xe5, 0xe5, 0x82, 0xe3, 0x2d, 0x43, 0xe1, 0x19, 0x28,
	0xa7, 0x24, 0x26, 0x5f, 0xf5, 0x84, 0xc2, 0xcd, 0x26, 0xec, 0xca, 0xcc, 0x2e, 0x7b, 0x33, 0x8a,
	0xb7, 0x00, 0xc2, 0x23, 0x50, 0xd5, 0x27, 0x3b, 0x4d, 0x71, 0xcc, 0xa9, 0x3a, 0x1b, 0x37, 0x8b,
	0x3a, 0x8b, 0x9a, 0xcc, 0xec, 0x6a, 0x6b, 0x65, 0xcf, 0x5b, 0xeb, 0x86, 0x2f, 0x40, 0x89, 0x0b,
	0xf5, 0x55, 0x84, 0x43, 0xf3, 0x96, 0xbe, 0x87, 0xfb, 0xea, 0x6d, 0x68, 0x4f, 0x6b, 0x7f, 0x33,
	0xfb, 0xee, 0xab, 0xd9, 0x55, 0x93, 0x60, 0x56, 0xf6, 0xe6, 0x22, 0xf7, 0x68, 0x34, 0xb6, 0x72,
	0x17, 0x63, 0x2b, 0x77, 0x39, 0xb6, 0x72, 0xdf, 0xa4, 0x65, 0x8c, 0xa4, 0x65, 0x5c, 0x48, 0xcb,
	0xb8, 0x94, 0x96, 0xf1, 0x5b, 0x5a, 0xc6, 0xf7, 0x3f, 0x56, 0xee, 0x43, 0x7d, 0xfb, 0x1f, 0xe8,
	0x5f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf7, 0x27, 0xa2, 0xcc, 0x9e, 0x04, 0x00, 0x00,
}

func (m *Lease) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Lease) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Lease) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Spec.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenerated(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size, err := m.ObjectMeta.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenerated(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *LeaseList) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LeaseList) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LeaseList) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Items) > 0 {
		for iNdEx := len(m.Items) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Items[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenerated(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size, err := m.ListMeta.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenerated(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *LeaseSpec) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LeaseSpec) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LeaseSpec) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Strategy != nil {
		i -= len(*m.Strategy)
		copy(dAtA[i:], *m.Strategy)
		i = encodeVarintGenerated(dAtA, i, uint64(len(*m.Strategy)))
		i--
		dAtA[i] = 0x32
	}
	if m.LeaseTransitions != nil {
		i = encodeVarintGenerated(dAtA, i, uint64(*m.LeaseTransitions))
		i--
		dAtA[i] = 0x28
	}
	if m.RenewTime != nil {
		{
			size, err := m.RenewTime.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenerated(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.AcquireTime != nil {
		{
			size, err := m.AcquireTime.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenerated(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.LeaseDurationSeconds != nil {
		i = encodeVarintGenerated(dAtA, i, uint64(*m.LeaseDurationSeconds))
		i--
		dAtA[i] = 0x10
	}
	if m.HolderIdentity != nil {
		i -= len(*m.HolderIdentity)
		copy(dAtA[i:], *m.HolderIdentity)
		i = encodeVarintGenerated(dAtA, i, uint64(len(*m.HolderIdentity)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenerated(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenerated(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Lease) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.ObjectMeta.Size()
	n += 1 + l + sovGenerated(uint64(l))
	l = m.Spec.Size()
	n += 1 + l + sovGenerated(uint64(l))
	return n
}

func (m *LeaseList) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.ListMeta.Size()
	n += 1 + l + sovGenerated(uint64(l))
	if len(m.Items) > 0 {
		for _, e := range m.Items {
			l = e.Size()
			n += 1 + l + sovGenerated(uint64(l))
		}
	}
	return n
}

func (m *LeaseSpec) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.HolderIdentity != nil {
		l = len(*m.HolderIdentity)
		n += 1 + l + sovGenerated(uint64(l))
	}
	if m.LeaseDurationSeconds != nil {
		n += 1 + sovGenerated(uint64(*m.LeaseDurationSeconds))
	}
	if m.AcquireTime != nil {
		l = m.AcquireTime.Size()
		n += 1 + l + sovGenerated(uint64(l))
	}
	if m.RenewTime != nil {
		l = m.RenewTime.Size()
		n += 1 + l + sovGenerated(uint64(l))
	}
	if m.LeaseTransitions != nil {
		n += 1 + sovGenerated(uint64(*m.LeaseTransitions))
	}
	if m.Strategy != nil {
		l = len(*m.Strategy)
		n += 1 + l + sovGenerated(uint64(l))
	}
	return n
}

func sovGenerated(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenerated(x uint64) (n int) {
	return sovGenerated(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Lease) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Lease{`,
		`ObjectMeta:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.ObjectMeta), "ObjectMeta", "v1.ObjectMeta", 1), `&`, ``, 1) + `,`,
		`Spec:` + strings.Replace(strings.Replace(this.Spec.String(), "LeaseSpec", "LeaseSpec", 1), `&`, ``, 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *LeaseList) String() string {
	if this == nil {
		return "nil"
	}
	repeatedStringForItems := "[]Lease{"
	for _, f := range this.Items {
		repeatedStringForItems += strings.Replace(strings.Replace(f.String(), "Lease", "Lease", 1), `&`, ``, 1) + ","
	}
	repeatedStringForItems += "}"
	s := strings.Join([]string{`&LeaseList{`,
		`ListMeta:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.ListMeta), "ListMeta", "v1.ListMeta", 1), `&`, ``, 1) + `,`,
		`Items:` + repeatedStringForItems + `,`,
		`}`,
	}, "")
	return s
}
func (this *LeaseSpec) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&LeaseSpec{`,
		`HolderIdentity:` + valueToStringGenerated(this.HolderIdentity) + `,`,
		`LeaseDurationSeconds:` + valueToStringGenerated(this.LeaseDurationSeconds) + `,`,
		`AcquireTime:` + strings.Replace(fmt.Sprintf("%v", this.AcquireTime), "MicroTime", "v1.MicroTime", 1) + `,`,
		`RenewTime:` + strings.Replace(fmt.Sprintf("%v", this.RenewTime), "MicroTime", "v1.MicroTime", 1) + `,`,
		`LeaseTransitions:` + valueToStringGenerated(this.LeaseTransitions) + `,`,
		`Strategy:` + valueToStringGenerated(this.Strategy) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringGenerated(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *Lease) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Lease: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Lease: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ObjectMeta", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ObjectMeta.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Spec", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Spec.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenerated(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenerated
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *LeaseList) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: LeaseList: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LeaseList: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ListMeta", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ListMeta.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Items", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Items = append(m.Items, Lease{})
			if err := m.Items[len(m.Items)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenerated(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenerated
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *LeaseSpec) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: LeaseSpec: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LeaseSpec: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HolderIdentity", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.HolderIdentity = &s
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LeaseDurationSeconds", wireType)
			}
			var v int32
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.LeaseDurationSeconds = &v
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AcquireTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.AcquireTime == nil {
				m.AcquireTime = &v1.MicroTime{}
			}
			if err := m.AcquireTime.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RenewTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.RenewTime == nil {
				m.RenewTime = &v1.MicroTime{}
			}
			if err := m.RenewTime.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LeaseTransitions", wireType)
			}
			var v int32
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.LeaseTransitions = &v
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Strategy", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := CoordinatedStrategy(dAtA[iNdEx:postIndex])
			m.Strategy = &s
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenerated(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenerated
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipGenerated(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthGenerated
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenerated
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenerated
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenerated        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenerated          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenerated = fmt.Errorf("proto: unexpected end of group")
)
