package jarviscore

import (
	"fmt"
	"sync"
	"time"

	"github.com/zhs007/jarviscore/basedef"
)

// FuncGetMsgState - getMsgState
type FuncGetMsgState func(addr string, msgid int64) int

type myReplyInfo struct {
	addr            string
	msgid           int64
	startTime       int64
	endTime         int64
	funcGetMsgState FuncGetMsgState
}

// wait4MyReplyMgr - Waiting for my reply to the msg queue
type wait4MyReplyMgr struct {
	mapWait4MyReply sync.Map
	mapEnd          sync.Map
}

// newWait4MyReplyMgr - new wait4MyReplyMgr
func newWait4MyReplyMgr() *wait4MyReplyMgr {
	mgr := &wait4MyReplyMgr{}

	return mgr
}

// addMsgInfo
func (mgr *wait4MyReplyMgr) addMsgInfo(addr string, msgid int64, funcGetMsgState FuncGetMsgState) error {
	if !IsValidNodeAddr(addr) {
		return ErrInvalidWait4MyReplyAddr
	}

	if msgid <= 0 {
		return ErrInvalidWait4MyReplyMsgID
	}

	myri := &myReplyInfo{
		addr:            addr,
		msgid:           msgid,
		startTime:       time.Now().Unix(),
		funcGetMsgState: funcGetMsgState,
	}

	k := fmt.Sprintf("%v:%v", addr, msgid)

	mgr.mapWait4MyReply.Store(k, myri)

	return nil
}

// getMsgState
func (mgr *wait4MyReplyMgr) getMsgState(addr string, msgid int64) int {
	k := fmt.Sprintf("%v:%v", addr, msgid)

	v, isok := mgr.mapWait4MyReply.Load(k)
	if !isok {
		ev, isok := mgr.mapEnd.Load(k)
		if !isok {
			return -1
		}

		myri, istypeok := ev.(*myReplyInfo)
		if istypeok {
			if myri.funcGetMsgState != nil {
				return myri.funcGetMsgState(addr, msgid)
			}

			return 1
		}

		return -1
	}

	myri, istypeok := v.(*myReplyInfo)
	if istypeok {
		if myri.funcGetMsgState != nil {
			return myri.funcGetMsgState(addr, msgid)
		}

		return 1
	}

	mgr.mapWait4MyReply.Delete(k)

	return -1
}

// delete
func (mgr *wait4MyReplyMgr) delete(addr string, msgid int64) {
	k := fmt.Sprintf("%v:%v", addr, msgid)

	v, isok := mgr.mapWait4MyReply.Load(k)
	if isok {
		myri, istypeok := v.(*myReplyInfo)
		if istypeok {
			myri.endTime = time.Now().Unix()

			mgr.mapEnd.Store(k, myri)
		}

		mgr.mapWait4MyReply.Delete(k)
	}
}

// clearEndCache
func (mgr *wait4MyReplyMgr) clearEndCache() {
	ct := time.Now().Unix()

	mgr.mapEnd.Range(func(k, v interface{}) bool {
		myri, ok := v.(*myReplyInfo)
		if ok && ct > myri.endTime+basedef.TimeClearEndMsgState {
			mgr.mapEnd.Delete(k)
		}

		return true
	})
}
