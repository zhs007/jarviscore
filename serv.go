package jarviscore

import (
	"context"
	"net"

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

func NewServer(servaddr string) (*jarvisServer, error) {
	lis, err := net.Listen("tcp", servaddr)
	if err != nil {
		return nil, err
	}

	grpcServ := grpc.NewServer()
	s := &jarvisServer{servaddr: servaddr, lis: lis, grpcServ: grpcServ, servchan: make(chan int)}
	pb.RegisterJarvisCoreServServer(grpcServ, s)

	return s, nil
}

func (s *jarvisServer) Start() (err error) {
	err = s.grpcServ.Serve(s.lis)
	s.servchan <- 0
	return
}

// Join implements jarviscorepb.JarvisCoreServ
func (s *jarvisServer) Join(ctx context.Context, in *pb.Join) (*pb.ReplyJoin, error) {
	return &pb.ReplyJoin{Code: pb.CODE_OK}, nil
}
