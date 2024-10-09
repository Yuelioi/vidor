package main

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/exp/rand"
)

// Task 定义一个任务
type Task struct {
	ID    string
	URL   string
	State int // 0: 空闲, 1: 下载中, 2: 完成
}

// TaskQueue 定义任务队列
type TaskQueue struct {
	workingTasks  []*Task
	queueTasks    []*Task
	finishedTasks []*Task
	limit         int
	mu            sync.Mutex
	taskChan      chan *Task
	downloadCtx   context.Context
	cancel        context.CancelFunc
	working       atomic.Bool
}

// NewTaskQueue 创建一个新的任务队列
func NewTaskQueue(limit int) *TaskQueue {
	ctx, cancel := context.WithCancel(context.Background())
	return &TaskQueue{
		workingTasks:  make([]*Task, 0),
		queueTasks:    make([]*Task, 0),
		finishedTasks: make([]*Task, 0),
		limit:         limit,
		taskChan:      make(chan *Task, limit),
		downloadCtx:   ctx,
		cancel:        cancel,
		working:       atomic.Bool{},
	}
}

func (tq *TaskQueue) StartNotify() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		tasks := slices.Concat(tq.workingTasks, tq.queueTasks, tq.finishedTasks)
		fmt.Printf("tasks: %v\n", tasks)
		if len(tq.queueTasks) == 0 && len(tq.workingTasks) == 0 {
			return
		}
	}
}

func (tq *TaskQueue) Add(task *Task) {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	tq.queueTasks = append(tq.queueTasks, task)

	fmt.Printf("tq.working.Load(): %v\n", tq.working.Load())

	if !tq.working.Load() {
		tq.working.Store(true)
		go func() {
			tq.Start()
			tq.working.Store(false)
		}()
	}
}

func (tq *TaskQueue) Start() {
	ticker := time.NewTicker(1 * time.Second)

	defer ticker.Stop()

	for range ticker.C {
		fmt.Println("检测中")
		// 尝试从通道中获取任务
		select {
		case task, ok := <-tq.taskChan:
			if !ok {
				continue
			}
			fmt.Println("获取任务成功")

			tq.mu.Lock()
			tq.workingTasks = append(tq.workingTasks, task)
			tq.mu.Unlock()

			// 在新的 goroutine 中处理下载任务
			go func(task *Task) {
				download(task)

				tq.mu.Lock()
				tq.finishedTasks = append(tq.finishedTasks, task)
				tq.workingTasks = removeTask(task, tq.workingTasks)
				tq.mu.Unlock()

			}(task)

		default:
			// 如果没有任务可用，则继续
			if len(tq.queueTasks) == 0 {
				return
			} else if len(tq.workingTasks) < tq.limit {
				tq.Reload()
			}
		}

	}
}

func download(task *Task) {
	fmt.Printf("%s 下载中\n", task.URL)
	time.Sleep(time.Second * time.Duration(rand.Intn(3)+5))
	fmt.Printf("%s 下载完成\n", task.URL)
	task.State = 2 // 设置为完成状态
}

func removeTask(task *Task, tasks []*Task) []*Task {
	for i, t := range tasks {
		if t.ID == task.ID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return tasks
		}
	}
	return tasks
}

func (tq *TaskQueue) Reload() {
	free := tq.limit - len(tq.workingTasks)
	tasksToReload := tq.queueTasks[:min(free, len(tq.queueTasks))]
	tq.queueTasks = tq.queueTasks[len(tasksToReload):]

	for _, task := range tasksToReload {
		tq.taskChan <- task
	}
}

// 主函数
func main() {
	taskQueue := NewTaskQueue(3)

	// 添加一些任务
	taskQueue.Add(&Task{ID: "task1", URL: "http://example.com/file1"})
	taskQueue.Add(&Task{ID: "task2", URL: "http://example.com/file2"})
	taskQueue.Add(&Task{ID: "task3", URL: "http://example.com/file3"})
	taskQueue.Add(&Task{ID: "task4", URL: "http://example.com/file4"})

	// 保持主程序运行
	time.Sleep(time.Second * 20)

}
