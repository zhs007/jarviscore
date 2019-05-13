package jarviscore

import (
	"context"

	"github.com/zhs007/jarviscore/base"
	pb "github.com/zhs007/jarviscore/proto"
	"go.uber.org/zap"
)

// Ctrl -
type Ctrl interface {
	Run(ctx context.Context, jarvisnode JarvisNode, srcAddr string, msgid int64, ci *pb.CtrlInfo) []*pb.JarvisMsg
}

// BuildReply2ForCtrl - build Reply2 for Ctrl
func BuildReply2ForCtrl(jarvisnode JarvisNode, srcAddr string, msgid int64,
	replytype pb.REPLYTYPE, info string, msgs []*pb.JarvisMsg) []*pb.JarvisMsg {

	msg, err1 := BuildReply2(jarvisnode,
		jarvisnode.GetMyInfo().Addr,
		srcAddr,
		replytype,
		info,
		msgid)

	if err1 != nil {
		jarvisbase.Warn("BuildReply2ForCtrl", zap.Error(err1))

		return msgs
	}

	return append(msgs, msg)
}

// BuildCtrlResultForCtrl - build CtrlResult for Ctrl
func BuildCtrlResultForCtrl(jarvisnode JarvisNode, srcAddr string, msgid int64,
	str string, errInfo string, msgs []*pb.JarvisMsg) []*pb.JarvisMsg {

	msg, err := BuildCtrlResult(jarvisnode, jarvisnode.GetMyInfo().Addr, srcAddr, msgid, str, errInfo)
	if err != nil {
		jarvisbase.Warn("BuildCtrlResultForCtrl", zap.Error(err))

		return msgs
	}

	return append(msgs, msg)
}

// BuildReplyRequestFileForCtrl - build ReplyRequestFile for Ctrl
func BuildReplyRequestFileForCtrl(jarvisnode JarvisNode, srcAddr string, msgid int64,
	fd *pb.FileData, msgs []*pb.JarvisMsg) []*pb.JarvisMsg {

	msg, err := BuildReplyRequestFile(jarvisnode, jarvisnode.GetMyInfo().Addr, srcAddr,
		fd, msgid)
	if err != nil {
		jarvisbase.Warn("BuildReplyRequestFileForCtrl", zap.Error(err))

		return msgs
	}

	return append(msgs, msg)
}
