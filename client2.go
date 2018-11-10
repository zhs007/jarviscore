package jarviscore

import (
	"context"
	"io"

	"go.uber.org/zap"

	"github.com/zhs007/jarviscore/base"
	pb "github.com/zhs007/jarviscore/proto"
	"google.golang.org/grpc"
)

type clientInfo2 struct {
	conn   *grpc.ClientConn
	client pb.JarvisCoreServClient
}

// jarvisClient2 -
type jarvisClient2 struct {
	node      *jarvisNode
	mapClient map[string]*clientInfo2
}

func newClient2(node *jarvisNode) *jarvisClient2 {
	return &jarvisClient2{
		node:      node,
		mapClient: make(map[string]*clientInfo2),
	}
}

func (c *jarvisClient2) sendMsg(ctx context.Context, msg *pb.JarvisMsg) error {
	ci2, ok := c.mapClient[msg.DestAddr]
	if !ok {
		return c.broadCastMsg(ctx, msg)
	}

	ci2.client.ProcMsg(ctx, msg)

	return nil
}

func (c *jarvisClient2) broadCastMsg(ctx context.Context, msg *pb.JarvisMsg) error {
	for _, v := range c.mapClient {
		v.client.ProcMsg(ctx, msg)
	}

	return nil
}

func (c *jarvisClient2) connectNode(ctx context.Context, ni *pb.NodeBaseInfo) error {
	if ni.ServAddr == c.node.myinfo.ServAddr {
		jarvisbase.Debug("jarvisClient2.connectNode", zap.Error(ErrServAddrIsMe))

		return ErrServAddrIsMe
	}

	conn, err := mgrconn.getConn(ni.ServAddr)
	if err != nil {
		jarvisbase.Debug("jarvisClient2.connectNode", zap.Error(err))

		return err
	}

	curctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ci := clientInfo2{
		conn:   conn,
		client: pb.NewJarvisCoreServClient(conn),
	}

	nbi := &pb.NodeBaseInfo{
		ServAddr: c.node.GetMyInfo().ServAddr,
		Addr:     c.node.GetMyInfo().Addr,
		Name:     c.node.GetMyInfo().Name,
	}

	msg := BuildConnNode(0, c.node.GetCoreDB().addr, ni.Addr, nbi)
	SignJarvisMsg(c.node.coredb.privKey, msg)

	stream, err1 := ci.client.ProcMsg(curctx, msg)
	if err1 != nil {
		jarvisbase.Debug("jarvisClient2.connectNode:ProcMsg", zap.Error(err1))

		mgrconn.delConn(ni.ServAddr)

		return err1
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			jarvisbase.Debug("jarvisClient2.connectNode:stream", zap.Error(err))
		} else {
			c.node.mgrJasvisMsg.sendMsg(msg, nil)
		}
	}

	return nil
}

func (c *jarvisClient2) connectRoot(ctx context.Context, servaddr string) error {
	if servaddr == c.node.myinfo.ServAddr {
		jarvisbase.Debug("jarvisClient2.connectRoot", zap.Error(ErrServAddrIsMe))

		return ErrServAddrIsMe
	}

	conn, err := mgrconn.getConn(servaddr)
	if err != nil {
		jarvisbase.Debug("jarvisClient2.connectRoot", zap.Error(err))

		return err
	}

	curctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ci := clientInfo2{
		conn:   conn,
		client: pb.NewJarvisCoreServClient(conn),
	}

	nbi := &pb.NodeBaseInfo{
		ServAddr: c.node.GetMyInfo().ServAddr,
		Addr:     c.node.GetMyInfo().Addr,
		Name:     c.node.GetMyInfo().Name,
	}

	msg := BuildConnectRoot(0, c.node.GetCoreDB().addr, "", nbi)
	SignJarvisMsg(c.node.coredb.privKey, msg)

	stream, err1 := ci.client.ProcMsg(curctx, msg)
	if err1 != nil {
		jarvisbase.Debug("jarvisClient2.connectRoot:ProcMsg", zap.Error(err1))

		mgrconn.delConn(servaddr)

		return err1
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			jarvisbase.Debug("jarvisClient2.connectRoot:stream", zap.Error(err))
		} else {
			c.node.mgrJasvisMsg.sendMsg(msg, nil)
		}
	}

	return nil
}
