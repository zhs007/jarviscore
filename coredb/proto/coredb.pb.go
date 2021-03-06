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
	return fileDescriptor_coredb_d26fe61598bfb6c9, []int{0}
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
	return fileDescriptor_coredb_d26fe61598bfb6c9, []int{0}
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
	// This is the connection address of the node
	ServAddr string `protobuf:"bytes,1,opt,name=servAddr,proto3" json:"servAddr,omitempty"`
	// This is the address of the node
	Addr string `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	// This is the node name
	Name string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	// This is the number of times I connected this node.
	ConnectNums int32 `protobuf:"varint,4,opt,name=connectNums,proto3" json:"connectNums,omitempty"`
	// This is the number of times I successfully connected this node.
	ConnectedNums int32 `protobuf:"varint,5,opt,name=connectedNums,proto3" json:"connectedNums,omitempty"`
	// ctrlID
	CtrlID int64 `protobuf:"varint,6,opt,name=ctrlID,proto3" json:"ctrlID,omitempty"` // Deprecated: Do not use.
	// This is the list of addresses connected to this node
	LstClientAddr []string `protobuf:"bytes,7,rep,name=lstClientAddr,proto3" json:"lstClientAddr,omitempty"`
	// This is the timestamp added for the first time.
	AddTime int64 `protobuf:"varint,8,opt,name=addTime,proto3" json:"addTime,omitempty"`
	// Is this node connected to me?
	ConnectMe bool `protobuf:"varint,9,opt,name=connectMe,proto3" json:"connectMe,omitempty"`
	// Am I connected to this node?
	ConnectNode bool `protobuf:"varint,10,opt,name=connectNode,proto3" json:"connectNode,omitempty"` // Deprecated: Do not use.
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
	NumsConnectFail int32 `protobuf:"varint,24,opt,name=numsConnectFail,proto3" json:"numsConnectFail,omitempty"`
	// Only if this value is 0, I will execute requestnodes
	LastMsgID4RequestNodes int64 `protobuf:"varint,25,opt,name=lastMsgID4RequestNodes,proto3" json:"lastMsgID4RequestNodes,omitempty"`
	// This is the version number of the node node data, which is usually a hash value.
	// This value will be set to lastNodesVersion after requestnodes.
	NodesVersion string `protobuf:"bytes,26,opt,name=nodesVersion,proto3" json:"nodesVersion,omitempty"`
	// This is the version number of the node node data, which is usually a hash value.
	LastNodesVersion     string   `protobuf:"bytes,27,opt,name=lastNodesVersion,proto3" json:"lastNodesVersion,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NodeInfo) Reset()         { *m = NodeInfo{} }
func (m *NodeInfo) String() string { return proto.CompactTextString(m) }
func (*NodeInfo) ProtoMessage()    {}
func (*NodeInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_coredb_d26fe61598bfb6c9, []int{1}
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

func (m *NodeInfo) GetLastMsgID4RequestNodes() int64 {
	if m != nil {
		return m.LastMsgID4RequestNodes
	}
	return 0
}

func (m *NodeInfo) GetNodesVersion() string {
	if m != nil {
		return m.NodesVersion
	}
	return ""
}

func (m *NodeInfo) GetLastNodesVersion() string {
	if m != nil {
		return m.LastNodesVersion
	}
	return ""
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
	return fileDescriptor_coredb_d26fe61598bfb6c9, []int{2}
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

// node info list v2
type NodeInfoList2 struct {
	Nodes                []*NodeInfo `protobuf:"bytes,1,rep,name=nodes,proto3" json:"nodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *NodeInfoList2) Reset()         { *m = NodeInfoList2{} }
func (m *NodeInfoList2) String() string { return proto.CompactTextString(m) }
func (*NodeInfoList2) ProtoMessage()    {}
func (*NodeInfoList2) Descriptor() ([]byte, []int) {
	return fileDescriptor_coredb_d26fe61598bfb6c9, []int{3}
}
func (m *NodeInfoList2) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeInfoList2.Unmarshal(m, b)
}
func (m *NodeInfoList2) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeInfoList2.Marshal(b, m, deterministic)
}
func (dst *NodeInfoList2) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeInfoList2.Merge(dst, src)
}
func (m *NodeInfoList2) XXX_Size() int {
	return xxx_messageInfo_NodeInfoList2.Size(m)
}
func (m *NodeInfoList2) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeInfoList2.DiscardUnknown(m)
}

var xxx_messageInfo_NodeInfoList2 proto.InternalMessageInfo

func (m *NodeInfoList2) GetNodes() []*NodeInfo {
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
	return fileDescriptor_coredb_d26fe61598bfb6c9, []int{4}
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
	proto.RegisterType((*NodeInfoList2)(nil), "coredbpb.NodeInfoList2")
	proto.RegisterType((*FileSystem)(nil), "coredbpb.FileSystem")
	proto.RegisterEnum("coredbpb.CONNECTTYPE", CONNECTTYPE_name, CONNECTTYPE_value)
}

func init() { proto.RegisterFile("coredb.proto", fileDescriptor_coredb_d26fe61598bfb6c9) }

var fileDescriptor_coredb_d26fe61598bfb6c9 = []byte{
	// 752 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x95, 0xdb, 0x6e, 0xdb, 0x46,
	0x10, 0x86, 0x4b, 0x51, 0xc7, 0x91, 0x6c, 0xcb, 0xeb, 0x43, 0xb7, 0x6e, 0x51, 0x10, 0x82, 0x51,
	0x10, 0x46, 0x61, 0xb4, 0x6e, 0x11, 0x20, 0x97, 0x8e, 0x24, 0x07, 0x82, 0x2d, 0x4a, 0x58, 0xcb,
	0x71, 0x72, 0x65, 0x50, 0xe2, 0x26, 0x21, 0x40, 0x91, 0xcc, 0xee, 0x4a, 0xb0, 0xdf, 0x22, 0x2f,
	0x9a, 0x8b, 0xbc, 0x41, 0xb0, 0xb3, 0xa2, 0x48, 0xc9, 0x09, 0x72, 0xc7, 0xf9, 0xfe, 0xe1, 0xf2,
	0x9f, 0xd9, 0x19, 0x10, 0x5a, 0xb3, 0x44, 0xf0, 0x60, 0x7a, 0x9e, 0x8a, 0x44, 0x25, 0xa4, 0x6e,
	0xa2, 0x74, 0xda, 0xf9, 0x6a, 0x41, 0x73, 0x2c, 0xc2, 0xa5, 0xaf, 0x78, 0xcf, 0x57, 0x3e, 0x39,
	0x86, 0x6a, 0x2a, 0xc2, 0x6b, 0xfe, 0x44, 0x2d, 0xc7, 0x72, 0x5b, 0x6c, 0x15, 0x21, 0x5f, 0x4c,
	0x35, 0x2f, 0xad, 0x38, 0x46, 0xe4, 0x4f, 0x80, 0x99, 0xe0, 0xbe, 0xe2, 0x93, 0x70, 0xce, 0xa9,
	0xed, 0x58, 0xae, 0xcd, 0x0a, 0x44, 0xeb, 0x49, 0x1c, 0x85, 0xb1, 0xd1, 0xcb, 0x46, 0xcf, 0x09,
	0x21, 0x50, 0xf6, 0x83, 0x40, 0xd0, 0x8a, 0x63, 0xb9, 0x0d, 0x86, 0xcf, 0xe4, 0x0f, 0x68, 0x48,
	0x25, 0xc6, 0xc6, 0x46, 0x15, 0x85, 0x1c, 0x64, 0xaa, 0x31, 0x53, 0xcb, 0x55, 0xe3, 0xa7, 0x03,
	0xad, 0x48, 0xaa, 0x89, 0x58, 0x48, 0xe5, 0x25, 0x01, 0xa7, 0x75, 0xc7, 0x76, 0x1b, 0x6c, 0x83,
	0x75, 0xbe, 0xd4, 0xa0, 0xae, 0x1f, 0x06, 0xf1, 0xfb, 0x84, 0x9c, 0x40, 0x5d, 0x72, 0xb1, 0xbc,
	0xd4, 0x26, 0x2c, 0x3c, 0x6d, 0x1d, 0xaf, 0xcd, 0x95, 0x0a, 0xe6, 0x08, 0x94, 0x63, 0x7f, 0x55,
	0x6a, 0x83, 0xe1, 0x33, 0x71, 0xa0, 0x39, 0x4b, 0xe2, 0x98, 0xcf, 0x94, 0xb7, 0x98, 0x4b, 0xac,
	0xb2, 0xc2, 0x8a, 0x88, 0x9c, 0xc2, 0xce, 0x2a, 0xe4, 0x01, 0xe6, 0x54, 0x30, 0x67, 0x13, 0x92,
	0x13, 0xa8, 0xce, 0x94, 0x88, 0x06, 0x3d, 0xac, 0xda, 0x7e, 0x55, 0xa2, 0x16, 0x5b, 0x11, 0x7d,
	0x42, 0x24, 0x55, 0x37, 0x0a, 0x79, 0xac, 0xd0, 0x6c, 0x0d, 0x2b, 0xdb, 0x84, 0x84, 0x42, 0xcd,
	0x0f, 0x02, 0xec, 0x75, 0x1d, 0x7b, 0x9d, 0x85, 0xba, 0x6d, 0xab, 0x8f, 0x0d, 0x39, 0x6d, 0x38,
	0x96, 0x5b, 0x67, 0x39, 0x20, 0xa7, 0x79, 0x05, 0xba, 0x6b, 0xa0, 0x75, 0xfc, 0x7c, 0x11, 0x13,
	0x17, 0xf6, 0xe2, 0x24, 0xe0, 0x93, 0xa7, 0x94, 0xbf, 0xe1, 0x42, 0x86, 0x49, 0x4c, 0x9b, 0xd8,
	0x86, 0x6d, 0xac, 0xbb, 0x9a, 0x21, 0xda, 0x32, 0x5d, 0xcd, 0x62, 0xd3, 0x2d, 0xb1, 0x3e, 0x61,
	0x07, 0xe5, 0x22, 0xc2, 0x5a, 0x7d, 0xa9, 0x6e, 0x79, 0x1c, 0x0c, 0xe5, 0x87, 0x41, 0x8f, 0xee,
	0x62, 0x2d, 0x9b, 0x50, 0xbb, 0xd1, 0xa0, 0x6b, 0x0c, 0x62, 0xcd, 0x7b, 0x98, 0xb7, 0x8d, 0xc9,
	0xdf, 0xb0, 0x5f, 0x40, 0xdc, 0xf4, 0xa7, 0x8d, 0xb9, 0xcf, 0x85, 0xad, 0xec, 0xa1, 0x99, 0xdc,
	0xfd, 0x67, 0xd9, 0x46, 0xd0, 0x7d, 0x8d, 0xa4, 0x7a, 0x2d, 0x92, 0x45, 0x2a, 0x29, 0xc1, 0x3b,
	0xc9, 0x81, 0x1e, 0xff, 0x80, 0xa7, 0x82, 0xcf, 0x7c, 0xc5, 0x03, 0x7a, 0x80, 0x6d, 0x2f, 0x90,
	0xac, 0x52, 0xc6, 0x67, 0x4b, 0x53, 0xe9, 0x61, 0x5e, 0xe9, 0x1a, 0x92, 0x7f, 0xa1, 0xae, 0xaf,
	0x01, 0xbb, 0x79, 0xe4, 0x58, 0xee, 0xee, 0xc5, 0xd1, 0x79, 0xb6, 0xc1, 0xe7, 0xdd, 0x91, 0xe7,
	0xf5, 0xbb, 0x93, 0xc9, 0xbb, 0x71, 0x9f, 0xad, 0xd3, 0xc8, 0x5f, 0xb0, 0xbb, 0xf4, 0xa3, 0x30,
	0xd0, 0x66, 0xf5, 0xdd, 0x49, 0x7a, 0x8c, 0xde, 0xb6, 0x28, 0xf9, 0x07, 0x0e, 0x54, 0x38, 0xe7,
	0x52, 0xf9, 0xf3, 0xb4, 0x97, 0x3b, 0xfd, 0x15, 0x6d, 0x7c, 0x4f, 0xc2, 0x21, 0x58, 0xcc, 0xe5,
	0xaa, 0x0b, 0x57, 0x7e, 0x18, 0x51, 0x8a, 0xc3, 0xbc, 0x8d, 0xc9, 0x0b, 0x38, 0xd6, 0x75, 0x60,
	0x0d, 0xff, 0x33, 0xfe, 0x69, 0xc1, 0xcd, 0x02, 0x4a, 0xfa, 0x1b, 0x1e, 0xff, 0x03, 0x55, 0xef,
	0xb0, 0x1e, 0x16, 0x99, 0x4d, 0xc8, 0x09, 0x4e, 0xc8, 0x06, 0x23, 0x67, 0xd0, 0xd6, 0x6f, 0x7b,
	0xc5, 0xbc, 0xdf, 0x31, 0xef, 0x19, 0xef, 0x7c, 0xb6, 0xa0, 0x95, 0xed, 0xfb, 0x4d, 0x28, 0x95,
	0xbe, 0x15, 0x19, 0xfb, 0xa9, 0xfc, 0x98, 0xa8, 0x41, 0x0f, 0xb7, 0xde, 0x66, 0x05, 0xa2, 0xa7,
	0x97, 0xc7, 0xc1, 0x20, 0x0e, 0xf8, 0x23, 0xee, 0x7e, 0x85, 0xad, 0x63, 0xad, 0xcd, 0xfd, 0x47,
	0xa3, 0xd9, 0x46, 0xcb, 0x62, 0xe2, 0x42, 0x05, 0x4d, 0xd2, 0xb2, 0x63, 0xbb, 0xcd, 0x0b, 0x92,
	0x5f, 0x52, 0xf6, 0x79, 0x66, 0x12, 0x3a, 0x2f, 0x61, 0xa7, 0xe8, 0xe8, 0x22, 0x7f, 0xd5, 0xfa,
	0xd9, 0xab, 0x67, 0x00, 0x57, 0x61, 0xc4, 0x6f, 0x9f, 0xa4, 0xe2, 0x73, 0x5c, 0xeb, 0x85, 0xd0,
	0x60, 0x5d, 0x49, 0x0e, 0xce, 0xde, 0x42, 0xb3, 0x30, 0x1e, 0xa4, 0x0d, 0xad, 0x3b, 0xef, 0xda,
	0x1b, 0xdd, 0x7b, 0x0f, 0x1a, 0xb7, 0x7f, 0x21, 0x7b, 0xd0, 0xec, 0x0d, 0x58, 0xbf, 0x3b, 0x31,
	0xc0, 0xd2, 0x29, 0x57, 0x23, 0x76, 0x7f, 0xc9, 0x7a, 0x0f, 0x23, 0xaf, 0xdb, 0x6f, 0x97, 0xc8,
	0x21, 0xb4, 0x33, 0x32, 0xbc, 0xbb, 0x99, 0x0c, 0xc6, 0x37, 0xfd, 0xb6, 0x3d, 0xad, 0xe2, 0x8f,
	0xe4, 0xbf, 0x6f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x6b, 0xa1, 0x8f, 0xe2, 0x58, 0x06, 0x00, 0x00,
}
