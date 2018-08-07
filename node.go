package jarviscore

import (
	"crypto/aes"

	"github.com/seehuhn/fortuna"
)

// JarvisNode -
type JarvisNode interface {
	Start(servaddr string, name string, token string) (err error)
}

// jarvisNode -
type jarvisNode struct {
	myinfo   BaseInfo
	client   jarvisClient
	serv     *jarvisServer
	gen      *fortuna.Generator
	lstother []*NodeInfo
}

const (
	nodeinfoCacheSize       = 32
	tokenLen                = 32
	randomMax         int64 = 0x7fffffffffffffff
	letterBytes             = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterBytesLen          = int64(len(letterBytes))
)

// NewNode -
func NewNode() JarvisNode {
	return &jarvisNode{lstother: make([]*NodeInfo, nodeinfoCacheSize)}
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

// GeneratorToken -
func (n *jarvisNode) GeneratorToken() string {
	b := make([]byte, tokenLen)
	for i := range b {
		b[i] = letterBytes[n.RandomInt64(letterBytesLen)]
	}

	return string(b)
}

// SetMyInfo -
func (n *jarvisNode) SetMyInfo(servaddr string, name string, token string) error {
	if token == "" {
		n.myinfo.Token = n.GeneratorToken()
	}

	n.myinfo.ServAddr = servaddr
	n.myinfo.Name = name

	return nil
}

// Close -
func (n *jarvisNode) Close() error {

	return nil
}

// Start -
func (n *jarvisNode) Start(servaddr string, name string, token string) (err error) {
	n.SetMyInfo(servaddr, name, token)

	n.serv, err = newServer(servaddr)
	if err != nil {
		return err
	}

	go n.serv.Start()

	<-n.serv.servchan

	return nil
}
