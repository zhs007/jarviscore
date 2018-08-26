package jarviscore

// NodeInfo -
type NodeInfo struct {
	baseinfo  BaseInfo
	mapclient map[string]BaseInfo
}

// NewNodeInfo -
func NewNodeInfo(bi *BaseInfo) *NodeInfo {
	return &NodeInfo{baseinfo: *bi, mapclient: make(map[string]BaseInfo)}
}
