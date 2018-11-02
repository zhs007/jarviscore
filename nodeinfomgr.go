package jarviscore

import (
	"sync"

	"go.uber.org/zap"

	"github.com/zhs007/jarviscore/base"
	"github.com/zhs007/jarviscore/coredb/proto"
	pb "github.com/zhs007/jarviscore/proto"
)

// const coredbMyNodeInfoPrefix = "ni:"

// nodeInfoMgr -
type nodeInfoMgr struct {
	sync.RWMutex
	mapNodeInfo map[string]*NodeInfo
	node        *jarvisNode
}

type funcEachNodeInfo func(*NodeInfo)

// nodeInfoMgr -
func newNodeInfoMgr(node *jarvisNode) *nodeInfoMgr {
	return &nodeInfoMgr{
		mapNodeInfo: make(map[string]*NodeInfo),
		node:        node,
	}
}

func (mgr *nodeInfoMgr) loadFromDB() {
	jarvisbase.Debug("nodeInfoMgr.loadFromDB")

	mgr.Lock()
	defer mgr.Unlock()

	for k := range mgr.mapNodeInfo {
		delete(mgr.mapNodeInfo, k)
	}

	mgr.node.coredb.foreachNodeEx(func(key string, val *coredbpb.NodeInfo) {
		bi := BaseInfo{
			Name:     val.Name,
			ServAddr: val.ServAddr,
			Addr:     val.Addr,
			NodeType: pb.NODETYPE_NORMAL,
		}

		mgr.addNodeInfo(&bi, true)
	})

	jarvisbase.Debug("nodeInfoMgr.loadFromDB end")

	// iter := mgr.node.coredb.db.NewIteratorWithPrefix([]byte(coredbMyNodeInfoPrefix))
	// for iter.Next() {
	// 	// Remember that the contents of the returned slice should not be modified, and
	// 	// only valid until the next call to Next.
	// 	// key := iter.Key()
	// 	value := iter.Value()

	// 	ni2db := &pb.NodeInfoInDB{}
	// 	err := proto.Unmarshal(value, ni2db)
	// 	if err != nil {
	// 		bi := BaseInfo{
	// 			Name:     ni2db.NodeInfo.Name,
	// 			ServAddr: ni2db.NodeInfo.ServAddr,
	// 			Addr:     ni2db.NodeInfo.Addr,
	// 			NodeType: ni2db.NodeInfo.NodeType,
	// 		}

	// 		mgr.addNodeInfo(&bi, true)
	// 	}
	// }

	// iter.Release()
	// err := iter.Error()
	// if err != nil {

	// }
}

func (mgr *nodeInfoMgr) _saveToDB(addr string) {
	jarvisbase.Debug("nodeInfoMgr.saveToDB")

	// mgr.RLock()
	// defer mgr.RUnlock()

	cni, ok := mgr.mapNodeInfo[addr]
	if !ok {
		return
	}

	err := mgr.node.coredb.saveNode(cni)
	if err != nil {
		jarvisbase.Error("nodeInfoMgr:saveToDB:saveNode", zap.Error(err))
	}

	// ni := &pb.NodeInfo{
	// 	Name:     cni.baseinfo.Name,
	// 	ServAddr: cni.baseinfo.ServAddr,
	// 	Addr:     cni.baseinfo.Addr,
	// 	NodeType: cni.baseinfo.NodeType,
	// }

	// ni2db := &pb.NodeInfoInDB{
	// 	NodeInfo:      ni,
	// 	ConnectNums:   int32(cni.connectNums),
	// 	ConnectedNums: int32(cni.connectedNums),
	// }

	// data, err := proto.Marshal(ni2db)
	// if err != nil {
	// 	return
	// }

	// err = mgr.node.coredb.db.Put(append([]byte(coredbMyNodeInfoPrefix), addr...), data)
	// if err != nil {
	// 	return
	// }
}

func (mgr *nodeInfoMgr) addNodeInfo(bi *BaseInfo, isload bool) {
	jarvisbase.Debug("nodeInfoMgr.addNodeInfo")

	mgr.Lock()
	defer mgr.Unlock()

	if _, ok := mgr.mapNodeInfo[bi.Addr]; ok {
		return
	}

	mgr.mapNodeInfo[bi.Addr] = NewNodeInfo(bi)

	if !isload {
		mgr._saveToDB(bi.Addr)
	}
}

func (mgr *nodeInfoMgr) chg2ConnectMe(addr string) {
	jarvisbase.Debug("nodeInfoMgr.chg2ConnectMe")

	mgr.RLock()
	defer mgr.RUnlock()

	if _, ok := mgr.mapNodeInfo[addr]; !ok {
		return
	}

	mgr.mapNodeInfo[addr].connectMe = true
}

func (mgr *nodeInfoMgr) chg2ConnectNode(addr string) {
	jarvisbase.Debug("nodeInfoMgr.chg2ConnectNode")

	mgr.RLock()
	defer mgr.RUnlock()

	if _, ok := mgr.mapNodeInfo[addr]; !ok {
		return
	}

	mgr.mapNodeInfo[addr].connectNode = true
}

func (mgr *nodeInfoMgr) hasNodeInfo(addr string) bool {
	jarvisbase.Debug("nodeInfoMgr.hasNodeInfo")

	mgr.RLock()
	defer mgr.RUnlock()

	if _, ok := mgr.mapNodeInfo[addr]; ok {
		return true
	}

	return false
}

func (mgr *nodeInfoMgr) getNodeConnectState(addr string) (bool, bool) {
	jarvisbase.Debug("nodeInfoMgr.getNodeConnectState")

	mgr.RLock()
	defer mgr.RUnlock()

	if _, ok := mgr.mapNodeInfo[addr]; !ok {
		return false, false
	}

	return mgr.mapNodeInfo[addr].connectMe, mgr.mapNodeInfo[addr].connectNode
}

func (mgr *nodeInfoMgr) foreach(oneach func(*NodeInfo)) {
	jarvisbase.Debug("nodeInfoMgr.foreach")

	mgr.RLock()
	defer mgr.RUnlock()

	for _, v := range mgr.mapNodeInfo {
		oneach(v)
	}

	jarvisbase.Debug("nodeInfoMgr.foreach end")
}

// onStartConnect
func (mgr *nodeInfoMgr) onStartConnect(addr string) {
	jarvisbase.Debug("nodeInfoMgr.onStartConnect")

	mgr.RLock()
	defer mgr.RUnlock()

	if _, ok := mgr.mapNodeInfo[addr]; !ok {
		return
	}

	mgr.mapNodeInfo[addr].connectNums++

	mgr._saveToDB(addr)
}

// onConnected
func (mgr *nodeInfoMgr) onConnected(addr string) {
	jarvisbase.Debug("nodeInfoMgr.onConnected")

	mgr.RLock()
	defer mgr.RUnlock()

	if _, ok := mgr.mapNodeInfo[addr]; !ok {
		return
	}

	mgr.mapNodeInfo[addr].connectedNums++

	mgr._saveToDB(addr)
}

// getCtrlID
func (mgr *nodeInfoMgr) getCtrlID(addr string) int64 {
	jarvisbase.Debug("nodeInfoMgr.getCtrlID")

	mgr.RLock()
	defer mgr.RUnlock()

	if _, ok := mgr.mapNodeInfo[addr]; !ok {
		return -1
	}

	return mgr.mapNodeInfo[addr].ctrlid + 1
}
