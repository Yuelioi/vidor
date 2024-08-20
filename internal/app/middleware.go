package app

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	cmdRuntime "runtime"

	pb "github.com/Yuelioi/vidor/internal/proto"
	"github.com/Yuelioi/vidor/internal/shared"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type MessageData struct {
	Message     string `json:"message"`
	MessageType string `json:"messageType"`
}
type TaskResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

/*
	获取主页选择下载详情列表

1. 获取下载器
2. 调用展示信息函数
3. 缓存数据
*/
func (a *App) ShowDownloadInfo(link string) *pb.VideoInfoResponse {
	var wg sync.WaitGroup

	// 清理上次查询缓存
	a.cache.Clear()

	// 获取下载器
	plugin := a.plugins[0]
	logger.Infof("检测到可用插件%s", plugin.Name)

	// 获取展示信息
	response, err := plugin.Service.GetVideoInfo(context.Background(), &pb.VideoInfoRequest{
		Url: link,
	})

	if err != nil {
		logger.Infof("获取视频信息失败%+v", err)
		return nil
	}
	fmt.Printf("Show Response: %v\n", response)

	// 缓存数据
	for _, info := range response.Tasks {
		wg.Add(1)
		go func(info *pb.Task) {
			defer wg.Done()
			a.cache.Set(info.Id, info)
		}(info)
	}
	wg.Wait()

	return response
}

/*
解析数据
*/
func (a *App) ParsePlaylist(ids []string) *pb.ParseResponse {

	// 获取缓存数据
	var wg sync.WaitGroup

	infos := []*pb.Task{}
	for _, id := range ids {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			info, ok := a.cache.Get(id)
			if ok {
				infos = append(infos, info)
			}
		}(id)
	}

	wg.Wait()

	// 获取下载器
	plugin := a.plugins[0]

	// 解析
	parseResponse, err := plugin.Service.ParseEpisodes(context.Background(), &pb.ParseRequest{Tasks: infos})

	if err != nil {
		return &pb.ParseResponse{}
	}

	return parseResponse
}

/*
	添加下载任务

1. 获取任务目标
2. 创建/添加到任务队列
3. 保存任务信息
*/
func (a *App) AddDownloadTasks(infos []*pb.Task) []shared.Part {

	_, err := a.plugins[0].Service.Download(context.Background(), &pb.DownloadRequest{
		Tasks: infos,
	})

	if err != nil {
		return []shared.Part{}
	}

	// if len(parts) == 0 {
	// 	return make([]shared.Part, 0)
	// }

	// var tasks = make([]*Task, 0)

	// for _, part := range parts {
	// 	if taskExists(a.tasks, part.URL) {
	// 		logger.Info("任务", "任务已存在", part.URL)
	// 		continue
	// 	} else {

	// 		task, err := createNewTask(part, a.config.DownloadDir, workName)
	// 		if err != nil {
	// 			a.Logger.Warnf("任务: 创建任务失败%s", err)
	// 			continue
	// 		}
	// 		tasks = append(tasks, task)
	// 		a.tasks = append(a.tasks, task)
	// 	}
	// }
	// // 添加到队列
	// if a.taskQueue == nil || a.taskQueue.state == Finished {
	// 	println("任务队列 重新创建")
	// 	a.taskQueue = NewTaskQueue(a, tasks)
	// } else {
	// 	println("任务队列 还在使用")
	// 	a.taskQueue.AddTasks(tasks)
	// }

	// if err := saveTasks(a.tasks, a.configDir); err != nil {
	// 	a.Logger.Warnf("添加任务:保存配置失败 %s", err)
	// }

	return []shared.Part{}

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
func (a *App) RemoveTask(uid string) bool {

	// for i, task := range a.tasks {
	// 	if task.part.TaskID == uid {

	// 		if checkTaskQueueWorking(a) {
	// 			if task.state == Queue {
	// 				// 1. 正在队列中
	// 				a.Logger.Info("任务移除(队列中):", task.part.Title)
	// 				a.taskQueue.removeQueueTasks([]*Task{task})
	// 			} else {
	// 				// 2. 正在下载中
	// 				a.Logger.Info("任务移除(下载中):", task.part.Title)
	// 				a.taskQueue.stopTask(task)
	// 			}
	// 		}

	// 		// 3. 直接删除
	// 		a.tasks = append(a.tasks[:i], a.tasks[i+1:]...)

	// 		if err := saveTasks(a.tasks, a.configDir); err != nil {
	// 			a.Logger.Warn(err)
	// 		}
	// 		return true
	// 	}
	// }
	return false
}

// 移除任务
// 移除完成任务: 去除a.tasks目标 并保存配置
// 移除下载中任务: 调用下载器StopDownload函数 关闭stopChan
// 移除队列中任务: 清理缓存队列的queueTasks
func (a *App) RemoveAllTask(parts []shared.Part) bool {

	// newTasks := make([]*Task, 0)
	// delQueueTasks := make([]*Task, 0)

	// partTaskIDs := make(map[string]shared.Part)

	// for _, part := range parts {
	// 	partTaskIDs[part.TaskID] = part
	// }

	// for i, task := range a.tasks {
	// 	if _, found := partTaskIDs[task.part.TaskID]; found {

	// 		if checkTaskQueueWorking(a) {
	// 			if task.state == Queue {
	// 				// 添加到待删队列
	// 				delQueueTasks = append(delQueueTasks, task)
	// 			} else {
	// 				// 直接调用停止函数
	// 				a.Logger.Info("任务移除(下载中):", task.part.Title)
	// 				task.downloader.Cancel(context.Background(), task.part)
	// 			}
	// 		}

	// 	} else {
	// 		newTasks = append(newTasks, a.tasks[i])
	// 	}
	// }

	// // 移除队列中任务
	// if checkTaskQueueWorking(a) {
	// 	a.taskQueue.removeQueueTasks(delQueueTasks)
	// }

	// // 保存任务清单
	// a.tasks = newTasks
	// if err := saveTasks(newTasks, a.configDir); err != nil {
	// 	a.Logger.Warn(err)
	// }
	return true
}

func checkTaskQueueWorking(a *App) bool {
	// return a.taskQueue != nil && a.taskQueue.state != Finished
	return true
}

func (a *App) SetDownloadDir(title string) string {
	home, _ := os.UserHomeDir()
	downloadsFolder := filepath.Join(home, "Downloads")

	target, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            title,
		DefaultDirectory: downloadsFolder,
	})

	if err != nil {
		a.Logger.Error(err)
		return ""
	}
	return target
}

func (a *App) OpenExplorer(dir string) error {
	dir = filepath.FromSlash(dir)
	var cmd *exec.Cmd

	switch cmdRuntime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", dir)
	case "darwin":
		cmd = exec.Command("open", dir)
	case "linux":
		cmd = exec.Command("xdg-open", dir)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}

func (a *App) OpenFileWithSystemPlayer(filePath string) error {
	var cmd *exec.Cmd

	switch cmdRuntime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", filePath)
	case "darwin":
		cmd = exec.Command("open", filePath)
	case "linux":
		cmd = exec.Command("xdg-open", filePath)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}

func (a *App) GetConfig() *Config {
	return a.config

}

// 获取前端任务片段
func (a *App) GetTaskParts() []shared.Part {
	// return tasksToParts(a.tasks)
	return []shared.Part{}
}

func (a *App) SaveConfig(config *Config) bool {
	a.config = config
	err := a.config.SaveConfig()
	if err != nil {
		a.Logger.Warnf("保存设置失败%s", err)
	} else {
		a.Logger.Info("保存设置成功")
	}

	return err == nil

}

// 任务转任务片段
func tasksToParts(tasks []*Task) []shared.Part {
	// parts := make([]shared.Part, len(tasks))
	// for i, task := range tasks {
	// 	parts[i] = *task.part
	// }
	// return parts
	return []shared.Part{}
}
