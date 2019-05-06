package jarviscore

import (
	"context"
	"strconv"
	"sync"

	"github.com/zhs007/jarviscore/base"
	"go.uber.org/zap"

	"github.com/zhs007/jarviscore/proto"
)

// FuncOnRangeProcMsgResult - onRangeProcMsgResult
type FuncOnRangeProcMsgResult func(prmd *ProcMsgResultData)

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
			taskinfo.Normal.Msg.ReplyType == jarviscorepb.REPLYTYPE_END &&
			taskinfo.Normal.Msg.ReplyMsgID > 0 {

			mgr.onEndMsg(taskinfo.Normal.Msg.SrcAddr,
				taskinfo.Normal.Msg.ReplyMsgID)
		}
	} else if taskinfo.Stream != nil {
		for i := 0; i < len(taskinfo.Stream.Msgs); i++ {
			if taskinfo.Stream.Msgs[i].Msg != nil && taskinfo.Stream.Msgs[i].Msg.MsgType == jarviscorepb.MSGTYPE_REPLY2 &&
				taskinfo.Stream.Msgs[i].Msg.ReplyType == jarviscorepb.REPLYTYPE_END &&
				taskinfo.Stream.Msgs[i].Msg.ReplyMsgID > 0 {

				mgr.onEndMsg(taskinfo.Stream.Msgs[i].Msg.SrcAddr,
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

	pmrd := NewProcMsgResultData(addr, msgid, onProcMsgResult)

	// jarvisbase.Info("procMsgResultMgr.startProcMsgResultData:Store",
	// 	zap.String("key", AppendString(addr, ":", strconv.FormatInt(msgid, 10))),
	// 	zap.Int("nums", mgr.countNums()))

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

// onEndMsg
func (mgr *procMsgResultMgr) onEndMsg(addr string, replymsgid int64) {
	d, err := mgr.getProcMsgResultData(addr, replymsgid)
	if err != nil {
		jarvisbase.Warn("procMsgResultMgr.onEndMsg:getProcMsgResultData",
			zap.Error(err),
			zap.Int64("replymsgid", replymsgid))
	}

	if d != nil {
		if d.OnMsgEnd() {
			mgr.mapWaitPush.Delete(AppendString(addr, ":", strconv.FormatInt(replymsgid, 10)))

			jarvisbase.Info("procMsgResultMgr.onEndMsg:Delete",
				zap.String("key", AppendString(addr, ":", strconv.FormatInt(replymsgid, 10))),
				zap.Int("nums", mgr.countNums()))
		}
	}
}

func (mgr *procMsgResultMgr) onPorcMsgResult(ctx context.Context, addr string, replymsgid int64,
	jarvisnode JarvisNode, result *JarvisMsgInfo) error {

	d, err := mgr.getProcMsgResultData(addr, replymsgid)
	if err != nil {
		jarvisbase.Warn("procMsgResultMgr.onPorcMsgResult:getProcMsgResultData",
			zap.Error(err))
	}

	if d != nil {
		err := d.OnPorcMsgResult(ctx, jarvisnode, result)

		if result.IsEnd() {
			// if result.JarvisResultType == JarvisResultTypeReplyStreamEnd {
			// if result.Err == nil && result.Msg == nil {
			if d.OnRecvEnd() {
				jarvisbase.Info("procMsgResultMgr.onPorcMsgResult:Delete",
					zap.String("key", AppendString(addr, ":", strconv.FormatInt(replymsgid, 10))),
					zap.Int("nums", mgr.countNums()))

				mgr.mapWaitPush.Delete(AppendString(addr, ":", strconv.FormatInt(replymsgid, 10)))
			}
		}

		return err
	}

	return nil
}

func (mgr *procMsgResultMgr) countNums() int {
	nums := 0

	mgr.mapWaitPush.Range(func(key interface{}, value interface{}) bool {
		nums++

		return true
	})

	return nums
}

func (mgr *procMsgResultMgr) delete(addr string, replymsgid int64) {
	jarvisbase.Info("procMsgResultMgr.delete",
		zap.String("key", AppendString(addr, ":", strconv.FormatInt(replymsgid, 10))),
		zap.Int("nums", mgr.countNums()))

	mgr.mapWaitPush.Delete(AppendString(addr, ":", strconv.FormatInt(replymsgid, 10)))
}

func (mgr *procMsgResultMgr) forEach(onRange FuncOnRangeProcMsgResult) {
	mgr.mapWaitPush.Range(func(key interface{}, value interface{}) bool {
		k, iskeyok := key.(string)
		if !iskeyok {
			jarvisbase.Info("procMsgResultMgr.delete:validKey",
				zap.Int("nums", mgr.countNums()))

			mgr.mapWaitPush.Delete(key)

			return true
		}

		prmd, isok := value.(*ProcMsgResultData)
		if isok {
			onRange(prmd)
		} else {
			jarvisbase.Info("procMsgResultMgr.delete",
				zap.String("key", string(k)),
				zap.Int("nums", mgr.countNums()))

			mgr.mapWaitPush.Delete(key)
		}

		return true
	})
}
