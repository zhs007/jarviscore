package jarviscore

import (
	"context"
	"net"
	"strings"
	"sync"

	"go.uber.org/zap"

	"github.com/zhs007/jarviscore/base"
	pb "github.com/zhs007/jarviscore/proto"
	"google.golang.org/grpc"
)

type serverClientChan (chan pb.ChannelInfo)

// jarvisServer
type jarvisServer struct {
	sync.RWMutex

	node        *jarvisNode
	lis         net.Listener
	grpcServ    *grpc.Server
	servchan    chan int
	mapChanInfo map[string]serverClientChan
}

// newServer -
func newServer(node *jarvisNode) (*jarvisServer, error) {
	lis, err := net.Listen("tcp", node.myinfo.BindAddr)
	if err != nil {
		jarvisbase.Error("newServer", zap.Error(err))

		return nil, err
	}

	jarvisbase.Info("Listen", zap.String("addr", node.myinfo.BindAddr))

	grpcServ := grpc.NewServer()
	s := &jarvisServer{
		node:        node,
		lis:         lis,
		grpcServ:    grpcServ,
		servchan:    make(chan int, 1),
		mapChanInfo: make(map[string]serverClientChan),
	}
	pb.RegisterJarvisCoreServServer(grpcServ, s)

	return s, nil
}

// Start -
func (s *jarvisServer) Start(ctx context.Context) (err error) {
	err = s.grpcServ.Serve(s.lis)

	// s.servchan <- 0

	return
}

// Stop -
func (s *jarvisServer) Stop() {
	s.Lock()
	for _, v := range s.mapChanInfo {
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

	if _, ok := s.mapChanInfo[token]; ok {
		return true
	}

	return false
}

// Join implements jarviscorepb.JarvisCoreServ
func (s *jarvisServer) Join(ctx context.Context, in *pb.Join) (*pb.ReplyJoin, error) {
	jarvisbase.Info("JarvisServ.Join",
		zap.String("Servaddr", in.ServAddr),
		zap.String("Addr", in.Addr),
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

	isvalidnode := (in.Addr != s.node.myinfo.Addr)

	// bi := BaseInfo{
	// 	Name:     in.Name,
	// 	ServAddr: in.ServAddr,
	// 	Token:    in.Token,
	// 	NodeType: in.NodeType,
	// }

	s.node.onNodeConnectMe(&BaseInfo{
		Name:     in.Name,
		ServAddr: peeripaddr,
		Addr:     in.Addr,
		NodeType: in.NodeType,
	})

	if isvalidnode {
		return &pb.ReplyJoin{
			Name:     s.node.myinfo.Name,
			Addr:     s.node.myinfo.Addr,
			NodeType: s.node.myinfo.NodeType,
		}, nil
	}

	return &pb.ReplyJoin{
		Err:      ErrAlreadyJoin.Error(),
		Name:     s.node.myinfo.Name,
		Addr:     s.node.myinfo.Addr,
		NodeType: s.node.myinfo.NodeType,
	}, nil
}

func (s *jarvisServer) broadcastNode(bi *BaseInfo) {
	s.RLock()
	defer s.RUnlock()

	ni := pb.NodeInfo{
		ServAddr: bi.ServAddr,
		Addr:     bi.Addr,
		Name:     bi.Name,
		NodeType: bi.NodeType,
	}

	ci := pb.ChannelInfo{
		ChannelType: pb.CHANNELTYPE_NODEINFO,
		Data: &pb.ChannelInfo_NodeInfo{
			NodeInfo: &ni,
		},
	}

	for _, v := range s.mapChanInfo {
		v <- ci
	}
}

func (s *jarvisServer) sendCtrl(ctrlinfo *pb.CtrlInfo) error {
	s.RLock()
	defer s.RUnlock()

	curChan, ok := s.mapChanInfo[ctrlinfo.DestAddr]
	if !ok {
		return ErrNotConnectMe
	}

	ci := pb.ChannelInfo{
		ChannelType: pb.CHANNELTYPE_CTRL,
		Data: &pb.ChannelInfo_CtrlInfo{
			CtrlInfo: ctrlinfo,
		},
	}

	curChan <- ci

	return nil
}

// Subscribe implements jarviscorepb.JarvisCoreServ
func (s *jarvisServer) Subscribe(in *pb.Subscribe, stream pb.JarvisCoreServ_SubscribeServer) error {
	if !s.node.hasNodeWithAddr(in.Addr) {
		return nil
	}

	s.node.mgrNodeInfo.foreach(func(cn *NodeInfo) {
		ni := pb.NodeInfo{
			ServAddr: cn.baseinfo.ServAddr,
			Addr:     cn.baseinfo.Addr,
			Name:     cn.baseinfo.Name,
			NodeType: cn.baseinfo.NodeType,
		}

		stream.SendMsg(&pb.ChannelInfo{
			ChannelType: in.ChannelType,
			Data:        &pb.ChannelInfo_NodeInfo{NodeInfo: &ni},
		})
	})

	s.Lock()
	if _, ok := s.mapChanInfo[in.Addr]; ok {
		close(s.mapChanInfo[in.Addr])
	}

	chanInfo := make(chan pb.ChannelInfo, 16)
	s.mapChanInfo[in.Addr] = chanInfo
	s.Unlock()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case ci, ok := <-chanInfo:
			if !ok {
				return nil
			}

			stream.SendMsg(&ci)

			// if ci.ChannelType == pb.CHANNELTYPE_NODEINFO {

			// }

			// ni := pb.NodeInfo{
			// 	ServAddr: bi.ServAddr,
			// 	Addr:     bi.Addr,
			// 	Name:     bi.Name,
			// 	NodeType: bi.NodeType,
			// }

			// stream.SendMsg(&pb.ChannelInfo{
			// 	ChannelType: in.ChannelType,
			// 	Data:        &pb.ChannelInfo_NodeInfo{NodeInfo: &ni},
			// })
		}

		// gs.Send(&pb.HelloReply{Message: "Hello " + in.Name})
	}
	// return nil
}

// RequestCtrl implements jarviscorepb.JarvisCoreServ
func (s *jarvisServer) RequestCtrl(ctx context.Context, in *pb.CtrlInfo) (*pb.BaseReply, error) {
	if in.DestAddr == s.node.myinfo.Addr {
		_, err := mgrCtrl.Run(in.CtrlType, in.Command)
		if err != nil {
			return &pb.BaseReply{}, nil
		}

		return &pb.BaseReply{}, nil
	}

	if s.hasClient(in.DestAddr) {
		return &pb.BaseReply{
			ReplyType: pb.REPLYTYPE_FORWARD,
		}, nil
	}

	return &pb.BaseReply{
		ReplyType: pb.REPLYTYPE_BROADCAST,
	}, nil
}

// ReplyCtrl implements jarviscorepb.JarvisCoreServ
func (s *jarvisServer) ReplyCtrl(ctx context.Context, in *pb.CtrlResult) (*pb.BaseReply, error) {
	return &pb.BaseReply{}, nil
}
