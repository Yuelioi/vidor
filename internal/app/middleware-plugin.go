package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/Yuelioi/vidor/internal/notify"
	"github.com/Yuelioi/vidor/internal/plugin"
	"github.com/Yuelioi/vidor/internal/tools"
	"github.com/Yuelioi/vidor/pkg/downloader"
	"github.com/go-resty/resty/v2"

	"golift.io/xtractr"
)

// 返回主机注册的插件
func (a *App) GetPlugins() map[string]plugin.Manifest {
	ms := make(map[string]plugin.Manifest, 0)

	for _, plugin := range a.plugins {
		mf := plugin.GetManifest()
		ms[mf.ID] = *mf
	}
	return ms
}

// 获取网络插件列表
func fetchPlugins(pluginsDir string) ([]*plugin.Manifest, error) {
	pluginsUrl := "https://cdn.yuelili.com/market/vidor/plugins.json"

	resp, err := resty.New().R().Get(pluginsUrl)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("链接失败")
	}

	var rawPlugins []map[string]interface{}
	err = json.Unmarshal(resp.Body(), &rawPlugins)
	if err != nil {
		return nil, err
	}

	plugins := make([]*plugin.Manifest, 0, len(rawPlugins))
	for _, rawPlugin := range rawPlugins {
		manifest := plugin.NewManifest(pluginsDir)

		manifestData, err := json.Marshal(rawPlugin)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(manifestData, manifest)
		if err != nil {
			return nil, err
		}

		plugins = append(plugins, manifest)
	}

	return plugins, nil
}

// 下载插件
//
// 1.下载
// 2.解压
// 3.注册到主机
func (a *App) DownloadPlugin(m plugin.Manifest) *plugin.Manifest {

	plugins, err := fetchPlugins(a.appDirs.Plugins)
	fmt.Printf("plugins: %v\n", plugins[0])
	if err != nil {
		a.notification.Send(a.ctx, notify.Notice{
			EventName:  "system.notice",
			Content:    "获取插件信息失败:" + err.Error(),
			NoticeType: "info",
			Provider:   "system",
		})
		return nil
	}

	var targetManifest *plugin.Manifest
	for _, manifest := range plugins {
		if m.ID == manifest.ID {
			targetManifest = manifest
		}
	}

	if len(targetManifest.DownloadURLs) == 0 {
		a.notification.Send(a.ctx, notify.Notice{
			EventName:  "system.notice",
			Content:    "获取插件链接失败",
			NoticeType: "info",
			Provider:   "system",
		})
		return nil
	}

	downUrl := targetManifest.DownloadURLs[0]
	name := tools.ExtractFileNameFromUrl(downUrl)
	name = tools.SanitizeFileName(name)

	zipPath := filepath.Join(a.appDirs.Temps, name)
	targetDir := filepath.Join(a.appDirs.Plugins, targetManifest.ID)
	err = os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		return nil
	}

	dl, err := downloader.New(
		context.Background(),
		downUrl,
		zipPath,
		true)
	if err != nil {
		a.notification.Send(a.ctx, notify.Notice{
			EventName:  "system.notice",
			Content:    "下载插件失败" + err.Error(),
			NoticeType: "info",
			Provider:   "system",
		})
		return nil
	}

	go dl.Download()

	tk := time.NewTicker(time.Second)
	defer tk.Stop()

	// 下载
	for range tk.C {
		targetManifest.Status = dl.Status
		if dl.State == 1 {
			// 下载中
			runtime.EventsEmit(a.ctx, "plugin.downloading", targetManifest)
		} else {
			// 下载失败
			runtime.EventsEmit(a.ctx, "plugin.downloading", targetManifest)
			break
		}
	}

	// 解压
	x := &xtractr.XFile{
		FilePath:  zipPath,
		OutputDir: targetDir,
	}

	_, files, _, err := xtractr.ExtractFile(x)
	if err != nil || files == nil {
		a.notification.Send(a.ctx, notify.Notice{
			EventName:  "system.notice",
			Content:    "解压插件信息失败" + err.Error(),
			NoticeType: "info",
			Provider:   "system",
		})
		return nil
	}

	// 注册
	a.registerPlugin(&m)

	return targetManifest
}

// 运行插件, 并建立连接
func (a *App) RunPlugin(m plugin.Manifest) *plugin.Manifest {
	plugin, ok := a.plugins[m.ID]
	if !ok {
		return nil
	}

	// 运行
	err := plugin.Run(context.Background())
	if err != nil {
		a.logger.Infof("插件开启失败:%s", err)
		return nil
	}

	err = plugin.Init(context.Background())
	if err != nil {
		a.logger.Infof("插件开启失败:%s", err)
		return nil
	}

	return plugin.GetManifest()
}

// 更新插件参数
func (a *App) UpdatePlugin(m plugin.Manifest) *plugin.Manifest {
	plugin, ok := a.plugins[m.ID]
	if !ok {
		return nil
	}

	// TODO 跟新
	err := plugin.Update(context.Background())
	if err != nil {
		a.logger.Infof("插件开启失败:%s", err)
		return nil
	}

	return plugin.GetManifest()
}

func (a *App) StopPlugin(m plugin.Manifest) *plugin.Manifest {
	plugin, ok := a.plugins[m.ID]
	if !ok {
		return nil
	}
	// 停止
	err := plugin.Shutdown(context.Background())
	if err != nil {
		return nil
	}

	return plugin.GetManifest()
}

// 启用插件, 但是不会运行
func (a *App) EnablePlugin(m plugin.Manifest) *plugin.Manifest {
	plugin, ok := a.plugins[m.ID]
	if !ok {
		return nil
	}
	manifest := plugin.GetManifest()
	// 保存配置
	p2 := a.SavePluginConfig(manifest.ID, manifest)
	if p2 != nil {
		return nil
	}
	return p2
}

// 关闭插件,并禁用插件
func (a *App) DisablePlugin(m plugin.Manifest) *plugin.Manifest {
	plugin, ok := a.plugins[m.ID]
	if !ok {
		a.logger.Infof("没有找到插件:%s", m.ID)
		return nil
	}

	manifest := plugin.GetManifest()

	// 关闭插件
	if manifest.State == 1 {
		err := plugin.Shutdown(context.Background())
		if err != nil {
			return nil
		}
	}

	// 禁用并保存配置
	manifest.Enable = false
	manifest.State = 3

	p2 := a.SavePluginConfig(manifest.ID, manifest)
	if p2 != nil {
		return nil
	}
	return p2
}

// 保存插件配置
func (a *App) SavePluginConfig(id string, m *plugin.Manifest) *plugin.Manifest {
	plugin, ok := a.plugins[id]
	if !ok {
		return nil
	}

	manifest := plugin.GetManifest()

	err := plugin.GetManifest().Save()
	if err != nil {
		return nil
	}
	return manifest
}
