package jarviscore

import (
	pb "github.com/zhs007/jarviscore/proto"
)

// BaseInfo -
type BaseInfo struct {
	Name     string
	ServAddr string
	Token    string
	NodeType pb.NODETYPE
}
