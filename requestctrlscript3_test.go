package jarviscore

import (
	"context"
	"fmt"
	"io/ioutil"
	"sync"
	"testing"
	"time"

	"github.com/zhs007/jarviscore/base"
	"github.com/zhs007/jarviscore/coredb/proto"
	pb "github.com/zhs007/jarviscore/proto"

	"go.uber.org/zap"
)

func outputErrRCS3(t *testing.T, err error, msg string, info string) {
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

func outputRCS3(t *testing.T, msg string) {
	jarvisbase.Info(msg)
	t.Logf(msg)
}

// funconcallRCS3
type funconcallRCS3 func(ctx context.Context, err error, obj *objRCS3) error

type mapnodeinfoRCS3 struct {
	iconn  bool
	connme bool
}

type nodeinfoRCS3 struct {
	mapAddr    sync.Map
	numsIConn  int
	numsConnMe int
}

func (ni *nodeinfoRCS3) onIConnectNode(node *coredbpb.NodeInfo) error {
	d, ok := ni.mapAddr.Load(node.Addr)
	if ok {
		mni, ok := d.(*mapnodeinfoRCS3)
		if !ok {
			return fmt.Errorf("nodeinfoRCS3.onIConnectNode:mapAddr2mapnodeinfo err")
		}

		if !mni.iconn {
			mni.iconn = true

			ni.numsIConn++
		}

		return nil
	}

	mni := &mapnodeinfoRCS3{
		iconn: true,
	}

	ni.mapAddr.Store(node.Addr, mni)
	ni.numsIConn++

	return nil
}

func (ni *nodeinfoRCS3) onNodeConnected(node *coredbpb.NodeInfo) error {
	d, ok := ni.mapAddr.Load(node.Addr)
	if ok {
		mni, ok := d.(*mapnodeinfoRCS3)
		if !ok {
			return fmt.Errorf("nodeinfoRCS3.onNodeConnected:mapAddr2mapnodeinfo err")
		}

		if !mni.connme {
			mni.connme = true

			ni.numsConnMe++
		}

		return nil
	}

	mni := &mapnodeinfoRCS3{
		connme: true,
	}

	ni.mapAddr.Store(node.Addr, mni)
	ni.numsConnMe++

	return nil
}

type objRCS3 struct {
	root               JarvisNode
	node1              JarvisNode
	node2              JarvisNode
	rootni             nodeinfoRCS3
	node1ni            nodeinfoRCS3
	node2ni            nodeinfoRCS3
	requestnodes       bool
	requestctrlnode1   bool
	requestctrlnode1ok bool
	err                error
}

func newObjRCS3() *objRCS3 {
	return &objRCS3{
		rootni:       nodeinfoRCS3{},
		node1ni:      nodeinfoRCS3{},
		node2ni:      nodeinfoRCS3{},
		requestnodes: false,
	}
}

func (obj *objRCS3) isDone() bool {
	if obj.rootni.numsConnMe != 2 || obj.rootni.numsIConn != 2 {
		return false
	}

	return obj.requestnodes && obj.requestctrlnode1 && obj.requestctrlnode1ok
}

func (obj *objRCS3) oncheck(ctx context.Context, funcCancel context.CancelFunc) error {
	if obj.rootni.numsConnMe == 2 &&
		obj.node1ni.numsConnMe >= 1 && obj.node1ni.numsIConn >= 1 &&
		obj.node2ni.numsConnMe >= 1 && obj.node2ni.numsIConn >= 1 &&
		!obj.requestnodes {
		err := obj.node1.RequestNodes(ctx, nil)
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
		obj.node1ni.numsIConn == 2 && obj.node2ni.numsIConn == 2 && !obj.requestctrlnode1 {

		dat, err := ioutil.ReadFile("./test/test.sh")
		if err != nil {
			obj.err = err

			funcCancel()
		}

		sf := &pb.FileData{
			Filename: "test.sh",
			File:     dat,
		}

		ci, err := BuildCtrlInfoForScriptFile3(sf, []string{
			"./test/test5090_requestscript3root",
			"./test/test5080_requestscript2root",
		})
		if err != nil {
			obj.err = err

			funcCancel()
		}

		err = obj.node1.RequestCtrl(ctx, obj.node2.GetMyInfo().Addr, ci,
			func(ctx context.Context, jarvisnode JarvisNode,
				lstResult []*JarvisMsgInfo) error {

				lastjmi := lstResult[len(lstResult)-1]
				if lastjmi.Err != nil {
					obj.err = lastjmi.Err

					funcCancel()
				} else if lastjmi.Msg != nil {
					jarvisbase.Info("objRCS3.oncheck:obj.node1.RequestCtrl",
						JSONMsg2Zap("msg", lastjmi.Msg))
				} else {
					obj.requestctrlnode1ok = true

					if obj.isDone() {
						funcCancel()
					}
				}

				return nil
			})
		if err != nil {
			return err
		}

		obj.requestctrlnode1 = true
	}

	return nil
}

func (obj *objRCS3) onIConn(ctx context.Context, funcCancel context.CancelFunc) error {
	return obj.oncheck(ctx, funcCancel)
}

func (obj *objRCS3) onConnMe(ctx context.Context, funcCancel context.CancelFunc) error {
	return obj.oncheck(ctx, funcCancel)
}

func (obj *objRCS3) makeString() string {
	return fmt.Sprintf("root(%v %v) node1(%v %v), node2(%v %v) requestnodes %v requestctrlnode1 %v requestctrlnode1ok %v root %v node1 %v node2 %v",
		obj.rootni.numsIConn, obj.rootni.numsConnMe,
		obj.node1ni.numsIConn, obj.node1ni.numsConnMe,
		obj.node2ni.numsIConn, obj.node2ni.numsConnMe,
		obj.requestnodes,
		obj.requestctrlnode1,
		obj.requestctrlnode1ok,
		obj.root.BuildStatus(),
		obj.node1.BuildStatus(),
		obj.node2.BuildStatus())
}

func startTestNodeRCS3(ctx context.Context, cfgfilename string, ni *nodeinfoRCS3, obj *objRCS3,
	oniconn funconcallRCS3, onconnme funconcallRCS3) (JarvisNode, error) {

	cfg, err := LoadConfig(cfgfilename)
	if err != nil {
		return nil, fmt.Errorf("startTestNode load config %v err is %v", cfgfilename, err)
	}

	curnode, err := NewNode(cfg)
	if err != nil {
		return nil, fmt.Errorf("startTestNode NewNode node %v", err)
	}

	curnode.SetNodeTypeInfo("testupdnode", "0.7.22")

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

func TestRequestCtrlScript3(t *testing.T) {
	rootcfg, err := LoadConfig("./test/test5090_requestscript3root.yaml")
	if err != nil {
		t.Fatalf("TestRequestCtrlScript3 load config %v err is %v", "./test/test5090_requestscript3root.yaml", err)

		return
	}

	InitJarvisCore(rootcfg)
	defer ReleaseJarvisCore()

	obj := newObjRCS3()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var errobj error

	oniconn := func(ctx context.Context, err error, obj *objRCS3) error {
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

	onconnme := func(ctx context.Context, err error, obj *objRCS3) error {
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

	obj.root, err = startTestNodeRCS3(ctx, "./test/test5090_requestscript3root.yaml", &obj.rootni, obj,
		oniconn, onconnme)
	if err != nil {
		outputErrRCS3(t, err, "TestRequestCtrlScript3 startTestNodeRCS3 root", "")

		return
	}

	obj.node1, err = startTestNodeRCS3(ctx, "./test/test5091_requestscript31.yaml", &obj.node1ni, obj,
		oniconn, onconnme)
	if err != nil {
		outputErrRCS3(t, err, "TestRequestCtrlScript3 startTestNodeRCS3 node1", "")

		return
	}

	obj.node2, err = startTestNodeRCS3(ctx, "./test/test5092_requestscript32.yaml", &obj.node2ni, obj,
		oniconn, onconnme)
	if err != nil {
		outputErrRCS3(t, err, "TestRequestCtrlScript3 startTestNodeRCS3 node2", "")

		return
	}

	go obj.root.Start(ctx)
	time.Sleep(time.Second * 1)
	go obj.node1.Start(ctx)
	go obj.node2.Start(ctx)

	<-ctx.Done()

	if errobj != nil {
		outputErrRCS3(t, errobj, "TestRequestCtrlScript3", "")

		return
	}

	if !obj.isDone() {
		outputErrRCS3(t, nil, "TestRequestCtrlScript3 no done", obj.makeString())

		return
	}

	outputRCS3(t, "TestRequestCtrlScript3 OK")
}
