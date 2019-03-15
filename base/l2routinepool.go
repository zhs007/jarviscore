package jarvisbase

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// L2BaseTask - level 2 basetask
type L2BaseTask struct {
	taskID   int64
	parentID string
}

// GetParentID - get parentID
func (task *L2BaseTask) GetParentID() string {
	return task.parentID
}

// GetTaskID - get taskID
func (task *L2BaseTask) GetTaskID() int64 {
	return task.taskID
}

// Init - init
func (task *L2BaseTask) Init(pool L2RoutinePool, parentID string) {
	task.taskID = pool.NewTaskID()
	task.parentID = parentID
}

// L2Task - level 2 task
type L2Task interface {
	// Run - run task
	Run(ctx context.Context) error
	// GetParentID - get parentID
	GetParentID() string
	// GetTaskID - get taskID
	GetTaskID() int64
}

// l2routine - l2routine
type l2routine struct {
	routineID   int
	parentID    string
	chanTask    chan L2Task
	chanRemove  chan *l2routine
	chanWaiting chan *l2routine
}

// isFull - is full
func (r *l2routine) isFull() bool {
	if cap(r.chanTask) < len(r.chanTask) {
		return false
	}

	return true
}

// sendTask - start a routine
func (r *l2routine) sendTask(task L2Task) bool {
	if r.parentID != "" && task.GetParentID() != r.parentID {
		Warn("l2routine:sendTask",
			zap.String("myparentid", r.parentID),
			zap.String("taskparentid", task.GetParentID()))

		return false
	}

	if cap(r.chanTask) < len(r.chanTask) {
		return false
	}

	r.parentID = task.GetParentID()

	r.chanTask <- task

	return true
}

// start - start a routine
func (r *l2routine) start(ctx context.Context) error {
	Debug("l2routine.Start...")

	for {
		isend := false
		select {
		case task, ok := <-r.chanTask:
			if !ok {
				isend = true
				break
			}

			Debug("l2routine.Start:get new task")

			if task != nil {
				task.Run(ctx)
			}

			Debug("l2routine.Start:", zap.Int("lasttask", len(r.chanTask)))

			if len(r.chanTask) == 0 {
				r.chanWaiting <- r
			}

		case <-ctx.Done():
			Debug("l2routine.Start:context done")
			isend = true
			break
		}

		if isend {
			break
		}
	}

	r.chanRemove <- r

	return nil
}

// L2RoutinePool - L2RoutinePool
type L2RoutinePool interface {
	// SendTask - send new task
	SendTask(task L2Task)
	// Start - start a routine pool
	Start(ctx context.Context, maxNums int) error
	// GetStatus - get status
	GetStatus() string
	// NewTaskID - new taskID
	NewTaskID() int64
}

// l2routinePool - l2routinePool
type l2routinePool struct {
	mapRoutine  map[string]*l2routine
	chanRemove  chan *l2routine
	chanWaiting chan *l2routine
	chanTask    chan L2Task
	lstTask     []L2Task
	maxNums     int
	lstWaiting  []*l2routine
	lstTotal    []*l2routine
	curTaskID   int64
	zeroTaskID  int64
}

// NewL2RoutinePool - new RoutinePool
func NewL2RoutinePool() L2RoutinePool {
	return &l2routinePool{
		chanRemove:  make(chan *l2routine, 128),
		chanWaiting: make(chan *l2routine, 128),
		chanTask:    make(chan L2Task, 256),
		mapRoutine:  make(map[string]*l2routine),
	}
}

// GetStatus - get status
func (pool *l2routinePool) GetStatus() string {
	return fmt.Sprintf("l2routinePool - mapRoutine %v, chanRemove %v, chanWaiting %v, chanTask %v, lstTask %v, maxNums %v, lstWaiting %v, lstTotal %v, curTaskID %v, zeroTaskID %v",
		len(pool.mapRoutine),
		len(pool.chanRemove),
		len(pool.chanWaiting),
		len(pool.chanTask),
		len(pool.lstTask),
		pool.maxNums,
		len(pool.lstWaiting),
		len(pool.lstTotal),
		pool.curTaskID,
		pool.zeroTaskID)
}

// SendTask - send new task
func (pool *l2routinePool) SendTask(task L2Task) {
	pool.chanTask <- task
}

// Start - start a routine pool
func (pool *l2routinePool) Start(ctx context.Context, maxNums int) error {
	Debug("l2routinePool.Start...")

	pool.maxNums = maxNums

	for {
		isend := false
		select {
		case task, ok := <-pool.chanTask:
			if !ok {
				isend = true
				break
			}

			Debug("l2routinePool.Start:new task")

			pool.run(ctx, task)
		case r, ok := <-pool.chanWaiting:
			if !ok {
				isend = true
				break
			}

			Debug("l2routinePool.Start:new waiting")

			delete(pool.mapRoutine, r.parentID)
			r.parentID = ""

			pool.lstWaiting = append(pool.lstWaiting, r)

			pool.onNewWaiting(ctx)
		case r, ok := <-pool.chanRemove:
			if !ok {
				isend = true
				break
			}

			Debug("l2routinePool.Start:new remove")

			delete(pool.mapRoutine, r.parentID)
			r.parentID = ""

			pool.removeRoutine(r)
		case <-ctx.Done():
			Debug("l2routinePool.Start:context done")

			isend = true
			break
		}

		if isend {
			break
		}
	}

	return nil
}

// onNewWaiting - on new waiting
func (pool *l2routinePool) onNewWaiting(ctx context.Context) error {
	if len(pool.lstTask) <= 0 {
		return nil
	}

	task := pool.lstTask[0]

	pid := task.GetParentID()
	cr, ok := pool.mapRoutine[pid]
	if ok {
		if !cr.isFull() {
			pool.lstTask = append(pool.lstTask[:0], pool.lstTask[1:]...)

			if !cr.sendTask(task) {
				pool.lstTask = append([]L2Task{task}, pool.lstTask...)
				delete(pool.mapRoutine, pid)
			}
		}

		return nil
	}

	pool.lstTask = append(pool.lstTask[:0], pool.lstTask[1:]...)

	if len(pool.lstWaiting) > 0 {
		r := pool.lstWaiting[0]

		pool.lstWaiting = append(pool.lstWaiting[:0], pool.lstWaiting[1:]...)

		pool.mapRoutine[pid] = r

		if !r.sendTask(task) {
			pool.lstTask = append([]L2Task{task}, pool.lstTask...)
			pool.lstWaiting = append(pool.lstWaiting, r)
			delete(pool.mapRoutine, pid)
		}

		return nil
	}

	if len(pool.lstTotal) < pool.maxNums {
		r := &l2routine{
			routineID:   len(pool.lstTotal) + 1,
			chanTask:    make(chan L2Task, 8),
			chanRemove:  pool.chanRemove,
			chanWaiting: pool.chanWaiting,
		}

		pool.lstTotal = append(pool.lstTotal, r)

		pool.mapRoutine[pid] = r

		go pool.startRountine(ctx, r)

		if !r.sendTask(task) {
			pool.lstTask = append([]L2Task{task}, pool.lstTask...)
			pool.lstWaiting = append(pool.lstWaiting, r)
			delete(pool.mapRoutine, pid)
		}

		return nil
	}

	pool.lstTask = append([]L2Task{task}, pool.lstTask...)

	return nil
}

// run - run with task
func (pool *l2routinePool) run(ctx context.Context, task L2Task) error {
	pool.lstTask = append(pool.lstTask, task)

	if len(pool.lstTask) == 1 {
		pool.zeroTaskID = task.GetTaskID()

		return pool.onNewWaiting(ctx)
	}

	return nil
}

// startRountine - start a new routine
func (pool *l2routinePool) startRountine(ctx context.Context, r *l2routine) error {
	curctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return r.start(curctx)
}

// removeRoutine - remove a routine
func (pool *l2routinePool) removeRoutine(r *l2routine) {
	i := pool.findRoutine(r)
	if i >= 0 {
		pool.lstTotal = append(pool.lstTotal[:i], pool.lstTotal[i+1:]...)
	}
}

// findRoutine - find a routine
func (pool *l2routinePool) findRoutine(r *l2routine) int {
	if len(pool.lstTotal) <= 0 {
		return -1
	}

	for i := range pool.lstTotal {
		if pool.lstTotal[i].routineID == r.routineID {
			return i
		}
	}

	return -1
}

// NewTaskID - new taskID
func (pool *l2routinePool) NewTaskID() int64 {
	pool.curTaskID++
	return pool.curTaskID
}
