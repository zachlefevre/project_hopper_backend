// Code generated by protoc-gen-go. DO NOT EDIT.
// source: algorithm.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type CreateAlgorithmCommand struct {
	Algorithm            *Algorithm `protobuf:"bytes,1,opt,name=algorithm,proto3" json:"algorithm,omitempty"`
	CreatedOn            int64      `protobuf:"varint,2,opt,name=createdOn,proto3" json:"createdOn,omitempty"`
	Id                   string     `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *CreateAlgorithmCommand) Reset()         { *m = CreateAlgorithmCommand{} }
func (m *CreateAlgorithmCommand) String() string { return proto.CompactTextString(m) }
func (*CreateAlgorithmCommand) ProtoMessage()    {}
func (*CreateAlgorithmCommand) Descriptor() ([]byte, []int) {
	return fileDescriptor_8c381a4f1e580eed, []int{0}
}

func (m *CreateAlgorithmCommand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateAlgorithmCommand.Unmarshal(m, b)
}
func (m *CreateAlgorithmCommand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateAlgorithmCommand.Marshal(b, m, deterministic)
}
func (m *CreateAlgorithmCommand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateAlgorithmCommand.Merge(m, src)
}
func (m *CreateAlgorithmCommand) XXX_Size() int {
	return xxx_messageInfo_CreateAlgorithmCommand.Size(m)
}
func (m *CreateAlgorithmCommand) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateAlgorithmCommand.DiscardUnknown(m)
}

var xxx_messageInfo_CreateAlgorithmCommand proto.InternalMessageInfo

func (m *CreateAlgorithmCommand) GetAlgorithm() *Algorithm {
	if m != nil {
		return m.Algorithm
	}
	return nil
}

func (m *CreateAlgorithmCommand) GetCreatedOn() int64 {
	if m != nil {
		return m.CreatedOn
	}
	return 0
}

func (m *CreateAlgorithmCommand) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type GetAlgorithmQuery struct {
	Algorithm            *Algorithm `protobuf:"bytes,1,opt,name=algorithm,proto3" json:"algorithm,omitempty"`
	CreatedOn            int64      `protobuf:"varint,2,opt,name=createdOn,proto3" json:"createdOn,omitempty"`
	Id                   string     `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *GetAlgorithmQuery) Reset()         { *m = GetAlgorithmQuery{} }
func (m *GetAlgorithmQuery) String() string { return proto.CompactTextString(m) }
func (*GetAlgorithmQuery) ProtoMessage()    {}
func (*GetAlgorithmQuery) Descriptor() ([]byte, []int) {
	return fileDescriptor_8c381a4f1e580eed, []int{1}
}

func (m *GetAlgorithmQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAlgorithmQuery.Unmarshal(m, b)
}
func (m *GetAlgorithmQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAlgorithmQuery.Marshal(b, m, deterministic)
}
func (m *GetAlgorithmQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAlgorithmQuery.Merge(m, src)
}
func (m *GetAlgorithmQuery) XXX_Size() int {
	return xxx_messageInfo_GetAlgorithmQuery.Size(m)
}
func (m *GetAlgorithmQuery) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAlgorithmQuery.DiscardUnknown(m)
}

var xxx_messageInfo_GetAlgorithmQuery proto.InternalMessageInfo

func (m *GetAlgorithmQuery) GetAlgorithm() *Algorithm {
	if m != nil {
		return m.Algorithm
	}
	return nil
}

func (m *GetAlgorithmQuery) GetCreatedOn() int64 {
	if m != nil {
		return m.CreatedOn
	}
	return 0
}

func (m *GetAlgorithmQuery) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type Algorithm struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version              string   `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	Id                   string   `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	Status               string   `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`
	FileIDs              []string `protobuf:"bytes,5,rep,name=fileIDs,proto3" json:"fileIDs,omitempty"`
	DatasetIDs           []string `protobuf:"bytes,6,rep,name=datasetIDs,proto3" json:"datasetIDs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Algorithm) Reset()         { *m = Algorithm{} }
func (m *Algorithm) String() string { return proto.CompactTextString(m) }
func (*Algorithm) ProtoMessage()    {}
func (*Algorithm) Descriptor() ([]byte, []int) {
	return fileDescriptor_8c381a4f1e580eed, []int{2}
}

func (m *Algorithm) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Algorithm.Unmarshal(m, b)
}
func (m *Algorithm) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Algorithm.Marshal(b, m, deterministic)
}
func (m *Algorithm) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Algorithm.Merge(m, src)
}
func (m *Algorithm) XXX_Size() int {
	return xxx_messageInfo_Algorithm.Size(m)
}
func (m *Algorithm) XXX_DiscardUnknown() {
	xxx_messageInfo_Algorithm.DiscardUnknown(m)
}

var xxx_messageInfo_Algorithm proto.InternalMessageInfo

func (m *Algorithm) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Algorithm) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *Algorithm) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Algorithm) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *Algorithm) GetFileIDs() []string {
	if m != nil {
		return m.FileIDs
	}
	return nil
}

func (m *Algorithm) GetDatasetIDs() []string {
	if m != nil {
		return m.DatasetIDs
	}
	return nil
}

// Events
type AlgorithmCreatedEvent struct {
	Algorithm            *Algorithm `protobuf:"bytes,1,opt,name=algorithm,proto3" json:"algorithm,omitempty"`
	Id                   string     `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *AlgorithmCreatedEvent) Reset()         { *m = AlgorithmCreatedEvent{} }
func (m *AlgorithmCreatedEvent) String() string { return proto.CompactTextString(m) }
func (*AlgorithmCreatedEvent) ProtoMessage()    {}
func (*AlgorithmCreatedEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_8c381a4f1e580eed, []int{3}
}

func (m *AlgorithmCreatedEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AlgorithmCreatedEvent.Unmarshal(m, b)
}
func (m *AlgorithmCreatedEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AlgorithmCreatedEvent.Marshal(b, m, deterministic)
}
func (m *AlgorithmCreatedEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AlgorithmCreatedEvent.Merge(m, src)
}
func (m *AlgorithmCreatedEvent) XXX_Size() int {
	return xxx_messageInfo_AlgorithmCreatedEvent.Size(m)
}
func (m *AlgorithmCreatedEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_AlgorithmCreatedEvent.DiscardUnknown(m)
}

var xxx_messageInfo_AlgorithmCreatedEvent proto.InternalMessageInfo

func (m *AlgorithmCreatedEvent) GetAlgorithm() *Algorithm {
	if m != nil {
		return m.Algorithm
	}
	return nil
}

func (m *AlgorithmCreatedEvent) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type MultipleAlgorithms struct {
	Algorithms           []*Algorithm `protobuf:"bytes,1,rep,name=algorithms,proto3" json:"algorithms,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *MultipleAlgorithms) Reset()         { *m = MultipleAlgorithms{} }
func (m *MultipleAlgorithms) String() string { return proto.CompactTextString(m) }
func (*MultipleAlgorithms) ProtoMessage()    {}
func (*MultipleAlgorithms) Descriptor() ([]byte, []int) {
	return fileDescriptor_8c381a4f1e580eed, []int{4}
}

func (m *MultipleAlgorithms) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MultipleAlgorithms.Unmarshal(m, b)
}
func (m *MultipleAlgorithms) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MultipleAlgorithms.Marshal(b, m, deterministic)
}
func (m *MultipleAlgorithms) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MultipleAlgorithms.Merge(m, src)
}
func (m *MultipleAlgorithms) XXX_Size() int {
	return xxx_messageInfo_MultipleAlgorithms.Size(m)
}
func (m *MultipleAlgorithms) XXX_DiscardUnknown() {
	xxx_messageInfo_MultipleAlgorithms.DiscardUnknown(m)
}

var xxx_messageInfo_MultipleAlgorithms proto.InternalMessageInfo

func (m *MultipleAlgorithms) GetAlgorithms() []*Algorithm {
	if m != nil {
		return m.Algorithms
	}
	return nil
}

func init() {
	proto.RegisterType((*CreateAlgorithmCommand)(nil), "pb.CreateAlgorithmCommand")
	proto.RegisterType((*GetAlgorithmQuery)(nil), "pb.GetAlgorithmQuery")
	proto.RegisterType((*Algorithm)(nil), "pb.Algorithm")
	proto.RegisterType((*AlgorithmCreatedEvent)(nil), "pb.AlgorithmCreatedEvent")
	proto.RegisterType((*MultipleAlgorithms)(nil), "pb.MultipleAlgorithms")
}

func init() { proto.RegisterFile("algorithm.proto", fileDescriptor_8c381a4f1e580eed) }

var fileDescriptor_8c381a4f1e580eed = []byte{
	// 362 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x93, 0xdb, 0x4a, 0xf3, 0x40,
	0x10, 0xc7, 0xbb, 0x49, 0xbf, 0x7e, 0x64, 0xb4, 0x16, 0x47, 0x5a, 0x96, 0x22, 0x12, 0x72, 0x55,
	0x10, 0x03, 0x56, 0x10, 0xaf, 0x84, 0x52, 0x45, 0xbc, 0x10, 0x31, 0xfa, 0x02, 0x5b, 0xb3, 0xc6,
	0x40, 0x4e, 0xec, 0x6e, 0x0b, 0xbe, 0x83, 0x6f, 0xe0, 0x5b, 0xf8, 0x84, 0x92, 0xb4, 0xdd, 0x9c,
	0x04, 0xf1, 0xc2, 0xbb, 0xce, 0xe1, 0xff, 0xfb, 0x77, 0x66, 0x27, 0x30, 0x60, 0x51, 0x90, 0x8a,
	0x50, 0xbd, 0xc6, 0x6e, 0x26, 0x52, 0x95, 0xa2, 0x91, 0x2d, 0x1c, 0x09, 0xa3, 0xb9, 0xe0, 0x4c,
	0xf1, 0xd9, 0xb6, 0x38, 0x4f, 0xe3, 0x98, 0x25, 0x3e, 0x1e, 0x83, 0xa5, 0x05, 0x94, 0xd8, 0x64,
	0xb2, 0x33, 0xed, 0xbb, 0xd9, 0xc2, 0xd5, 0x8d, 0x5e, 0x59, 0xc7, 0x43, 0xb0, 0x9e, 0x0b, 0x8c,
	0x7f, 0x9f, 0x50, 0xc3, 0x26, 0x13, 0xd3, 0x2b, 0x13, 0xb8, 0x07, 0x46, 0xe8, 0x53, 0xd3, 0x26,
	0x13, 0xcb, 0x33, 0x42, 0xdf, 0x49, 0x60, 0xff, 0x86, 0x2b, 0x0d, 0x7a, 0x58, 0x72, 0xf1, 0xf6,
	0x97, 0x7e, 0x1f, 0x04, 0x2c, 0x8d, 0x41, 0x84, 0x6e, 0xc2, 0x62, 0x5e, 0x78, 0x58, 0x5e, 0xf1,
	0x1b, 0x29, 0xfc, 0x5f, 0x71, 0x21, 0xc3, 0x74, 0x4d, 0xb3, 0xbc, 0x6d, 0xd8, 0x64, 0xe1, 0x08,
	0x7a, 0x52, 0x31, 0xb5, 0x94, 0xb4, 0x5b, 0xe4, 0x36, 0x51, 0x4e, 0x78, 0x09, 0x23, 0x7e, 0x7b,
	0x25, 0xe9, 0x3f, 0xdb, 0xcc, 0x09, 0x9b, 0x10, 0x8f, 0x00, 0x7c, 0xa6, 0x98, 0xe4, 0x2a, 0x2f,
	0xf6, 0x8a, 0x62, 0x25, 0xe3, 0x3c, 0xc1, 0xb0, 0x5c, 0xfe, 0x7a, 0x86, 0xeb, 0x15, 0x4f, 0xd4,
	0xef, 0x36, 0xb2, 0xfe, 0x9f, 0x86, 0x9e, 0x79, 0x0e, 0x78, 0xb7, 0x8c, 0x54, 0x98, 0x45, 0xe5,
	0xd3, 0x4a, 0x3c, 0x01, 0xd0, 0x12, 0x49, 0x89, 0x6d, 0xb6, 0x99, 0x95, 0x86, 0xe9, 0x3b, 0x01,
	0xd4, 0x95, 0x59, 0x10, 0x08, 0x1e, 0x30, 0xc5, 0xf1, 0x12, 0x06, 0x8d, 0xa3, 0xc1, 0x71, 0x0e,
	0xf9, 0xfe, 0x92, 0xc6, 0x75, 0x03, 0xa7, 0x83, 0xe7, 0xb0, 0x5b, 0x7d, 0x7f, 0x1c, 0xe6, 0x0d,
	0xad, 0x8b, 0x68, 0xe9, 0xa6, 0x9f, 0x04, 0x0e, 0xea, 0x3d, 0x8f, 0x2a, 0x15, 0x1c, 0xdd, 0x06,
	0xaf, 0x2e, 0x6c, 0xfb, 0x5f, 0x40, 0xbf, 0xda, 0x2f, 0x9b, 0x82, 0x51, 0x1e, 0xb6, 0xb7, 0xe7,
	0x74, 0xf0, 0xb4, 0x3d, 0xf9, 0x0f, 0x66, 0x8b, 0x5e, 0xf1, 0xb1, 0x9d, 0x7d, 0x05, 0x00, 0x00,
	0xff, 0xff, 0xef, 0xf2, 0x58, 0xb6, 0x7f, 0x03, 0x00, 0x00,
}
