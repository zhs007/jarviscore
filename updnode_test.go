package jarviscore

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/zhs007/jarviscore/base"
	"github.com/zhs007/jarviscore/coredb/proto"

	"go.uber.org/zap"
)

func outputErrUN(t *testing.T, err error, msg string, info string) {
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

func outputUN(t *testing.T, msg string) {
	jarvisbase.Info(msg)
	t.Logf(msg)
}

// funconcallUN
type funconcallUN func(ctx context.Context, err error, obj *objUN) error

type mapnodeinfoUN struct {
	iconn  bool
	connme bool
}

type nodeinfoUN struct {
	mapAddr    sync.Map
	numsIConn  int
	numsConnMe int
}

func (ni *nodeinfoUN) onIConnectNode(node *coredbpb.NodeInfo) error {
	d, ok := ni.mapAddr.Load(node.Addr)
	if ok {
		mni, ok := d.(*mapnodeinfoUN)
		if !ok {
			return fmt.Errorf("nodeinfoUN.onIConnectNode:mapAddr2mapnodeinfo err")
		}

		if !mni.iconn {
			mni.iconn = true

			ni.numsIConn++
		}

		return nil
	}

	mni := &mapnodeinfoUN{
		iconn: true,
	}

	ni.mapAddr.Store(node.Addr, mni)
	ni.numsIConn++

	return nil
}

func (ni *nodeinfoUN) onNodeConnected(node *coredbpb.NodeInfo) error {
	d, ok := ni.mapAddr.Load(node.Addr)
	if ok {
		mni, ok := d.(*mapnodeinfoUN)
		if !ok {
			return fmt.Errorf("nodeinfoUN.onNodeConnected:mapAddr2mapnodeinfo err")
		}

		if !mni.connme {
			mni.connme = true

			ni.numsConnMe++
		}

		return nil
	}

	mni := &mapnodeinfoUN{
		connme: true,
	}

	ni.mapAddr.Store(node.Addr, mni)
	ni.numsConnMe++

	return nil
}

type objUN struct {
	root             JarvisNode
	node1            JarvisNode
	node2            JarvisNode
	rootni           nodeinfoUN
	node1ni          nodeinfoUN
	node2ni          nodeinfoUN
	requestnodes     bool
	updnodefromnode1 bool
	endupdnodes      bool
}

func newObjUN() *objUN {
	return &objUN{
		rootni:       nodeinfoUN{},
		node1ni:      nodeinfoUN{},
		node2ni:      nodeinfoUN{},
		requestnodes: false,
	}
}

func (obj *objUN) isDone() bool {
	if obj.rootni.numsConnMe != 2 || obj.rootni.numsIConn != 2 {
		return false
	}

	return obj.endupdnodes
}

func (obj *objUN) onIConn(ctx context.Context, funcCancel context.CancelFunc) error {
	if obj.rootni.numsConnMe == 2 && !obj.requestnodes {
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

	if obj.node1ni.numsConnMe == 2 && obj.node2ni.numsConnMe == 2 && !obj.updnodefromnode1 {
		err := obj.node1.UpdateAllNodes(ctx, "testupdnode", "0.7.25", nil,
			func(ctx context.Context, jarvisnode JarvisNode, numsNode int, lstResult []*ClientGroupProcMsgResults) error {

				jarvisbase.Info("objUN.onIConn:node1.UpdateAllNodes",
					zap.Int("numsNode", numsNode),
					jarvisbase.JSON("lstResult", lstResult))

				if numsNode == 2 && len(lstResult) == 2 && len(lstResult[0].Results) > 0 && len(lstResult[1].Results) > 0 &&
					lstResult[0].Results[len(lstResult[0].Results)-1].Msg == nil &&
					lstResult[1].Results[len(lstResult[1].Results)-1].Msg == nil {

					if (len(lstResult[0].Results) == 2 && len(lstResult[1].Results) == 4) ||
						(len(lstResult[1].Results) == 2 && len(lstResult[0].Results) == 4) {

						obj.endupdnodes = true

						funcCancel()

						return nil
					}
				}

				return nil
			})
		if err != nil {
			return err
		}

		obj.updnodefromnode1 = true
	}

	return nil
}

func (obj *objUN) onConnMe(ctx context.Context) error {
	return nil
}

func (obj *objUN) makeString() string {
	return fmt.Sprintf("root(%v %v) node1(%v %v), node2(%v %v)",
		obj.rootni.numsIConn, obj.rootni.numsConnMe,
		obj.node1ni.numsIConn, obj.node1ni.numsConnMe,
		obj.node2ni.numsIConn, obj.node2ni.numsConnMe)
}

func startTestNodeUN(ctx context.Context, cfgfilename string, ni *nodeinfoUN, obj *objUN, oniconn funconcallUN, onconnme funconcallUN) (JarvisNode, error) {
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

	obj := newObjUN()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var errobj error

	oniconn := func(ctx context.Context, err error, obj *objUN) error {
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

	onconnme := func(ctx context.Context, err error, obj *objUN) error {
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

	obj.root, err = startTestNodeUN(ctx, "./test/test5000_updnoderoot.yaml", &obj.rootni, obj,
		oniconn, onconnme)
	if err != nil {
		outputErrUN(t, err, "TestUpdNode startTestNodeUN root", "")

		return
	}

	obj.node1, err = startTestNodeUN(ctx, "./test/test5001_updnode1.yaml", &obj.node1ni, obj,
		oniconn, onconnme)
	if err != nil {
		outputErrUN(t, err, "TestUpdNode startTestNodeUN node1", "")

		return
	}

	obj.node2, err = startTestNodeUN(ctx, "./test/test5002_updnode2.yaml", &obj.node2ni, obj,
		oniconn, onconnme)
	if err != nil {
		outputErrUN(t, err, "TestUpdNode startTestNodeUN node2", "")

		return
	}

	go obj.root.Start(ctx)
	go obj.node1.Start(ctx)
	go obj.node2.Start(ctx)

	<-ctx.Done()

	if errobj != nil {
		outputErrUN(t, errobj, "TestUpdNode", "")

		return
	}

	if !obj.isDone() {
		outputErrUN(t, nil, "TestUpdNode no done", obj.makeString())

		return
	}

	outputUN(t, "TestUpdNode OK")
}
