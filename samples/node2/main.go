package main

import (
	"context"

	"github.com/zhs007/jarviscore"
	pb "github.com/zhs007/jarviscore/proto"
)

func main() {
	cfg := jarviscore.Config{
		DBPath:         "./dat",
		LogPath:        "./log",
		AnkaDBHttpServ: "127.0.0.1:8889",
		AnkaDBEngine:   "leveldb",
		DefPeerAddr:    "127.0.0.1:7788",
		LogConsole:     true,
		LogLevel:       "debug",
	}

	myinfo := jarviscore.BaseInfo{
		Name:     "node002",
		BindAddr: ":7789",
		ServAddr: ":7789",
		NodeType: pb.NODETYPE_NORMAL,
	}

	jarviscore.InitJarvisCore(cfg)
	defer jarviscore.ReleaseJarvisCore()

	node := jarviscore.NewNode(myinfo)
	// defer node.Stop()

	node.Start(context.Background())
}
