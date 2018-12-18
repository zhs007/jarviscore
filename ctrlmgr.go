package jarviscore

import (
	"sync"

	pb "github.com/zhs007/jarviscore/proto"
)

// ctrlMgr -
type ctrlMgr struct {
	sync.RWMutex
	mapCtrl map[string](Ctrl)
}

func (mgr *ctrlMgr) Reg(ctrltype string, ctrl Ctrl) {
	mgr.Lock()
	defer mgr.Unlock()

	mgr.mapCtrl[ctrltype] = ctrl
}

func (mgr *ctrlMgr) Run(ci *pb.CtrlInfo) ([]byte, error) {
	mgr.RLock()
	defer mgr.RUnlock()

	if c, ok := mgr.mapCtrl[ci.CtrlType]; ok {
		return c.Run(ci)
	}

	return nil, ErrNoCtrlCmd
}
