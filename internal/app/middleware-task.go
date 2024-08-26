package app

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "github.com/Yuelioi/vidor/internal/proto"
	"github.com/Yuelioi/vidor/internal/task"
)

/*
	添加下载任务

1. 获取任务目标
2. 创建/添加到任务队列
3. 保存任务信息
*/
func (app *App) AddDownloadTasks(taskMaps []taskMap) bool {

	// 获取任务
	tasks := []*pb.Task{}
	for _, taskMap := range taskMaps {
		cacheTask, ok := app.cache.Task(taskMap.id)
		if !ok {
			continue
		}

		// 将 formatIds 转换为集合，便于快速查找
		formatSet := make(map[string]struct{})
		for _, formatId := range taskMap.formatIds {
			formatSet[formatId] = struct{}{}
		}

		// 过滤掉不符合条件的 formats
		filterSegments(cacheTask.Segments, formatSet)

		tasks = append(tasks, cacheTask)
	}

	// 获取缓存下载器
	plugin := app.cache.Downloader()

	// 清除任务缓存
	app.cache.ClearTasks()

	stream, err := plugin.Service.Download(context.Background(), &pb.DownloadRequest{
		Tasks: tasks,
	})

	if err != nil {
		return false
	}

	for {
		progress, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error receiving progress: %v", err)
		}
		fmt.Printf("Download progress: %s - %s\n", progress.Id, progress.Speed)
	}

	return true

	// return tasksToParts(tasks)
}

/*
移除单个任务
 1. 下载中
    调用download.Stop, handleTask续会自动检测tq.queueTasks是否为空, 需要refill等
 2. 队列中
    直接删除对应任务
 3. 已完成
    直接删除对应任务
*/
func (app *App) RemoveTask(uid string) bool {

	// for i, task := range app.tasks {
	// 	if task.part.TaskID == uid {

	// 		if checkTaskQueueWorking(a) {
	// 			if task.state == Queue {
	// 				// 1. 正在队列中
	// 				app.logger.Info("任务移除(队列中):", task.part.Title)
	// 				app.taskQueue.removeQueueTasks([]*Task{task})
	// 			} else {
	// 				// 2. 正在下载中
	// 				app.logger.Info("任务移除(下载中):", task.part.Title)
	// 				app.taskQueue.stopTask(task)
	// 			}
	// 		}

	// 		// 3. 直接删除
	// 		app.tasks = append(app.tasks[:i], app.tasks[i+1:]...)

	// 		if err := saveTasks(app.tasks, app.configDir); err != nil {
	// 			app.logger.Warn(err)
	// 		}
	// 		return true
	// 	}
	// }
	return false
}

// 移除任务
// 移除完成任务: 去除app.tasks目标 并保存配置
// 移除下载中任务: 调用下载器StopDownload函数 关闭stopChan
// 移除队列中任务: 清理缓存队列的queueTasks
func (app *App) RemoveAllTask(parts []Part) bool {

	// newTasks := make([]*Task, 0)
	// delQueueTasks := make([]*Task, 0)

	// partTaskIDs := make(map[string]Part)

	// for _, part := range parts {
	// 	partTaskIDs[part.TaskID] = part
	// }

	// for i, task := range app.tasks {
	// 	if _, found := partTaskIDs[task.part.TaskID]; found {

	// 		if checkTaskQueueWorking(a) {
	// 			if task.state == Queue {
	// 				// 添加到待删队列
	// 				delQueueTasks = append(delQueueTasks, task)
	// 			} else {
	// 				// 直接调用停止函数
	// 				app.logger.Info("任务移除(下载中):", task.part.Title)
	// 				task.downloader.Cancel(context.Background(), task.part)
	// 			}
	// 		}

	// 	} else {
	// 		newTasks = append(newTasks, app.tasks[i])
	// 	}
	// }

	// // 移除队列中任务
	// if checkTaskQueueWorking(a) {
	// 	app.taskQueue.removeQueueTasks(delQueueTasks)
	// }

	// // 保存任务清单
	// app.tasks = newTasks
	// if err := saveTasks(newTasks, app.configDir); err != nil {
	// 	app.logger.Warn(err)
	// }
	return true
}

// 获取前端任务片段
func (app *App) TaskParts() []Part {
	// return tasksToParts(app.tasks)
	return []Part{}
}

// 任务转任务片段
func tasksToParts(tasks []*task.Task) []Part {
	// parts := make([]Part, len(tasks))
	// for i, task := range tasks {
	// 	parts[i] = *task.part
	// }
	// return parts
	return []Part{}
}

func checkTaskQueueWorking(app *App) bool {
	// return app.taskQueue != nil && app.taskQueue.state != Finished
	return true
}
