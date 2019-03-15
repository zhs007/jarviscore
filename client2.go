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

// ClientProcMsgResult - result for client.ProcMsg
type ClientProcMsgResult struct {
	Msg *pb.JarvisMsg `json:"msg"`
	Err error         `json:"err"`
}

// FuncOnProcMsgResult - on procmsg recv the message
type FuncOnProcMsgResult func(ctx context.Context, jarvisnode JarvisNode,
	lstResult []*ClientProcMsgResult) error

// IsClientProcMsgResultEnd - is end
func IsClientProcMsgResultEnd(lstResult []*ClientProcMsgResult) bool {
	return len(lstResult) > 0 && lstResult[len(lstResult)-1].Msg == nil
}

// ClientGroupProcMsgResults - result for FuncOnSendMsgResult
type ClientGroupProcMsgResults struct {
	Results []*ClientProcMsgResult `json:"results"`
}

// FuncOnGroupSendMsgResult - on group sendmsg recv the messages
type FuncOnGroupSendMsgResult func(ctx context.Context, jarvisnode JarvisNode,
	numsNode int, lstResult []*ClientGroupProcMsgResults) error

// CountClientGroupProcMsgResultsEnd - count number for end
func CountClientGroupProcMsgResultsEnd(lstResult []*ClientGroupProcMsgResults) int {
	nums := 0
	for i := 0; i < len(lstResult); i++ {
		cr := lstResult[i]
		if len(cr.Results) > 0 && cr.Results[len(cr.Results)-1].Msg == nil {
			nums++
		}
	}

	return nums
}

type clientTask struct {
	jarvisbase.L2BaseTask

	servaddr     string
	addr         string
	client       *jarvisClient2
	msg          *pb.JarvisMsg
	node         *coredbpb.NodeInfo
	funcOnResult FuncOnProcMsgResult
}

func (task *clientTask) Run(ctx context.Context) error {
	if task.msg == nil {
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

	return task.client._sendMsg(ctx, task.msg, task.funcOnResult)
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
func (c *jarvisClient2) BuildMsgPoolStatus() *pb.L2PoolInfo {
	return c.poolMsg.BuildStatus()
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
func (c *jarvisClient2) addSendMsgTask(msg *pb.JarvisMsg, node *coredbpb.NodeInfo, funcOnResult FuncOnProcMsgResult) {
	task := &clientTask{
		msg:          msg,
		client:       c,
		node:         node,
		funcOnResult: funcOnResult,
		addr:         msg.DestAddr,
	}

	task.Init(c.poolMsg, msg.DestAddr)

	c.poolMsg.SendTask(task)
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

func (c *jarvisClient2) _sendMsg(ctx context.Context, smsg *pb.JarvisMsg, funcOnResult FuncOnProcMsgResult) error {
	var lstResult []*ClientProcMsgResult

	_, ok := c.mapClient.Load(smsg.DestAddr)
	if !ok {
		jarvisbase.Warn("jarvisClient2._sendMsg:getValidClientConn", zap.Error(ErrNotConnectedNode))

		if funcOnResult != nil {
			lstResult = append(lstResult, &ClientProcMsgResult{
				Err: ErrNotConnectedNode,
			})

			funcOnResult(ctx, c.node, lstResult)
		}

		return ErrNotConnectedNode
		// return c._broadCastMsg(ctx, smsg)
	}

	jarvisbase.Debug("jarvisClient2._sendMsg", jarvisbase.JSON("msg", smsg))

	ci2, err := c._getValidClientConn(smsg.DestAddr)
	if err != nil {
		jarvisbase.Warn("jarvisClient2._sendMsg:getValidClientConn", zap.Error(err))

		if funcOnResult != nil {
			lstResult = append(lstResult, &ClientProcMsgResult{
				Err: err,
			})

			funcOnResult(ctx, c.node, lstResult)
		}

		return err
	}

	err = c._signJarvisMsg(smsg)
	if err != nil {
		jarvisbase.Warn("jarvisClient2._sendMsg:_signJarvisMsg", zap.Error(err))

		if funcOnResult != nil {
			lstResult = append(lstResult, &ClientProcMsgResult{
				Err: err,
			})

			funcOnResult(ctx, c.node, lstResult)
		}

		return err
	}

	stream, err := ci2.client.ProcMsg(ctx, smsg)
	if err != nil {
		jarvisbase.Warn("jarvisClient2._sendMsg:ProcMsg", zap.Error(err))

		if funcOnResult != nil {
			lstResult = append(lstResult, &ClientProcMsgResult{
				Err: err,
			})

			funcOnResult(ctx, c.node, lstResult)
		}

		return err
	}

	for {
		getmsg, err := stream.Recv()
		if err == io.EOF {
			jarvisbase.Debug("jarvisClient2._sendMsg:stream eof")

			if funcOnResult != nil {
				lstResult = append(lstResult, &ClientProcMsgResult{})

				funcOnResult(ctx, c.node, lstResult)
			}

			break
		}

		if err != nil {
			jarvisbase.Warn("jarvisClient2._sendMsg:stream", zap.Error(err))

			if funcOnResult != nil {
				lstResult = append(lstResult, &ClientProcMsgResult{
					Err: err,
				})

				funcOnResult(ctx, c.node, lstResult)
			}

			break
		} else {
			jarvisbase.Debug("jarvisClient2._sendMsg:stream", jarvisbase.JSON("msg", getmsg))

			c.node.PostMsg(getmsg, nil, nil, nil)

			if funcOnResult != nil {
				lstResult = append(lstResult, &ClientProcMsgResult{
					Msg: getmsg,
				})

				funcOnResult(ctx, c.node, lstResult)
			}
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

func (c *jarvisClient2) _connectNode(ctx context.Context, servaddr string, node *coredbpb.NodeInfo, funcOnResult FuncOnProcMsgResult) error {
	var lstResult []*ClientProcMsgResult

	if node != nil {
		if node.Addr == c.node.myinfo.Addr {
			jarvisbase.Warn("jarvisClient2._connectNode:checkNodeAddr",
				zap.Error(ErrServAddrIsMe),
				zap.String("addr", c.node.myinfo.Addr),
				zap.String("bindaddr", c.node.myinfo.BindAddr),
				zap.String("servaddr", c.node.myinfo.ServAddr))

			if funcOnResult != nil {
				lstResult = append(lstResult, &ClientProcMsgResult{
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
				lstResult = append(lstResult, &ClientProcMsgResult{
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
			lstResult = append(lstResult, &ClientProcMsgResult{
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
			lstResult = append(lstResult, &ClientProcMsgResult{
				Err: ErrServAddrIsMe,
			})

			funcOnResult(ctx, c.node, lstResult)
		}

		return ErrServAddrIsMe
	}

	conn, err := mgrconn.getConn(servaddr)
	if err != nil {
		jarvisbase.Warn("jarvisClient2._connectNode", zap.Error(err))

		if funcOnResult != nil {
			lstResult = append(lstResult, &ClientProcMsgResult{
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
			lstResult = append(lstResult, &ClientProcMsgResult{
				Err: err,
			})

			funcOnResult(ctx, c.node, lstResult)
		}

		return err
	}

	err = c._signJarvisMsg(msg)
	if err != nil {
		jarvisbase.Warn("jarvisClient2._connectNode:_signJarvisMsg", zap.Error(err))

		if funcOnResult != nil {
			lstResult = append(lstResult, &ClientProcMsgResult{
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
			lstResult = append(lstResult, &ClientProcMsgResult{
				Err: err1,
			})

			funcOnResult(ctx, c.node, lstResult)
		}

		err := conn.Close()
		if err != nil {
			jarvisbase.Warn("jarvisClient2._connectNode:Close", zap.Error(err1))
		}

		mgrconn.delConn(servaddr)

		c.fsa.onConnFail(servaddr)

		return err1
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			jarvisbase.Debug("jarvisClient2._connectNode:stream eof")

			if funcOnResult != nil {
				lstResult = append(lstResult, &ClientProcMsgResult{})

				funcOnResult(ctx, c.node, lstResult)
			}

			break
		}

		if err != nil {
			grpcerr, ok := status.FromError(err)
			if ok && grpcerr.Code() == codes.Unimplemented {
				jarvisbase.Warn("jarvisClient2._connectNode:ProcMsg:Unimplemented", zap.Error(err))

				if funcOnResult != nil {
					lstResult = append(lstResult, &ClientProcMsgResult{
						Err: err,
					})

					funcOnResult(ctx, c.node, lstResult)
				}

				err1 := conn.Close()
				if err1 != nil {
					jarvisbase.Warn("jarvisClient2._connectNode:Close", zap.Error(err))
				}

				mgrconn.delConn(servaddr)

				c.fsa.onConnFail(servaddr)

				return err
			}

			jarvisbase.Warn("jarvisClient2._connectNode:stream",
				zap.Error(err),
				zap.String("servaddr", servaddr))

			if funcOnResult != nil {
				lstResult = append(lstResult, &ClientProcMsgResult{
					Err: err,
				})

				funcOnResult(ctx, c.node, lstResult)
			}

			return err
		}

		jarvisbase.Debug("jarvisClient2._connectNode:stream", jarvisbase.JSON("msg", msg))

		if msg.MsgType == pb.MSGTYPE_REPLY_CONNECT {
			ni := msg.GetNodeInfo()

			c.mapClient.Store(ni.Addr, ci)
		}

		c.node.PostMsg(msg, nil, nil, nil)

		if funcOnResult != nil {
			lstResult = append(lstResult, &ClientProcMsgResult{
				Msg: msg,
			})

			funcOnResult(ctx, c.node, lstResult)
		}

	}

	return nil
}

func (c *jarvisClient2) _signJarvisMsg(msg *pb.JarvisMsg) error {
	msg.MsgID = c.node.GetCoreDB().GetNewSendMsgID(msg.DestAddr)
	msg.CurTime = time.Now().Unix()

	return SignJarvisMsg(c.node.GetCoreDB().GetPrivateKey(), msg)
}
