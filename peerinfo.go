package jarviscore

// // peerAddrArr -
// type peerAddrArr struct {
// 	PeerAddr []string `yaml:"peeraddr"`
// }

// type peerInfo struct {
// 	peeraddr      string
// 	connectnums   int
// 	connectednums int
// }

// type peerInfoSlice []peerInfo

// func (s peerInfoSlice) Len() int {
// 	return len(s)
// }

// func (s peerInfoSlice) Swap(i, j int) {
// 	s[i], s[j] = s[j], s[i]
// }

// func (s peerInfoSlice) Less(i, j int) bool {
// 	return s[j].connectnums-s[j].connectednums > s[i].connectnums-s[i].connectednums
// }

// func loadPeerAddrFile(filename string) (*peerAddrArr, error) {
// 	buf, err := loadFile(filename)
// 	if err != nil {
// 		return nil, err
// 	}

// 	paa := peerAddrArr{}
// 	err1 := yaml.Unmarshal(buf, &paa)
// 	if err1 != nil {
// 		return nil, err1
// 	}

// 	return &paa, nil
// }

// func (paa *peerAddrArr) rmPeerAddrIndex(i int) {
// 	if i >= 0 && i < len(paa.PeerAddr) {
// 		narr := append(paa.PeerAddr[:i], paa.PeerAddr[i+1:]...)
// 		paa.PeerAddr = narr
// 	}
// }

// func (paa *peerAddrArr) rmPeerAddr(peeraddr string) {
// 	for i := 0; i < len(paa.PeerAddr); {
// 		if paa.PeerAddr[i] == peeraddr {
// 			paa.rmPeerAddrIndex(i)
// 		} else {
// 			i++
// 		}
// 	}
// }

// func (paa *peerAddrArr) insPeerAddr(peeraddr string) {
// 	if len(peeraddr) == 0 {
// 		return
// 	}

// 	for i := 0; i < len(paa.PeerAddr); i++ {
// 		if paa.PeerAddr[i] == peeraddr {
// 			return
// 		}
// 	}

// 	paa.PeerAddr = append(paa.PeerAddr, peeraddr)
// }
