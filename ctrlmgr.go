package jarviscore

import (
	"sync"

	pb "github.com/zhs007/jarviscore/proto"
)

// ctrlMgr -
type ctrlMgr struct {
	sync.RWMutex
	mapCtrl map[pb.CTRLTYPE](Ctrl)
}

func (mgr *ctrlMgr) Reg(ctrltype pb.CTRLTYPE, ctrl Ctrl) {
	mgr.Lock()
	defer mgr.Unlock()

	mgrCtrl.mapCtrl[ctrltype] = ctrl
}

func (mgr *ctrlMgr) Run(ctrltype pb.CTRLTYPE, command []byte) ([]byte, error) {
	mgr.RLock()
	defer mgr.RUnlock()

	if c, ok := mgr.mapCtrl[ctrltype]; ok {
		return c.Run(command)
	}

	return nil, ErrNoCtrlCmd
}
