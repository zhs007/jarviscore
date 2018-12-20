package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/zhs007/jarviscore/base"
	"go.uber.org/zap"

	"github.com/zhs007/jarviscore"
	"github.com/zhs007/jarviscore/coredb/proto"
	pb "github.com/zhs007/jarviscore/proto"
)

func main() {
	cfg, err := jarviscore.LoadConfig("config.yaml")
	if err != nil {
		fmt.Printf("load config.yaml err " + err.Error())
	}
	// cfg := jarviscore.Config{
	// 	DBPath:         "./dat",
	// 	LogPath:        "./log",
	// 	AnkaDBHttpServ: "127.0.0.1:8880",
	// 	AnkaDBEngine:   "leveldb",
	// 	DefPeerAddr:    "jarvis.heyalgo.io:7788",
	// 	LogConsole:     true,
	// 	LogLevel:       "debug",
	// 	LstTrustNode:   []string{"1JJaKpZGhYPuVHc1EKiiHZEswPAB5SybW5"},
	// }

	// myinfo := jarviscore.BaseInfo{
	// 	Name:     cfg.BaseNodeInfo.NodeName,
	// 	BindAddr: cfg.BaseNodeInfo.BindAddr,
	// 	ServAddr: cfg.BaseNodeInfo.ServAddr,
	// }

	jarviscore.InitJarvisCore(cfg)
	defer jarviscore.ReleaseJarvisCore()

	// pubip := jarviscore.GetHTTPPulicIP()
	// log.Debug(pubip)
	// ip := net.ParseIP(":7788")
	// log.Debug(ip.String())

	// ip1 := net.ParseIP("127.0.0.1:7788")
	// log.Debug(ip1.String())

	node := jarviscore.NewNode(cfg)
	node.RegNodeEventFunc(jarviscore.EventOnIConnectNode, onIConnectNode)
	node.RegNodeEventFunc(jarviscore.EventOnNodeConnected, onNodeConnected)
	node.RegMsgEventFunc(jarviscore.EventOnCtrl, onCtrl)
	node.RegMsgEventFunc(jarviscore.EventOnCtrlResult, onCtrlResult)
	// defer node.Stop()

	node.Start(context.Background())
}

func sendCtrl(ctx context.Context, jarvisnode jarviscore.JarvisNode, node *coredbpb.NodeInfo) error {
	dat, err := ioutil.ReadFile("./test/test.sh")
	if err != nil {
		jarvisbase.Warn("load script file", zap.Error(err))

		return err
	}

	ci, err := jarviscore.BuildCtrlInfoForScriptFile(1, "test.sh", dat, "")
	if err != nil {
		jarvisbase.Warn("BuildCtrlInfoForScriptFile", zap.Error(err))

		return err
	}

	err = jarvisnode.RequestCtrl(ctx, node.Addr, ci)
	if err != nil {
		jarvisbase.Warn("BuildCtrlInfoForScriptFile", zap.Error(err))

		return err
	}

	return nil
}

// onIConnectNode - func event
func onIConnectNode(ctx context.Context, jarvisnode jarviscore.JarvisNode, node *coredbpb.NodeInfo) error {
	return sendCtrl(ctx, jarvisnode, node)
}

// onNodeConnected - func event
func onNodeConnected(ctx context.Context, jarvisnode jarviscore.JarvisNode, node *coredbpb.NodeInfo) error {
	if jarvisnode.IsConnected(node.Addr) {
		return sendCtrl(ctx, jarvisnode, node)
	}

	return nil
}

func onCtrl(ctx context.Context, jarvisnode jarviscore.JarvisNode, msg *pb.JarvisMsg) error {
	jarvisbase.Info("onCtrl")

	return nil
}

func onCtrlResult(ctx context.Context, jarvisnode jarviscore.JarvisNode, msg *pb.JarvisMsg) error {
	jarvisbase.Info("onCtrlResult")

	return nil
}
