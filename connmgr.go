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
	mapConn sync.Map
}

func (mgr *connMgr) _getConn(servaddr string) *grpc.ClientConn {
	val, ok := mgr.mapConn.Load(servaddr)
	if ok {
		conn, typeok := val.(*grpc.ClientConn)
		if typeok {
			return conn
		}
	}

	return nil
}

func (mgr *connMgr) isValidConn(servaddr string) bool {
	conn := mgr._getConn(servaddr)
	if conn == nil {
		return false
	}

	cs := conn.GetState()
	if !(cs == connectivity.Shutdown || cs == connectivity.TransientFailure) {
		return true
	}

	return false
}

func (mgr *connMgr) getConn(servaddr string) (*grpc.ClientConn, error) {
	conn := mgr._getConn(servaddr)
	if conn != nil {
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

	mgr.mapConn.Store(servaddr, conn)

	return conn, nil
}

func (mgr *connMgr) delConn(servaddr string) {
	mgr.mapConn.Delete(servaddr)
}
