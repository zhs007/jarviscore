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

	"go.uber.org/zap"
)

func outputErrRCS(t *testing.T, err error, msg string, info string) {
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

func outputRCS(t *testing.T, msg string) {
	jarvisbase.Info(msg)
	t.Logf(msg)
}

// funconcallRCS
type funconcallRCS func(ctx context.Context, err error, obj *objRCS) error

type mapnodeinfoRCS struct {
	iconn  bool
	connme bool
}

type nodeinfoRCS struct {
	mapAddr    sync.Map
	numsIConn  int
	numsConnMe int
}

func (ni *nodeinfoRCS) onIConnectNode(node *coredbpb.NodeInfo) error {
	d, ok := ni.mapAddr.Load(node.Addr)
	if ok {
		mni, ok := d.(*mapnodeinfoRCS)
		if !ok {
			return fmt.Errorf("nodeinfoRCS.onIConnectNode:mapAddr2mapnodeinfo err")
		}

		if !mni.iconn {
			mni.iconn = true

			ni.numsIConn++
		}

		return nil
	}

	mni := &mapnodeinfoRCS{
		iconn: true,
	}

	ni.mapAddr.Store(node.Addr, mni)
	ni.numsIConn++

	return nil
}

func (ni *nodeinfoRCS) onNodeConnected(node *coredbpb.NodeInfo) error {
	d, ok := ni.mapAddr.Load(node.Addr)
	if ok {
		mni, ok := d.(*mapnodeinfoRCS)
		if !ok {
			return fmt.Errorf("nodeinfoRCS.onNodeConnected:mapAddr2mapnodeinfo err")
		}

		if !mni.connme {
			mni.connme = true

			ni.numsConnMe++
		}

		return nil
	}

	mni := &mapnodeinfoRCS{
		connme: true,
	}

	ni.mapAddr.Store(node.Addr, mni)
	ni.numsConnMe++

	return nil
}

type objRCS struct {
	root               JarvisNode
	node1              JarvisNode
	node2              JarvisNode
	rootni             nodeinfoRCS
	node1ni            nodeinfoRCS
	node2ni            nodeinfoRCS
	requestnodes       bool
	requestctrlnode1   bool
	requestctrlnode1ok bool
	err                error
}

func newObjRCS() *objRCS {
	return &objRCS{
		rootni:       nodeinfoRCS{},
		node1ni:      nodeinfoRCS{},
		node2ni:      nodeinfoRCS{},
		requestnodes: false,
	}
}

func (obj *objRCS) isDone() bool {
	if obj.rootni.numsConnMe != 2 || obj.rootni.numsIConn != 2 {
		return false
	}

	return obj.requestnodes && obj.requestctrlnode1 && obj.requestctrlnode1ok
}

func (obj *objRCS) oncheck(ctx context.Context, funcCancel context.CancelFunc) error {
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
		obj.node1ni.numsIConn == 2 && obj.node2ni.numsIConn == 2 && !obj.requestctrlnode1 {

		dat, err := ioutil.ReadFile("./test/test.sh")
		if err != nil {
			obj.err = err

			funcCancel()
		}

		ci, err := BuildCtrlInfoForScriptFile("test.sh", dat, "", "test.script")
		if err != nil {
			obj.err = err

			funcCancel()
		}

		err = obj.node1.RequestCtrl(ctx, obj.node2.GetMyInfo().Addr, ci,
			func(ctx context.Context, jarvisnode JarvisNode,
				lstResult []*JarvisMsgInfo) error {

				lastjmi := lstResult[len(lstResult)-1]
				if IsClientProcMsgResultEnd(lstResult) {
					// if lastjmi.IsEnd() {
					obj.requestctrlnode1ok = true

					if obj.isDone() {
						funcCancel()
					}
				} else if lastjmi.Err != nil {
					obj.err = lastjmi.Err

					funcCancel()
				} else if lastjmi.Msg != nil {
					jarvisbase.Info("objRCS.oncheck:obj.node1.RequestCtrl",
						JSONMsg2Zap("msg", lastjmi.Msg))
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

func (obj *objRCS) onIConn(ctx context.Context, funcCancel context.CancelFunc) error {
	return obj.oncheck(ctx, funcCancel)
}

func (obj *objRCS) onConnMe(ctx context.Context, funcCancel context.CancelFunc) error {
	return obj.oncheck(ctx, funcCancel)
}

func (obj *objRCS) makeString() string {
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

func startTestNodeRCS(ctx context.Context, cfgfilename string, ni *nodeinfoRCS, obj *objRCS,
	oniconn funconcallRCS, onconnme funconcallRCS) (JarvisNode, error) {

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

func TestRequestCtrlScript(t *testing.T) {
	rootcfg, err := LoadConfig("./test/test5070_requestscriptroot.yaml")
	if err != nil {
		t.Fatalf("TestRequestCtrlScript load config %v err is %v", "./test/test5070_requestscriptroot.yaml", err)

		return
	}

	InitJarvisCore(rootcfg)
	defer ReleaseJarvisCore()

	obj := newObjRCS()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var errobj error

	oniconn := func(ctx context.Context, err error, obj *objRCS) error {
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

	onconnme := func(ctx context.Context, err error, obj *objRCS) error {
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

	obj.root, err = startTestNodeRCS(ctx, "./test/test5070_requestscriptroot.yaml", &obj.rootni, obj,
		oniconn, onconnme)
	if err != nil {
		outputErrRCS(t, err, "TestRequestCtrlScript startTestNodeRCS root", "")

		return
	}

	obj.node1, err = startTestNodeRCS(ctx, "./test/test5071_requestscript1.yaml", &obj.node1ni, obj,
		oniconn, onconnme)
	if err != nil {
		outputErrRCS(t, err, "TestRequestCtrlScript startTestNodeRCS node1", "")

		return
	}

	obj.node2, err = startTestNodeRCS(ctx, "./test/test5072_requestscript2.yaml", &obj.node2ni, obj,
		oniconn, onconnme)
	if err != nil {
		outputErrRCS(t, err, "TestRequestCtrlScript startTestNodeRCS node2", "")

		return
	}

	go obj.root.Start(ctx)
	time.Sleep(time.Second * 1)
	go obj.node1.Start(ctx)
	go obj.node2.Start(ctx)

	<-ctx.Done()

	if errobj != nil {
		outputErrRCS(t, errobj, "TestRequestCtrlScript", "")

		return
	}

	if !obj.isDone() {
		outputErrRCS(t, nil, "TestRequestCtrlScript no done", obj.makeString())

		return
	}

	outputRCS(t, "TestRequestCtrlScript OK")
}
