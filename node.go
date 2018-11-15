package jarviscore

import (
	"context"

	"github.com/zhs007/jarviscore/base"
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
	GetCoreDB() *CoreDB
	// SendCtrl - send ctrl to jarvisnode with addr
	SendCtrl(ctx context.Context, addr string, ci *pb.CtrlInfo) error

	// AddCtrl2List - add ctrl msg to tasklist
	AddCtrl2List(addr string, ci *pb.CtrlInfo) error

	// OnMsg - proc JarvisMsg
	OnMsg(ctx context.Context, msg *pb.JarvisMsg, stream pb.JarvisCoreServ_ProcMsgServer) error

	// GetMyInfo - get my nodeinfo
	GetMyInfo() *BaseInfo
}

// jarvisNode -
type jarvisNode struct {
	myinfo       BaseInfo
	coredb       *CoreDB
	mgrJasvisMsg *jarvisMsgMgr
	mgrClient2   *jarvisClient2
	serv2        *jarvisServer2
}

const (
	nodeinfoCacheSize       = 32
	randomMax         int64 = 0x7fffffffffffffff
	stateNormal             = 0
	stateStart              = 1
	stateEnd                = 2
)

// NewNode -
func NewNode(baseinfo BaseInfo) JarvisNode {
	db, err := newCoreDB()
	if err != nil {
		jarvisbase.Error("NewNode:newCoreDB", zap.Error(err))

		return nil
	}

	node := &jarvisNode{
		myinfo: baseinfo,
		coredb: db,
	}

	err = node.coredb.loadPrivateKeyEx()
	if err != nil {
		jarvisbase.Error("NewNode:loadPrivateKey", zap.Error(err))

		return nil
	}

	node.coredb.loadAllNodes()

	node.myinfo.Addr = node.coredb.privKey.ToAddress()
	node.myinfo.Name = config.BaseNodeInfo.NodeName
	node.myinfo.BindAddr = config.BaseNodeInfo.BindAddr
	node.myinfo.ServAddr = config.BaseNodeInfo.ServAddr

	// mgrJasvisMsg
	node.mgrJasvisMsg = newJarvisMsgMgr(node)

	// mgrClient2
	node.mgrClient2 = newClient2(node)

	return node
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
	go n.coredb.ankaDB.Start(coredbctx)

	msgmgrctx, msgmgrcancel := context.WithCancel(ctx)
	defer msgmgrcancel()
	go n.mgrJasvisMsg.start(msgmgrctx)

	jarvisbase.Info("StartServer", zap.String("ServAddr", n.myinfo.ServAddr))
	n.serv2, err = newServer2(n)
	if err != nil {
		return err
	}

	servctx, servcancel := context.WithCancel(ctx)
	defer servcancel()
	go n.serv2.Start(servctx)

	n.connectNode(config.RootServAddr)
	n.connectAllNodes()

	for {
		select {
		case <-ctx.Done():
			n.Stop()
			return nil
		}
	}
}

// GetCoreDB - get coredb
func (n *jarvisNode) GetCoreDB() *CoreDB {
	return n.coredb
}

// SendCtrl - send ctrl to jarvisnode with addr
func (n *jarvisNode) SendCtrl(ctx context.Context, addr string, ci *pb.CtrlInfo) error {
	msg, err := BuildRequestCtrl(n.coredb.privKey, 0, n.myinfo.Addr, addr, ci)
	if err != nil {
		jarvisbase.Debug("jarvisNode.SendCtrl", zap.Error(err))

		return err
	}

	err = n.mgrClient2.sendMsg(ctx, msg)
	if err != nil {
		jarvisbase.Debug("jarvisNode.SendCtrl:sendMsg", zap.Error(err))

		return err
	}

	return nil
	// return n.requestCtrl(ctx, addr, ctrltype, []byte(command))
}

// OnMsg - proc JarvisMsg
func (n *jarvisNode) OnMsg(ctx context.Context, msg *pb.JarvisMsg, stream pb.JarvisCoreServ_ProcMsgServer) error {
	jarvisbase.Debug("jarvisNode.OnMsg", jarvisbase.JSON("msg", msg))

	// is timeout
	if IsTimeOut(msg) {
		jarvisbase.Debug("jarvisNode.OnMsg", zap.Error(ErrJarvisMsgTimeOut))

		return nil
	}

	// proc local msg
	if msg.MsgType == pb.MSGTYPE_LOCAL_CONNECT_OTHER {
		// verify msg
		err := VerifyJarvisMsg(msg)
		if err != nil {
			jarvisbase.Debug("jarvisNode.OnMsg", zap.Error(err))

			return nil
		}

		return n.onMsgLocalConnect(ctx, msg)
	}

	// proc connect msg
	if msg.MsgType == pb.MSGTYPE_CONNECT_NODE {
		// verify msg
		err := VerifyJarvisMsg(msg)
		if err != nil {
			jarvisbase.Debug("jarvisNode.OnMsg", zap.Error(err))

			return nil
		}

		return n.onMsgConnectNode(ctx, msg, stream)
	}

	// if is not my msg, broadcast msg
	if n.myinfo.Addr != msg.DestAddr {
		n.mgrClient2.broadCastMsg(ctx, msg)
	} else {
		// verify msg
		err := VerifyJarvisMsg(msg)
		if err != nil {
			jarvisbase.Debug("jarvisNode.OnMsg", zap.Error(err))

			return nil
		}

		if msg.MsgType == pb.MSGTYPE_NODE_INFO {
			return n.onMsgNodeInfo(ctx, msg)
		} else if msg.MsgType == pb.MSGTYPE_REPLY_CONNECT {
			return n.onMsgReplyConnect(ctx, msg)
		} else if msg.MsgType == pb.MSGTYPE_REQUEST_CTRL {
			return n.onMsgRequestCtrl(ctx, msg)
		} else if msg.MsgType == pb.MSGTYPE_REPLY {
			return n.onMsgReply(ctx, msg)
		} else if msg.MsgType == pb.MSGTYPE_REPLY_CTRL_RESULT {
			return n.onMsgCtrlResult(ctx, msg)
		} else if msg.MsgType == pb.MSGTYPE_LOCAL_SENDMSG {
			return n.onMsgLocalSendMsg(ctx, msg)
		}

	}

	return nil
}

// onMsgLocalConnect
func (n *jarvisNode) onMsgLocalConnect(ctx context.Context, msg *pb.JarvisMsg) error {
	ci := msg.GetConnInfo()

	cn := n.coredb.findNodeWithServAddr(ci.ServAddr)
	if cn == nil {
		return n.mgrClient2.connectNode(ctx, ci.ServAddr)
	} else if !cn.ConnectNode {
		return n.mgrClient2.connectNode(ctx, cn.ServAddr)
	}

	return nil
}

// onMsgNodeInfo
func (n *jarvisNode) onMsgNodeInfo(ctx context.Context, msg *pb.JarvisMsg) error {
	ni := msg.GetNodeInfo()
	cn := n.coredb.getNode(ni.Addr)
	if cn == nil {
		err := n.coredb.insNode(ni)
		if err != nil {
			jarvisbase.Debug("jarvisNode.onMsgNodeInfo:insNode", zap.Error(err))

			return err
		}

		return n.mgrClient2.connectNode(ctx, ni.ServAddr)
	} else if !cn.ConnectNode {
		return n.mgrClient2.connectNode(ctx, ni.ServAddr)
	}

	return nil
}

// onMsgConnectNode
func (n *jarvisNode) onMsgConnectNode(ctx context.Context, msg *pb.JarvisMsg, stream pb.JarvisCoreServ_ProcMsgServer) error {
	if stream == nil {
		jarvisbase.Debug("jarvisNode.onMsgConnectNode", zap.Error(ErrStreamNil))

		return ErrStreamNil
	}

	ci := msg.GetConnInfo()

	mni := &pb.NodeBaseInfo{
		ServAddr: n.myinfo.ServAddr,
		Addr:     n.myinfo.Addr,
		Name:     n.myinfo.Name,
	}

	sendmsg, err := BuildReplyConn(n.coredb.privKey, 0, n.myinfo.Addr, ci.MyInfo.Addr, mni)
	if err != nil {
		jarvisbase.Debug("jarvisNode.onMsgConnectNode:BuildReplyConn", zap.Error(err))

		return err
	}
	// SignJarvisMsg(n.coredb.privKey, sendmsg)
	// jarvisbase.Debug("jarvisNode.onMsgConnectNode:sendmsg", jarvisbase.JSON("msg", sendmsg))
	err = stream.Send(sendmsg)
	if err != nil {
		jarvisbase.Debug("jarvisNode.onMsgConnectNode:sendmsg", zap.Error(err))

		return err
	}

	cn := n.coredb.getNode(ci.MyInfo.Addr)
	if cn == nil {
		err := n.coredb.insNode(ci.MyInfo)
		if err != nil {
			jarvisbase.Debug("jarvisNode.onMsgConnectNode:insNode", zap.Error(err))

			return err
		}

		cn = n.coredb.getNode(ci.MyInfo.Addr)

		cn.ConnectMe = true

		msg, err := BuildLocalConnectOther(n.coredb.privKey, 0, n.myinfo.Addr, ci.MyInfo.Addr,
			ci.MyInfo.ServAddr, ci.MyInfo)
		if err != nil {
			jarvisbase.Debug("jarvisNode.onMsgConnectNode:BuildLocalConnectOther", zap.Error(err))

			return err
		}
		// SignJarvisMsg(n.coredb.privKey, msg)
		n.mgrJasvisMsg.sendMsg(msg, nil, nil)
	} else if !cn.ConnectNode {
		msg, err := BuildLocalConnectOther(n.coredb.privKey, 0, n.myinfo.Addr, ci.MyInfo.Addr,
			ci.MyInfo.ServAddr, ci.MyInfo)
		if err != nil {
			jarvisbase.Debug("jarvisNode.onMsgConnectNode:BuildLocalConnectOther", zap.Error(err))

			return err
		}
		// SignJarvisMsg(n.coredb.privKey, msg)
		n.mgrJasvisMsg.sendMsg(msg, nil, nil)
	}

	jarvisbase.Debug("jarvisNode.onMsgConnectNode:end")

	return nil
}

// onMsgReplyConnect
func (n *jarvisNode) onMsgReplyConnect(ctx context.Context, msg *pb.JarvisMsg) error {
	ni := msg.GetNodeInfo()
	cn := n.coredb.getNode(ni.Addr)
	if cn == nil {
		err := n.coredb.insNode(ni)
		if err != nil {
			jarvisbase.Debug("jarvisNode.onMsgReplyConnect:insNode", zap.Error(err))

			return err
		}

		cn = n.coredb.getNode(ni.Addr)
	}

	cn.ConnectNode = true

	n.coredb.updNodeBaseInfo(ni)

	return nil
}

// GetMyInfo - get my nodeinfo
func (n *jarvisNode) GetMyInfo() *BaseInfo {
	return &n.myinfo
}

// connectAllNodes - connect all nodes
func (n *jarvisNode) connectAllNodes() error {
	// nbi := &pb.NodeBaseInfo{
	// 	ServAddr: n.myinfo.ServAddr,
	// 	Addr:     n.myinfo.Addr,
	// 	Name:     n.myinfo.Name,
	// }

	for _, node := range n.coredb.mapNodes {
		n.connectNode(node.ServAddr)

		// msg := BuildLocalConnectOther(0, n.myinfo.Addr, node.Addr, nbi)
		// SignJarvisMsg(n.coredb.privKey, msg)
		// n.mgrJasvisMsg.sendMsg(msg, nil)
	}

	return nil
}

// connectNode - connect node
func (n *jarvisNode) connectNode(servaddr string) error {
	nbi := &pb.NodeBaseInfo{
		ServAddr: n.myinfo.ServAddr,
		Addr:     n.myinfo.Addr,
		Name:     n.myinfo.Name,
	}

	msg, err := BuildLocalConnectOther(n.coredb.privKey, 0, n.myinfo.Addr, "", servaddr, nbi)
	if err != nil {
		jarvisbase.Debug("jarvisNode.connectNode:BuildLocalConnectOther", zap.Error(err))

		return err
	}

	// SignJarvisMsg(n.coredb.privKey, msg)
	n.mgrJasvisMsg.sendMsg(msg, nil, nil)

	return nil
}

// onMsgRequestCtrl
func (n *jarvisNode) onMsgRequestCtrl(ctx context.Context, msg *pb.JarvisMsg) error {
	sendmsg, err := BuildReply(n.coredb.privKey, 0, n.myinfo.Addr, msg.SrcAddr, pb.REPLYTYPE_OK)
	if err != nil {
		jarvisbase.Debug("jarvisNode.onMsgRequestCtrl:BuildReply", zap.Error(err))

		return err
	}

	err = n.mgrClient2.sendMsg(ctx, sendmsg)
	if err != nil {
		jarvisbase.Debug("jarvisNode.onMsgRequestCtrl:sendMsg", zap.Error(err))

		return err
	}

	ci := msg.GetCtrlInfo()
	ret, err := mgrCtrl.Run(ci)
	if err != nil {
		sendmsg2, err := BuildCtrlResult(n.coredb.privKey, 0, n.myinfo.Addr, msg.SrcAddr, ci.CtrlID, err.Error())
		err = n.mgrClient2.sendMsg(ctx, sendmsg2)
		if err != nil {
			jarvisbase.Debug("jarvisNode.onMsgRequestCtrl:sendMsg", zap.Error(err))

			return err
		}

		return nil
	}

	sendmsg2, err := BuildCtrlResult(n.coredb.privKey, 0, n.myinfo.Addr, msg.SrcAddr, ci.CtrlID, string(ret))
	err = n.mgrClient2.sendMsg(ctx, sendmsg2)
	if err != nil {
		jarvisbase.Debug("jarvisNode.onMsgRequestCtrl:sendMsg", zap.Error(err))

		return err
	}

	return nil
}

// onMsgReply
func (n *jarvisNode) onMsgReply(ctx context.Context, msg *pb.JarvisMsg) error {
	return nil
}

// onMsgCtrlResult
func (n *jarvisNode) onMsgCtrlResult(ctx context.Context, msg *pb.JarvisMsg) error {
	return nil
}

// onMsgLocalSendMsg
func (n *jarvisNode) onMsgLocalSendMsg(ctx context.Context, msg *pb.JarvisMsg) error {
	sendmsg := msg.GetMsg()

	err := n.mgrClient2.sendMsg(ctx, sendmsg)
	if err != nil {
		jarvisbase.Debug("jarvisNode.onMsgLocalSendMsg:sendMsg", zap.Error(err))

		return err
	}

	return nil
}

// AddCtrl2List - add ctrl msg to tasklist
func (n *jarvisNode) AddCtrl2List(addr string, ci *pb.CtrlInfo) error {
	msg, err := BuildRequestCtrl(n.coredb.privKey, 0, n.myinfo.Addr, addr, ci)
	if err != nil {
		jarvisbase.Debug("jarvisNode.AddCtrl2List", zap.Error(err))

		return err
	}

	// SignJarvisMsg(n.coredb.privKey, msg)
	n.mgrJasvisMsg.sendMsg(msg, nil, nil)

	return nil
}
