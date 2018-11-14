package jarviscore

import pb "github.com/zhs007/jarviscore/proto"

// Ctrl -
type Ctrl interface {
	Run(ci *pb.CtrlInfo) (result []byte, err error)
}
