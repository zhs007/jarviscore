package jarviscore

import (
	"context"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/zhs007/jarviscore/base"
)

type taskMgr struct {
	sync.RWMutex

	finishTask int
	sendTask   int
}

func (mgr *taskMgr) send() {
	mgr.Lock()
	defer mgr.Unlock()

	mgr.sendTask++
}

func (mgr *taskMgr) finish() {
	mgr.Lock()
	defer mgr.Unlock()

	mgr.finishTask++
}

func (mgr *taskMgr) isok() bool {
	mgr.Lock()
	defer mgr.Unlock()

	return mgr.finishTask == mgr.sendTask
}

type taskTest struct {
	index int
	mgr   *taskMgr
}

func (task *taskTest) Run(ctx context.Context) error {
	// time.Sleep(1 * time.Second)
	jarvisbase.Debug("run ", zap.Int("index", task.index))
	task.mgr.finish()

	return nil
}

func makeTask(mgr *taskMgr, pool RoutinePool, cancel context.CancelFunc) {
	time.Sleep(3 * time.Second)

	for i := 0; i < 10000; i++ {
		task := &taskTest{
			index: i,
			mgr:   mgr,
		}

		pool.SendTask(task)
		mgr.send()
	}

	time.Sleep(3 * time.Second)
	jarvisbase.Debug("end")
	// fmt.Print("end\n")
	cancel()
}

func Test_RountinePool(t *testing.T) {
	jarvisbase.InitLogger(zap.DebugLevel, true, "")
	jarvisbase.Debug("start...")
	// fmt.Print("haha\n")

	pool := NewRoutinePool()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mgr := &taskMgr{}
	go makeTask(mgr, pool, cancel)

	pool.Start(ctx, 128)

	if !mgr.isok() {
		t.Fail()
	}
}
