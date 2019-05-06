package jarviscore

import (
	"context"
	// pb "github.com/zhs007/jarviscore/proto"
)

// FuncOnProcMsgResult - on procmsg recv the message
type FuncOnProcMsgResult func(ctx context.Context, jarvisnode JarvisNode,
	lstResult []*JarvisMsgInfo) error

// IsClientProcMsgResultEnd - is end
func IsClientProcMsgResultEnd(lstResult []*JarvisMsgInfo) bool {
	if len(lstResult) > 0 {
		hasend := false

		for _, v := range lstResult {
			if v.Msg != nil && v.IsEnd() {
				hasend = true

				break
			}
		}

		if hasend && lstResult[len(lstResult)-1].IsEndOrIGI() {
			return true
		}
	}

	return false
	// return len(lstResult) > 0 && lstResult[len(lstResult)-1].JarvisResultType == JarvisResultTypeReplyStreamEnd
}

// ProcMsgResultData -
type ProcMsgResultData struct {
	addr            string
	msgid           int64
	onProcMsgResult FuncOnProcMsgResult
	lstResult       []*JarvisMsgInfo
	endOnMsg        bool
	endRecv         bool
}

// NewProcMsgResultData - new ProcMsgResultData
func NewProcMsgResultData(addr string, msgid int64, onProcMsgResult FuncOnProcMsgResult) *ProcMsgResultData {
	return &ProcMsgResultData{
		addr:            addr,
		msgid:           msgid,
		onProcMsgResult: onProcMsgResult,
	}
}

// OnRecvEnd - on RecvEnd
func (pmrd *ProcMsgResultData) OnRecvEnd() bool {
	pmrd.endRecv = true

	return pmrd.endOnMsg && pmrd.endRecv
}

// OnMsgEnd - on MsgEnd
func (pmrd *ProcMsgResultData) OnMsgEnd() bool {
	pmrd.endOnMsg = true

	return pmrd.endOnMsg && pmrd.endRecv
}

// OnPorcMsgResult - on PorcMsgResult
func (pmrd *ProcMsgResultData) OnPorcMsgResult(ctx context.Context, jarvisnode JarvisNode, result *JarvisMsgInfo) error {
	pmrd.lstResult = append(pmrd.lstResult, result)

	return pmrd.onProcMsgResult(ctx, jarvisnode, pmrd.lstResult)
}

// ClientGroupProcMsgResults - result for FuncOnSendMsgResult
type ClientGroupProcMsgResults struct {
	Results []*JarvisMsgInfo `json:"results"`
}

// FuncOnGroupSendMsgResult - on group sendmsg recv the messages
type FuncOnGroupSendMsgResult func(ctx context.Context, jarvisnode JarvisNode,
	numsNode int, lstResult []*ClientGroupProcMsgResults) error

// CountClientGroupProcMsgResultsEnd - count number for end
func CountClientGroupProcMsgResultsEnd(lstResult []*ClientGroupProcMsgResults) int {
	nums := 0
	for i := 0; i < len(lstResult); i++ {
		cr := lstResult[i]
		if len(cr.Results) > 0 && cr.Results[len(cr.Results)-1].Msg == nil && cr.Results[len(cr.Results)-1].Err == nil {
			nums++
		}
	}

	return nums
}
