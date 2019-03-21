package jarviscore

import (
	"context"

	"github.com/zhs007/jarviscore/base"
)

type jarvisMsgTask struct {
	mgr      *jarvisMsgMgr
	chanEnd  chan int
	taskinfo JarvisTask
}

func (task *jarvisMsgTask) Run(ctx context.Context) error {
	err := task.mgr.node.OnMsg(ctx, task.taskinfo)

	if task.chanEnd != nil {
		task.chanEnd <- 0
	}

	return err
}

// jarvisMsgMgr - jarvis msg mgr
type jarvisMsgMgr struct {
	pool jarvisbase.RoutinePool
	node JarvisNode
}

// newJarvisMsgMgr - new jarvisMsgMgr
func newJarvisMsgMgr(node JarvisNode) *jarvisMsgMgr {
	mgr := &jarvisMsgMgr{
		pool: jarvisbase.NewRoutinePool(),
		node: node,
	}

	return mgr
}

// sendMsg - send a ctrl msg
func (mgr *jarvisMsgMgr) sendMsg(normal *NormalTaskInfo, chanEnd chan int) {

	task := &jarvisMsgTask{
		taskinfo: JarvisTask{
			Normal: normal,
		},
		mgr:     mgr,
		chanEnd: chanEnd,
	}

	mgr.pool.SendTask(task)
}

// sendStreamMsg - send stream msg
func (mgr *jarvisMsgMgr) sendStreamMsg(stream *StreamTaskInfo, chanEnd chan int) {

	task := &jarvisMsgTask{
		taskinfo: JarvisTask{
			Stream: stream,
		},
		mgr:     mgr,
		chanEnd: chanEnd,
	}

	mgr.pool.SendTask(task)
}

// start - start goroutine to proc ctrl msg
func (mgr *jarvisMsgMgr) start(ctx context.Context) error {
	return mgr.pool.Start(ctx, 1)
}
