package jarviscore

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	jarvisbase "github.com/zhs007/jarviscore/base"
	coredbpb "github.com/zhs007/jarviscore/coredb/proto"
	"go.uber.org/zap"
)

// funconcallRN
type funconcallRN func(ctx context.Context, err error, obj *objRN) error

type mapnodeinfoRN struct {
	iconn  bool
	connme bool
}

type nodeinfoRN struct {
	mapAddr    sync.Map
	numsIConn  int
	numsConnMe int
}

func (ni *nodeinfoRN) onIConnectNode(node *coredbpb.NodeInfo) error {
	d, ok := ni.mapAddr.Load(node.Addr)
	if ok {
		mni, ok := d.(*mapnodeinfoRN)
		if !ok {
			return fmt.Errorf("nodeinfoRN.onIConnectNode:mapAddr2mapnodeinfo err")
		}

		if !mni.iconn {
			mni.iconn = true

			ni.numsIConn++
		}

		return nil
	}

	mni := &mapnodeinfoRN{
		iconn: true,
	}

	ni.mapAddr.Store(node.Addr, mni)
	ni.numsIConn++

	return nil
}

func (ni *nodeinfoRN) onNodeConnected(node *coredbpb.NodeInfo) error {
	d, ok := ni.mapAddr.Load(node.Addr)
	if ok {
		mni, ok := d.(*mapnodeinfoRN)
		if !ok {
			return fmt.Errorf("nodeinfoRN.onNodeConnected:mapAddr2mapnodeinfo err")
		}

		if !mni.connme {
			mni.connme = true

			ni.numsConnMe++
		}

		return nil
	}

	mni := &mapnodeinfoRN{
		connme: true,
	}

	ni.mapAddr.Store(node.Addr, mni)
	ni.numsConnMe++

	return nil
}

type objRN struct {
	sync.RWMutex

	root         JarvisNode
	node1        JarvisNode
	node2        JarvisNode
	rootni       nodeinfoRN
	node1ni      nodeinfoRN
	node2ni      nodeinfoRN
	requestnodes bool
	rqnodes1     bool
	rqnodes2     bool
}

func newObjRN() *objRN {
	return &objRN{
		rootni:       nodeinfoRN{},
		node1ni:      nodeinfoRN{},
		node2ni:      nodeinfoRN{},
		requestnodes: false,
	}
}

func (obj *objRN) isDone() bool {
	if obj.rootni.numsConnMe != 2 || obj.rootni.numsIConn != 2 {
		return false
	}

	return obj.requestnodes && obj.rqnodes1 && obj.rqnodes2
}

// cancelFunc - cancel function
type cancelFunc func(err error)

func (obj *objRN) oncheck(ctx context.Context, funcCancel cancelFunc) error {
	obj.Lock()
	defer obj.Unlock()

	if obj.rootni.numsConnMe == 2 &&
		obj.node1ni.numsConnMe >= 1 && obj.node1ni.numsIConn >= 1 &&
		obj.node2ni.numsConnMe >= 1 && obj.node2ni.numsIConn >= 1 &&
		!obj.requestnodes {

		err := obj.node1.RequestNodes(ctx, true, func(ctx context.Context, jarvisnode JarvisNode,
			numsNode int, lstResult []*ClientGroupProcMsgResults) error {

			jarvisbase.Info("objRN.onIConn:node1.RequestNodes",
				zap.Int("numsNode", numsNode),
				jarvisbase.JSON("lstResult", lstResult))

			if numsNode == 1 && len(lstResult) == 1 && IsClientProcMsgResultEnd(lstResult[0].Results) {
				// if numsNode == 1 && len(lstResult) == 1 && lstResult[0].Results[len(lstResult[0].Results)-1].IsEnd() {
				obj.rqnodes1 = true

				if obj.isDone() {
					funcCancel(nil)
				}
			}

			return nil
		})
		if err != nil {
			jarvisbase.Info("objRN.onIConn:node1.RequestNodes",
				zap.Error(err))

			funcCancel(err)

			return err
		}

		err = obj.node2.RequestNodes(ctx, true, func(ctx context.Context, jarvisnode JarvisNode,
			numsNode int, lstResult []*ClientGroupProcMsgResults) error {

			jarvisbase.Info("objRN.onIConn:node2.RequestNodes",
				zap.Int("numsNode", numsNode),
				jarvisbase.JSON("lstResult", lstResult))

			if numsNode == 1 && len(lstResult) == 1 && IsClientProcMsgResultEnd(lstResult[0].Results) {
				// if numsNode == 1 && len(lstResult) == 1 && lstResult[0].Results[len(lstResult[0].Results)-1].IsEnd() {
				obj.rqnodes2 = true

				if obj.isDone() {
					funcCancel(nil)
				}
			}

			return nil
		})
		if err != nil {
			jarvisbase.Info("objRN.onIConn:node1.RequestNodes",
				zap.Error(err))

			funcCancel(err)

			return err
		}

		obj.requestnodes = true
	}

	if obj.isDone() {
		funcCancel(nil)
	}

	return nil
}

func (obj *objRN) onIConn(ctx context.Context, funcCancel cancelFunc) error {
	return obj.oncheck(ctx, funcCancel)
}

func (obj *objRN) onConnMe(ctx context.Context, funcCancel cancelFunc) error {
	return obj.oncheck(ctx, funcCancel)
}

func (obj *objRN) makeString() string {
	return fmt.Sprintf("root(%v %v) node1(%v %v), node2(%v %v) requestnodes %v rqnodes1 %v rqnodes2 %v",
		obj.rootni.numsIConn, obj.rootni.numsConnMe,
		obj.node1ni.numsIConn, obj.node1ni.numsConnMe,
		obj.node2ni.numsIConn, obj.node2ni.numsConnMe,
		obj.requestnodes, obj.rqnodes1, obj.rqnodes2)
}

func startTestNodeRN(ctx context.Context, cfgfilename string, ni *nodeinfoRN, obj *objRN, oniconn funconcallRN, onconnme funconcallRN) (JarvisNode, error) {
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

func TestRequestNode(t *testing.T) {
	rootcfg, err := LoadConfig("./test/test5010_reqnoderoot.yaml")
	if err != nil {
		t.Fatalf("TestRequestNode load config %v err is %v", "./test/test5010_reqnoderoot.yaml", err)

		return
	}

	InitJarvisCore(rootcfg, "testnode", "1.2.3")
	defer ReleaseJarvisCore()

	obj := newObjRN()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var errobj error

	oniconn := func(ctx context.Context, err error, obj *objRN) error {
		if err != nil {
			errobj = err

			cancel()

			return nil
		}

		err1 := obj.onIConn(ctx, func(err error) {
			if err != nil {
				errobj = err
			}

			cancel()
		})

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

	onconnme := func(ctx context.Context, err error, obj *objRN) error {
		if err != nil {
			errobj = err

			cancel()

			return nil
		}

		err1 := obj.onConnMe(ctx, func(err error) {
			if err != nil {
				errobj = err
			}

			cancel()
		})
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

	obj.root, err = startTestNodeRN(ctx, "./test/test5010_reqnoderoot.yaml", &obj.rootni, obj,
		oniconn, onconnme)
	if err != nil {
		t.Fatalf("TestRequestNode startTestNodeRN root err %v", err)

		return
	}

	obj.node1, err = startTestNodeRN(ctx, "./test/test5011_reqnode1.yaml", &obj.node1ni, obj,
		oniconn, onconnme)
	if err != nil {
		t.Fatalf("TestRequestNode startTestNodeRN node1 err %v", err)

		return
	}

	obj.node2, err = startTestNodeRN(ctx, "./test/test5012_reqnode2.yaml", &obj.node2ni, obj,
		oniconn, onconnme)
	if err != nil {
		t.Fatalf("TestRequestNode startTestNodeRN node2 err %v", err)

		return
	}

	go obj.root.Start(ctx)
	time.Sleep(time.Second * 1)
	go obj.node1.Start(ctx)
	go obj.node2.Start(ctx)

	<-ctx.Done()

	if errobj != nil {
		t.Fatalf("TestRequestNode err %v", errobj)

		return
	}

	if !obj.isDone() {
		t.Fatalf("TestRequestNode no done %v", obj.makeString())

		return
	}

	t.Logf("TestRequestNode OK")
}
