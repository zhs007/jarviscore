// Code generated by protoc-gen-go. DO NOT EDIT.
// source: jarviscore.proto

package jarviscorepb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import any "github.com/golang/protobuf/ptypes/any"

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

// jarvis msg type
type MSGTYPE int32

const (
	// null type
	MSGTYPE_NULL_TYPE MSGTYPE = 0
	// connect jarvis node
	MSGTYPE_CONNECT_NODE MSGTYPE = 1
	// request a node's child nodes
	MSGTYPE_REQUEST_NODES MSGTYPE = 2
	// forward message
	MSGTYPE_FORWARD_MSG MSGTYPE = 3
	// request a ctrl msg
	MSGTYPE_REQUEST_CTRL MSGTYPE = 4
	// reply a ctrl result
	MSGTYPE_REPLY_CTRL_RESULT MSGTYPE = 5
	// trusted node request you trust a other node
	MSGTYPE_TRUST_NODE MSGTYPE = 6
	// trusted node request you remove a trusted node
	MSGTYPE_RM_TRUST_NODE MSGTYPE = 7
	// reply msg
	MSGTYPE_REPLY MSGTYPE = 8
	// node info
	MSGTYPE_NODE_INFO MSGTYPE = 9
	// reply connect
	MSGTYPE_REPLY_CONNECT MSGTYPE = 10
	// local connect other
	MSGTYPE_LOCAL_CONNECT_OTHER MSGTYPE = 11
	// local sendmsg
	MSGTYPE_LOCAL_SENDMSG MSGTYPE = 12
	// local request nodes child
	MSGTYPE_LOCAL_REQUEST_NODES MSGTYPE = 13
	// transfer file
	MSGTYPE_TRANSFER_FILE MSGTYPE = 14
)

var MSGTYPE_name = map[int32]string{
	0:  "NULL_TYPE",
	1:  "CONNECT_NODE",
	2:  "REQUEST_NODES",
	3:  "FORWARD_MSG",
	4:  "REQUEST_CTRL",
	5:  "REPLY_CTRL_RESULT",
	6:  "TRUST_NODE",
	7:  "RM_TRUST_NODE",
	8:  "REPLY",
	9:  "NODE_INFO",
	10: "REPLY_CONNECT",
	11: "LOCAL_CONNECT_OTHER",
	12: "LOCAL_SENDMSG",
	13: "LOCAL_REQUEST_NODES",
	14: "TRANSFER_FILE",
}
var MSGTYPE_value = map[string]int32{
	"NULL_TYPE":           0,
	"CONNECT_NODE":        1,
	"REQUEST_NODES":       2,
	"FORWARD_MSG":         3,
	"REQUEST_CTRL":        4,
	"REPLY_CTRL_RESULT":   5,
	"TRUST_NODE":          6,
	"RM_TRUST_NODE":       7,
	"REPLY":               8,
	"NODE_INFO":           9,
	"REPLY_CONNECT":       10,
	"LOCAL_CONNECT_OTHER": 11,
	"LOCAL_SENDMSG":       12,
	"LOCAL_REQUEST_NODES": 13,
	"TRANSFER_FILE":       14,
}

func (x MSGTYPE) String() string {
	return proto.EnumName(MSGTYPE_name, int32(x))
}
func (MSGTYPE) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_jarviscore_6f9bb307b6786fc3, []int{0}
}

// REPLYTYPE - reply type
type REPLYTYPE int32

const (
	// reply ok
	REPLYTYPE_OK REPLYTYPE = 0
	// forward
	REPLYTYPE_FORWARD REPLYTYPE = 1
	// timeout
	REPLYTYPE_TIMEOUT REPLYTYPE = 2
	// is me
	REPLYTYPE_ISME REPLYTYPE = 3
	// error
	REPLYTYPE_ERROR REPLYTYPE = 4
)

var REPLYTYPE_name = map[int32]string{
	0: "OK",
	1: "FORWARD",
	2: "TIMEOUT",
	3: "ISME",
	4: "ERROR",
}
var REPLYTYPE_value = map[string]int32{
	"OK":      0,
	"FORWARD": 1,
	"TIMEOUT": 2,
	"ISME":    3,
	"ERROR":   4,
}

func (x REPLYTYPE) String() string {
	return proto.EnumName(REPLYTYPE_name, int32(x))
}
func (REPLYTYPE) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_jarviscore_6f9bb307b6786fc3, []int{1}
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
	return fileDescriptor_jarviscore_6f9bb307b6786fc3, []int{0}
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

// ConnectInfo - used in LOCAL_CONNECT_OTHER, CONNECT_NODE
type ConnectInfo struct {
	ServAddr             string        `protobuf:"bytes,1,opt,name=servAddr,proto3" json:"servAddr,omitempty"`
	MyInfo               *NodeBaseInfo `protobuf:"bytes,2,opt,name=myInfo,proto3" json:"myInfo,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *ConnectInfo) Reset()         { *m = ConnectInfo{} }
func (m *ConnectInfo) String() string { return proto.CompactTextString(m) }
func (*ConnectInfo) ProtoMessage()    {}
func (*ConnectInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_jarviscore_6f9bb307b6786fc3, []int{1}
}
func (m *ConnectInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConnectInfo.Unmarshal(m, b)
}
func (m *ConnectInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConnectInfo.Marshal(b, m, deterministic)
}
func (dst *ConnectInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConnectInfo.Merge(dst, src)
}
func (m *ConnectInfo) XXX_Size() int {
	return xxx_messageInfo_ConnectInfo.Size(m)
}
func (m *ConnectInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ConnectInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ConnectInfo proto.InternalMessageInfo

func (m *ConnectInfo) GetServAddr() string {
	if m != nil {
		return m.ServAddr
	}
	return ""
}

func (m *ConnectInfo) GetMyInfo() *NodeBaseInfo {
	if m != nil {
		return m.MyInfo
	}
	return nil
}

type ReplyJoin struct {
	Addr                 string   `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReplyJoin) Reset()         { *m = ReplyJoin{} }
func (m *ReplyJoin) String() string { return proto.CompactTextString(m) }
func (*ReplyJoin) ProtoMessage()    {}
func (*ReplyJoin) Descriptor() ([]byte, []int) {
	return fileDescriptor_jarviscore_6f9bb307b6786fc3, []int{2}
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

// FileData - file data
type FileData struct {
	File                 []byte   `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
	Filename             string   `protobuf:"bytes,2,opt,name=filename,proto3" json:"filename,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FileData) Reset()         { *m = FileData{} }
func (m *FileData) String() string { return proto.CompactTextString(m) }
func (*FileData) ProtoMessage()    {}
func (*FileData) Descriptor() ([]byte, []int) {
	return fileDescriptor_jarviscore_6f9bb307b6786fc3, []int{3}
}
func (m *FileData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FileData.Unmarshal(m, b)
}
func (m *FileData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FileData.Marshal(b, m, deterministic)
}
func (dst *FileData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FileData.Merge(dst, src)
}
func (m *FileData) XXX_Size() int {
	return xxx_messageInfo_FileData.Size(m)
}
func (m *FileData) XXX_DiscardUnknown() {
	xxx_messageInfo_FileData.DiscardUnknown(m)
}

var xxx_messageInfo_FileData proto.InternalMessageInfo

func (m *FileData) GetFile() []byte {
	if m != nil {
		return m.File
	}
	return nil
}

func (m *FileData) GetFilename() string {
	if m != nil {
		return m.Filename
	}
	return ""
}

// CtrlScriptData - used in CtrlInfo.dat
type CtrlScriptData struct {
	File                 []byte   `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
	DestPath             string   `protobuf:"bytes,2,opt,name=destPath,proto3" json:"destPath,omitempty"`
	Filename             string   `protobuf:"bytes,3,opt,name=filename,proto3" json:"filename,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CtrlScriptData) Reset()         { *m = CtrlScriptData{} }
func (m *CtrlScriptData) String() string { return proto.CompactTextString(m) }
func (*CtrlScriptData) ProtoMessage()    {}
func (*CtrlScriptData) Descriptor() ([]byte, []int) {
	return fileDescriptor_jarviscore_6f9bb307b6786fc3, []int{4}
}
func (m *CtrlScriptData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CtrlScriptData.Unmarshal(m, b)
}
func (m *CtrlScriptData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CtrlScriptData.Marshal(b, m, deterministic)
}
func (dst *CtrlScriptData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CtrlScriptData.Merge(dst, src)
}
func (m *CtrlScriptData) XXX_Size() int {
	return xxx_messageInfo_CtrlScriptData.Size(m)
}
func (m *CtrlScriptData) XXX_DiscardUnknown() {
	xxx_messageInfo_CtrlScriptData.DiscardUnknown(m)
}

var xxx_messageInfo_CtrlScriptData proto.InternalMessageInfo

func (m *CtrlScriptData) GetFile() []byte {
	if m != nil {
		return m.File
	}
	return nil
}

func (m *CtrlScriptData) GetDestPath() string {
	if m != nil {
		return m.DestPath
	}
	return ""
}

func (m *CtrlScriptData) GetFilename() string {
	if m != nil {
		return m.Filename
	}
	return ""
}

// CtrlInfo -
type CtrlInfo struct {
	CtrlID               int64    `protobuf:"varint,1,opt,name=ctrlID,proto3" json:"ctrlID,omitempty"`
	CtrlType             string   `protobuf:"bytes,2,opt,name=ctrlType,proto3" json:"ctrlType,omitempty"`
	Command              string   `protobuf:"bytes,3,opt,name=command,proto3" json:"command,omitempty"`
	Params               []string `protobuf:"bytes,4,rep,name=params,proto3" json:"params,omitempty"`
	Description          string   `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
	Dat                  *any.Any `protobuf:"bytes,1000,opt,name=dat,proto3" json:"dat,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CtrlInfo) Reset()         { *m = CtrlInfo{} }
func (m *CtrlInfo) String() string { return proto.CompactTextString(m) }
func (*CtrlInfo) ProtoMessage()    {}
func (*CtrlInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_jarviscore_6f9bb307b6786fc3, []int{5}
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

func (m *CtrlInfo) GetCtrlID() int64 {
	if m != nil {
		return m.CtrlID
	}
	return 0
}

func (m *CtrlInfo) GetCtrlType() string {
	if m != nil {
		return m.CtrlType
	}
	return ""
}

func (m *CtrlInfo) GetCommand() string {
	if m != nil {
		return m.Command
	}
	return ""
}

func (m *CtrlInfo) GetParams() []string {
	if m != nil {
		return m.Params
	}
	return nil
}

func (m *CtrlInfo) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *CtrlInfo) GetDat() *any.Any {
	if m != nil {
		return m.Dat
	}
	return nil
}

type CtrlResult struct {
	CtrlID               int64    `protobuf:"varint,1,opt,name=ctrlID,proto3" json:"ctrlID,omitempty"`
	CtrlResult           string   `protobuf:"bytes,2,opt,name=ctrlResult,proto3" json:"ctrlResult,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CtrlResult) Reset()         { *m = CtrlResult{} }
func (m *CtrlResult) String() string { return proto.CompactTextString(m) }
func (*CtrlResult) ProtoMessage()    {}
func (*CtrlResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_jarviscore_6f9bb307b6786fc3, []int{6}
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

func (m *CtrlResult) GetCtrlID() int64 {
	if m != nil {
		return m.CtrlID
	}
	return 0
}

func (m *CtrlResult) GetCtrlResult() string {
	if m != nil {
		return m.CtrlResult
	}
	return ""
}

// JarvisMsg - jarvis base msg
//      sign(msgID + msgType + destAddr + curTime + srcAddr + data)
type JarvisMsg struct {
	MsgID     int64     `protobuf:"varint,1,opt,name=msgID,proto3" json:"msgID,omitempty"`
	CurTime   int64     `protobuf:"varint,2,opt,name=curTime,proto3" json:"curTime,omitempty"`
	SignR     []byte    `protobuf:"bytes,3,opt,name=signR,proto3" json:"signR,omitempty"`
	SignS     []byte    `protobuf:"bytes,4,opt,name=signS,proto3" json:"signS,omitempty"`
	PubKey    []byte    `protobuf:"bytes,5,opt,name=pubKey,proto3" json:"pubKey,omitempty"`
	SrcAddr   string    `protobuf:"bytes,6,opt,name=srcAddr,proto3" json:"srcAddr,omitempty"`
	MyAddr    string    `protobuf:"bytes,7,opt,name=myAddr,proto3" json:"myAddr,omitempty"`
	DestAddr  string    `protobuf:"bytes,8,opt,name=destAddr,proto3" json:"destAddr,omitempty"`
	MsgType   MSGTYPE   `protobuf:"varint,9,opt,name=msgType,proto3,enum=jarviscorepb.MSGTYPE" json:"msgType,omitempty"`
	ReplyType REPLYTYPE `protobuf:"varint,10,opt,name=replyType,proto3,enum=jarviscorepb.REPLYTYPE" json:"replyType,omitempty"`
	Err       string    `protobuf:"bytes,11,opt,name=err,proto3" json:"err,omitempty"`
	// Types that are valid to be assigned to Data:
	//	*JarvisMsg_NodeInfo
	//	*JarvisMsg_CtrlInfo
	//	*JarvisMsg_CtrlResult
	//	*JarvisMsg_ConnInfo
	//	*JarvisMsg_Msg
	//	*JarvisMsg_File
	Data                 isJarvisMsg_Data `protobuf_oneof:"data"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *JarvisMsg) Reset()         { *m = JarvisMsg{} }
func (m *JarvisMsg) String() string { return proto.CompactTextString(m) }
func (*JarvisMsg) ProtoMessage()    {}
func (*JarvisMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_jarviscore_6f9bb307b6786fc3, []int{7}
}
func (m *JarvisMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JarvisMsg.Unmarshal(m, b)
}
func (m *JarvisMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JarvisMsg.Marshal(b, m, deterministic)
}
func (dst *JarvisMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JarvisMsg.Merge(dst, src)
}
func (m *JarvisMsg) XXX_Size() int {
	return xxx_messageInfo_JarvisMsg.Size(m)
}
func (m *JarvisMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_JarvisMsg.DiscardUnknown(m)
}

var xxx_messageInfo_JarvisMsg proto.InternalMessageInfo

type isJarvisMsg_Data interface {
	isJarvisMsg_Data()
}

type JarvisMsg_NodeInfo struct {
	NodeInfo *NodeBaseInfo `protobuf:"bytes,100,opt,name=nodeInfo,proto3,oneof"`
}
type JarvisMsg_CtrlInfo struct {
	CtrlInfo *CtrlInfo `protobuf:"bytes,101,opt,name=ctrlInfo,proto3,oneof"`
}
type JarvisMsg_CtrlResult struct {
	CtrlResult *CtrlResult `protobuf:"bytes,102,opt,name=ctrlResult,proto3,oneof"`
}
type JarvisMsg_ConnInfo struct {
	ConnInfo *ConnectInfo `protobuf:"bytes,103,opt,name=connInfo,proto3,oneof"`
}
type JarvisMsg_Msg struct {
	Msg *JarvisMsg `protobuf:"bytes,104,opt,name=msg,proto3,oneof"`
}
type JarvisMsg_File struct {
	File *FileData `protobuf:"bytes,105,opt,name=file,proto3,oneof"`
}

func (*JarvisMsg_NodeInfo) isJarvisMsg_Data()   {}
func (*JarvisMsg_CtrlInfo) isJarvisMsg_Data()   {}
func (*JarvisMsg_CtrlResult) isJarvisMsg_Data() {}
func (*JarvisMsg_ConnInfo) isJarvisMsg_Data()   {}
func (*JarvisMsg_Msg) isJarvisMsg_Data()        {}
func (*JarvisMsg_File) isJarvisMsg_Data()       {}

func (m *JarvisMsg) GetData() isJarvisMsg_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *JarvisMsg) GetMsgID() int64 {
	if m != nil {
		return m.MsgID
	}
	return 0
}

func (m *JarvisMsg) GetCurTime() int64 {
	if m != nil {
		return m.CurTime
	}
	return 0
}

func (m *JarvisMsg) GetSignR() []byte {
	if m != nil {
		return m.SignR
	}
	return nil
}

func (m *JarvisMsg) GetSignS() []byte {
	if m != nil {
		return m.SignS
	}
	return nil
}

func (m *JarvisMsg) GetPubKey() []byte {
	if m != nil {
		return m.PubKey
	}
	return nil
}

func (m *JarvisMsg) GetSrcAddr() string {
	if m != nil {
		return m.SrcAddr
	}
	return ""
}

func (m *JarvisMsg) GetMyAddr() string {
	if m != nil {
		return m.MyAddr
	}
	return ""
}

func (m *JarvisMsg) GetDestAddr() string {
	if m != nil {
		return m.DestAddr
	}
	return ""
}

func (m *JarvisMsg) GetMsgType() MSGTYPE {
	if m != nil {
		return m.MsgType
	}
	return MSGTYPE_NULL_TYPE
}

func (m *JarvisMsg) GetReplyType() REPLYTYPE {
	if m != nil {
		return m.ReplyType
	}
	return REPLYTYPE_OK
}

func (m *JarvisMsg) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

func (m *JarvisMsg) GetNodeInfo() *NodeBaseInfo {
	if x, ok := m.GetData().(*JarvisMsg_NodeInfo); ok {
		return x.NodeInfo
	}
	return nil
}

func (m *JarvisMsg) GetCtrlInfo() *CtrlInfo {
	if x, ok := m.GetData().(*JarvisMsg_CtrlInfo); ok {
		return x.CtrlInfo
	}
	return nil
}

func (m *JarvisMsg) GetCtrlResult() *CtrlResult {
	if x, ok := m.GetData().(*JarvisMsg_CtrlResult); ok {
		return x.CtrlResult
	}
	return nil
}

func (m *JarvisMsg) GetConnInfo() *ConnectInfo {
	if x, ok := m.GetData().(*JarvisMsg_ConnInfo); ok {
		return x.ConnInfo
	}
	return nil
}

func (m *JarvisMsg) GetMsg() *JarvisMsg {
	if x, ok := m.GetData().(*JarvisMsg_Msg); ok {
		return x.Msg
	}
	return nil
}

func (m *JarvisMsg) GetFile() *FileData {
	if x, ok := m.GetData().(*JarvisMsg_File); ok {
		return x.File
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*JarvisMsg) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _JarvisMsg_OneofMarshaler, _JarvisMsg_OneofUnmarshaler, _JarvisMsg_OneofSizer, []interface{}{
		(*JarvisMsg_NodeInfo)(nil),
		(*JarvisMsg_CtrlInfo)(nil),
		(*JarvisMsg_CtrlResult)(nil),
		(*JarvisMsg_ConnInfo)(nil),
		(*JarvisMsg_Msg)(nil),
		(*JarvisMsg_File)(nil),
	}
}

func _JarvisMsg_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*JarvisMsg)
	// data
	switch x := m.Data.(type) {
	case *JarvisMsg_NodeInfo:
		b.EncodeVarint(100<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.NodeInfo); err != nil {
			return err
		}
	case *JarvisMsg_CtrlInfo:
		b.EncodeVarint(101<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.CtrlInfo); err != nil {
			return err
		}
	case *JarvisMsg_CtrlResult:
		b.EncodeVarint(102<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.CtrlResult); err != nil {
			return err
		}
	case *JarvisMsg_ConnInfo:
		b.EncodeVarint(103<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ConnInfo); err != nil {
			return err
		}
	case *JarvisMsg_Msg:
		b.EncodeVarint(104<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Msg); err != nil {
			return err
		}
	case *JarvisMsg_File:
		b.EncodeVarint(105<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.File); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("JarvisMsg.Data has unexpected type %T", x)
	}
	return nil
}

func _JarvisMsg_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*JarvisMsg)
	switch tag {
	case 100: // data.nodeInfo
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(NodeBaseInfo)
		err := b.DecodeMessage(msg)
		m.Data = &JarvisMsg_NodeInfo{msg}
		return true, err
	case 101: // data.ctrlInfo
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(CtrlInfo)
		err := b.DecodeMessage(msg)
		m.Data = &JarvisMsg_CtrlInfo{msg}
		return true, err
	case 102: // data.ctrlResult
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(CtrlResult)
		err := b.DecodeMessage(msg)
		m.Data = &JarvisMsg_CtrlResult{msg}
		return true, err
	case 103: // data.connInfo
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ConnectInfo)
		err := b.DecodeMessage(msg)
		m.Data = &JarvisMsg_ConnInfo{msg}
		return true, err
	case 104: // data.msg
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(JarvisMsg)
		err := b.DecodeMessage(msg)
		m.Data = &JarvisMsg_Msg{msg}
		return true, err
	case 105: // data.file
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(FileData)
		err := b.DecodeMessage(msg)
		m.Data = &JarvisMsg_File{msg}
		return true, err
	default:
		return false, nil
	}
}

func _JarvisMsg_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*JarvisMsg)
	// data
	switch x := m.Data.(type) {
	case *JarvisMsg_NodeInfo:
		s := proto.Size(x.NodeInfo)
		n += 2 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *JarvisMsg_CtrlInfo:
		s := proto.Size(x.CtrlInfo)
		n += 2 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *JarvisMsg_CtrlResult:
		s := proto.Size(x.CtrlResult)
		n += 2 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *JarvisMsg_ConnInfo:
		s := proto.Size(x.ConnInfo)
		n += 2 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *JarvisMsg_Msg:
		s := proto.Size(x.Msg)
		n += 2 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *JarvisMsg_File:
		s := proto.Size(x.File)
		n += 2 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*NodeBaseInfo)(nil), "jarviscorepb.NodeBaseInfo")
	proto.RegisterType((*ConnectInfo)(nil), "jarviscorepb.ConnectInfo")
	proto.RegisterType((*ReplyJoin)(nil), "jarviscorepb.ReplyJoin")
	proto.RegisterType((*FileData)(nil), "jarviscorepb.FileData")
	proto.RegisterType((*CtrlScriptData)(nil), "jarviscorepb.CtrlScriptData")
	proto.RegisterType((*CtrlInfo)(nil), "jarviscorepb.CtrlInfo")
	proto.RegisterType((*CtrlResult)(nil), "jarviscorepb.CtrlResult")
	proto.RegisterType((*JarvisMsg)(nil), "jarviscorepb.JarvisMsg")
	proto.RegisterEnum("jarviscorepb.MSGTYPE", MSGTYPE_name, MSGTYPE_value)
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
	ProcMsg(ctx context.Context, in *JarvisMsg, opts ...grpc.CallOption) (JarvisCoreServ_ProcMsgClient, error)
}

type jarvisCoreServClient struct {
	cc *grpc.ClientConn
}

func NewJarvisCoreServClient(cc *grpc.ClientConn) JarvisCoreServClient {
	return &jarvisCoreServClient{cc}
}

func (c *jarvisCoreServClient) ProcMsg(ctx context.Context, in *JarvisMsg, opts ...grpc.CallOption) (JarvisCoreServ_ProcMsgClient, error) {
	stream, err := c.cc.NewStream(ctx, &_JarvisCoreServ_serviceDesc.Streams[0], "/jarviscorepb.JarvisCoreServ/procMsg", opts...)
	if err != nil {
		return nil, err
	}
	x := &jarvisCoreServProcMsgClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type JarvisCoreServ_ProcMsgClient interface {
	Recv() (*JarvisMsg, error)
	grpc.ClientStream
}

type jarvisCoreServProcMsgClient struct {
	grpc.ClientStream
}

func (x *jarvisCoreServProcMsgClient) Recv() (*JarvisMsg, error) {
	m := new(JarvisMsg)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// JarvisCoreServServer is the server API for JarvisCoreServ service.
type JarvisCoreServServer interface {
	ProcMsg(*JarvisMsg, JarvisCoreServ_ProcMsgServer) error
}

func RegisterJarvisCoreServServer(s *grpc.Server, srv JarvisCoreServServer) {
	s.RegisterService(&_JarvisCoreServ_serviceDesc, srv)
}

func _JarvisCoreServ_ProcMsg_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(JarvisMsg)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(JarvisCoreServServer).ProcMsg(m, &jarvisCoreServProcMsgServer{stream})
}

type JarvisCoreServ_ProcMsgServer interface {
	Send(*JarvisMsg) error
	grpc.ServerStream
}

type jarvisCoreServProcMsgServer struct {
	grpc.ServerStream
}

func (x *jarvisCoreServProcMsgServer) Send(m *JarvisMsg) error {
	return x.ServerStream.SendMsg(m)
}

var _JarvisCoreServ_serviceDesc = grpc.ServiceDesc{
	ServiceName: "jarviscorepb.JarvisCoreServ",
	HandlerType: (*JarvisCoreServServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "procMsg",
			Handler:       _JarvisCoreServ_ProcMsg_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "jarviscore.proto",
}

func init() { proto.RegisterFile("jarviscore.proto", fileDescriptor_jarviscore_6f9bb307b6786fc3) }

var fileDescriptor_jarviscore_6f9bb307b6786fc3 = []byte{
	// 863 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x54, 0xdb, 0x8e, 0xe2, 0x46,
	0x10, 0x05, 0xcc, 0xcd, 0x05, 0x43, 0x7a, 0x3b, 0x7b, 0xf1, 0xce, 0x43, 0x34, 0xe2, 0x25, 0xa3,
	0x4d, 0xc4, 0x44, 0xb3, 0x89, 0x12, 0xed, 0x4b, 0xc4, 0x80, 0x59, 0xd8, 0x05, 0x3c, 0xdb, 0x36,
	0x8a, 0x56, 0x4a, 0x84, 0x3c, 0x76, 0x8f, 0xd7, 0x11, 0xb6, 0x51, 0xdb, 0x8c, 0xc4, 0x57, 0xe4,
	0x13, 0xf2, 0x21, 0xf9, 0xa1, 0x7c, 0x46, 0xd4, 0xe5, 0x0b, 0x66, 0xb2, 0xb3, 0xf3, 0x44, 0x9f,
	0xaa, 0x73, 0x0e, 0x55, 0xae, 0xea, 0x06, 0xf2, 0xa7, 0x2d, 0xee, 0xfc, 0xd8, 0x89, 0x04, 0x1f,
	0x6c, 0x45, 0x94, 0x44, 0xb4, 0x7b, 0x88, 0x6c, 0x6f, 0x4e, 0x5f, 0x7a, 0x51, 0xe4, 0x6d, 0xf8,
	0x05, 0xe6, 0x6e, 0x76, 0xb7, 0x17, 0x76, 0xb8, 0x4f, 0x89, 0x7d, 0x06, 0xdd, 0x65, 0xe4, 0xf2,
	0x2b, 0x3b, 0xe6, 0xb3, 0xf0, 0x36, 0xa2, 0xa7, 0xd0, 0x8e, 0xb9, 0xb8, 0x1b, 0xba, 0xae, 0xd0,
	0xaa, 0x67, 0xd5, 0x73, 0x95, 0x15, 0x98, 0x52, 0xa8, 0xdb, 0x32, 0x5e, 0xc3, 0x38, 0x9e, 0x65,
	0x2c, 0xb4, 0x03, 0xae, 0x29, 0x69, 0x4c, 0x9e, 0xfb, 0x7f, 0x40, 0x67, 0x14, 0x85, 0x21, 0x77,
	0x92, 0x47, 0x2d, 0x2f, 0xa1, 0x19, 0xec, 0x25, 0x0b, 0x4d, 0x3b, 0x97, 0xa7, 0x83, 0x72, 0xe1,
	0x83, 0x72, 0x69, 0x2c, 0x63, 0xf6, 0x5f, 0x83, 0xca, 0xf8, 0x76, 0xb3, 0x7f, 0x17, 0xf9, 0x61,
	0x51, 0x53, 0xf5, 0x33, 0x35, 0xd5, 0x4a, 0x35, 0xbd, 0x81, 0xf6, 0xc4, 0xdf, 0xf0, 0xb1, 0x9d,
	0xd8, 0x32, 0x7f, 0xeb, 0x6f, 0x38, 0x6a, 0xba, 0x0c, 0xcf, 0xb2, 0x48, 0xf9, 0x5b, 0xd2, 0x15,
	0xb8, 0xff, 0x3b, 0xf4, 0x46, 0x89, 0xd8, 0x98, 0x8e, 0xf0, 0xb7, 0xc9, 0x97, 0x1c, 0x5c, 0x1e,
	0x27, 0xd7, 0x76, 0xf2, 0x29, 0x77, 0xc8, 0xf1, 0x91, 0xbb, 0x72, 0xcf, 0xfd, 0x9f, 0x2a, 0xb4,
	0xa5, 0x3d, 0x7e, 0xab, 0xe7, 0xd0, 0x74, 0xe4, 0x79, 0x8c, 0xd6, 0x0a, 0xcb, 0x90, 0x34, 0x90,
	0x27, 0x6b, 0xbf, 0x2d, 0xca, 0xcb, 0x31, 0xd5, 0xa0, 0xe5, 0x44, 0x41, 0x60, 0x87, 0x6e, 0xe6,
	0x9d, 0x43, 0xe9, 0xb6, 0xb5, 0x85, 0x1d, 0xc4, 0x5a, 0xfd, 0x4c, 0x39, 0x57, 0x59, 0x86, 0xe8,
	0x19, 0x74, 0x5c, 0x1e, 0x63, 0x3b, 0x7e, 0x14, 0x6a, 0x0d, 0x54, 0x95, 0x43, 0xf4, 0x5b, 0x50,
	0x5c, 0x3b, 0xd1, 0xfe, 0x6d, 0xe1, 0x54, 0x9e, 0x0e, 0xd2, 0x05, 0x1a, 0xe4, 0x0b, 0x34, 0x18,
	0x86, 0x7b, 0x26, 0x19, 0xfd, 0x31, 0x80, 0x2c, 0x9e, 0xf1, 0x78, 0xb7, 0x49, 0x1e, 0x2c, 0xff,
	0x1b, 0x00, 0xa7, 0x60, 0x65, 0x0d, 0x94, 0x22, 0xfd, 0xbf, 0x1a, 0xa0, 0xbe, 0xc3, 0xc1, 0x2f,
	0x62, 0x8f, 0x3e, 0x85, 0x46, 0x10, 0x7b, 0x85, 0x49, 0x0a, 0xb0, 0xcd, 0x9d, 0xb0, 0xfc, 0x6c,
	0x40, 0x0a, 0xcb, 0xa1, 0xe4, 0xc7, 0xbe, 0x17, 0x32, 0x6c, 0xbf, 0xcb, 0x52, 0x90, 0x47, 0x4d,
	0xad, 0x7e, 0x88, 0x9a, 0xf8, 0x49, 0x76, 0x37, 0xef, 0xf9, 0x1e, 0xbb, 0xee, 0xb2, 0x0c, 0x49,
	0xf7, 0x58, 0x38, 0xb8, 0xa3, 0xcd, 0xf4, 0x23, 0x66, 0x50, 0x2a, 0x82, 0x3d, 0x26, 0x5a, 0x98,
	0xc8, 0x50, 0x3e, 0x6f, 0xcc, 0xb4, 0x0f, 0xf3, 0xc6, 0xdc, 0x05, 0xb4, 0x82, 0xd8, 0xc3, 0x69,
	0xa9, 0x67, 0xd5, 0xf3, 0xde, 0xe5, 0xb3, 0xe3, 0xbd, 0x5e, 0x98, 0x6f, 0xad, 0x8f, 0xd7, 0x3a,
	0xcb, 0x59, 0xf4, 0x27, 0x50, 0x85, 0xdc, 0x69, 0x94, 0x00, 0x4a, 0x5e, 0x1c, 0x4b, 0x98, 0x7e,
	0x3d, 0xff, 0x88, 0xa2, 0x03, 0x93, 0x12, 0x50, 0xb8, 0x10, 0x5a, 0x07, 0xff, 0x5e, 0x1e, 0xe9,
	0x2f, 0xd0, 0x0e, 0x23, 0x17, 0x2f, 0x8c, 0xe6, 0x3e, 0x76, 0xa5, 0xa6, 0x15, 0x56, 0xb0, 0xe9,
	0x8f, 0xe9, 0x8a, 0xa1, 0x92, 0xa3, 0xf2, 0xf9, 0xb1, 0x32, 0x5f, 0x52, 0xa9, 0xca, 0x99, 0xf4,
	0xcd, 0xd1, 0x64, 0x6f, 0x51, 0xa7, 0xfd, 0x5f, 0x97, 0xe6, 0xa7, 0x95, 0xf2, 0xd4, 0xe9, 0xcf,
	0xd0, 0x76, 0xa2, 0x30, 0xc4, 0x7f, 0xf4, 0x50, 0xf9, 0xf2, 0x9e, 0xf2, 0xf0, 0x8a, 0xe0, 0x9f,
	0x66, 0x64, 0xfa, 0x1d, 0x28, 0x41, 0xec, 0x69, 0x9f, 0x50, 0x73, 0xef, 0x3b, 0x15, 0x6b, 0x34,
	0xad, 0x30, 0xc9, 0xa2, 0xdf, 0x67, 0x77, 0xd5, 0xff, 0x5c, 0x4f, 0xf9, 0x9b, 0x30, 0xad, 0xa4,
	0xb7, 0xf8, 0xaa, 0x09, 0x75, 0xd7, 0x4e, 0xec, 0x57, 0x7f, 0xd7, 0xa0, 0x95, 0x4d, 0x89, 0x9e,
	0x80, 0xba, 0x5c, 0xcd, 0xe7, 0x6b, 0x09, 0x48, 0x85, 0x12, 0xe8, 0x8e, 0x8c, 0xe5, 0x52, 0x1f,
	0x59, 0xeb, 0xa5, 0x31, 0xd6, 0x49, 0x95, 0x3e, 0x81, 0x13, 0xa6, 0x7f, 0x58, 0xe9, 0x66, 0x1a,
	0x31, 0x49, 0x8d, 0x7e, 0x05, 0x9d, 0x89, 0xc1, 0x7e, 0x1b, 0xb2, 0xf1, 0x7a, 0x61, 0xbe, 0x25,
	0x8a, 0x54, 0xe5, 0x9c, 0x91, 0xc5, 0xe6, 0xa4, 0x4e, 0x9f, 0xc1, 0x13, 0x1c, 0x2a, 0xe2, 0x35,
	0xd3, 0xcd, 0xd5, 0xdc, 0x22, 0x0d, 0xda, 0x03, 0xb0, 0xd8, 0x2a, 0xb3, 0x22, 0x4d, 0x34, 0x5f,
	0xac, 0x4b, 0xa1, 0x16, 0x55, 0xa1, 0x81, 0x4a, 0xd2, 0xc6, 0xda, 0x8c, 0xb1, 0xbe, 0x9e, 0x2d,
	0x27, 0x06, 0x51, 0xd3, 0x4a, 0xd0, 0x33, 0xad, 0x90, 0x00, 0x7d, 0x01, 0x5f, 0xcf, 0x8d, 0xd1,
	0x70, 0x9e, 0x87, 0xd6, 0x86, 0x35, 0xd5, 0x19, 0xe9, 0x48, 0x6e, 0x9a, 0x30, 0xf5, 0xe5, 0x58,
	0x16, 0xd9, 0x3d, 0x70, 0x8f, 0xdb, 0x39, 0x91, 0x5c, 0x8b, 0x0d, 0x97, 0xe6, 0x44, 0x67, 0xeb,
	0xc9, 0x6c, 0xae, 0x93, 0xde, 0xab, 0x2b, 0x50, 0x8b, 0x9d, 0xa4, 0x4d, 0xa8, 0x19, 0xef, 0x49,
	0x85, 0x76, 0xa0, 0x95, 0xb5, 0x4d, 0xaa, 0x12, 0x58, 0xb3, 0x85, 0x6e, 0xac, 0x2c, 0x52, 0xa3,
	0x6d, 0xa8, 0xcf, 0xcc, 0x85, 0x4e, 0x14, 0x59, 0xbd, 0xce, 0x98, 0xc1, 0x48, 0xfd, 0xf2, 0x03,
	0xf4, 0xd2, 0x79, 0x8d, 0x22, 0xc1, 0x4d, 0x2e, 0xee, 0xe8, 0xaf, 0xd0, 0xda, 0x8a, 0xc8, 0x91,
	0xcf, 0xc0, 0x43, 0x83, 0x3d, 0x7d, 0x28, 0xd1, 0xaf, 0xfc, 0x50, 0xbd, 0x69, 0xe2, 0x1b, 0xf5,
	0xfa, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x84, 0x3b, 0x35, 0xae, 0x14, 0x07, 0x00, 0x00,
}
