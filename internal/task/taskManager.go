package task

import (
	"context"
	"fmt"
	"slices"
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
}

func NewTaskManager(limit int, manager *plugin.PluginManager, ctx context.Context) *TaskManager {
	taskNotify := notify.NewTaskNotification(ctx)

	return &TaskManager{
		workingTasks:  make([]*pb.Task, 0),
		queueTasks:    make([]*pb.Task, 0),
		finishedTasks: make([]*pb.Task, 0),
		pm:            manager,
		queue:         []*TaskQueue{},
		limit:         limit,
		notification:  taskNotify,
	}
}

func (tm *TaskManager) AddTasks(tasks ...*pb.Task) {

	for _, task := range tasks {

		if len(tm.queue) < tm.limit {
			// 选择下载器
			p, err := tm.pm.Select(task.Url)
			if err != nil {
				task.Status = err.Error()
				tm.notification.UpdateTask(task)
				continue
			}
			tq := NewTaskQueue(p, func(completedTask *pb.Task, err error) {
				tm.taskCompleted(completedTask, err)
			})
			tm.queue = append(tm.queue, tq)
			tm.workingTasks = append(tm.workingTasks, task)
			tq.work(task)
		} else {
			tm.queueTasks = append(tm.queueTasks, task)
		}
	}

}

func (tm *TaskManager) SetLimit(newLimit int) {
	tm.limit = newLimit
	tm.adjustRunningTasks()
}

func (tm *TaskManager) adjustRunningTasks() {
	// 如果当前正在工作的任务数超过新的限制，则停止多余的任务
	if len(tm.workingTasks) > tm.limit {
		for i := len(tm.workingTasks) - 1; i >= tm.limit; i-- {
			// 停止最后一个任务并移至队列
			tm.queueTasks = append(tm.queueTasks, tm.workingTasks[i])
			// tm.removeWorkingTask(tm.workingTasks[i])
		}
		tm.workingTasks = tm.workingTasks[:tm.limit]
	} else if len(tm.workingTasks) < tm.limit && len(tm.queueTasks) > 0 {
		// 如果有空闲槽位且有待处理的任务，则开始新的任务
		for len(tm.workingTasks) < tm.limit && len(tm.queueTasks) > 0 {
			nextTask := tm.queueTasks[0] // 取出第一个待处理任务
			tm.queueTasks = tm.queueTasks[1:]
			tm.AddTasks(nextTask) // 重新添加任务以启动它

		}
	}
}

func (tm *TaskManager) handleError(task *pb.Task, err error) {
	// 从workingTasks中移除出错的任务
	for i, t := range tm.workingTasks {
		if t == task {
			tm.workingTasks = append(tm.workingTasks[:i], tm.workingTasks[i+1:]...)
			break
		}
	}

	// 可以选择将任务放回队列或进行其他处理
	tm.queueTasks = append(tm.queueTasks, task)

	// 发送错误通知
	task.Status = err.Error()
	tm.notification.UpdateTask(task)
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
