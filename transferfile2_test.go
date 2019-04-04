package jarviscore

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/zhs007/jarviscore/proto"

	"github.com/zhs007/jarviscore/base"
	"github.com/zhs007/jarviscore/coredb/proto"

	"go.uber.org/zap"
)

func sendfile22node(ctx context.Context, fn string, destfn string, srcnode JarvisNode, destaddr string, funcOnResult FuncOnProcMsgResult) error {
	dat, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}

	fd := &jarviscorepb.FileData{
		File:     dat,
		Filename: destfn,
	}

	return srcnode.SendFile2(ctx, destaddr, fd, funcOnResult)
}

func randfillFileTF2(fn string, len int) error {
	f, err := os.Create(fn)
	if err != nil {
		return err
	}

	defer f.Close()

	l4 := len / 4
	l1 := len % 4

	for i := 0; i < l4; i++ {
		d := []byte{
			byte(rand.Intn(256)),
			byte(rand.Intn(256)),
			byte(rand.Intn(256)),
			byte(rand.Intn(256)),
		}

		n, err := f.Write(d)
		if err != nil {
			return nil
		}

		if n != 4 {
			return fmt.Errorf("randfillFile len err")
		}
	}

	for i := 0; i < l1; i++ {
		d := []byte{
			byte(rand.Intn(256)),
		}

		n, err := f.Write(d)
		if err != nil {
			return nil
		}

		if n != 1 {
			return fmt.Errorf("randfillFile len err")
		}
	}

	return nil
}

func outputErrTF2(t *testing.T, err error, msg string, info string) {
	if err == nil && info == "" {
		jarvisbase.Error(msg)
		t.Fatalf(msg)

		return
	} else if err == nil {
		jarvisbase.Error(msg, zap.String("info", info))
		t.Fatalf(msg+" info %v", info)

		return
	}

	jarvisbase.Error(msg, zap.Error(err))
	t.Fatalf(msg+" err %v", err)
}

func outputTF2(t *testing.T, msg string) {
	jarvisbase.Info(msg)
	t.Logf(msg)
}

// funconcallTF2
type funconcallTF2 func(ctx context.Context, err error, obj *objTF2) error

type mapnodeinfoTF2 struct {
	iconn  bool
	connme bool
}

type nodeinfoTF2 struct {
	mapAddr    sync.Map
	numsIConn  int
	numsConnMe int
}

func (ni *nodeinfoTF2) onIConnectNode(node *coredbpb.NodeInfo) error {
	d, ok := ni.mapAddr.Load(node.Addr)
	if ok {
		mni, ok := d.(*mapnodeinfoTF2)
		if !ok {
			return fmt.Errorf("nodeinfoRF2.onIConnectNode:mapAddr2mapnodeinfo err")
		}

		if !mni.iconn {
			mni.iconn = true

			ni.numsIConn++
		}

		return nil
	}

	mni := &mapnodeinfoTF2{
		iconn: true,
	}

	ni.mapAddr.Store(node.Addr, mni)
	ni.numsIConn++

	return nil
}

func (ni *nodeinfoTF2) onNodeConnected(node *coredbpb.NodeInfo) error {
	d, ok := ni.mapAddr.Load(node.Addr)
	if ok {
		mni, ok := d.(*mapnodeinfoTF2)
		if !ok {
			return fmt.Errorf("nodeinfoRF2.onNodeConnected:mapAddr2mapnodeinfo err")
		}

		if !mni.connme {
			mni.connme = true

			ni.numsConnMe++
		}

		return nil
	}

	mni := &mapnodeinfoTF2{
		connme: true,
	}

	ni.mapAddr.Store(node.Addr, mni)
	ni.numsConnMe++

	return nil
}

type objTF2 struct {
	root            JarvisNode
	node1           JarvisNode
	node2           JarvisNode
	rootni          nodeinfoTF2
	node1ni         nodeinfoTF2
	node2ni         nodeinfoTF2
	requestnodes    bool
	transferfile1   bool
	transferfile1ok bool
	transferfile2   bool
	transferfile2ok bool
	err             error
}

func newObjTF2() *objTF2 {
	return &objTF2{
		rootni:       nodeinfoTF2{},
		node1ni:      nodeinfoTF2{},
		node2ni:      nodeinfoTF2{},
		requestnodes: false,
	}
}

func (obj *objTF2) isDone() bool {
	if obj.rootni.numsConnMe != 2 || obj.rootni.numsIConn != 2 {
		return false
	}

	return obj.transferfile1 && obj.transferfile1ok && obj.transferfile2 && obj.transferfile2ok
}

func (obj *objTF2) oncheck(ctx context.Context, funcCancel context.CancelFunc) error {
	if obj.rootni.numsConnMe == 2 &&
		obj.node1ni.numsConnMe >= 1 && obj.node1ni.numsIConn >= 1 &&
		obj.node2ni.numsConnMe >= 1 && obj.node2ni.numsIConn >= 1 &&
		!obj.requestnodes {

		err := obj.node1.GetCoreDB().TrustNode(obj.node2.GetMyInfo().Addr)
		if err != nil {
			jarvisbase.Warn("objUN.oncheck:node1.TrustNode",
				zap.Error(err))

			return err
		}

		err = obj.node2.GetCoreDB().TrustNode(obj.node1.GetMyInfo().Addr)
		if err != nil {
			jarvisbase.Warn("objUN.oncheck:node2.TrustNode",
				zap.Error(err))

			return err
		}

		err = obj.node1.RequestNodes(ctx, nil)
		if err != nil {
			return err
		}

		err = obj.node2.RequestNodes(ctx, nil)
		if err != nil {
			return err
		}

		obj.requestnodes = true
	}

	if obj.node1ni.numsConnMe == 2 && obj.node2ni.numsConnMe == 2 &&
		obj.node1ni.numsIConn == 2 && obj.node2ni.numsIConn == 2 && !obj.transferfile1 {

		curresultnums := 0

		err := sendfile22node(ctx, "./test/tf2001.dat", "./test/node1_tf2001.dat", obj.node1, obj.node2.GetMyInfo().Addr,
			func(ctx context.Context, jarvisnode JarvisNode,
				lstResult []*JarvisMsgInfo) error {

				for i := len(lstResult) - 1; i < len(lstResult); i++ {
					if lstResult[i].Msg != nil {
						jarvisbase.Info("sendfile22node obj.node1", JSONMsg2Zap("result", lstResult[i].Msg))
					}
				}

				// jarvisbase.Info("obj.node1.RequestFile", jarvisbase.JSON("result", lstResult))

				if len(lstResult) > curresultnums {
					for ; curresultnums < len(lstResult); curresultnums++ {
						if lstResult[curresultnums].Err != nil {
							obj.err = fmt.Errorf("lstResult[%v].Err %v", curresultnums, lstResult[curresultnums].Err)

							curresultnums++

							funcCancel()

							return nil
						}

						if lstResult[curresultnums].Msg != nil &&
							lstResult[curresultnums].Msg.MsgType == jarviscorepb.MSGTYPE_REPLY2 {

							if lstResult[curresultnums].Msg.ReplyType == jarviscorepb.REPLYTYPE_END {
								obj.transferfile1ok = true
							}
						}

						if lstResult[curresultnums].JarvisResultType == JarvisResultTypeReplyStreamEnd {
							// if lstResult[curresultnums].Err == nil && lstResult[curresultnums].Msg == nil {
							if obj.isDone() {
								funcCancel()
							}

							return nil
						}
					}
				}

				return nil
			})
		if err != nil {
			obj.err = err

			return err
		}

		obj.transferfile1 = true
	}

	if obj.node1ni.numsConnMe == 2 && obj.node2ni.numsConnMe == 2 &&
		obj.node1ni.numsIConn == 2 && obj.node2ni.numsIConn == 2 && !obj.transferfile2 {

		curresultnums := 0

		err := sendfile22node(ctx, "./test/tf2002.dat", "./test/node2_tf2002.dat", obj.node2, obj.node1.GetMyInfo().Addr,
			func(ctx context.Context, jarvisnode JarvisNode,
				lstResult []*JarvisMsgInfo) error {

				for i := len(lstResult) - 1; i < len(lstResult); i++ {
					if lstResult[i].Msg != nil {
						jarvisbase.Info("sendfile22node obj.node2", JSONMsg2Zap("result", lstResult[i].Msg))
					}
				}

				// jarvisbase.Info("obj.node1.RequestFile", jarvisbase.JSON("result", lstResult))

				if len(lstResult) > curresultnums {
					for ; curresultnums < len(lstResult); curresultnums++ {
						if lstResult[curresultnums].Err != nil {
							obj.err = fmt.Errorf("lstResult[%v].Err %v", curresultnums, lstResult[curresultnums].Err)

							curresultnums++

							funcCancel()

							return nil
						}

						if lstResult[curresultnums].Msg != nil &&
							lstResult[curresultnums].Msg.MsgType == jarviscorepb.MSGTYPE_REPLY2 {

							if lstResult[curresultnums].Msg.ReplyType == jarviscorepb.REPLYTYPE_END {
								obj.transferfile2ok = true
							}
						}

						if lstResult[curresultnums].JarvisResultType == JarvisResultTypeReplyStreamEnd {
							// if lstResult[curresultnums].Err == nil && lstResult[curresultnums].Msg == nil {
							if obj.isDone() {
								funcCancel()
							}

							return nil
						}
					}
				}

				return nil
			})
		if err != nil {
			obj.err = err

			return err
		}

		obj.transferfile2 = true
	}

	return nil
}

func (obj *objTF2) onIConn(ctx context.Context, funcCancel context.CancelFunc) error {
	return obj.oncheck(ctx, funcCancel)
}

func (obj *objTF2) onConnMe(ctx context.Context, funcCancel context.CancelFunc) error {
	return obj.oncheck(ctx, funcCancel)
}

func (obj *objTF2) makeString() string {
	return fmt.Sprintf("root(%v %v) node1(%v %v), node2(%v %v) requestnodes %v transferfile1 %v transferfile1ok %v transferfile2 %v transferfile2ok %v root %v node1 %v node2 %v",
		obj.rootni.numsIConn, obj.rootni.numsConnMe,
		obj.node1ni.numsIConn, obj.node1ni.numsConnMe,
		obj.node2ni.numsIConn, obj.node2ni.numsConnMe,
		obj.requestnodes,
		obj.transferfile1,
		obj.transferfile1ok,
		obj.transferfile2,
		obj.transferfile2ok,
		obj.root.BuildStatus(),
		obj.node1.BuildStatus(),
		obj.node2.BuildStatus())
}

func startTestNodeTF2(ctx context.Context, cfgfilename string, ni *nodeinfoTF2, obj *objTF2,
	oniconn funconcallTF2, onconnme funconcallTF2) (JarvisNode, error) {

	cfg, err := LoadConfig(cfgfilename)
	if err != nil {
		return nil, fmt.Errorf("startTestNode load config %v err is %v", cfgfilename, err)
	}

	curnode, err := NewNode(cfg)
	if err != nil {
		return nil, fmt.Errorf("startTestNode NewNode node %v", err)
	}

	curnode.SetNodeTypeInfo("testreqfile", "0.7.22")

	curnode.RegNodeEventFunc(EventOnIConnectNode,
		func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
			err := ni.onIConnectNode(node)

			oniconn(ctx, err, obj)

			return nil
		})

	curnode.RegNodeEventFunc(EventOnNodeConnected,
		func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error {
			err := ni.onNodeConnected(node)

			onconnme(ctx, err, obj)

			return nil
		})

	return curnode, nil
}

func TestTransferFile2(t *testing.T) {
	randfillFileTF2("./test/tf2001.dat", 2*1024*1024)
	randfillFileTF2("./test/tf2002.dat", 10*1024*1024)

	rootcfg, err := LoadConfig("./test/test5060_transferfile2root.yaml")
	if err != nil {
		t.Fatalf("TestTransferFile2 load config %v err is %v", "./test/test5060_transferfile2root.yaml", err)

		return
	}

	InitJarvisCore(rootcfg)
	defer ReleaseJarvisCore()

	obj := newObjTF2()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var errobj error

	oniconn := func(ctx context.Context, err error, obj *objTF2) error {
		if err != nil {
			errobj = err

			cancel()

			return nil
		}

		err1 := obj.onIConn(ctx, cancel)
		if err1 != nil {
			errobj = err1

			cancel()

			return nil
		}

		if obj.isDone() {
			cancel()

			return nil
		}

		return nil
	}

	onconnme := func(ctx context.Context, err error, obj *objTF2) error {
		if err != nil {
			errobj = err

			cancel()

			return nil
		}

		err1 := obj.onConnMe(ctx, cancel)
		if err1 != nil {
			errobj = err1

			cancel()

			return nil
		}

		if obj.isDone() {
			cancel()

			return nil
		}

		return nil
	}

	obj.root, err = startTestNodeTF2(ctx, "./test/test5060_transferfile2root.yaml", &obj.rootni, obj,
		oniconn, onconnme)
	if err != nil {
		outputErrTF2(t, err, "TestTransferFile2 startTestNodeTF2 root", "")

		return
	}

	obj.node1, err = startTestNodeTF2(ctx, "./test/test5061_transferfile21.yaml", &obj.node1ni, obj,
		oniconn, onconnme)
	if err != nil {
		outputErrTF2(t, err, "TestTransferFile2 startTestNodeTF2 node1", "")

		return
	}

	obj.node2, err = startTestNodeTF2(ctx, "./test/test5062_transferfile22.yaml", &obj.node2ni, obj,
		oniconn, onconnme)
	if err != nil {
		outputErrTF2(t, err, "TestTransferFile2 startTestNodeTF2 node2", "")

		return
	}

	go obj.root.Start(ctx)
	time.Sleep(time.Second * 1)
	go obj.node1.Start(ctx)
	go obj.node2.Start(ctx)

	<-ctx.Done()

	if errobj != nil {
		outputErrTF2(t, errobj, "TestTransferFile2", "")

		return
	}

	if obj.err != nil {
		outputErrTF2(t, obj.err, "TestTransferFile2", "")

		return
	}

	if !obj.isDone() {
		outputErrTF2(t, nil, "TestTransferFile2 no done", obj.makeString())

		return
	}

	outputTF2(t, "TestTransferFile2 OK")
}
