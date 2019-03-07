package jarvisbase

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"
)

type l2taskMgr struct {
	sync.RWMutex

	finishTask int
	sendTask   int
	totalTask  int
	mapIndex   map[int]int

	cancel context.CancelFunc
}

func (mgr *l2taskMgr) init(pid int) {
	mgr.Lock()
	defer mgr.Unlock()

	mgr.mapIndex[pid] = 0
}

func (mgr *l2taskMgr) send() {
	mgr.Lock()
	defer mgr.Unlock()

	mgr.sendTask++
}

func (mgr *l2taskMgr) finish(pid int) {
	mgr.Lock()
	defer mgr.Unlock()

	mgr.mapIndex[pid]++
	mgr.finishTask++
}

func (mgr *l2taskMgr) check(pid int, index int) bool {
	mgr.RLock()
	defer mgr.RUnlock()

	if mgr.mapIndex[pid] == index {
		return true
	}

	Debug("l2taskMgr.check",
		zap.Int("pid", pid),
		zap.Int("mgr.mapIndex", mgr.mapIndex[pid]),
		zap.Int("index", index))

	return false
}

func (mgr *l2taskMgr) isok() bool {
	mgr.RLock()
	defer mgr.RUnlock()

	Debug("l2taskMgr:isok",
		zap.Int("total", mgr.totalTask),
		zap.Int("send", mgr.sendTask),
		zap.Int("finish", mgr.finishTask))

	return mgr.finishTask == mgr.sendTask && mgr.finishTask == mgr.totalTask
}

type l2taskTest struct {
	index int
	j     int
	mgr   *l2taskMgr
	pid   string
}

func (task *l2taskTest) Run(ctx context.Context) error {
	// time.Sleep(1 * time.Second)

	if !task.mgr.check(task.j, task.index) {
		task.mgr.cancel()

		return nil
	}

	Debug("run ", zap.Int("index", task.index), zap.Int("j", task.j))
	task.mgr.finish(task.j)

	if task.mgr.isok() {
		task.mgr.cancel()

		return nil
	}

	return nil
}

// GetParentID - get parentID
func (task *l2taskTest) GetParentID() string {
	return task.pid
}

func l2makeTask(mgr *l2taskMgr, pool L2RoutinePool, maxpid int, maxc int) {
	time.Sleep(3 * time.Second)

	mgr.totalTask = maxpid * maxc

	for i := 0; i < maxc; i++ {
		for j := 0; j < maxpid; j++ {
			if i == 0 {
				mgr.init(j)
			}

			task := &l2taskTest{
				index: i,
				j:     j,
				mgr:   mgr,
				pid:   fmt.Sprintf("pid%v", j),
			}

			pool.SendTask(task)
			mgr.send()
		}
	}

	// time.Sleep(3 * time.Second)
	// Debug("end")
	// fmt.Print("end\n")
	// cancel()
}

func TestL2RountinePool(t *testing.T) {
	InitLogger(zap.InfoLevel, true, "")
	Debug("start...")
	// fmt.Print("haha\n")

	pool := NewL2RoutinePool()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	mgr := &l2taskMgr{
		cancel:   cancel,
		mapIndex: make(map[int]int),
	}
	go l2makeTask(mgr, pool, 1000, 102)

	pool.Start(ctx, 128)

	if !mgr.isok() {
		t.Fail()
	}
}
