package jarviscore

import (
	"context"

	"github.com/zhs007/jarviscore/coredb"
	coredbpb "github.com/zhs007/jarviscore/coredb/proto"
	pb "github.com/zhs007/jarviscore/proto"
)

// JarvisNode -
type JarvisNode interface {
	// Start - start jarvis node
	Start(ctx context.Context) (err error)
	// Stop - stop jarvis node
	Stop() (err error)
	// GetCoreDB - get jarvis node coredb
	GetCoreDB() *coredb.CoreDB

	// GetConfig - get config
	GetConfig() *Config

	// BuildStatus - build jarviscorepb.JarvisNodeStatus
	BuildStatus() *pb.JarvisNodeStatus

	// RequestCtrl - send ctrl to jarvisnode with addr
	RequestCtrl(ctx context.Context, addr string, ci *pb.CtrlInfo,
		funcOnResult FuncOnProcMsgResult) error
	// SendFile - send filedata to jarvisnode with addr
	SendFile(ctx context.Context, addr string, fd *pb.FileData,
		funcOnResult FuncOnProcMsgResult) error
	// SendFile2 - send filedata to jarvisnode with addr
	SendFile2(ctx context.Context, addr string, fd *pb.FileData,
		funcOnResult FuncOnProcMsgResult) error
	// RequestFile - request node send filedata to me
	RequestFile(ctx context.Context, addr string, rf *pb.RequestFile,
		funcOnResult FuncOnProcMsgResult) error
	// RequestNodes - request nodes
	RequestNodes(ctx context.Context, isNeedLocalHost bool, funcOnResult FuncOnGroupSendMsgResult) error
	// UpdateNode - update node
	UpdateNode(ctx context.Context, addr string, nodetype string, nodetypever string, isOnlyRestart bool,
		funcOnResult FuncOnProcMsgResult) error
	// UpdateAllNodes - update all nodes
	UpdateAllNodes(ctx context.Context, nodetype string, nodetypever string, isOnlyRestart bool,
		funcOnResult FuncOnGroupSendMsgResult) error
	// ClearLogs - clear logs
	ClearLogs(ctx context.Context, addr string, funcOnResult FuncOnProcMsgResult) error
	// ClearAllLogs - clear logs
	ClearAllLogs(ctx context.Context, funcOnResult FuncOnGroupSendMsgResult) error

	// AddNodeBaseInfo - add nodeinfo
	AddNodeBaseInfo(nbi *pb.NodeBaseInfo) error

	// OnMsg - proc JarvisMsg
	OnMsg(ctx context.Context, task *JarvisMsgTask) error

	// GetMyInfo - get my nodeinfo
	GetMyInfo() *BaseInfo

	// RegNodeEventFunc - reg event handle
	RegNodeEventFunc(event string, eventfunc FuncNodeEvent) error
	// RegMsgEventFunc - reg event handle
	RegMsgEventFunc(event string, eventfunc FuncMsgEvent) error

	// IsConnected - is connected this node
	IsConnected(addr string) bool

	// FindNodeWithName - find node with name
	FindNodeWithName(name string) *coredbpb.NodeInfo
	// FindNode - find node
	FindNode(addr string) *coredbpb.NodeInfo

	// SetNodeTypeInfo - set node type and version
	SetNodeTypeInfo(nodetype string, nodetypeversion string)

	// RegCtrl - register a ctrl
	RegCtrl(ctrltype string, ctrl Ctrl) error

	// PostMsg - like windows postMessage
	PostMsg(normal *NormalMsgTaskInfo, chanEnd chan int)
	// PostStreamMsg - like windows postMessage
	PostStreamMsg(stream *StreamMsgTaskInfo, chanEnd chan int)

	// SendStreamMsg - send stream message to other node
	SendStreamMsg(addr string, msgs []*pb.JarvisMsg, funcOnResult FuncOnProcMsgResult)
	// SendMsg - send a message to other node
	SendMsg(addr string, msg *pb.JarvisMsg, funcOnResult FuncOnProcMsgResult)

	// ConnectNode - connect node
	ConnectNode(node *coredbpb.NodeInfo, funcOnResult FuncOnProcMsgResult) error
	// ConnectNodeWithServAddr - connect node
	ConnectNodeWithServAddr(servaddr string, funcOnResult FuncOnProcMsgResult) error

	// OnClientProcMsg - on Client.ProcMsg
	OnClientProcMsg(addr string, msgid int64, onProcMsgResult FuncOnProcMsgResult) error
	// OnReplyProcMsg - on reply
	OnReplyProcMsg(ctx context.Context, addr string, replymsgid int64, jrt int, msg *pb.JarvisMsg, err error) error
}
