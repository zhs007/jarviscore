package jarviscore

import (
	"context"
	"io"
	"net"
	"sync"

	"github.com/zhs007/jarviscore/coredb/proto"

	"go.uber.org/zap"

	"github.com/zhs007/jarviscore/base"
	pb "github.com/zhs007/jarviscore/proto"
	"google.golang.org/grpc"
)

type clientTask struct {
	client   *jarvisClient2
	msg      *pb.JarvisMsg
	servaddr string
	node     *coredbpb.NodeInfo
}

func (task *clientTask) Run(ctx context.Context) error {
	if task.msg == nil {
		err := task.client._connectNode(ctx, task.servaddr)
		if err != nil && task.node != nil {
			if err == ErrServAddrIsMe || err == ErrInvalidServAddr {
				task.client.node.mgrEvent.onNodeEvent(ctx, EventOnDeprecateNode, task.node)
			}

			task.client.node.mgrEvent.onNodeEvent(ctx, EventOnIConnectNodeFail, task.node)
		}

		return err
	}

	return task.client._sendMsg(ctx, task.msg)
}

type clientInfo2 struct {
	conn     *grpc.ClientConn
	client   pb.JarvisCoreServClient
	servAddr string
}

// jarvisClient2 -
type jarvisClient2 struct {
	// sync.RWMutex
	// sync.Map

	pool      jarvisbase.RoutinePool
	node      *jarvisNode
	mapClient sync.Map
	// mapClient map[string]*clientInfo2
}

func newClient2(node *jarvisNode) *jarvisClient2 {
	return &jarvisClient2{
		node: node,
		// mapClient: make(map[string]*clientInfo2),
		pool: jarvisbase.NewRoutinePool(),
	}
}

// start - start goroutine to proc client task
func (c *jarvisClient2) start(ctx context.Context) error {
	return c.pool.Start(ctx, 128)
}

// addTask - add a client task
func (c *jarvisClient2) addTask(msg *pb.JarvisMsg, servaddr string, node *coredbpb.NodeInfo) {
	task := &clientTask{
		msg:      msg,
		servaddr: servaddr,
		client:   c,
		node:     node,
	}

	c.pool.SendTask(task)
}

func (c *jarvisClient2) isConnected(addr string) bool {
	// c.Lock()
	// defer c.Unlock()

	_, ok := c.mapClient.Load(addr)
	return ok
}

func (c *jarvisClient2) _getValidClientConn(addr string) (*clientInfo2, error) {
	// c.Lock()
	// defer c.Unlock()

	mi, ok := c.mapClient.Load(addr)
	if ok {
		ci, ok := mi.(*clientInfo2)
		if ok {
			if mgrconn.isValidConn(ci.servAddr) {
				return ci, nil
			}
		}
	}

	ci, ok := mi.(*clientInfo2)
	if ok {
		if mgrconn.isValidConn(ci.servAddr) {
			return ci, nil
		}
	}

	conn, err := mgrconn.getConn(ci.servAddr)
	if err != nil {
		jarvisbase.Warn("jarvisClient2.getValidClientConn", zap.Error(err))

		return nil, err
	}

	nci := &clientInfo2{
		conn:     conn,
		client:   pb.NewJarvisCoreServClient(conn),
		servAddr: ci.servAddr,
	}

	c.mapClient.Store(addr, nci)

	return nci, nil
}

func (c *jarvisClient2) _sendMsg(ctx context.Context, smsg *pb.JarvisMsg) error {
	// c.Lock()
	// defer c.Unlock()

	_, ok := c.mapClient.Load(smsg.DestAddr)
	// c.Unlock()
	if !ok {
		return c._broadCastMsg(ctx, smsg)
	}

	jarvisbase.Debug("jarvisClient2._sendMsg", jarvisbase.JSON("msg", smsg))

	ci2, err := c._getValidClientConn(smsg.DestAddr)
	if err != nil {
		jarvisbase.Warn("jarvisClient2._sendMsg:getValidClientConn", zap.Error(err))

		return err
	}

	stream, err := ci2.client.ProcMsg(ctx, smsg)
	if err != nil {
		jarvisbase.Warn("jarvisClient2._sendMsg:ProcMsg", zap.Error(err))

		return err
	}

	for {
		getmsg, err := stream.Recv()
		if err == io.EOF {
			jarvisbase.Debug("jarvisClient2._sendMsg:stream eof")

			c.node.mgrEvent.onNodeEvent(ctx, EventOnEndRequestNode, c.node.GetCoreDB().GetNode(smsg.DestAddr))

			break
		}

		if err != nil {
			jarvisbase.Warn("jarvisClient2._sendMsg:stream", zap.Error(err))

			break
		} else {
			jarvisbase.Debug("jarvisClient2._sendMsg:stream", jarvisbase.JSON("msg", getmsg))

			c.node.PostMsg(getmsg, nil, nil)
		}
	}

	return nil
}

func (c *jarvisClient2) _broadCastMsg(ctx context.Context, msg *pb.JarvisMsg) error {
	jarvisbase.Debug("jarvisClient2._broadCastMsg", jarvisbase.JSON("msg", msg))

	// c.mapClient.Range(func(key, v interface{}) bool {
	// 	ci, ok := v.(*clientInfo2)
	// 	if ok {
	// 		stream, err := ci.client.ProcMsg(ctx, msg)
	// 		if err != nil {
	// 			jarvisbase.Warn("jarvisClient2._broadCastMsg:ProcMsg", zap.Error(err))

	// 			return true
	// 		}

	// 		for {
	// 			msg, err := stream.Recv()
	// 			if err == io.EOF {
	// 				jarvisbase.Debug("jarvisClient2._broadCastMsg:stream eof")

	// 				break
	// 			}

	// 			if err != nil {
	// 				jarvisbase.Warn("jarvisClient2._broadCastMsg:stream", zap.Error(err))

	// 				break
	// 			} else {
	// 				jarvisbase.Debug("jarvisClient2._broadCastMsg:stream", jarvisbase.JSON("msg", msg))

	// 				c.node.mgrJasvisMsg.sendMsg(msg, nil, nil)
	// 			}
	// 		}
	// 	}

	// 	return true
	// })

	return nil
}

func (c *jarvisClient2) _connectNode(ctx context.Context, servaddr string) error {
	_, _, err := net.SplitHostPort(servaddr)
	if err != nil {
		jarvisbase.Warn("jarvisClient2._connectNode:checkServAddr", zap.Error(err))

		return ErrInvalidServAddr
	}

	if IsMyServAddr(servaddr, c.node.myinfo.BindAddr) {
		jarvisbase.Warn("jarvisClient2._connectNode", zap.Error(ErrServAddrIsMe))

		return ErrServAddrIsMe
	}

	conn, err := mgrconn.getConn(servaddr)
	if err != nil {
		jarvisbase.Warn("jarvisClient2._connectNode", zap.Error(err))

		return err
	}

	curctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ci := &clientInfo2{
		conn:     conn,
		client:   pb.NewJarvisCoreServClient(conn),
		servAddr: servaddr,
	}

	nbi := &pb.NodeBaseInfo{
		ServAddr:        c.node.GetMyInfo().ServAddr,
		Addr:            c.node.GetMyInfo().Addr,
		Name:            c.node.GetMyInfo().Name,
		NodeType:        c.node.GetMyInfo().NodeType,
		NodeTypeVersion: c.node.GetMyInfo().NodeTypeVersion,
		CoreVersion:     c.node.GetMyInfo().CoreVersion,
	}

	msg, err := BuildConnNode(c.node.coredb.GetPrivateKey(), 0, c.node.GetMyInfo().Addr, "", servaddr, nbi)
	if err != nil {
		jarvisbase.Warn("jarvisClient2._connectNode:BuildConnNode", zap.Error(err))

		return err
	}

	stream, err1 := ci.client.ProcMsg(curctx, msg)
	if err1 != nil {
		jarvisbase.Warn("jarvisClient2._connectNode:ProcMsg", zap.Error(err1))

		err := conn.Close()
		if err != nil {
			jarvisbase.Warn("jarvisClient2._connectNode:Close", zap.Error(err1))
		}

		mgrconn.delConn(servaddr)

		return err1
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			jarvisbase.Debug("jarvisClient2._connectNode:stream eof")

			break
		}

		if err != nil {
			jarvisbase.Warn("jarvisClient2._connectNode:stream", zap.Error(err))

			break
		} else {
			jarvisbase.Debug("jarvisClient2._connectNode:stream", jarvisbase.JSON("msg", msg))

			if msg.MsgType == pb.MSGTYPE_REPLY_CONNECT {
				ni := msg.GetNodeInfo()

				c.mapClient.Store(ni.Addr, ci)
				// } else if msg.MsgType == pb.MSGTYPE_REPLY_CONNECT2 {
				// 	rc2 := msg.GetReplyConn2()

				// 	c.mapClient.Store(rc2.Nbi.Addr, ci)
			}

			c.node.PostMsg(msg, nil, nil)
		}
	}

	return nil
}
