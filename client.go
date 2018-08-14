package jarviscore

import (
	"context"
	"sync"
	"time"

	pb "github.com/zhs007/jarviscore/proto"
	"google.golang.org/grpc"
)

// jarvisClient -
type jarvisClient struct {
	peeraddrmgr *peerAddrMgr
	mapConn     map[string]*grpc.ClientConn
	mapClient   map[string]pb.JarvisCoreServClient
	clientchan  chan int
	wg          sync.WaitGroup
}

func newClient() *jarvisClient {
	return &jarvisClient{
		mapConn:    make(map[string]*grpc.ClientConn),
		mapClient:  make(map[string]pb.JarvisCoreServClient),
		clientchan: make(chan int, 1)}
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

	var curconn *grpc.ClientConn
	if _, ok := c.mapConn[servaddr]; ok {
		curconn = c.mapConn[servaddr]
	} else {
		conn, err := grpc.Dial(servaddr, grpc.WithInsecure())
		if err != nil {
			warnLog("JarvisClient.connect", err)

			return err
		}
		curconn = conn
	}

	jarvisclient := pb.NewJarvisCoreServClient(curconn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err1 := jarvisclient.Join(ctx, &pb.Join{Servaddr: myinfo.ServAddr, Token: myinfo.Token, Name: myinfo.Name, Nodetype: myinfo.NodeType})
	if err1 != nil {
		warnLog("JarvisClient.connect", err1)

		return err1
	}

	if r.Code == pb.CODE_OK {
		return nil
	}

	return nil
}
