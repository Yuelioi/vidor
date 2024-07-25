package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	cmdRuntime "runtime"

	"github.com/Yuelioi/vidor/shared"
	"github.com/Yuelioi/vidor/utils"
	"github.com/google/uuid"
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

func (a *App) ShowDownloadInfo(link string) *shared.PlaylistInfo {
	downloader, err := a.newDownloader(link)
	if err != nil {
		return &shared.PlaylistInfo{}
	}

	pi, err := (*downloader).ShowInfo(link, *a.config, a.Callback)
	if err != nil {
		a.Logger.Warn(err)
		return &shared.PlaylistInfo{}
	}

	a.Logger.Infof("获取视频元数据成功%s", link)
	return pi
}

func (a *App) AddDownloadTasks(parts []shared.Part, workName string) []shared.Part {
	if len(parts) == 0 {
		return []shared.Part{}
	}

	var tasks = make([]Task, 0)

	for _, part := range parts {
		if a.taskExists(part.Url) {
			logger.Infof("添加任务:任务已存在 %s", part.Url)
		} else {
			task, err := a.createNewTask(part, workName)
			if err != nil {
				a.Logger.Warnf("添加任务:创建任务失败%s", err)
				continue
			}
			tasks = append(tasks, *task)
			a.tasks = append(a.tasks, *task)
		}
	}
	a.ensureTaskQueue(tasks)

	if err := saveTasks(a.tasks, a.configDir); err != nil {
		// 貌似会不同步, 但是一般不会出问题
		a.Logger.Warnf("添加任务:保存配置失败 %s", err)
	}

	return tasksToParts(tasks)
}

// 检测/创建任务队列
func (a *App) ensureTaskQueue(tasks []Task) {
	if a.taskQueue == nil {
		NewTaskQueue(a, tasks)
	} else {
		a.taskQueue.AddTasks(tasks)
	}
}

func (a *App) createNewTask(part shared.Part, workName string) (*Task, error) {
	downloader, err := a.newDownloader(part.Url)

	if err != nil {
		return nil, err
	}

	return &Task{
		downloader: downloader,
		part: &shared.Part{
			UID:         uuid.New().String(),
			DownloadDir: filepath.Join(a.config.DownloadDir, workName),
			Url:         part.Url,
			Title:       part.Title,
			Thumbnail:   part.Thumbnail,
			Status:      shared.TaskStatus.Queue,
			Quality:     part.Quality,
			CreatedAt:   time.Now(),
			State:       shared.TaskStatus.Queue,
		},
	}, nil
}

func (a *App) taskExists(url string) bool {
	for _, existingTask := range a.tasks {
		if existingTask.part.Url == url {
			return true
		}
	}
	return false
}

func (a *App) RemoveTask(uid string) bool {
	for i, task := range a.tasks {
		if task.part.UID == uid {
			fmt.Printf("task.part.UID: %v\n", task.part.UID)
			if task.downloader != nil {
				(*task.downloader).StopDownload(task.part, a.Callback)
				// 关闭完记得关闭下载器
				a.tasks[i].downloader = nil
			}
			a.tasks = append(a.tasks[:i], a.tasks[i+1:]...)
			if err := saveTasks(a.tasks, a.configDir); err != nil {
				a.Logger.Warn(err)
			}
			return true
		}
	}
	return false
}

// 移除
func (a *App) RemoveAllTask(parts []shared.Part) bool {

	newTasks := make([]Task, 0)
	partUIDs := make(map[string]shared.Part)

	for _, part := range parts {
		partUIDs[part.UID] = part
	}
	fmt.Println("正在移除", len(parts))
	var retainedTasks []Task

	// 有任务时需要加锁
	if a.taskQueue != nil {
		a.taskQueue.mu.Lock()
		defer a.taskQueue.mu.Unlock()
	}

	for i, task := range a.tasks {
		if _, found := partUIDs[task.part.UID]; found {
			if task.downloader != nil {
				if task.part.State == shared.TaskStatus.Queue {
					retainedTasks = append(retainedTasks, task)
				} else {
					(*task.downloader).StopDownload(task.part, a.Callback)
				}
				a.tasks[i].downloader = nil
			}
		} else {
			newTasks = append(newTasks, a.tasks[i])
		}
	}

	// 移除正在下载的任务
	if a.taskQueue != nil {
		a.taskQueue.RemoveTasks(retainedTasks)
	}

	if err := saveTasks(newTasks, a.configDir); err != nil {
		a.Logger.Warn(err)
		return false
	}
	a.tasks = newTasks
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

func (a *App) SetFFmpegPath(title string) string {
	home, _ := os.UserHomeDir()
	downloadsFolder := filepath.Join(home, "Downloads")

	target, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
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

// 任务转任务片段
func tasksToParts(tasks []Task) []shared.Part {
	parts := make([]shared.Part, len(tasks))
	for i, task := range tasks {
		parts[i] = *task.part
	}
	return parts
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
