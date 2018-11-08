package jarviscore

import (
	"context"
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

// Join implements jarviscorepb.JarvisCoreServ
func (s *jarvisServer2) Join(ctx context.Context, in *pb.Join) (*pb.ReplyJoin, error) {
	return &pb.ReplyJoin{}, nil
}

// Subscribe implements jarviscorepb.JarvisCoreServ
func (s *jarvisServer2) Subscribe(in *pb.Subscribe, stream pb.JarvisCoreServ_SubscribeServer) error {
	return nil
}

// RequestCtrl implements jarviscorepb.JarvisCoreServ
func (s *jarvisServer2) RequestCtrl(ctx context.Context, in *pb.CtrlInfo) (*pb.BaseReply, error) {
	return &pb.BaseReply{}, nil
}

// ReplyCtrl implements jarviscorepb.JarvisCoreServ
func (s *jarvisServer2) ReplyCtrl(ctx context.Context, in *pb.CtrlResult) (*pb.BaseReply, error) {
	return &pb.BaseReply{}, nil
}

// GetMyServAddr implements jarviscorepb.JarvisCoreServ
func (s *jarvisServer2) GetMyServAddr(ctx context.Context, in *pb.ServAddr) (*pb.ServAddr, error) {
	return &pb.ServAddr{}, nil
}

// Trust implements jarviscorepb.JarvisCoreServ
func (s *jarvisServer2) Trust(ctx context.Context, in *pb.TrustNode) (*pb.BaseReply, error) {
	return &pb.BaseReply{}, nil
}

// Trust implements jarviscorepb.JarvisCoreServ
func (s *jarvisServer2) Join2(ctx context.Context, in *pb.JarvisMsg) (*pb.JarvisMsg, error) {
	return &pb.JarvisMsg{}, nil
}

// SendMsg implements jarviscorepb.JarvisCoreServ
func (s *jarvisServer2) SendMsg(in *pb.JarvisMsg, stream pb.JarvisCoreServ_SendMsgServer) error {
	s.node.mgrJasvisMsg.sendMsg(in)

	return nil
}
