package jarviscore

import (
	"time"

	"github.com/zhs007/jarviscore/base"
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
	// // JarvisResultTypeReplyStreamEnd - reply stream end
	// JarvisResultTypeReplyStreamEnd = 4
)

// NormalTaskInfo - normal task info
type NormalTaskInfo struct {
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
	msgs    []*pb.JarvisMsg
	procMsg pb.JarvisCoreServ_ProcMsgServer
	isSent  bool
	// procMsgStream pb.JarvisCoreServ_ProcMsgStreamServer
}

// NewJarvisMsgReplyStream - new JarvisMsgReplyStream
func NewJarvisMsgReplyStream(procMsg pb.JarvisCoreServ_ProcMsgServer) *JarvisMsgReplyStream {
	return &JarvisMsgReplyStream{
		procMsg: procMsg,
		isSent:  false,
		// procMsgStream: procMsgStream,
	}
}

// // IsValid - is valid stream
// func (stream *JarvisMsgReplyStream) IsValid() bool {
// 	return true //stream.procMsg != nil || stream.procMsgStream != nil
// }

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

	// if stream.procMsg == nil && stream.procMsgStream == nil {

	// 	jarvisbase.Warn("JarvisMsgReplyStream.ReplyMsg", zap.Error(ErrStreamNil))

	// 	return ErrStreamNil
	// }

	// sendmsg.CurTime = time.Now().Unix()

	// err := SignJarvisMsg(jn.GetCoreDB().GetPrivateKey(), sendmsg)
	// if err != nil {
	// 	jarvisbase.Warn("JarvisMsgReplyStream.ReplyMsg:SignJarvisMsg", zap.Error(err))

	// 	return err
	// }

	// if stream.procMsg != nil {
	// 	err = stream.procMsg.Send(sendmsg)
	// 	if err != nil {
	// 		jarvisbase.Warn("JarvisMsgReplyStream.ReplyMsg:procMsg.sendmsg", zap.Error(err))

	// 		return err
	// 	}
	// }

	// if stream.procMsgStream != nil {
	// 	err = stream.procMsgStream.Send(sendmsg)
	// 	if err != nil {
	// 		jarvisbase.Warn("JarvisMsgReplyStream.ReplyMsg:procMsgStream.sendmsg", zap.Error(err))

	// 		return err
	// 	}
	// }

	return nil
}
