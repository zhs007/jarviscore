package jarviscore

import "google.golang.org/grpc"

var mgrconn *connMgr

func init() {
	mgrconn = &connMgr{mapConn: make(map[string]*grpc.ClientConn)}
}
