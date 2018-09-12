package jarviscore

import (
	"crypto/aes"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/protobuf/proto"

	"github.com/seehuhn/fortuna"
	"github.com/zhs007/jarviscore/crypto"
	"github.com/zhs007/jarviscore/db"
	"github.com/zhs007/jarviscore/log"
	pb "github.com/zhs007/jarviscore/proto"
	"go.uber.org/zap"
)

// JarvisNode -
type JarvisNode interface {
	Start() (err error)
	Stop() (err error)
}

// jarvisNode -
type jarvisNode struct {
	myinfo      BaseInfo
	client      *jarvisClient
	serv        *jarvisServer
	gen         *fortuna.Generator
	mgrNodeInfo *nodeInfoMgr
	// mgrpeeraddr *peerAddrMgr
	mgrNodeCtrl *nodeCtrlMgr
	signalchan  chan os.Signal
	servstate   int
	clientstate int
	nodechan    chan int
	coredb      jarvisdb.Database
	privKey     *jarviscrypto.PrivateKey
	// wg          sync.WaitGroup
}

const (
	nodeinfoCacheSize       = 32
	tokenLen                = 32
	randomMax         int64 = 0x7fffffffffffffff
	letterBytes             = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterBytesLen          = int64(len(letterBytes))
	stateNormal             = 0
	stateStart              = 1
	stateEnd                = 2
	coredbMyPrivKey         = "myprivkey"
)

// NewNode -
func NewNode(baseinfo BaseInfo) JarvisNode {
	db, err := jarvisdb.NewJarvisLDB(getRealPath("coredb"), 16, 16)
	if err != nil {
		errorLog("NewNode:NewJarvisLDB", err)
		return nil
	}

	// db.Get([]byte(coredbMyPrivKey))

	node := &jarvisNode{
		myinfo:      baseinfo,
		signalchan:  make(chan os.Signal, 1),
		mgrNodeCtrl: newNodeCtrlMgr(),
		coredb:      db,
	}

	err = node.loadPrivateKey()
	if err != nil {
		errorLog("NewNode:loadPrivateKey", err)

		return nil
	}

	node.mgrNodeInfo = newNodeInfoMgr(node)
	node.mgrNodeInfo.loadFromDB()

	// node.myinfo = baseinfo
	node.myinfo.Addr = node.privKey.ToAddress()

	// signal.Notify(node.signalchan)
	signal.Notify(node.signalchan, os.Interrupt, os.Kill, syscall.SIGSTOP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGTSTP)

	// var token string
	// var prikey *privatekey

	// prikey, err = loadPrivateKeyFile()
	// if err != nil {
	// 	token = node.generatorToken()

	// 	log.Info("generatorToken", zap.String("Token", token))

	// 	prikey = &privatekey{Token: token}
	// 	savePrivateKeyFile(prikey)
	// }

	// token = prikey.Token
	// node.setMyInfo(baseinfo.ServAddr, baseinfo.BindAddr, baseinfo.Name, token)

	return node
}

func (n *jarvisNode) savePrivateKey() error {
	privkey := &pb.PrivateKey{
		PriKey: n.privKey.ToPrivateBytes(),
	}

	data, err := proto.Marshal(privkey)
	if err != nil {
		return err
	}

	err = n.coredb.Put([]byte(coredbMyPrivKey), data)
	if err != nil {
		return err
	}

	return nil
}

func (n *jarvisNode) loadPrivateKey() error {
	dat, err := n.coredb.Get([]byte(coredbMyPrivKey))
	if err != nil {
		n.privKey = jarviscrypto.GenerateKey()

		return n.savePrivateKey()
	}

	pbprivkey := &pb.PrivateKey{}
	err = proto.Unmarshal(dat, pbprivkey)
	if err != nil {
		n.privKey = jarviscrypto.GenerateKey()

		return n.savePrivateKey()
	}

	privkey := jarviscrypto.NewPrivateKey()
	err = privkey.FromBytes(pbprivkey.PriKey)
	if err != nil {
		n.privKey = jarviscrypto.GenerateKey()

		return n.savePrivateKey()
	}

	n.privKey = privkey

	return nil
}

// RandomInt64 -
func (n *jarvisNode) RandomInt64(maxval int64) int64 {
	if n.gen == nil {
		n.gen = fortuna.NewGenerator(aes.NewCipher)
	}

	var rt = int64((randomMax / maxval) * maxval)
	var cr = n.gen.Int63()
	for cr >= rt {
		cr = n.gen.Int63()
	}

	return cr % maxval
}

// generatorToken -
func (n *jarvisNode) generatorToken() string {
	b := make([]byte, tokenLen)
	for i := range b {
		b[i] = letterBytes[n.RandomInt64(letterBytesLen)]
	}

	return string(b)
}

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

// StopWithSignal -
func (n *jarvisNode) StopWithSignal(signal string) error {
	log.Info("StopWithSignal", zap.String("signal", signal))

	n.Stop()

	return nil
}

// Stop -
func (n *jarvisNode) Stop() error {
	if n.serv != nil {
		n.serv.Stop()
	}

	if n.client != nil {
		n.client.Stop()
	}

	// n.mgrpeeraddr.savePeerAddrFile()
	n.mgrNodeCtrl.save()

	return nil
}

// func (n *jarvisNode) waitSignal() {
// 	c := make(chan os.Signal, 1)
// 	signal.Notify(c, os.Interrupt, os.Kill)

// 	s := <-c
// 	log.Info("Signal", zap.String("signal", s.String()))
// }

func (n *jarvisNode) onStateChg() bool {
	if n.servstate == stateEnd && n.clientstate == stateEnd {
		return true
	}

	return false
}

func (n *jarvisNode) waitEnd() {
	for {
		select {
		case signal := <-n.signalchan:
			n.StopWithSignal(signal.String())
		case <-n.serv.servchan:
			log.Info("ServEnd")
			n.servstate = stateEnd
			if n.onStateChg() {
				return
			}
		case <-n.client.clientchan:
			log.Info("ClientEnd")
			n.clientstate = stateEnd
			if n.onStateChg() {
				return
			}
		case <-n.nodechan:
			log.Info("SafeEnd")
		}
	}
}

// Start -
func (n *jarvisNode) Start() (err error) {
	// n.mgrpeeraddr, err = newPeerAddrMgr(config.DefPeerAddr)
	// if err != nil {
	// 	return err
	// }

	log.Info("StartServer", zap.String("ServAddr", n.myinfo.ServAddr))
	n.serv, err = newServer(n)
	if err != nil {
		return err
	}

	n.client = newClient(n)

	go n.serv.Start()
	go n.client.Start()

	n.waitEnd()

	return nil
}

func (n *jarvisNode) hasNodeWithAddr(addr string) bool {
	if addr == n.myinfo.Addr {
		return true
	}

	return n.mgrNodeInfo.hasNodeInfo(addr)
}

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

	n.mgrNodeInfo.addNodeInfo(bi, false)
	n.mgrNodeInfo.chg2ConnectMe(bi.Addr)

	_, connNode := n.mgrNodeInfo.getNodeConnectState(bi.Addr)
	if !connNode {
		n.client.pushNewConnect(bi)
	}
}

// onIConnectNode
func (n *jarvisNode) onIConnectNode(bi *BaseInfo) {
	if bi.Addr == n.myinfo.Addr {
		return
	}

	n.mgrNodeInfo.addNodeInfo(bi, false)
	n.mgrNodeInfo.chg2ConnectNode(bi.Addr)

	n.serv.broadcastNode(bi)
}

// onGetNewNode
func (n *jarvisNode) onGetNewNode(bi *BaseInfo) {
	if bi.Addr == n.myinfo.Addr {
		return
	}

	n.mgrNodeInfo.addNodeInfo(bi, false)
	_, connNode := n.mgrNodeInfo.getNodeConnectState(bi.Addr)
	if !connNode {
		n.client.pushNewConnect(bi)
	}
}
