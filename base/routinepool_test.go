package jarvisbase

import (
	"context"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"
)

type taskMgr struct {
	sync.RWMutex

	finishTask int
	sendTask   int
	totalTask  int

	cancel context.CancelFunc
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

	return mgr.finishTask == mgr.sendTask && mgr.finishTask == mgr.totalTask
}

type taskTest struct {
	index int
	mgr   *taskMgr
}

func (task *taskTest) Run(ctx context.Context) error {
	// time.Sleep(1 * time.Second)
	Debug("run ", zap.Int("index", task.index))
	task.mgr.finish()

	if task.mgr.isok() {
		task.mgr.cancel()

		return nil
	}

	return nil
}

func makeTask(mgr *taskMgr, pool RoutinePool, maxtasks int) {
	time.Sleep(3 * time.Second)

	mgr.totalTask = maxtasks

	for i := 0; i < maxtasks; i++ {
		task := &taskTest{
			index: i,
			mgr:   mgr,
		}

		pool.SendTask(task)
		mgr.send()
	}

	// time.Sleep(3 * time.Second)
	// Debug("end")
	// fmt.Print("end\n")
	// cancel()
}

func TestRountinePool(t *testing.T) {
	InitLogger(zap.InfoLevel, true, "")
	Debug("start...")
	// fmt.Print("haha\n")

	pool := NewRoutinePool()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	mgr := &taskMgr{cancel: cancel}
	go makeTask(mgr, pool, 10000)

	pool.Start(ctx, 128)

	if !mgr.isok() {
		t.Fail()
	}
}
