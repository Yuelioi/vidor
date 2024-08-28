package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	pb "github.com/Yuelioi/vidor/internal/proto"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// 获取主页选择下载详情列表
//
//   - 1. 获取下载器
//   - 2. 调用展示信息函数
//   - 3. 缓存信息数据
func (app *App) ShowDownloadInfo(link string) *pb.InfoResponse {
	// 清理上次查询任务缓存
	app.cache.ClearTasks()

	// 获取下载器
	plugin, err := app.selectPlugin(link)
	if err != nil {
		app.logger.Infof("未找到可用插件%+v", err)
		runtime.EventsEmit(app.ctx, "system.message", &Notice{
			Message:     "未找到可用插件",
			MessageType: "info",
		})
		return nil
	}

	app.logger.Infof("获取视频信息失败%+v", err)
	runtime.EventsEmit(app.ctx, "system.message", &Notice{
		Message:     fmt.Sprintf("获取视频信息失败%s", plugin.Name),
		MessageType: "info",
	})

	// 储存下载器
	app.cache.SetDownloader(plugin)

	// 传递上下文
	ctx := context.Background()
	ctx = app.GetConfig().InjectMetadata(ctx)
	ctx = plugin.InjectMetadata(ctx)

	// 获取展示信息
	response, err := plugin.Service.GetInfo(ctx, &pb.InfoRequest{
		Url: link,
	})

	if err != nil {
		app.logger.Infof("获取视频信息失败%+v", err)
		runtime.EventsEmit(app.ctx, "system.message", &Notice{
			Message:     fmt.Sprintf("获取视频信息失败%+v", err),
			MessageType: "error",
		})
		return nil
	}

	// 缓存任务数据
	app.cache.AddTasks(response.Tasks)

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
func (app *App) ParsePlaylist(ids []string) *pb.TasksResponse {

	// 获取任务缓存数据
	tasks, err := app.cache.Tasks(ids)

	fmt.Printf("tasks: %v\n", tasks)
	if err != nil {
		return nil
	}

	// 获取缓存下载器
	plugin := app.cache.Downloader()

	// 传递上下文
	ctx := context.Background()

	// 解析
	TasksResponse, err := plugin.Service.Parse(ctx, &pb.TasksRequest{Tasks: tasks})

	if err != nil {
		return nil
	}

	fmt.Println("TasksResponse", TasksResponse)
	// 更新数据

	// 缓存任务
	app.cache.AddTasks(TasksResponse.Tasks)
	return TasksResponse
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
