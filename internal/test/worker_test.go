package test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/Yuelioi/vidor/utils"
)

const (
	Queue = iota
	Working
	Finished
)

/*
创建新的任务队列
并添加任务
*/
func NewTaskQueue(tasks []*Task, DownloadLimit int) *TaskQueue {
	limit := utils.ClampInt(DownloadLimit, 1, 10)

	tq := &TaskQueue{
		state:      Queue,
		tasks:      make(chan *Task, limit),
		queueTasks: make([]*Task, 0),
		done:       make(chan struct{}),
	}

	tq.AddTasks(tasks)

	for i := 0; i < limit; i++ {
		go tq.worker()
	}

	go func() {
		tq.wg.Wait()
	}()

	return tq

}

// 任务工作函数
// 影响: tq.queueTasks tq.tasks

func (tq *TaskQueue) worker() {
	for {

		select {
		case task, ok := <-tq.tasks:
			if !ok {
				// 如果tq.tasks关闭，直接退出
				tq.Close()
				return
			}
			// 处理任务
			tq.state = Working
			task.state = Working
			tq.handleTask(task)

		case <-tq.done:
			tq.Close()
			return

		default:
			if len(tq.queueTasks) > 0 {
				tq.reFillTasks()
			} else {
				tq.Close()
				return
			}
			time.Sleep(1 * time.Second)
		}
	}
}

/*
添加任务
 1. 添加到任务通道
 2. 超额的添加到tq.queueTasks

影响: tq.queueTasks
*/
func (tq *TaskQueue) AddTasks(tasks []*Task) {
	tq.queueTasks = append(tq.queueTasks, tasks...)
}

// 重新填充任务
//
// 影响: tq.queueTasks tq.tasks 已加锁
func (tq *TaskQueue) reFillTasks() {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	availableCapacity := cap(tq.tasks) - len(tq.tasks)
	slice := utils.MinInt(availableCapacity, len(tq.queueTasks))

	for _, task := range tq.queueTasks {
		if availableCapacity <= 0 {
			break
		}

		task.state = Working
		tq.tasks <- task
		// 添加到任务通道 wg+1
		tq.tasksRemaining.Add(1)
		tq.wg.Add(1)
		availableCapacity--
	}
	tq.queueTasks = tq.queueTasks[slice:]
}

// 移除队列中的任务
// 影响: tq.queueTasks 已加锁
func (tq *TaskQueue) removeQueueTasks(tasks []*Task) {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	tasksSet := make(map[string]struct{})
	for _, task := range tasks {
		tasksSet[task.part.TaskID] = struct{}{}
	}

	fmt.Printf("tasks: %v\n", tasks)

	newTasks := make([]*Task, 0)

	// 基于索引重新规划 queueTasks
	for _, task := range tq.queueTasks {
		if _, ok := tasksSet[task.part.TaskID]; ok {
			continue
		}
		newTasks = append(newTasks, task)
	}

	tq.queueTasks = newTasks
}

/*
处理任务

	1.正常下载流程
	2.完成需要检测任务队列
*/
func (tq *TaskQueue) handleTask(task *Task) {

	defer func() {
		// if err := recover(); err != nil {
		// 	tq.app.Logger.Errorf("Task handling panic: %v", err)
		// }
	}()

	defer func() {
		tq.tasksRemaining.Add(-1)
		task.state = Finished
		// 下载完成 检测任务队列
		if len(tq.queueTasks) > 0 {
			fmt.Printf("任务%s 完成: 准备重新填充\n", task.part.TaskID)
			tq.reFillTasks()
			tq.state = Finished
		} else if tq.tasksRemaining.Load() != 0 {
			fmt.Printf("任务%s 完成, 等待后续任务下载完毕: \n", task.part.TaskID)
			tq.state = Finished

		} else {
			fmt.Printf("任务%s 完成: 关闭下载队列\n\n", task.part.TaskID)
			tq.state = Finished
			close(tq.done)
		}
		tq.wg.Done()
	}()
	task.state = Working
	fmt.Printf("任务%s working\n", task.part.TaskID)
	time.Sleep(time.Second * time.Duration(rand.Intn(3)))

}

func (tq *TaskQueue) Close() {
	fmt.Println("子协程已退出")
	tq.state = Finished
}

func TestNewTaskQueue(t *testing.T) {
	tasks := []*Task{
		{part: Part{TaskID: "1"}},
		{part: Part{TaskID: "2"}},
		{part: Part{TaskID: "3"}},
		// {part: Part{TaskID: "4"}},
		// {part: Part{TaskID: "5"}},
	}
	NewTaskQueue(tasks, 1)
	time.Sleep(time.Second * 20)
}

func TestTaskQueue_AddTasks(t *testing.T) {
	tasks1 := []*Task{
		{part: Part{TaskID: "1"}},
		{part: Part{TaskID: "2"}},
	}
	tq := NewTaskQueue(tasks1, 2)

	tasks2 := []*Task{
		// {part: Part{TaskID: "1"}},
		// {part: Part{TaskID: "2"}},
		{part: Part{TaskID: "3"}},
		{part: Part{TaskID: "4"}},
	}
	tq.AddTasks(tasks2)
	time.Sleep(time.Second * 20)

}

func TestTaskQueue_removeQueueTasks(t *testing.T) {
	tasks := []*Task{
		{part: Part{TaskID: "1"}},
		{part: Part{TaskID: "2"}},
		{part: Part{TaskID: "3"}},
		{part: Part{TaskID: "4"}},
		{part: Part{TaskID: "5"}},
	}
	tq := NewTaskQueue(tasks, 2)

	tq.removeQueueTasks([]*Task{tasks[3], tasks[2]})

	time.Sleep(time.Second * 10)

}

func TestTaskQueue_Close(t *testing.T) {
	tasks := []*Task{
		{part: Part{TaskID: "1"}},
		{part: Part{TaskID: "2"}},
		{part: Part{TaskID: "3"}},
		{part: Part{TaskID: "4"}},
		{part: Part{TaskID: "5"}},
	}

	tq := NewTaskQueue(tasks, 2)
	time.Sleep(time.Second * 2)
	tq.Close()
	time.Sleep(time.Second * 10)
}
