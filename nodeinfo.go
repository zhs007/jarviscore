package jarviscore

// NodeInfo -
type NodeInfo struct {
	baseinfo    BaseInfo
	mapclient   map[string]BaseInfo
	connectMe   bool
	connectNode bool
}

// NewNodeInfo -
func NewNodeInfo(bi *BaseInfo) *NodeInfo {
	return &NodeInfo{
		baseinfo:  *bi,
		mapclient: make(map[string]BaseInfo),
	}
}

func (ni *NodeInfo) onNodeConnectMe() {
	ni.connectMe = true
}

func (ni *NodeInfo) onIConnectNode() {
	ni.connectNode = true
}
