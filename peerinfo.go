package jarviscore

import (
	"gopkg.in/yaml.v2"
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

func loadPeerAddrFile(filename string) (*peeraddrarr, error) {
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

func (paa *peeraddrarr) insPeerAddr(peeraddr string) {
	for i := 0; i < len(paa.PeerAddr); i++ {
		if paa.PeerAddr[i] == peeraddr {
			return
		}
	}

	paa.PeerAddr = append(paa.PeerAddr, peeraddr)
}
