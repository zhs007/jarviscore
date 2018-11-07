package jarviscore

import (
	"sync"
)

// ctrlMgr -
type ctrlMgr struct {
	sync.RWMutex
	mapCtrl map[string](Ctrl)
}

func (mgr *ctrlMgr) Reg(ctrltype string, ctrl Ctrl) {
	mgr.Lock()
	defer mgr.Unlock()

	mgrCtrl.mapCtrl[ctrltype] = ctrl
}

func (mgr *ctrlMgr) Run(ctrltype string, command []byte) ([]byte, error) {
	mgr.RLock()
	defer mgr.RUnlock()

	if c, ok := mgr.mapCtrl[ctrltype]; ok {
		return c.Run(command)
	}

	return nil, ErrNoCtrlCmd
}
