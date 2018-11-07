package main

import (
	"context"
	"fmt"

	"github.com/zhs007/jarviscore"
)

func main() {
	cfg, err := jarviscore.LoadConfig("config.yaml")
	if err != nil {
		fmt.Printf("load config.yaml err " + err.Error())
	}

	// cfg := jarviscore.Config{
	// 	DBPath:         "./dat",
	// 	LogPath:        "./log",
	// 	AnkaDBHttpServ: "127.0.0.1:8889",
	// 	AnkaDBEngine:   "leveldb",
	// 	DefPeerAddr:    "127.0.0.1:7788",
	// 	LogConsole:     true,
	// 	LogLevel:       "debug",
	// }

	myinfo := jarviscore.BaseInfo{
		Name:     cfg.BaseNodeInfo.NodeName,
		BindAddr: cfg.BaseNodeInfo.BindAddr,
		ServAddr: cfg.BaseNodeInfo.ServAddr,
	}

	jarviscore.InitJarvisCore(*cfg)
	defer jarviscore.ReleaseJarvisCore()

	node := jarviscore.NewNode(myinfo)
	// defer node.Stop()

	node.Start(context.Background())
}
