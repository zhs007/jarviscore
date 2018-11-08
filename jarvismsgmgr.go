package jarviscore

import (
	"context"

	pb "github.com/zhs007/jarviscore/proto"
)

type jarvisMsgTask struct {
	msg *pb.JarvisMsg
	mgr *jarvisMsgMgr
}

func (task *jarvisMsgTask) Run(ctx context.Context) error {
	return task.mgr.node.OnMsg(ctx, task.msg)
}

// jarvisMsgMgr - jarvis msg mgr
type jarvisMsgMgr struct {
	pool RoutinePool
	node JarvisNode
}

// newJarvisMsgMgr - new jarvisMsgMgr
func newJarvisMsgMgr(node JarvisNode) *jarvisMsgMgr {
	mgr := &jarvisMsgMgr{
		pool: NewRoutinePool(),
		node: node,
	}

	return mgr
}

// sendMsg - send a ctrl msg
func (mgr *jarvisMsgMgr) sendMsg(msg *pb.JarvisMsg) {
	task := &jarvisMsgTask{
		msg: msg,
	}

	mgr.pool.SendTask(task)
}

// start - start goroutine to proc ctrl msg
func (mgr *jarvisMsgMgr) start(ctx context.Context) error {
	return mgr.pool.Start(ctx, 128)
}
