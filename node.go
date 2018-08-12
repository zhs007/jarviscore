package jarviscore

import (
	"crypto/aes"
	"os"
	"os/signal"

	"github.com/seehuhn/fortuna"
	"github.com/zhs007/jarviscore/log"
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
	client      jarvisClient
	serv        *jarvisServer
	gen         *fortuna.Generator
	lstother    []*NodeInfo
	peeraddrmgr *peeraddrmgr
	signalchan  chan os.Signal
	// wg          sync.WaitGroup
}

const (
	nodeinfoCacheSize       = 32
	tokenLen                = 32
	randomMax         int64 = 0x7fffffffffffffff
	letterBytes             = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterBytesLen          = int64(len(letterBytes))
)

// NewNode -
func NewNode(baseinfo BaseInfo) JarvisNode {
	node := &jarvisNode{lstother: make([]*NodeInfo, nodeinfoCacheSize), signalchan: make(chan os.Signal, 1)}
	signal.Notify(node.signalchan, os.Interrupt, os.Kill)

	node.setMyInfo(baseinfo.ServAddr, baseinfo.Name, baseinfo.Token)

	return node
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

// setMyInfo -
func (n *jarvisNode) setMyInfo(servaddr string, name string, token string) error {
	if token == "" {
		n.myinfo.Token = n.generatorToken()

		log.Info("generatorToken", zap.String("Token", n.myinfo.Token))
	}

	n.myinfo.ServAddr = servaddr
	n.myinfo.Name = name

	return nil
}

// StopWithSignal -
func (n *jarvisNode) StopWithSignal(signal string) error {
	log.Info("StopWithSignal", zap.String("signal", signal))

	n.Stop()

	return nil
}

// Stop -
func (n *jarvisNode) Stop() error {
	n.peeraddrmgr.savePeerAddrFile()

	return nil
}

// func (n *jarvisNode) waitSignal() {
// 	c := make(chan os.Signal, 1)
// 	signal.Notify(c, os.Interrupt, os.Kill)

// 	s := <-c
// 	log.Info("Signal", zap.String("signal", s.String()))
// }

// Start -
func (n *jarvisNode) Start() (err error) {
	n.peeraddrmgr, err = newPeerAddrMgr(config.PeerAddrFile, config.DefPeerAddr)
	if err != nil {
		return err
	}

	log.Info("StartServer", zap.String("ServAddr", n.myinfo.ServAddr))
	n.serv, err = newServer(n.myinfo.ServAddr)
	if err != nil {
		return err
	}

	go n.serv.Start()

	select {
	case signal := <-n.signalchan:
		n.StopWithSignal(signal.String())
	case <-n.serv.servchan:
		log.Info("ServEnd")
	}

	return nil
}
