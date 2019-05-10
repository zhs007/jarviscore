package coredb

import (
	"strconv"
	"time"

	"github.com/zhs007/jarviscore/coredb/proto"
	"github.com/zhs007/jarviscore/proto"
)

const keyMyPrivateData = "myprivatedata"
const prefixKeyNodeInfo = "ni:"
const prefixKeyNormalTask = "task:normal:"
const prefixKeyServiceTask = "task:service:"

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

func makeTaskKey(task *jarviscorepb.JarvisTask) string {
	if task.TaskType == jarviscorepb.TASKTYPE_NORMAL {
		return prefixKeyNormalTask + task.Name + ":" + strconv.FormatInt(task.CurTime, 10)
	}

	return prefixKeyServiceTask + task.Name
}
