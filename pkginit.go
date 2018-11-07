package jarviscore

import (
	"google.golang.org/grpc"
)

var mgrconn *connMgr
var mgrCtrl *ctrlMgr

func init() {
	mgrconn = &connMgr{mapConn: make(map[string]*grpc.ClientConn)}

	mgrCtrl = &ctrlMgr{mapCtrl: make(map[string](Ctrl))}
	mgrCtrl.Reg("shell", &CtrlShell{})
}
