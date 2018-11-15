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

	jarvisbase.Debug("jarvisClient2.sendMsg", jarvisbase.JSON("msg", msg))

	ci2.client.ProcMsg(ctx, msg)

	return nil
}

func (c *jarvisClient2) broadCastMsg(ctx context.Context, msg *pb.JarvisMsg) error {
	jarvisbase.Debug("jarvisClient2.broadCastMsg", jarvisbase.JSON("msg", msg))

	for _, v := range c.mapClient {
		v.client.ProcMsg(ctx, msg)
	}

	return nil
}

func (c *jarvisClient2) connectNode(ctx context.Context, servaddr string) error {
	if servaddr == c.node.myinfo.ServAddr {
		jarvisbase.Debug("jarvisClient2.connectNode", zap.Error(ErrServAddrIsMe))

		return ErrServAddrIsMe
	}

	conn, err := mgrconn.getConn(servaddr)
	if err != nil {
		jarvisbase.Debug("jarvisClient2.connectNode", zap.Error(err))

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
		jarvisbase.Debug("jarvisClient2.connectNode:BuildConnNode", zap.Error(err))

		return err
	}
	// jarvisbase.Debug("jarvisClient2.connectNode:ProcMsg", jarvisbase.JSON("msg", msg))

	// SignJarvisMsg(c.node.coredb.privKey, msg)

	// jarvisbase.Debug("jarvisClient2.connectNode:ProcMsg", jarvisbase.JSON("msg", msg))

	stream, err1 := ci.client.ProcMsg(curctx, msg)
	if err1 != nil {
		jarvisbase.Debug("jarvisClient2.connectNode:ProcMsg", zap.Error(err1))

		mgrconn.delConn(servaddr)

		return err1
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			jarvisbase.Debug("jarvisClient2.connectNode:stream eof")

			break
		}

		if err != nil {
			jarvisbase.Debug("jarvisClient2.connectNode:stream", zap.Error(err))

			break
		} else {
			jarvisbase.Debug("jarvisClient2.connectNode:stream", jarvisbase.JSON("msg", msg))

			if msg.MsgType == pb.MSGTYPE_REPLY_CONNECT {
				ni := msg.GetNodeInfo()

				c.mapClient[ni.Addr] = ci
			}

			c.node.mgrJasvisMsg.sendMsg(msg, nil, nil)
		}
	}

	return nil
}
