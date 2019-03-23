package jarviscore

import (
	"context"
	"strconv"
	"sync"

	"github.com/zhs007/jarviscore/base"
	"go.uber.org/zap"

	"github.com/zhs007/jarviscore/proto"
)

// procMsgResultMgr - procMsg result manager
type procMsgResultMgr struct {
	node        JarvisNode
	mapWaitPush sync.Map
}

// newProcMsgResultMgr - new procMsgResultMgr
func newProcMsgResultMgr(node JarvisNode) *procMsgResultMgr {
	mgr := &procMsgResultMgr{
		node: node,
	}

	return mgr
}

// onProcMsg
func (mgr *procMsgResultMgr) onProcMsg(ctx context.Context, taskinfo *JarvisTask) error {
	if taskinfo.Normal != nil {
		if taskinfo.Normal.Msg.MsgType == jarviscorepb.MSGTYPE_REPLY2 &&
			taskinfo.Normal.Msg.ReplyType == jarviscorepb.REPLYTYPE_END {

			mgr.endProcMsgResultData(taskinfo.Normal.Msg.SrcAddr,
				taskinfo.Normal.Msg.ReplyMsgID)
		}
	} else if taskinfo.Stream != nil {
		for i := 0; i < len(taskinfo.Stream.Msgs); i++ {
			if taskinfo.Stream.Msgs[i].Msg.MsgType == jarviscorepb.MSGTYPE_REPLY2 &&
				taskinfo.Stream.Msgs[i].Msg.ReplyType == jarviscorepb.REPLYTYPE_END {

				mgr.endProcMsgResultData(taskinfo.Stream.Msgs[i].Msg.SrcAddr,
					taskinfo.Stream.Msgs[i].Msg.ReplyMsgID)
			}
		}
	}

	return nil
}

// startProcMsgResultData
func (mgr *procMsgResultMgr) startProcMsgResultData(addr string, msgid int64, onProcMsgResult FuncOnProcMsgResult) error {

	d, _ := mgr.getProcMsgResultData(addr, msgid)
	if d != nil {
		return ErrDuplicateProcMsgResultData
	}

	pmrd := NewProcMsgResultData(onProcMsgResult)

	mgr.mapWaitPush.Store(AppendString(addr, ":", strconv.FormatInt(msgid, 10)), pmrd)

	return nil
}

// getProcMsgResultData
func (mgr *procMsgResultMgr) getProcMsgResultData(addr string, msgid int64) (*ProcMsgResultData, error) {

	v, ok := mgr.mapWaitPush.Load(AppendString(addr, ":", strconv.FormatInt(msgid, 10)))
	if ok {
		d, typeok := v.(*ProcMsgResultData)
		if typeok {
			return d, nil
		}

		return nil, ErrInvalidProcMsgResultData
	}

	return nil, ErrNoProcMsgResultData
}

// endProcMsgResultData
func (mgr *procMsgResultMgr) endProcMsgResultData(addr string, msgid int64) {
	mgr.mapWaitPush.Delete(AppendString(addr, ":", strconv.FormatInt(msgid, 10)))
}

func (mgr *procMsgResultMgr) onPorcMsgResult(ctx context.Context, addr string, replymsgid int64,
	jarvisnode JarvisNode, result *JarvisMsgInfo) error {

	d, err := mgr.getProcMsgResultData(addr, replymsgid)
	if err != nil {
		jarvisbase.Warn("procMsgResultMgr.onPorcMsgResult:getProcMsgResultData",
			zap.Error(err))
	}

	if d != nil {
		return d.OnPorcMsgResult(ctx, jarvisnode, result)
	}

	return nil
}
