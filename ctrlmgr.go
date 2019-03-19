package jarviscore

import (
	"sync"

	pb "github.com/zhs007/jarviscore/proto"
)

// ctrlMgr -
type ctrlMgr struct {
	mapCtrl sync.Map
}

func (mgr *ctrlMgr) Reg(ctrltype string, ctrl Ctrl) {
	mgr.mapCtrl.Store(ctrltype, ctrl)
}

func (mgr *ctrlMgr) getCtrl(ctrltype string) Ctrl {
	val, ok := mgr.mapCtrl.Load(ctrltype)
	if ok {
		ctrl, typeok := val.(Ctrl)
		if typeok {
			return ctrl
		}
	}

	return nil
}

func (mgr *ctrlMgr) Run(ci *pb.CtrlInfo) ([]byte, error) {
	c := mgr.getCtrl(ci.CtrlType)
	if c != nil {
		return c.Run(ci)
	}

	return nil, ErrNoCtrlCmd
}
