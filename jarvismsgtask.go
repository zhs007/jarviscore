package jarviscore

import pb "github.com/zhs007/jarviscore/proto"

// NormalTaskInfo - normal task info
type NormalTaskInfo struct {
	Msg      *pb.JarvisMsg
	Stream   pb.JarvisCoreServ_ProcMsgServer
	OnResult FuncOnProcMsgResult
}

// JarvisMsgInfo - JarvisMsg information
type JarvisMsgInfo struct {
	Msg *pb.JarvisMsg `json:"msg"`
	Err error         `json:"err"`
}

// StreamTaskInfo - stream task info
type StreamTaskInfo struct {
	Msgs     []JarvisMsgInfo
	Stream   pb.JarvisCoreServ_ProcMsgStreamServer
	OnResult FuncOnProcMsgResult
}

// JarvisTask - jarvis task
type JarvisTask struct {
	Normal *NormalTaskInfo
	Stream *StreamTaskInfo
}
