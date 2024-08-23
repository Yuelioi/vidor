package app

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	cmdRuntime "runtime"

	pb "github.com/Yuelioi/vidor/internal/proto"
	"github.com/golang/protobuf/ptypes/empty"
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

func (a *App) InitPlugin() *empty.Empty {
	return nil
}

func (a *App) UpdatePlugin() *empty.Empty {
	return nil
}

/*
	获取主页选择下载详情列表

1. 获取下载器
2. 调用展示信息函数
3. 缓存数据
*/
func (a *App) ShowDownloadInfo(link string) *pb.InfoResponse {

	// 清理上次查询任务缓存
	a.cache.ClearTasks()

	// 获取下载器
	plugin, err := a.selectPlugin(link)
	if err != nil {
		return &pb.InfoResponse{}
	}
	logger.Infof("检测到可用插件%s", plugin.Name)

	// 储存下载器
	a.cache.SetDownloader(plugin)

	// 传递上下文
	ctx := context.Background()

	// 获取展示信息
	response, err := plugin.service.GetInfo(ctx, &pb.InfoRequest{
		Url: link,
	})

	if err != nil {
		logger.Infof("获取视频信息失败%+v", err)
		return nil
	}
	fmt.Printf("Show Response: %v\n", response)

	// 缓存任务数据
	a.cache.SetTasks(response.Tasks)

	return response
}

type taskMap struct {
	id        string
	formatIds []string
}

// 过滤 segments 中的 formats
func filterSegments(segments []*pb.Segment, formatSet map[string]struct{}) {
	for _, seg := range segments {
		filteredFormats := []*pb.Format{}
		for _, format := range seg.Formats {
			if _, exists := formatSet[format.Id]; exists {
				filteredFormats = append(filteredFormats, format)
			}
		}
		seg.Formats = filteredFormats
	}
}

/*
解析数据
*/
func (a *App) ParsePlaylist(ids []string) *pb.ParseResponse {

	// 获取任务缓存数据
	tasks, err := a.cache.Tasks(ids)

	fmt.Printf("tasks: %v\n", tasks)
	if err != nil {
		return &pb.ParseResponse{}
	}

	// 获取缓存下载器
	plugin := a.cache.Downloader()

	// 传递上下文
	ctx := context.Background()

	// 解析
	parseResponse, err := plugin.service.Parse(ctx, &pb.ParseRequest{Tasks: tasks})

	if err != nil {
		return &pb.ParseResponse{}
	}

	fmt.Println("parseResponse", parseResponse)
	// 更新数据

	// 缓存任务
	a.cache.SetTasks(parseResponse.Tasks)
	return parseResponse
}

/*
	添加下载任务

1. 获取任务目标
2. 创建/添加到任务队列
3. 保存任务信息
*/
func (a *App) AddDownloadTasks(taskMaps []taskMap) bool {

	// 获取任务
	tasks := []*pb.Task{}
	for _, taskMap := range taskMaps {
		cacheTask, ok := a.cache.Task(taskMap.id)
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
	plugin := a.cache.Downloader()

	// 清除任务缓存
	a.cache.ClearTasks()

	stream, err := plugin.service.Download(context.Background(), &pb.DownloadRequest{
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
func (a *App) RemoveAllTask(parts []Part) bool {

	// newTasks := make([]*Task, 0)
	// delQueueTasks := make([]*Task, 0)

	// partTaskIDs := make(map[string]Part)

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

func (a *App) GetPlugins() []*Plugin {
	return a.plugins
}

// 获取前端任务片段
func (a *App) GetTaskParts() []Part {
	// return tasksToParts(a.tasks)
	return []Part{}
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
func tasksToParts(tasks []*Task) []Part {
	// parts := make([]Part, len(tasks))
	// for i, task := range tasks {
	// 	parts[i] = *task.part
	// }
	// return parts
	return []Part{}
}
