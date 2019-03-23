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

// CONNECTTYPE - connect type
type CONNECTTYPE int32

const (
	// unknown connection
	CONNECTTYPE_UNKNOWN_CONN CONNECTTYPE = 0
	// direct connection
	CONNECTTYPE_DIRECT_CONN CONNECTTYPE = 1
	// forward once
	CONNECTTYPE_FORWARD_ONCE CONNECTTYPE = 2
	// forward multiple times
	CONNECTTYPE_FORWARD_MULTIPLE CONNECTTYPE = 3
)

var CONNECTTYPE_name = map[int32]string{
	0: "UNKNOWN_CONN",
	1: "DIRECT_CONN",
	2: "FORWARD_ONCE",
	3: "FORWARD_MULTIPLE",
}
var CONNECTTYPE_value = map[string]int32{
	"UNKNOWN_CONN":     0,
	"DIRECT_CONN":      1,
	"FORWARD_ONCE":     2,
	"FORWARD_MULTIPLE": 3,
}

func (x CONNECTTYPE) String() string {
	return proto.EnumName(CONNECTTYPE_name, int32(x))
}
func (CONNECTTYPE) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_coredb_71aed0bbcb80d5a0, []int{0}
}

// private data
type PrivateData struct {
	PriKey               []byte   `protobuf:"bytes,1,opt,name=priKey,proto3" json:"priKey,omitempty"`
	PubKey               []byte   `protobuf:"bytes,2,opt,name=pubKey,proto3" json:"pubKey,omitempty"`
	CreateTime           int64    `protobuf:"varint,3,opt,name=createTime,proto3" json:"createTime,omitempty"`
	OnlineTime           int64    `protobuf:"varint,4,opt,name=onlineTime,proto3" json:"onlineTime,omitempty"`
	Addr                 string   `protobuf:"bytes,5,opt,name=addr,proto3" json:"addr,omitempty"`
	StrPriKey            string   `protobuf:"bytes,6,opt,name=strPriKey,proto3" json:"strPriKey,omitempty"`
	StrPubKey            string   `protobuf:"bytes,7,opt,name=strPubKey,proto3" json:"strPubKey,omitempty"`
	LstTrustNode         []string `protobuf:"bytes,8,rep,name=lstTrustNode,proto3" json:"lstTrustNode,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PrivateData) Reset()         { *m = PrivateData{} }
func (m *PrivateData) String() string { return proto.CompactTextString(m) }
func (*PrivateData) ProtoMessage()    {}
func (*PrivateData) Descriptor() ([]byte, []int) {
	return fileDescriptor_coredb_71aed0bbcb80d5a0, []int{0}
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

func (m *PrivateData) GetStrPriKey() string {
	if m != nil {
		return m.StrPriKey
	}
	return ""
}

func (m *PrivateData) GetStrPubKey() string {
	if m != nil {
		return m.StrPubKey
	}
	return ""
}

func (m *PrivateData) GetLstTrustNode() []string {
	if m != nil {
		return m.LstTrustNode
	}
	return nil
}

// node info
type NodeInfo struct {
	ServAddr      string   `protobuf:"bytes,1,opt,name=servAddr,proto3" json:"servAddr,omitempty"`
	Addr          string   `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	Name          string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	ConnectNums   int32    `protobuf:"varint,4,opt,name=connectNums,proto3" json:"connectNums,omitempty"`
	ConnectedNums int32    `protobuf:"varint,5,opt,name=connectedNums,proto3" json:"connectedNums,omitempty"`
	CtrlID        int64    `protobuf:"varint,6,opt,name=ctrlID,proto3" json:"ctrlID,omitempty"` // Deprecated: Do not use.
	LstClientAddr []string `protobuf:"bytes,7,rep,name=lstClientAddr,proto3" json:"lstClientAddr,omitempty"`
	AddTime       int64    `protobuf:"varint,8,opt,name=addTime,proto3" json:"addTime,omitempty"`
	ConnectMe     bool     `protobuf:"varint,9,opt,name=connectMe,proto3" json:"connectMe,omitempty"`
	ConnectNode   bool     `protobuf:"varint,10,opt,name=connectNode,proto3" json:"connectNode,omitempty"` // Deprecated: Do not use.
	// nodetype version
	NodeTypeVersion string `protobuf:"bytes,11,opt,name=nodeTypeVersion,proto3" json:"nodeTypeVersion,omitempty"`
	// node type
	NodeType string `protobuf:"bytes,12,opt,name=nodeType,proto3" json:"nodeType,omitempty"`
	// jarviscore version
	CoreVersion string `protobuf:"bytes,13,opt,name=coreVersion,proto3" json:"coreVersion,omitempty"`
	// current message id
	// I send message for this node, the msgid is lastSendMsgID + 1
	LastSendMsgID int64 `protobuf:"varint,14,opt,name=lastSendMsgID,proto3" json:"lastSendMsgID,omitempty"`
	// last connect time
	LastConnectTime int64 `protobuf:"varint,15,opt,name=lastConnectTime,proto3" json:"lastConnectTime,omitempty"`
	// last connected time
	LastConnectedTime int64 `protobuf:"varint,16,opt,name=lastConnectedTime,proto3" json:"lastConnectedTime,omitempty"`
	// last connect me time
	LastConnectMeTime int64 `protobuf:"varint,17,opt,name=lastConnectMeTime,proto3" json:"lastConnectMeTime,omitempty"`
	// groups
	LstGroups []string `protobuf:"bytes,18,rep,name=lstGroups,proto3" json:"lstGroups,omitempty"`
	// deprecated
	Deprecated bool `protobuf:"varint,19,opt,name=deprecated,proto3" json:"deprecated,omitempty"`
	// last receive message id
	LastRecvMsgID int64 `protobuf:"varint,20,opt,name=lastRecvMsgID,proto3" json:"lastRecvMsgID,omitempty"`
	// connection type
	ConnType CONNECTTYPE `protobuf:"varint,21,opt,name=connType,proto3,enum=coredbpb.CONNECTTYPE" json:"connType,omitempty"`
	// When I can't connect directly to this node, this is the list of nodes I know
	// that can be directly connected to this node.
	ValidConnNodes []string `protobuf:"bytes,22,rep,name=validConnNodes,proto3" json:"validConnNodes,omitempty"`
	// This node was ignored before this time,
	// this property is only used for deprecated to be false
	TimestampDeprecated int64 `protobuf:"varint,23,opt,name=timestampDeprecated,proto3" json:"timestampDeprecated,omitempty"`
	// This property is the number of failures to connect to the node.
	// If there is a successful connection, this property will be reset to 0.
	NumsConnectFail      int32    `protobuf:"varint,24,opt,name=numsConnectFail,proto3" json:"numsConnectFail,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NodeInfo) Reset()         { *m = NodeInfo{} }
func (m *NodeInfo) String() string { return proto.CompactTextString(m) }
func (*NodeInfo) ProtoMessage()    {}
func (*NodeInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_coredb_71aed0bbcb80d5a0, []int{1}
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

// Deprecated: Do not use.
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

func (m *NodeInfo) GetConnectMe() bool {
	if m != nil {
		return m.ConnectMe
	}
	return false
}

// Deprecated: Do not use.
func (m *NodeInfo) GetConnectNode() bool {
	if m != nil {
		return m.ConnectNode
	}
	return false
}

func (m *NodeInfo) GetNodeTypeVersion() string {
	if m != nil {
		return m.NodeTypeVersion
	}
	return ""
}

func (m *NodeInfo) GetNodeType() string {
	if m != nil {
		return m.NodeType
	}
	return ""
}

func (m *NodeInfo) GetCoreVersion() string {
	if m != nil {
		return m.CoreVersion
	}
	return ""
}

func (m *NodeInfo) GetLastSendMsgID() int64 {
	if m != nil {
		return m.LastSendMsgID
	}
	return 0
}

func (m *NodeInfo) GetLastConnectTime() int64 {
	if m != nil {
		return m.LastConnectTime
	}
	return 0
}

func (m *NodeInfo) GetLastConnectedTime() int64 {
	if m != nil {
		return m.LastConnectedTime
	}
	return 0
}

func (m *NodeInfo) GetLastConnectMeTime() int64 {
	if m != nil {
		return m.LastConnectMeTime
	}
	return 0
}

func (m *NodeInfo) GetLstGroups() []string {
	if m != nil {
		return m.LstGroups
	}
	return nil
}

func (m *NodeInfo) GetDeprecated() bool {
	if m != nil {
		return m.Deprecated
	}
	return false
}

func (m *NodeInfo) GetLastRecvMsgID() int64 {
	if m != nil {
		return m.LastRecvMsgID
	}
	return 0
}

func (m *NodeInfo) GetConnType() CONNECTTYPE {
	if m != nil {
		return m.ConnType
	}
	return CONNECTTYPE_UNKNOWN_CONN
}

func (m *NodeInfo) GetValidConnNodes() []string {
	if m != nil {
		return m.ValidConnNodes
	}
	return nil
}

func (m *NodeInfo) GetTimestampDeprecated() int64 {
	if m != nil {
		return m.TimestampDeprecated
	}
	return 0
}

func (m *NodeInfo) GetNumsConnectFail() int32 {
	if m != nil {
		return m.NumsConnectFail
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
	return fileDescriptor_coredb_71aed0bbcb80d5a0, []int{2}
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

// FileSystem
type FileSystem struct {
	CurFileID            int64    `protobuf:"varint,1,opt,name=curFileID,proto3" json:"curFileID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FileSystem) Reset()         { *m = FileSystem{} }
func (m *FileSystem) String() string { return proto.CompactTextString(m) }
func (*FileSystem) ProtoMessage()    {}
func (*FileSystem) Descriptor() ([]byte, []int) {
	return fileDescriptor_coredb_71aed0bbcb80d5a0, []int{3}
}
func (m *FileSystem) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FileSystem.Unmarshal(m, b)
}
func (m *FileSystem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FileSystem.Marshal(b, m, deterministic)
}
func (dst *FileSystem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FileSystem.Merge(dst, src)
}
func (m *FileSystem) XXX_Size() int {
	return xxx_messageInfo_FileSystem.Size(m)
}
func (m *FileSystem) XXX_DiscardUnknown() {
	xxx_messageInfo_FileSystem.DiscardUnknown(m)
}

var xxx_messageInfo_FileSystem proto.InternalMessageInfo

func (m *FileSystem) GetCurFileID() int64 {
	if m != nil {
		return m.CurFileID
	}
	return 0
}

func init() {
	proto.RegisterType((*PrivateData)(nil), "coredbpb.PrivateData")
	proto.RegisterType((*NodeInfo)(nil), "coredbpb.NodeInfo")
	proto.RegisterType((*NodeInfoList)(nil), "coredbpb.NodeInfoList")
	proto.RegisterType((*FileSystem)(nil), "coredbpb.FileSystem")
	proto.RegisterEnum("coredbpb.CONNECTTYPE", CONNECTTYPE_name, CONNECTTYPE_value)
}

func init() { proto.RegisterFile("coredb.proto", fileDescriptor_coredb_71aed0bbcb80d5a0) }

var fileDescriptor_coredb_71aed0bbcb80d5a0 = []byte{
	// 688 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x94, 0xef, 0x6e, 0xda, 0x30,
	0x14, 0xc5, 0x17, 0x02, 0x14, 0x2e, 0xb4, 0xa5, 0xee, 0x9f, 0x59, 0xd5, 0x34, 0x45, 0xa8, 0x9a,
	0xa2, 0x6a, 0xaa, 0xb6, 0xee, 0x09, 0x3a, 0xa0, 0x13, 0x6a, 0x09, 0xc8, 0xa5, 0xeb, 0xf6, 0xa9,
	0x0a, 0xc4, 0xdb, 0x22, 0x85, 0x24, 0xb2, 0x0d, 0x6a, 0xdf, 0x62, 0x2f, 0xb9, 0x07, 0xd8, 0x1b,
	0x4c, 0xbe, 0x26, 0x7f, 0xa0, 0xfb, 0x96, 0xfb, 0x3b, 0x57, 0xf6, 0xb9, 0xc7, 0x76, 0xa0, 0x3d,
	0x4f, 0x04, 0x0f, 0x66, 0x17, 0xa9, 0x48, 0x54, 0x42, 0x1a, 0xa6, 0x4a, 0x67, 0xdd, 0xbf, 0x16,
	0xb4, 0x26, 0x22, 0x5c, 0xf9, 0x8a, 0xf7, 0x7d, 0xe5, 0x93, 0x13, 0xa8, 0xa7, 0x22, 0xbc, 0xe1,
	0xcf, 0xd4, 0x72, 0x2c, 0xb7, 0xcd, 0xd6, 0x15, 0xf2, 0xe5, 0x4c, 0xf3, 0xca, 0x9a, 0x63, 0x45,
	0xde, 0x02, 0xcc, 0x05, 0xf7, 0x15, 0x9f, 0x86, 0x0b, 0x4e, 0x6d, 0xc7, 0x72, 0x6d, 0x56, 0x22,
	0x5a, 0x4f, 0xe2, 0x28, 0x8c, 0x8d, 0x5e, 0x35, 0x7a, 0x41, 0x08, 0x81, 0xaa, 0x1f, 0x04, 0x82,
	0xd6, 0x1c, 0xcb, 0x6d, 0x32, 0xfc, 0x26, 0x6f, 0xa0, 0x29, 0x95, 0x98, 0x18, 0x1b, 0x75, 0x14,
	0x0a, 0x90, 0xa9, 0xc6, 0xcc, 0x4e, 0xa1, 0x1a, 0x3f, 0x5d, 0x68, 0x47, 0x52, 0x4d, 0xc5, 0x52,
	0x2a, 0x2f, 0x09, 0x38, 0x6d, 0x38, 0xb6, 0xdb, 0x64, 0x1b, 0xac, 0xfb, 0xa7, 0x0e, 0x0d, 0xfd,
	0x31, 0x8c, 0x7f, 0x24, 0xe4, 0x14, 0x1a, 0x92, 0x8b, 0xd5, 0x95, 0x36, 0x61, 0xe1, 0x6a, 0x79,
	0x9d, 0x9b, 0xab, 0x94, 0xcc, 0x11, 0xa8, 0xc6, 0xfe, 0x7a, 0xd4, 0x26, 0xc3, 0x6f, 0xe2, 0x40,
	0x6b, 0x9e, 0xc4, 0x31, 0x9f, 0x2b, 0x6f, 0xb9, 0x90, 0x38, 0x65, 0x8d, 0x95, 0x11, 0x39, 0x83,
	0xdd, 0x75, 0xc9, 0x03, 0xec, 0xa9, 0x61, 0xcf, 0x26, 0x24, 0xa7, 0x50, 0x9f, 0x2b, 0x11, 0x0d,
	0xfb, 0x38, 0xb5, 0xfd, 0xb9, 0x42, 0x2d, 0xb6, 0x26, 0x7a, 0x85, 0x48, 0xaa, 0x5e, 0x14, 0xf2,
	0x58, 0xa1, 0xd9, 0x1d, 0x9c, 0x6c, 0x13, 0x12, 0x0a, 0x3b, 0x7e, 0x10, 0x60, 0xd6, 0x0d, 0xcc,
	0x3a, 0x2b, 0x75, 0x6c, 0xeb, 0xcd, 0x46, 0x9c, 0x36, 0x1d, 0xcb, 0x6d, 0xb0, 0x02, 0x90, 0xb3,
	0x62, 0x02, 0x9d, 0x1a, 0x68, 0x1d, 0xb7, 0x2f, 0x63, 0xe2, 0xc2, 0x7e, 0x9c, 0x04, 0x7c, 0xfa,
	0x9c, 0xf2, 0xaf, 0x5c, 0xc8, 0x30, 0x89, 0x69, 0x0b, 0x63, 0xd8, 0xc6, 0x3a, 0xd5, 0x0c, 0xd1,
	0xb6, 0x49, 0x35, 0xab, 0x4d, 0x5a, 0x22, 0x5f, 0x61, 0x17, 0xe5, 0x32, 0xc2, 0x59, 0x7d, 0xa9,
	0xee, 0x78, 0x1c, 0x8c, 0xe4, 0xcf, 0x61, 0x9f, 0xee, 0xe1, 0x2c, 0x9b, 0x50, 0xbb, 0xd1, 0xa0,
	0x67, 0x0c, 0xe2, 0xcc, 0xfb, 0xd8, 0xb7, 0x8d, 0xc9, 0x7b, 0x38, 0x28, 0x21, 0x6e, 0xf2, 0xe9,
	0x60, 0xef, 0x4b, 0x61, 0xab, 0x7b, 0x64, 0x6e, 0xee, 0xc1, 0x8b, 0x6e, 0x23, 0xe8, 0x5c, 0x23,
	0xa9, 0xbe, 0x88, 0x64, 0x99, 0x4a, 0x4a, 0xf0, 0x4c, 0x0a, 0xa0, 0xaf, 0x7f, 0xc0, 0x53, 0xc1,
	0xe7, 0xbe, 0xe2, 0x01, 0x3d, 0xc4, 0xd8, 0x4b, 0x24, 0x9b, 0x94, 0xf1, 0xf9, 0xca, 0x4c, 0x7a,
	0x54, 0x4c, 0x9a, 0x43, 0xf2, 0x11, 0x1a, 0xfa, 0x18, 0x30, 0xcd, 0x63, 0xc7, 0x72, 0xf7, 0x2e,
	0x8f, 0x2f, 0xb2, 0x17, 0x7c, 0xd1, 0x1b, 0x7b, 0xde, 0xa0, 0x37, 0x9d, 0x7e, 0x9f, 0x0c, 0x58,
	0xde, 0x46, 0xde, 0xc1, 0xde, 0xca, 0x8f, 0xc2, 0x40, 0x9b, 0xd5, 0x67, 0x27, 0xe9, 0x09, 0x7a,
	0xdb, 0xa2, 0xe4, 0x03, 0x1c, 0xaa, 0x70, 0xc1, 0xa5, 0xf2, 0x17, 0x69, 0xbf, 0x70, 0xfa, 0x1a,
	0x6d, 0xfc, 0x4f, 0xc2, 0x4b, 0xb0, 0x5c, 0xc8, 0x75, 0x0a, 0xd7, 0x7e, 0x18, 0x51, 0x8a, 0x97,
	0x79, 0x1b, 0x77, 0x7f, 0x5b, 0xd0, 0xce, 0xde, 0xd9, 0x6d, 0x28, 0x95, 0x4e, 0x43, 0xc6, 0x7e,
	0x2a, 0x7f, 0x25, 0x6a, 0xd8, 0xc7, 0xd7, 0x66, 0xb3, 0x12, 0xd1, 0xb7, 0x86, 0xc7, 0xc1, 0x30,
	0x0e, 0xf8, 0x13, 0xbe, 0xb9, 0x1a, 0xcb, 0x6b, 0xad, 0x2d, 0xfc, 0x27, 0xa3, 0xd9, 0x46, 0xcb,
	0x6a, 0xe2, 0x42, 0x2d, 0xc6, 0x19, 0xab, 0x8e, 0xed, 0xb6, 0x2e, 0x49, 0x11, 0x4e, 0xb6, 0x3d,
	0x33, 0x0d, 0xdd, 0x73, 0x80, 0xeb, 0x30, 0xe2, 0x77, 0xcf, 0x52, 0xf1, 0x05, 0xbe, 0x89, 0xa5,
	0xd0, 0x20, 0xb7, 0x53, 0x80, 0xf3, 0x6f, 0xd0, 0x2a, 0x65, 0x4b, 0x3a, 0xd0, 0xbe, 0xf7, 0x6e,
	0xbc, 0xf1, 0x83, 0xf7, 0xa8, 0x71, 0xe7, 0x15, 0xd9, 0x87, 0x56, 0x7f, 0xc8, 0x06, 0xbd, 0xa9,
	0x01, 0x96, 0x6e, 0xb9, 0x1e, 0xb3, 0x87, 0x2b, 0xd6, 0x7f, 0x1c, 0x7b, 0xbd, 0x41, 0xa7, 0x42,
	0x8e, 0xa0, 0x93, 0x91, 0xd1, 0xfd, 0xed, 0x74, 0x38, 0xb9, 0x1d, 0x74, 0xec, 0x59, 0x1d, 0xff,
	0xc2, 0x9f, 0xfe, 0x05, 0x00, 0x00, 0xff, 0xff, 0x38, 0x22, 0x9c, 0xe1, 0x95, 0x05, 0x00, 0x00,
}
