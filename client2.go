package jarviscore

import (
	"context"
	"io"
	"sync"
	"time"

	"google.golang.org/grpc/codes"

	"github.com/zhs007/jarviscore/coredb"

	"github.com/zhs007/jarviscore/coredb/proto"

	"go.uber.org/zap"

	"github.com/zhs007/jarviscore/base"
	pb "github.com/zhs007/jarviscore/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type clientTask struct {
	jarvisbase.L2BaseTask

	servaddr     string
	addr         string
	client       *jarvisClient2
	msg          *pb.JarvisMsg
	msgs         []*pb.JarvisMsg
	node         *coredbpb.NodeInfo
	funcOnResult FuncOnProcMsgResult
}

func (task *clientTask) Run(ctx context.Context) error {

	if task.msgs != nil {
		return task.client._sendMsgStream(ctx, task.addr, task.msgs, task.funcOnResult)
	}

	if task.msg != nil {
		return task.client._sendMsg(ctx, task.msg, task.funcOnResult)
	}

	err := task.client._connectNode(ctx, task.servaddr, task.node, task.funcOnResult)
	if err != nil && task.node != nil {
		if err == ErrServAddrIsMe || err == ErrInvalidServAddr {
			task.client.node.mgrEvent.onNodeEvent(ctx, EventOnDeprecateNode, task.node)
		}

		task.client.node.mgrEvent.onNodeEvent(ctx, EventOnIConnectNodeFail, task.node)
	}

	if err != nil {
		if task.node != nil {
			jarvisbase.Warn("clientTask.Run:_connectNode",
				zap.Error(err),
				zap.String("servaddr", task.servaddr),
				jarvisbase.JSON("node", task.node))
		} else {
			jarvisbase.Warn("clientTask.Run:_connectNode",
				zap.Error(err),
				zap.String("servaddr", task.servaddr))
		}
	}

	task.client.onConnTaskEnd(task.servaddr)

	return err
}

// // GetParentID - get parentID
// func (task *clientTask) GetParentID() string {
// 	if task.msg == nil {
// 		return task.addr
// 	}

// 	return ""
// }

type clientInfo2 struct {
	conn     *grpc.ClientConn
	client   pb.JarvisCoreServClient
	servAddr string
}

// jarvisClient2 -
type jarvisClient2 struct {
	poolMsg         jarvisbase.L2RoutinePool
	poolConn        jarvisbase.RoutinePool
	node            *jarvisNode
	mapClient       sync.Map
	fsa             *failservaddr
	mapConnServAddr sync.Map
}

func newClient2(node *jarvisNode) *jarvisClient2 {
	return &jarvisClient2{
		node:     node,
		poolConn: jarvisbase.NewRoutinePool(),
		poolMsg:  jarvisbase.NewL2RoutinePool(),
		fsa:      newFailServAddr(),
	}
}

// BuildStatus - build status
func (c *jarvisClient2) BuildNodeStatus(ns *pb.JarvisNodeStatus) {
	ns.MsgPool = c.poolMsg.BuildStatus()

	c.mapClient.Range(func(key, v interface{}) bool {
		addr, keyok := key.(string)
		// ci, vok := v.(*clientInfo2)
		if keyok {
			ni := c.node.GetCoreDB().GetNode(addr)
			if ni != nil {
				nbi := &pb.NodeBaseInfo{
					ServAddr:        ni.ServAddr,
					Addr:            ni.Addr,
					Name:            ni.Name,
					NodeTypeVersion: ni.NodeTypeVersion,
					NodeType:        ni.NodeType,
					CoreVersion:     ni.CoreVersion,
				}

				ns.LstConnected = append(ns.LstConnected, nbi)
			}
		}

		return true
	})
}

// start - start goroutine to proc client task
func (c *jarvisClient2) start(ctx context.Context) error {
	go c.poolConn.Start(ctx, 128)
	go c.poolMsg.Start(ctx, 128)

	<-ctx.Done()

	return nil
}

// onConnTaskEnd - on ConnTask end
func (c *jarvisClient2) onConnTaskEnd(servaddr string) {
	c.mapConnServAddr.Delete(servaddr)
}

// addConnTask - add a client task
func (c *jarvisClient2) addConnTask(servaddr string, node *coredbpb.NodeInfo, funcOnResult FuncOnProcMsgResult) {

	_, ok := c.mapConnServAddr.Load(servaddr)
	if ok {
		return
	}

	c.mapConnServAddr.Store(servaddr, 0)

	task := &clientTask{
		servaddr:     servaddr,
		client:       c,
		node:         node,
		funcOnResult: funcOnResult,
	}

	c.poolConn.SendTask(task)
}

// addSendMsgTask - add a client send message task
func (c *jarvisClient2) addSendMsgTask(msg *pb.JarvisMsg, addr string, funcOnResult FuncOnProcMsgResult) {
	task := &clientTask{
		msg:          msg,
		client:       c,
		funcOnResult: funcOnResult,
		addr:         addr,
	}

	task.Init(c.poolMsg, addr)

	c.poolMsg.SendTask(task)
}

// addSendMsgStreamTask - add a client send message stream task
func (c *jarvisClient2) addSendMsgStreamTask(msgs []*pb.JarvisMsg, addr string, funcOnResult FuncOnProcMsgResult) {
	task := &clientTask{
		msgs:         msgs,
		client:       c,
		funcOnResult: funcOnResult,
		addr:         addr,
	}

	task.Init(c.poolMsg, addr)

	c.poolMsg.SendTask(task)
}

func (c *jarvisClient2) isConnected(addr string) bool {
	_, ok := c.mapClient.Load(addr)

	return ok
}

func (c *jarvisClient2) _getClientConn(addr string) *clientInfo2 {
	mi, ok := c.mapClient.Load(addr)
	if ok {
		ci, typeok := mi.(*clientInfo2)
		if typeok {
			return ci
		}
	}

	return nil
}

func (c *jarvisClient2) _getValidClientConn(addr string) (*clientInfo2, error) {
	ci := c._getClientConn(addr)
	if ci != nil {
		if mgrConn.isValidConn(ci.servAddr) {
			return ci, nil
		}
	}

	cn := c.node.GetCoreDB().GetNode(addr)
	if cn == nil {
		jarvisbase.Warn("jarvisClient2.getValidClientConn:GetNode", zap.Error(ErrCannotFindNodeWithAddr))

		return nil, ErrCannotFindNodeWithAddr
	}

	conn, err := mgrConn.getConn(cn.ServAddr)
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

func (c *jarvisClient2) _sendMsg(ctx context.Context, smsg *pb.JarvisMsg, funcOnResult FuncOnProcMsgResult) error {
	// var lstResult []*JarvisMsgInfo

	newsendmsgid := c._getNewSendMsgID(smsg.DestAddr)
	destaddr := smsg.DestAddr

	if funcOnResult != nil {
		err := c.node.OnClientProcMsg(destaddr, newsendmsgid, funcOnResult)
		if err != nil {
			jarvisbase.Warn("jarvisClient2._sendMsg:OnClientProcMsg",
				zap.Error(err),
				JSONMsg2Zap("msg", smsg),
				zap.Int64("newsendmsgid", newsendmsgid))
		}
	}

	_, ok := c.mapClient.Load(smsg.DestAddr)
	if !ok {
		jarvisbase.Warn("jarvisClient2._sendMsg:mapClient",
			zap.Error(ErrNotConnectedNode),
			JSONMsg2Zap("msg", smsg))

		if funcOnResult != nil {
			c.node.OnReplyProcMsg(ctx, destaddr, newsendmsgid, nil, ErrNotConnectedNode)
		}

		// return ErrNotConnectedNode
		// return c._broadCastMsg(ctx, smsg)
	}

	jarvisbase.Debug("jarvisClient2._sendMsg",
		JSONMsg2Zap("msg", smsg))

	ci2, err := c._getValidClientConn(smsg.DestAddr)
	if err != nil {
		jarvisbase.Warn("jarvisClient2._sendMsg:getValidClientConn", zap.Error(err))

		if funcOnResult != nil {
			c.node.OnReplyProcMsg(ctx, destaddr, newsendmsgid, nil, err)
		}

		return err
	}

	err = c._signJarvisMsg(smsg, newsendmsgid, false)
	if err != nil {
		jarvisbase.Warn("jarvisClient2._sendMsg:_signJarvisMsg", zap.Error(err))

		if funcOnResult != nil {
			c.node.OnReplyProcMsg(ctx, destaddr, newsendmsgid, nil, err)
		}

		return err
	}

	stream, err := ci2.client.ProcMsg(ctx, smsg)
	if err != nil {
		jarvisbase.Warn("jarvisClient2._sendMsg:ProcMsg", zap.Error(err))

		if funcOnResult != nil {
			c.node.OnReplyProcMsg(ctx, destaddr, newsendmsgid, nil, err)
		}

		return err
	}

	for {
		getmsg, err := stream.Recv()
		if err == io.EOF {
			jarvisbase.Debug("jarvisClient2._sendMsg:stream eof")

			if funcOnResult != nil {
				c.node.OnReplyProcMsg(ctx, destaddr, newsendmsgid, nil, nil)
			}

			break
		}

		if err != nil {
			jarvisbase.Warn("jarvisClient2._sendMsg:stream", zap.Error(err))

			if funcOnResult != nil {
				c.node.OnReplyProcMsg(ctx, destaddr, newsendmsgid, nil, err)
			}

			break
		} else {
			jarvisbase.Debug("jarvisClient2._sendMsg:stream",
				JSONMsg2Zap("msg", getmsg))

			c.node.PostMsg(&NormalTaskInfo{
				Msg: getmsg,
			}, nil)

			if funcOnResult != nil {
				c.node.OnReplyProcMsg(ctx, destaddr, newsendmsgid, getmsg, nil)
			}
		}
	}

	return nil
}

func (c *jarvisClient2) _broadCastMsg(ctx context.Context, msg *pb.JarvisMsg) error {

	jarvisbase.Debug("jarvisClient2._broadCastMsg",
		JSONMsg2Zap("msg", msg))

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

func (c *jarvisClient2) _connectNode(ctx context.Context, servaddr string, node *coredbpb.NodeInfo, funcOnResult FuncOnProcMsgResult) error {
	var lstResult []*JarvisMsgInfo

	if node != nil {
		if node.Addr == c.node.myinfo.Addr {
			jarvisbase.Warn("jarvisClient2._connectNode:checkNodeAddr",
				zap.Error(ErrServAddrIsMe),
				zap.String("addr", c.node.myinfo.Addr),
				zap.String("bindaddr", c.node.myinfo.BindAddr),
				zap.String("servaddr", c.node.myinfo.ServAddr))

			if funcOnResult != nil {
				lstResult = append(lstResult, &JarvisMsgInfo{
					Err: ErrServAddrIsMe,
				})

				funcOnResult(ctx, c.node, lstResult)
			}

			return ErrServAddrIsMe
		}

		if coredb.IsDeprecatedNode(node) {
			jarvisbase.Warn("jarvisClient2._connectNode:IsDeprecatedNode",
				zap.Error(ErrDeprecatedNode),
				zap.String("addr", c.node.myinfo.Addr),
				zap.String("bindaddr", c.node.myinfo.BindAddr),
				zap.String("servaddr", c.node.myinfo.ServAddr))

			if funcOnResult != nil {
				lstResult = append(lstResult, &JarvisMsgInfo{
					Err: ErrDeprecatedNode,
				})

				funcOnResult(ctx, c.node, lstResult)
			}

			return ErrDeprecatedNode
		}
	}

	if !IsValidServAddr(servaddr) {
		jarvisbase.Warn("jarvisClient2._connectNode",
			zap.Error(ErrInvalidServAddr),
			zap.String("addr", c.node.myinfo.Addr),
			zap.String("servaddr", servaddr))

		if funcOnResult != nil {
			lstResult = append(lstResult, &JarvisMsgInfo{
				Err: ErrInvalidServAddr,
			})

			funcOnResult(ctx, c.node, lstResult)
		}

		return ErrInvalidServAddr
	}

	if IsMyServAddr(servaddr, c.node.myinfo.BindAddr) {
		jarvisbase.Warn("jarvisClient2._connectNode",
			zap.Error(ErrServAddrIsMe),
			zap.String("addr", c.node.myinfo.Addr),
			zap.String("bindaddr", c.node.myinfo.BindAddr),
			zap.String("servaddr", c.node.myinfo.ServAddr))

		if funcOnResult != nil {
			lstResult = append(lstResult, &JarvisMsgInfo{
				Err: ErrServAddrIsMe,
			})

			funcOnResult(ctx, c.node, lstResult)
		}

		return ErrServAddrIsMe
	}

	conn, err := mgrConn.getConn(servaddr)
	if err != nil {
		jarvisbase.Warn("jarvisClient2._connectNode", zap.Error(err))

		if funcOnResult != nil {
			lstResult = append(lstResult, &JarvisMsgInfo{
				Err: err,
			})

			funcOnResult(ctx, c.node, lstResult)
		}

		return err
	}

	if c.fsa.isFailServAddr(servaddr) {
		return ErrServAddrConnFail
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

		if funcOnResult != nil {
			lstResult = append(lstResult, &JarvisMsgInfo{
				Err: err,
			})

			funcOnResult(ctx, c.node, lstResult)
		}

		return err
	}

	err = c._signJarvisMsg(msg, 0, false)
	if err != nil {
		jarvisbase.Warn("jarvisClient2._connectNode:_signJarvisMsg", zap.Error(err))

		if funcOnResult != nil {
			lstResult = append(lstResult, &JarvisMsgInfo{
				Err: err,
			})

			funcOnResult(ctx, c.node, lstResult)
		}

		return err
	}

	stream, err1 := ci.client.ProcMsg(curctx, msg)
	if err1 != nil {
		jarvisbase.Warn("jarvisClient2._connectNode:ProcMsg", zap.Error(err1))

		if funcOnResult != nil {
			lstResult = append(lstResult, &JarvisMsgInfo{
				Err: err1,
			})

			funcOnResult(ctx, c.node, lstResult)
		}

		err := conn.Close()
		if err != nil {
			jarvisbase.Warn("jarvisClient2._connectNode:Close", zap.Error(err1))
		}

		mgrConn.delConn(servaddr)

		c.fsa.onConnFail(servaddr)

		return err1
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			jarvisbase.Debug("jarvisClient2._connectNode:stream eof")

			if funcOnResult != nil {
				lstResult = append(lstResult, &JarvisMsgInfo{})

				funcOnResult(ctx, c.node, lstResult)
			}

			break
		}

		if err != nil {
			grpcerr, ok := status.FromError(err)
			if ok && grpcerr.Code() == codes.Unimplemented {
				jarvisbase.Warn("jarvisClient2._connectNode:ProcMsg:Unimplemented", zap.Error(err))

				if funcOnResult != nil {
					lstResult = append(lstResult, &JarvisMsgInfo{
						Err: err,
					})

					funcOnResult(ctx, c.node, lstResult)
				}

				err1 := conn.Close()
				if err1 != nil {
					jarvisbase.Warn("jarvisClient2._connectNode:Close", zap.Error(err))
				}

				mgrConn.delConn(servaddr)

				c.fsa.onConnFail(servaddr)

				return err
			}

			jarvisbase.Warn("jarvisClient2._connectNode:stream",
				zap.Error(err),
				zap.String("servaddr", servaddr))

			if funcOnResult != nil {
				lstResult = append(lstResult, &JarvisMsgInfo{
					Err: err,
				})

				funcOnResult(ctx, c.node, lstResult)
			}

			return err
		}

		jarvisbase.Debug("jarvisClient2._connectNode:stream",
			JSONMsg2Zap("msg", msg))

		if msg.MsgType == pb.MSGTYPE_REPLY_CONNECT {
			ni := msg.GetNodeInfo()

			c.mapClient.Store(ni.Addr, ci)
		}

		c.node.PostMsg(&NormalTaskInfo{
			Msg: msg,
		}, nil)

		if funcOnResult != nil {
			lstResult = append(lstResult, &JarvisMsgInfo{
				Msg: msg,
			})

			funcOnResult(ctx, c.node, lstResult)
		}

	}

	return nil
}

func (c *jarvisClient2) _getNewSendMsgID(destaddr string) int64 {
	return c.node.GetCoreDB().GetNewSendMsgID(destaddr)
}

func (c *jarvisClient2) _signJarvisMsg(msg *pb.JarvisMsg, newsendmsgid int64, isstream bool) error {
	if isstream {
		msg.StreamMsgID = newsendmsgid
	} else {
		msg.MsgID = newsendmsgid
	}

	msg.CurTime = time.Now().Unix()
	msg.LastMsgID = c.node.GetCoreDB().GetCurRecvMsgID(msg.DestAddr)

	return SignJarvisMsg(c.node.GetCoreDB().GetPrivateKey(), msg)
}

func (c *jarvisClient2) _procRecvMsgStream(ctx context.Context,
	stream pb.JarvisCoreServ_ProcMsgStreamClient, funcOnResult FuncOnProcMsgResult,
	chanEnd chan int, destaddr string, newsendmsgid int64) {

	for {
		getmsg, err := stream.Recv()
		if err == io.EOF {
			jarvisbase.Debug("jarvisClient2._sendMsg:stream eof")

			if funcOnResult != nil {
				c.node.OnReplyProcMsg(ctx, destaddr, newsendmsgid, nil, nil)
			}

			break
		}

		if err != nil {
			jarvisbase.Warn("jarvisClient2._sendMsg:stream", zap.Error(err))

			if funcOnResult != nil {
				c.node.OnReplyProcMsg(ctx, destaddr, newsendmsgid, nil, err)
			}

			break
		} else {
			jarvisbase.Debug("jarvisClient2._sendMsg:stream",
				JSONMsg2Zap("msg", getmsg))

			c.node.PostMsg(&NormalTaskInfo{
				Msg: getmsg,
			}, nil)

			if funcOnResult != nil {
				c.node.OnReplyProcMsg(ctx, destaddr, newsendmsgid, getmsg, nil)
			}

		}
	}

	chanEnd <- 0
}

func (c *jarvisClient2) _sendMsgStream(ctx context.Context, destAddr string, smsgs []*pb.JarvisMsg, funcOnResult FuncOnProcMsgResult) error {

	newsendmsgid := c._getNewSendMsgID(destAddr)
	destaddr := smsgs[0].DestAddr

	if funcOnResult != nil {
		err := c.node.OnClientProcMsg(destaddr, newsendmsgid, funcOnResult)
		if err != nil {
			jarvisbase.Warn("jarvisClient2._sendMsgStream:OnClientProcMsg",
				zap.Error(err),
				zap.String("destaddr", destaddr),
				zap.Int64("newsendmsgid", newsendmsgid))
		}
	}

	// var lstResult []*JarvisMsgInfo

	_, ok := c.mapClient.Load(destAddr)
	if !ok {
		jarvisbase.Warn("jarvisClient2._sendMsgStream:mapClient",
			zap.Error(ErrNotConnectedNode))

		if funcOnResult != nil {
			c.node.OnReplyProcMsg(ctx, destaddr, newsendmsgid, nil, ErrNotConnectedNode)
		}
	}

	ci2, err := c._getValidClientConn(destAddr)
	if err != nil {
		jarvisbase.Warn("jarvisClient2._sendMsgStream:getValidClientConn", zap.Error(err))

		if funcOnResult != nil {
			c.node.OnReplyProcMsg(ctx, destaddr, newsendmsgid, nil, err)
		}

		return err
	}

	chanEnd := make(chan int)
	stream, err := ci2.client.ProcMsgStream(ctx)
	if err != nil {
		jarvisbase.Warn("jarvisClient2._sendMsgStream:ProcMsgStream", zap.Error(err))

		if funcOnResult != nil {
			c.node.OnReplyProcMsg(ctx, destaddr, newsendmsgid, nil, err)
		}

		return err
	}

	if stream == nil {
		jarvisbase.Warn("jarvisClient2._sendMsgStream:ProcMsgStream", zap.Error(ErrProcMsgStreamNil))

		if funcOnResult != nil {
			c.node.OnReplyProcMsg(ctx, destaddr, newsendmsgid, nil, ErrProcMsgStreamNil)
		}

		return err
	}

	go c._procRecvMsgStream(ctx, stream, funcOnResult, chanEnd, destaddr, newsendmsgid)

	for i := 0; i < len(smsgs); i++ {
		err := c._signJarvisMsg(smsgs[i], newsendmsgid, true)
		if err != nil {
			jarvisbase.Warn("jarvisClient2._sendMsgStream:_signJarvisMsg", zap.Error(err))

			if funcOnResult != nil {
				c.node.OnReplyProcMsg(ctx, destaddr, newsendmsgid, nil, err)
			}

			return err
		}

		err = stream.Send(smsgs[i])
		if err != nil {
			jarvisbase.Warn("jarvisClient2._sendMsgStream:ProcMsg", zap.Error(err))

			if funcOnResult != nil {
				c.node.OnReplyProcMsg(ctx, destaddr, newsendmsgid, nil, err)
			}

			return err
		}
	}

	err = stream.CloseSend()
	if err != nil {
		jarvisbase.Warn("jarvisClient2._sendMsgStream:CloseSend", zap.Error(err))

		if funcOnResult != nil {
			c.node.OnReplyProcMsg(ctx, destaddr, newsendmsgid, nil, err)
		}

		return err
	}

	<-chanEnd

	return nil
}
