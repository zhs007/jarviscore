package jarviscore

import (
	"context"
	"io"
	"net"

	"go.uber.org/zap"

	"github.com/zhs007/jarviscore/base"
	pb "github.com/zhs007/jarviscore/proto"
	"google.golang.org/grpc"
)

type clientTask struct {
	client   *jarvisClient2
	msg      *pb.JarvisMsg
	servaddr string
}

func (task *clientTask) Run(ctx context.Context) error {
	if task.msg == nil {
		return task.client._connectNode(ctx, task.servaddr)
	}

	return task.client._sendMsg(ctx, task.msg)
}

type clientInfo2 struct {
	conn   *grpc.ClientConn
	client pb.JarvisCoreServClient
}

// jarvisClient2 -
type jarvisClient2 struct {
	pool      jarvisbase.RoutinePool
	node      *jarvisNode
	mapClient map[string]*clientInfo2
}

func newClient2(node *jarvisNode) *jarvisClient2 {
	return &jarvisClient2{
		node:      node,
		mapClient: make(map[string]*clientInfo2),
		pool:      jarvisbase.NewRoutinePool(),
	}
}

// start - start goroutine to proc client task
func (c *jarvisClient2) start(ctx context.Context) error {
	return c.pool.Start(ctx, 128)
}

// addTask - add a client task
func (c *jarvisClient2) addTask(msg *pb.JarvisMsg, servaddr string) {
	task := &clientTask{
		msg:      msg,
		servaddr: servaddr,
		client:   c,
	}

	c.pool.SendTask(task)
}

func (c *jarvisClient2) isConnected(addr string) bool {
	_, ok := c.mapClient[addr]
	return ok
}

func (c *jarvisClient2) _sendMsg(ctx context.Context, msg *pb.JarvisMsg) error {
	ci2, ok := c.mapClient[msg.DestAddr]
	if !ok {
		return c._broadCastMsg(ctx, msg)
	}

	jarvisbase.Debug("jarvisClient2._sendMsg", jarvisbase.JSON("msg", msg))

	ci2.client.ProcMsg(ctx, msg)

	return nil
}

func (c *jarvisClient2) _broadCastMsg(ctx context.Context, msg *pb.JarvisMsg) error {
	jarvisbase.Debug("jarvisClient2._broadCastMsg", jarvisbase.JSON("msg", msg))

	for _, v := range c.mapClient {
		v.client.ProcMsg(ctx, msg)
	}

	return nil
}

func (c *jarvisClient2) _connectNode(ctx context.Context, servaddr string) error {
	_, _, err := net.SplitHostPort(servaddr)
	if err != nil {
		jarvisbase.Debug("jarvisClient2._connectNode:checkServAddr", zap.Error(err))

		return err
	}

	if IsMyServAddr(servaddr, c.node.myinfo.BindAddr) {
		jarvisbase.Debug("jarvisClient2._connectNode", zap.Error(ErrServAddrIsMe))

		return ErrServAddrIsMe
	}

	conn, err := mgrconn.getConn(servaddr)
	if err != nil {
		jarvisbase.Debug("jarvisClient2._connectNode", zap.Error(err))

		return err
	}

	curctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ci := &clientInfo2{
		conn:   conn,
		client: pb.NewJarvisCoreServClient(conn),
	}

	nbi := &pb.NodeBaseInfo{
		ServAddr: c.node.GetMyInfo().ServAddr,
		Addr:     c.node.GetMyInfo().Addr,
		Name:     c.node.GetMyInfo().Name,
	}

	msg, err := BuildConnNode(c.node.coredb.privKey, 0, c.node.GetMyInfo().Addr, "", servaddr, nbi)
	if err != nil {
		jarvisbase.Debug("jarvisClient2._connectNode:BuildConnNode", zap.Error(err))

		return err
	}

	stream, err1 := ci.client.ProcMsg(curctx, msg)
	if err1 != nil {
		jarvisbase.Debug("jarvisClient2._connectNode:ProcMsg", zap.Error(err1))

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
			jarvisbase.Debug("jarvisClient2._connectNode:stream", zap.Error(err))

			break
		} else {
			jarvisbase.Debug("jarvisClient2._connectNode:stream", jarvisbase.JSON("msg", msg))

			if msg.MsgType == pb.MSGTYPE_REPLY_CONNECT {
				ni := msg.GetNodeInfo()

				c.mapClient[ni.Addr] = ci
			}

			c.node.mgrJasvisMsg.sendMsg(msg, nil, nil)
		}
	}

	return nil
}
