package jarviscore

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/zhs007/jarviscore/base"
	"github.com/zhs007/jarviscore/basedef"
	"github.com/zhs007/jarviscore/coredb"
	"github.com/zhs007/jarviscore/coredb/proto"
	pb "github.com/zhs007/jarviscore/proto"
	"go.uber.org/zap"
)

// JarvisNode -
type JarvisNode interface {
	// Start - start jarvis node
	Start(ctx context.Context) (err error)
	// Stop - stop jarvis node
	Stop() (err error)
	// GetCoreDB - get jarvis node coredb
	GetCoreDB() *coredb.CoreDB

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
	OnMsg(ctx context.Context, msg *pb.JarvisMsg, stream pb.JarvisCoreServ_ProcMsgServer, funcOnResult FuncOnProcMsgResult) error

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

// jarvisNode -
type jarvisNode struct {
	myinfo       BaseInfo
	coredb       *coredb.CoreDB
	mgrJasvisMsg *jarvisMsgMgr
	mgrClient2   *jarvisClient2
	serv2        *jarvisServer2
	mgrEvent     *eventMgr
	cfg          *Config
	mgrCtrl      *ctrlMgr
	// mgrRequest   *requestMgr
}

const (
	nodeinfoCacheSize       = 32
	randomMax         int64 = 0x7fffffffffffffff
	stateNormal             = 0
	stateStart              = 1
	stateEnd                = 2
)

// NewNode -
func NewNode(cfg *Config) (JarvisNode, error) {
	jarvisbase.Info("jarviscore version is " + basedef.VERSION)

	if !IsValidNodeName(cfg.BaseNodeInfo.NodeName) {
		// jarvisbase.Error("NewNode:IsValidNodeName", zap.Error(ErrInvalidNodeName))

		return nil, ErrInvalidNodeName
	}

	db, err := coredb.NewCoreDB(cfg.AnkaDB.DBPath, cfg.AnkaDB.HTTPServ, cfg.AnkaDB.Engine)
	if err != nil {
		// jarvisbase.Error("NewNode:newCoreDB", zap.Error(err))

		return nil, err
	}

	node := &jarvisNode{
		myinfo: BaseInfo{
			Name:        cfg.BaseNodeInfo.NodeName,
			BindAddr:    cfg.BaseNodeInfo.BindAddr,
			ServAddr:    cfg.BaseNodeInfo.ServAddr,
			CoreVersion: basedef.VERSION,
		},
		coredb: db,
		cfg:    cfg,
		mgrCtrl: &ctrlMgr{
			mapCtrl: make(map[string](Ctrl)),
		},
		// mgrRequest: &requestMgr{},
	}

	node.mgrCtrl.Reg(CtrlTypeShell, &CtrlShell{})
	node.mgrCtrl.Reg(CtrlTypeScriptFile, &CtrlScriptFile{})
	node.mgrCtrl.Reg(CtrlTypeScriptFile2, &CtrlScriptFile2{})

	// event
	node.mgrEvent = newEventMgr(node)
	node.mgrEvent.regNodeEventFunc(EventOnNodeConnected, onNodeConnected)
	node.mgrEvent.regNodeEventFunc(EventOnIConnectNode, onIConnectNode)
	node.mgrEvent.regNodeEventFunc(EventOnDeprecateNode, onDeprecateNode)
	node.mgrEvent.regNodeEventFunc(EventOnIConnectNodeFail, onIConnectNodeFail)

	err = node.coredb.Init()
	if err != nil {
		// jarvisbase.Error("NewNode:Init", zap.Error(err))

		return nil, err
	}

	node.myinfo.Addr = node.coredb.GetPrivateKey().ToAddress()
	node.myinfo.Name = cfg.BaseNodeInfo.NodeName
	node.myinfo.BindAddr = cfg.BaseNodeInfo.BindAddr
	node.myinfo.ServAddr = cfg.BaseNodeInfo.ServAddr

	jarvisbase.Info("jarviscore.NewNode",
		zap.String("Addr", node.myinfo.Addr),
		zap.String("Name", node.myinfo.Name),
		zap.String("BindAddr", node.myinfo.BindAddr),
		zap.String("ServAddr", node.myinfo.ServAddr))

	// mgrJasvisMsg
	node.mgrJasvisMsg = newJarvisMsgMgr(node)

	// mgrClient2
	node.mgrClient2 = newClient2(node)

	return node, nil
}

// Stop -
func (n *jarvisNode) Stop() error {
	if n.serv2 != nil {
		n.serv2.Stop()
	}

	return nil
}

// Start -
func (n *jarvisNode) Start(ctx context.Context) (err error) {
	coredbctx, coredbcancel := context.WithCancel(ctx)
	defer coredbcancel()
	go n.coredb.Start(coredbctx)

	msgmgrctx, msgmgrcancel := context.WithCancel(ctx)
	defer msgmgrcancel()
	go n.mgrJasvisMsg.start(msgmgrctx)

	clientctx, clientcancel := context.WithCancel(ctx)
	defer clientcancel()
	go n.mgrClient2.start(clientctx)

	jarvisbase.Info("StartServer", zap.String("ServAddr", n.myinfo.ServAddr))
	n.serv2, err = newServer2(n)
	if err != nil {
		return err
	}

	servctx, servcancel := context.WithCancel(ctx)
	defer servcancel()
	go n.serv2.Start(servctx)

	n.ConnectNodeWithServAddr(n.cfg.RootServAddr, nil)
	jarvisbase.Info("StartServer:connectRoot",
		zap.String("RootServAddr", n.cfg.RootServAddr))

	n.connectAllNodes()
	jarvisbase.Info("StartServer:connectAllNodes")

	tickerRequestChild := time.NewTicker(time.Duration(n.cfg.TimeRequestChild) * time.Second)

	for {
		select {
		case <-tickerRequestChild.C:
			n.onTimerRequestNodes(ctx)
		case <-ctx.Done():
			n.Stop()
			return nil
		}
	}
}

// GetCoreDB - get coredb
func (n *jarvisNode) GetCoreDB() *coredb.CoreDB {
	return n.coredb
}

// OnMsg - proc JarvisMsg
func (n *jarvisNode) OnMsg(ctx context.Context, msg *pb.JarvisMsg, stream pb.JarvisCoreServ_ProcMsgServer, funcOnResult FuncOnProcMsgResult) error {
	jarvisbase.Debug("jarvisNode.OnMsg", jarvisbase.JSON("msg", msg))

	// is timeout
	if IsTimeOut(msg) {
		jarvisbase.Warn("jarvisNode.OnMsg", zap.Error(ErrJarvisMsgTimeOut))

		n.replyStream2(msg, stream, pb.REPLYTYPE_ERROR, ErrJarvisMsgTimeOut.Error())

		return nil
	}

	// // proc local msg
	// if msg.MsgType == pb.MSGTYPE_LOCAL_SENDMSG {

	// 	// verify msg
	// 	err := VerifyJarvisMsg(msg)
	// 	if err != nil {
	// 		jarvisbase.Warn("jarvisNode.OnMsg",
	// 			zap.Error(err),
	// 			jarvisbase.JSON("msg", msg))

	// 		n.replyStream2(msg, stream, pb.REPLYTYPE_ERROR, err.Error())

	// 		return nil
	// 	}

	// 	if msg.MsgType == pb.MSGTYPE_LOCAL_SENDMSG {
	// 		return n.onMsgLocalSendMsg(ctx, msg, funcOnResult)
	// 	}
	// }

	// proc connect msg
	if msg.MsgType == pb.MSGTYPE_CONNECT_NODE {
		// verify msg
		err := VerifyJarvisMsg(msg)
		if err != nil {
			jarvisbase.Warn("jarvisNode.OnMsg",
				zap.Error(err),
				jarvisbase.JSON("msg", msg))

			n.replyStream2(msg, stream, pb.REPLYTYPE_ERROR, err.Error())

			return nil
		}

		return n.onMsgConnectNode(ctx, msg, stream)
	}

	// if is not my msg, broadcast msg
	if n.myinfo.Addr != msg.DestAddr {
		//!!! 先不考虑转发协议
		// n.mgrClient2.addTask(msg, "", nil, nil)
	} else {
		// verify msg
		err := VerifyJarvisMsg(msg)
		if err != nil {
			jarvisbase.Warn("jarvisNode.OnMsg:VerifyJarvisMsg",
				zap.Error(err),
				jarvisbase.JSON("msg", msg))

			n.replyStream2(msg, stream, pb.REPLYTYPE_ERROR, err.Error())

			return nil
		}

		err = n.checkMsgID(ctx, msg)
		if err != nil {
			jarvisbase.Warn("jarvisNode.OnMsg:checkMsgID", zap.Error(err))

			if err == ErrInvalidMsgID {
				n.replyStream2(msg, stream, pb.REPLYTYPE_ERRMSGID, "")
			} else {
				n.replyStream2(msg, stream, pb.REPLYTYPE_ERROR, err.Error())
			}

			return nil
		}

		if msg.MsgType == pb.MSGTYPE_NODE_INFO {
			return n.onMsgNodeInfo(ctx, msg)
		} else if msg.MsgType == pb.MSGTYPE_REPLY_CONNECT {
			return n.onMsgReplyConnect(ctx, msg)
		} else if msg.MsgType == pb.MSGTYPE_REQUEST_CTRL {
			return n.onMsgRequestCtrl(ctx, msg, stream, funcOnResult)
		} else if msg.MsgType == pb.MSGTYPE_REPLY_CTRL_RESULT {
			return n.onMsgCtrlResult(ctx, msg)
		} else if msg.MsgType == pb.MSGTYPE_REQUEST_NODES {
			return n.onMsgRequestNodes(ctx, msg, stream)
		} else if msg.MsgType == pb.MSGTYPE_TRANSFER_FILE {
			return n.onMsgTransferFile(ctx, msg, stream)
		} else if msg.MsgType == pb.MSGTYPE_REQUEST_FILE {
			return n.onMsgRequestFile(ctx, msg, stream)
		} else if msg.MsgType == pb.MSGTYPE_REPLY_REQUEST_FILE {
			return n.onMsgReplyRequestFile(ctx, msg)
		} else if msg.MsgType == pb.MSGTYPE_REPLY_TRANSFER_FILE {
			return n.onMsgReplyTransferFile(ctx, msg)
		} else if msg.MsgType == pb.MSGTYPE_REPLY2 {
			return n.onMsgReply2(ctx, msg)
		} else if msg.MsgType == pb.MSGTYPE_UPDATENODE {
			return n.onMsgUpdateNode(ctx, msg, stream)
		}

	}

	return nil
}

// // onMsgLocalConnect
// func (n *jarvisNode) onMsgLocalConnect(ctx context.Context, msg *pb.JarvisMsg, funcOnResult FuncOnProcMsgResult) error {
// 	ci := msg.GetConnInfo()

// 	// if is me, return
// 	if ci.ServAddr == n.myinfo.ServAddr {
// 		return nil
// 	}

// 	cn := n.coredb.FindNodeWithServAddr(ci.ServAddr)
// 	if cn == nil {
// 		n.mgrClient2.addConnTask(ci.ServAddr, nil, funcOnResult)

// 		return nil
// 	}

// 	// if is me, return
// 	if cn.Addr == n.myinfo.Addr {
// 		return nil
// 	}

// 	// if it is deprecated, return
// 	if cn.Deprecated {
// 		return nil
// 	}

// 	if cn.ConnType == coredbpb.CONNECTTYPE_UNKNOWN_CONN {
// 		n.mgrClient2.addConnTask(cn.ServAddr, cn, funcOnResult)

// 		return nil
// 	}

// 	return nil
// }

// onMsgNodeInfo
func (n *jarvisNode) onMsgNodeInfo(ctx context.Context, msg *pb.JarvisMsg) error {
	ni := msg.GetNodeInfo()
	return n.AddNodeBaseInfo(ni)
}

// onMsgConnectNode
func (n *jarvisNode) onMsgConnectNode(ctx context.Context, msg *pb.JarvisMsg, stream pb.JarvisCoreServ_ProcMsgServer) error {
	if stream == nil {
		jarvisbase.Warn("jarvisNode.onMsgConnectNode", zap.Error(ErrStreamNil))

		return ErrStreamNil
	}

	ci := msg.GetConnInfo()

	mni := &pb.NodeBaseInfo{
		ServAddr:        n.myinfo.ServAddr,
		Addr:            n.myinfo.Addr,
		Name:            n.myinfo.Name,
		NodeTypeVersion: n.myinfo.NodeTypeVersion,
		NodeType:        n.myinfo.NodeType,
		CoreVersion:     n.myinfo.CoreVersion,
	}

	sendmsg, err := BuildReplyConn(n, n.myinfo.Addr, ci.MyInfo.Addr, mni)
	if err != nil {
		jarvisbase.Warn("jarvisNode.onMsgConnectNode:BuildReplyConn", zap.Error(err))

		return err
	}

	err = n.sendMsg2ClientStream(stream, sendmsg)
	if err != nil {
		jarvisbase.Warn("jarvisNode.onMsgConnectNode:sendMsg2ClientStream", zap.Error(err))

		return err
	}

	cn := n.coredb.GetNode(ci.MyInfo.Addr)
	if cn == nil {
		err := n.coredb.UpdNodeBaseInfo(ci.MyInfo)
		if err != nil {
			jarvisbase.Warn("jarvisNode.onMsgConnectNode:UpdNodeBaseInfo", zap.Error(err))

			return err
		}

		cn = n.coredb.GetNode(ci.MyInfo.Addr)
	} else {
		err := n.coredb.UpdNodeBaseInfo(ci.MyInfo)
		if err != nil {
			jarvisbase.Warn("jarvisNode.onMsgConnectNode:UpdNodeBaseInfo", zap.Error(err))

			return err
		}
	}

	n.mgrEvent.onNodeEvent(ctx, EventOnNodeConnected, cn)

	jarvisbase.Debug("jarvisNode.onMsgConnectNode:end")

	return nil
}

// onMsgReplyConnect
func (n *jarvisNode) onMsgReplyConnect(ctx context.Context, msg *pb.JarvisMsg) error {
	ni := msg.GetNodeInfo()
	cn := n.coredb.GetNode(ni.Addr)
	if cn == nil {
		err := n.coredb.UpdNodeBaseInfo(ni)
		if err != nil {
			jarvisbase.Warn("jarvisNode.onMsgReplyConnect:InsNode", zap.Error(err))

			return err
		}

		cn = n.coredb.GetNode(ni.Addr)
	} else {
		n.coredb.UpdNodeBaseInfo(ni)
	}

	if msg.LastMsgID > 0 {
		cn.LastSendMsgID = msg.LastMsgID
	}

	n.mgrEvent.onNodeEvent(ctx, EventOnIConnectNode, cn)

	return nil
}

// GetMyInfo - get my nodeinfo
func (n *jarvisNode) GetMyInfo() *BaseInfo {
	return &n.myinfo
}

// connectAllNodes - connect all nodes
func (n *jarvisNode) connectAllNodes() error {
	n.coredb.ForEachMapNodes(func(key string, node *coredbpb.NodeInfo) error {
		n.ConnectNode(node, nil)
		return nil
	})

	return nil
}

// ConnectNodeWithServAddr - connect node
func (n *jarvisNode) ConnectNodeWithServAddr(servaddr string, funcOnResult FuncOnProcMsgResult) error {
	// if is me, return
	if servaddr == n.myinfo.ServAddr {
		return nil
	}

	cn := n.coredb.FindNodeWithServAddr(servaddr)
	if cn == nil {
		n.mgrClient2.addConnTask(servaddr, nil, funcOnResult)

		return nil
	}

	// if is me, return
	if cn.Addr == n.myinfo.Addr {
		return nil
	}

	// if it is deprecated, return
	if coredb.IsDeprecatedNode(cn) {
		return nil
	}

	if cn.ConnType == coredbpb.CONNECTTYPE_UNKNOWN_CONN {
		cn.ConnectNums++
		cn.LastConnectTime = time.Now().Unix()

		n.mgrClient2.addConnTask(cn.ServAddr, cn, funcOnResult)

		return nil
	}

	// nbi := &pb.NodeBaseInfo{
	// 	ServAddr:        n.myinfo.ServAddr,
	// 	Addr:            n.myinfo.Addr,
	// 	Name:            n.myinfo.Name,
	// 	NodeTypeVersion: n.myinfo.NodeTypeVersion,
	// 	NodeType:        n.myinfo.NodeType,
	// 	CoreVersion:     n.myinfo.CoreVersion,
	// }

	// msg, err := BuildLocalConnectOther(n, n.myinfo.Addr, "", servaddr, nbi)
	// if err != nil {
	// 	jarvisbase.Warn("jarvisNode.ConnectNodeWithServAddr:BuildLocalConnectOther", zap.Error(err))

	// 	return err
	// }

	// n.PostMsg(msg, nil, nil, funcOnResult)

	return nil
}

// ConnectNode - connect node
func (n *jarvisNode) ConnectNode(node *coredbpb.NodeInfo, funcOnResult FuncOnProcMsgResult) error {
	// if is me, return
	if node.Addr == n.myinfo.Addr {
		return nil
	}

	// if it is deprecated, return
	if coredb.IsDeprecatedNode(node) {
		return nil
	}

	if node.ConnType == coredbpb.CONNECTTYPE_UNKNOWN_CONN {
		node.ConnectNums++
		node.LastConnectTime = time.Now().Unix()

		n.mgrClient2.addConnTask(node.ServAddr, node, funcOnResult)

		return nil
	}

	// nbi := &pb.NodeBaseInfo{
	// 	ServAddr:        n.myinfo.ServAddr,
	// 	Addr:            n.myinfo.Addr,
	// 	Name:            n.myinfo.Name,
	// 	NodeTypeVersion: n.myinfo.NodeTypeVersion,
	// 	NodeType:        n.myinfo.NodeType,
	// 	CoreVersion:     n.myinfo.CoreVersion,
	// }

	// msg, err := BuildLocalConnectOther(n, n.myinfo.Addr, node.Addr, node.ServAddr, nbi)
	// if err != nil {
	// 	jarvisbase.Warn("jarvisNode.ConnectNode:BuildLocalConnectOther", zap.Error(err))

	// 	return err
	// }

	// n.PostMsg(msg, nil, nil, funcOnResult)

	return nil
}

// onMsgRequestCtrl
func (n *jarvisNode) onMsgRequestCtrl(ctx context.Context, msg *pb.JarvisMsg,
	stream pb.JarvisCoreServ_ProcMsgServer, funcOnResult FuncOnProcMsgResult) error {

	jarvisbase.Info("jarvisNode.onMsgRequestCtrl:recvmsg",
		jarvisbase.JSON("msg", msg))

	n.replyStream2(msg, stream, pb.REPLYTYPE_ISME, "")

	n.mgrEvent.onMsgEvent(ctx, EventOnCtrl, msg)

	ci := msg.GetCtrlInfo()
	ret, err := n.mgrCtrl.Run(ci)
	if err != nil {
		sendmsg2, err := BuildCtrlResult(n, n.myinfo.Addr, msg.SrcAddr, ci.CtrlID, msg.MsgID, err.Error())
		if err != nil {
			jarvisbase.Warn("jarvisNode.onMsgRequestCtrl:BuildCtrlResult", zap.Error(err))

			return err
		}

		jarvisbase.Info("jarvisNode.onMsgRequestCtrl",
			jarvisbase.JSON("msg", msg),
			jarvisbase.JSON("sendmsg", sendmsg2))

		n.mgrClient2.addSendMsgTask(sendmsg2, nil, funcOnResult)

		return nil
	}

	sendmsg2, err := BuildCtrlResult(n, n.myinfo.Addr, msg.SrcAddr, ci.CtrlID, msg.MsgID, string(ret))

	jarvisbase.Info("jarvisNode.onMsgRequestCtrl",
		jarvisbase.JSON("msg", msg),
		jarvisbase.JSON("sendmsg", sendmsg2))

	n.mgrClient2.addSendMsgTask(sendmsg2, nil, funcOnResult)

	return nil
}

// onMsgReply2
func (n *jarvisNode) onMsgReply2(ctx context.Context, msg *pb.JarvisMsg) error {
	// if msg.ReplyMsgID > 0 {
	// 	n.mgrRequest.onReplyRequest(ctx, n, msg)
	// }

	if msg.ReplyType == pb.REPLYTYPE_ERRMSGID {
	}

	if msg.LastMsgID > 0 {
		cn := n.coredb.GetNode(msg.SrcAddr)
		if cn != nil && cn.LastSendMsgID != msg.LastMsgID {
			cn.LastSendMsgID = msg.LastMsgID

			n.coredb.UpdNodeInfo(msg.SrcAddr)
		}
	}

	return nil
}

// onMsgReplyTransferFile
func (n *jarvisNode) onMsgReplyTransferFile(ctx context.Context, msg *pb.JarvisMsg) error {
	// if msg.ReplyMsgID > 0 {
	// 	n.mgrRequest.onReplyRequest(ctx, n, msg)
	// }

	n.mgrEvent.onMsgEvent(ctx, EventOnReplyTransferFile, msg)

	return nil
}

// onMsgCtrlResult
func (n *jarvisNode) onMsgCtrlResult(ctx context.Context, msg *pb.JarvisMsg) error {
	n.mgrEvent.onMsgEvent(ctx, EventOnCtrlResult, msg)

	return nil
}

// onMsgLocalSendMsg
func (n *jarvisNode) onMsgLocalSendMsg(ctx context.Context, msg *pb.JarvisMsg,
	funcOnResult FuncOnProcMsgResult) error {

	sendmsg := msg.GetMsg()

	n.mgrClient2.addSendMsgTask(sendmsg, nil, funcOnResult)

	return nil
}

// onMsgUpdateNode
func (n *jarvisNode) onMsgUpdateNode(ctx context.Context, msg *pb.JarvisMsg, stream pb.JarvisCoreServ_ProcMsgServer) error {

	if n.cfg.AutoUpdate {
		n.replyStream2(msg, stream, pb.REPLYTYPE_ISME, "")

		n.mgrEvent.onMsgEvent(ctx, EventOnUpdateNode, msg)

		curscript, outstring, err := updateNode(&UpdateNodeParam{
			NewVersion: "v" + msg.GetUpdateNode().NodeTypeVersion,
		}, n.cfg.UpdateScript)
		if err != nil {
			n.replyStream2(msg, stream, pb.REPLYTYPE_ERROR, err.Error())

			return err
		}

		n.replyStream2(msg, stream, pb.REPLYTYPE_OK, curscript)
		n.replyStream2(msg, stream, pb.REPLYTYPE_OK, outstring)
	} else {
		n.replyStream2(msg, stream, pb.REPLYTYPE_ERROR, ErrAutoUpdateClosed.Error())
	}

	return nil
}

// onNodeConnected - func event
func onNodeConnected(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {

	node.ConnectMe = true
	node.NumsConnectFail = 0
	node.TimestampDeprecated = 0

	jarvisbase.Debug("jarvisNode.onMsgConnectNode:ConnType",
		zap.Int32("ConnType", int32(node.ConnType)))

	if node.ConnType == coredbpb.CONNECTTYPE_UNKNOWN_CONN {
		err := jarvisnode.ConnectNode(node, nil)
		if err != nil {
			jarvisbase.Warn("jarvisNode.onNodeConnected:ConnectNode",
				zap.Error(err))

			return err
		}

		err = jarvisnode.GetCoreDB().UpdNodeInfo(node.Addr)
		if err != nil {
			jarvisbase.Warn("jarvisNode.onNodeConnected:UpdNodeInfo",
				zap.Error(err))

			return err
		}
	} else {
		err := jarvisnode.GetCoreDB().UpdNodeInfo(node.Addr)
		if err != nil {
			jarvisbase.Warn("jarvisNode.onNodeConnected:UpdNodeInfo", zap.Error(err))

			return err
		}
	}

	return nil
}

// onIConnectNode - func event
func onIConnectNode(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
	jarvisbase.Debug("onIConnectNode")

	node.ConnectedNums++
	node.LastConnectedTime = time.Now().Unix()
	node.ConnType = coredbpb.CONNECTTYPE_DIRECT_CONN

	node.TimestampDeprecated = 0
	node.NumsConnectFail = 0

	err := jarvisnode.GetCoreDB().UpdNodeInfo(node.Addr)
	if err != nil {
		jarvisbase.Warn("jarvisNode.onIConnectNode:UpdNodeInfo", zap.Error(err))
	}

	return nil
}

// onDeprecateNode - func event
func onDeprecateNode(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
	jarvisbase.Debug("onDeprecateNode")

	if !node.Deprecated {
		jarvisbase.Info("onDeprecateNode",
			zap.String("addr", node.Addr),
			zap.String("servaddr", node.ServAddr))

		node.Deprecated = true

		err := jarvisnode.GetCoreDB().UpdNodeInfo(node.Addr)
		if err != nil {
			jarvisbase.Warn("jarvisNode.onDeprecateNode:UpdNodeInfo",
				zap.Error(err))
		}
	}

	return nil
}

// onIConnectNodeFail - func event
func onIConnectNodeFail(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
	jarvisbase.Debug("onIConnectNodeFail")

	if !node.Deprecated {
		node.NumsConnectFail++

		if node.NumsConnectFail%3 == 0 {
			ci := int(node.NumsConnectFail / 3)
			if ci >= len(basedef.LastTimeDeprecated) {
				ci = len(basedef.LastTimeDeprecated) - 1
			}

			node.TimestampDeprecated = time.Now().Unix() + basedef.LastTimeDeprecated[ci]

			jarvisbase.Info("onIConnectNodeFail",
				zap.String("addr", node.Addr),
				zap.String("servaddr", node.ServAddr),
				zap.Int32("NumsConnectFail", node.NumsConnectFail),
				zap.Int64("TimestampDeprecated", node.TimestampDeprecated))
		}

		err := jarvisnode.GetCoreDB().UpdNodeInfo(node.Addr)
		if err != nil {
			jarvisbase.Warn("jarvisNode.onIConnectNodeFail:UpdNodeInfo",
				zap.Error(err))
		}
	}

	return nil
}

// RegNodeEventFunc - reg event handle
func (n *jarvisNode) RegNodeEventFunc(event string, eventfunc FuncNodeEvent) error {
	return n.mgrEvent.regNodeEventFunc(event, eventfunc)
}

// RegMsgEventFunc - reg event handle
func (n *jarvisNode) RegMsgEventFunc(event string, eventfunc FuncMsgEvent) error {
	return n.mgrEvent.regMsgEventFunc(event, eventfunc)
}

// IsConnected - is connected this node
func (n *jarvisNode) IsConnected(addr string) bool {
	return n.mgrClient2.isConnected(addr)
}

// onMsgLocalRequesrNodes
func (n *jarvisNode) onMsgLocalRequesrNodes(ctx context.Context, msg *pb.JarvisMsg,
	funcOnResult FuncOnProcMsgResult) error {

	// numsSend := 0
	// numsRecv := 0
	// var totalResults []*ResultSendMsg

	// //!! 在网络IO很快的时候，假设一共有2个节点，但第一个节点很快返回的话，可能还没全部发送完成，就产生回调
	// //!! 所以这里分2次遍历
	// n.coredb.ForEachMapNodes(func(key string, v *coredbpb.NodeInfo) error {
	// 	if !v.Deprecated && n.mgrClient2.isConnected(v.Addr) {
	// 		numsSend++
	// 	}

	// 	return nil
	// })

	// n.coredb.ForEachMapNodes(func(key string, v *coredbpb.NodeInfo) error {
	// 	jarvisbase.Debug(fmt.Sprintf("jarvisNode.onMsgLocalRequesrNodes %v", v))

	// 	if !v.Deprecated && n.mgrClient2.isConnected(v.Addr) {
	// 		sendmsg, err := BuildRequestNodes(n, n.myinfo.Addr, v.Addr)
	// 		if err != nil {
	// 			jarvisbase.Warn("jarvisNode.onMsgLocalRequesrNodes:BuildRequestNodes", zap.Error(err))

	// 			return nil
	// 		}

	// 		n.mgrEvent.onNodeEvent(ctx, EventOnRequestNode, v)

	// 		n.mgrClient2.addTask(sendmsg, "", nil,
	// 			func(ctx context.Context, jarvisnode JarvisNode, lstResult []*ResultSendMsg) error {
	// 				numsRecv++

	// 				if len(lstResult) != 1 {
	// 					jarvisbase.Error("jarvisNode.onMsgLocalRequesrNodes:FuncOnSendMsgResult", zap.Int("len", len(lstResult)))

	// 					totalResults = append(totalResults,
	// 						&ResultSendMsg{
	// 							Err: ErrFuncOnSendMsgResultLength,
	// 						})
	// 				} else {
	// 					totalResults = append(totalResults, lstResult[0])
	// 				}

	// 				if funcOnResult != nil && numsSend == numsRecv {
	// 					funcOnResult(ctx, jarvisnode, totalResults)
	// 				}

	// 				return nil
	// 			})
	// 	}

	// 	return nil
	// })

	return nil
}

// RequestCtrl - send ctrl to jarvisnode with addr
func (n *jarvisNode) RequestCtrl(ctx context.Context, addr string, ci *pb.CtrlInfo,
	funcOnResult FuncOnProcMsgResult) error {

	sendmsg, err := BuildRequestCtrl(n, n.myinfo.Addr, addr, ci)
	if err != nil {
		jarvisbase.Warn("jarvisNode.RequestCtrl", zap.Error(err))

		return err
	}

	n.mgrClient2.addSendMsgTask(sendmsg, nil, funcOnResult)

	// msg, err := BuildLocalSendMsg(n, n.myinfo.Addr, "", sendmsg)
	// if err != nil {
	// 	jarvisbase.Warn("jarvisNode.RequestCtrl:BuildLocalSendMsg", zap.Error(err))

	// 	return err
	// }

	// if funcReply != nil {
	// 	n.mgrRequest.addRequestData(msg, funcReply)
	// }

	// n.PostMsg(msg, nil, nil, funcOnResult)

	return nil
}

// SendFile - send filedata to jarvisnode with addr
func (n *jarvisNode) SendFile(ctx context.Context, addr string, fd *pb.FileData,
	funcOnResult FuncOnProcMsgResult) error {

	sendmsg, err := BuildTransferFile(n, n.myinfo.Addr, addr, fd)
	if err != nil {
		jarvisbase.Warn("jarvisNode.SendFile:BuildTransferFile", zap.Error(err))

		return err
	}

	n.mgrClient2.addSendMsgTask(sendmsg, nil, funcOnResult)

	// msg, err := BuildLocalSendMsg(n, n.myinfo.Addr, "", sendmsg)
	// if err != nil {
	// 	jarvisbase.Warn("jarvisNode.SendFile:BuildLocalSendMsg", zap.Error(err))

	// 	return err
	// }

	// if funcReply != nil {
	// 	n.mgrRequest.addRequestData(msg, funcReply)
	// }

	// n.PostMsg(msg, nil, nil, funcOnResult)

	return nil
}

// onTimerRequestNodes
func (n *jarvisNode) onTimerRequestNodes(ctx context.Context) error {
	jarvisbase.Debug("jarvisNode.onTimerRequestNodes")

	return n.RequestNodes(ctx, nil)
}

// RequestNode - update node
func (n *jarvisNode) RequestNode(ctx context.Context, addr string,
	funcOnResult FuncOnProcMsgResult) error {

	ni := n.coredb.GetNode(addr)
	if ni != nil && !coredb.IsDeprecatedNode(ni) {
		sendmsg, err := BuildRequestNodes(n, n.myinfo.Addr, ni.Addr)
		if err != nil {
			jarvisbase.Warn("jarvisNode.RequestNode:BuildRequestNodes", zap.Error(err))

			return nil
		}

		n.mgrClient2.addSendMsgTask(sendmsg, nil, funcOnResult)

		// msg, err := BuildLocalSendMsg(n, n.myinfo.Addr, "", sendmsg)
		// if err != nil {
		// 	jarvisbase.Warn("jarvisNode.RequestNode:BuildLocalSendMsg", zap.Error(err))

		// 	return err
		// }

		// n.PostMsg(msg, nil, nil, funcOnResult)
	}

	return nil
}

// RequestNodes - request nodes
func (n *jarvisNode) RequestNodes(ctx context.Context, funcOnResult FuncOnGroupSendMsgResult) error {

	numsSend := 0

	var totalResults []*ClientGroupProcMsgResults

	//!! 在网络IO很快的时候，假设一共有2个节点，但第一个节点很快返回的话，可能还没全部发送完成，就产生回调
	//!! 所以这里分2次遍历
	n.coredb.ForEachMapNodes(func(key string, v *coredbpb.NodeInfo) error {
		if !coredb.IsDeprecatedNode(v) && n.mgrClient2.isConnected(v.Addr) {
			numsSend++
		}

		return nil
	})

	n.coredb.ForEachMapNodes(func(key string, v *coredbpb.NodeInfo) error {
		jarvisbase.Debug(fmt.Sprintf("jarvisNode.RequestNodes %v", v))

		if !coredb.IsDeprecatedNode(v) && n.mgrClient2.isConnected(v.Addr) {
			curResult := &ClientGroupProcMsgResults{}
			totalResults = append(totalResults, curResult)

			err := n.RequestNode(ctx, v.Addr,
				func(ctx context.Context, jarvisnode JarvisNode, lstResult []*ClientProcMsgResult) error {
					curResult.Results = lstResult

					if funcOnResult != nil {
						funcOnResult(ctx, jarvisnode, numsSend, totalResults)
					}

					return nil
				})
			if err != nil {
				jarvisbase.Warn("jarvisNode.RequestNodes:RequestNode", zap.Error(err))

				return nil
			}
		}

		return nil
	})

	return nil
}

// onMsgRequestNodes
func (n *jarvisNode) onMsgRequestNodes(ctx context.Context, msg *pb.JarvisMsg, stream pb.JarvisCoreServ_ProcMsgServer) error {
	// jarvisbase.Debug("jarvisNode.onMsgRequestNodes")

	if stream == nil {
		jarvisbase.Warn("jarvisNode.onMsgRequestNodes", zap.Error(ErrStreamNil))

		return ErrStreamNil
	}

	n.coredb.ForEachMapNodes(func(key string, v *coredbpb.NodeInfo) error {
		//!!! don't broadcast the localhost and deprecated node
		if IsLocalHostAddr(v.ServAddr) && v.Deprecated {
			return nil
		}

		mni := &pb.NodeBaseInfo{
			ServAddr: v.ServAddr,
			Addr:     v.Addr,
			Name:     v.Name,
		}

		// jarvisbase.Debug("jarvisNode.onMsgRequestNodes", jarvisbase.JSON("node", mni))

		sendmsg, err := BuildNodeInfo(n, n.myinfo.Addr, msg.SrcAddr, mni)
		if err != nil {
			jarvisbase.Warn("jarvisNode.onMsgRequestNodes:BuildNodeInfo", zap.Error(err))

			return err
		}

		err = n.sendMsg2ClientStream(stream, sendmsg)
		if err != nil {
			jarvisbase.Warn("jarvisNode.onMsgRequestNodes:sendMsg2ClientStream", zap.Error(err))

			return err
		}

		return nil
	})

	return nil
}

// FindNodeWithName - find node with name
func (n *jarvisNode) FindNodeWithName(name string) *coredbpb.NodeInfo {
	return n.coredb.FindMapNode(name)
}

// replyStream2
func (n *jarvisNode) replyStream2(msg *pb.JarvisMsg, stream pb.JarvisCoreServ_ProcMsgServer,
	rt pb.REPLYTYPE, strErr string) error {

	if stream == nil {
		return nil
	}

	sendmsg, err := BuildReply2(n, n.myinfo.Addr, msg.SrcAddr, rt, strErr, msg.MsgID)
	if err != nil {
		jarvisbase.Warn("jarvisNode.replyStream2:BuildReply2", zap.Error(err))

		return err
	}

	err = n.sendMsg2ClientStream(stream, sendmsg)
	if err != nil {
		jarvisbase.Warn("jarvisNode.replyStream2:SendMsg", zap.Error(err))

		return err
	}

	return nil
}

// replyTransferFile
func (n *jarvisNode) replyTransferFile(msg *pb.JarvisMsg, stream pb.JarvisCoreServ_ProcMsgServer,
	md5str string) error {

	sendmsg, err := BuildReplyTransferFile(n, n.myinfo.Addr, msg.SrcAddr, md5str, msg.MsgID)
	if err != nil {
		jarvisbase.Warn("jarvisNode.replyTransferFile:BuildReply2", zap.Error(err))

		return err
	}

	n.sendMsg2ClientStream(stream, sendmsg)
	// err = stream.SendMsg(sendmsg)
	if err != nil {
		jarvisbase.Warn("jarvisNode.replyTransferFile:sendMsg2ClientStream", zap.Error(err))

		return err
	}

	return nil
}

// onMsgTransferFile
func (n *jarvisNode) onMsgTransferFile(ctx context.Context, msg *pb.JarvisMsg,
	stream pb.JarvisCoreServ_ProcMsgServer) error {

	// jarvisbase.Debug("jarvisNode.onMsgTransferFile")

	if stream == nil {
		jarvisbase.Warn("jarvisNode.onMsgTransferFile", zap.Error(ErrStreamNil))

		return ErrStreamNil
	}

	fd := msg.GetFile()

	err := StoreLocalFile(fd)
	if err != nil {
		err1 := n.replyStream2(msg, stream, pb.REPLYTYPE_ERROR, err.Error())
		if err1 != nil {
			jarvisbase.Warn("jarvisNode.onMsgTransferFile:replyStream err", zap.Error(err1))

			return err1
		}

		return err
	}

	n.mgrEvent.onMsgEvent(ctx, EventOnTransferFile, msg)

	md5str := GetMD5String(fd.File)

	err = n.replyTransferFile(msg, stream, md5str)
	if err != nil {
		jarvisbase.Warn("jarvisNode.onMsgTransferFile:replyTransferFile", zap.Error(err))

		return err
	}

	return nil
}

// SetNodeTypeInfo - set node type and version
func (n *jarvisNode) SetNodeTypeInfo(nodetype string, nodetypeversion string) {
	n.myinfo.NodeType = nodetype
	n.myinfo.NodeTypeVersion = nodetypeversion
}

// onMsgRequestFile
func (n *jarvisNode) onMsgRequestFile(ctx context.Context, msg *pb.JarvisMsg,
	stream pb.JarvisCoreServ_ProcMsgServer) error {

	// jarvisbase.Debug("jarvisNode.onMsgRequestFile")

	if stream == nil {
		jarvisbase.Warn("jarvisNode.onMsgRequestFile", zap.Error(ErrStreamNil))

		return ErrStreamNil
	}

	n.mgrEvent.onMsgEvent(ctx, EventOnRequestFile, msg)

	rf := msg.GetRequestFile()

	buf, err := ioutil.ReadFile(rf.Filename)
	if err != nil {
		jarvisbase.Warn("jarvisNode.onMsgRequestFile:ReadFile", zap.Error(err))

		n.replyStream2(msg, stream, pb.REPLYTYPE_ERROR, err.Error())

		return err
	}

	fd := &pb.FileData{
		File:     buf,
		Filename: rf.Filename,
	}

	sendmsg, err := BuildReplyRequestFile(n, n.myinfo.Addr, msg.SrcAddr, fd)
	if err != nil {
		jarvisbase.Warn("jarvisNode.onMsgRequestFile:BuildReplyRequestFile", zap.Error(err))

		n.replyStream2(msg, stream, pb.REPLYTYPE_ERROR, err.Error())

		return err
	}

	err = n.sendMsg2ClientStream(stream, sendmsg)
	if err != nil {
		jarvisbase.Warn("jarvisNode.onMsgRequestFile:sendMsg2ClientStream", zap.Error(err))

		return err
	}

	return nil
}

// onMsgReplyRequestFile
func (n *jarvisNode) onMsgReplyRequestFile(ctx context.Context, msg *pb.JarvisMsg) error {
	// jarvisbase.Debug("jarvisNode.onMsgReplyRequestFile")

	// if msg.ReplyMsgID > 0 {
	// 	n.mgrRequest.onReplyRequest(ctx, n, msg)
	// }

	n.mgrEvent.onMsgEvent(ctx, EventOnReplyRequestFile, msg)

	return nil
}

// RequestFile - request node send filedata to me
func (n *jarvisNode) RequestFile(ctx context.Context, addr string, rf *pb.RequestFile,
	funcOnResult FuncOnProcMsgResult) error {

	sendmsg, err := BuildRequestFile(n, n.myinfo.Addr, addr, rf)
	if err != nil {
		jarvisbase.Warn("jarvisNode.RequestFile", zap.Error(err))

		return err
	}

	n.mgrClient2.addSendMsgTask(sendmsg, nil, funcOnResult)

	// msg, err := BuildLocalSendMsg(n, n.myinfo.Addr, "", sendmsg)
	// if err != nil {
	// 	jarvisbase.Warn("jarvisNode.RequestFile:BuildLocalSendMsg", zap.Error(err))

	// 	return err
	// }

	// if funcReply != nil {
	// 	n.mgrRequest.addRequestData(msg, funcReply)
	// }

	// n.PostMsg(msg, nil, nil, funcOnResult)

	return nil
}

// RegCtrl - register a ctrl
func (n *jarvisNode) RegCtrl(ctrltype string, ctrl Ctrl) error {
	n.mgrCtrl.Reg(ctrltype, ctrl)

	return nil
}

// PostMsg - like windows postMessage
func (n *jarvisNode) PostMsg(msg *pb.JarvisMsg, stream pb.JarvisCoreServ_ProcMsgServer,
	chanEnd chan int, funcOnResult FuncOnProcMsgResult) {

	n.mgrJasvisMsg.sendMsg(msg, stream, chanEnd, funcOnResult)
}

// AddNodeBaseInfo - add nodeinfo
func (n *jarvisNode) AddNodeBaseInfo(nbi *pb.NodeBaseInfo) error {

	cn := n.coredb.GetNode(nbi.Addr)
	if cn == nil {
		err := n.coredb.UpdNodeBaseInfo(nbi)
		if err != nil {
			jarvisbase.Warn("jarvisNode.AddNodeBaseInfo:UpdNodeBaseInfo", zap.Error(err))

			return err
		}

		cn = n.coredb.GetNode(nbi.Addr)
		if cn == nil {
			jarvisbase.Warn("jarvisNode.AddNodeBaseInfo:GetNode", zap.Error(ErrAssertGetNode))

			return ErrAssertGetNode
		}

		n.ConnectNode(cn, nil)

		return nil
	} else if cn.ConnType == coredbpb.CONNECTTYPE_UNKNOWN_CONN {
		n.ConnectNode(cn, nil)

		return nil
	}

	return nil
}

// checkMsgID
func (n *jarvisNode) checkMsgID(ctx context.Context, msg *pb.JarvisMsg) error {
	if msg.MsgType == pb.MSGTYPE_REPLY_CONNECT {
		return nil
	}

	// it is my message
	cn := n.GetCoreDB().GetNode(msg.SrcAddr)
	if cn == nil {
		return ErrUnknowNode
	}

	if msg.MsgID <= cn.LastRecvMsgID {
		jarvisbase.Warn("jarvisNode.checkMsgID",
			zap.String("destaddr", msg.DestAddr),
			zap.String("srcaddr", msg.SrcAddr),
			zap.Int64("msgid", msg.MsgID),
			zap.Int64("lasrrevmsgid", cn.LastRecvMsgID),
			jarvisbase.JSON("msg", msg))

		return ErrInvalidMsgID
	}

	if msg.LastMsgID > 0 {
		n.GetCoreDB().UpdMsgID(msg.SrcAddr, msg.LastMsgID, msg.MsgID)
	} else {
		n.GetCoreDB().UpdRecvMsgID(msg.SrcAddr, msg.MsgID)
	}

	return nil
}

// UpdateNode - update node
func (n *jarvisNode) UpdateNode(ctx context.Context, addr string, nodetype string, nodetypever string,
	funcOnResult FuncOnProcMsgResult) error {

	sendmsg, err := BuildUpdateNode(n, n.myinfo.Addr, addr, nodetype, nodetypever)
	if err != nil {
		jarvisbase.Warn("jarvisNode.RequestFile", zap.Error(err))

		return err
	}

	n.mgrClient2.addSendMsgTask(sendmsg, nil, funcOnResult)

	// msg, err := BuildLocalSendMsg(n, n.myinfo.Addr, "", sendmsg)
	// if err != nil {
	// 	jarvisbase.Warn("jarvisNode.RequestFile:BuildLocalSendMsg", zap.Error(err))

	// 	return err
	// }

	// if funcReply != nil {
	// 	n.mgrRequest.addRequestData(msg, funcReply)
	// }

	// n.PostMsg(msg, nil, nil, funcOnResult)

	return nil
}

// UpdateAllNodes - update all nodes
func (n *jarvisNode) UpdateAllNodes(ctx context.Context, nodetype string, nodetypever string,
	funcOnResult FuncOnGroupSendMsgResult) error {

	numsSend := 0

	var totalResults []*ClientGroupProcMsgResults

	//!! 在网络IO很快的时候，假设一共有2个节点，但第一个节点很快返回的话，可能还没全部发送完成，就产生回调
	//!! 所以这里分2次遍历
	n.coredb.ForEachMapNodes(func(addr string, ni *coredbpb.NodeInfo) error {
		if ni.NodeType == nodetype && ni.NodeTypeVersion != nodetypever {
			numsSend++
		}

		return nil
	})

	n.coredb.ForEachMapNodes(func(addr string, ni *coredbpb.NodeInfo) error {
		if ni.NodeType == nodetype && ni.NodeTypeVersion != nodetypever {

			curResult := &ClientGroupProcMsgResults{}
			totalResults = append(totalResults, curResult)

			err := n.UpdateNode(ctx, addr, nodetype, nodetypever,
				func(ctx context.Context, jarvisnode JarvisNode, lstResult []*ClientProcMsgResult) error {
					curResult.Results = lstResult
					// numsRecv++

					// jarvisbase.Debug("jarvisNode.UpdateAllNodes:FuncOnSendMsgResult",
					// 	zap.Int("numsRecv", numsRecv),
					// 	zap.Int("numsSend", numsSend))

					// totalResults = append(totalResults, &ClientGroupProcMsgResults{
					// 	Results: lstResult,
					// })

					// if len(lstResult) != 1 {
					// 	jarvisbase.Error("jarvisNode.UpdateAllNodes:FuncOnSendMsgResult", zap.Int("len", len(lstResult)))

					// 	totalResults = append(totalResults,
					// 		&ResultSendMsg{
					// 			Err: ErrFuncOnSendMsgResultLength,
					// 		})
					// } else {
					// 	totalResults = append(totalResults, lstResult[0])
					// }

					if funcOnResult != nil {
						funcOnResult(ctx, jarvisnode, numsSend, totalResults)
					}

					return nil
				})

			if err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}

// sendMsg2ClientStream
func (n *jarvisNode) sendMsg2ClientStream(stream pb.JarvisCoreServ_ProcMsgServer, sendmsg *pb.JarvisMsg) error {
	if stream == nil {
		jarvisbase.Warn("jarvisNode.sendMsg2ClientStream", zap.Error(ErrStreamNil))

		return ErrStreamNil
	}

	sendmsg.MsgID = n.GetCoreDB().GetNewSendMsgID(sendmsg.DestAddr)
	sendmsg.CurTime = time.Now().Unix()

	err := SignJarvisMsg(n.GetCoreDB().GetPrivateKey(), sendmsg)
	if err != nil {
		jarvisbase.Warn("jarvisNode.sendMsg2ClientStream:SignJarvisMsg", zap.Error(err))

		return err
	}

	err = stream.Send(sendmsg)
	if err != nil {
		jarvisbase.Warn("jarvisNode.sendMsg2ClientStream:sendmsg", zap.Error(err))

		return err
	}

	return nil
}
