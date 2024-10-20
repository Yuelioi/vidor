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
	queues        []*TaskQueue
	workingTasks  []*pb.Task
	queueTasks    []*pb.Task
	finishedTasks []*pb.Task
	notifyCancel  context.CancelFunc
	notification  *notify.TaskNotification
	taskQueuePool *sync.Pool
}

func NewTaskManager(limit int, manager *plugin.PluginManager, ctx context.Context) *TaskManager {
	taskNotify := notify.NewTaskNotification(ctx)

	return &TaskManager{
		pm:            manager,
		limit:         limit,
		queues:        []*TaskQueue{},
		workingTasks:  make([]*pb.Task, 0),
		queueTasks:    make([]*pb.Task, 0),
		finishedTasks: make([]*pb.Task, 0),
		notification:  taskNotify,
		taskQueuePool: &sync.Pool{},
	}
}

// 最初添加任务
func (tm *TaskManager) AddTasks(tasks ...*pb.Task) {
	for _, task := range tasks {
		if len(tm.queues) < tm.limit {
			// 开始任务
			tm.startNotify()
			tm.processTask(task)
		} else {
			// 加入候补
			tm.queueTasks = append(tm.queueTasks, task)
		}
	}

}

// 添加任务后 处理任务
func (tm *TaskManager) processTask(task *pb.Task) error {
	tq, err := tm.getTaskQueue(task)
	if tq == nil {
		return err
	}

	defer tm.putTaskQueue(tq)

	// 注册任务队列
	tm.queues = append(tm.queues, tq)
	// 更新工作任务
	tm.workingTasks = append(tm.workingTasks, task)
	// 开始工作
	tq.work(func(completedTask *pb.Task, err error) {
		tm.taskCompleted(completedTask, err)
	})

	return nil
}

// 获取/新建任务队列
func (tm *TaskManager) getTaskQueue(task *pb.Task) (*TaskQueue, error) {
	tq := tm.taskQueuePool.Get()

	if tq == nil {
		// 获取下载器
		p, err := tm.pm.SelectDownloader(task.Url)
		if err != nil {
			task.Status = err.Error()
			tm.notification.UpdateTask(task)
			return nil, err
		}

		tq = NewTaskQueue(p, task)
		return tq.(*TaskQueue), nil
	}
	// 初始化tq
	return tq.(*TaskQueue), nil
}

// 存任务队列
func (tm *TaskManager) putTaskQueue(tq *TaskQueue) {
	tm.taskQueuePool.Put(tq)
}

func (tm *TaskManager) SetLimit(newLimit int) {
	tm.limit = newLimit
}

// 任务善后
func (tm *TaskManager) taskCompleted(task *pb.Task, err error) error {
	// 更新任务状态
	if err != nil {
		task.Status = err.Error()
	}

	tm.workingTasks = removeTask(task, tm.workingTasks)
	tm.finishedTasks = append(tm.finishedTasks, task)

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

func (tm *TaskManager) startNextQueuedTask() {
	if len(tm.queueTasks) > 0 && len(tm.workingTasks) < tm.limit {
		nextTask := tm.queueTasks[0]
		tm.queueTasks = tm.queueTasks[1:]
		tm.AddTasks(nextTask)
		return
	}

	if len(tm.queueTasks) == 0 && len(tm.workingTasks) == 0 {
		tm.stopNotify()
	}
}

func (tm *TaskManager) Notify(ctx context.Context) {
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
	case <-ctx.Done():
		println("已关闭通知")
		break
	}
}

func (tm *TaskManager) startNotify() {
	if tm.notifyCancel == nil {
		ctx, cancel := context.WithCancel(context.Background())
		tm.notifyCancel = cancel
		go tm.Notify(ctx)
	}

}
func (tm *TaskManager) stopNotify() {
	if tm.notifyCancel != nil {
		tm.notifyCancel()
		tm.notifyCancel = nil
	}
}
