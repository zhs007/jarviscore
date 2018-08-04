package jarviscore

import (
	"crypto/aes"

	"github.com/seehuhn/fortuna"
)

// Node -
type Node struct {
	myinfo NodeInfo
	client Client
	serv   server
	gen    *fortuna.Generator
}

const tokenLen = 32
const randomMax int64 = 0x7fffffffffffffff
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const letterBytesLen = int64(len(letterBytes))

// NewNode -
func NewNode() *Node {
	return &Node{}
}

// RandomInt64 -
func (n *Node) RandomInt64(maxval int64) int64 {
	if n.gen == nil {
		n.gen = fortuna.NewGenerator(aes.NewCipher)
	}

	var rt = int64((randomMax / maxval) * maxval)
	var cr = n.gen.Int63()
	for cr >= rt {
		cr = n.gen.Int63()
	}

	return cr
}

// GeneratorToken -
func (n *Node) GeneratorToken() string {
	if n.gen == nil {
		n.gen = fortuna.NewGenerator(aes.NewCipher)
	}

	b := make([]byte, tokenLen)
	for i := range b {
		b[i] = letterBytes[n.RandomInt64(letterBytesLen)]
	}

	return string(b)
}

// SetMyInfo -
func (n *Node) SetMyInfo(servaddr string, name string, token string) error {
	if token == "" {
		n.myinfo.Token = n.GeneratorToken()
	}

	n.myinfo.ServAddr = servaddr
	n.myinfo.Name = name

	return nil
}
