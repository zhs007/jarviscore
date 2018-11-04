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
		AnkaDBHttpServ: "127.0.0.1:8880",
		AnkaDBEngine:   "leveldb",
		DefPeerAddr:    "jarvis.heyalgo.io:7788",
		LogConsole:     true,
		LogLevel:       "debug",
		LstTrustNode:   []string{"1JJaKpZGhYPuVHc1EKiiHZEswPAB5SybW5"},
	}

	myinfo := jarviscore.BaseInfo{
		Name:     "node001",
		BindAddr: ":7788",
		ServAddr: ":7788",
		NodeType: pb.NODETYPE_NORMAL,
	}

	jarviscore.InitJarvisCore(cfg)
	defer jarviscore.ReleaseJarvisCore()

	// pubip := jarviscore.GetHTTPPulicIP()
	// log.Debug(pubip)
	// ip := net.ParseIP(":7788")
	// log.Debug(ip.String())

	// ip1 := net.ParseIP("127.0.0.1:7788")
	// log.Debug(ip1.String())

	node := jarviscore.NewNode(myinfo)
	// defer node.Stop()

	node.Start(context.Background())
}
