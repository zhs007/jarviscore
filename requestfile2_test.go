package jarviscore

import (
	"context"
	"fmt"
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

func randfillFile2(fn string, len int) error {
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
			return fmt.Errorf("randfillFile2 len err")
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
			return fmt.Errorf("randfillFile2 len err")
		}
	}

	return nil
}

func outputErrRF2(t *testing.T, err error, msg string, info string) {
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

func outputRF2(t *testing.T, msg string) {
	jarvisbase.Info(msg)
	t.Logf(msg)
}

// funconcallRF2
type funconcallRF2 func(ctx context.Context, err error, obj *objRF2) error

type mapnodeinfoRF2 struct {
	iconn  bool
	connme bool
}

type nodeinfoRF2 struct {
	mapAddr    sync.Map
	numsIConn  int
	numsConnMe int
}

func (ni *nodeinfoRF2) onIConnectNode(node *coredbpb.NodeInfo) error {
	d, ok := ni.mapAddr.Load(node.Addr)
	if ok {
		mni, ok := d.(*mapnodeinfoRF2)
		if !ok {
			return fmt.Errorf("nodeinfoRF2.onIConnectNode:mapAddr2mapnodeinfo err")
		}

		if !mni.iconn {
			mni.iconn = true

			ni.numsIConn++
		}

		return nil
	}

	mni := &mapnodeinfoRF2{
		iconn: true,
	}

	ni.mapAddr.Store(node.Addr, mni)
	ni.numsIConn++

	return nil
}

func (ni *nodeinfoRF2) onNodeConnected(node *coredbpb.NodeInfo) error {
	d, ok := ni.mapAddr.Load(node.Addr)
	if ok {
		mni, ok := d.(*mapnodeinfoRF2)
		if !ok {
			return fmt.Errorf("nodeinfoRF2.onNodeConnected:mapAddr2mapnodeinfo err")
		}

		if !mni.connme {
			mni.connme = true

			ni.numsConnMe++
		}

		return nil
	}

	mni := &mapnodeinfoRF2{
		connme: true,
	}

	ni.mapAddr.Store(node.Addr, mni)
	ni.numsConnMe++

	return nil
}

type objRF2 struct {
	root           JarvisNode
	node1          JarvisNode
	node2          JarvisNode
	rootni         nodeinfoRF2
	node1ni        nodeinfoRF2
	node2ni        nodeinfoRF2
	requestnodes   bool
	requestfile1   bool
	requestfile1ok bool
	requestfile2   bool
	requestfile2ok bool
	err            error
}

func newObjRF2() *objRF2 {
	return &objRF2{
		rootni:       nodeinfoRF2{},
		node1ni:      nodeinfoRF2{},
		node2ni:      nodeinfoRF2{},
		requestnodes: false,
	}
}

func (obj *objRF2) isDone() bool {
	if obj.rootni.numsConnMe != 2 || obj.rootni.numsIConn != 2 {
		return false
	}

	return obj.requestfile1 && obj.requestfile1ok && obj.requestfile2 && obj.requestfile2ok
}

func (obj *objRF2) oncheck(ctx context.Context, funcCancel context.CancelFunc) error {
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
		obj.node1ni.numsIConn == 2 && obj.node2ni.numsIConn == 2 && !obj.requestfile1 {

		curresultnums := 0

		rf := &jarviscorepb.RequestFile{
			Filename: "./test/rf2001.dat",
		}

		err := obj.node1.RequestFile(ctx, obj.node2.GetMyInfo().Addr, rf,
			func(ctx context.Context, jarvisnode JarvisNode,
				lstResult []*JarvisMsgInfo) error {

				if len(lstResult) > 0 {
					if lstResult[len(lstResult)-1].Msg != nil {
						jarvisbase.Info("obj.node1.RequestFile", JSONMsg2Zap("result", lstResult[len(lstResult)-1].Msg))
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
							lstResult[curresultnums].Msg.MsgType == jarviscorepb.MSGTYPE_REPLY_REQUEST_FILE {

							fd := lstResult[curresultnums].Msg.GetFile()
							if fd == nil {
								obj.err = ErrNoFileData

								funcCancel()

								return nil
							}

							if fd.Md5String == "" {
								obj.err = ErrFileDataNoMD5String

								funcCancel()

								return nil
							}

							if fd.Md5String != GetMD5String(fd.File) {
								obj.err = ErrInvalidFileDataMD5String

								funcCancel()

								return nil
							}
						}

						if IsClientProcMsgResultEnd(lstResult) {
							// if lstResult[curresultnums].IsEnd() {
							// if lstResult[curresultnums].JarvisResultType == JarvisResultTypeReplyStreamEnd {
							// if lstResult[curresultnums].Err == nil && lstResult[curresultnums].Msg == nil {
							obj.requestfile1ok = true

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

		obj.requestfile1 = true
	}

	if obj.node1ni.numsConnMe == 2 && obj.node2ni.numsConnMe == 2 &&
		obj.node1ni.numsIConn == 2 && obj.node2ni.numsIConn == 2 && !obj.requestfile2 {

		curresultnums := 0

		rf := &jarviscorepb.RequestFile{
			Filename: "./test/rf2002.dat",
		}

		err := obj.node2.RequestFile(ctx, obj.node1.GetMyInfo().Addr, rf,
			func(ctx context.Context, jarvisnode JarvisNode,
				lstResult []*JarvisMsgInfo) error {

				if len(lstResult) > 0 {
					if lstResult[len(lstResult)-1].Msg != nil {
						jarvisbase.Info("obj.node2.RequestFile", JSONMsg2Zap("result", lstResult[len(lstResult)-1].Msg))
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
							lstResult[curresultnums].Msg.MsgType == jarviscorepb.MSGTYPE_REPLY_REQUEST_FILE {

							fd := lstResult[curresultnums].Msg.GetFile()
							if fd == nil {
								obj.err = ErrNoFileData

								funcCancel()

								return nil
							}

							if fd.Md5String == "" {
								obj.err = ErrFileDataNoMD5String

								funcCancel()

								return nil
							}

							if fd.Md5String != GetMD5String(fd.File) {
								obj.err = ErrInvalidFileDataMD5String

								funcCancel()

								return nil
							}
						}

						if IsClientProcMsgResultEnd(lstResult) {
							// if lstResult[curresultnums].IsEnd() {
							// if lstResult[curresultnums].JarvisResultType == JarvisResultTypeReplyStreamEnd {
							// if lstResult[curresultnums].Err == nil && lstResult[curresultnums].Msg == nil {
							obj.requestfile2ok = true

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

		obj.requestfile2 = true
	}

	return nil
}

func (obj *objRF2) onIConn(ctx context.Context, funcCancel context.CancelFunc) error {
	return obj.oncheck(ctx, funcCancel)
}

func (obj *objRF2) onConnMe(ctx context.Context, funcCancel context.CancelFunc) error {
	return obj.oncheck(ctx, funcCancel)
}

func (obj *objRF2) makeString() string {
	return fmt.Sprintf("root(%v %v) node1(%v %v), node2(%v %v) requestnodes %v requestfile1 %v requestfile1ok %v requestfile2 %v requestfile2ok %v root %v node1 %v node2 %v",
		obj.rootni.numsIConn, obj.rootni.numsConnMe,
		obj.node1ni.numsIConn, obj.node1ni.numsConnMe,
		obj.node2ni.numsIConn, obj.node2ni.numsConnMe,
		obj.requestnodes,
		obj.requestfile1,
		obj.requestfile1ok,
		obj.requestfile2,
		obj.requestfile2ok,
		obj.root.BuildStatus(),
		obj.node1.BuildStatus(),
		obj.node2.BuildStatus())
}

func startTestNodeRF2(ctx context.Context, cfgfilename string, ni *nodeinfoRF2, obj *objRF2,
	oniconn funconcallRF2, onconnme funconcallRF2) (JarvisNode, error) {

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

func TestRequestFile2(t *testing.T) {
	randfillFile2("./test/rf2001.dat", 2*1024*1024)
	randfillFile2("./test/rf2002.dat", 10*1024*1024)

	rootcfg, err := LoadConfig("./test/test5030_reqfile2root.yaml")
	if err != nil {
		t.Fatalf("TestRequestFile2 load config %v err is %v", "./test/test5030_reqfile2root.yaml", err)

		return
	}

	InitJarvisCore(rootcfg)
	defer ReleaseJarvisCore()

	obj := newObjRF2()

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	var errobj error

	oniconn := func(ctx context.Context, err error, obj *objRF2) error {
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

	onconnme := func(ctx context.Context, err error, obj *objRF2) error {
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

	obj.root, err = startTestNodeRF2(ctx, "./test/test5030_reqfile2root.yaml", &obj.rootni, obj,
		oniconn, onconnme)
	if err != nil {
		outputErrRF2(t, err, "TestRequestFile2 startTestNodeRF2 root", "")

		return
	}

	obj.node1, err = startTestNodeRF2(ctx, "./test/test5031_reqfile21.yaml", &obj.node1ni, obj,
		oniconn, onconnme)
	if err != nil {
		outputErrRF2(t, err, "TestRequestFile2 startTestNodeRF2 node1", "")

		return
	}

	obj.node2, err = startTestNodeRF2(ctx, "./test/test5032_reqfile22.yaml", &obj.node2ni, obj,
		oniconn, onconnme)
	if err != nil {
		outputErrRF2(t, err, "TestRequestFile2 startTestNodeRF2 node2", "")

		return
	}

	go obj.root.Start(ctx)
	time.Sleep(time.Second * 1)
	go obj.node1.Start(ctx)
	go obj.node2.Start(ctx)

	<-ctx.Done()

	if errobj != nil {
		outputErrRF2(t, errobj, "TestRequestFile2", "")

		return
	}

	if obj.err != nil {
		outputErrRF2(t, obj.err, "TestRequestFile2", "")

		return
	}

	if !obj.isDone() {
		outputErrRF2(t, nil, "TestRequestFile2 no done", obj.makeString())

		return
	}

	outputRF2(t, "TestRequestFile2 OK")
}
