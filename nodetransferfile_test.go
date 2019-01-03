package jarviscore

import (
	"context"
	"io/ioutil"
	"path"
	"testing"
	"time"

	"github.com/zhs007/jarviscore/coredb/proto"
	pb "github.com/zhs007/jarviscore/proto"
)

func transferFile(ctx context.Context, jarvisnode JarvisNode, destaddr string, fn string) (string, error) {
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		return "", err
	}

	md5str := GetMD5String(b)

	jarvisnode.SendFile(ctx, destaddr, &pb.FileData{
		File:     b,
		Filename: path.Join("./test/desttf.dat"),
		// DestPath: path.Join("./test/desttf.dat"),
	}, nil)

	return md5str, nil
}

func TestCheckNodeTransferFile(t *testing.T) {
	cfg1, err := LoadConfig("./test/node9.yaml")
	if err != nil {
		t.Fatalf("TestCheckNodeTransferFile load config %v", err)
	}

	cfg2, err := LoadConfig("./test/node10.yaml")
	if err != nil {
		t.Fatalf("TestCheckNodeTransferFile load config %v", err)
	}

	InitJarvisCore(cfg1)
	defer ReleaseJarvisCore()

	node1 := NewNode(cfg1)
	addr1 := node1.GetCoreDB().GetPrivateKey().ToAddress()

	node2 := NewNode(cfg2)
	addr2 := node2.GetCoreDB().GetPrivateKey().ToAddress()

	cp := 0
	tfp := 0

	mapICN1 := make(map[string]string)
	mapNC1 := make(map[string]string)
	mapICN2 := make(map[string]string)
	mapNC2 := make(map[string]string)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	errstr := ""
	md5str := ""

	funcStartTransferFile := func(ctx context.Context, cancel context.CancelFunc) {
		if cp == 4 && tfp == 0 {
			md5str, err = transferFile(ctx, node1, node2.GetMyInfo().Addr, "./test/test.sh")
			tfp++
		} else if cp == 4 && tfp == 3 {
			cancel()
		}
	}

	node1.RegNodeEventFunc(EventOnIConnectNode,
		func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
			if addr2 != node.Addr {
				errstr = "TestCheckNodeTransferFile node addr fail"

				cancel()
			}

			_, ok := mapICN1[node.Addr]
			if !ok {
				mapICN1[node.Addr] = node.Addr

				cp++

				funcStartTransferFile(ctx, cancel)
			}

			return nil
		})

	node1.RegNodeEventFunc(EventOnNodeConnected,
		func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
			if addr2 != node.Addr {
				errstr = "TestCheckNodeTransferFile node addr fail"

				cancel()
			}

			_, ok := mapNC1[node.Addr]
			if !ok {
				mapNC1[node.Addr] = node.Addr

				cp++

				funcStartTransferFile(ctx, cancel)
			}

			return nil
		})

	node1.RegMsgEventFunc(EventOnReplyTransferFile,
		func(ctx context.Context, jarvisnode JarvisNode, msg *pb.JarvisMsg) error {
			rtf := msg.GetReplyTransferFile()
			if rtf == nil {
				errstr = "no rtf"

				cancel()
			}

			if rtf.Md5String != md5str {
				errstr = "md5 fail"

				cancel()
			}

			tfp++
			if cp == 4 && tfp == 3 {
				cancel()
			}

			return nil
		})

	node2.RegNodeEventFunc(EventOnIConnectNode,
		func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
			if addr1 != node.Addr {
				errstr = "TestCheckNodeTransferFile node addr fail"

				cancel()
			}

			_, ok := mapICN2[node.Addr]
			if !ok {
				mapICN2[node.Addr] = node.Addr

				cp++

				funcStartTransferFile(ctx, cancel)
			}

			return nil
		})

	node2.RegNodeEventFunc(EventOnNodeConnected,
		func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
			if addr1 != node.Addr {
				errstr = "TestCheckNodeTransferFile node addr fail"

				cancel()
			}

			_, ok := mapNC2[node.Addr]
			if !ok {
				mapNC2[node.Addr] = node.Addr

				cp++

				funcStartTransferFile(ctx, cancel)
			}

			return nil
		})

	node2.RegMsgEventFunc(EventOnTransferFile,
		func(ctx context.Context, jarvisnode JarvisNode, msg *pb.JarvisMsg) error {
			tfp++

			return nil
		})

	go node1.Start(ctx)
	go node2.Start(ctx)

	<-ctx.Done()

	node1.GetCoreDB().Close()
	node2.GetCoreDB().Close()

	if errstr != "" {
		t.Fatalf("TestCheckNodeTransferFile err %v", errstr)
	}

	if cp != 4 {
		t.Fatalf("TestCheckNodeTransferFile need some time %v", cp)
	}

	if tfp != 3 {
		t.Fatalf("TestCheckNodeTransferFile ctrl need some time %v", tfp)
	}

	t.Logf("TestCheckNodeTransferFile OK")
}
