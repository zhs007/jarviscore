package jarviscore

import (
	"context"
	"io/ioutil"
	"time"

	"github.com/zhs007/jarviscore/base"
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
	RequestCtrl(ctx context.Context, addr string, ci *pb.CtrlInfo) error
	// SendFile - send filedata to jarvisnode with addr
	SendFile(ctx context.Context, addr string, fd *pb.FileData) error
	// RequestFile - request node send filedata to me
	RequestFile(ctx context.Context, addr string, rf *pb.RequestFile) error
	// RequestNodes - request nodes
	RequestNodes() error

	// AddCtrl2List - add ctrl msg to tasklist
	AddCtrl2List(addr string, ci *pb.CtrlInfo) error

	// AddNodeBaseInfo - add nodeinfo
	AddNodeBaseInfo(nbi *pb.NodeBaseInfo) error

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

	// SetNodeTypeInfo - set node type and version
	SetNodeTypeInfo(nodetype string, nodetypeversion string)

	// RegCtrl - register a ctrl
	RegCtrl(ctrltype string, ctrl Ctrl) error

	// PostMsg - like windows postMessage
	PostMsg(msg *pb.JarvisMsg, stream pb.JarvisCoreServ_ProcMsgServer, chanEnd chan int)
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
	jarvisbase.Info("jarviscore version is " + VERSION)

	if !IsValidNodeName(cfg.BaseNodeInfo.NodeName) {
		jarvisbase.Error("NewNode:IsValidNodeName", zap.Error(ErrInvalidNodeName))

		return nil
	}

	db, err := coredb.NewCoreDB(cfg.AnkaDB.DBPath, cfg.AnkaDB.HTTPServ, cfg.AnkaDB.Engine)
	if err != nil {
		jarvisbase.Error("NewNode:newCoreDB", zap.Error(err))

		return nil
	}

	node := &jarvisNode{
		myinfo: BaseInfo{
			Name:        cfg.BaseNodeInfo.NodeName,
			BindAddr:    cfg.BaseNodeInfo.BindAddr,
			ServAddr:    cfg.BaseNodeInfo.ServAddr,
			CoreVersion: VERSION,
		},
		coredb: db,
		cfg:    cfg,
		mgrCtrl: &ctrlMgr{
			mapCtrl: make(map[string](Ctrl)),
		},
	}

	node.mgrCtrl.Reg(CtrlTypeShell, &CtrlShell{})
	node.mgrCtrl.Reg(CtrlTypeScriptFile, &CtrlScriptFile{})
	node.mgrCtrl.Reg(CtrlTypeScriptFile2, &CtrlScriptFile2{})

	// event
	node.mgrEvent = newEventMgr(node)
	node.mgrEvent.regNodeEventFunc(EventOnNodeConnected, onNodeConnected)
	node.mgrEvent.regNodeEventFunc(EventOnIConnectNode, onIConnectNode)
	node.mgrEvent.regNodeEventFunc(EventOnDeprecateNode, onDeprecateNode)

	err = node.coredb.LoadPrivateKeyEx()
	if err != nil {
		jarvisbase.Error("NewNode:loadPrivateKey", zap.Error(err))

		return nil
	}

	node.coredb.LoadAllNodes()

	node.myinfo.Addr = node.coredb.GetPrivateKey().ToAddress()
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
func (n *jarvisNode) GetCoreDB() *coredb.CoreDB {
	return n.coredb
}

// sendCtrl - send ctrl to jarvisnode with addr
func (n *jarvisNode) sendCtrl(ctx context.Context, addr string, ci *pb.CtrlInfo) error {
	msg, err := BuildRequestCtrl(n.coredb.GetPrivateKey(), 0, n.myinfo.Addr, addr, ci)
	if err != nil {
		jarvisbase.Warn("jarvisNode.SendCtrl", zap.Error(err))

		return err
	}

	n.mgrClient2.addTask(msg, "", nil)

	// jarvisbase.Debug("jarvisNode.SendCtrl", jarvisbase.JSON("msg", msg))

	return nil
	// return n.requestCtrl(ctx, addr, ctrltype, []byte(command))
}

// OnMsg - proc JarvisMsg
func (n *jarvisNode) OnMsg(ctx context.Context, msg *pb.JarvisMsg, stream pb.JarvisCoreServ_ProcMsgServer) error {
	jarvisbase.Debug("jarvisNode.OnMsg", jarvisbase.JSON("msg", msg))

	// is timeout
	if IsTimeOut(msg) {
		jarvisbase.Warn("jarvisNode.OnMsg", zap.Error(ErrJarvisMsgTimeOut))

		return nil
	}

	// proc local msg
	if msg.MsgType == pb.MSGTYPE_LOCAL_CONNECT_OTHER ||
		msg.MsgType == pb.MSGTYPE_LOCAL_SENDMSG ||
		msg.MsgType == pb.MSGTYPE_LOCAL_REQUEST_NODES {

		// verify msg
		err := VerifyJarvisMsg(msg)
		if err != nil {
			jarvisbase.Warn("jarvisNode.OnMsg", zap.Error(err))

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
			jarvisbase.Warn("jarvisNode.OnMsg", zap.Error(err))

			return nil
		}

		return n.onMsgConnectNode(ctx, msg, stream)
	}

	// if is not my msg, broadcast msg
	if n.myinfo.Addr != msg.DestAddr {
		n.mgrClient2.addTask(msg, "", nil)
	} else {
		// verify msg
		err := VerifyJarvisMsg(msg)
		if err != nil {
			jarvisbase.Warn("jarvisNode.OnMsg", zap.Error(err))

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
		} else if msg.MsgType == pb.MSGTYPE_TRANSFER_FILE {
			return n.onMsgTransferFile(ctx, msg, stream)
		} else if msg.MsgType == pb.MSGTYPE_REQUEST_FILE {
			return n.onMsgRequestFile(ctx, msg, stream)
		} else if msg.MsgType == pb.MSGTYPE_REPLY_REQUEST_FILE {
			return n.onMsgReplyRequestFile(ctx, msg)
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

	cn := n.coredb.FindNodeWithServAddr(ci.ServAddr)
	if cn == nil {
		n.mgrClient2.addTask(nil, ci.ServAddr, nil)

		return nil
	}

	// if is me, return
	if cn.Addr == n.myinfo.Addr {
		return nil
	}

	if !cn.ConnectNode {
		n.mgrClient2.addTask(nil, cn.ServAddr, cn)

		return nil
	}

	return nil
}

// onMsgNodeInfo
func (n *jarvisNode) onMsgNodeInfo(ctx context.Context, msg *pb.JarvisMsg) error {
	ni := msg.GetNodeInfo()
	return n.AddNodeBaseInfo(ni)
	// cn := n.coredb.GetNode(ni.Addr)
	// if cn == nil {
	// 	err := n.coredb.UpdNodeBaseInfo(ni)
	// 	if err != nil {
	// 		jarvisbase.Warn("jarvisNode.onMsgNodeInfo:UpdNodeBaseInfo", zap.Error(err))

	// 		return err
	// 	}

	// 	n.mgrClient2.addTask(nil, ni.ServAddr, n.coredb.GetNode(ni.Addr))

	// 	return nil
	// } else if !cn.ConnectNode {
	// 	n.mgrClient2.addTask(nil, ni.ServAddr, cn)

	// 	return nil
	// }

	// return nil
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

	sendmsg, err := BuildReplyConn(n.coredb.GetPrivateKey(), 0, n.myinfo.Addr, ci.MyInfo.Addr, mni)
	if err != nil {
		jarvisbase.Warn("jarvisNode.onMsgConnectNode:BuildReplyConn", zap.Error(err))

		return err
	}
	// SignJarvisMsg(n.coredb.privKey, sendmsg)
	// jarvisbase.Debug("jarvisNode.onMsgConnectNode:sendmsg", jarvisbase.JSON("msg", sendmsg))
	err = stream.Send(sendmsg)
	if err != nil {
		jarvisbase.Warn("jarvisNode.onMsgConnectNode:sendmsg", zap.Error(err))

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
	}

	n.mgrEvent.onNodeEvent(ctx, EventOnNodeConnected, cn)
	// if !cn.ConnectMe {
	// 	n.mgrEvent.onNodeEvent(ctx, EventOnNodeConnected, cn)

	// 	n.coredb.UpdNodeInfo(cn.Addr)
	// } else {
	// 	n.mgrEvent.onNodeEvent(ctx, EventOnNodeConnected, cn)
	// }
	// cn.ConnectMe = true
	// n.mgrEvent.onNodeEvent(ctx, EventOnNodeConnected, cn)

	// msg, err := BuildLocalConnectOther(n.coredb.GetPrivateKey(), 0, n.myinfo.Addr, ci.MyInfo.Addr,
	// 	ci.MyInfo.ServAddr, ci.MyInfo)
	// if err != nil {
	// 	jarvisbase.Warn("jarvisNode.onMsgConnectNode:BuildLocalConnectOther", zap.Error(err))

	// 	return err
	// }
	// // SignJarvisMsg(n.coredb.privKey, msg)
	// n.PostMsg(msg, nil, nil)
	// } else if !cn.ConnectNode {
	// 	msg, err := BuildLocalConnectOther(n.coredb.GetPrivateKey(), 0, n.myinfo.Addr, ci.MyInfo.Addr,
	// 		ci.MyInfo.ServAddr, ci.MyInfo)
	// 	if err != nil {
	// 		jarvisbase.Warn("jarvisNode.onMsgConnectNode:BuildLocalConnectOther", zap.Error(err))

	// 		return err
	// 	}
	// 	// SignJarvisMsg(n.coredb.privKey, msg)
	// 	n.PostMsg(msg, nil, nil)
	// }

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

	// if !cn.ConnectNode {
	// 	n.mgrEvent.onNodeEvent(ctx, EventOnIConnectNode, cn)

	// 	err := n.coredb.UpdNodeInfo(cn.Addr)
	// 	if err != nil {
	// 		jarvisbase.Warn("jarvisNode.onMsgReplyConnect:UpdNodeInfo", zap.Error(err))
	// 	}
	// } else {
	// n.mgrEvent.onNodeEvent(ctx, EventOnIConnectNode, cn)
	// }

	n.mgrEvent.onNodeEvent(ctx, EventOnIConnectNode, cn)

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

	n.coredb.ForEachMapNodes(func(key string, node *coredbpb.NodeInfo) error {
		n.connectNode(node.ServAddr)
		return nil
	})
	// for _, node := range n.coredb.mapNodes {
	// 	n.connectNode(node.ServAddr)

	// 	// msg := BuildLocalConnectOther(0, n.myinfo.Addr, node.Addr, nbi)
	// 	// SignJarvisMsg(n.coredb.privKey, msg)
	// 	// n.mgrJasvisMsg.sendMsg(msg, nil)
	// }

	return nil
}

// connectNode - connect node
func (n *jarvisNode) connectNode(servaddr string) error {
	nbi := &pb.NodeBaseInfo{
		ServAddr:        n.myinfo.ServAddr,
		Addr:            n.myinfo.Addr,
		Name:            n.myinfo.Name,
		NodeTypeVersion: n.myinfo.NodeTypeVersion,
		NodeType:        n.myinfo.NodeType,
		CoreVersion:     n.myinfo.CoreVersion,
	}

	msg, err := BuildLocalConnectOther(n.coredb.GetPrivateKey(), 0, n.myinfo.Addr, "", servaddr, nbi)
	if err != nil {
		jarvisbase.Warn("jarvisNode.connectNode:BuildLocalConnectOther", zap.Error(err))

		return err
	}

	// SignJarvisMsg(n.coredb.privKey, msg)
	n.PostMsg(msg, nil, nil)

	return nil
}

// onMsgRequestCtrl
func (n *jarvisNode) onMsgRequestCtrl(ctx context.Context, msg *pb.JarvisMsg) error {
	sendmsg, err := BuildReply(n.coredb.GetPrivateKey(), 0, n.myinfo.Addr, msg.SrcAddr, pb.REPLYTYPE_OK, "")
	if err != nil {
		jarvisbase.Warn("jarvisNode.onMsgRequestCtrl:BuildReply", zap.Error(err))

		return err
	}

	n.mgrClient2.addTask(sendmsg, "", nil)

	n.mgrEvent.onMsgEvent(ctx, EventOnCtrl, msg)

	ci := msg.GetCtrlInfo()
	ret, err := n.mgrCtrl.Run(ci)
	if err != nil {
		sendmsg2, err := BuildCtrlResult(n.coredb.GetPrivateKey(), 0, n.myinfo.Addr, msg.SrcAddr, ci.CtrlID, err.Error())
		if err != nil {
			jarvisbase.Warn("jarvisNode.onMsgRequestCtrl:BuildCtrlResult", zap.Error(err))

			return err
		}

		n.mgrClient2.addTask(sendmsg2, "", nil)

		return nil
	}

	sendmsg2, err := BuildCtrlResult(n.coredb.GetPrivateKey(), 0, n.myinfo.Addr, msg.SrcAddr, ci.CtrlID, string(ret))
	n.mgrClient2.addTask(sendmsg2, "", nil)

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

	n.mgrClient2.addTask(sendmsg, "", nil)

	return nil
}

// AddCtrl2List - add ctrl msg to tasklist
func (n *jarvisNode) AddCtrl2List(addr string, ci *pb.CtrlInfo) error {
	msg, err := BuildRequestCtrl(n.coredb.GetPrivateKey(), 0, n.myinfo.Addr, addr, ci)
	if err != nil {
		jarvisbase.Warn("jarvisNode.AddCtrl2List", zap.Error(err))

		return err
	}

	// SignJarvisMsg(n.coredb.privKey, msg)
	n.PostMsg(msg, nil, nil)

	return nil
}

// onNodeConnected - func event
func onNodeConnected(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {

	if !node.ConnectMe {
		node.ConnectMe = true

		err := jarvisnode.GetCoreDB().UpdNodeInfo(node.Addr)
		if err != nil {
			jarvisbase.Warn("jarvisNode.onNodeConnected:UpdNodeInfo", zap.Error(err))
		}
	}

	if !node.ConnectNode {
		msg, err := BuildLocalConnectOther(jarvisnode.GetCoreDB().GetPrivateKey(), 0,
			jarvisnode.GetMyInfo().Addr, node.Addr,
			node.ServAddr, GetNodeBaseInfo(node))
		if err != nil {
			jarvisbase.Warn("jarvisNode.onMsgConnectNode:BuildLocalConnectOther", zap.Error(err))

			return err
		}

		jarvisnode.PostMsg(msg, nil, nil)
	}

	return nil
}

// onIConnectNode - func event
func onIConnectNode(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
	jarvisbase.Debug("onIConnectNode")

	if !node.ConnectNode {
		node.ConnectNode = true

		err := jarvisnode.GetCoreDB().UpdNodeInfo(node.Addr)
		if err != nil {
			jarvisbase.Warn("jarvisNode.onIConnectNode:UpdNodeInfo", zap.Error(err))
		}
	}

	return nil
}

// onDeprecateNode - func event
func onDeprecateNode(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
	jarvisbase.Debug("onDeprecateNode")

	if !node.Deprecated {
		node.Deprecated = true

		err := jarvisnode.GetCoreDB().UpdNodeInfo(node.Addr)
		if err != nil {
			jarvisbase.Warn("jarvisNode.onDeprecateNode:UpdNodeInfo", zap.Error(err))
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
func (n *jarvisNode) onMsgLocalRequesrNodes(ctx context.Context, msg *pb.JarvisMsg) error {
	// sendmsg := msg.GetMsg()
	// jarvisbase.Debug("jarvisNode.onMsgLocalRequesrNodes")

	n.coredb.ForEachMapNodes(func(key string, v *coredbpb.NodeInfo) error {
		if !v.Deprecated && n.mgrClient2.isConnected(v.Addr) {
			sendmsg, err := BuildRequestNodes(n.coredb.GetPrivateKey(), 0, n.myinfo.Addr, v.Addr)
			if err != nil {
				jarvisbase.Warn("jarvisNode.onMsgLocalRequesrNodes:BuildRequestNodes", zap.Error(err))

				return nil
			}

			n.mgrEvent.onNodeEvent(ctx, EventOnRequestNode, v)

			n.mgrClient2.addTask(sendmsg, "", nil)
		}

		return nil
	})

	// for _, v := range n.coredb.mapNodes {
	// 	// jarvisbase.Debug("jarvisNode.onMsgLocalRequesrNodes", jarvisbase.JSON("node", v))

	// 	if n.mgrClient2.isConnected(v.Addr) {
	// 		sendmsg, err := BuildRequestNodes(n.coredb.privKey, 0, n.myinfo.Addr, v.Addr)
	// 		if err != nil {
	// 			jarvisbase.Warn("jarvisNode.onMsgLocalRequesrNodes:BuildRequestNodes", zap.Error(err))

	// 			continue
	// 		}

	// 		n.mgrClient2.addTask(sendmsg, "")
	// 	}
	// }

	return nil
}

// RequestCtrl - send ctrl to jarvisnode with addr
func (n *jarvisNode) RequestCtrl(ctx context.Context, addr string, ci *pb.CtrlInfo) error {
	sendmsg, err := BuildRequestCtrl(n.coredb.GetPrivateKey(), 0, n.myinfo.Addr, addr, ci)
	if err != nil {
		jarvisbase.Warn("jarvisNode.RequestCtrl", zap.Error(err))

		return err
	}

	msg, err := BuildLocalSendMsg(n.coredb.GetPrivateKey(), 0, n.myinfo.Addr, "", sendmsg)
	if err != nil {
		jarvisbase.Warn("jarvisNode.RequestCtrl:BuildLocalSendMsg", zap.Error(err))

		return err
	}

	n.PostMsg(msg, nil, nil)

	return nil
}

// SendFile - send filedata to jarvisnode with addr
func (n *jarvisNode) SendFile(ctx context.Context, addr string, fd *pb.FileData) error {
	sendmsg, err := BuildFileData(n.coredb.GetPrivateKey(), 0, n.myinfo.Addr, addr, fd)
	if err != nil {
		jarvisbase.Warn("jarvisNode.SendFile", zap.Error(err))

		return err
	}

	msg, err := BuildLocalSendMsg(n.coredb.GetPrivateKey(), 0, n.myinfo.Addr, "", sendmsg)
	if err != nil {
		jarvisbase.Warn("jarvisNode.SendFile:BuildLocalSendMsg", zap.Error(err))

		return err
	}

	n.PostMsg(msg, nil, nil)

	return nil
}

// onTimerRequestNodes
func (n *jarvisNode) onTimerRequestNodes() error {
	jarvisbase.Debug("jarvisNode.onTimerRequestNodes")

	return n.RequestNodes()
	// msg, err := BuildLocalRequestNodes(n.coredb.GetPrivateKey(), 0, n.myinfo.Addr, "")
	// if err != nil {
	// 	jarvisbase.Warn("jarvisNode.onTimerRequestNodes:BuildLocalRequestNodes", zap.Error(err))

	// 	return err
	// }

	// n.PostMsg(msg, nil, nil)

	// return nil
}

// RequestNodes - request nodes
func (n *jarvisNode) RequestNodes() error {
	msg, err := BuildLocalRequestNodes(n.coredb.GetPrivateKey(), 0, n.myinfo.Addr, "")
	if err != nil {
		jarvisbase.Warn("jarvisNode.onTimerRequestNodes:BuildLocalRequestNodes", zap.Error(err))

		return err
	}

	n.PostMsg(msg, nil, nil)

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
		mni := &pb.NodeBaseInfo{
			ServAddr: v.ServAddr,
			Addr:     v.Addr,
			Name:     v.Name,
		}

		// jarvisbase.Debug("jarvisNode.onMsgRequestNodes", jarvisbase.JSON("node", mni))

		sendmsg, err := BuildNodeInfo(n.coredb.GetPrivateKey(), 0, n.myinfo.Addr, msg.SrcAddr, mni)
		if err != nil {
			jarvisbase.Warn("jarvisNode.onMsgRequestNodes:BuildNodeInfo", zap.Error(err))

			return err
		}

		err = stream.Send(sendmsg)
		if err != nil {
			jarvisbase.Warn("jarvisNode.onMsgRequestNodes:sendmsg", zap.Error(err))

			return err
		}

		return nil
	})

	// for _, v := range n.coredb.mapNodes {
	// 	mni := &pb.NodeBaseInfo{
	// 		ServAddr: v.ServAddr,
	// 		Addr:     v.Addr,
	// 		Name:     v.Name,
	// 	}

	// 	// jarvisbase.Debug("jarvisNode.onMsgRequestNodes", jarvisbase.JSON("node", mni))

	// 	sendmsg, err := BuildNodeInfo(n.coredb.privKey, 0, n.myinfo.Addr, msg.SrcAddr, mni)
	// 	if err != nil {
	// 		jarvisbase.Warn("jarvisNode.onMsgRequestNodes:BuildNodeInfo", zap.Error(err))

	// 		return err
	// 	}

	// 	err = stream.Send(sendmsg)
	// 	if err != nil {
	// 		jarvisbase.Warn("jarvisNode.onMsgRequestNodes:sendmsg", zap.Error(err))

	// 		return err
	// 	}
	// }

	return nil
}

// FindNodeWithName - find node with name
func (n *jarvisNode) FindNodeWithName(name string) *coredbpb.NodeInfo {
	n.coredb.Lock()
	defer n.coredb.Unlock()

	return n.coredb.FindMapNode(name)
	// for _, v := range n.coredb.mapNodes {
	// 	if v.Name == name {
	// 		return v
	// 	}
	// }

	// return nil
}

// replyStream
func (n *jarvisNode) replyStream(msg *pb.JarvisMsg, stream pb.JarvisCoreServ_ProcMsgServer,
	rt pb.REPLYTYPE, strErr string) error {

	sendmsg, err := BuildReply(n.coredb.GetPrivateKey(), 0, n.myinfo.Addr, msg.SrcAddr, rt, strErr)
	if err != nil {
		jarvisbase.Warn("jarvisNode.replyStream:BuildReply", zap.Error(err))

		return err
	}

	err = stream.SendMsg(sendmsg)
	if err != nil {
		jarvisbase.Warn("jarvisNode.replyStream:SendMsg", zap.Error(err))

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
		err1 := n.replyStream(msg, stream, pb.REPLYTYPE_ERROR, err.Error())
		if err1 != nil {
			jarvisbase.Warn("jarvisNode.onMsgTransferFile:replyStream err", zap.Error(err1))

			return err1
		}

		return err
	}

	// f, err := os.Create(fd.Filename)
	// if err != nil {
	// 	err1 := n.replyStream(msg, stream, pb.REPLYTYPE_ERROR, err.Error())
	// 	if err1 != nil {
	// 		jarvisbase.Warn("jarvisNode.onMsgTransferFile:replyStream err", zap.Error(err1))

	// 		return err1
	// 	}

	// 	return err
	// }

	// defer f.Close()

	// f.Write(fd.File)
	// f.Close()

	err = n.replyStream(msg, stream, pb.REPLYTYPE_OK, "")
	if err != nil {
		jarvisbase.Warn("jarvisNode.onMsgTransferFile:replyStream", zap.Error(err))

		return err
	}

	// sendmsg, err := BuildReply(n.coredb.privKey, 0, n.myinfo.Addr, msg.SrcAddr, pb.REPLYTYPE_OK, "")
	// if err != nil {
	// 	jarvisbase.Debug("jarvisNode.onMsgTransferFile:BuildReply", zap.Error(err))

	// 	return err
	// }

	// err = stream.SendMsg(sendmsg)
	// if err != nil {
	// 	jarvisbase.Debug("jarvisNode.onMsgTransferFile:SendMsg", zap.Error(err))

	// 	return err
	// }

	// for _, v := range n.coredb.mapNodes {
	// 	mni := &pb.NodeBaseInfo{
	// 		ServAddr: v.ServAddr,
	// 		Addr:     v.Addr,
	// 		Name:     v.Name,
	// 	}

	// 	jarvisbase.Debug("jarvisNode.onMsgRequestNodes", jarvisbase.JSON("node", mni))

	// 	sendmsg, err := BuildNodeInfo(n.coredb.privKey, 0, n.myinfo.Addr, msg.SrcAddr, mni)
	// 	if err != nil {
	// 		jarvisbase.Debug("jarvisNode.onMsgRequestNodes:BuildNodeInfo", zap.Error(err))

	// 		return err
	// 	}

	// 	err = stream.Send(sendmsg)
	// 	if err != nil {
	// 		jarvisbase.Debug("jarvisNode.onMsgRequestNodes:sendmsg", zap.Error(err))

	// 		return err
	// 	}
	// }

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

	rf := msg.GetRequestFile()

	buf, err := ioutil.ReadFile(rf.Filename)
	if err != nil {
		jarvisbase.Warn("jarvisNode.onMsgRequestFile:ReadFile", zap.Error(err))

		return err
	}

	fd := &pb.FileData{
		File:     buf,
		Filename: rf.Filename,
	}

	sendmsg, err := BuildReplyRequestFile(n.coredb.GetPrivateKey(), 0, n.myinfo.Addr, msg.SrcAddr, fd)
	if err != nil {
		jarvisbase.Warn("jarvisNode.onMsgRequestFile:BuildReplyRequestFile", zap.Error(err))

		return err
	}

	err = stream.Send(sendmsg)
	if err != nil {
		jarvisbase.Warn("jarvisNode.onMsgRequestFile:Send", zap.Error(err))

		return err
	}

	return nil
}

// onMsgReplyRequestFile
func (n *jarvisNode) onMsgReplyRequestFile(ctx context.Context, msg *pb.JarvisMsg) error {
	// jarvisbase.Debug("jarvisNode.onMsgReplyRequestFile")

	n.mgrEvent.onMsgEvent(ctx, EventOnReplyRequestFile, msg)

	return nil
}

// RequestFile - request node send filedata to me
func (n *jarvisNode) RequestFile(ctx context.Context, addr string, rf *pb.RequestFile) error {
	sendmsg, err := BuildRequestFile(n.coredb.GetPrivateKey(), 0, n.myinfo.Addr, addr, rf)
	if err != nil {
		jarvisbase.Warn("jarvisNode.RequestFile", zap.Error(err))

		return err
	}

	msg, err := BuildLocalSendMsg(n.coredb.GetPrivateKey(), 0, n.myinfo.Addr, "", sendmsg)
	if err != nil {
		jarvisbase.Warn("jarvisNode.RequestFile:BuildLocalSendMsg", zap.Error(err))

		return err
	}

	n.PostMsg(msg, nil, nil)

	return nil
}

// RegCtrl - register a ctrl
func (n *jarvisNode) RegCtrl(ctrltype string, ctrl Ctrl) error {
	n.mgrCtrl.Reg(ctrltype, ctrl)

	return nil
}

// PostMsg - like windows postMessage
func (n *jarvisNode) PostMsg(msg *pb.JarvisMsg, stream pb.JarvisCoreServ_ProcMsgServer, chanEnd chan int) {
	n.mgrJasvisMsg.sendMsg(msg, stream, chanEnd)
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

		n.mgrClient2.addTask(nil, nbi.ServAddr, n.coredb.GetNode(nbi.Addr))

		return nil
	} else if !cn.ConnectNode {
		n.mgrClient2.addTask(nil, nbi.ServAddr, cn)

		return nil
	}

	return nil
}
