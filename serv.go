package jarviscore

import (
	"context"
	"net"
	"strings"
	"sync"

	"go.uber.org/zap"

	"github.com/zhs007/jarviscore/log"
	pb "github.com/zhs007/jarviscore/proto"
	"google.golang.org/grpc"
)

type serverClientChan (chan BaseInfo)

// jarvisServer
type jarvisServer struct {
	sync.RWMutex

	node            *jarvisNode
	lis             net.Listener
	grpcServ        *grpc.Server
	servchan        chan int
	mapChanNodeInfo map[string]serverClientChan
}

// newServer -
func newServer(node *jarvisNode) (*jarvisServer, error) {
	lis, err := net.Listen("tcp", node.myinfo.BindAddr)
	if err != nil {
		errorLog("newServer", err)

		return nil, err
	}

	log.Info("Listen", zap.String("addr", node.myinfo.BindAddr))

	grpcServ := grpc.NewServer()
	s := &jarvisServer{
		node:            node,
		lis:             lis,
		grpcServ:        grpcServ,
		servchan:        make(chan int, 1),
		mapChanNodeInfo: make(map[string]serverClientChan),
	}
	pb.RegisterJarvisCoreServServer(grpcServ, s)

	return s, nil
}

// Start -
func (s *jarvisServer) Start() (err error) {
	err = s.grpcServ.Serve(s.lis)

	s.servchan <- 0

	return
}

// Stop -
func (s *jarvisServer) Stop() {
	s.Lock()
	for _, v := range s.mapChanNodeInfo {
		close(v)
	}
	s.Unlock()

	s.lis.Close()

	return
}

// hasClient
func (s *jarvisServer) hasClient(token string) bool {
	s.RLock()
	defer s.RUnlock()

	if _, ok := s.mapChanNodeInfo[token]; ok {
		return true
	}

	return false
}

// Join implements jarviscorepb.JarvisCoreServ
func (s *jarvisServer) Join(ctx context.Context, in *pb.Join) (*pb.ReplyJoin, error) {
	log.Info("JarvisServ.Join",
		zap.String("Servaddr", in.ServAddr),
		zap.String("Token", in.Token),
		zap.String("Name", in.Name),
		zap.Int("Nodetype", int(in.NodeType)))

	peeripaddr := in.ServAddr
	addrSlice := strings.Split(in.ServAddr, ":")
	if len(addrSlice) == 2 {
		if (addrSlice[0] == "" || addrSlice[0] == "0.0.0.0") && addrSlice[1] != "" {
			clientip := getGRPCClientIP(ctx)
			if clientip != "" {
				peeripaddr = clientip + ":" + addrSlice[1]
			}
		}
	}

	isvalidnode := (in.Token != s.node.myinfo.Token)

	// bi := BaseInfo{
	// 	Name:     in.Name,
	// 	ServAddr: in.ServAddr,
	// 	Token:    in.Token,
	// 	NodeType: in.NodeType,
	// }

	s.node.onNodeConnectMe(&BaseInfo{
		Name:     in.Name,
		ServAddr: peeripaddr,
		Token:    in.Token,
		NodeType: in.NodeType,
	})

	if isvalidnode {
		return &pb.ReplyJoin{
			Code:     pb.CODE_OK,
			Name:     s.node.myinfo.Name,
			Token:    s.node.myinfo.Token,
			NodeType: s.node.myinfo.NodeType,
		}, nil
	}

	return &pb.ReplyJoin{
		Code:     pb.CODE_ALREADY_IN,
		Name:     s.node.myinfo.Name,
		Token:    s.node.myinfo.Token,
		NodeType: s.node.myinfo.NodeType,
	}, nil
}

func (s *jarvisServer) broadcastNode(bi *BaseInfo) {
	s.RLock()
	defer s.RUnlock()

	for _, v := range s.mapChanNodeInfo {
		v <- *bi
	}
}

// Subscribe implements jarviscorepb.JarvisCoreServ
func (s *jarvisServer) Subscribe(in *pb.Subscribe, stream pb.JarvisCoreServ_SubscribeServer) error {
	if !s.node.hasNodeToken(in.Token) {
		return nil
	}

	s.node.mgrNodeInfo.foreach(func(cn *NodeInfo) {
		ni := pb.NodeInfo{
			ServAddr: cn.baseinfo.ServAddr,
			Token:    cn.baseinfo.Token,
			Name:     cn.baseinfo.Name,
			NodeType: cn.baseinfo.NodeType,
		}

		stream.SendMsg(&pb.ChannelInfo{
			ChannelType: in.ChannelType,
			Data:        &pb.ChannelInfo_NodeInfo{NodeInfo: &ni},
		})
	})

	s.Lock()
	if _, ok := s.mapChanNodeInfo[in.Token]; ok {
		close(s.mapChanNodeInfo[in.Token])
	}

	chanNodeInfo := make(chan BaseInfo)
	s.mapChanNodeInfo[in.Token] = chanNodeInfo
	s.Unlock()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case bi, ok := <-chanNodeInfo:
			if !ok {
				return nil
			}

			ni := pb.NodeInfo{
				ServAddr: bi.ServAddr,
				Token:    bi.Token,
				Name:     bi.Name,
				NodeType: bi.NodeType,
			}

			stream.SendMsg(&pb.ChannelInfo{
				ChannelType: in.ChannelType,
				Data:        &pb.ChannelInfo_NodeInfo{NodeInfo: &ni},
			})
		}

		// gs.Send(&pb.HelloReply{Message: "Hello " + in.Name})
	}
	// return nil
}

// RequestCtrl implements jarviscorepb.JarvisCoreServ
func (s *jarvisServer) RequestCtrl(ctx context.Context, in *pb.CtrlInfo) (*pb.BaseReply, error) {
	if in.DestToken == s.node.myinfo.Token {
		_, err := mgrCtrl.Run(in.CtrlType, in.Command)
		if err != nil {
			return &pb.BaseReply{
				Code: pb.CODE_OK,
			}, nil
		}

		return &pb.BaseReply{
			Code: pb.CODE_OK,
		}, nil
	}

	if s.hasClient(in.DestToken) {
		return &pb.BaseReply{
			Code: pb.CODE_FORWARD_MSG,
		}, nil
	}

	return &pb.BaseReply{
		Code: pb.CODE_BROADCAST_MSG,
	}, nil
}

// ReplyCtrl implements jarviscorepb.JarvisCoreServ
func (s *jarvisServer) ReplyCtrl(ctx context.Context, in *pb.CtrlResult) (*pb.BaseReply, error) {
	return &pb.BaseReply{
		Code: pb.CODE_OK,
	}, nil
}
