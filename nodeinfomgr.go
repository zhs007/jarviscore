package jarviscore

import "sync"

// nodeInfoMgr -
type nodeInfoMgr struct {
	sync.RWMutex
	mapNodeInfo map[string]NodeInfo
}

type funcEachNodeInfo func(*NodeInfo)

// nodeInfoMgr -
func newNodeInfoMgr() *nodeInfoMgr {
	return &nodeInfoMgr{mapNodeInfo: make(map[string]NodeInfo)}
}

func (mgr *nodeInfoMgr) addNodeInfo(bi *BaseInfo) {
	mgr.Lock()
	defer mgr.Unlock()

	mgr.mapNodeInfo[bi.Token] = NewNodeInfo(bi)
}

func (mgr *nodeInfoMgr) hasNodeInfo(token string) bool {
	mgr.RLock()
	defer mgr.RUnlock()

	if _, ok := mgr.mapNodeInfo[token]; ok {
		return true
	}

	return false
}

func (mgr *nodeInfoMgr) foreach(oneach func(*NodeInfo)) {
	mgr.RLock()
	defer mgr.RUnlock()

	for _, v := range mgr.mapNodeInfo {
		oneach(&v)
	}
}
