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

// ResultSendMsg - result for FuncOnSendMsgResult
type ResultSendMsg struct {
	Msgs []*pb.JarvisMsg
	Err  error
}

// FuncOnSendMsgResult - on sendmsg recv all the messages
type FuncOnSendMsgResult func(ctx context.Context, jarvisnode JarvisNode, lstResult []*ResultSendMsg) error

type clientTask struct {
	client       *jarvisClient2
	msg          *pb.JarvisMsg
	servaddr     string
	node         *coredbpb.NodeInfo
	funcOnResult FuncOnSendMsgResult
}

func (task *clientTask) Run(ctx context.Context) error {
	if task.msg == nil {
		err := task.client._connectNode(ctx, task.servaddr, task.funcOnResult)
		if err != nil && task.node != nil {
			if err == ErrServAddrIsMe || err == ErrInvalidServAddr {
				task.client.node.mgrEvent.onNodeEvent(ctx, EventOnDeprecateNode, task.node)
			}

			task.client.node.mgrEvent.onNodeEvent(ctx, EventOnIConnectNodeFail, task.node)
		}

		return err
	}

	return task.client._sendMsg(ctx, task.msg, task.funcOnResult)
}

type clientInfo2 struct {
	conn     *grpc.ClientConn
	client   pb.JarvisCoreServClient
	servAddr string
}

// jarvisClient2 -
type jarvisClient2 struct {
	pool      jarvisbase.RoutinePool
	node      *jarvisNode
	mapClient sync.Map
}

func newClient2(node *jarvisNode) *jarvisClient2 {
	return &jarvisClient2{
		node: node,
		pool: jarvisbase.NewRoutinePool(),
	}
}

// start - start goroutine to proc client task
func (c *jarvisClient2) start(ctx context.Context) error {
	return c.pool.Start(ctx, 128)
}

// addTask - add a client task
func (c *jarvisClient2) addTask(msg *pb.JarvisMsg, servaddr string, node *coredbpb.NodeInfo, funcOnResult FuncOnSendMsgResult) {
	task := &clientTask{
		msg:          msg,
		servaddr:     servaddr,
		client:       c,
		node:         node,
		funcOnResult: funcOnResult,
	}

	c.pool.SendTask(task)
}

func (c *jarvisClient2) isConnected(addr string) bool {
	_, ok := c.mapClient.Load(addr)

	return ok
}

func (c *jarvisClient2) _getValidClientConn(addr string) (*clientInfo2, error) {
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

func (c *jarvisClient2) _sendMsg(ctx context.Context, smsg *pb.JarvisMsg, funcOnResult FuncOnSendMsgResult) error {
	_, ok := c.mapClient.Load(smsg.DestAddr)

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

	var lstMsg []*pb.JarvisMsg

	for {
		getmsg, err := stream.Recv()
		if err == io.EOF {
			jarvisbase.Debug("jarvisClient2._sendMsg:stream eof")

			if funcOnResult != nil {
				funcOnResult(ctx, c.node, []*ResultSendMsg{
					&ResultSendMsg{
						Msgs: lstMsg,
						Err:  nil,
					},
				})
			}
			// c.node.mgrEvent.onNodeEvent(ctx, EventOnEndRequestNode, c.node.GetCoreDB().GetNode(smsg.DestAddr))

			break
		}

		if err != nil {
			jarvisbase.Warn("jarvisClient2._sendMsg:stream", zap.Error(err))

			if funcOnResult != nil {
				funcOnResult(ctx, c.node, []*ResultSendMsg{
					&ResultSendMsg{
						Msgs: lstMsg,
						Err:  err,
					},
				})
			}

			break
		} else {
			jarvisbase.Debug("jarvisClient2._sendMsg:stream", jarvisbase.JSON("msg", getmsg))

			lstMsg = append(lstMsg, getmsg)

			c.node.PostMsg(getmsg, nil, nil, funcOnResult)
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

func (c *jarvisClient2) _connectNode(ctx context.Context, servaddr string, funcOnResult FuncOnSendMsgResult) error {
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

	msg, err := BuildConnNode(c.node, c.node.GetMyInfo().Addr, "", servaddr, nbi)
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

	var lstMsg []*pb.JarvisMsg

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			jarvisbase.Debug("jarvisClient2._connectNode:stream eof")

			if funcOnResult != nil {
				funcOnResult(ctx, c.node, []*ResultSendMsg{
					&ResultSendMsg{
						Msgs: lstMsg,
						Err:  nil,
					},
				})
			}

			break
		}

		if err != nil {
			jarvisbase.Warn("jarvisClient2._connectNode:stream", zap.Error(err))

			if funcOnResult != nil {
				funcOnResult(ctx, c.node, []*ResultSendMsg{
					&ResultSendMsg{
						Msgs: lstMsg,
						Err:  err,
					},
				})
			}

			break
		} else {
			jarvisbase.Debug("jarvisClient2._connectNode:stream", jarvisbase.JSON("msg", msg))

			lstMsg = append(lstMsg, msg)

			if msg.MsgType == pb.MSGTYPE_REPLY_CONNECT {
				ni := msg.GetNodeInfo()

				c.mapClient.Store(ni.Addr, ci)
			}

			c.node.PostMsg(msg, nil, nil, funcOnResult)
		}
	}

	return nil
}
