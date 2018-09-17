package jarviscore

import (
	"github.com/zhs007/jarviscore/err"
	"google.golang.org/grpc"
)

// connMgr -
type connMgr struct {
	mapConn map[string]*grpc.ClientConn
}

func (mgr connMgr) getConn(servaddr string) (*grpc.ClientConn, error) {
	if conn, ok := mgr.mapConn[servaddr]; ok {
		return conn, nil
	}

	conn, err := grpc.Dial(servaddr, grpc.WithInsecure())
	if err != nil {
		jarviserr.WarnLog("connMgr.getConn", err)

		return nil, err
	}

	mgr.mapConn[servaddr] = conn

	return conn, nil
}

func (mgr connMgr) delConn(servaddr string) {
	_, ok := mgr.mapConn[servaddr]
	if ok {
		delete(mgr.mapConn, servaddr)
	}
}
