package jarviscore

import (
	"sync"

	"google.golang.org/grpc/connectivity"

	"github.com/zhs007/jarviscore/base"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// connMgr -
type connMgr struct {
	sync.RWMutex

	mapConn map[string]*grpc.ClientConn
}

func (mgr *connMgr) isValidConn(servaddr string) bool {
	mgr.Lock()
	defer mgr.Unlock()

	if conn, ok := mgr.mapConn[servaddr]; ok {
		cs := conn.GetState()
		if !(cs == connectivity.Shutdown || cs == connectivity.TransientFailure) {
			return true
		}
	}

	return false
}

func (mgr *connMgr) getConn(servaddr string) (*grpc.ClientConn, error) {
	mgr.Lock()
	defer mgr.Unlock()

	if conn, ok := mgr.mapConn[servaddr]; ok {
		cs := conn.GetState()
		if !(cs == connectivity.Shutdown || cs == connectivity.TransientFailure) {
			return conn, nil
		}

		jarvisbase.Warn("connMgr.getConn:InvalidConn",
			zap.Int("state", int(cs)))

		err := conn.Close()
		if err != nil {
			jarvisbase.Warn("connMgr.getConn:Close", zap.Error(err))
		}
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
