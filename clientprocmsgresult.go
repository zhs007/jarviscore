package jarviscore

import (
	"context"
)

// FuncOnProcMsgResult - on procmsg recv the message
type FuncOnProcMsgResult func(ctx context.Context, jarvisnode JarvisNode,
	lstResult []*JarvisMsgInfo) error

// IsClientProcMsgResultEnd - is end
func IsClientProcMsgResultEnd(lstResult []*JarvisMsgInfo) bool {
	return len(lstResult) > 0 && lstResult[len(lstResult)-1].Msg == nil
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
		if len(cr.Results) > 0 && cr.Results[len(cr.Results)-1].Msg == nil {
			nums++
		}
	}

	return nums
}
