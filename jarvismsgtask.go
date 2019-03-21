package jarviscore

import (
	"time"

	"github.com/zhs007/jarviscore/base"
	pb "github.com/zhs007/jarviscore/proto"
	"go.uber.org/zap"
)

// NormalTaskInfo - normal task info
type NormalTaskInfo struct {
	Msg         *pb.JarvisMsg
	ReplyStream *JarvisMsgReplyStream
	OnResult    FuncOnProcMsgResult
}

// JarvisMsgInfo - JarvisMsg information
type JarvisMsgInfo struct {
	Msg *pb.JarvisMsg `json:"msg"`
	Err error         `json:"err"`
}

// StreamTaskInfo - stream task info
type StreamTaskInfo struct {
	Msgs        []JarvisMsgInfo
	ReplyStream *JarvisMsgReplyStream
	OnResult    FuncOnProcMsgResult
}

// JarvisTask - jarvis task
type JarvisTask struct {
	Normal *NormalTaskInfo
	Stream *StreamTaskInfo
}

// JarvisMsgReplyStream - reply JarvisMsg stream
type JarvisMsgReplyStream struct {
	procMsg       pb.JarvisCoreServ_ProcMsgServer
	procMsgStream pb.JarvisCoreServ_ProcMsgStreamServer
}

// NewJarvisMsgReplyStream - new JarvisMsgReplyStream
func NewJarvisMsgReplyStream(procMsg pb.JarvisCoreServ_ProcMsgServer, procMsgStream pb.JarvisCoreServ_ProcMsgStreamServer) *JarvisMsgReplyStream {
	return &JarvisMsgReplyStream{
		procMsg:       procMsg,
		procMsgStream: procMsgStream,
	}
}

// IsValid - is valid stream
func (stream *JarvisMsgReplyStream) IsValid() bool {
	return stream.procMsg != nil || stream.procMsgStream != nil
}

// ReplyMsg - reply JarvisMsg
func (stream *JarvisMsgReplyStream) ReplyMsg(jn JarvisNode, sendmsg *pb.JarvisMsg) error {
	if stream.procMsg == nil && stream.procMsgStream == nil {

		jarvisbase.Warn("JarvisMsgReplyStream.ReplyMsg", zap.Error(ErrStreamNil))

		return ErrStreamNil
	}

	sendmsg.CurTime = time.Now().Unix()

	err := SignJarvisMsg(jn.GetCoreDB().GetPrivateKey(), sendmsg)
	if err != nil {
		jarvisbase.Warn("JarvisMsgReplyStream.ReplyMsg:SignJarvisMsg", zap.Error(err))

		return err
	}

	if stream.procMsg != nil {
		err = stream.procMsg.Send(sendmsg)
		if err != nil {
			jarvisbase.Warn("JarvisMsgReplyStream.ReplyMsg:procMsg.sendmsg", zap.Error(err))

			return err
		}
	}

	if stream.procMsgStream != nil {
		err = stream.procMsgStream.Send(sendmsg)
		if err != nil {
			jarvisbase.Warn("JarvisMsgReplyStream.ReplyMsg:procMsgStream.sendmsg", zap.Error(err))

			return err
		}
	}

	return nil
}
