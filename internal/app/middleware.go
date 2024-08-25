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
	"github.com/go-resty/resty/v2"
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

// 插件

func (app *App) DownloadPlugin(p *Plugin) *Plugin {
	pluginDir := fmt.Sprintf("https://cdn.yuelili.com/market/vidor/plugins/", p.ID)

	client := &resty.Client{}
	resp, err := client.R().Get(pluginDir + "/" + p.Location)

	if err != nil {
		return nil
	}
	fmt.Println(resp)

	return nil
}

// 运行插件, 并建立连接
func (app *App) RunPlugin(p *Plugin) *Plugin {
	plugin, ok := app.plugins[p.ID]
	if !ok {
		return nil
	}
	// 运行
	err := plugin.Run(app.config)
	if err != nil {
		app.logger.Infof("插件开启失败:%s", err)
		return nil
	}

	err = plugin.Init()
	if err != nil {
		app.logger.Infof("插件开启失败:%s", err)
		return nil
	}

	return plugin
}

// 更新插件参数
func (app *App) UpdatePlugin(p *Plugin) *Plugin {
	plugin, ok := app.plugins[p.ID]
	if !ok {
		return nil
	}
	// TODO 跟新
	err := plugin.Run(app.config)
	if err != nil {
		app.logger.Infof("插件开启失败:%s", err)
		return nil
	}

	return plugin
}

func (app *App) StopPlugin(p *Plugin) *Plugin {
	plugin, ok := app.plugins[p.ID]
	if !ok {
		return nil
	}
	// 停止
	_, err := plugin.service.Shutdown(context.Background(), nil)
	if err != nil {
		return nil
	}
	p.State = 3
	return p
}

// 启用插件, 但是不会运行
func (app *App) EnablePlugin(p *Plugin) (*Plugin, string) {
	plugin, ok := app.plugins[p.ID]
	if !ok {
		app.logger.Infof("没有找到插件:%s", p.ID)
		return nil, fmt.Sprintf("没有找到插件:%s", p.ID)
	}
	plugin.Enable = true
	// 保存配置
	p2, err := app.SavePluginConfig(plugin.ID, plugin.PluginConfig)
	if err != nil {
		return nil, fmt.Sprintf("保存插件配置失败:%s", p.ID)
	}
	return p2, fmt.Sprintf("保存插件配置失败:%s", p.ID)
}

// 关闭插件,并禁用插件
func (app *App) DisablePlugin(p *Plugin) *Plugin {
	plugin, ok := app.plugins[p.ID]
	if !ok {
		app.logger.Infof("没有找到插件:%s", p.ID)
		return nil
	}

	// 关闭插件
	if plugin.State == 1 {
		_, err := plugin.service.Shutdown(context.Background(), nil)
		if err != nil {
			return nil
		}
	}

	// 禁用并保存配置
	plugin.Enable = false
	plugin.State = 3

	p2, err := app.SavePluginConfig(plugin.ID, plugin.PluginConfig)
	if err != nil {
		return nil
	}
	return p2
}

// 保存插件配置
func (app *App) SavePluginConfig(id string, pluginConfig *PluginConfig) (*Plugin, error) {
	plugin, ok := app.plugins[id]
	if !ok {
		return nil, pluginConfigSaveError
	}

	err := app.UpdatePluginsConfig(id, pluginConfig).config.Save()
	if err != nil {
		return nil, err
	}
	return plugin, nil
}

/*
	获取主页选择下载详情列表

1. 获取下载器
2. 调用展示信息函数
3. 缓存数据
*/
func (app *App) ShowDownloadInfo(link string) *pb.InfoResponse {
	// 清理上次查询任务缓存
	app.cache.ClearTasks()

	// 获取下载器
	plugin, err := app.selectPlugin(link)
	if err != nil {
		return &pb.InfoResponse{}
	}
	app.logger.Infof("检测到可用插件%s", plugin.Name)

	// 储存下载器
	app.cache.SetDownloader(plugin)

	// 传递上下文
	ctx := context.Background()

	// 获取展示信息
	response, err := plugin.service.GetInfo(ctx, &pb.InfoRequest{
		Url: link,
	})

	if err != nil {
		app.logger.Infof("获取视频信息失败%+v", err)
		return nil
	}
	fmt.Printf("Show Response: %v\n", response)

	// 缓存任务数据
	app.cache.SetTasks(response.Tasks)

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
func (app *App) ParsePlaylist(ids []string) *pb.ParseResponse {

	// 获取任务缓存数据
	tasks, err := app.cache.Tasks(ids)

	fmt.Printf("tasks: %v\n", tasks)
	if err != nil {
		return &pb.ParseResponse{}
	}

	// 获取缓存下载器
	plugin := app.cache.Downloader()

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
	app.cache.SetTasks(parseResponse.Tasks)
	return parseResponse
}

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

func checkTaskQueueWorking(app *App) bool {
	// return app.taskQueue != nil && app.taskQueue.state != Finished
	return true
}

func (app *App) SetDownloadDir(title string) string {
	home, _ := os.UserHomeDir()
	downloadsFolder := filepath.Join(home, "Downloads")

	target, err := runtime.OpenDirectoryDialog(app.ctx, runtime.OpenDialogOptions{
		Title:            title,
		DefaultDirectory: downloadsFolder,
	})

	if err != nil {
		app.logger.Error(err)
		return ""
	}
	return target
}

func (app *App) OpenExplorer(dir string) error {
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

func (app *App) OpenFileWithSystemPlayer(filePath string) error {
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

func (app *App) GetConfig() *Config {
	return app.config
}

func (app *App) GetPlugins() map[string]*Plugin {
	return app.plugins
}

// 获取前端任务片段
func (app *App) TaskParts() []Part {
	// return tasksToParts(app.tasks)
	return []Part{}
}

// 保存配置文件到本地
func (app *App) SaveConfig(config *Config) bool {

	// 保存配置文件
	err := app.config.Save()
	if err != nil {
		app.logger.Warnf("保存设置失败%s", err)
	} else {
		app.logger.Info("保存设置成功")
	}
	return err == nil

}

// 修改系统配置
func (app *App) UpdateSystemConfig(systemConfig *SystemConfig) *App {
	app.config.SystemConfig = systemConfig
	return app
}

// 修改插件配置
func (app *App) UpdatePluginsConfig(id string, pluginConfig *PluginConfig) *App {
	plugin, ok := app.plugins[id]
	if ok {
		plugin.PluginConfig = pluginConfig
	}
	return app
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
