package jarviscore

import (
	"github.com/zhs007/jarviscore/base"
	pb "github.com/zhs007/jarviscore/proto"
	"go.uber.org/zap"
)

// Ctrl -
type Ctrl interface {
	Run(jarvisnode JarvisNode, srcAddr string, msgid int64, ci *pb.CtrlInfo) []*pb.JarvisMsg
}

// BuildReply2ForCtrl - build Reply2 for Ctrl
func BuildReply2ForCtrl(jarvisnode JarvisNode, srcAddr string, msgid int64,
	replytype pb.REPLYTYPE, info string) []*pb.JarvisMsg {

	msg, err1 := BuildReply2(jarvisnode,
		jarvisnode.GetMyInfo().Addr,
		srcAddr,
		replytype,
		info,
		msgid)

	if err1 != nil {
		jarvisbase.Warn("BuildReply2ForCtrl", zap.Error(err1))

		return nil
	}

	return []*pb.JarvisMsg{msg}
}

// BuildCtrlResultForCtrl - build CtrlResult for Ctrl
func BuildCtrlResultForCtrl(jarvisnode JarvisNode, srcAddr string, msgid int64, str string) []*pb.JarvisMsg {

	msg, err := BuildCtrlResult(jarvisnode, jarvisnode.GetMyInfo().Addr, srcAddr, msgid, str)
	if err != nil {
		jarvisbase.Warn("BuildCtrlResultForCtrl", zap.Error(err))

		return nil
	}

	return []*pb.JarvisMsg{msg}
}
