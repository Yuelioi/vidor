package app

// APP 中间件
// 前端 >> 后端 下载任务

import (
	pb "github.com/Yuelioi/vidor/internal/proto"
)

func (app *App) GetTasks() []*pb.Task {
	return nil
}

/*
	添加下载任务

1. 获取任务目标
2. 创建/添加到任务队列
3. 保存任务信息
*/
func (a *App) AddDownloadTasks(tasks []*pb.Task) bool {
	a.taskManager.AddTasks(tasks...)
	return true
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
func (app *App) RemoveTask(id string) bool {
	return false
}

// 移除任务
// 移除完成任务: 去除app.tasks目标 并保存配置
// 移除下载中任务: 调用下载器StopDownload函数 关闭stopChan
// 移除队列中任务: 清理缓存队列的queueTasks
func (app *App) RemoveAllTask(ids []string) bool {
	return true
}
