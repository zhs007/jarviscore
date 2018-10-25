package coredb

const keyMyPrivateData = "myprivatedata"
const prefixKeyNodeInfo = "ni:"

func makeNodeInfoKeyID(addr string) string {
	return prefixKeyNodeInfo + addr
}
