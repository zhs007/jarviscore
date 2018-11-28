package jarviscore

import (
	"sync"

	"github.com/zhs007/jarviscore/base"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// connMgr -
type connMgr struct {
	sync.RWMutex

	mapConn map[string]*grpc.ClientConn
}

func (mgr *connMgr) getConn(servaddr string) (*grpc.ClientConn, error) {
	mgr.Lock()
	defer mgr.Unlock()

	if conn, ok := mgr.mapConn[servaddr]; ok {
		return conn, nil
	}

	conn, err := grpc.Dial(servaddr, grpc.WithInsecure())
	if err != nil {
		jarvisbase.Warn("connMgr.getConn", zap.Error(err))

		return nil, err
	}

	mgr.mapConn[servaddr] = conn

	return conn, nil
}

func (mgr *connMgr) delConn(servaddr string) {
	mgr.Lock()
	defer mgr.Unlock()

	_, ok := mgr.mapConn[servaddr]
	if ok {
		delete(mgr.mapConn, servaddr)
	}
}
