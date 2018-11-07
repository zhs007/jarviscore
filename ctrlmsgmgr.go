package jarviscore

import (
	"context"

	pb "github.com/zhs007/jarviscore/proto"
)

// ctrlMsg - ctrl msg
type ctrlMsg struct {
	ctrlID     int64
	srcAddr    string
	ctrl       *pb.CtrlInfo
	ctrlResult *pb.CtrlResult
}

// ctrlMsgMgr - ctrl msg mgr
type ctrlMsgMgr struct {
	mapCtrlMsg map[int64]ctrlMsg
	chanMsg    chan ctrlMsg
}

// newCtrlMsgMgr - new ctrlMsgMgr
func newCtrlMsgMgr() *ctrlMsgMgr {
	mgr := &ctrlMsgMgr{
		mapCtrlMsg: make(map[int64]ctrlMsg),
		chanMsg:    make(chan ctrlMsg, 256),
	}

	return mgr
}

// sendMsg - send a ctrl msg
func (mgr *ctrlMsgMgr) sendMsg(msg ctrlMsg) {
	mgr.chanMsg <- msg
}

// start - start goroutine to proc ctrl msg
func (mgr *ctrlMsgMgr) start(ctx context.Context) error {
	for {
		select {
		case msg, ok := <-mgr.chanMsg:
			if !ok {
				return nil
			}

			mymsg, mapok := mgr.mapCtrlMsg[msg.ctrlID]
			if mapok {
				if msg.ctrlResult != nil {
					mymsg.ctrlResult = msg.ctrlResult
				}
			} else {
				mgr.mapCtrlMsg[msg.ctrlID] = msg
			}
		case <-ctx.Done():
			return nil
		}
	}
}
