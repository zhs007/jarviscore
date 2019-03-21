package jarviscore

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/zhs007/jarviscore/base"
	"github.com/zhs007/jarviscore/basedef"
	"github.com/zhs007/jarviscore/coredb"
	"github.com/zhs007/jarviscore/coredb/proto"
	pb "github.com/zhs007/jarviscore/proto"
	"go.uber.org/zap"
)

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
		coredb:  db,
		cfg:     cfg,
		mgrCtrl: &ctrlMgr{},
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

// onNormalMsg - proc JarvisMsg
func (n *jarvisNode) onNormalMsg(ctx context.Context, normal NormalTaskInfo) error {
	jarvisbase.Debug("jarvisNode.onNormalMsg",
		JSONMsg2Zap("msg", normal.Msg))

	// is timeout
	if IsTimeOut(normal.Msg) {
		jarvisbase.Warn("jarvisNode.onNormalMsg",
			zap.Error(ErrJarvisMsgTimeOut),
			JSONMsg2Zap("msg", normal.Msg))

		n.replyStream2(normal.Msg, normal.Stream, pb.REPLYTYPE_ERROR, ErrJarvisMsgTimeOut.Error())

		return nil
	}

	// proc connect msg
	if normal.Msg.MsgType == pb.MSGTYPE_CONNECT_NODE {
		// verify msg
		err := VerifyJarvisMsg(normal.Msg)
		if err != nil {
			jarvisbase.Warn("jarvisNode.onNormalMsg",
				zap.Error(err),
				JSONMsg2Zap("msg", normal.Msg))

			n.replyStream2(normal.Msg, normal.Stream, pb.REPLYTYPE_ERROR, err.Error())

			return nil
		}

		return n.onMsgConnectNode(ctx, normal.Msg, normal.Stream)
	}

	// if is not my msg, broadcast msg
	if n.myinfo.Addr != normal.Msg.DestAddr {
		//!!! 先不考虑转发协议
		// n.mgrClient2.addTask(msg, "", nil, nil)
	} else {
		// verify msg
		err := VerifyJarvisMsg(normal.Msg)
		if err != nil {
			jarvisbase.Warn("jarvisNode.onNormalMsg:VerifyJarvisMsg",
				zap.Error(err),
				JSONMsg2Zap("msg", normal.Msg))

			n.replyStream2(normal.Msg, normal.Stream, pb.REPLYTYPE_ERROR, err.Error())

			return nil
		}

		if normal.Stream != nil {
			err = n.checkMsgID(ctx, normal.Msg)
			if err != nil {
				jarvisbase.Warn("jarvisNode.onNormalMsg:checkMsgID", zap.Error(err))

				if err == ErrInvalidMsgID {
					n.replyStream2(normal.Msg, normal.Stream, pb.REPLYTYPE_ERRMSGID, "")
				} else {
					n.replyStream2(normal.Msg, normal.Stream, pb.REPLYTYPE_ERROR, err.Error())
				}

				return nil
			}
		}

		if normal.Msg.MsgType == pb.MSGTYPE_NODE_INFO {
			return n.onMsgNodeInfo(ctx, normal.Msg)
		} else if normal.Msg.MsgType == pb.MSGTYPE_REPLY_CONNECT {
			return n.onMsgReplyConnect(ctx, normal.Msg)
		} else if normal.Msg.MsgType == pb.MSGTYPE_REQUEST_CTRL {
			return n.onMsgRequestCtrl(ctx, normal.Msg, normal.Stream, normal.OnResult)
		} else if normal.Msg.MsgType == pb.MSGTYPE_REPLY_CTRL_RESULT {
			return n.onMsgCtrlResult(ctx, normal.Msg)
		} else if normal.Msg.MsgType == pb.MSGTYPE_REQUEST_NODES {
			return n.onMsgRequestNodes(ctx, normal.Msg, normal.Stream)
		} else if normal.Msg.MsgType == pb.MSGTYPE_TRANSFER_FILE {
			return n.onMsgTransferFile(ctx, normal.Msg, normal.Stream)
		} else if normal.Msg.MsgType == pb.MSGTYPE_REQUEST_FILE {
			return n.onMsgRequestFile(ctx, normal.Msg, normal.Stream)
		} else if normal.Msg.MsgType == pb.MSGTYPE_REPLY_REQUEST_FILE {
			return n.onMsgReplyRequestFile(ctx, normal.Msg)
		} else if normal.Msg.MsgType == pb.MSGTYPE_REPLY_TRANSFER_FILE {
			return n.onMsgReplyTransferFile(ctx, normal.Msg)
		} else if normal.Msg.MsgType == pb.MSGTYPE_REPLY2 {
			return n.onMsgReply2(ctx, normal.Msg)
		} else if normal.Msg.MsgType == pb.MSGTYPE_UPDATENODE {
			return n.onMsgUpdateNode(ctx, normal.Msg, normal.Stream)
		}

	}

	return nil
}

// OnMsg - proc JarvisMsg
func (n *jarvisNode) OnMsg(ctx context.Context, task JarvisTask) error {

	if task.Normal != nil {
		return n.onNormalMsg(ctx, *task.Normal)
	}

	return nil
}

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

	ts := time.Now().Unix()
	if cn.ConnType == coredbpb.CONNECTTYPE_UNKNOWN_CONN && ts > basedef.TimeReconnect+cn.LastConnectTime {
		cn.ConnectNums++
		cn.LastConnectTime = ts

		n.mgrClient2.addConnTask(cn.ServAddr, cn, funcOnResult)

		return nil
	}

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

	ts := time.Now().Unix()
	if node.ConnType == coredbpb.CONNECTTYPE_UNKNOWN_CONN && ts > basedef.TimeReconnect+node.LastConnectTime {
		node.ConnectNums++
		node.LastConnectTime = ts

		n.mgrClient2.addConnTask(node.ServAddr, node, funcOnResult)

		return nil
	}

	return nil
}

// replyCtrlResult
func (n *jarvisNode) replyCtrlResult(ctx context.Context, msg *pb.JarvisMsg, info string) error {

	ci := msg.GetCtrlInfo()
	if ci == nil {
		jarvisbase.Warn("jarvisNode.replyCtrlResult:GetCtrlInfo", zap.Error(ErrNoCtrlInfo))

		n.reply2(msg, pb.REPLYTYPE_ERROR, ErrNoCtrlInfo.Error())

		return ErrNoCtrlInfo
	}

	sendmsg2, err := BuildCtrlResult(n, n.myinfo.Addr, msg.SrcAddr, ci.CtrlID, msg.MsgID, info)

	if err != nil {
		jarvisbase.Warn("jarvisNode.replyCtrlResult:BuildCtrlResult", zap.Error(err))

		n.reply2(msg, pb.REPLYTYPE_ERROR, err.Error())

		return err
	}

	jarvisbase.Info("jarvisNode.replyCtrlResult",
		JSONMsg2Zap("msg", msg),
		JSONMsg2Zap("sendmsg", sendmsg2))

	n.mgrClient2.addSendMsgTask(sendmsg2, nil, nil)

	return nil
}

// reply2
func (n *jarvisNode) reply2(msg *pb.JarvisMsg, rt pb.REPLYTYPE, strErr string) error {

	sendmsg, err := BuildReply2(n, n.myinfo.Addr, msg.SrcAddr, rt, strErr, msg.MsgID)
	if err != nil {
		jarvisbase.Warn("jarvisNode.reply2:BuildReply2", zap.Error(err))

		return err
	}

	n.mgrClient2.addSendMsgTask(sendmsg, nil, nil)

	return nil
}

// runRequestCtrl
func (n *jarvisNode) runRequestCtrl(ctx context.Context, msg *pb.JarvisMsg,
	stream pb.JarvisCoreServ_ProcMsgServer, funcOnResult FuncOnProcMsgResult) {

	ci := msg.GetCtrlInfo()
	ret, err := n.mgrCtrl.Run(ci)
	if err != nil {
		if ret != nil {
			n.replyCtrlResult(ctx, msg, string(ret))
		}

		n.replyCtrlResult(ctx, msg, err.Error())

		return
	}

	n.replyCtrlResult(ctx, msg, string(ret))
}

// onMsgRequestCtrl
func (n *jarvisNode) onMsgRequestCtrl(ctx context.Context, msg *pb.JarvisMsg,
	stream pb.JarvisCoreServ_ProcMsgServer, funcOnResult FuncOnProcMsgResult) error {

	jarvisbase.Info("jarvisNode.onMsgRequestCtrl:recvmsg",
		JSONMsg2Zap("msg", msg))

	n.replyStream2(msg, stream, pb.REPLYTYPE_ISME, "")

	n.mgrEvent.onMsgEvent(ctx, EventOnCtrl, msg)

	go n.runRequestCtrl(ctx, msg, stream, funcOnResult)

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

// // onMsgLocalSendMsg
// func (n *jarvisNode) onMsgLocalSendMsg(ctx context.Context, msg *pb.JarvisMsg,
// 	funcOnResult FuncOnProcMsgResult) error {

// 	sendmsg := msg.GetMsg()

// 	n.mgrClient2.addSendMsgTask(sendmsg, nil, funcOnResult)

// 	return nil
// }

// onMsgUpdateNode
func (n *jarvisNode) onMsgUpdateNode(ctx context.Context, msg *pb.JarvisMsg, stream pb.JarvisCoreServ_ProcMsgServer) error {

	jarvisbase.Info("jarvisNode.onMsgUpdateNode",
		JSONMsg2Zap("msg", msg))

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

// RequestCtrl - send ctrl to jarvisnode with addr
func (n *jarvisNode) RequestCtrl(ctx context.Context, addr string, ci *pb.CtrlInfo,
	funcOnResult FuncOnProcMsgResult) error {

	sendmsg, err := BuildRequestCtrl(n, n.myinfo.Addr, addr, ci)
	if err != nil {
		jarvisbase.Warn("jarvisNode.RequestCtrl", zap.Error(err))

		return err
	}

	n.mgrClient2.addSendMsgTask(sendmsg, nil, funcOnResult)

	return nil
}

// SendFile - send filedata to jarvisnode with addr
func (n *jarvisNode) SendFile(ctx context.Context, addr string, fd *pb.FileData,
	funcOnResult FuncOnProcMsgResult) error {

	// ProcFileDataWithBuff(fd.File, )

	sendmsg, err := BuildTransferFile(n, n.myinfo.Addr, addr, fd)
	if err != nil {
		jarvisbase.Warn("jarvisNode.SendFile:BuildTransferFile", zap.Error(err))

		return err
	}

	n.mgrClient2.addSendMsgTask(sendmsg, nil, funcOnResult)

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
				func(ctx context.Context, jarvisnode JarvisNode, lstResult []*JarvisMsgInfo) error {
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

// FindNode - find node
func (n *jarvisNode) FindNode(addr string) *coredbpb.NodeInfo {
	return n.coredb.GetNode(addr)
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
	if err != nil {
		jarvisbase.Warn("jarvisNode.replyTransferFile:sendMsg2ClientStream", zap.Error(err))

		return err
	}

	return nil
}

// onMsgTransferFile
func (n *jarvisNode) onMsgTransferFile(ctx context.Context, msg *pb.JarvisMsg,
	stream pb.JarvisCoreServ_ProcMsgServer) error {

	if stream == nil {
		jarvisbase.Warn("jarvisNode.onMsgTransferFile", zap.Error(ErrStreamNil))

		return ErrStreamNil
	}

	n.replyStream2(msg, stream, pb.REPLYTYPE_IGOTIT, "")

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

// replyFile
func (n *jarvisNode) replyFile(ctx context.Context, msg *pb.JarvisMsg, rf *pb.RequestFile,
	stream pb.JarvisCoreServ_ProcMsgServer) error {

	err := ProcFileData(rf.Filename, func(fd *pb.FileData, isend bool) error {
		jarvisbase.Info("jarvisNode.replyFile",
			zap.Int("buflen", len(fd.File)),
			zap.Int64("filelen", fd.TotalLength))

		fd.Filename = rf.Filename

		sendmsg, err := BuildReplyRequestFile(n, n.myinfo.Addr, msg.SrcAddr, fd, msg.MsgID)
		if err != nil {
			jarvisbase.Warn("jarvisNode.replyFile:BuildReplyRequestFile", zap.Error(err))

			n.replyStream2(msg, stream, pb.REPLYTYPE_ERROR, err.Error())

			return err
		}

		err = n.sendMsg2ClientStream(stream, sendmsg)
		if err != nil {
			jarvisbase.Warn("jarvisNode.replyFile:sendMsg2ClientStream", zap.Error(err))

			return err
		}

		return nil
	})
	if err != nil {
		n.replyStream2(msg, stream, pb.REPLYTYPE_ERROR, err.Error())

		return err
	}

	return nil
}

// replyPartFile
func (n *jarvisNode) replyPartFile(ctx context.Context, msg *pb.JarvisMsg, rf *pb.RequestFile,
	stream pb.JarvisCoreServ_ProcMsgServer) error {

	fl, err := GetFileLength(rf.Filename)
	if err != nil {
		jarvisbase.Warn("jarvisNode.replyPartFile:GetFileLength", zap.Error(err))

		n.replyStream2(msg, stream, pb.REPLYTYPE_ERROR, err.Error())

		return err
	}

	fdata, err := os.Open(rf.Filename)
	if err != nil {
		jarvisbase.Warn("jarvisNode.replyPartFile:Open", zap.Error(err))

		n.replyStream2(msg, stream, pb.REPLYTYPE_ERROR, err.Error())

		return err
	}

	defer fdata.Close()

	off, err := fdata.Seek(rf.Start, 0)
	if err != nil {
		jarvisbase.Warn("jarvisNode.replyPartFile:Open", zap.Error(err))

		n.replyStream2(msg, stream, pb.REPLYTYPE_ERROR, err.Error())

		return err
	}

	if off != rf.Start {
		jarvisbase.Warn("jarvisNode.replyPartFile:Open", zap.Error(ErrInvalidSeekFileOffset))

		n.replyStream2(msg, stream, pb.REPLYTYPE_ERROR, ErrInvalidSeekFileOffset.Error())

		return ErrInvalidSeekFileOffset
	}

	buf := make([]byte, rf.Length)

	len, err := fdata.Read(buf)
	if err != nil {
		jarvisbase.Warn("jarvisNode.replyPartFile:Read", zap.Error(err))

		n.replyStream2(msg, stream, pb.REPLYTYPE_ERROR, err.Error())

		return err
	}

	fd := &pb.FileData{
		File:        buf,
		Filename:    rf.Filename,
		Ft:          pb.FileType_FT_BINARY,
		Start:       rf.Start,
		Length:      int64(len),
		TotalLength: fl,
		Md5String:   GetMD5String(buf),
	}

	sendmsg, err := BuildReplyRequestFile(n, n.myinfo.Addr, msg.SrcAddr, fd, msg.MsgID)
	if err != nil {
		jarvisbase.Warn("jarvisNode.replyPartFile:BuildReplyRequestFile", zap.Error(err))

		n.replyStream2(msg, stream, pb.REPLYTYPE_ERROR, err.Error())

		return err
	}

	err = n.sendMsg2ClientStream(stream, sendmsg)
	if err != nil {
		jarvisbase.Warn("jarvisNode.replyPartFile:sendMsg2ClientStream", zap.Error(err))

		return err
	}

	return nil
}

// onMsgRequestFile
func (n *jarvisNode) onMsgRequestFile(ctx context.Context, msg *pb.JarvisMsg,
	stream pb.JarvisCoreServ_ProcMsgServer) error {

	if stream == nil {
		jarvisbase.Warn("jarvisNode.onMsgRequestFile", zap.Error(ErrStreamNil))

		return ErrStreamNil
	}

	n.replyStream2(msg, stream, pb.REPLYTYPE_IGOTIT, "")

	n.mgrEvent.onMsgEvent(ctx, EventOnRequestFile, msg)

	rf := msg.GetRequestFile()

	if rf.Length > 0 {
		if rf.Length > basedef.BigFileLength {
			rf.Length = basedef.BigFileLength
		}

		return n.replyPartFile(ctx, msg, rf, stream)
	}

	return n.replyFile(ctx, msg, rf, stream)
}

// onMsgReplyRequestFile
func (n *jarvisNode) onMsgReplyRequestFile(ctx context.Context, msg *pb.JarvisMsg) error {

	jarvisbase.Info("jarvisNode.onMsgReplyRequestFile")

	fd := msg.GetFile()
	if fd == nil {
		jarvisbase.Warn("jarvisNode.onMsgReplyRequestFile", zap.Error(ErrNoFileData))

		return ErrNoFileData
	}

	if fd.Md5String == "" {
		jarvisbase.Warn("jarvisNode.onMsgReplyRequestFile", zap.Error(ErrFileDataNoMD5String))

		return ErrFileDataNoMD5String
	}

	if fd.Md5String != GetMD5String(fd.File) {
		jarvisbase.Warn("jarvisNode.onMsgReplyRequestFile", zap.Error(ErrInvalidFileDataMD5String))

		return ErrInvalidFileDataMD5String
	}

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

	return nil
}

// RegCtrl - register a ctrl
func (n *jarvisNode) RegCtrl(ctrltype string, ctrl Ctrl) error {
	n.mgrCtrl.Reg(ctrltype, ctrl)

	return nil
}

// PostMsg - like windows postMessage
func (n *jarvisNode) PostMsg(normal *NormalTaskInfo, chanEnd chan int) {

	n.mgrJasvisMsg.sendMsg(normal, chanEnd)
}

// PostStreamMsg - like windows postMessage
func (n *jarvisNode) PostStreamMsg(stream *StreamTaskInfo, chanEnd chan int) {

	n.mgrJasvisMsg.sendStreamMsg(stream, chanEnd)
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
			JSONMsg2Zap("msg", msg))

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
				func(ctx context.Context, jarvisnode JarvisNode, lstResult []*JarvisMsgInfo) error {
					curResult.Results = lstResult

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

	// sendmsg.MsgID = n.GetCoreDB().GetNewSendMsgID(sendmsg.DestAddr)
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

// BuildStatus - build jarviscorepb.JarvisNodeStatus
func (n *jarvisNode) BuildStatus() *pb.JarvisNodeStatus {
	ns := &pb.JarvisNodeStatus{
		MyBaseInfo: &pb.NodeBaseInfo{
			ServAddr:        n.myinfo.ServAddr,
			Addr:            n.myinfo.Addr,
			Name:            n.myinfo.Name,
			NodeTypeVersion: n.myinfo.NodeTypeVersion,
			NodeType:        n.myinfo.NodeType,
			CoreVersion:     n.myinfo.CoreVersion,
		},
	}

	n.mgrClient2.BuildNodeStatus(ns)

	return ns
}
