// Code generated by protoc-gen-go. DO NOT EDIT.
// source: jarviscore.proto

package jarviscorepb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
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

type CHANNELTYPE int32

const (
	CHANNELTYPE_NODEINFO CHANNELTYPE = 0
	CHANNELTYPE_CTRL     CHANNELTYPE = 1
)

var CHANNELTYPE_name = map[int32]string{
	0: "NODEINFO",
	1: "CTRL",
}
var CHANNELTYPE_value = map[string]int32{
	"NODEINFO": 0,
	"CTRL":     1,
}

func (x CHANNELTYPE) String() string {
	return proto.EnumName(CHANNELTYPE_name, int32(x))
}
func (CHANNELTYPE) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_jarviscore_148b13d2a45d4b27, []int{0}
}

type REPLYTYPE int32

const (
	REPLYTYPE_NONE      REPLYTYPE = 0
	REPLYTYPE_FORWARD   REPLYTYPE = 1
	REPLYTYPE_BROADCAST REPLYTYPE = 2
)

var REPLYTYPE_name = map[int32]string{
	0: "NONE",
	1: "FORWARD",
	2: "BROADCAST",
}
var REPLYTYPE_value = map[string]int32{
	"NONE":      0,
	"FORWARD":   1,
	"BROADCAST": 2,
}

func (x REPLYTYPE) String() string {
	return proto.EnumName(REPLYTYPE_name, int32(x))
}
func (REPLYTYPE) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_jarviscore_148b13d2a45d4b27, []int{1}
}

type NodeBaseInfo struct {
	ServAddr             string   `protobuf:"bytes,1,opt,name=servAddr,proto3" json:"servAddr,omitempty"`
	Addr                 string   `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NodeBaseInfo) Reset()         { *m = NodeBaseInfo{} }
func (m *NodeBaseInfo) String() string { return proto.CompactTextString(m) }
func (*NodeBaseInfo) ProtoMessage()    {}
func (*NodeBaseInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_jarviscore_148b13d2a45d4b27, []int{0}
}
func (m *NodeBaseInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeBaseInfo.Unmarshal(m, b)
}
func (m *NodeBaseInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeBaseInfo.Marshal(b, m, deterministic)
}
func (dst *NodeBaseInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeBaseInfo.Merge(dst, src)
}
func (m *NodeBaseInfo) XXX_Size() int {
	return xxx_messageInfo_NodeBaseInfo.Size(m)
}
func (m *NodeBaseInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeBaseInfo.DiscardUnknown(m)
}

var xxx_messageInfo_NodeBaseInfo proto.InternalMessageInfo

func (m *NodeBaseInfo) GetServAddr() string {
	if m != nil {
		return m.ServAddr
	}
	return ""
}

func (m *NodeBaseInfo) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *NodeBaseInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type CtrlDataInDB struct {
	Ctrlid               int64    `protobuf:"varint,1,opt,name=ctrlid,proto3" json:"ctrlid,omitempty"`
	CtrlType             string   `protobuf:"bytes,2,opt,name=ctrlType,proto3" json:"ctrlType,omitempty"`
	ForwordAddr          string   `protobuf:"bytes,3,opt,name=forwordAddr,proto3" json:"forwordAddr,omitempty"`
	Command              []byte   `protobuf:"bytes,4,opt,name=command,proto3" json:"command,omitempty"`
	Result               []byte   `protobuf:"bytes,5,opt,name=result,proto3" json:"result,omitempty"`
	ForwordNums          int32    `protobuf:"varint,6,opt,name=forwordNums,proto3" json:"forwordNums,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CtrlDataInDB) Reset()         { *m = CtrlDataInDB{} }
func (m *CtrlDataInDB) String() string { return proto.CompactTextString(m) }
func (*CtrlDataInDB) ProtoMessage()    {}
func (*CtrlDataInDB) Descriptor() ([]byte, []int) {
	return fileDescriptor_jarviscore_148b13d2a45d4b27, []int{1}
}
func (m *CtrlDataInDB) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CtrlDataInDB.Unmarshal(m, b)
}
func (m *CtrlDataInDB) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CtrlDataInDB.Marshal(b, m, deterministic)
}
func (dst *CtrlDataInDB) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CtrlDataInDB.Merge(dst, src)
}
func (m *CtrlDataInDB) XXX_Size() int {
	return xxx_messageInfo_CtrlDataInDB.Size(m)
}
func (m *CtrlDataInDB) XXX_DiscardUnknown() {
	xxx_messageInfo_CtrlDataInDB.DiscardUnknown(m)
}

var xxx_messageInfo_CtrlDataInDB proto.InternalMessageInfo

func (m *CtrlDataInDB) GetCtrlid() int64 {
	if m != nil {
		return m.Ctrlid
	}
	return 0
}

func (m *CtrlDataInDB) GetCtrlType() string {
	if m != nil {
		return m.CtrlType
	}
	return ""
}

func (m *CtrlDataInDB) GetForwordAddr() string {
	if m != nil {
		return m.ForwordAddr
	}
	return ""
}

func (m *CtrlDataInDB) GetCommand() []byte {
	if m != nil {
		return m.Command
	}
	return nil
}

func (m *CtrlDataInDB) GetResult() []byte {
	if m != nil {
		return m.Result
	}
	return nil
}

func (m *CtrlDataInDB) GetForwordNums() int32 {
	if m != nil {
		return m.ForwordNums
	}
	return 0
}

type Join struct {
	ServAddr             string   `protobuf:"bytes,1,opt,name=servAddr,proto3" json:"servAddr,omitempty"`
	Addr                 string   `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Join) Reset()         { *m = Join{} }
func (m *Join) String() string { return proto.CompactTextString(m) }
func (*Join) ProtoMessage()    {}
func (*Join) Descriptor() ([]byte, []int) {
	return fileDescriptor_jarviscore_148b13d2a45d4b27, []int{2}
}
func (m *Join) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Join.Unmarshal(m, b)
}
func (m *Join) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Join.Marshal(b, m, deterministic)
}
func (dst *Join) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Join.Merge(dst, src)
}
func (m *Join) XXX_Size() int {
	return xxx_messageInfo_Join.Size(m)
}
func (m *Join) XXX_DiscardUnknown() {
	xxx_messageInfo_Join.DiscardUnknown(m)
}

var xxx_messageInfo_Join proto.InternalMessageInfo

func (m *Join) GetServAddr() string {
	if m != nil {
		return m.ServAddr
	}
	return ""
}

func (m *Join) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *Join) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type ReplyJoin struct {
	Err                  string   `protobuf:"bytes,1,opt,name=err,proto3" json:"err,omitempty"`
	Addr                 string   `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReplyJoin) Reset()         { *m = ReplyJoin{} }
func (m *ReplyJoin) String() string { return proto.CompactTextString(m) }
func (*ReplyJoin) ProtoMessage()    {}
func (*ReplyJoin) Descriptor() ([]byte, []int) {
	return fileDescriptor_jarviscore_148b13d2a45d4b27, []int{3}
}
func (m *ReplyJoin) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReplyJoin.Unmarshal(m, b)
}
func (m *ReplyJoin) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReplyJoin.Marshal(b, m, deterministic)
}
func (dst *ReplyJoin) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReplyJoin.Merge(dst, src)
}
func (m *ReplyJoin) XXX_Size() int {
	return xxx_messageInfo_ReplyJoin.Size(m)
}
func (m *ReplyJoin) XXX_DiscardUnknown() {
	xxx_messageInfo_ReplyJoin.DiscardUnknown(m)
}

var xxx_messageInfo_ReplyJoin proto.InternalMessageInfo

func (m *ReplyJoin) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

func (m *ReplyJoin) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *ReplyJoin) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type BaseReply struct {
	Err                  string    `protobuf:"bytes,1,opt,name=err,proto3" json:"err,omitempty"`
	ReplyType            REPLYTYPE `protobuf:"varint,2,opt,name=replyType,proto3,enum=jarviscorepb.REPLYTYPE" json:"replyType,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *BaseReply) Reset()         { *m = BaseReply{} }
func (m *BaseReply) String() string { return proto.CompactTextString(m) }
func (*BaseReply) ProtoMessage()    {}
func (*BaseReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_jarviscore_148b13d2a45d4b27, []int{4}
}
func (m *BaseReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BaseReply.Unmarshal(m, b)
}
func (m *BaseReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BaseReply.Marshal(b, m, deterministic)
}
func (dst *BaseReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BaseReply.Merge(dst, src)
}
func (m *BaseReply) XXX_Size() int {
	return xxx_messageInfo_BaseReply.Size(m)
}
func (m *BaseReply) XXX_DiscardUnknown() {
	xxx_messageInfo_BaseReply.DiscardUnknown(m)
}

var xxx_messageInfo_BaseReply proto.InternalMessageInfo

func (m *BaseReply) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

func (m *BaseReply) GetReplyType() REPLYTYPE {
	if m != nil {
		return m.ReplyType
	}
	return REPLYTYPE_NONE
}

type Subscribe struct {
	ChannelType          CHANNELTYPE `protobuf:"varint,1,opt,name=channelType,proto3,enum=jarviscorepb.CHANNELTYPE" json:"channelType,omitempty"`
	Addr                 string      `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Subscribe) Reset()         { *m = Subscribe{} }
func (m *Subscribe) String() string { return proto.CompactTextString(m) }
func (*Subscribe) ProtoMessage()    {}
func (*Subscribe) Descriptor() ([]byte, []int) {
	return fileDescriptor_jarviscore_148b13d2a45d4b27, []int{5}
}
func (m *Subscribe) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Subscribe.Unmarshal(m, b)
}
func (m *Subscribe) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Subscribe.Marshal(b, m, deterministic)
}
func (dst *Subscribe) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Subscribe.Merge(dst, src)
}
func (m *Subscribe) XXX_Size() int {
	return xxx_messageInfo_Subscribe.Size(m)
}
func (m *Subscribe) XXX_DiscardUnknown() {
	xxx_messageInfo_Subscribe.DiscardUnknown(m)
}

var xxx_messageInfo_Subscribe proto.InternalMessageInfo

func (m *Subscribe) GetChannelType() CHANNELTYPE {
	if m != nil {
		return m.ChannelType
	}
	return CHANNELTYPE_NODEINFO
}

func (m *Subscribe) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

type CtrlInfo struct {
	Ctrlid               int64    `protobuf:"varint,1,opt,name=ctrlid,proto3" json:"ctrlid,omitempty"`
	DestAddr             string   `protobuf:"bytes,2,opt,name=destAddr,proto3" json:"destAddr,omitempty"`
	SrcAddr              string   `protobuf:"bytes,3,opt,name=srcAddr,proto3" json:"srcAddr,omitempty"`
	MyAddr               string   `protobuf:"bytes,4,opt,name=myAddr,proto3" json:"myAddr,omitempty"`
	CtrlType             string   `protobuf:"bytes,5,opt,name=ctrlType,proto3" json:"ctrlType,omitempty"`
	Command              []byte   `protobuf:"bytes,6,opt,name=command,proto3" json:"command,omitempty"`
	ForwordNums          int32    `protobuf:"varint,7,opt,name=forwordNums,proto3" json:"forwordNums,omitempty"`
	SignR                []byte   `protobuf:"bytes,8,opt,name=signR,proto3" json:"signR,omitempty"`
	SignS                []byte   `protobuf:"bytes,9,opt,name=signS,proto3" json:"signS,omitempty"`
	PubKey               []byte   `protobuf:"bytes,10,opt,name=pubKey,proto3" json:"pubKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CtrlInfo) Reset()         { *m = CtrlInfo{} }
func (m *CtrlInfo) String() string { return proto.CompactTextString(m) }
func (*CtrlInfo) ProtoMessage()    {}
func (*CtrlInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_jarviscore_148b13d2a45d4b27, []int{6}
}
func (m *CtrlInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CtrlInfo.Unmarshal(m, b)
}
func (m *CtrlInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CtrlInfo.Marshal(b, m, deterministic)
}
func (dst *CtrlInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CtrlInfo.Merge(dst, src)
}
func (m *CtrlInfo) XXX_Size() int {
	return xxx_messageInfo_CtrlInfo.Size(m)
}
func (m *CtrlInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_CtrlInfo.DiscardUnknown(m)
}

var xxx_messageInfo_CtrlInfo proto.InternalMessageInfo

func (m *CtrlInfo) GetCtrlid() int64 {
	if m != nil {
		return m.Ctrlid
	}
	return 0
}

func (m *CtrlInfo) GetDestAddr() string {
	if m != nil {
		return m.DestAddr
	}
	return ""
}

func (m *CtrlInfo) GetSrcAddr() string {
	if m != nil {
		return m.SrcAddr
	}
	return ""
}

func (m *CtrlInfo) GetMyAddr() string {
	if m != nil {
		return m.MyAddr
	}
	return ""
}

func (m *CtrlInfo) GetCtrlType() string {
	if m != nil {
		return m.CtrlType
	}
	return ""
}

func (m *CtrlInfo) GetCommand() []byte {
	if m != nil {
		return m.Command
	}
	return nil
}

func (m *CtrlInfo) GetForwordNums() int32 {
	if m != nil {
		return m.ForwordNums
	}
	return 0
}

func (m *CtrlInfo) GetSignR() []byte {
	if m != nil {
		return m.SignR
	}
	return nil
}

func (m *CtrlInfo) GetSignS() []byte {
	if m != nil {
		return m.SignS
	}
	return nil
}

func (m *CtrlInfo) GetPubKey() []byte {
	if m != nil {
		return m.PubKey
	}
	return nil
}

type CtrlResult struct {
	Ctrlid               int64    `protobuf:"varint,1,opt,name=ctrlid,proto3" json:"ctrlid,omitempty"`
	DestAddr             string   `protobuf:"bytes,2,opt,name=destAddr,proto3" json:"destAddr,omitempty"`
	SrcAddr              string   `protobuf:"bytes,3,opt,name=srcAddr,proto3" json:"srcAddr,omitempty"`
	MyAddr               string   `protobuf:"bytes,4,opt,name=myAddr,proto3" json:"myAddr,omitempty"`
	CtrlResult           []byte   `protobuf:"bytes,5,opt,name=ctrlResult,proto3" json:"ctrlResult,omitempty"`
	ForwordNums          int32    `protobuf:"varint,6,opt,name=forwordNums,proto3" json:"forwordNums,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CtrlResult) Reset()         { *m = CtrlResult{} }
func (m *CtrlResult) String() string { return proto.CompactTextString(m) }
func (*CtrlResult) ProtoMessage()    {}
func (*CtrlResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_jarviscore_148b13d2a45d4b27, []int{7}
}
func (m *CtrlResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CtrlResult.Unmarshal(m, b)
}
func (m *CtrlResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CtrlResult.Marshal(b, m, deterministic)
}
func (dst *CtrlResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CtrlResult.Merge(dst, src)
}
func (m *CtrlResult) XXX_Size() int {
	return xxx_messageInfo_CtrlResult.Size(m)
}
func (m *CtrlResult) XXX_DiscardUnknown() {
	xxx_messageInfo_CtrlResult.DiscardUnknown(m)
}

var xxx_messageInfo_CtrlResult proto.InternalMessageInfo

func (m *CtrlResult) GetCtrlid() int64 {
	if m != nil {
		return m.Ctrlid
	}
	return 0
}

func (m *CtrlResult) GetDestAddr() string {
	if m != nil {
		return m.DestAddr
	}
	return ""
}

func (m *CtrlResult) GetSrcAddr() string {
	if m != nil {
		return m.SrcAddr
	}
	return ""
}

func (m *CtrlResult) GetMyAddr() string {
	if m != nil {
		return m.MyAddr
	}
	return ""
}

func (m *CtrlResult) GetCtrlResult() []byte {
	if m != nil {
		return m.CtrlResult
	}
	return nil
}

func (m *CtrlResult) GetForwordNums() int32 {
	if m != nil {
		return m.ForwordNums
	}
	return 0
}

type ChannelInfo struct {
	ChannelType CHANNELTYPE `protobuf:"varint,1,opt,name=channelType,proto3,enum=jarviscorepb.CHANNELTYPE" json:"channelType,omitempty"`
	// Types that are valid to be assigned to Data:
	//	*ChannelInfo_NodeInfo
	//	*ChannelInfo_CtrlInfo
	Data                 isChannelInfo_Data `protobuf_oneof:"data"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *ChannelInfo) Reset()         { *m = ChannelInfo{} }
func (m *ChannelInfo) String() string { return proto.CompactTextString(m) }
func (*ChannelInfo) ProtoMessage()    {}
func (*ChannelInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_jarviscore_148b13d2a45d4b27, []int{8}
}
func (m *ChannelInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChannelInfo.Unmarshal(m, b)
}
func (m *ChannelInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChannelInfo.Marshal(b, m, deterministic)
}
func (dst *ChannelInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChannelInfo.Merge(dst, src)
}
func (m *ChannelInfo) XXX_Size() int {
	return xxx_messageInfo_ChannelInfo.Size(m)
}
func (m *ChannelInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ChannelInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ChannelInfo proto.InternalMessageInfo

type isChannelInfo_Data interface {
	isChannelInfo_Data()
}

type ChannelInfo_NodeInfo struct {
	NodeInfo *NodeBaseInfo `protobuf:"bytes,2,opt,name=nodeInfo,proto3,oneof"`
}
type ChannelInfo_CtrlInfo struct {
	CtrlInfo *CtrlInfo `protobuf:"bytes,3,opt,name=ctrlInfo,proto3,oneof"`
}

func (*ChannelInfo_NodeInfo) isChannelInfo_Data() {}
func (*ChannelInfo_CtrlInfo) isChannelInfo_Data() {}

func (m *ChannelInfo) GetData() isChannelInfo_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *ChannelInfo) GetChannelType() CHANNELTYPE {
	if m != nil {
		return m.ChannelType
	}
	return CHANNELTYPE_NODEINFO
}

func (m *ChannelInfo) GetNodeInfo() *NodeBaseInfo {
	if x, ok := m.GetData().(*ChannelInfo_NodeInfo); ok {
		return x.NodeInfo
	}
	return nil
}

func (m *ChannelInfo) GetCtrlInfo() *CtrlInfo {
	if x, ok := m.GetData().(*ChannelInfo_CtrlInfo); ok {
		return x.CtrlInfo
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*ChannelInfo) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _ChannelInfo_OneofMarshaler, _ChannelInfo_OneofUnmarshaler, _ChannelInfo_OneofSizer, []interface{}{
		(*ChannelInfo_NodeInfo)(nil),
		(*ChannelInfo_CtrlInfo)(nil),
	}
}

func _ChannelInfo_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*ChannelInfo)
	// data
	switch x := m.Data.(type) {
	case *ChannelInfo_NodeInfo:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.NodeInfo); err != nil {
			return err
		}
	case *ChannelInfo_CtrlInfo:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.CtrlInfo); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("ChannelInfo.Data has unexpected type %T", x)
	}
	return nil
}

func _ChannelInfo_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*ChannelInfo)
	switch tag {
	case 2: // data.nodeInfo
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(NodeBaseInfo)
		err := b.DecodeMessage(msg)
		m.Data = &ChannelInfo_NodeInfo{msg}
		return true, err
	case 3: // data.ctrlInfo
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(CtrlInfo)
		err := b.DecodeMessage(msg)
		m.Data = &ChannelInfo_CtrlInfo{msg}
		return true, err
	default:
		return false, nil
	}
}

func _ChannelInfo_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*ChannelInfo)
	// data
	switch x := m.Data.(type) {
	case *ChannelInfo_NodeInfo:
		s := proto.Size(x.NodeInfo)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *ChannelInfo_CtrlInfo:
		s := proto.Size(x.CtrlInfo)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type ServAddr struct {
	ServAddr             string   `protobuf:"bytes,1,opt,name=servAddr,proto3" json:"servAddr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ServAddr) Reset()         { *m = ServAddr{} }
func (m *ServAddr) String() string { return proto.CompactTextString(m) }
func (*ServAddr) ProtoMessage()    {}
func (*ServAddr) Descriptor() ([]byte, []int) {
	return fileDescriptor_jarviscore_148b13d2a45d4b27, []int{9}
}
func (m *ServAddr) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServAddr.Unmarshal(m, b)
}
func (m *ServAddr) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServAddr.Marshal(b, m, deterministic)
}
func (dst *ServAddr) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServAddr.Merge(dst, src)
}
func (m *ServAddr) XXX_Size() int {
	return xxx_messageInfo_ServAddr.Size(m)
}
func (m *ServAddr) XXX_DiscardUnknown() {
	xxx_messageInfo_ServAddr.DiscardUnknown(m)
}

var xxx_messageInfo_ServAddr proto.InternalMessageInfo

func (m *ServAddr) GetServAddr() string {
	if m != nil {
		return m.ServAddr
	}
	return ""
}

type TrustNode struct {
	Addr                 string   `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TrustNode) Reset()         { *m = TrustNode{} }
func (m *TrustNode) String() string { return proto.CompactTextString(m) }
func (*TrustNode) ProtoMessage()    {}
func (*TrustNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_jarviscore_148b13d2a45d4b27, []int{10}
}
func (m *TrustNode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TrustNode.Unmarshal(m, b)
}
func (m *TrustNode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TrustNode.Marshal(b, m, deterministic)
}
func (dst *TrustNode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TrustNode.Merge(dst, src)
}
func (m *TrustNode) XXX_Size() int {
	return xxx_messageInfo_TrustNode.Size(m)
}
func (m *TrustNode) XXX_DiscardUnknown() {
	xxx_messageInfo_TrustNode.DiscardUnknown(m)
}

var xxx_messageInfo_TrustNode proto.InternalMessageInfo

func (m *TrustNode) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func init() {
	proto.RegisterType((*NodeBaseInfo)(nil), "jarviscorepb.NodeBaseInfo")
	proto.RegisterType((*CtrlDataInDB)(nil), "jarviscorepb.CtrlDataInDB")
	proto.RegisterType((*Join)(nil), "jarviscorepb.Join")
	proto.RegisterType((*ReplyJoin)(nil), "jarviscorepb.ReplyJoin")
	proto.RegisterType((*BaseReply)(nil), "jarviscorepb.BaseReply")
	proto.RegisterType((*Subscribe)(nil), "jarviscorepb.Subscribe")
	proto.RegisterType((*CtrlInfo)(nil), "jarviscorepb.CtrlInfo")
	proto.RegisterType((*CtrlResult)(nil), "jarviscorepb.CtrlResult")
	proto.RegisterType((*ChannelInfo)(nil), "jarviscorepb.ChannelInfo")
	proto.RegisterType((*ServAddr)(nil), "jarviscorepb.ServAddr")
	proto.RegisterType((*TrustNode)(nil), "jarviscorepb.TrustNode")
	proto.RegisterEnum("jarviscorepb.CHANNELTYPE", CHANNELTYPE_name, CHANNELTYPE_value)
	proto.RegisterEnum("jarviscorepb.REPLYTYPE", REPLYTYPE_name, REPLYTYPE_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// JarvisCoreServClient is the client API for JarvisCoreServ service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type JarvisCoreServClient interface {
	Join(ctx context.Context, in *Join, opts ...grpc.CallOption) (*ReplyJoin, error)
	RequestCtrl(ctx context.Context, in *CtrlInfo, opts ...grpc.CallOption) (*BaseReply, error)
	ReplyCtrl(ctx context.Context, in *CtrlResult, opts ...grpc.CallOption) (*BaseReply, error)
	Subscribe(ctx context.Context, in *Subscribe, opts ...grpc.CallOption) (JarvisCoreServ_SubscribeClient, error)
	GetMyServAddr(ctx context.Context, in *ServAddr, opts ...grpc.CallOption) (*ServAddr, error)
	Trust(ctx context.Context, in *TrustNode, opts ...grpc.CallOption) (*BaseReply, error)
}

type jarvisCoreServClient struct {
	cc *grpc.ClientConn
}

func NewJarvisCoreServClient(cc *grpc.ClientConn) JarvisCoreServClient {
	return &jarvisCoreServClient{cc}
}

func (c *jarvisCoreServClient) Join(ctx context.Context, in *Join, opts ...grpc.CallOption) (*ReplyJoin, error) {
	out := new(ReplyJoin)
	err := c.cc.Invoke(ctx, "/jarviscorepb.JarvisCoreServ/join", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jarvisCoreServClient) RequestCtrl(ctx context.Context, in *CtrlInfo, opts ...grpc.CallOption) (*BaseReply, error) {
	out := new(BaseReply)
	err := c.cc.Invoke(ctx, "/jarviscorepb.JarvisCoreServ/requestCtrl", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jarvisCoreServClient) ReplyCtrl(ctx context.Context, in *CtrlResult, opts ...grpc.CallOption) (*BaseReply, error) {
	out := new(BaseReply)
	err := c.cc.Invoke(ctx, "/jarviscorepb.JarvisCoreServ/replyCtrl", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jarvisCoreServClient) Subscribe(ctx context.Context, in *Subscribe, opts ...grpc.CallOption) (JarvisCoreServ_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &_JarvisCoreServ_serviceDesc.Streams[0], "/jarviscorepb.JarvisCoreServ/subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &jarvisCoreServSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type JarvisCoreServ_SubscribeClient interface {
	Recv() (*ChannelInfo, error)
	grpc.ClientStream
}

type jarvisCoreServSubscribeClient struct {
	grpc.ClientStream
}

func (x *jarvisCoreServSubscribeClient) Recv() (*ChannelInfo, error) {
	m := new(ChannelInfo)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *jarvisCoreServClient) GetMyServAddr(ctx context.Context, in *ServAddr, opts ...grpc.CallOption) (*ServAddr, error) {
	out := new(ServAddr)
	err := c.cc.Invoke(ctx, "/jarviscorepb.JarvisCoreServ/getMyServAddr", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jarvisCoreServClient) Trust(ctx context.Context, in *TrustNode, opts ...grpc.CallOption) (*BaseReply, error) {
	out := new(BaseReply)
	err := c.cc.Invoke(ctx, "/jarviscorepb.JarvisCoreServ/trust", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// JarvisCoreServServer is the server API for JarvisCoreServ service.
type JarvisCoreServServer interface {
	Join(context.Context, *Join) (*ReplyJoin, error)
	RequestCtrl(context.Context, *CtrlInfo) (*BaseReply, error)
	ReplyCtrl(context.Context, *CtrlResult) (*BaseReply, error)
	Subscribe(*Subscribe, JarvisCoreServ_SubscribeServer) error
	GetMyServAddr(context.Context, *ServAddr) (*ServAddr, error)
	Trust(context.Context, *TrustNode) (*BaseReply, error)
}

func RegisterJarvisCoreServServer(s *grpc.Server, srv JarvisCoreServServer) {
	s.RegisterService(&_JarvisCoreServ_serviceDesc, srv)
}

func _JarvisCoreServ_Join_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Join)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JarvisCoreServServer).Join(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jarviscorepb.JarvisCoreServ/Join",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JarvisCoreServServer).Join(ctx, req.(*Join))
	}
	return interceptor(ctx, in, info, handler)
}

func _JarvisCoreServ_RequestCtrl_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CtrlInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JarvisCoreServServer).RequestCtrl(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jarviscorepb.JarvisCoreServ/RequestCtrl",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JarvisCoreServServer).RequestCtrl(ctx, req.(*CtrlInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _JarvisCoreServ_ReplyCtrl_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CtrlResult)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JarvisCoreServServer).ReplyCtrl(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jarviscorepb.JarvisCoreServ/ReplyCtrl",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JarvisCoreServServer).ReplyCtrl(ctx, req.(*CtrlResult))
	}
	return interceptor(ctx, in, info, handler)
}

func _JarvisCoreServ_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Subscribe)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(JarvisCoreServServer).Subscribe(m, &jarvisCoreServSubscribeServer{stream})
}

type JarvisCoreServ_SubscribeServer interface {
	Send(*ChannelInfo) error
	grpc.ServerStream
}

type jarvisCoreServSubscribeServer struct {
	grpc.ServerStream
}

func (x *jarvisCoreServSubscribeServer) Send(m *ChannelInfo) error {
	return x.ServerStream.SendMsg(m)
}

func _JarvisCoreServ_GetMyServAddr_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServAddr)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JarvisCoreServServer).GetMyServAddr(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jarviscorepb.JarvisCoreServ/GetMyServAddr",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JarvisCoreServServer).GetMyServAddr(ctx, req.(*ServAddr))
	}
	return interceptor(ctx, in, info, handler)
}

func _JarvisCoreServ_Trust_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TrustNode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JarvisCoreServServer).Trust(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jarviscorepb.JarvisCoreServ/Trust",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JarvisCoreServServer).Trust(ctx, req.(*TrustNode))
	}
	return interceptor(ctx, in, info, handler)
}

var _JarvisCoreServ_serviceDesc = grpc.ServiceDesc{
	ServiceName: "jarviscorepb.JarvisCoreServ",
	HandlerType: (*JarvisCoreServServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "join",
			Handler:    _JarvisCoreServ_Join_Handler,
		},
		{
			MethodName: "requestCtrl",
			Handler:    _JarvisCoreServ_RequestCtrl_Handler,
		},
		{
			MethodName: "replyCtrl",
			Handler:    _JarvisCoreServ_ReplyCtrl_Handler,
		},
		{
			MethodName: "getMyServAddr",
			Handler:    _JarvisCoreServ_GetMyServAddr_Handler,
		},
		{
			MethodName: "trust",
			Handler:    _JarvisCoreServ_Trust_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "subscribe",
			Handler:       _JarvisCoreServ_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "jarviscore.proto",
}

func init() { proto.RegisterFile("jarviscore.proto", fileDescriptor_jarviscore_148b13d2a45d4b27) }

var fileDescriptor_jarviscore_148b13d2a45d4b27 = []byte{
	// 700 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x55, 0x51, 0x6b, 0xdb, 0x3c,
	0x14, 0x8d, 0x1a, 0x27, 0xb5, 0xaf, 0xd3, 0x62, 0xc4, 0x47, 0x3f, 0x2f, 0x0f, 0x5b, 0x30, 0x6c,
	0x94, 0x3e, 0x94, 0xad, 0x5b, 0x61, 0xd0, 0x97, 0x39, 0x4e, 0x4a, 0xdb, 0x75, 0x4e, 0x51, 0x02,
	0xa3, 0xb0, 0x17, 0xc7, 0x56, 0xbb, 0x94, 0xc4, 0xce, 0x24, 0xa7, 0x23, 0xff, 0x61, 0x0f, 0xfb,
	0x31, 0x7b, 0xd8, 0xf3, 0x7e, 0xd9, 0x90, 0x1c, 0x2b, 0x8e, 0xdb, 0x86, 0xc2, 0xc6, 0xde, 0x7c,
	0x8f, 0x74, 0x8f, 0x74, 0x75, 0xce, 0xbd, 0x06, 0xeb, 0x26, 0x60, 0xb7, 0x23, 0x1e, 0x26, 0x8c,
	0xee, 0x4f, 0x59, 0x92, 0x26, 0xb8, 0xb1, 0x44, 0xa6, 0x43, 0x87, 0x40, 0xc3, 0x4f, 0x22, 0xda,
	0x0e, 0x38, 0x3d, 0x8d, 0xaf, 0x12, 0xdc, 0x04, 0x9d, 0x53, 0x76, 0xeb, 0x46, 0x11, 0xb3, 0x51,
	0x0b, 0xed, 0x1a, 0x44, 0xc5, 0x18, 0x83, 0x16, 0x08, 0x7c, 0x43, 0xe2, 0xf2, 0x5b, 0x60, 0x71,
	0x30, 0xa1, 0x76, 0x35, 0xc3, 0xc4, 0xb7, 0xf3, 0x13, 0x41, 0xc3, 0x4b, 0xd9, 0xb8, 0x13, 0xa4,
	0xc1, 0x69, 0xdc, 0x69, 0xe3, 0x1d, 0xa8, 0x87, 0x29, 0x1b, 0x8f, 0x22, 0x49, 0x59, 0x25, 0x8b,
	0x48, 0x1c, 0x26, 0xbe, 0x06, 0xf3, 0x29, 0x5d, 0x90, 0xaa, 0x18, 0xb7, 0xc0, 0xbc, 0x4a, 0xd8,
	0xd7, 0x84, 0x45, 0xf2, 0x2e, 0x19, 0x7f, 0x11, 0xc2, 0x36, 0x6c, 0x86, 0xc9, 0x64, 0x12, 0xc4,
	0x91, 0xad, 0xb5, 0xd0, 0x6e, 0x83, 0xe4, 0xa1, 0x38, 0x8f, 0x51, 0x3e, 0x1b, 0xa7, 0x76, 0x4d,
	0x2e, 0x2c, 0xa2, 0x02, 0xa7, 0x3f, 0x9b, 0x70, 0xbb, 0xde, 0x42, 0xbb, 0x35, 0x52, 0x84, 0x9c,
	0x33, 0xd0, 0xce, 0x92, 0x51, 0xfc, 0x57, 0x9e, 0xa1, 0x0b, 0x06, 0xa1, 0xd3, 0xf1, 0x5c, 0x12,
	0x5a, 0x50, 0xa5, 0x2c, 0xe7, 0x12, 0x9f, 0x8f, 0xa6, 0x19, 0x80, 0x21, 0xd4, 0x91, 0x54, 0xf7,
	0xd0, 0x1c, 0x82, 0xc1, 0xc4, 0x92, 0x7a, 0xc4, 0xed, 0x83, 0xff, 0xf7, 0x8b, 0x12, 0xef, 0x93,
	0xee, 0xc5, 0xf9, 0xe5, 0xe0, 0xf2, 0xa2, 0x4b, 0x96, 0x3b, 0x9d, 0x4f, 0x60, 0xf4, 0x67, 0x43,
	0x1e, 0xb2, 0xd1, 0x90, 0xe2, 0x23, 0x30, 0xc3, 0xcf, 0x41, 0x1c, 0xd3, 0x4c, 0x0a, 0x24, 0x59,
	0x9e, 0xac, 0xb2, 0x78, 0x27, 0xae, 0xef, 0x77, 0xcf, 0x25, 0x4f, 0x71, 0xf7, 0x7d, 0x75, 0x38,
	0xdf, 0x36, 0x40, 0x17, 0x0e, 0x90, 0x96, 0x5a, 0xa3, 0x7e, 0x44, 0x79, 0xea, 0x2e, 0x93, 0x55,
	0x2c, 0xb4, 0xe5, 0x2c, 0x2c, 0x28, 0x9f, 0x87, 0x82, 0x6d, 0x32, 0x97, 0x0b, 0x9a, 0x5c, 0x58,
	0x44, 0x2b, 0x5e, 0xaa, 0x95, 0xbc, 0x54, 0x70, 0x4a, 0x7d, 0xd5, 0x29, 0x25, 0x47, 0x6c, 0xde,
	0x71, 0x04, 0xfe, 0x0f, 0x6a, 0x7c, 0x74, 0x1d, 0x13, 0x5b, 0x97, 0x99, 0x59, 0x90, 0xa3, 0x7d,
	0xdb, 0x58, 0xa2, 0x7d, 0x71, 0xb7, 0xe9, 0x6c, 0xf8, 0x9e, 0xce, 0x6d, 0xc8, 0x7c, 0x97, 0x45,
	0xce, 0x0f, 0x04, 0x20, 0x9e, 0x83, 0x64, 0x36, 0xfc, 0x37, 0x0f, 0xf2, 0x14, 0x20, 0x54, 0x67,
	0x2e, 0x1a, 0xa1, 0x80, 0x3c, 0xa2, 0x19, 0x7e, 0x21, 0x30, 0xbd, 0x4c, 0x69, 0x29, 0xe4, 0x1f,
	0xd9, 0xe4, 0x2d, 0xe8, 0x71, 0x12, 0xc9, 0x21, 0x23, 0x8b, 0x33, 0x0f, 0x9a, 0xab, 0x99, 0xc5,
	0x31, 0x74, 0x52, 0x21, 0x6a, 0x37, 0x7e, 0x93, 0x29, 0x2b, 0x33, 0xab, 0x32, 0x73, 0xa7, 0x74,
	0xe6, 0x62, 0x55, 0x64, 0xe5, 0x3b, 0xdb, 0x75, 0xd0, 0xa2, 0x20, 0x0d, 0x9c, 0x17, 0xa0, 0xf7,
	0xf3, 0xce, 0x5d, 0xd3, 0xd5, 0xce, 0x33, 0x30, 0x06, 0x6c, 0xc6, 0x53, 0x71, 0x0d, 0xe5, 0x69,
	0xb4, 0xf4, 0xf4, 0xde, 0x73, 0x30, 0x0b, 0xc5, 0xe1, 0x06, 0xe8, 0x7e, 0xaf, 0xd3, 0x3d, 0xf5,
	0x8f, 0x7b, 0x56, 0x05, 0xeb, 0xa0, 0x79, 0x03, 0x72, 0x6e, 0xa1, 0xbd, 0x57, 0x60, 0xa8, 0x86,
	0x13, 0xb0, 0xdf, 0xf3, 0xbb, 0x56, 0x05, 0x9b, 0xb0, 0x79, 0xdc, 0x23, 0x1f, 0x5d, 0xd2, 0xb1,
	0x10, 0xde, 0x02, 0xa3, 0x4d, 0x7a, 0x6e, 0xc7, 0x73, 0xfb, 0x03, 0x6b, 0xe3, 0xe0, 0x7b, 0x15,
	0xb6, 0xcf, 0x64, 0x41, 0x5e, 0xc2, 0xa8, 0xb8, 0x2d, 0x3e, 0x04, 0xed, 0x46, 0x8c, 0x0d, 0xbc,
	0x5a, 0xa9, 0x18, 0x25, 0xcd, 0x72, 0x7b, 0xe7, 0x33, 0xc6, 0xa9, 0xe0, 0x77, 0x60, 0x32, 0xfa,
	0x65, 0x46, 0x79, 0x2a, 0xde, 0x04, 0x3f, 0xf0, 0x4e, 0x65, 0x06, 0x35, 0x5e, 0x24, 0x43, 0x36,
	0x24, 0x64, 0xbe, 0x7d, 0x37, 0x3f, 0x33, 0xcf, 0x3a, 0x06, 0x0f, 0x0c, 0xae, 0x26, 0x4b, 0x69,
	0x9f, 0x1a, 0x39, 0xcd, 0xb2, 0x6d, 0x96, 0x36, 0x73, 0x2a, 0x2f, 0x11, 0x76, 0x61, 0xeb, 0x9a,
	0xa6, 0x1f, 0xe6, 0x4a, 0xba, 0x52, 0x29, 0x39, 0xde, 0x7c, 0x00, 0x77, 0x2a, 0xf8, 0x08, 0x6a,
	0xa9, 0x10, 0xb4, 0x7c, 0x07, 0xa5, 0xf2, 0x9a, 0x22, 0x86, 0x75, 0xf9, 0xaf, 0x7c, 0xfd, 0x3b,
	0x00, 0x00, 0xff, 0xff, 0x2d, 0x03, 0xdd, 0x98, 0x3f, 0x07, 0x00, 0x00,
}
