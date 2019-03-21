package jarviscore

import (
	"context"

	"github.com/zhs007/jarviscore/coredb"
	"github.com/zhs007/jarviscore/coredb/proto"
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

	// BuildStatus - build jarviscorepb.JarvisNodeStatus
	BuildStatus() *pb.JarvisNodeStatus

	// RequestCtrl - send ctrl to jarvisnode with addr
	RequestCtrl(ctx context.Context, addr string, ci *pb.CtrlInfo,
		funcOnResult FuncOnProcMsgResult) error
	// SendFile - send filedata to jarvisnode with addr
	SendFile(ctx context.Context, addr string, fd *pb.FileData,
		funcOnResult FuncOnProcMsgResult) error
	// RequestFile - request node send filedata to me
	RequestFile(ctx context.Context, addr string, rf *pb.RequestFile,
		funcOnResult FuncOnProcMsgResult) error
	// RequestNodes - request nodes
	RequestNodes(ctx context.Context, funcOnResult FuncOnGroupSendMsgResult) error
	// UpdateNode - update node
	UpdateNode(ctx context.Context, addr string, nodetype string, nodetypever string,
		funcOnResult FuncOnProcMsgResult) error
	// UpdateAllNodes - update all nodes
	UpdateAllNodes(ctx context.Context, nodetype string, nodetypever string,
		funcOnResult FuncOnGroupSendMsgResult) error

	// AddNodeBaseInfo - add nodeinfo
	AddNodeBaseInfo(nbi *pb.NodeBaseInfo) error

	// OnMsg - proc JarvisMsg
	OnMsg(ctx context.Context, task JarvisTask) error

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
	PostMsg(msg *pb.JarvisMsg, stream pb.JarvisCoreServ_ProcMsgServer, chanEnd chan int,
		funcOnResult FuncOnProcMsgResult)

	// ConnectNode - connect node
	ConnectNode(node *coredbpb.NodeInfo, funcOnResult FuncOnProcMsgResult) error
	// ConnectNodeWithServAddr - connect node
	ConnectNodeWithServAddr(servaddr string, funcOnResult FuncOnProcMsgResult) error
}
