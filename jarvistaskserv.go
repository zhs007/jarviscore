package jarviscore

import (
	"context"
	"net"

	"go.uber.org/zap"

	"github.com/zhs007/jarviscore/base"
	pb "github.com/zhs007/jarviscore/proto"
	"google.golang.org/grpc"
)

// jarvisTaskServ
type jarvisTaskServ struct {
	node     *jarvisNode
	lis      net.Listener
	grpcServ *grpc.Server
}

// newJarvisTask -
func newJarvisTask(node *jarvisNode) (*jarvisTaskServ, error) {
	lis, err := net.Listen("tcp", node.GetConfig().TaskServ.BindAddr)
	if err != nil {
		jarvisbase.Error("newJarvisTask", zap.Error(err))

		return nil, err
	}

	jarvisbase.Info("newJarvisTask:Listen",
		zap.String("addr", node.GetConfig().TaskServ.BindAddr))

	grpcServ := grpc.NewServer()

	s := &jarvisTaskServ{
		node:     node,
		lis:      lis,
		grpcServ: grpcServ,
	}

	pb.RegisterJarvisTaskServServer(grpcServ, s)

	return s, nil
}

// Start -
func (s *jarvisTaskServ) Start(ctx context.Context) error {
	return s.grpcServ.Serve(s.lis)
}

// Stop -
func (s *jarvisTaskServ) Stop() {
	s.lis.Close()

	return
}

func (s *jarvisTaskServ) UpdTask(context.Context, *pb.JarvisTask) (*pb.ReplyUpdTask, error) {
	return &pb.ReplyUpdTask{
		IsOK: true,
	}, nil
}
