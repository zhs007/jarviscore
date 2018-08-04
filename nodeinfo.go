package jarviscore

import (
	pb "github.com/zhs007/jarviscore/proto"
)

type NodeInfo struct {
	Name     string
	ServAddr string
	Token    string
	NodeType pb.NODETYPE
}
