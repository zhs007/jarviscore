package jarviscore

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"github.com/zhs007/jarviscore/coredb/proto"
	pb "github.com/zhs007/jarviscore/proto"
)

func sendCtrl(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
	dat, err := ioutil.ReadFile("./test/test.sh")
	if err != nil {
		return err
	}

	ci, err := BuildCtrlInfoForScriptFile(1, "test.sh", dat, "")
	if err != nil {
		return err
	}

	err = jarvisnode.RequestCtrl(ctx, node.Addr, ci)
	if err != nil {
		return err
	}

	return nil
}

func TestCheckNodeCtrl(t *testing.T) {
	cfg1, err := LoadConfig("./test/node3.yaml")
	if err != nil {
		t.Fatalf("TestCheckNodeCtrl load config %v", err)
	}

	cfg2, err := LoadConfig("./test/node4.yaml")
	if err != nil {
		t.Fatalf("TestCheckNodeCtrl load config %v", err)
	}

	InitJarvisCore(cfg1)
	defer ReleaseJarvisCore()

	node1 := NewNode(cfg1)
	addr1 := node1.GetCoreDB().GetPrivateKey().ToAddress()

	node2 := NewNode(cfg2)
	addr2 := node2.GetCoreDB().GetPrivateKey().ToAddress()

	cp := 0
	ctrlp := 0

	node1.RegNodeEventFunc(EventOnIConnectNode,
		func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
			if addr2 != node.Addr {
				t.Fatalf("TestCheckNodeCtrl node addr fail")
			}

			if ctrlp == 0 {
				err := sendCtrl(ctx, jarvisnode, node)
				if err != nil {
					t.Fatalf("TestCheckNodeCtrl sendctrl err %v", err)
				}

				ctrlp++
			}

			cp++

			return nil
		})

	node1.RegNodeEventFunc(EventOnNodeConnected,
		func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
			if addr2 != node.Addr {
				t.Fatalf("TestCheckNodeCtrl node addr fail")
			}

			cp++

			return nil
		})

	node1.RegMsgEventFunc(EventOnCtrlResult,
		func(ctx context.Context, jarvisnode JarvisNode, msg *pb.JarvisMsg) error {
			ctrlp++

			return nil
		})

	node2.RegNodeEventFunc(EventOnIConnectNode,
		func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
			if addr1 != node.Addr {
				t.Fatalf("TestCheckNodeCtrl node addr fail")
			}

			cp++

			return nil
		})

	node2.RegNodeEventFunc(EventOnNodeConnected,
		func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
			if addr1 != node.Addr {
				t.Fatalf("TestCheckNodeCtrl node addr fail")
			}

			cp++

			return nil
		})

	node2.RegMsgEventFunc(EventOnCtrl,
		func(ctx context.Context, jarvisnode JarvisNode, msg *pb.JarvisMsg) error {
			ctrlp++

			return nil
		})

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	go node1.Start(ctx)
	go node2.Start(ctx)

	<-ctx.Done()

	node1.GetCoreDB().Close()
	node2.GetCoreDB().Close()

	if cp != 4 {
		t.Fatalf("TestCheckNodeCtrl need some time %v", cp)
	}

	if ctrlp != 3 {
		t.Fatalf("TestCheckNodeCtrl ctrl need some time %v", ctrlp)
	}

	// node.RegMsgEventFunc(EventOnCtrl, onCtrl)
	// node.RegMsgEventFunc(EventOnCtrlResult, onCtrlResult)

	// node.Start(context.Background())

	t.Logf("TestCheckNodeCtrl OK")
}
