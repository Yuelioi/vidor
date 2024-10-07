package app

import (
	"context"
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

	// 获取下载器
	p, err := a.manager.Select(link)
	if err != nil {
		a.notification.Send(notify.Notice{
			EventName:  "system.notice",
			Content:    "未找到可用插件",
			NoticeType: "info",
			Provider:   p.Manifest.Name},
		)
		return nil
	}

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
		a.notification.Send(notify.Notice{
			EventName:  "system.notice",
			Content:    "获取视频信息失败" + err.Error(),
			NoticeType: "info",
			Provider:   p.Manifest.Name},
		)
		return nil
	}

	// 本地化封面
	response.Cover = "/files/" + response.Cover

	// 设置工作文件夹
	response.DownloaderDir = a.config.DownloadDir

	return response
}

/*
解析数据
*/
func (a *App) ParsePlaylist(tasks []*pb.Task) *pb.TasksResponse {

	// 获取缓存下载器
	plugin := a.cache.Downloader()

	// 传递上下文
	ctx := context.Background()

	// 解析
	tasksResponse, err := plugin.Service.Parse(ctx, &pb.TasksRequest{Tasks: tasks})
	if err != nil {
		return nil
	}

	// 缓存任务
	return tasksResponse
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
