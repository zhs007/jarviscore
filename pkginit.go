package jarviscore

import (
	pb "github.com/zhs007/jarviscore/proto"
	"google.golang.org/grpc"
)

var mgrconn *connMgr
var mgrCtrl *ctrlMgr

func init() {
	mgrconn = &connMgr{mapConn: make(map[string]*grpc.ClientConn)}
	mgrCtrl = &ctrlMgr{mapCtrl: make(map[pb.CTRLTYPE](*Ctrl))}
}
