package jarviscore

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/zhs007/jarviscore/coredb/proto"
	"github.com/zhs007/jarviscore/proto"
)

func TestCheckPrivateKey(t *testing.T) {
	cfg, err := LoadConfig("./test/node1.yaml")
	if err != nil {
		t.Fatalf("TestCheckPrivateKey load config %v", err)

		return
	}

	InitJarvisCore(cfg)
	defer ReleaseJarvisCore()

	node1, err := NewNode(cfg)
	if err != nil {
		t.Fatalf("TestCheckPrivateKey NewNode node1 %v", err)

		return
	}

	addr1 := node1.GetCoreDB().GetPrivateKey().ToAddress()
	node1.GetCoreDB().Close()

	node2, err := NewNode(cfg)
	if err != nil {
		t.Fatalf("TestCheckPrivateKey NewNode node2 %v", err)

		return
	}

	addr2 := node2.GetCoreDB().GetPrivateKey().ToAddress()
	node2.GetCoreDB().Close()

	if addr1 != addr2 {
		t.Fatalf("TestCheckPrivateKey addr fail")
	}

	t.Logf("TestCheckPrivateKey OK")
}

func TestCheckNode(t *testing.T) {
	cfg1, err := LoadConfig("./test/node1.yaml")
	if err != nil {
		t.Fatalf("TestCheckNode load config %v", err)
	}

	cfg2, err := LoadConfig("./test/node2.yaml")
	if err != nil {
		t.Fatalf("TestCheckNode load config %v", err)
	}

	InitJarvisCore(cfg1)
	defer ReleaseJarvisCore()

	node1, err := NewNode(cfg1)
	if err != nil {
		t.Fatalf("TestCheckNode NewNode node1 %v", err)

		return
	}

	addr1 := node1.GetCoreDB().GetPrivateKey().ToAddress()

	node2, err := NewNode(cfg2)
	if err != nil {
		t.Fatalf("TestCheckNode NewNode node2 %v", err)

		return
	}

	addr2 := node2.GetCoreDB().GetPrivateKey().ToAddress()

	cp := 0

	mapICN1 := make(map[string]string)
	mapNC1 := make(map[string]string)
	mapICN2 := make(map[string]string)
	mapNC2 := make(map[string]string)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	errstr := ""

	node1.RegNodeEventFunc(EventOnIConnectNode,
		func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
			if addr2 != node.Addr {
				errstr = "TestCheckNode node1 EventOnIConnectNode addr fail"
				cancel()

				return nil
			}

			if node.ConnType != coredbpb.CONNECTTYPE_DIRECT_CONN {
				errstr = fmt.Sprintf("TestCheckNode node1 EventOnIConnectNode node.ConnType %v %v",
					node.ConnType, coredbpb.CONNECTTYPE_DIRECT_CONN)

				cancel()

				return nil
			}

			_, ok := mapICN1[node.Addr]
			if !ok {
				mapICN1[node.Addr] = node.Addr

				cp++

				if cp == 4 {
					cancel()

					return nil
				}
			}

			return nil
		})

	node1.RegNodeEventFunc(EventOnNodeConnected,
		func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
			if addr2 != node.Addr {
				errstr = "TestCheckNode node1 EventOnNodeConnected addr fail"
				cancel()

				return nil
			}

			if !node.ConnectMe {
				errstr = "TestCheckNode node1 EventOnNodeConnected node.ConnectMe"
				cancel()

				return nil
			}

			_, ok := mapNC1[node.Addr]
			if !ok {
				mapNC1[node.Addr] = node.Addr

				cp++

				if cp == 4 {
					cancel()
				}
			}

			return nil
		})

	node2.RegNodeEventFunc(EventOnIConnectNode,
		func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
			if addr1 != node.Addr {
				errstr = "TestCheckNode node2 EventOnIConnectNode addr fail"
				cancel()

				return nil
			}

			if node.ConnType != coredbpb.CONNECTTYPE_DIRECT_CONN {
				errstr = fmt.Sprintf("TestCheckNode node2 EventOnIConnectNode node.ConnType %v %v",
					node.ConnType, coredbpb.CONNECTTYPE_DIRECT_CONN)

				cancel()

				return nil
			}

			_, ok := mapICN2[node.Addr]
			if !ok {
				mapICN2[node.Addr] = node.Addr

				cp++

				if cp == 4 {
					cancel()

					return nil
				}
			}

			return nil
		})

	node2.RegNodeEventFunc(EventOnNodeConnected,
		func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
			if addr1 != node.Addr {
				errstr = "TestCheckNode node2 EventOnNodeConnected addr fail"
				cancel()

				return nil
			}

			if !node.ConnectMe {
				errstr = "TestCheckNode node2 EventOnNodeConnected node.ConnectMe"
				cancel()

				return nil
			}

			_, ok := mapNC2[node.Addr]
			if !ok {
				mapNC2[node.Addr] = node.Addr

				cp++

				if cp == 4 {
					cancel()
				}
			}

			return nil
		})

	go node1.Start(ctx)
	go node2.Start(ctx)

	<-ctx.Done()

	node1.GetCoreDB().Close()
	node2.GetCoreDB().Close()

	if errstr != "" {
		t.Fatalf(errstr)
	}

	if cp != 4 {
		t.Fatalf("TestCheckNode need some time %v", cp)
	}

	t.Logf("TestCheckNode OK")
}

func TestConnectNodeFail(t *testing.T) {
	cfg1, err := LoadConfig("./test/node5.yaml")
	if err != nil {
		t.Fatalf("TestConnectNodeFail load config %v", err)
	}

	InitJarvisCore(cfg1)
	defer ReleaseJarvisCore()

	node1, err := NewNode(cfg1)
	if err != nil {
		t.Fatalf("TestConnectNodeFail NewNode node1 %v", err)

		return
	}

	nbi := &jarviscorepb.NodeBaseInfo{
		ServAddr:        "127.0.0.1:7898",
		Addr:            "1JJaKpZGhYPuVHc1EKiiHZEswPAB5SybW5",
		Name:            "test005",
		NodeTypeVersion: "v0.1.0",
		NodeType:        "test",
		CoreVersion:     "v0.1.0",
	}
	node1.AddNodeBaseInfo(nbi)

	isfail := false
	errstr := ""

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	node1.RegNodeEventFunc(EventOnNodeConnected,
		func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
			return nil
		})

	node1.RegNodeEventFunc(EventOnIConnectNodeFail,
		func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
			if !node.Deprecated {
				errstr = "EventOnIConnectNodeFail fail"
				cancel()
			}

			isfail = true

			cancel()

			return nil
		})

	go node1.Start(ctx)

	<-ctx.Done()

	node1.GetCoreDB().Close()

	if errstr != "" {
		t.Fatalf(errstr)
	}

	if !isfail {
		t.Fatalf("TestConnectNodeFail no fail.")
	}

	t.Logf("TestConnectNodeFail OK")
}

// func _requestNode(jarvisnode JarvisNode) {
// 	jarvisnode.RequestNodes(nil)
// }

// func TestRequestNodes(t *testing.T) {
// 	cfg1, err := LoadConfig("./test/node7.yaml")
// 	if err != nil {
// 		t.Fatalf("TestRequestNodes load config %v", err)
// 	}

// 	cfg2, err := LoadConfig("./test/node8.yaml")
// 	if err != nil {
// 		t.Fatalf("TestRequestNodes load config %v", err)
// 	}

// 	InitJarvisCore(cfg1)
// 	defer ReleaseJarvisCore()

// 	node1, err := NewNode(cfg1)
// 	if err != nil {
// 		t.Fatalf("TestRequestNodes NewNode node1 %v", err)

// 		return
// 	}

// 	addr1 := node1.GetCoreDB().GetPrivateKey().ToAddress()

// 	node2, err := NewNode(cfg2)
// 	if err != nil {
// 		t.Fatalf("TestRequestNodes NewNode node2 %v", err)

// 		return
// 	}

// 	addr2 := node2.GetCoreDB().GetPrivateKey().ToAddress()

// 	cp := 0

// 	rn := 0
// 	rne := 0

// 	mapICN1 := make(map[string]string)
// 	mapNC1 := make(map[string]string)
// 	mapICN2 := make(map[string]string)
// 	mapNC2 := make(map[string]string)

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	errstr := ""

// 	funcOnRequestNodes := func(ctx context.Context, jarvisnode JarvisNode, lstResult []*ClientGroupProcMsgResults) error {
// 		if len(lstResult) != 1 {
// 			errstr = "TestRequestNodes node1 funcOnRequestNodes addr fail"
// 			cancel()

// 			return nil
// 		}

// 		rne = rne + 1
// 		if rne == rn {
// 			cancel()

// 			return nil
// 		}

// 		return nil
// 	}

// 	node1.RegNodeEventFunc(EventOnIConnectNode,
// 		func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
// 			if addr2 != node.Addr {
// 				errstr = "TestRequestNodes node1 EventOnIConnectNode addr fail"
// 				cancel()

// 				return nil
// 			}

// 			if node.ConnType != coredbpb.CONNECTTYPE_DIRECT_CONN {

// 				errstr = fmt.Sprintf("TestCheckNode node1 EventOnIConnectNode node.ConnType %v %v",
// 					node.ConnType, coredbpb.CONNECTTYPE_DIRECT_CONN)

// 				cancel()

// 				return nil
// 			}

// 			_, ok := mapICN1[node.Addr]
// 			if !ok {
// 				mapICN1[node.Addr] = node.Addr

// 				cp++

// 				if cp == 4 {
// 					node1.RequestNodes(ctx, funcOnRequestNodes)
// 				}
// 			}

// 			return nil
// 		})

// 	node1.RegNodeEventFunc(EventOnNodeConnected,
// 		func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
// 			if addr2 != node.Addr {
// 				errstr = "TestRequestNodes node1 EventOnNodeConnected addr fail"
// 				cancel()

// 				return nil
// 			}

// 			if !node.ConnectMe {
// 				errstr = "TestRequestNodes node1 EventOnNodeConnected node.ConnectMe"
// 				cancel()

// 				return nil
// 			}

// 			_, ok := mapNC1[node.Addr]
// 			if !ok {
// 				mapNC1[node.Addr] = node.Addr

// 				cp++

// 				if cp == 4 {
// 					node1.RequestNodes(funcOnRequestNodes)
// 				}
// 			}

// 			return nil
// 		})

// 	node1.RegNodeEventFunc(EventOnRequestNode,
// 		func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
// 			if addr2 != node.Addr {
// 				errstr = "TestRequestNodes node1 EventOnRequestNode addr fail"
// 				cancel()

// 				return nil
// 			}

// 			rn = rn + 1

// 			return nil
// 		})

// 	// node1.RegNodeEventFunc(EventOnEndRequestNode,
// 	// 	func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
// 	// 		if addr2 != node.Addr {
// 	// 			errstr = "TestRequestNodes node1 EventOnEndRequestNode addr fail"
// 	// 			cancel()

// 	// 			return nil
// 	// 		}

// 	// 		rne = rne + 1
// 	// 		if rne == rn {
// 	// 			cancel()

// 	// 			return nil
// 	// 		}

// 	// 		return nil
// 	// 	})

// 	node2.RegNodeEventFunc(EventOnIConnectNode,
// 		func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
// 			if addr1 != node.Addr {
// 				errstr = "TestRequestNodes node2 EventOnIConnectNode addr fail"
// 				cancel()

// 				return nil
// 			}

// 			if node.ConnType != coredbpb.CONNECTTYPE_DIRECT_CONN {
// 				errstr = fmt.Sprintf("TestCheckNode node2 EventOnIConnectNode node.ConnType %v %v",
// 					node.ConnType, coredbpb.CONNECTTYPE_DIRECT_CONN)

// 				cancel()

// 				return nil
// 			}

// 			_, ok := mapICN2[node.Addr]
// 			if !ok {
// 				mapICN2[node.Addr] = node.Addr

// 				cp++

// 				if cp == 4 {
// 					node1.RequestNodes(funcOnRequestNodes)
// 				}
// 			}

// 			return nil
// 		})

// 	node2.RegNodeEventFunc(EventOnNodeConnected,
// 		func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
// 			if addr1 != node.Addr {
// 				errstr = "TestRequestNodes node2 EventOnNodeConnected addr fail"
// 				cancel()

// 				return nil
// 			}

// 			if !node.ConnectMe {
// 				errstr = "TestRequestNodes node2 EventOnNodeConnected node.ConnectMe"
// 				cancel()

// 				return nil
// 			}

// 			_, ok := mapNC2[node.Addr]
// 			if !ok {
// 				mapNC2[node.Addr] = node.Addr

// 				cp++

// 				if cp == 4 {
// 					node1.RequestNodes(funcOnRequestNodes)
// 				}
// 			}

// 			return nil
// 		})

// 	go node1.Start(ctx)
// 	go node2.Start(ctx)

// 	<-ctx.Done()

// 	node1.GetCoreDB().Close()
// 	node2.GetCoreDB().Close()

// 	if errstr != "" {
// 		t.Fatalf(errstr)
// 	}

// 	if cp != 4 {
// 		t.Fatalf("TestRequestNodes need some time cp:%v", cp)
// 	}

// 	if rn != rne || rn == 0 {
// 		t.Fatalf("TestRequestNodes need some time %v %v", rn, rne)
// 	}

// 	t.Logf("TestRequestNodes OK")
// }
