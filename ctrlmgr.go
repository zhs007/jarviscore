package jarviscore

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"github.com/zhs007/jarviscore/base"
	pb "github.com/zhs007/jarviscore/proto"
)

// ctrlMgr -
type ctrlMgr struct {
	mapCtrl sync.Map
}

func (mgr *ctrlMgr) Reg(ctrltype string, ctrl Ctrl) {
	mgr.mapCtrl.Store(ctrltype, ctrl)
}

func (mgr *ctrlMgr) getCtrl(ctrltype string) Ctrl {
	val, ok := mgr.mapCtrl.Load(ctrltype)
	if ok {
		ctrl, typeok := val.(Ctrl)
		if typeok {
			return ctrl
		}
	}

	return nil
}

func (mgr *ctrlMgr) Run(ctx context.Context, jarvisnode JarvisNode, srcAddr string, msgid int64, ci *pb.CtrlInfo) []*pb.JarvisMsg {
	c := mgr.getCtrl(ci.CtrlType)
	if c != nil {
		return c.Run(ctx, jarvisnode, srcAddr, msgid, ci)
	}

	msg, err := BuildReply2(jarvisnode,
		jarvisnode.GetMyInfo().Addr,
		srcAddr,
		pb.REPLYTYPE_ERROR,
		ErrNoCtrl.Error(),
		msgid)

	if err != nil {
		jarvisbase.Warn("ctrlMgr.Run:BuildReply2", zap.Error(err))

		return nil
	}

	return []*pb.JarvisMsg{msg}
}
