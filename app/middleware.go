package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	cmdRuntime "runtime"

	"github.com/Yuelioi/vidor/shared"
	"github.com/Yuelioi/vidor/utils"
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
*/
func (a *App) ShowDownloadInfo(link string) *shared.PlaylistInfo {
	downloader, err := newDownloader(a.downloaders, a.Notice, link)
	// 没有下载器 直接返回空
	if err != nil {
		a.Logger.Info(err)
		return new(shared.PlaylistInfo)
	}

	pi, err := downloader.ShowInfo(link, *a.config, a.Callback)
	if err != nil {
		a.Logger.Warn(err)
		return new(shared.PlaylistInfo)
	}

	a.Logger.Infof("下载: 获取视频元数据成功%s", link)
	return pi
}

/*
	添加下载任务

1. 获取任务目标
2. 创建/添加到任务队列
3. 保存任务信息
*/
func (a *App) AddDownloadTasks(parts []shared.Part, workName string) []shared.Part {
	if len(parts) == 0 {
		return make([]shared.Part, 0)
	}

	var tasks = make([]*Task, 0)

	for _, part := range parts {
		if taskExists(a.tasks, part.Url) {
			logger.Info("任务", "任务已存在", part.Url)
			continue
		} else {

			task, err := createNewTask(part, a.config.DownloadDir, workName)
			if err != nil {
				a.Logger.Warnf("任务: 创建任务失败%s", err)
				continue
			}
			tasks = append(tasks, task)
			a.tasks = append(a.tasks, task)
		}
	}
	// 添加到队列
	if a.taskQueue == nil || a.taskQueue.state == Finished {
		println("任务队列 重新创建")
		a.taskQueue = NewTaskQueue(a, tasks)
	} else {
		println("任务队列 还在使用")
		a.taskQueue.AddTasks(tasks)
	}

	if err := saveTasks(a.tasks, a.configDir); err != nil {
		a.Logger.Warnf("添加任务:保存配置失败 %s", err)
	}

	return tasksToParts(tasks)
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

	for i, task := range a.tasks {
		if task.part.TaskID == uid {

			if checkTaskQueueWorking(a) {
				if task.state == Queue {
					// 1. 正在队列中
					a.Logger.Info("任务移除(队列中):", task.part.Title)
					a.taskQueue.removeQueueTasks([]*Task{task})
				} else {
					// 2. 正在下载中
					a.Logger.Info("任务移除(下载中):", task.part.Title)
					a.taskQueue.stopTask(task)
				}
			}

			// 3. 直接删除
			a.tasks = append(a.tasks[:i], a.tasks[i+1:]...)

			if err := saveTasks(a.tasks, a.configDir); err != nil {
				a.Logger.Warn(err)
			}
			return true
		}
	}
	return false
}

// 移除任务
// 移除完成任务: 去除a.tasks目标 并保存配置
// 移除下载中任务: 调用下载器StopDownload函数 关闭stopChan
// 移除队列中任务: 清理缓存队列的queueTasks
func (a *App) RemoveAllTask(parts []shared.Part) bool {

	newTasks := make([]*Task, 0)
	delQueueTasks := make([]*Task, 0)

	partTaskIDs := make(map[string]shared.Part)

	for _, part := range parts {
		partTaskIDs[part.TaskID] = part
	}

	for i, task := range a.tasks {
		if _, found := partTaskIDs[task.part.TaskID]; found {

			if checkTaskQueueWorking(a) {
				if task.state == Queue {
					// 添加到待删队列
					delQueueTasks = append(delQueueTasks, task)
				} else {
					// 直接调用停止函数
					a.Logger.Info("任务移除(下载中):", task.part.Title)
					task.downloader.StopDownload(task.part, a.Callback)
				}
			}

		} else {
			newTasks = append(newTasks, a.tasks[i])
		}
	}

	// 移除队列中任务
	if checkTaskQueueWorking(a) {
		a.taskQueue.removeQueueTasks(delQueueTasks)
	}

	// 保存任务清单
	a.tasks = newTasks
	if err := saveTasks(newTasks, a.configDir); err != nil {
		a.Logger.Warn(err)
	}
	return true
}

func checkTaskQueueWorking(a *App) bool {
	return a.taskQueue != nil && a.taskQueue.state != Finished
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

func (a *App) SetFFmpegPath(title string) string {
	home, _ := os.UserHomeDir()
	downloadsFolder := filepath.Join(home, "Downloads")

	target, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:                title,
		DefaultDirectory:     downloadsFolder,
		CanCreateDirectories: true,
	})

	if err != nil {
		a.Logger.Error(err)
		return ""
	}

	if err := utils.SetFFmpegPath(target); err != nil {
		a.Logger.Error(err)
		return ""
	}

	return target
}

func (a *App) CheckFFmpeg(target string) bool {
	if err := utils.SetFFmpegPath(target); err != nil {
		a.Logger.Error(err)
		return false
	}

	return true
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

func (a *App) GetConfig() *shared.Config {
	return a.config
}

// 获取前端任务片段
func (a *App) GetTaskParts() []shared.Part {
	return tasksToParts(a.tasks)
}

func (a *App) SaveConfig(config *shared.Config) bool {
	a.config = config
	err := saveConfig(a.configDir, *config)
	if err != nil {
		a.Logger.Warnf("保存设置失败%s", err)
	} else {
		a.Logger.Info("保存设置成功")
	}

	return err == nil
}

// 任务转任务片段
func tasksToParts(tasks []*Task) []shared.Part {
	parts := make([]shared.Part, len(tasks))
	for i, task := range tasks {
		parts[i] = *task.part
	}
	return parts
}
