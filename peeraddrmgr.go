package jarviscore

import (
	"gopkg.in/yaml.v2"

	"github.com/zhs007/jarviscore/errcode"
)

// peeraddrarr -
type peeraddrarr struct {
	PeerAddr []string `yaml:"peeraddr"`
}

type peerinfo struct {
	peeraddr      string
	connectnums   int
	connectednums int
}

// peeraddrmgr
type peeraddrmgr struct {
	arr *peeraddrarr
	lst []peerinfo
}

func loadPeerAddr(filename string) (*peeraddrarr, error) {
	buf, err := loadFile(filename)
	if err != nil {
		return nil, err
	}

	paa := peeraddrarr{}
	err1 := yaml.Unmarshal(buf, &paa)
	if err1 != nil {
		return nil, err1
	}

	return &paa, nil
}

func newPeerAddrMgr(filename string) (*peeraddrmgr, error) {
	mgr := &peeraddrmgr{}

	var err error
	mgr.arr, err = loadPeerAddr(filename)
	if err != nil {
		return nil, err
	}

	arrlen := len(mgr.arr.PeerAddr)
	if arrlen == 0 {
		return nil, newError(jarviserrcode.PEERADDREMPTY)
	}

	mgr.lst = make([]peerinfo, arrlen)

	return mgr, nil
}

func (paa *peeraddrarr) rmPeerAddrIndex(i int) {
	if i >= 0 && i < len(paa.PeerAddr) {
		narr := append(paa.PeerAddr[:i], paa.PeerAddr[i+1:]...)
		paa.PeerAddr = narr
	}
}

func (paa *peeraddrarr) rmPeerAddr(peeraddr string) {
	for i := 0; i < len(paa.PeerAddr); {
		if paa.PeerAddr[i] == peeraddr {
			paa.rmPeerAddrIndex(i)
		} else {
			i++
		}
	}
}

func (mgr *peeraddrmgr) canConnect(peeraddr string) bool {
	for i := 0; i < len(mgr.lst); i++ {
		if peeraddr == mgr.lst[i].peeraddr {
			return false
		}
	}

	return true
}
