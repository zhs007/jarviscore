package jarviscore

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/zhs007/jarviscore/coredb/proto"
)

// funconcall
type funconcall func(ctx context.Context, err error, obj *updnodeobj) error

type mapnodeinfo struct {
	iconn  bool
	connme bool
}

type updnodenodeinfo struct {
	mapAddr    sync.Map
	numsIConn  int
	numsConnMe int
}

func (ni *updnodenodeinfo) onIConnectNode(node *coredbpb.NodeInfo) error {
	d, ok := ni.mapAddr.Load(node.Addr)
	if ok {
		mni, ok := d.(*mapnodeinfo)
		if !ok {
			return fmt.Errorf("updnodenodeinfo.onIConnectNode:mapAddr2mapnodeinfo err")
		}

		if !mni.iconn {
			mni.iconn = true

			ni.numsIConn++
		}

		return nil
	}

	mni := &mapnodeinfo{
		iconn: true,
	}

	ni.mapAddr.Store(node.Addr, mni)
	ni.numsIConn++

	return nil
}

func (ni *updnodenodeinfo) onNodeConnected(node *coredbpb.NodeInfo) error {
	d, ok := ni.mapAddr.Load(node.Addr)
	if ok {
		mni, ok := d.(*mapnodeinfo)
		if !ok {
			return fmt.Errorf("updnodenodeinfo.onNodeConnected:mapAddr2mapnodeinfo err")
		}

		if !mni.connme {
			mni.connme = true

			ni.numsConnMe++
		}

		return nil
	}

	mni := &mapnodeinfo{
		connme: true,
	}

	ni.mapAddr.Store(node.Addr, mni)
	ni.numsConnMe++

	return nil
}

type updnodeobj struct {
	root             JarvisNode
	node1            JarvisNode
	node2            JarvisNode
	rootni           updnodenodeinfo
	node1ni          updnodenodeinfo
	node2ni          updnodenodeinfo
	requestnodes     bool
	updnodefromnode1 bool
	endupdnodes      bool
}

func newUpdNodeObj() *updnodeobj {
	return &updnodeobj{
		rootni:       updnodenodeinfo{},
		node1ni:      updnodenodeinfo{},
		node2ni:      updnodenodeinfo{},
		requestnodes: false,
	}
}

func (obj *updnodeobj) isDone() bool {
	if obj.rootni.numsConnMe != 2 || obj.rootni.numsIConn != 2 {
		return false
	}

	return obj.endupdnodes
}

func (obj *updnodeobj) onIConn(ctx context.Context) error {
	if obj.rootni.numsConnMe == 2 && !obj.requestnodes {
		err := obj.node1.RequestNodes(nil)
		if err != nil {
			return err
		}

		err = obj.node2.RequestNodes(nil)
		if err != nil {
			return err
		}

		obj.requestnodes = true
	}

	if obj.node1ni.numsConnMe == 2 && obj.node2ni.numsConnMe == 2 && !obj.updnodefromnode1 {
		err := obj.node1.UpdateAllNodes(ctx, "testupdnode", "0.7.25", nil,
			func(ctx context.Context, jarvisnode JarvisNode, lstResult []*ResultSendMsg) error {
				if len(lstResult) != 2 {
					return fmt.Errorf("UpdateAllNodes:len %v", len(lstResult))
				}

				obj.endupdnodes = true

				return nil
			})
		if err != nil {
			return err
		}

		obj.updnodefromnode1 = true
	}

	return nil
}

func (obj *updnodeobj) onConnMe(ctx context.Context) error {
	return nil
}

func (obj *updnodeobj) makeString() string {
	return fmt.Sprintf("root(%v %v) node1(%v %v), node2(%v %v)",
		obj.rootni.numsIConn, obj.rootni.numsConnMe,
		obj.node1ni.numsIConn, obj.node1ni.numsConnMe,
		obj.node2ni.numsIConn, obj.node2ni.numsConnMe)
}

func startTestNode(ctx context.Context, cfgfilename string, ni *updnodenodeinfo, obj *updnodeobj, oniconn funconcall, onconnme funconcall) (JarvisNode, error) {
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

func TestUpdNode(t *testing.T) {
	rootcfg, err := LoadConfig("./test/test5000_updnoderoot.yaml")
	if err != nil {
		t.Fatalf("TestUpdNode load config %v err is %v", "./test/test5000_updnoderoot.yaml", err)

		return
	}

	InitJarvisCore(rootcfg)
	defer ReleaseJarvisCore()

	obj := newUpdNodeObj()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var errobj error

	oniconn := func(ctx context.Context, err error, obj *updnodeobj) error {
		if err != nil {
			errobj = err

			cancel()

			return nil
		}

		err1 := obj.onIConn(ctx)
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

	onconnme := func(ctx context.Context, err error, obj *updnodeobj) error {
		if err != nil {
			errobj = err

			cancel()

			return nil
		}

		err1 := obj.onConnMe(ctx)
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

	obj.root, err = startTestNode(ctx, "./test/test5000_updnoderoot.yaml", &obj.rootni, obj,
		oniconn, onconnme)
	if err != nil {
		t.Fatalf("TestUpdNode startTestNode root err %v", err)

		return
	}

	obj.node1, err = startTestNode(ctx, "./test/test5001_updnode1.yaml", &obj.node1ni, obj,
		oniconn, onconnme)
	if err != nil {
		t.Fatalf("TestUpdNode startTestNode node1 err %v", err)

		return
	}

	obj.node2, err = startTestNode(ctx, "./test/test5002_updnode2.yaml", &obj.node2ni, obj,
		oniconn, onconnme)
	if err != nil {
		t.Fatalf("TestUpdNode startTestNode node2 err %v", err)

		return
	}

	go obj.root.Start(ctx)
	go obj.node1.Start(ctx)
	go obj.node2.Start(ctx)

	<-ctx.Done()

	if errobj != nil {
		t.Fatalf("TestUpdNode err %v", errobj)

		return
	}

	if !obj.isDone() {
		t.Fatalf("TestUpdNode no done %v", obj.makeString())

		return
	}

	t.Logf("TestRequestNodes OK")
}
