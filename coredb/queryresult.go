package coredb

// ResultPrivateKey -
type ResultPrivateKey struct {
	PrivateKey struct {
		PriKey     string `json:"priKey"`
		PubKey     string `json:"pubKey"`
		CreateTime int64  `json:"createTime"`
		OnlineTime int64  `json:"onlineTime"`
		Addr       string `json:"addr"`
	} `json:"privateKey"`
}

// ResultPrivateData -
type ResultPrivateData struct {
	PrivateData struct {
		PriKey     string `json:"priKey"`
		PubKey     string `json:"pubKey"`
		CreateTime int64  `json:"createTime"`
		OnlineTime int64  `json:"onlineTime"`
		Addr       string `json:"addr"`
	} `json:"privateData"`
}

// ResultNodeInfo -
type ResultNodeInfo struct {
	NodeInfo struct {
		Addr          string   `json:"addr"`
		ServAddr      string   `json:"servAddr"`
		Name          string   `json:"name"`
		ConnectNums   int      `json:"connectNums"`
		ConnectedNums int      `json:"connectedNums"`
		CtrlID        int64    `json:"ctrlID"`
		LstClientAddr []string `json:"lstClientAddr"`
		AddTime       int64    `json:"addTime"`
	} `json:"nodeInfo"`
}

// ResultNodeInfos -
type ResultNodeInfos struct {
	SnapshotID int64 `json:"snapshotID"`
	EndIndex   int   `json:"endIndex"`
	MaxIndex   int   `json:"maxIndex"`

	Nodes []struct {
		Addr          string   `json:"addr"`
		ServAddr      string   `json:"servAddr"`
		Name          string   `json:"name"`
		ConnectNums   int      `json:"connectNums"`
		ConnectedNums int      `json:"connectedNums"`
		CtrlID        int64    `json:"ctrlID"`
		LstClientAddr []string `json:"lstClientAddr"`
		AddTime       int64    `json:"addTime"`
	} `json:"nodes"`
}
