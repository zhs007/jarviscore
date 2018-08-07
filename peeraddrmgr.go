package jarviscore

import (
	"github.com/zhs007/jarviscore/errcode"
)

// peeraddrmgr
type peeraddrmgr struct {
	arr *peeraddrarr
	lst []peerinfo
}

func newPeerAddrMgr(peeraddrfile string, defpeeraddr string) (*peeraddrmgr, error) {
	mgr := &peeraddrmgr{}

	var err error
	mgr.arr, err = loadPeerAddrFile(peeraddrfile)
	if err != nil {
		if len(defpeeraddr) > 0 {
			warnLog("loadPeerAddrFile", err)

			mgr.arr = &peeraddrarr{}
			mgr.arr.insPeerAddr(defpeeraddr)
		} else {
			errorLog("loadPeerAddrFile", err)

			return nil, err
		}
	}

	arrlen := len(mgr.arr.PeerAddr)
	if arrlen == 0 {
		return nil, newError(jarviserrcode.PEERADDREMPTY)
	}

	mgr.lst = make([]peerinfo, arrlen)

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
