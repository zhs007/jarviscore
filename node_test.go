package jarviscore

import (
	"context"
	"testing"
	"time"

	"github.com/zhs007/jarviscore/coredb/proto"
)

func TestCheckPrivateKey(t *testing.T) {
	cfg, err := LoadConfig("./test/node1.yaml")
	if err != nil {
		t.Fatalf("TestCheckPrivateKey load config %v", err)
	}

	InitJarvisCore(cfg)
	defer ReleaseJarvisCore()

	node1 := NewNode(cfg)
	addr1 := node1.GetCoreDB().GetPrivateKey().ToAddress()
	node1.GetCoreDB().Close()

	node2 := NewNode(cfg)
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

	node1 := NewNode(cfg1)
	addr1 := node1.GetCoreDB().GetPrivateKey().ToAddress()

	node2 := NewNode(cfg2)
	addr2 := node2.GetCoreDB().GetPrivateKey().ToAddress()

	cp := 0

	node1.RegNodeEventFunc(EventOnIConnectNode,
		func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
			if addr2 != node.Addr {
				t.Fatalf("TestCheckNode node addr fail")
			}

			cp++

			return nil
		})

	node1.RegNodeEventFunc(EventOnNodeConnected,
		func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
			if addr2 != node.Addr {
				t.Fatalf("TestCheckNode node addr fail")
			}

			cp++

			return nil
		})

	node2.RegNodeEventFunc(EventOnIConnectNode,
		func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
			if addr1 != node.Addr {
				t.Fatalf("TestCheckNode node addr fail")
			}

			cp++

			return nil
		})

	node2.RegNodeEventFunc(EventOnNodeConnected,
		func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
			if addr1 != node.Addr {
				t.Fatalf("TestCheckNode node addr fail")
			}

			cp++

			return nil
		})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go node1.Start(ctx)
	go node2.Start(ctx)

	<-ctx.Done()

	node1.GetCoreDB().Close()
	node2.GetCoreDB().Close()

	if cp != 4 {
		t.Fatalf("TestCheckNode need some time")
	}

	t.Logf("TestCheckPrivateKey OK")
}
