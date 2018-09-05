package jarviscore

import (
	"context"
	"io"
	"time"

	"go.uber.org/zap"

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

	node         *jarvisNode
	mgrpeeraddr  *peerAddrMgr
	mapClient    map[string]*clientInfo
	clientchan   chan int
	activitychan chan int
	activitynums int
	waitchan     chan int
	connectchan  chan string
}

func newClient(node *jarvisNode) *jarvisClient {
	return &jarvisClient{
		node:         node,
		mapClient:    make(map[string]*clientInfo),
		clientchan:   make(chan int, 1),
		activitychan: make(chan int, 16),
		activitynums: 0,
		waitchan:     make(chan int, 1),
		connectchan:  make(chan string, 16),
	}
}

func (c *jarvisClient) pushNewConnect(servaddr string) {
	c.connectchan <- servaddr
}

func (c *jarvisClient) onConnectFail(addr string) {
}

func (c *jarvisClient) onStop() {
	c.clientchan <- 0
}

func (c *jarvisClient) Start(mgrpeeraddr *peerAddrMgr) {
	ctx, cancel := context.WithCancel(context.Background())

	c.mgrpeeraddr = mgrpeeraddr

	for _, v := range mgrpeeraddr.arr.PeerAddr {
		go c.connect(ctx, v)
		time.Sleep(time.Second)
	}

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
func (c *jarvisClient) connect(ctx context.Context, servaddr string) error {
	log.Info("jarvisClient.connect")

	c.activitychan <- 1
	defer c.onConnectEnd()

	c.mgrpeeraddr.onStartConnect(servaddr)

	conn, err := mgrconn.getConn(servaddr)
	if err != nil {
		warnLog("JarvisClient.connect:getConn", err)

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
		Token:    c.node.myinfo.Token,
		Name:     c.node.myinfo.Name,
		NodeType: c.node.myinfo.NodeType})
	if err1 != nil {
		warnLog("JarvisClient.connect:Join", err1)

		mgrconn.delConn(servaddr)

		return err1
	}

	if r.Code == pb.CODE_OK {
		log.Info("jarvisClient.connect:OK")

		c.node.onIConnectNode(&BaseInfo{
			Name:     r.Name,
			Token:    r.Token,
			NodeType: r.NodeType,
			ServAddr: servaddr,
		})

		c.mgrpeeraddr.onConnected(servaddr)

		c.subscribe(ctx, &ci, pb.CHANNELTYPE_NODEINFO)

		return nil
	}

	log.Info("JarvisClient.Join", zap.Int("code", int(r.Code)))

	mgrconn.delConn(servaddr)

	return nil
}

func (c *jarvisClient) subscribe(ctx context.Context, ci *clientInfo, ct pb.CHANNELTYPE) error {
	log.Info("jarvisClient.subscribe")

	curctx, cancel := context.WithCancel(ctx)
	defer cancel()

	stream, err := ci.client.Subscribe(curctx, &pb.Subscribe{
		ChannelType: ct,
		Token:       c.node.myinfo.Token,
	})
	if err != nil {
		warnLog("JarvisClient.subscribe:Subscribe", err)
		return err
	}

	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			log.Info("jarvisClient.subscribe:stream.EOF")
			break
		}
		if err != nil {
			warnLog("JarvisClient.subscribe:stream.Recv", err)

			return err
		}

		if reply.ChannelType == ct {
			log.Info("jarvisClient.subscribe:NodeInfo")

			ni := reply.GetNodeInfo()
			if ni != nil {
				log.Info("JarvisClient.subscribe:NODEINFO",
					zap.String("Servaddr", ni.ServAddr),
					zap.String("Token", ni.Token),
					zap.String("Name", ni.Name),
					zap.Int("Nodetype", int(ni.NodeType)))

				c.node.onGetNewNode(&BaseInfo{
					Name:     ni.Name,
					ServAddr: ni.ServAddr,
					Token:    ni.Token,
					NodeType: ni.NodeType,
				})
			}
		}
		// log.Printf("Greeting: %s", reply.Message)
	}

	return nil
}

func (c *jarvisClient) sendCtrl(token string) error {
	return nil
}

func (c *jarvisClient) sendCtrlResult(token string) error {
	return nil
}
