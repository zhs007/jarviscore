package jarviscore

import (
	"context"
	"time"

	pb "github.com/zhs007/jarviscore/proto"
	"google.golang.org/grpc"
)

// JarvisClient -
type JarvisClient struct {
	mapConn   map[string]*grpc.ClientConn
	mapClient map[string]pb.JarvisCoreServClient
}

//
func (c *JarvisClient) Connect(servaddr string) error {
	var curconn *grpc.ClientConn
	if _, ok := c.mapConn[servaddr]; ok {
		curconn = c.mapConn[servaddr]
	} else {
		conn, err := grpc.Dial(servaddr, grpc.WithInsecure())
		if err != nil {
			return err
		}
		curconn = conn
	}

	jarvisclient := pb.NewJarvisCoreServClient(curconn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err1 := jarvisclient.Join(ctx, &pb.Join{Servaddr: "", Token: "", Name: "", Nodetype: pb.NODETYPE_NORMAL})
	if err1 != nil {
		return err1
	}

	if r.Code == pb.CODE_OK {
		return nil
	}

	return nil
}
