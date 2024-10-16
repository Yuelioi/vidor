package task

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/Yuelioi/vidor/internal/notify"
	"github.com/Yuelioi/vidor/internal/plugin"
	pb "github.com/Yuelioi/vidor/internal/proto"
)

// TaskQueue 任务队列接口
type TaskManager struct {
	pm            *plugin.PluginManager
	limit         int
	queue         []*TaskQueue
	workingTasks  []*pb.Task
	queueTasks    []*pb.Task
	finishedTasks []*pb.Task
	notifyChan    chan struct{}
	notification  *notify.TaskNotification
	taskQueuePool *sync.Pool
}

func NewTaskManager(limit int, manager *plugin.PluginManager, ctx context.Context) *TaskManager {
	taskNotify := notify.NewTaskNotification(ctx)

	return &TaskManager{
		pm:            manager,
		limit:         limit,
		queue:         []*TaskQueue{},
		workingTasks:  make([]*pb.Task, 0),
		queueTasks:    make([]*pb.Task, 0),
		finishedTasks: make([]*pb.Task, 0),
		notifyChan:    make(chan struct{}),
		notification:  taskNotify,
		taskQueuePool: &sync.Pool{},
	}
}

func (tm *TaskManager) AddTasks(tasks ...*pb.Task) {
	for _, task := range tasks {
		if len(tm.queue) < tm.limit {

			tm.processTask(task)

		} else {
			tm.queueTasks = append(tm.queueTasks, task)
		}
	}

}

func (tm *TaskManager) processTask(task *pb.Task) error {
	tq := tm.getTaskQueue(task)
	defer tm.putTaskQueue(tq)

	// 	tq := NewTaskQueue(p, func(completedTask *pb.Task, err error) {
	// 	tm.taskCompleted(completedTask, err)
	// })
	tm.queue = append(tm.queue, tq)
	tm.workingTasks = append(tm.workingTasks, task)
	tq.work()

	// 处理任务...

	return nil
}

// 获取/新建任务队列, 可能为nil
func (tm *TaskManager) getTaskQueue(task *pb.Task) *TaskQueue {
	tq := tm.taskQueuePool.Get().(*TaskQueue)
	if tq == nil {
		p, err := tm.pm.Select(task.Url)
		if err != nil {
			task.Status = err.Error()
			tm.notification.UpdateTask(task)
			return nil
		}

		tq = NewTaskQueue(p, task, func(completedTask *pb.Task, err error) {
			tm.taskCompleted(completedTask, err)
		})
	}
	// 初始化tq
	return tq
}

func (tm *TaskManager) putTaskQueue(tq *TaskQueue) {
	tm.taskQueuePool.Put(tq)
}

func (tm *TaskManager) SetLimit(newLimit int) {
	tm.limit = newLimit
}

func (tm *TaskManager) taskCompleted(task *pb.Task, err error) error {
	// 更新状态
	tm.moveTaskToFinished(task, err == nil)

	// 处理下一个队列中的任务
	tm.startNextQueuedTask()
	return nil
}

func removeTask(task *pb.Task, tasks []*pb.Task) []*pb.Task {
	for i, t := range tasks {
		if t.Id == task.Id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return tasks
		}
	}
	return tasks
}

func (tm *TaskManager) moveTaskToFinished(task *pb.Task, removeFromWorking bool) {
	if removeFromWorking {
		tm.workingTasks = removeTask(task, tm.workingTasks)
	}
	tm.finishedTasks = append(tm.finishedTasks, task)
}

func (tm *TaskManager) startNextQueuedTask() {
	if len(tm.queueTasks) > 0 && len(tm.workingTasks) < tm.limit {
		nextTask := tm.queueTasks[0]
		tm.queueTasks = tm.queueTasks[1:]
		tm.AddTasks(nextTask)
	} else {
		close(tm.notifyChan)

	}
}

func (tm *TaskManager) StartNotify() {
	fmt.Println("开始下载通知")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	select {
	case <-ticker.C:
		for _, task := range tm.workingTasks {
			task.State = 1
		}
		for _, task := range tm.queueTasks {
			task.State = 2
		}
		for _, task := range tm.finishedTasks {
			task.State = 3
		}

		tasks := slices.Concat(tm.workingTasks, tm.queueTasks, tm.finishedTasks)
		tm.notification.UpdateTasks(tasks)
	case <-tm.notifyChan:
		break
	}
}
