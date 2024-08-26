package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	pb "github.com/Yuelioi/vidor/internal/proto"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type MessageData struct {
	Message     string `json:"message"`
	MessageType string `json:"messageType"`
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
	response, err := plugin.Service.GetInfo(ctx, &pb.InfoRequest{
		Url: link,
	})

	if err != nil {
		app.logger.Infof("获取视频信息失败%+v", err)
		return nil
	}
	fmt.Printf("Show Response: %v\n", response)

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
	parseResponse, err := plugin.Service.Parse(ctx, &pb.ParseRequest{Tasks: tasks})

	if err != nil {
		return &pb.ParseResponse{}
	}

	fmt.Println("parseResponse", parseResponse)
	// 更新数据

	// 缓存任务
	app.cache.AddTasks(parseResponse.Tasks)
	return parseResponse
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
