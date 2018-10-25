// Code generated by protoc-gen-go. DO NOT EDIT.
// source: coredb.proto

package coredbpb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// private data
type PrivateData struct {
	PriKey               []byte   `protobuf:"bytes,1,opt,name=priKey,proto3" json:"priKey,omitempty"`
	PubKey               []byte   `protobuf:"bytes,2,opt,name=pubKey,proto3" json:"pubKey,omitempty"`
	CreateTime           int64    `protobuf:"varint,3,opt,name=createTime,proto3" json:"createTime,omitempty"`
	OnlineTime           int64    `protobuf:"varint,4,opt,name=onlineTime,proto3" json:"onlineTime,omitempty"`
	Addr                 string   `protobuf:"bytes,5,opt,name=addr,proto3" json:"addr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PrivateData) Reset()         { *m = PrivateData{} }
func (m *PrivateData) String() string { return proto.CompactTextString(m) }
func (*PrivateData) ProtoMessage()    {}
func (*PrivateData) Descriptor() ([]byte, []int) {
	return fileDescriptor_coredb_94e91060a685ce41, []int{0}
}
func (m *PrivateData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PrivateData.Unmarshal(m, b)
}
func (m *PrivateData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PrivateData.Marshal(b, m, deterministic)
}
func (dst *PrivateData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PrivateData.Merge(dst, src)
}
func (m *PrivateData) XXX_Size() int {
	return xxx_messageInfo_PrivateData.Size(m)
}
func (m *PrivateData) XXX_DiscardUnknown() {
	xxx_messageInfo_PrivateData.DiscardUnknown(m)
}

var xxx_messageInfo_PrivateData proto.InternalMessageInfo

func (m *PrivateData) GetPriKey() []byte {
	if m != nil {
		return m.PriKey
	}
	return nil
}

func (m *PrivateData) GetPubKey() []byte {
	if m != nil {
		return m.PubKey
	}
	return nil
}

func (m *PrivateData) GetCreateTime() int64 {
	if m != nil {
		return m.CreateTime
	}
	return 0
}

func (m *PrivateData) GetOnlineTime() int64 {
	if m != nil {
		return m.OnlineTime
	}
	return 0
}

func (m *PrivateData) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

// node info
type NodeInfo struct {
	ServAddr             string   `protobuf:"bytes,1,opt,name=servAddr,proto3" json:"servAddr,omitempty"`
	Addr                 string   `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	ConnectNums          int32    `protobuf:"varint,4,opt,name=connectNums,proto3" json:"connectNums,omitempty"`
	ConnectedNums        int32    `protobuf:"varint,5,opt,name=connectedNums,proto3" json:"connectedNums,omitempty"`
	CtrlID               int64    `protobuf:"varint,6,opt,name=ctrlID,proto3" json:"ctrlID,omitempty"`
	LstClientAddr        []string `protobuf:"bytes,7,rep,name=lstClientAddr,proto3" json:"lstClientAddr,omitempty"`
	AddTime              int64    `protobuf:"varint,8,opt,name=addTime,proto3" json:"addTime,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NodeInfo) Reset()         { *m = NodeInfo{} }
func (m *NodeInfo) String() string { return proto.CompactTextString(m) }
func (*NodeInfo) ProtoMessage()    {}
func (*NodeInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_coredb_94e91060a685ce41, []int{1}
}
func (m *NodeInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeInfo.Unmarshal(m, b)
}
func (m *NodeInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeInfo.Marshal(b, m, deterministic)
}
func (dst *NodeInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeInfo.Merge(dst, src)
}
func (m *NodeInfo) XXX_Size() int {
	return xxx_messageInfo_NodeInfo.Size(m)
}
func (m *NodeInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeInfo.DiscardUnknown(m)
}

var xxx_messageInfo_NodeInfo proto.InternalMessageInfo

func (m *NodeInfo) GetServAddr() string {
	if m != nil {
		return m.ServAddr
	}
	return ""
}

func (m *NodeInfo) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *NodeInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *NodeInfo) GetConnectNums() int32 {
	if m != nil {
		return m.ConnectNums
	}
	return 0
}

func (m *NodeInfo) GetConnectedNums() int32 {
	if m != nil {
		return m.ConnectedNums
	}
	return 0
}

func (m *NodeInfo) GetCtrlID() int64 {
	if m != nil {
		return m.CtrlID
	}
	return 0
}

func (m *NodeInfo) GetLstClientAddr() []string {
	if m != nil {
		return m.LstClientAddr
	}
	return nil
}

func (m *NodeInfo) GetAddTime() int64 {
	if m != nil {
		return m.AddTime
	}
	return 0
}

// node info list
type NodeInfoList struct {
	SnapshotID           int64       `protobuf:"varint,1,opt,name=snapshotID,proto3" json:"snapshotID,omitempty"`
	EndIndex             int32       `protobuf:"varint,2,opt,name=endIndex,proto3" json:"endIndex,omitempty"`
	MaxIndex             int32       `protobuf:"varint,3,opt,name=maxIndex,proto3" json:"maxIndex,omitempty"`
	Nodes                []*NodeInfo `protobuf:"bytes,4,rep,name=nodes,proto3" json:"nodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *NodeInfoList) Reset()         { *m = NodeInfoList{} }
func (m *NodeInfoList) String() string { return proto.CompactTextString(m) }
func (*NodeInfoList) ProtoMessage()    {}
func (*NodeInfoList) Descriptor() ([]byte, []int) {
	return fileDescriptor_coredb_94e91060a685ce41, []int{2}
}
func (m *NodeInfoList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeInfoList.Unmarshal(m, b)
}
func (m *NodeInfoList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeInfoList.Marshal(b, m, deterministic)
}
func (dst *NodeInfoList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeInfoList.Merge(dst, src)
}
func (m *NodeInfoList) XXX_Size() int {
	return xxx_messageInfo_NodeInfoList.Size(m)
}
func (m *NodeInfoList) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeInfoList.DiscardUnknown(m)
}

var xxx_messageInfo_NodeInfoList proto.InternalMessageInfo

func (m *NodeInfoList) GetSnapshotID() int64 {
	if m != nil {
		return m.SnapshotID
	}
	return 0
}

func (m *NodeInfoList) GetEndIndex() int32 {
	if m != nil {
		return m.EndIndex
	}
	return 0
}

func (m *NodeInfoList) GetMaxIndex() int32 {
	if m != nil {
		return m.MaxIndex
	}
	return 0
}

func (m *NodeInfoList) GetNodes() []*NodeInfo {
	if m != nil {
		return m.Nodes
	}
	return nil
}

func init() {
	proto.RegisterType((*PrivateData)(nil), "coredbpb.PrivateData")
	proto.RegisterType((*NodeInfo)(nil), "coredbpb.NodeInfo")
	proto.RegisterType((*NodeInfoList)(nil), "coredbpb.NodeInfoList")
}

func init() { proto.RegisterFile("coredb.proto", fileDescriptor_coredb_94e91060a685ce41) }

var fileDescriptor_coredb_94e91060a685ce41 = []byte{
	// 330 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x92, 0x41, 0x4e, 0xc3, 0x30,
	0x10, 0x45, 0xe5, 0xa6, 0x69, 0xd3, 0x69, 0xd9, 0x78, 0x81, 0x22, 0x16, 0x28, 0xaa, 0x58, 0x64,
	0xd5, 0x05, 0x9c, 0x00, 0xd1, 0x4d, 0x04, 0xaa, 0x90, 0xc5, 0x05, 0x9c, 0x78, 0x10, 0x91, 0x52,
	0x3b, 0x72, 0xdc, 0xaa, 0xdc, 0x02, 0x2e, 0xca, 0x19, 0x90, 0xc7, 0x49, 0x9b, 0xee, 0xf2, 0xdf,
	0xff, 0xa3, 0xcc, 0x7c, 0x19, 0x56, 0x95, 0xb1, 0xa8, 0xca, 0x4d, 0x6b, 0x8d, 0x33, 0x3c, 0x09,
	0xaa, 0x2d, 0xd7, 0xbf, 0x0c, 0x96, 0xef, 0xb6, 0x3e, 0x4a, 0x87, 0x5b, 0xe9, 0x24, 0xbf, 0x85,
	0x59, 0x6b, 0xeb, 0x57, 0xfc, 0x4e, 0x59, 0xc6, 0xf2, 0x95, 0xe8, 0x15, 0xf1, 0x43, 0xe9, 0xf9,
	0xa4, 0xe7, 0xa4, 0xf8, 0x3d, 0x40, 0x65, 0x51, 0x3a, 0xfc, 0xa8, 0xf7, 0x98, 0x46, 0x19, 0xcb,
	0x23, 0x31, 0x22, 0xde, 0x37, 0xba, 0xa9, 0x75, 0xf0, 0xa7, 0xc1, 0xbf, 0x10, 0xce, 0x61, 0x2a,
	0x95, 0xb2, 0x69, 0x9c, 0xb1, 0x7c, 0x21, 0xe8, 0x7b, 0xfd, 0xc7, 0x20, 0xd9, 0x19, 0x85, 0x85,
	0xfe, 0x34, 0xfc, 0x0e, 0x92, 0x0e, 0xed, 0xf1, 0xd9, 0x87, 0x18, 0x85, 0xce, 0xfa, 0x3c, 0x3c,
	0xb9, 0x0c, 0x7b, 0xa6, 0x65, 0xbf, 0xca, 0x42, 0xd0, 0x37, 0xcf, 0x60, 0x59, 0x19, 0xad, 0xb1,
	0x72, 0xbb, 0xc3, 0xbe, 0xa3, 0x2d, 0x62, 0x31, 0x46, 0xfc, 0x01, 0x6e, 0x7a, 0x89, 0x8a, 0x32,
	0x31, 0x65, 0xae, 0xa1, 0x2f, 0xa1, 0x72, 0xb6, 0x29, 0xb6, 0xe9, 0x8c, 0x0e, 0xe9, 0x95, 0x9f,
	0x6e, 0x3a, 0xf7, 0xd2, 0xd4, 0xa8, 0x1d, 0x2d, 0x3a, 0xcf, 0xa2, 0x7c, 0x21, 0xae, 0x21, 0x4f,
	0x61, 0x2e, 0x95, 0xa2, 0x1e, 0x12, 0x1a, 0x1f, 0xe4, 0xfa, 0x87, 0xc1, 0x6a, 0x38, 0xf8, 0xad,
	0xee, 0x9c, 0x6f, 0xad, 0xd3, 0xb2, 0xed, 0xbe, 0x8c, 0x2b, 0xb6, 0x74, 0x76, 0x24, 0x46, 0xc4,
	0x97, 0x82, 0x5a, 0x15, 0x5a, 0xe1, 0x89, 0x8e, 0x8f, 0xc5, 0x59, 0x7b, 0x6f, 0x2f, 0x4f, 0xc1,
	0x8b, 0x82, 0x37, 0x68, 0x9e, 0x43, 0xac, 0x8d, 0x42, 0x5f, 0x41, 0x94, 0x2f, 0x1f, 0xf9, 0x66,
	0x78, 0x07, 0x9b, 0xe1, 0xf7, 0x22, 0x04, 0xca, 0x19, 0x3d, 0x94, 0xa7, 0xff, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x13, 0xaa, 0xf4, 0xe8, 0x38, 0x02, 0x00, 0x00,
}
