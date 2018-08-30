package jarviscore

import (
	"context"
	"io"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/zhs007/jarviscore/log"
	pb "github.com/zhs007/jarviscore/proto"
	"google.golang.org/grpc"
)

// clientState -
type clientState int

const (
	// clientStateNormal - normal
	clientStateNormal clientState = iota
	// clientStateRunning - running
	clientStateRunning
	// clientStateStop - stop
	clientStateStop
	// clientStateStopped - stopped
	clientStateStopped
)

type clientInfo struct {
	// ctx    *context.Context
	conn   *grpc.ClientConn
	client pb.JarvisCoreServClient
}

// jarvisClient -
type jarvisClient struct {
	sync.RWMutex

	node         *jarvisNode
	mgrpeeraddr  *peerAddrMgr
	mapClient    map[string]*clientInfo
	clientchan   chan int
	activitychan chan int
	activitynums int
	state        clientState
	// wg          sync.WaitGroup
}

func newClient(node *jarvisNode) *jarvisClient {
	return &jarvisClient{
		node:         node,
		mapClient:    make(map[string]*clientInfo),
		clientchan:   make(chan int, 1),
		activitychan: make(chan int, 16),
		activitynums: 0,
		state:        clientStateNormal,
	}
}

func (c *jarvisClient) onConnectFail(addr string) {
}

func (c *jarvisClient) Start(mgrpeeraddr *peerAddrMgr) error {
	c.Lock()
	if c.state != clientStateNormal {
		c.Unlock()
		return newError(int(pb.CODE_CLIENT_ALREADY_START))
	}

	c.state = clientStateRunning
	c.Unlock()

	c.mgrpeeraddr = mgrpeeraddr

	for _, v := range mgrpeeraddr.arr.PeerAddr {
		go c.connect(v)
		time.Sleep(time.Second)
	}

	c.waitEnd()

	c.Lock()
	c.state = clientStateStopped
	c.Unlock()

	c.clientchan <- 0

	return nil
}

func (c *jarvisClient) Stop() {
	c.RLock()
	defer c.RUnlock()

	if c.state == clientStateRunning {
		return
	}

	if c.activitynums == 0 {
		close(c.activitychan)
	}

	c.Lock()
	c.state = clientStateStop
	c.Unlock()
}

func (c *jarvisClient) onConnectEnd() {
	c.activitychan <- -1
}

func (c *jarvisClient) waitEnd() {
	for {
		select {
		case v, ok := <-c.activitychan:
			c.activitynums += v

			c.RLock()
			if c.state == clientStateStop && c.activitynums == 0 {
				c.RUnlock()
				break
			}
			c.RUnlock()

			if !ok {
				break
			}
		}
	}
}

//
func (c *jarvisClient) connect(servaddr string) error {
	c.activitychan <- 1
	defer c.onConnectEnd()

	c.mgrpeeraddr.onStartConnect(servaddr)

	conn, err := mgrconn.getConn(servaddr)
	if err != nil {
		warnLog("JarvisClient.connect:getConn", err)

		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ci := clientInfo{
		// ctx: &ctx,
		conn:   conn,
		client: pb.NewJarvisCoreServClient(conn),
	}

	r, err1 := ci.client.Join(ctx, &pb.Join{
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
		c.node.addNode(&BaseInfo{
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
	curctx, cancel := context.WithCancel(ctx)
	defer cancel()

	stream, err := ci.client.Subscribe(curctx, &pb.Subscribe{
		ChannelType: ct,
		Token:       c.node.myinfo.Token,
	})
	if err != nil {
		return err
	}

	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if reply.ChannelType == ct {
			ni := reply.GetNodeInfo()
			if ni != nil {
				log.Info("JarvisClient.subscribe:NODEINFO",
					zap.String("Servaddr", ni.ServAddr),
					zap.String("Token", ni.Token),
					zap.String("Name", ni.Name),
					zap.Int("Nodetype", int(ni.NodeType)))

				bi := BaseInfo{
					Name:     ni.Name,
					ServAddr: ni.ServAddr,
					Token:    ni.Token,
					NodeType: ni.NodeType,
				}

				c.node.onAddNode(&bi)
			}
		}
		// log.Printf("Greeting: %s", reply.Message)
	}

	return nil
}
