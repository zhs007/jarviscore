package coredb

import (
	"time"

	"github.com/zhs007/jarviscore/coredb/proto"
)

const keyMyPrivateData = "myprivatedata"
const prefixKeyNodeInfo = "ni:"

func makeNodeInfoKeyID(addr string) string {
	return prefixKeyNodeInfo + addr
}

// IsDeprecatedNode - is a deprecated node?
func IsDeprecatedNode(ni *coredbpb.NodeInfo) bool {
	if ni.Deprecated {
		return true
	}

	if ni.TimestampDeprecated > 0 {
		ct := time.Now().Unix()
		if ct <= ni.TimestampDeprecated {
			return true
		}
	}

	return false
}
