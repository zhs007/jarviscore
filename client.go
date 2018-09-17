package jarviscore

import (
	"context"
	"io"
	"time"

	"go.uber.org/zap"

	"github.com/zhs007/jarviscore/err"
	"github.com/zhs007/jarviscore/log"
	pb "github.com/zhs007/jarviscore/proto"
	"google.golang.org/grpc"
)

type clientInfo struct {
	conn   *grpc.ClientConn
	client pb.JarvisCoreServClient
}

// jarvisClient -
type jarvisClient struct {
	// sync.RWMutex

	node *jarvisNode
	// mgrpeeraddr  *peerAddrMgr
	mapClient    map[string]*clientInfo
	clientchan   chan int
	activitychan chan int
	activitynums int
	waitchan     chan int
	connectchan  chan *BaseInfo
}

func newClient(node *jarvisNode) *jarvisClient {
	return &jarvisClient{
		node:         node,
		mapClient:    make(map[string]*clientInfo),
		clientchan:   make(chan int, 1),
		activitychan: make(chan int, 16),
		activitynums: 0,
		waitchan:     make(chan int, 1),
		connectchan:  make(chan *BaseInfo, 16),
	}
}

func (c *jarvisClient) pushNewConnect(bi *BaseInfo) {
	c.connectchan <- bi
}

func (c *jarvisClient) onConnectFail(addr string) {
}

func (c *jarvisClient) onStop() {
	c.clientchan <- 0
}

func (c *jarvisClient) startConnectAllNode() {
	c.node.mgrNodeInfo.foreach(func(cn *NodeInfo) {
		c.pushNewConnect(&cn.baseinfo)
	})
}

func (c *jarvisClient) Start() {
	ctx, cancel := context.WithCancel(context.Background())

	go c.connectRoot(ctx, config.DefPeerAddr)
	go c.startConnectAllNode()
	// c.mgrpeeraddr = mgrpeeraddr

	// for _, v := range mgrpeeraddr.arr.PeerAddr {
	// 	go c.connect(ctx, v)
	// 	time.Sleep(time.Second)
	// }

	stopping := false

	for {
		select {
		case <-c.waitchan:
			cancel()

			if c.activitynums == 0 {
				c.onStop()
				break
			}

			stopping = true
		case v := <-c.connectchan:
			go c.connect(ctx, v)
			time.Sleep(time.Second)
		case v, ok := <-c.activitychan:
			if !ok {
				c.onStop()
				break
			}

			c.activitynums += v

			if stopping && c.activitynums == 0 {
				c.onStop()
				break
			}
		}
	}
}

func (c *jarvisClient) Stop() {
	c.waitchan <- 0
}

func (c *jarvisClient) onConnectEnd() {
	c.activitychan <- -1
}

//
func (c *jarvisClient) connectRoot(ctx context.Context, servaddr string) error {
	log.Info("jarvisClient.connectRoot", zap.String("servaddr", servaddr))

	c.activitychan <- 1
	defer c.onConnectEnd()

	// c.node.mgrNodeInfo.onStartConnect(bi.Addr)
	// c.mgrpeeraddr.onStartConnect(servaddr)

	conn, err := mgrconn.getConn(servaddr)
	if err != nil {
		jarviserr.WarnLog("JarvisClient.connectRoot:getConn", err)

		return err
	}

	curctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ci := clientInfo{
		conn:   conn,
		client: pb.NewJarvisCoreServClient(conn),
	}

	r, err1 := ci.client.Join(curctx, &pb.Join{
		ServAddr: c.node.myinfo.ServAddr,
		Addr:     c.node.myinfo.Addr,
		Name:     c.node.myinfo.Name,
		NodeType: c.node.myinfo.NodeType})
	if err1 != nil {
		jarviserr.WarnLog("JarvisClient.connectRoot:Join", err1)

		mgrconn.delConn(servaddr)

		return err1
	}

	if r.Code == pb.CODE_OK {
		log.Info("jarvisClient.connectRoot:OK")

		c.node.onIConnectNode(&BaseInfo{
			Name:     r.Name,
			Addr:     r.Addr,
			NodeType: r.NodeType,
			ServAddr: servaddr,
		})

		// c.node.mgrNodeInfo.onConnected(bi.Addr)
		// c.mgrpeeraddr.onConnected(servaddr)

		c.subscribe(ctx, &ci, pb.CHANNELTYPE_NODEINFO)

		return nil
	}

	log.Info("JarvisClient.connectRoot:Join", zap.Int("code", int(r.Code)))

	mgrconn.delConn(servaddr)

	return nil
}

//
func (c *jarvisClient) connect(ctx context.Context, bi *BaseInfo) error {
	log.Info("jarvisClient.connect", zap.String("servaddr", bi.ServAddr))

	c.activitychan <- 1
	defer c.onConnectEnd()

	c.node.mgrNodeInfo.onStartConnect(bi.Addr)
	// c.mgrpeeraddr.onStartConnect(servaddr)

	conn, err := mgrconn.getConn(bi.ServAddr)
	if err != nil {
		jarviserr.WarnLog("JarvisClient.connect:getConn", err)

		return err
	}

	curctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ci := clientInfo{
		conn:   conn,
		client: pb.NewJarvisCoreServClient(conn),
	}

	r, err1 := ci.client.Join(curctx, &pb.Join{
		ServAddr: c.node.myinfo.ServAddr,
		Addr:     c.node.myinfo.Addr,
		Name:     c.node.myinfo.Name,
		NodeType: c.node.myinfo.NodeType})
	if err1 != nil {
		jarviserr.WarnLog("JarvisClient.connect:Join", err1)

		mgrconn.delConn(bi.ServAddr)

		return err1
	}

	if r.Code == pb.CODE_OK {
		log.Info("jarvisClient.connect:OK")

		c.node.onIConnectNode(&BaseInfo{
			Name:     r.Name,
			Addr:     r.Addr,
			NodeType: r.NodeType,
			ServAddr: bi.ServAddr,
		})

		c.node.mgrNodeInfo.onConnected(bi.Addr)
		// c.mgrpeeraddr.onConnected(servaddr)

		c.subscribe(ctx, &ci, pb.CHANNELTYPE_NODEINFO)

		return nil
	}

	log.Info("JarvisClient.Join", zap.Int("code", int(r.Code)))

	mgrconn.delConn(bi.ServAddr)

	return nil
}

func (c *jarvisClient) subscribe(ctx context.Context, ci *clientInfo, ct pb.CHANNELTYPE) error {
	log.Info("jarvisClient.subscribe")

	curctx, cancel := context.WithCancel(ctx)
	defer cancel()

	stream, err := ci.client.Subscribe(curctx, &pb.Subscribe{
		ChannelType: ct,
		Addr:        c.node.myinfo.Addr,
	})
	if err != nil {
		jarviserr.WarnLog("JarvisClient.subscribe:Subscribe", err)
		return err
	}

	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			log.Info("jarvisClient.subscribe:stream.EOF")
			break
		}
		if err != nil {
			jarviserr.WarnLog("JarvisClient.subscribe:stream.Recv", err)

			return err
		}

		if reply.ChannelType == ct {
			log.Info("jarvisClient.subscribe:NodeInfo")

			ni := reply.GetNodeInfo()
			if ni != nil {
				log.Info("JarvisClient.subscribe:NODEINFO",
					zap.String("Servaddr", ni.ServAddr),
					zap.String("Addr", ni.Addr),
					zap.String("Name", ni.Name),
					zap.Int("Nodetype", int(ni.NodeType)))

				c.node.onGetNewNode(&BaseInfo{
					Name:     ni.Name,
					ServAddr: ni.ServAddr,
					Addr:     ni.Addr,
					NodeType: ni.NodeType,
				})
			}
		}
		// log.Printf("Greeting: %s", reply.Message)
	}

	return nil
}

func (c *jarvisClient) sendCtrl(ctx context.Context, ci *pb.CtrlInfo) error {
	curclient, ok := c.mapClient[ci.DestAddr]
	if !ok {
		return jarviserr.NewError(pb.CODE_NOT_CONNECTNODE)
	}

	curctx, cancel := context.WithCancel(ctx)
	defer cancel()

	r, err1 := curclient.client.RequestCtrl(curctx, ci)
	if err1 != nil {
		return err1
	}

	if r.Code == pb.CODE_OK {
		return nil
	}

	return nil
}

func (c *jarvisClient) sendCtrlResult(token string) error {
	return nil
}
