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

// jarvisServer2
type jarvisServer2 struct {
	node     *jarvisNode
	lis      net.Listener
	grpcServ *grpc.Server
}

// newServer2 -
func newServer2(node *jarvisNode) (*jarvisServer2, error) {
	lis, err := net.Listen("tcp", node.myinfo.BindAddr)
	if err != nil {
		jarvisbase.Error("newServer", zap.Error(err))

		return nil, err
	}

	jarvisbase.Info("Listen", zap.String("addr", node.myinfo.BindAddr))

	grpcServ := grpc.NewServer()

	s := &jarvisServer2{
		node:     node,
		lis:      lis,
		grpcServ: grpcServ,
	}

	pb.RegisterJarvisCoreServServer(grpcServ, s)

	return s, nil
}

// Start -
func (s *jarvisServer2) Start(ctx context.Context) error {
	return s.grpcServ.Serve(s.lis)
}

// Stop -
func (s *jarvisServer2) Stop() {
	s.lis.Close()

	return
}

// ProcMsg implements jarviscorepb.JarvisCoreServ
func (s *jarvisServer2) ProcMsg(in *pb.JarvisMsg, stream pb.JarvisCoreServ_ProcMsgServer) error {
	if IsSyncMsg(in) {
		chanEnd := make(chan int)

		s.node.PostMsg(&NormalTaskInfo{
			Msg:         in,
			ReplyStream: NewJarvisMsgReplyStream(stream),
		}, chanEnd)

		<-chanEnd

		return nil
	}

	// if isme
	if in.SrcAddr == s.node.myinfo.Addr {
		jarvisbase.Warn("jarvisServer2.ProcMsg:isme",
			JSONMsg2Zap("msg", in))

		err := s.replyStream2ProcMsg(in.SrcAddr, in.MsgID,
			stream, pb.REPLYTYPE_ISME, "")
		if err != nil {
			jarvisbase.Warn("jarvisServer2.ProcMsg:replyStream2ProcMsg", zap.Error(err))

			return err
		}

		return nil
	}

	if in.DestAddr == s.node.myinfo.Addr {
		err := s.replyStream2ProcMsg(in.SrcAddr, in.MsgID,
			stream, pb.REPLYTYPE_IGOTIT, "")
		if err != nil {
			jarvisbase.Warn("jarvisServer2.ProcMsg:replyStream2ProcMsg:IGOTIT", zap.Error(err))

			return err
		}
	}

	// chanEnd := make(chan int)

	s.node.PostMsg(&NormalTaskInfo{
		Msg:         in,
		ReplyStream: NewJarvisMsgReplyStream(nil),
	}, nil)

	// <-chanEnd

	return nil
}

// ProcMsgStream implements jarviscorepb.JarvisCoreServ
func (s *jarvisServer2) ProcMsgStream(stream pb.JarvisCoreServ_ProcMsgStreamServer) error {

	var lstmsgs []JarvisMsgInfo
	var firstmsg *pb.JarvisMsg

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if firstmsg == nil && in != nil {
			firstmsg = in
		}

		if err != nil {
			jarvisbase.Warn("jarvisServer2.ProcMsgStream:stream.Recv",
				zap.Error(err))

			if firstmsg != nil {
				err := s.replyStream2ProcMsgStream(firstmsg.SrcAddr, firstmsg.MsgID,
					stream, pb.REPLYTYPE_ERROR, err.Error())

				if err != nil {
					jarvisbase.Warn("jarvisServer2.ProcMsg:replyStream2ProcMsgStream",
						zap.Error(err))
				}
			}

			lstmsgs = append(lstmsgs, JarvisMsgInfo{
				JarvisResultType: JarvisResultTypeLocalError,
				Err:              err,
			})

			break
		}

		lstmsgs = append(lstmsgs, JarvisMsgInfo{
			JarvisResultType: JarvisResultTypeReply,
			Msg:              in,
		})

		if in.ReplyMsgID > 0 {
			s.node.OnReplyProcMsg(stream.Context(), in.SrcAddr, in.ReplyMsgID, JarvisResultTypeReply, in, nil)
		}
	}

	if firstmsg == nil {
		jarvisbase.Warn("jarvisServer2.ProcMsg:firstmsg")

		return nil
	}

	// if isme
	if firstmsg.SrcAddr == s.node.myinfo.Addr {
		jarvisbase.Warn("jarvisServer2.ProcMsg:isme",
			JSONMsg2Zap("msg", firstmsg))

		err := s.replyStream2ProcMsgStream(firstmsg.SrcAddr, firstmsg.MsgID,
			stream, pb.REPLYTYPE_ISME, "")
		if err != nil {
			jarvisbase.Warn("jarvisServer2.ProcMsgStream:replyStream2ProcMsgStream", zap.Error(err))

			return err
		}

		return nil
	}

	if firstmsg.DestAddr == s.node.myinfo.Addr {
		err := s.replyStream2ProcMsg(firstmsg.SrcAddr, firstmsg.MsgID,
			stream, pb.REPLYTYPE_IGOTIT, "")
		if err != nil {
			jarvisbase.Warn("jarvisServer2.ProcMsg:replyStream2ProcMsgStream:IGOTIT", zap.Error(err))

			return err
		}
	}

	// chanEnd := make(chan int)

	s.node.PostStreamMsg(&StreamTaskInfo{
		Msgs:        lstmsgs,
		ReplyStream: NewJarvisMsgReplyStream(nil),
	}, nil)

	// <-chanEnd

	return nil
}

// replyStream2ProcMsgStream
func (s *jarvisServer2) replyStream2ProcMsgStream(addr string, replyMsgID int64,
	stream pb.JarvisCoreServ_ProcMsgStreamServer, rt pb.REPLYTYPE, strErr string) error {

	sendmsg, err := BuildReply2(s.node, s.node.myinfo.Addr, addr, rt, strErr, replyMsgID)
	if err != nil {
		jarvisbase.Warn("jarvisServer2.replyStream2ProcMsgStream:BuildReply2", zap.Error(err))

		return err
	}

	err = stream.Send(sendmsg)
	if err != nil {
		jarvisbase.Warn("jarvisServer2.replyStream2ProcMsgStream:SendMsg", zap.Error(err))

		return err
	}

	return nil
}

// replyStream2ProcMsg
func (s *jarvisServer2) replyStream2ProcMsg(addr string, replyMsgID int64,
	stream pb.JarvisCoreServ_ProcMsgServer, rt pb.REPLYTYPE, strErr string) error {

	sendmsg, err := BuildReply2(s.node, s.node.myinfo.Addr, addr, rt, strErr, replyMsgID)
	if err != nil {
		jarvisbase.Warn("jarvisServer2.replyStream2ProcMsg:BuildReply2", zap.Error(err))

		return err
	}

	err = stream.Send(sendmsg)
	if err != nil {
		jarvisbase.Warn("jarvisServer2.replyStream2ProcMsg:SendMsg", zap.Error(err))

		return err
	}

	return nil
}
