package jarviscore

import (
	"context"
	"time"

	"github.com/zhs007/jarviscore/base"
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
	GetCoreDB() *CoreDB
	// RequestCtrl - send ctrl to jarvisnode with addr
	RequestCtrl(ctx context.Context, addr string, ci *pb.CtrlInfo) error

	// AddCtrl2List - add ctrl msg to tasklist
	AddCtrl2List(addr string, ci *pb.CtrlInfo) error

	// OnMsg - proc JarvisMsg
	OnMsg(ctx context.Context, msg *pb.JarvisMsg, stream pb.JarvisCoreServ_ProcMsgServer) error

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
}

// jarvisNode -
type jarvisNode struct {
	myinfo       BaseInfo
	coredb       *CoreDB
	mgrJasvisMsg *jarvisMsgMgr
	mgrClient2   *jarvisClient2
	serv2        *jarvisServer2
	mgrEvent     *eventMgr
	cfg          *Config
}

const (
	nodeinfoCacheSize       = 32
	randomMax         int64 = 0x7fffffffffffffff
	stateNormal             = 0
	stateStart              = 1
	stateEnd                = 2
)

// NewNode -
func NewNode(cfg *Config) JarvisNode {
	db, err := newCoreDB(cfg)
	if err != nil {
		jarvisbase.Error("NewNode:newCoreDB", zap.Error(err))

		return nil
	}

	node := &jarvisNode{
		myinfo: BaseInfo{
			Name:     cfg.BaseNodeInfo.NodeName,
			BindAddr: cfg.BaseNodeInfo.BindAddr,
			ServAddr: cfg.BaseNodeInfo.ServAddr,
		},
		coredb: db,
		cfg:    cfg,
	}

	// event
	node.mgrEvent = newEventMgr(node)
	node.mgrEvent.regNodeEventFunc(EventOnNodeConnected, onNodeConnected)
	node.mgrEvent.regNodeEventFunc(EventOnIConnectNode, onIConnectNode)

	err = node.coredb.loadPrivateKeyEx()
	if err != nil {
		jarvisbase.Error("NewNode:loadPrivateKey", zap.Error(err))

		return nil
	}

	node.coredb.loadAllNodes()

	node.myinfo.Addr = node.coredb.privKey.ToAddress()
	node.myinfo.Name = cfg.BaseNodeInfo.NodeName
	node.myinfo.BindAddr = cfg.BaseNodeInfo.BindAddr
	node.myinfo.ServAddr = cfg.BaseNodeInfo.ServAddr

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

	n.connectNode(n.cfg.RootServAddr)
	n.connectAllNodes()

	tickerRequestChild := time.NewTicker(time.Duration(n.cfg.TimeRequestChild) * time.Second)

	for {
		select {
		case <-tickerRequestChild.C:
			n.onTimerRequestNodes()
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

// sendCtrl - send ctrl to jarvisnode with addr
func (n *jarvisNode) sendCtrl(ctx context.Context, addr string, ci *pb.CtrlInfo) error {
	msg, err := BuildRequestCtrl(n.coredb.privKey, 0, n.myinfo.Addr, addr, ci)
	if err != nil {
		jarvisbase.Debug("jarvisNode.SendCtrl", zap.Error(err))

		return err
	}

	n.mgrClient2.addTask(msg, "")

	jarvisbase.Debug("jarvisNode.SendCtrl", jarvisbase.JSON("msg", msg))

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
	if msg.MsgType == pb.MSGTYPE_LOCAL_CONNECT_OTHER ||
		msg.MsgType == pb.MSGTYPE_LOCAL_SENDMSG ||
		msg.MsgType == pb.MSGTYPE_LOCAL_REQUEST_NODES {

		// verify msg
		err := VerifyJarvisMsg(msg)
		if err != nil {
			jarvisbase.Debug("jarvisNode.OnMsg", zap.Error(err))

			return nil
		}

		if msg.MsgType == pb.MSGTYPE_LOCAL_CONNECT_OTHER {
			return n.onMsgLocalConnect(ctx, msg)
		} else if msg.MsgType == pb.MSGTYPE_LOCAL_SENDMSG {
			return n.onMsgLocalSendMsg(ctx, msg)
		} else if msg.MsgType == pb.MSGTYPE_LOCAL_REQUEST_NODES {
			return n.onMsgLocalRequesrNodes(ctx, msg)
		}
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
		n.mgrClient2.addTask(msg, "")
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
		} else if msg.MsgType == pb.MSGTYPE_REQUEST_NODES {
			return n.onMsgRequestNodes(ctx, msg, stream)
		}

	}

	return nil
}

// onMsgLocalConnect
func (n *jarvisNode) onMsgLocalConnect(ctx context.Context, msg *pb.JarvisMsg) error {
	ci := msg.GetConnInfo()

	// if is me, return
	if ci.ServAddr == n.myinfo.ServAddr {
		return nil
	}

	cn := n.coredb.findNodeWithServAddr(ci.ServAddr)
	if cn == nil {
		n.mgrClient2.addTask(nil, ci.ServAddr)

		return nil
	}

	// if is me, return
	if cn.Addr == n.myinfo.Addr {
		return nil
	}

	if !cn.ConnectNode {
		n.mgrClient2.addTask(nil, cn.ServAddr)

		return nil
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

		n.mgrClient2.addTask(nil, ni.ServAddr)

		return nil
	} else if !cn.ConnectNode {
		n.mgrClient2.addTask(nil, ni.ServAddr)

		return nil
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

		// cn.ConnectMe = true
		n.mgrEvent.onNodeEvent(ctx, EventOnNodeConnected, cn)

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

	// cn.ConnectNode = true
	n.mgrEvent.onNodeEvent(ctx, EventOnIConnectNode, cn)

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

	n.mgrClient2.addTask(sendmsg, "")

	n.mgrEvent.onMsgEvent(ctx, EventOnCtrl, msg)

	ci := msg.GetCtrlInfo()
	ret, err := mgrCtrl.Run(ci)
	if err != nil {
		sendmsg2, err := BuildCtrlResult(n.coredb.privKey, 0, n.myinfo.Addr, msg.SrcAddr, ci.CtrlID, err.Error())
		if err != nil {
			jarvisbase.Debug("jarvisNode.onMsgRequestCtrl:BuildCtrlResult", zap.Error(err))

			return err
		}

		n.mgrClient2.addTask(sendmsg2, "")

		return nil
	}

	sendmsg2, err := BuildCtrlResult(n.coredb.privKey, 0, n.myinfo.Addr, msg.SrcAddr, ci.CtrlID, string(ret))
	n.mgrClient2.addTask(sendmsg2, "")

	return nil
}

// onMsgReply
func (n *jarvisNode) onMsgReply(ctx context.Context, msg *pb.JarvisMsg) error {
	return nil
}

// onMsgCtrlResult
func (n *jarvisNode) onMsgCtrlResult(ctx context.Context, msg *pb.JarvisMsg) error {
	n.mgrEvent.onMsgEvent(ctx, EventOnCtrlResult, msg)

	return nil
}

// onMsgLocalSendMsg
func (n *jarvisNode) onMsgLocalSendMsg(ctx context.Context, msg *pb.JarvisMsg) error {
	sendmsg := msg.GetMsg()

	n.mgrClient2.addTask(sendmsg, "")

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

// onNodeConnected - func event
func onNodeConnected(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
	node.ConnectMe = true

	return nil
}

// onIConnectNode - func event
func onIConnectNode(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
	jarvisbase.Debug("onIConnectNode")

	node.ConnectNode = true

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
func (n *jarvisNode) onMsgLocalRequesrNodes(ctx context.Context, msg *pb.JarvisMsg) error {
	// sendmsg := msg.GetMsg()
	jarvisbase.Debug("jarvisNode.onMsgLocalRequesrNodes")

	for _, v := range n.coredb.mapNodes {
		jarvisbase.Debug("jarvisNode.onMsgLocalRequesrNodes", jarvisbase.JSON("node", v))

		if n.mgrClient2.isConnected(v.Addr) {
			sendmsg, err := BuildRequestNodes(n.coredb.privKey, 0, n.myinfo.Addr, v.Addr)
			if err != nil {
				jarvisbase.Debug("jarvisNode.onMsgLocalRequesrNodes:BuildRequestNodes", zap.Error(err))

				continue
			}

			n.mgrClient2.addTask(sendmsg, "")
		}
	}

	return nil
}

// RequestCtrl - send ctrl to jarvisnode with addr
func (n *jarvisNode) RequestCtrl(ctx context.Context, addr string, ci *pb.CtrlInfo) error {
	sendmsg, err := BuildRequestCtrl(n.coredb.privKey, 0, n.myinfo.Addr, addr, ci)
	if err != nil {
		jarvisbase.Debug("jarvisNode.RequestCtrl", zap.Error(err))

		return err
	}

	msg, err := BuildLocalSendMsg(n.coredb.privKey, 0, n.myinfo.Addr, "", sendmsg)
	if err != nil {
		jarvisbase.Debug("jarvisNode.RequestCtrl:BuildLocalSendMsg", zap.Error(err))

		return err
	}

	n.mgrJasvisMsg.sendMsg(msg, nil, nil)

	return nil
}

// onTimerRequestNodes
func (n *jarvisNode) onTimerRequestNodes() error {
	jarvisbase.Debug("jarvisNode.onTimerRequestNodes")

	msg, err := BuildLocalRequestNodes(n.coredb.privKey, 0, n.myinfo.Addr, "")
	if err != nil {
		jarvisbase.Debug("jarvisNode.onTimerRequestNodes:BuildLocalRequestNodes", zap.Error(err))

		return err
	}

	n.mgrJasvisMsg.sendMsg(msg, nil, nil)

	return nil
}

// onMsgRequestNodes
func (n *jarvisNode) onMsgRequestNodes(ctx context.Context, msg *pb.JarvisMsg, stream pb.JarvisCoreServ_ProcMsgServer) error {
	jarvisbase.Debug("jarvisNode.onMsgRequestNodes")

	if stream == nil {
		jarvisbase.Debug("jarvisNode.onMsgRequestNodes", zap.Error(ErrStreamNil))

		return ErrStreamNil
	}

	for _, v := range n.coredb.mapNodes {
		mni := &pb.NodeBaseInfo{
			ServAddr: v.ServAddr,
			Addr:     v.Addr,
			Name:     v.Name,
		}

		jarvisbase.Debug("jarvisNode.onMsgRequestNodes", jarvisbase.JSON("node", mni))

		sendmsg, err := BuildNodeInfo(n.coredb.privKey, 0, n.myinfo.Addr, msg.SrcAddr, mni)
		if err != nil {
			jarvisbase.Debug("jarvisNode.onMsgRequestNodes:BuildNodeInfo", zap.Error(err))

			return err
		}

		err = stream.Send(sendmsg)
		if err != nil {
			jarvisbase.Debug("jarvisNode.onMsgRequestNodes:sendmsg", zap.Error(err))

			return err
		}
	}

	return nil
}

// FindNodeWithName - find node with name
func (n *jarvisNode) FindNodeWithName(name string) *coredbpb.NodeInfo {
	n.coredb.Lock()
	defer n.coredb.Unlock()

	for _, v := range n.coredb.mapNodes {
		if v.Name == name {
			return v
		}
	}

	return nil
}
