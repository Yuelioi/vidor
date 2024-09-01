package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/Yuelioi/vidor/internal/globals"
	"github.com/Yuelioi/vidor/internal/plugin"
	pb "github.com/Yuelioi/vidor/internal/proto"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// 基于链接获取下载器
func (app *App) selectDownloadPlugin(url string) (*plugin.DownloadPlugin, error) {
	for _, p := range app.plugins {

		base := p.GetManifest()
		if base.Type == "downloader" {
			downloadPlugin, ok := p.(*plugin.DownloadPlugin)
			if !ok {
				return nil, nil
			}

			for _, match := range downloadPlugin.Manifest.Matches {
				reg, err := regexp.Compile(match)
				if err != nil {
					return nil, errors.New("插件正则表达式编译失败: " + err.Error())
				}
				if reg.MatchString(url) {
					return downloadPlugin, nil
				}
			}
		}

	}
	return nil, globals.ErrPluginNotFound
}

// 获取主页选择下载详情列表
//
//   - 1. 获取下载器
//   - 2. 调用展示信息函数
//   - 3. 缓存信息数据
func (a *App) ShowDownloadInfo(link string) *pb.InfoResponse {
	// 清理上次查询任务缓存
	a.cache.ClearTasks()

	// 获取下载器
	p, err := a.selectDownloadPlugin(link)
	if err != nil {
		a.logger.Infof("未找到可用插件%+v", err)
		runtime.EventsEmit(a.ctx, "system.message", &Notice{
			Message:     "未找到可用插件",
			MessageType: "info",
		})
		return nil
	}

	a.logger.Infof("获取视频信息失败%+v", err)
	runtime.EventsEmit(a.ctx, "system.message", &Notice{
		Message:     fmt.Sprintf("获取视频信息失败%s", p.Manifest.Name),
		MessageType: "info",
	})

	// 储存下载器
	a.cache.SetDownloader(p)

	// 传递上下文
	ctx := context.Background()
	ctx = a.GetConfig().InjectMetadata(ctx)
	ctx = plugin.InjectMetadata(ctx, p.Manifest.PluginConfig.Settings)

	// 获取展示信息
	response, err := p.Service.GetInfo(ctx, &pb.InfoRequest{
		Url: link,
	})

	if err != nil {
		a.logger.Infof("获取视频信息失败%+v", err)
		runtime.EventsEmit(a.ctx, "system.message", &Notice{
			Message:     fmt.Sprintf("获取视频信息失败%+v", err),
			MessageType: "error",
		})
		return nil
	}

	// 缓存任务数据
	a.cache.AddTasks(response.Tasks)

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
	home, _ := os.UserHomeDir()
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
