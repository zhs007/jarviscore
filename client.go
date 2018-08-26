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

type clientInfo struct {
	// ctx    *context.Context
	conn   *grpc.ClientConn
	client pb.JarvisCoreServClient
}

// jarvisClient -
type jarvisClient struct {
	node        *jarvisNode
	mgrpeeraddr *peerAddrMgr
	mapClient   map[string]*clientInfo
	clientchan  chan int
	wg          sync.WaitGroup
}

func newClient(node *jarvisNode) *jarvisClient {
	return &jarvisClient{
		node:       node,
		mapClient:  make(map[string]*clientInfo),
		clientchan: make(chan int, 1)}
}

func (c *jarvisClient) onConnectFail(addr string) {
}

func (c *jarvisClient) Start(mgrpeeraddr *peerAddrMgr) error {
	c.mgrpeeraddr = mgrpeeraddr

	for _, v := range mgrpeeraddr.arr.PeerAddr {
		go c.connect(v)
		time.Sleep(time.Second)
	}

	c.wg.Wait()

	c.clientchan <- 0

	return nil
}

//
func (c *jarvisClient) connect(servaddr string) error {
	c.wg.Add(1)
	defer c.wg.Done()

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

		if reply.ChannelType == ct && reply.NodeInfo != nil {
			log.Info("JarvisClient.subscribe:NODEINFO",
				zap.String("Servaddr", reply.NodeInfo.ServAddr),
				zap.String("Token", reply.NodeInfo.Token),
				zap.String("Name", reply.NodeInfo.Name),
				zap.Int("Nodetype", int(reply.NodeInfo.NodeType)))

			bi := BaseInfo{
				Name:     reply.NodeInfo.Name,
				ServAddr: reply.NodeInfo.ServAddr,
				Token:    reply.NodeInfo.Token,
				NodeType: reply.NodeInfo.NodeType,
			}

			c.node.onAddNode(&bi)
		}
		// log.Printf("Greeting: %s", reply.Message)
	}

	return nil
}
