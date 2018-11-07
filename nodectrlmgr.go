package jarviscore

import (
	"sync"
)

// nodeCtrlMgr -
type nodeCtrlMgr struct {
	sync.RWMutex
	mapNodeCtrl map[string](*nodeCtrlInfo)
}

func newNodeCtrlMgr() *nodeCtrlMgr {
	return &nodeCtrlMgr{
		mapNodeCtrl: make(map[string](*nodeCtrlInfo)),
	}
}

// func (mgr *nodeCtrlMgr) clear() {
// 	if len(mgr.mapNodeCtrl) == 0 {
// 		return
// 	}

// 	for k := range mgr.mapNodeCtrl {
// 		delete(mgr.mapNodeCtrl, k)
// 	}

// 	mgr.mapNodeCtrl = nil
// 	mgr.mapNodeCtrl = make(map[string](*nodeCtrlInfo))
// }

func (mgr *nodeCtrlMgr) isNeedRun(addr string, ctrlid int64) bool {
	mgr.RLock()
	defer mgr.RUnlock()

	if v, ok := mgr.mapNodeCtrl[addr]; ok {
		if v.hasCtrl(ctrlid) {
			return false
		}

		return true
	}

	mgr.Lock()
	defer mgr.Unlock()

	nci := newNodeCtrlInfo(addr) //loadNodeCtrlInfo(getRealPath(token + ".yaml"))
	if nci == nil {
		// mgr.mapNodeCtrl[addr] = nci //newNodeCtrlInfo()

		return true
	}

	mgr.mapNodeCtrl[addr] = nci

	if nci.hasCtrl(ctrlid) {
		return false
	}

	return true
}

func (mgr *nodeCtrlMgr) addCtrl(addr string, ctrlid int64, ctrltype string, command []byte, forwordAddr string, forwordNums int32) {
	mgr.Lock()
	defer mgr.Unlock()

	if v, ok := mgr.mapNodeCtrl[addr]; ok {
		v.addCtrl(ctrlid, ctrltype, command, forwordAddr, forwordNums)
	}
}

func (mgr *nodeCtrlMgr) setCtrlResult(addr string, ctrlid int64, result []byte) {
	mgr.Lock()
	defer mgr.Unlock()

	if v, ok := mgr.mapNodeCtrl[addr]; ok {
		v.setCtrlResult(ctrlid, result)
	}
}

// func (mgr *nodeCtrlMgr) save() {
// 	mgr.Lock()
// 	defer mgr.Unlock()

// 	for k, v := range mgr.mapNodeCtrl {
// 		v.save(getRealPath(k + ".yaml"))
// 	}
// }
