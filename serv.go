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
	node     *jarvisNode
	lis      net.Listener
	grpcServ *grpc.Server
	servchan chan int
}

// NewServer -
func newServer(node *jarvisNode) (*jarvisServer, error) {
	lis, err := net.Listen("tcp", node.myinfo.ServAddr)
	if err != nil {
		errorLog("newServer", err)

		return nil, err
	}

	log.Info("Listen", zap.String("addr", node.myinfo.ServAddr))

	grpcServ := grpc.NewServer()
	s := &jarvisServer{node: node, lis: lis, grpcServ: grpcServ, servchan: make(chan int, 1)}
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
	log.Info("JarvisServ.Join",
		zap.String("Servaddr", in.Servaddr),
		zap.String("Token", in.Token),
		zap.String("Name", in.Name),
		zap.Int("Nodetype", int(in.Nodetype)))

	s.node.onAddNode(in.Servaddr)

	return &pb.ReplyJoin{Code: pb.CODE_OK}, nil
}
