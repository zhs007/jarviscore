package jarviscore

import (
	"io/ioutil"
	"path"

	"github.com/zhs007/jarviscore/errcode"

	"gopkg.in/yaml.v2"
)

// peeraddrmgr
type peeraddrmgr struct {
	arr *peeraddrarr
	lst []peerinfo
}

func newPeerAddrMgr(peeraddrfile string, defpeeraddr string) (*peeraddrmgr, error) {
	mgr := &peeraddrmgr{}

	var err error
	mgr.arr, err = loadPeerAddrFile(path.Join(config.RunPath, peeraddrfile))
	if err != nil {
		if len(defpeeraddr) > 0 {
			warnLog("loadPeerAddrFile", err)

			mgr.arr = &peeraddrarr{}
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

	mgr.lst = make([]peerinfo, 0, arrlen)

	// log.Debug("arrlen", zap.Int("len", len(mgr.lst)))

	return mgr, nil
}

func (mgr *peeraddrmgr) canConnect(peeraddr string) bool {
	for i := 0; i < len(mgr.lst); i++ {
		if peeraddr == mgr.lst[i].peeraddr {
			return false
		}
	}

	return true
}

func (mgr *peeraddrmgr) savePeerAddrFile() error {
	arr := &peeraddrarr{}

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
