package jarviscore

// NodeInfo -
type NodeInfo struct {
	baseinfo  BaseInfo
	mapclient map[string]BaseInfo
}

// NewNodeInfo -
func NewNodeInfo() *NodeInfo {
	return &NodeInfo{mapclient: make(map[string]BaseInfo)}
}
