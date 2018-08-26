package jarviscore

import (
	"context"
	"sync"
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

func (c *jarvisClient) Start(mgrpeeraddr *peerAddrMgr, myinfo *BaseInfo) error {
	c.mgrpeeraddr = mgrpeeraddr

	for _, v := range mgrpeeraddr.arr.PeerAddr {
		go c.connect(v, myinfo)
		time.Sleep(time.Second)
	}

	c.wg.Wait()

	c.clientchan <- 0

	return nil
}

//
func (c *jarvisClient) connect(servaddr string, myinfo *BaseInfo) error {
	c.wg.Add(1)
	defer c.wg.Done()

	c.mgrpeeraddr.onStartConnect(servaddr)

	conn, err := mgrconn.getConn(servaddr)
	if err != nil {
		warnLog("JarvisClient.connect:getConn", err)

		return err
	}

	ci := clientInfo{conn: conn, client: pb.NewJarvisCoreServClient(conn)}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	r, err1 := ci.client.Join(ctx, &pb.Join{ServAddr: myinfo.ServAddr, Token: myinfo.Token, Name: myinfo.Name, NodeType: myinfo.NodeType})
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

		return nil
	}

	log.Info("JarvisClient.Join", zap.Int("code", int(r.Code)))

	mgrconn.delConn(servaddr)

	return nil
}
