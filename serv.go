package jarviscore

import (
	"context"
	"net"

	"go.uber.org/zap"

	"github.com/zhs007/jarviscore/log"
	pb "github.com/zhs007/jarviscore/proto"
	"google.golang.org/grpc"
)

// jarvisServer
type jarvisServer struct {
	servaddr string
	lis      net.Listener
	grpcServ *grpc.Server
	servchan chan int
}

// NewServer -
func newServer(servaddr string) (*jarvisServer, error) {
	lis, err := net.Listen("tcp", servaddr)
	if err != nil {
		errorLog("newServer", err)

		return nil, err
	}

	log.Info("Listen", zap.String("addr", servaddr))

	grpcServ := grpc.NewServer()
	s := &jarvisServer{servaddr: servaddr, lis: lis, grpcServ: grpcServ, servchan: make(chan int, 1)}
	pb.RegisterJarvisCoreServServer(grpcServ, s)

	return s, nil
}

func (s *jarvisServer) Start() (err error) {
	err = s.grpcServ.Serve(s.lis)

	s.servchan <- 0

	return
}

func (s *jarvisServer) Stop() {
	s.lis.Close()

	return
}

// Join implements jarviscorepb.JarvisCoreServ
func (s *jarvisServer) Join(ctx context.Context, in *pb.Join) (*pb.ReplyJoin, error) {
	return &pb.ReplyJoin{Code: pb.CODE_OK}, nil
}
