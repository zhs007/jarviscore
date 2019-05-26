package jarviscore

import (
	"time"

	jarvisbase "github.com/zhs007/jarviscore/base"
	pb "github.com/zhs007/jarviscore/proto"
	"go.uber.org/zap"
)

const (
	// JarvisResultTypeSend - send
	JarvisResultTypeSend = 1
	// JarvisResultTypeReply - reply
	JarvisResultTypeReply = 2
	// JarvisResultTypeLocalError - error
	JarvisResultTypeLocalError = 3
	// JarvisResultTypeRemoved - removed
	JarvisResultTypeRemoved = 5
	// JarvisResultTypeLocalErrorEnd - error & end
	JarvisResultTypeLocalErrorEnd = 6
)

// NormalMsgTaskInfo - normal message task info
type NormalMsgTaskInfo struct {
	Msg         *pb.JarvisMsg
	ReplyStream *JarvisMsgReplyStream
	OnResult    FuncOnProcMsgResult
}

// JarvisMsgInfo - JarvisMsg information
type JarvisMsgInfo struct {
	JarvisResultType int           `json:"jarvisresulttype"`
	Msg              *pb.JarvisMsg `json:"msg"`
	Err              error         `json:"err"`
}

// IsEnd - is end msg
func (jmi *JarvisMsgInfo) IsEnd() bool {
	return jmi.Msg != nil && jmi.Msg.MsgType == pb.MSGTYPE_REPLY2 && jmi.Msg.ReplyType == pb.REPLYTYPE_END
}

// IsErrorEnd - is error end msg
func (jmi *JarvisMsgInfo) IsErrorEnd() bool {
	return jmi.JarvisResultType == JarvisResultTypeLocalErrorEnd
}

// IsEndOrIGI - is end msg or IGOTIT
func (jmi *JarvisMsgInfo) IsEndOrIGI() bool {
	return jmi.Msg != nil && jmi.Msg.MsgType == pb.MSGTYPE_REPLY2 &&
		(jmi.Msg.ReplyType == pb.REPLYTYPE_END || jmi.Msg.ReplyType == pb.REPLYTYPE_IGOTIT)
}

// StreamMsgTaskInfo - stream message task info
type StreamMsgTaskInfo struct {
	Msgs        []JarvisMsgInfo
	ReplyStream *JarvisMsgReplyStream
	OnResult    FuncOnProcMsgResult
}

// JarvisMsgTask - jarvis message task
type JarvisMsgTask struct {
	Normal *NormalMsgTaskInfo
	Stream *StreamMsgTaskInfo
}

// JarvisMsgReplyStream - reply JarvisMsg stream
type JarvisMsgReplyStream struct {
	msgs    []*pb.JarvisMsg
	procMsg pb.JarvisCoreServ_ProcMsgServer
	isSent  bool
}

// NewJarvisMsgReplyStream - new JarvisMsgReplyStream
func NewJarvisMsgReplyStream(procMsg pb.JarvisCoreServ_ProcMsgServer) *JarvisMsgReplyStream {
	return &JarvisMsgReplyStream{
		procMsg: procMsg,
		isSent:  false,
	}
}

// ReplyMsg - reply JarvisMsg
func (stream *JarvisMsgReplyStream) ReplyMsg(jn JarvisNode, sendmsg *pb.JarvisMsg) error {
	if sendmsg == nil {
		jarvisbase.Warn("JarvisMsgReplyStream.ReplyMsg",
			zap.Error(ErrInvalidJarvisMsgReplyStreamSendMsg))

		return ErrInvalidJarvisMsgReplyStreamSendMsg
	}

	if stream.procMsg != nil {
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

		return nil
	}

	if stream.isSent {
		jarvisbase.Warn("JarvisMsgReplyStream.ReplyMsg",
			zap.Error(ErrJarvisMsgReplyStreamSent))

		return ErrJarvisMsgReplyStreamSent
	}

	if len(stream.msgs) > 0 {
		if stream.msgs[0].DestAddr != sendmsg.DestAddr {
			jarvisbase.Warn("JarvisMsgReplyStream.ReplyMsg",
				zap.Error(ErrInvalidJarvisMsgReplyStreamDestAddr))

			return ErrInvalidJarvisMsgReplyStreamDestAddr
		}

		if stream.msgs[0].ReplyMsgID != sendmsg.ReplyMsgID {
			jarvisbase.Warn("JarvisMsgReplyStream.ReplyMsg",
				zap.Error(ErrInvalidJarvisMsgReplyStreamReplyMsgID))

			return ErrInvalidJarvisMsgReplyStreamReplyMsgID
		}
	}

	stream.msgs = append(stream.msgs, sendmsg)

	if sendmsg.MsgType == pb.MSGTYPE_REPLY2 && sendmsg.ReplyType == pb.REPLYTYPE_END {
		jn.SendStreamMsg(sendmsg.DestAddr, stream.msgs, nil)

		stream.isSent = true
	}

	return nil
}
