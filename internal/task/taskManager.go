package task

import (
	"fmt"
	"sync"
)

// TaskQueue 任务队列接口
type TaskQueue interface {
	AddTask(id string)
	RemoveTask(id string)
	List()
}

// Task 任务结构体
type Task struct {
	id string
}

// taskQueue 任务队列实现
type taskQueue struct {
	tasks map[string]Task
	mu    sync.RWMutex
}

// New 创建一个新的任务队列
func New() TaskQueue {
	return &taskQueue{
		tasks: make(map[string]Task),
	}
}

// AddTask 添加任务到队列
func (tq *taskQueue) AddTask(id string) {
	tq.mu.Lock()
	defer tq.mu.Unlock()
	tq.tasks[id] = Task{id: id}
}

// RemoveTask 从队列中移除任务
func (tq *taskQueue) RemoveTask(id string) {
	tq.mu.Lock()
	defer tq.mu.Unlock()
	delete(tq.tasks, id)
}

// List 列出所有任务
func (tq *taskQueue) List() {
	tq.mu.RLock()
	defer tq.mu.RUnlock()
	for _, task := range tq.tasks {
		fmt.Println(task.id)
	}
}

func main() {
	tq := New()
	tq.AddTask("task1")
	tq.AddTask("task2")
	tq.List() // 输出：task1 task2
	tq.RemoveTask("task1")
	tq.List() // 输出：task2
}
