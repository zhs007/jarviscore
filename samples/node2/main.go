package main

import (
	"github.com/zhs007/jarviscore"
	pb "github.com/zhs007/jarviscore/proto"
)

func main() {
	cfg := jarviscore.Config{
		RunPath:      "./",
		PeerAddrFile: "peeraddr.yaml",
		DefPeerAddr:  "127.0.0.1:7788",
	}

	myinfo := jarviscore.BaseInfo{
		Name:     "node002",
		BindAddr: ":7789",
		ServAddr: "127.0.0.1:7789",
		NodeType: pb.NODETYPE_NORMAL,
	}

	jarviscore.InitJarvisCore(cfg)
	defer jarviscore.ReleaseJarvisCore()

	node := jarviscore.NewNode(myinfo)
	defer node.Stop()

	node.Start()
}
