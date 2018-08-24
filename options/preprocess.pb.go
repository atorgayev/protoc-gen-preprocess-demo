// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: options/preprocess.proto

/*
Package preprocess is a generated protocol buffer package.

It is generated from these files:
	options/preprocess.proto

It has these top-level messages:
	PreprocessFieldOptions
	PreprocessString
*/
package preprocess

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type PreprocessFieldOptions struct {
	String_          *PreprocessString `protobuf:"bytes,1,opt,name=string" json:"string,omitempty"`
	Test             *string           `protobuf:"bytes,2,opt,name=test" json:"test,omitempty"`
	XXX_unrecognized []byte            `json:"-"`
}

func (m *PreprocessFieldOptions) Reset()                    { *m = PreprocessFieldOptions{} }
func (m *PreprocessFieldOptions) String() string            { return proto.CompactTextString(m) }
func (*PreprocessFieldOptions) ProtoMessage()               {}
func (*PreprocessFieldOptions) Descriptor() ([]byte, []int) { return fileDescriptorPreprocess, []int{0} }

func (m *PreprocessFieldOptions) GetString_() *PreprocessString {
	if m != nil {
		return m.String_
	}
	return nil
}

func (m *PreprocessFieldOptions) GetTest() string {
	if m != nil && m.Test != nil {
		return *m.Test
	}
	return ""
}

type PreprocessString struct {
	TrimSpaces       *bool  `protobuf:"varint,1,opt,name=trim_spaces,json=trimSpaces" json:"trim_spaces,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *PreprocessString) Reset()                    { *m = PreprocessString{} }
func (m *PreprocessString) String() string            { return proto.CompactTextString(m) }
func (*PreprocessString) ProtoMessage()               {}
func (*PreprocessString) Descriptor() ([]byte, []int) { return fileDescriptorPreprocess, []int{1} }

func (m *PreprocessString) GetTrimSpaces() bool {
	if m != nil && m.TrimSpaces != nil {
		return *m.TrimSpaces
	}
	return false
}

var E_Field = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.FieldOptions)(nil),
	ExtensionType: (*PreprocessFieldOptions)(nil),
	Field:         11111,
	Name:          "preprocess.field",
	Tag:           "bytes,11111,opt,name=field",
	Filename:      "options/preprocess.proto",
}

func init() {
	proto.RegisterType((*PreprocessFieldOptions)(nil), "preprocess.PreprocessFieldOptions")
	proto.RegisterType((*PreprocessString)(nil), "preprocess.PreprocessString")
	proto.RegisterExtension(E_Field)
}

func init() { proto.RegisterFile("options/preprocess.proto", fileDescriptorPreprocess) }

var fileDescriptorPreprocess = []byte{
	// 212 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0xc8, 0x2f, 0x28, 0xc9,
	0xcc, 0xcf, 0x2b, 0xd6, 0x2f, 0x28, 0x4a, 0x2d, 0x28, 0xca, 0x4f, 0x4e, 0x2d, 0x2e, 0xd6, 0x2b,
	0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x42, 0x88, 0x48, 0x29, 0xa4, 0xe7, 0xe7, 0xa7, 0xe7, 0xa4,
	0xea, 0x83, 0x65, 0x92, 0x4a, 0xd3, 0xf4, 0x53, 0x52, 0x8b, 0x93, 0x8b, 0x32, 0x0b, 0x4a, 0xf2,
	0x8b, 0x20, 0xaa, 0x95, 0x92, 0xb8, 0xc4, 0x02, 0xe0, 0xea, 0xdd, 0x32, 0x53, 0x73, 0x52, 0xfc,
	0x21, 0x06, 0x0b, 0x99, 0x70, 0xb1, 0x15, 0x97, 0x14, 0x65, 0xe6, 0xa5, 0x4b, 0x30, 0x2a, 0x30,
	0x6a, 0x70, 0x1b, 0xc9, 0xe8, 0x21, 0x59, 0x85, 0xd0, 0x13, 0x0c, 0x56, 0x13, 0x04, 0x55, 0x2b,
	0x24, 0xc4, 0xc5, 0x52, 0x92, 0x5a, 0x5c, 0x22, 0xc1, 0xa4, 0xc0, 0xa8, 0xc1, 0x19, 0x04, 0x66,
	0x2b, 0x19, 0x73, 0x09, 0xa0, 0xab, 0x17, 0x92, 0xe7, 0xe2, 0x2e, 0x29, 0xca, 0xcc, 0x8d, 0x2f,
	0x2e, 0x48, 0x4c, 0x4e, 0x2d, 0x06, 0x5b, 0xc1, 0x11, 0xc4, 0x05, 0x12, 0x0a, 0x06, 0x8b, 0x58,
	0x45, 0x70, 0xb1, 0xa6, 0x81, 0x9c, 0x23, 0x24, 0xab, 0x07, 0xf1, 0x84, 0x1e, 0xcc, 0x13, 0x7a,
	0xc8, 0xce, 0x94, 0x78, 0x1e, 0x06, 0x76, 0x9d, 0x12, 0x76, 0xd7, 0x21, 0x2b, 0x0d, 0x82, 0x18,
	0xe8, 0xc4, 0x13, 0x85, 0x14, 0x44, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x01, 0x89, 0xfe, 0x36,
	0x49, 0x01, 0x00, 0x00,
}
