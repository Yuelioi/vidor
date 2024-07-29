package test

import (
	"fmt"
	"testing"
	"time"
)

func (a *App) AddDownloadTasks(parts []Part, workName string) {
	if len(parts) == 0 {
		return
	}

	var tasks = make([]*Task, 0)

	for _, part := range parts {
		if taskExists(a.tasks, part.Url) {
			continue
		} else {

			task := &Task{
				part:  part,
				state: 0,
			}

			tasks = append(tasks, task)
			a.tasks = append(a.tasks, task)
		}
	}
	// 添加到队列
	if a.taskQueue == nil || a.taskQueue.state == Finished {
		println("任务队列 重新创建")
		a.taskQueue = NewTaskQueue(tasks, 1)
	} else {
		println("任务队列 还在使用")
		a.taskQueue.AddTasks(tasks)
	}

}

func (a *App) RemoveTask(uid string) bool {

	for i, task := range a.tasks {
		if task.part.UID == uid {
			if task.state == Working {
				if task.state == Queue {
					// 1. 正在队列中
					a.Logger.Info("任务移除(队列中):", task.part.UID)
					a.taskQueue.removeQueueTasks([]*Task{task})
				} else {
					// 2. 正在下载中
					a.Logger.Info("任务移除(下载中):", task.part.UID)
				}
			}

			// 3. 直接删除
			a.tasks = append(a.tasks[:i], a.tasks[i+1:]...)

			return true
		}
	}
	return false
}

func (a *App) RemoveAllTask(parts []Part) bool {

	newTasks := make([]*Task, 0)
	delQueueTasks := make([]*Task, 0)

	partUIDs := make(map[string]Part)

	for _, part := range parts {
		partUIDs[part.UID] = part
	}

	for i, task := range a.tasks {
		if _, found := partUIDs[task.part.UID]; found {
			if task.state == Working {
				if task.state == Queue {
					// 添加到待删队列
					delQueueTasks = append(delQueueTasks, task)
				} else {
					// 直接调用停止函数
					a.Logger.Info("任务移除(下载中):", task.part.UID)
				}
			}
		} else {
			newTasks = append(newTasks, a.tasks[i])
		}
	}

	// 移除队列中任务
	if a.taskQueue != nil && a.taskQueue.state != Finished {
		a.taskQueue.removeQueueTasks(delQueueTasks)
	}

	// 保存任务清单
	a.tasks = newTasks

	return true
}

func Test_AddDownloadTasks(t *testing.T) {
	a := App{}

	parts := []Part{
		{UID: "1", Url: "1"},
		{UID: "2", Url: "2"},
		{UID: "3", Url: "3"},
	}
	parts2 := []Part{
		{UID: "4", Url: "4"},
		{UID: "5", Url: "5"},
		{UID: "6", Url: "6"},
	}
	a.AddDownloadTasks(parts, "workName")

	// time.Sleep(time.Second * 10)
	a.AddDownloadTasks(parts2, "workName")

	time.Sleep(time.Second * 20)
}
func Test_RemoveTask(t *testing.T) {
	a := App{}

	parts := []Part{
		{UID: "1", Url: "1"},
		{UID: "2", Url: "2"},
		{UID: "3", Url: "3"},
	}
	parts2 := []Part{
		{UID: "4", Url: "4"},
		{UID: "5", Url: "5"},
		{UID: "6", Url: "6"},
	}
	a.AddDownloadTasks(parts, "workName")
	a.AddDownloadTasks(parts2, "workName")
	time.Sleep(time.Second * 3)

	a.taskQueue.removeQueueTasks(a.tasks)

	time.Sleep(time.Second * 10)
}

func Test_RemoveAllTask(t *testing.T) {

	a := App{}

	parts := []Part{
		{UID: "1", Url: "1"},
		{UID: "2", Url: "2"},
		{UID: "3", Url: "3"},
	}
	parts2 := []Part{
		{UID: "4", Url: "4"},
		{UID: "5", Url: "5"},
		{UID: "6", Url: "6"},
	}
	a.AddDownloadTasks(parts, "workName")
	a.AddDownloadTasks(parts2, "workName")
	fmt.Println("开始移除队列任务")

	queueTasks := make([]*Task, 0)
	for _, task := range a.tasks {
		println(task.state)
		if task.state == Queue {
			queueTasks = append(queueTasks, task)
		}
	}
	a.taskQueue.removeQueueTasks(queueTasks)

	time.Sleep(time.Second * 20)
}

func taskExists(tasks []*Task, url string) bool {
	for _, existingTask := range tasks {
		if existingTask.part.Url == url {
			return true
		}
	}
	return false
}
