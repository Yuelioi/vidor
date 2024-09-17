package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Yuelioi/vidor/internal/notify"
	"github.com/Yuelioi/vidor/internal/plugin"
	pb "github.com/Yuelioi/vidor/internal/proto"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// 获取主页选择下载详情列表
//
//   - 1. 获取下载器
//   - 2. 调用展示信息函数
//   - 3. 缓存信息数据
func (a *App) ShowDownloadInfo(link string) *pb.InfoResponse {
	// 清理上次查询任务缓存
	a.cache.ClearTasks()

	// 获取下载器
	p, err := a.manager.Select(link)
	if err != nil {
		a.notification.Send(a.ctx, notify.Notice{
			EventName:  "plugin.show",
			Content:    "未找到可用插件",
			NoticeType: "info",
			Provider:   p.Manifest.Name},
		)
		return nil
	}

	a.notification.Send(a.ctx, notify.Notice{
		EventName:  "plugin.show",
		Content:    "获取视频信息失败",
		NoticeType: "info",
		Provider:   p.Manifest.Name},
	)

	// 储存下载器
	a.cache.SetDownloader(p)

	// 传递上下文
	ctx := context.Background()
	ctx = a.GetConfig().InjectMetadata(ctx)
	ctx = plugin.InjectMetadata(ctx, p.Manifest.Settings)

	// 获取展示信息
	response, err := p.Service.GetInfo(ctx, &pb.InfoRequest{
		Url: link,
	})

	if err != nil {
		a.notification.Send(a.ctx, notify.Notice{
			EventName:  "plugin.show",
			Content:    "获取视频信息失败",
			NoticeType: "info",
			Provider:   p.Manifest.Name},
		)
		return nil
	}

	// 缓存任务数据
	a.cache.AddTasks(response.Tasks)

	return response
}

/*
解析数据
*/
func (a *App) ParsePlaylist(ids []string) *pb.TasksResponse {

	// 获取任务缓存数据
	tasks, err := a.cache.Tasks(ids)

	fmt.Printf("tasks: %v\n", tasks)
	if err != nil {
		return nil
	}

	// 获取缓存下载器
	plugin := a.cache.Downloader()

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
	a.cache.AddTasks(TasksResponse.Tasks)
	return TasksResponse
}

func (a *App) SetDownloadDir(title string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		a.logger.Warnf("获取用户文件夹失败:%s", err)
		return ""
	}
	downloadsFolder := filepath.Join(home, "Downloads")

	target, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            title,
		DefaultDirectory: downloadsFolder,
	})

	if err != nil {
		a.logger.Error(err)
		return ""
	}
	return target
}
