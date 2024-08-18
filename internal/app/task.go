package app

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

// 保存单个任务
func saveTask(srcTask *Task, tasks []*Task, configDir string) error {
	// parts := make([]shared.Part, 0)

	// // 修改/更新
	// for _, task := range tasks {
	// 	if srcTask.part.TaskID == task.part.TaskID {
	// 		parts = append(parts, *srcTask.part)
	// 	} else {
	// 		parts = append(parts, *task.part)
	// 	}
	// }

	// tasksData, err := json.MarshalIndent(parts, "", "  ")
	// if err != nil {
	// 	logger.Error("" + err.Error())
	// }

	// err = os.WriteFile(filepath.Join(configDir, "tasks.json"), tasksData, 0644)
	// if err != nil {
	// 	logger.Error(err)
	// 	return err
	// }
	return nil
}

// task更新时 保存
func saveTasks(tasks []*Task, configDir string) error {
	// parts := make([]shared.Part, 0)

	// for _, task := range tasks {
	// 	part := *task.part
	// 	parts = append(parts, part)
	// }

	// tasksData, err := json.MarshalIndent(parts, "", "  ")
	// if err != nil {
	// 	logger.Error(err)
	// 	return err
	// }

	// err = os.WriteFile(filepath.Join(configDir, "tasks.json"), tasksData, 0644)
	// if err != nil {
	// 	logger.Error(err)
	// 	return err
	// }
	return nil
}

// func (a *App) loadTasks() error {
// 	tasks := make([]*Task, 0)

// 	configFile := filepath.Join(a.appInfo.configDir, "tasks.json")
// 	configData, err := os.ReadFile(configFile)
// 	if err != nil {
// 		logger.Errorf("Cannot read/find task file: %v", err)
// 		a.tasks = tasks
// 		return err
// 	}

// 	parts := make([]shared.Part, 0)
// 	err = json.Unmarshal(configData, &parts)
// 	if err != nil {
// 		logger.Errorf("Cannot convert task data: %v", err)
// 		a.tasks = tasks
// 		return err
// 	}

// 	for _, part := range parts {
// 		// 过滤掉不存在的任务
// 		if _, err = os.Stat(part.DownloadDir); err == nil {
// 			newPart := part
// 			tasks = append(tasks, &Task{
// 				part: &newPart,
// 			})
// 		}
// 	}
// 	a.tasks = tasks
// 	saveTasks(tasks, a.configDir)
// 	return nil
// }
