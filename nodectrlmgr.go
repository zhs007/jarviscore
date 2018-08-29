package jarviscore

import "sync"

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

func (mgr *nodeCtrlMgr) clear() {
	if len(mgr.mapNodeCtrl) == 0 {
		return
	}

	for k := range mgr.mapNodeCtrl {
		delete(mgr.mapNodeCtrl, k)
	}

	mgr.mapNodeCtrl = nil
	mgr.mapNodeCtrl = make(map[string](*nodeCtrlInfo))
}

func (mgr *nodeCtrlMgr) isNeedRun(token string, ctrlid int32) bool {
	mgr.RLock()
	defer mgr.RUnlock()

	if v, ok := mgr.mapNodeCtrl[token]; ok {
		if v.hasCtrl(ctrlid) {
			return false
		}

		return true
	}

	mgr.Lock()
	defer mgr.Unlock()

	mapci, err := loadNodeCtrlInfo(getRealPath(token + ".yaml"))
	if err != nil {
		mgr.mapNodeCtrl[token] = newNodeCtrlInfo()

		return true
	}

	mgr.mapNodeCtrl[token] = mapci

	if mapci.hasCtrl(ctrlid) {
		return false
	}

	return true
}

func (mgr *nodeCtrlMgr) addCtrl(token string, ctrlid int32, command string) {
	mgr.Lock()
	defer mgr.Unlock()

	if v, ok := mgr.mapNodeCtrl[token]; ok {
		v.addCtrl(ctrlid, command)
	}
}

func (mgr *nodeCtrlMgr) setCtrlResult(token string, ctrlid int32, result string) {
	mgr.Lock()
	defer mgr.Unlock()

	if v, ok := mgr.mapNodeCtrl[token]; ok {
		v.setCtrlResult(ctrlid, result)
	}
}

func (mgr *nodeCtrlMgr) save() {
	mgr.Lock()
	defer mgr.Unlock()

	for k, v := range mgr.mapNodeCtrl {
		v.save(getRealPath(k + ".yaml"))
	}
}
