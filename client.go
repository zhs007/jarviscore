package jarviscore

import (
	"context"
	"sync"
	"time"

	pb "github.com/zhs007/jarviscore/proto"
	"google.golang.org/grpc"
)

type clientInfo struct {
	conn   *grpc.ClientConn
	client pb.JarvisCoreServClient
}

// jarvisClient -
type jarvisClient struct {
	peeraddrmgr *peerAddrMgr
	mapClient   map[string]*clientInfo
	clientchan  chan int
	wg          sync.WaitGroup
}

func newClient() *jarvisClient {
	return &jarvisClient{
		mapClient:  make(map[string]*clientInfo),
		clientchan: make(chan int, 1)}
}

func (c *jarvisClient) onConnectFail(addr string) {
}

func (c *jarvisClient) Start(peeraddrmgr *peerAddrMgr, myinfo *BaseInfo) error {
	c.peeraddrmgr = peeraddrmgr

	for _, v := range peeraddrmgr.arr.PeerAddr {
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

	conn, err := mgrconn.getConn(servaddr)
	if err != nil {
		warnLog("JarvisClient.connect", err)

		return err
	}

	ci := clientInfo{conn: conn, client: pb.NewJarvisCoreServClient(conn)}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	r, err1 := ci.client.Join(ctx, &pb.Join{Servaddr: myinfo.ServAddr, Token: myinfo.Token, Name: myinfo.Name, Nodetype: myinfo.NodeType})
	if err1 != nil {
		warnLog("JarvisClient.connect", err1)

		return err1
	}

	if r.Code == pb.CODE_OK {
		return nil
	}

	return nil
}
