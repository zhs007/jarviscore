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
	// if isme
	if in.SrcAddr == s.node.myinfo.Addr {
		jarvisbase.Warn("jarvisServer2.ProcMsg:isme",
			JSONMsg2Zap("msg", in))

		err := s.node.replyStream2(in.SrcAddr, in.MsgID,
			NewJarvisMsgReplyStream(stream, nil), pb.REPLYTYPE_ISME, "")
		if err != nil {
			jarvisbase.Warn("jarvisServer2.ProcMsg:isme:err", zap.Error(err))

			return err
		}

		return nil
	}

	chanEnd := make(chan int)

	s.node.PostMsg(&NormalTaskInfo{
		Msg: in,
		ReplyStream: &JarvisMsgReplyStream{
			procMsg: stream,
		},
	}, chanEnd)

	<-chanEnd

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
				err := s.node.replyStream2(firstmsg.SrcAddr, firstmsg.MsgID,
					NewJarvisMsgReplyStream(nil, stream), pb.REPLYTYPE_ERROR, err.Error())

				if err != nil {
					jarvisbase.Warn("jarvisServer2.ProcMsg:replyStream2:err",
						zap.Error(err))
				}
			}

			lstmsgs = append(lstmsgs, JarvisMsgInfo{
				Err: err,
			})

			break
		}

		lstmsgs = append(lstmsgs, JarvisMsgInfo{
			Msg: in,
		})

		if in.ReplyMsgID > 0 {
			s.node.OnReplyProcMsg(stream.Context(), in.SrcAddr, in.ReplyMsgID, in, nil)
		}
	}

	chanEnd := make(chan int)

	s.node.PostStreamMsg(&StreamTaskInfo{
		Msgs: lstmsgs,
		ReplyStream: &JarvisMsgReplyStream{
			procMsgStream: stream,
		},
	}, chanEnd)

	<-chanEnd

	return nil
}
