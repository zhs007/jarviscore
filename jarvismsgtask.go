package jarviscore

import pb "github.com/zhs007/jarviscore/proto"

// NormalTaskInfo - normal task info
type NormalTaskInfo struct {
	Msg      *pb.JarvisMsg
	Stream   pb.JarvisCoreServ_ProcMsgServer
	OnResult FuncOnProcMsgResult
}

// StreamTaskInfo - stream task info
type StreamTaskInfo struct {
	Msgs     []*pb.JarvisMsg
	Stream   pb.JarvisCoreServ_ProcMsgStreamServer
	OnResult FuncOnProcMsgResult
}

// JarvisTask - jarvis task
type JarvisTask struct {
	Normal *NormalTaskInfo
	Stream *StreamTaskInfo
}
