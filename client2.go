package jarviscore

import (
	"context"

	pb "github.com/zhs007/jarviscore/proto"
	"google.golang.org/grpc"
)

type clientInfo2 struct {
	conn   *grpc.ClientConn
	client pb.JarvisCoreServClient
}

// jarvisClient2 -
type jarvisClient2 struct {
	// sync.RWMutex

	node      JarvisNode
	mapClient map[string]*clientInfo2
	// clientchan   chan int
	// activitychan chan int
	// activitynums int
	// waitchan     chan int
	// connectchan  chan *BaseInfo
}

func newClient2(node JarvisNode) *jarvisClient2 {
	return &jarvisClient2{
		node:      node,
		mapClient: make(map[string]*clientInfo2),
		// clientchan:   make(chan int, 1),
		// activitychan: make(chan int, 16),
		// activitynums: 0,
		// waitchan:     make(chan int, 1),
		// connectchan:  make(chan *BaseInfo, 16),
	}
}

func (c *jarvisClient2) sendMsg(ctx context.Context, msg *pb.JarvisMsg) error {
	ci2, ok := c.mapClient[msg.DestAddr]
	if !ok {
		return c.broadCastMsg(ctx, msg)
	}

	ci2.client.ProcMsg(ctx, msg)

	return nil
}

func (c *jarvisClient2) broadCastMsg(ctx context.Context, msg *pb.JarvisMsg) error {
	for _, v := range c.mapClient {
		v.client.ProcMsg(ctx, msg)
	}

	return nil
}

// func (c *jarvisClient2) pushNewConnect(bi *BaseInfo) {
// 	c.connectchan <- bi
// }

// func (c *jarvisClient2) onConnectFail(addr string) {
// }

// func (c *jarvisClient2) onStop() {
// 	c.clientchan <- 0
// }

// func (c *jarvisClient2) startConnectAllNode() {
// 	c.node.mgrNodeInfo.foreach(func(cn *coredbpb.NodeInfo) {
// 		bi := &BaseInfo{
// 			Name:     cn.Name,
// 			ServAddr: cn.ServAddr,
// 			Addr:     cn.Addr,
// 		}
// 		c.pushNewConnect(bi)
// 	})
// }

// func (c *jarvisClient2) Start(ctx context.Context) {
// 	ctx, cancel := context.WithCancel(ctx)

// 	go c.connectRoot(ctx, config.RootServAddr)
// 	go c.startConnectAllNode()
// 	// c.mgrpeeraddr = mgrpeeraddr

// 	// for _, v := range mgrpeeraddr.arr.PeerAddr {
// 	// 	go c.connect(ctx, v)
// 	// 	time.Sleep(time.Second)
// 	// }

// 	stopping := false

// 	for {
// 		select {
// 		case <-c.waitchan:
// 			cancel()

// 			if c.activitynums == 0 {
// 				c.onStop()
// 				break
// 			}

// 			stopping = true
// 		case v := <-c.connectchan:
// 			go c.connect(ctx, v)
// 			time.Sleep(time.Second)
// 		case v, ok := <-c.activitychan:
// 			if !ok {
// 				c.onStop()
// 				break
// 			}

// 			c.activitynums += v

// 			if stopping && c.activitynums == 0 {
// 				c.onStop()
// 				break
// 			}
// 		}
// 	}
// }

// func (c *jarvisClient2) Stop() {
// 	c.waitchan <- 0
// }

// func (c *jarvisClient2) onConnectEnd() {
// 	c.activitychan <- -1
// }

// //
// func (c *jarvisClient2) connectRoot(ctx context.Context, servaddr string) error {
// 	jarvisbase.Info("jarvisClient.connectRoot", zap.String("servaddr", servaddr))

// 	c.activitychan <- 1
// 	defer c.onConnectEnd()

// 	// c.node.mgrNodeInfo.onStartConnect(bi.Addr)
// 	// c.mgrpeeraddr.onStartConnect(servaddr)

// 	conn, err := mgrconn.getConn(servaddr)
// 	if err != nil {
// 		jarvisbase.Warn("JarvisClient.connectRoot:getConn", zap.Error(err))

// 		return err
// 	}

// 	// jarvisbase.Info("jarvisClient.connectRoot:getConn:ok")

// 	curctx, cancel := context.WithCancel(ctx)
// 	defer cancel()

// 	ci := clientInfo{
// 		conn:   conn,
// 		client: pb.NewJarvisCoreServClient(conn),
// 	}

// 	// jarvisbase.Info("jarvisClient.connectRoot:getConn:clientinfo ok")

// 	r, err1 := ci.client.Join(curctx, &pb.Join{
// 		ServAddr: c.node.myinfo.ServAddr,
// 		Addr:     c.node.myinfo.Addr,
// 		Name:     c.node.myinfo.Name,
// 	})
// 	if err1 != nil {
// 		jarvisbase.Warn("JarvisClient.connectRoot:Join", zap.Error(err1))

// 		mgrconn.delConn(servaddr)

// 		return err1
// 	}

// 	jarvisbase.Info("jarvisClient.connectRoot",
// 		jarvisbase.JSON("result", r))

// 	if r.Err == "" {
// 		jarvisbase.Info("jarvisClient.connectRoot:OK")

// 		c.mapClient[r.Addr] = &ci

// 		c.node.onIConnectNode(&BaseInfo{
// 			Name:     r.Name,
// 			Addr:     r.Addr,
// 			ServAddr: servaddr,
// 		})

// 		// c.node.requestCtrl(ctx, r.Addr, pb.CTRLTYPE_SHELL, []byte("whami"))

// 		// c.node.mgrNodeInfo.onConnected(bi.Addr)
// 		// c.mgrpeeraddr.onConnected(servaddr)

// 		c.subscribe(ctx, &ci, pb.CHANNELTYPE_NODEINFO)

// 		return nil
// 	}

// 	jarvisbase.Info("JarvisClient.connectRoot:Join",
// 		zap.String("err", r.Err))

// 	mgrconn.delConn(servaddr)

// 	return nil
// }

// //
// func (c *jarvisClient2) connect(ctx context.Context, bi *BaseInfo) error {
// 	jarvisbase.Info("jarvisClient.connect", zap.String("servaddr", bi.ServAddr))

// 	c.activitychan <- 1
// 	defer c.onConnectEnd()

// 	c.node.mgrNodeInfo.onStartConnect(bi.Addr)
// 	// c.mgrpeeraddr.onStartConnect(servaddr)

// 	conn, err := mgrconn.getConn(bi.ServAddr)
// 	if err != nil {
// 		jarvisbase.Warn("JarvisClient.connect:getConn", zap.Error(err))

// 		return err
// 	}

// 	curctx, cancel := context.WithCancel(ctx)
// 	defer cancel()

// 	ci := clientInfo{
// 		conn:   conn,
// 		client: pb.NewJarvisCoreServClient(conn),
// 	}

// 	r, err1 := ci.client.Join(curctx, &pb.Join{
// 		ServAddr: c.node.myinfo.ServAddr,
// 		Addr:     c.node.myinfo.Addr,
// 		Name:     c.node.myinfo.Name,
// 	})
// 	if err1 != nil {
// 		jarvisbase.Warn("JarvisClient.connect:Join", zap.Error(err1))

// 		mgrconn.delConn(bi.ServAddr)

// 		return err1
// 	}

// 	if r.Err == "" {
// 		jarvisbase.Info("jarvisClient.connect:OK")

// 		c.node.onIConnectNode(&BaseInfo{
// 			Name:     r.Name,
// 			Addr:     r.Addr,
// 			ServAddr: bi.ServAddr,
// 		})

// 		c.node.mgrNodeInfo.onConnected(bi.Addr)
// 		// c.mgrpeeraddr.onConnected(servaddr)

// 		c.subscribe(ctx, &ci, pb.CHANNELTYPE_NODEINFO)

// 		return nil
// 	}

// 	jarvisbase.Info("JarvisClient.Join", zap.String("err", r.Err))

// 	mgrconn.delConn(bi.ServAddr)

// 	return nil
// }

// func (c *jarvisClient2) subscribe(ctx context.Context, ci *clientInfo, ct pb.CHANNELTYPE) error {
// 	jarvisbase.Info("jarvisClient.subscribe")

// 	curctx, cancel := context.WithCancel(ctx)
// 	defer cancel()

// 	stream, err := ci.client.Subscribe(curctx, &pb.Subscribe{
// 		ChannelType: ct,
// 		Addr:        c.node.myinfo.Addr,
// 	})
// 	if err != nil {
// 		jarvisbase.Warn("JarvisClient.subscribe:Subscribe", zap.Error(err))
// 		return err
// 	}

// 	for {
// 		reply, err := stream.Recv()
// 		if err == io.EOF {
// 			jarvisbase.Info("jarvisClient.subscribe:stream.EOF")
// 			break
// 		}
// 		if err != nil {
// 			jarvisbase.Warn("JarvisClient.subscribe:stream.Recv", zap.Error(err))

// 			return err
// 		}

// 		if reply.ChannelType == ct {
// 			jarvisbase.Info("jarvisClient.subscribe:NodeInfo")

// 			ni := reply.GetNodeInfo()
// 			if ni != nil {
// 				jarvisbase.Info("JarvisClient.subscribe:NODEINFO",
// 					zap.String("Servaddr", ni.ServAddr),
// 					zap.String("Addr", ni.Addr),
// 					zap.String("Name", ni.Name),
// 				)

// 				c.node.onGetNewNode(&BaseInfo{
// 					Name:     ni.Name,
// 					ServAddr: ni.ServAddr,
// 					Addr:     ni.Addr,
// 				})
// 			}
// 		}
// 		// log.Printf("Greeting: %s", reply.Message)
// 	}

// 	return nil
// }

// func (c *jarvisClient2) sendCtrl(ctx context.Context, ci *pb.CtrlInfo) error {
// 	curclient, ok := c.mapClient[ci.DestAddr]
// 	if !ok {
// 		return ErrNotConnectNode
// 	}

// 	curctx, cancel := context.WithCancel(ctx)
// 	defer cancel()

// 	r, err1 := curclient.client.RequestCtrl(curctx, ci)
// 	if err1 != nil {
// 		return err1
// 	}

// 	if r.Err == "" {
// 		return nil
// 	}

// 	return nil
// }

// func (c *jarvisClient2) sendCtrlResult(token string) error {
// 	return nil
// }
