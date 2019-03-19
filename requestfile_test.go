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

func randfillFile(fn string, len int) error {
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

func outputErrRF(t *testing.T, err error, msg string, info string) {
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

func outputRF(t *testing.T, msg string) {
	jarvisbase.Info(msg)
	t.Logf(msg)
}

// funconcallRF
type funconcallRF func(ctx context.Context, err error, obj *objRF) error

type mapnodeinfoRF struct {
	iconn  bool
	connme bool
}

type nodeinfoRF struct {
	mapAddr    sync.Map
	numsIConn  int
	numsConnMe int
}

func (ni *nodeinfoRF) onIConnectNode(node *coredbpb.NodeInfo) error {
	d, ok := ni.mapAddr.Load(node.Addr)
	if ok {
		mni, ok := d.(*mapnodeinfoRF)
		if !ok {
			return fmt.Errorf("nodeinfoRF.onIConnectNode:mapAddr2mapnodeinfo err")
		}

		if !mni.iconn {
			mni.iconn = true

			ni.numsIConn++
		}

		return nil
	}

	mni := &mapnodeinfoRF{
		iconn: true,
	}

	ni.mapAddr.Store(node.Addr, mni)
	ni.numsIConn++

	return nil
}

func (ni *nodeinfoRF) onNodeConnected(node *coredbpb.NodeInfo) error {
	d, ok := ni.mapAddr.Load(node.Addr)
	if ok {
		mni, ok := d.(*mapnodeinfoRF)
		if !ok {
			return fmt.Errorf("nodeinfoRF.onNodeConnected:mapAddr2mapnodeinfo err")
		}

		if !mni.connme {
			mni.connme = true

			ni.numsConnMe++
		}

		return nil
	}

	mni := &mapnodeinfoRF{
		connme: true,
	}

	ni.mapAddr.Store(node.Addr, mni)
	ni.numsConnMe++

	return nil
}

type objRF struct {
	root           JarvisNode
	node1          JarvisNode
	node2          JarvisNode
	rootni         nodeinfoRF
	node1ni        nodeinfoRF
	node2ni        nodeinfoRF
	requestnodes   bool
	requestfile1   bool
	requestfile1ok bool
	requestfile2   bool
	requestfile2ok bool
	err            error
}

func newObjRF() *objRF {
	return &objRF{
		rootni:       nodeinfoRF{},
		node1ni:      nodeinfoRF{},
		node2ni:      nodeinfoRF{},
		requestnodes: false,
	}
}

func (obj *objRF) isDone() bool {
	if obj.rootni.numsConnMe != 2 || obj.rootni.numsIConn != 2 {
		return false
	}

	return obj.requestfile1 && obj.requestfile1ok && obj.requestfile2 && obj.requestfile2ok
}

func (obj *objRF) oncheckrequestfile2(ctx context.Context, funcCancel context.CancelFunc) error {
	if obj.requestfile1 && obj.requestfile1ok &&
		!obj.requestfile2 {

		curresultnums := 0

		rf := &jarviscorepb.RequestFile{
			Filename: "./test/rf002.dat",
		}

		err := obj.node2.RequestFile(ctx, obj.node1.GetMyInfo().Addr, rf,
			func(ctx context.Context, jarvisnode JarvisNode,
				lstResult []*ClientProcMsgResult) error {

				for i := 0; i < len(lstResult); i++ {
					if lstResult[i].Msg != nil {
						jarvisbase.Info("obj.node2.RequestFile", JSONMsg2Zap("result", lstResult[i].Msg))
					}
				}

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

						if lstResult[curresultnums].Err == nil && lstResult[curresultnums].Msg == nil {

							var lst []*jarviscorepb.FileData
							totalmd5 := ""
							totallen := int64(0)
							for i := 0; i < len(lstResult); i++ {
								if lstResult[i].Msg != nil && lstResult[i].Msg.GetFile() != nil {
									lst = append(lst, lstResult[i].Msg.GetFile())

									if lstResult[i].Msg.GetFile().FileMD5String != "" {
										totalmd5 = lstResult[i].Msg.GetFile().FileMD5String
									}

									totallen = lstResult[i].Msg.GetFile().TotalLength
								}
							}

							tl, md5str, err := CountMD5String(lst)
							if err != nil {
								obj.err = err

								funcCancel()

								return nil
							}

							if tl != totallen {
								obj.err = fmt.Errorf("obj.node2.RequestFile length %v %v", tl, totallen)

								funcCancel()

								return nil
							}

							if md5str != totalmd5 {
								obj.err = fmt.Errorf("obj.node2.RequestFile md5check %v %v", md5str, totalmd5)

								funcCancel()

								return nil
							}

							obj.requestfile2ok = true

							funcCancel()

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

func (obj *objRF) oncheck(ctx context.Context, funcCancel context.CancelFunc) error {
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
		obj.node1ni.numsIConn == 2 && obj.node2ni.numsIConn == 2 && !obj.requestfile1 {

		curresultnums := 0

		rf := &jarviscorepb.RequestFile{
			Filename: "./test/rf001.dat",
		}

		err := obj.node1.RequestFile(ctx, obj.node2.GetMyInfo().Addr, rf,
			func(ctx context.Context, jarvisnode JarvisNode,
				lstResult []*ClientProcMsgResult) error {

				for i := 0; i < len(lstResult); i++ {
					if lstResult[i].Msg != nil {
						jarvisbase.Info("obj.node1.RequestFile", JSONMsg2Zap("result", lstResult[i].Msg))
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

						if lstResult[curresultnums].Err == nil && lstResult[curresultnums].Msg == nil {
							obj.requestfile1ok = true

							obj.oncheckrequestfile2(ctx, funcCancel)

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

	return nil
}

func (obj *objRF) onIConn(ctx context.Context, funcCancel context.CancelFunc) error {
	return obj.oncheck(ctx, funcCancel)
}

func (obj *objRF) onConnMe(ctx context.Context, funcCancel context.CancelFunc) error {
	return obj.oncheck(ctx, funcCancel)
}

func (obj *objRF) makeString() string {
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

func startTestNodeRF(ctx context.Context, cfgfilename string, ni *nodeinfoRF, obj *objRF,
	oniconn funconcallRF, onconnme funconcallRF) (JarvisNode, error) {

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

func TestRequestFile(t *testing.T) {
	randfillFile("./test/rf001.dat", 2*1024*1024)
	randfillFile("./test/rf002.dat", 10*1024*1024)

	rootcfg, err := LoadConfig("./test/test5020_reqfileroot.yaml")
	if err != nil {
		t.Fatalf("TestRequestFile load config %v err is %v", "./test/test5020_reqfileroot.yaml", err)

		return
	}

	InitJarvisCore(rootcfg)
	defer ReleaseJarvisCore()

	obj := newObjRF()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var errobj error

	oniconn := func(ctx context.Context, err error, obj *objRF) error {
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

	onconnme := func(ctx context.Context, err error, obj *objRF) error {
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

	obj.root, err = startTestNodeRF(ctx, "./test/test5020_reqfileroot.yaml", &obj.rootni, obj,
		oniconn, onconnme)
	if err != nil {
		outputErrRF(t, err, "TestRequestFile startTestNodeRF root", "")

		return
	}

	obj.node1, err = startTestNodeRF(ctx, "./test/test5021_reqfile1.yaml", &obj.node1ni, obj,
		oniconn, onconnme)
	if err != nil {
		outputErrRF(t, err, "TestRequestFile startTestNodeRF node1", "")

		return
	}

	obj.node2, err = startTestNodeRF(ctx, "./test/test5022_reqfile2.yaml", &obj.node2ni, obj,
		oniconn, onconnme)
	if err != nil {
		outputErrRF(t, err, "TestRequestFile startTestNodeRF node2", "")

		return
	}

	go obj.root.Start(ctx)
	time.Sleep(time.Second * 1)
	go obj.node1.Start(ctx)
	go obj.node2.Start(ctx)

	<-ctx.Done()

	if errobj != nil {
		outputErrRF(t, errobj, "TestRequestFile", "")

		return
	}

	if obj.err != nil {
		outputErrRF(t, obj.err, "TestRequestFile", "")

		return
	}

	if !obj.isDone() {
		outputErrRF(t, nil, "TestRequestFile no done", obj.makeString())

		return
	}

	outputRF(t, "TestRequestFile OK")
}
