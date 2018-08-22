package jarviscore

import (
	"io/ioutil"
	"path"
	"sort"

	"github.com/zhs007/jarviscore/errcode"

	"gopkg.in/yaml.v2"
)

// peerAddrMgr
type peerAddrMgr struct {
	// load & save
	arr *peerAddrArr
	lst peerInfoSlice
}

func newPeerAddrMgr(peeraddrfile string, defpeeraddr string) (*peerAddrMgr, error) {
	mgr := &peerAddrMgr{}

	var err error
	mgr.arr, err = loadPeerAddrFile(path.Join(config.RunPath, peeraddrfile))
	if err != nil {
		if len(defpeeraddr) > 0 {
			warnLog("loadPeerAddrFile", err)

			mgr.arr = &peerAddrArr{}
			// log.Debug("arrlen", zap.Int("len", len(mgr.arr.PeerAddr)))

			mgr.arr.insPeerAddr(defpeeraddr)

			// log.Debug("arrlen", zap.Int("len", len(mgr.arr.PeerAddr)))
		} else {
			errorLog("loadPeerAddrFile", err)

			return nil, err
		}
	}

	arrlen := len(mgr.arr.PeerAddr)
	if arrlen == 0 {
		return nil, newError(jarviserrcode.PEERADDREMPTY)
	}

	mgr.lst = make([]peerInfo, 0, arrlen)

	// log.Debug("arrlen", zap.Int("len", len(mgr.lst)))

	return mgr, nil
}

func (mgr *peerAddrMgr) canConnect(peeraddr string) bool {
	for i := 0; i < len(mgr.lst); i++ {
		if peeraddr == mgr.lst[i].peeraddr {
			return false
		}
	}

	return true
}

func (mgr *peerAddrMgr) savePeerAddrFile() error {
	arr := &peerAddrArr{}

	sort.Sort(peerInfoSlice(mgr.lst))

	for _, v := range mgr.lst {
		arr.insPeerAddr(v.peeraddr)
	}

	for _, v := range mgr.arr.PeerAddr {
		arr.insPeerAddr(v)
	}

	d, err := yaml.Marshal(arr)
	if err != nil {
		return err
	}

	ioutil.WriteFile(path.Join(config.RunPath, "peeraddr.yaml"), d, 0755)

	return nil
}

// onStartConnect
func (mgr *peerAddrMgr) onStartConnect(peeraddr string) {
	for i := 0; i < len(mgr.lst); i++ {
		if peeraddr == mgr.lst[i].peeraddr {
			mgr.lst[i].connectnums++

			return
		}
	}

	mgr.lst = append(mgr.lst, peerInfo{peeraddr: peeraddr, connectnums: 1})

	return
}

// onConnected
func (mgr *peerAddrMgr) onConnected(peeraddr string) {
	for i := 0; i < len(mgr.lst); i++ {
		if peeraddr == mgr.lst[i].peeraddr {
			mgr.lst[i].connectednums++

			return
		}
	}
}
