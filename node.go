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
	SendCtrl(ctx context.Context, addr string, ctrltype string, command string) error

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

// func (n *jarvisNode) savePrivateKey() error {
// 	privkey := &pb.PrivateKey{
// 		PriKey: n.privKey.ToPrivateBytes(),
// 	}

// 	data, err := proto.Marshal(privkey)
// 	if err != nil {
// 		return err
// 	}

// 	err = n.coredb.Put([]byte(coredbMyPrivKey), data)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (n *jarvisNode) loadPrivateKey() error {
// 	dat, err := n.coredb.Get([]byte(coredbMyPrivKey))
// 	if err != nil {
// 		n.privKey = jarviscrypto.GenerateKey()

// 		return n.savePrivateKey()
// 	}

// 	pbprivkey := &pb.PrivateKey{}
// 	err = proto.Unmarshal(dat, pbprivkey)
// 	if err != nil {
// 		n.privKey = jarviscrypto.GenerateKey()

// 		return n.savePrivateKey()
// 	}

// 	privkey := jarviscrypto.NewPrivateKey()
// 	err = privkey.FromBytes(pbprivkey.PriKey)
// 	if err != nil {
// 		n.privKey = jarviscrypto.GenerateKey()

// 		return n.savePrivateKey()
// 	}

// 	n.privKey = privkey

// 	return nil
// }

// // RandomInt64 -
// func (n *jarvisNode) RandomInt64(maxval int64) int64 {
// 	if n.gen == nil {
// 		n.gen = fortuna.NewGenerator(aes.NewCipher)
// 	}

// 	var rt = int64((randomMax / maxval) * maxval)
// 	var cr = n.gen.Int63()
// 	for cr >= rt {
// 		cr = n.gen.Int63()
// 	}

// 	return cr % maxval
// }

// // generatorToken -
// func (n *jarvisNode) generatorToken() string {
// 	b := make([]byte, tokenLen)
// 	for i := range b {
// 		b[i] = letterBytes[n.RandomInt64(letterBytesLen)]
// 	}

// 	return string(b)
// }

// // setMyInfo -
// func (n *jarvisNode) setMyInfo(servaddr string, bindaddr string, name string, token string) error {
// 	// if token == "" {
// 	// 	n.myinfo.Token = n.generatorToken()

// 	// 	log.Info("generatorToken", zap.String("Token", n.myinfo.Token))
// 	// }

// 	n.myinfo.ServAddr = servaddr
// 	n.myinfo.BindAddr = bindaddr
// 	n.myinfo.Name = name
// 	n.myinfo.Token = token

// 	return nil
// }

// // StopWithSignal -
// func (n *jarvisNode) StopWithSignal(signal string) error {
// 	jarvisbase.Info("StopWithSignal", zap.String("signal", signal))

// 	n.Stop()

// 	return nil
// }

// Stop -
func (n *jarvisNode) Stop() error {
	if n.serv2 != nil {
		n.serv2.Stop()
	}

	return nil
}

// func (n *jarvisNode) waitSignal() {
// 	c := make(chan os.Signal, 1)
// 	signal.Notify(c, os.Interrupt, os.Kill)

// 	s := <-c
// 	log.Info("Signal", zap.String("signal", s.String()))
// }

// func (n *jarvisNode) onStateChg() bool {
// 	if n.servstate == stateEnd && n.clientstate == stateEnd {
// 		return true
// 	}

// 	return false
// }

// func (n *jarvisNode) waitEnd() {
// 	for {
// 		select {
// 		case signal := <-n.signalchan:
// 			n.StopWithSignal(signal.String())
// 		case <-n.serv.servchan:
// 			jarvisbase.Info("ServEnd")
// 			n.servstate = stateEnd
// 			if n.onStateChg() {
// 				return
// 			}
// 		case <-n.client.clientchan:
// 			jarvisbase.Info("ClientEnd")
// 			n.clientstate = stateEnd
// 			if n.onStateChg() {
// 				return
// 			}
// 		case <-n.nodechan:
// 			jarvisbase.Info("SafeEnd")
// 		}
// 	}
// }

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

	for {
		select {
		case <-ctx.Done():
			n.Stop()
			return nil
		}
	}
}

// func (n *jarvisNode) hasNodeWithAddr(addr string) bool {
// 	if addr == n.myinfo.Addr {
// 		return true
// 	}

// 	return n.mgrNodeInfo.hasNodeInfo(addr)
// }

// // onAddNode
// func (n *jarvisNode) onAddNode(bi *BaseInfo) {
// 	if n.hasNodeToken(bi.Token) {
// 		return
// 	}

// 	if n.mgrpeeraddr.canConnect(bi.ServAddr) {
// 		go n.client.connect(bi.ServAddr)
// 	}
// }

// // addNode
// func (n *jarvisNode) addNode(bi *BaseInfo) {
// 	if n.hasNodeToken(bi.Token) {
// 		return
// 	}

// 	n.mgrNodeInfo.addNodeInfo(bi)
// }

// onNodeConnectMe
func (n *jarvisNode) onNodeConnectMe(bi *BaseInfo) {
	if bi.Addr == n.myinfo.Addr {
		return
	}

	// n.mgrNodeInfo.addNodeInfo(bi)
	// n.mgrNodeInfo.chg2ConnectMe(bi.Addr)

	// _, connNode := n.mgrNodeInfo.getNodeConnectState(bi.Addr)
	// if !connNode {
	// 	n.client.pushNewConnect(bi)
	// }
}

// onIConnectNode
func (n *jarvisNode) onIConnectNode(bi *BaseInfo) {
	if bi.Addr == n.myinfo.Addr {
		return
	}

	// n.mgrNodeInfo.addNodeInfo(bi)
	// n.mgrNodeInfo.chg2ConnectNode(bi.Addr)

	// n.serv.broadcastNode(bi)
}

// onGetNewNode
func (n *jarvisNode) onGetNewNode(bi *BaseInfo) {
	if bi.Addr == n.myinfo.Addr {
		return
	}

	// n.mgrNodeInfo.addNodeInfo(bi)
	// _, connNode := n.mgrNodeInfo.getNodeConnectState(bi.Addr)
	// if !connNode {
	// 	n.client.pushNewConnect(bi)
	// }
}

// requestCtrl
func (n *jarvisNode) requestCtrl(ctx context.Context, addr string, ctrltype string, command []byte) error {
	// ctrlid := n.mgrNodeInfo.getCtrlID(addr)
	// if ctrlid < 0 {
	// 	return ErrCoreDBNoAddr
	// }

	// buf := append([]byte(addr), command...)

	// r, s, err := n.coredb.privKey.Sign(buf)
	// if err != nil {
	// 	return ErrSign
	// }

	// // pk := jarviscrypto.NewPublicKey()
	// // pk.FromBytes(n.coredb.privKey.ToPublicBytes())
	// // if !pk.Verify(buf, r, s) {
	// // 	return ErrPublicKeyVerify
	// // }

	// ci := &pb.CtrlInfo{
	// 	Ctrlid:      ctrlid,
	// 	DestAddr:    addr,
	// 	SrcAddr:     n.myinfo.Addr,
	// 	CtrlType:    ctrltype,
	// 	Command:     command,
	// 	ForwordNums: 0,
	// 	SignR:       r.Bytes(),
	// 	SignS:       s.Bytes(),
	// 	PubKey:      n.coredb.privKey.ToPublicBytes(),
	// }

	// connMe, connNode := n.mgrNodeInfo.getNodeConnectState(addr)
	// if connMe {
	// 	if n.serv.sendCtrl(ci) == nil {
	// 		return nil
	// 	}
	// }

	// if connNode {
	// 	if n.client.sendCtrl(ctx, ci) == nil {
	// 		return nil
	// 	}
	// }

	return nil
}

// GetCoreDB - get coredb
func (n *jarvisNode) GetCoreDB() *CoreDB {
	return n.coredb
}

// SendCtrl - send ctrl to jarvisnode with addr
func (n *jarvisNode) SendCtrl(ctx context.Context, addr string, ctrltype string, command string) error {
	return n.requestCtrl(ctx, addr, ctrltype, []byte(command))
}

// OnMsg - proc JarvisMsg
func (n *jarvisNode) OnMsg(ctx context.Context, msg *pb.JarvisMsg, stream pb.JarvisCoreServ_ProcMsgServer) error {
	// is timeout
	if IsTimeOut(msg) {
		jarvisbase.Debug("jarvisNode.OnMsg", zap.Error(ErrJarvisMsgTimeOut))

		return nil
	}

	// if is not my msg, broadcast msg
	if n.coredb.addr != msg.DestAddr {
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
		} else if msg.MsgType == pb.MSGTYPE_CONNECT_NODE {
			return n.onMsgConnectNode(ctx, msg, stream)
		} else if msg.MsgType == pb.MSGTYPE_REPLY_CONNECT {
			return n.onMsgReplyConnect(ctx, msg)
		}
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

		return n.mgrClient2.connectNode(ctx, ni)
	} else if !cn.ConnectNode {
		return n.mgrClient2.connectNode(ctx, ni)
	}

	return nil
}

// onMsgConnectNode
func (n *jarvisNode) onMsgConnectNode(ctx context.Context, msg *pb.JarvisMsg, stream pb.JarvisCoreServ_ProcMsgServer) error {
	if stream == nil {
		jarvisbase.Debug("jarvisNode.onMsgConnectNode", zap.Error(ErrStreamNil))

		return ErrStreamNil
	}

	ni := msg.GetNodeInfo()

	mni := &pb.NodeBaseInfo{
		ServAddr: n.myinfo.ServAddr,
		Addr:     n.myinfo.Addr,
		Name:     n.myinfo.Name,
	}
	sendmsg := BuildReplyConn(0, n.myinfo.Addr, ni.Addr, mni)
	stream.Send(sendmsg)

	cn := n.coredb.getNode(ni.Addr)
	if cn == nil {
		err := n.coredb.insNode(ni)
		if err != nil {
			jarvisbase.Debug("jarvisNode.onMsgConnectNode:insNode", zap.Error(err))

			return err
		}

		return n.mgrClient2.connectNode(ctx, ni)
	} else if !cn.ConnectNode {
		return n.mgrClient2.connectNode(ctx, ni)
	}

	return nil
}

// onMsgReplyConnect
func (n *jarvisNode) onMsgReplyConnect(ctx context.Context, msg *pb.JarvisMsg) error {
	ni := msg.GetNodeInfo()
	cn := n.coredb.getNode(ni.Addr)
	if cn == nil {
		jarvisbase.Debug("jarvisNode.onMsgReplyConnect", zap.Error(ErrCoreDBHasNotNode))

		return ErrCoreDBHasNotNode
	}

	cn.ConnectNode = true

	n.coredb.updNodeBaseInfo(ni)

	return nil
}

// GetMyInfo - get my nodeinfo
func (n *jarvisNode) GetMyInfo() *BaseInfo {
	return &n.myinfo
}
