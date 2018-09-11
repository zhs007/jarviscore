package jarviscore

import "sync"

// nodeInfoMgr -
type nodeInfoMgr struct {
	sync.RWMutex
	mapNodeInfo map[string]*NodeInfo
}

type funcEachNodeInfo func(*NodeInfo)

// nodeInfoMgr -
func newNodeInfoMgr() *nodeInfoMgr {
	return &nodeInfoMgr{mapNodeInfo: make(map[string]*NodeInfo)}
}

func (mgr *nodeInfoMgr) addNodeInfo(bi *BaseInfo) {
	mgr.Lock()
	defer mgr.Unlock()

	if _, ok := mgr.mapNodeInfo[bi.Addr]; ok {
		return
	}

	mgr.mapNodeInfo[bi.Addr] = NewNodeInfo(bi)
}

func (mgr *nodeInfoMgr) chg2ConnectMe(token string) {
	mgr.RLock()
	defer mgr.RUnlock()

	if _, ok := mgr.mapNodeInfo[token]; !ok {
		return
	}

	mgr.mapNodeInfo[token].connectMe = true
}

func (mgr *nodeInfoMgr) chg2ConnectNode(token string) {
	mgr.RLock()
	defer mgr.RUnlock()

	if _, ok := mgr.mapNodeInfo[token]; !ok {
		return
	}

	mgr.mapNodeInfo[token].connectNode = true
}

func (mgr *nodeInfoMgr) hasNodeInfo(token string) bool {
	mgr.RLock()
	defer mgr.RUnlock()

	if _, ok := mgr.mapNodeInfo[token]; ok {
		return true
	}

	return false
}

func (mgr *nodeInfoMgr) getNodeConnectState(token string) (bool, bool) {
	mgr.RLock()
	defer mgr.RUnlock()

	if _, ok := mgr.mapNodeInfo[token]; !ok {
		return false, false
	}

	return mgr.mapNodeInfo[token].connectMe, mgr.mapNodeInfo[token].connectNode
}

func (mgr *nodeInfoMgr) foreach(oneach func(*NodeInfo)) {
	mgr.RLock()
	defer mgr.RUnlock()

	for _, v := range mgr.mapNodeInfo {
		oneach(v)
	}
}
