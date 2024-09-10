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

	"github.com/Yuelioi/vidor/internal/config"
	"github.com/Yuelioi/vidor/internal/plugin"
	"github.com/Yuelioi/vidor/internal/tools"
	"github.com/Yuelioi/vidor/pkg/downloader"
	"github.com/go-resty/resty/v2"

	"golift.io/xtractr"
)

// 返回主机注册的插件
func (app *App) GetPlugins() map[string]plugin.Manifest {
	ms := make(map[string]plugin.Manifest, 0)

	for _, plugin := range app.plugins {
		mf := plugin.GetManifest()
		ms[mf.ID] = *mf
	}

	return ms
}

// 获取网络插件列表
func fetchPlugins(pluginDir string) ([]*plugin.Manifest, error) {
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
		manifest := plugin.NewManifest(pluginDir)

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
func (app *App) DownloadPlugin(m plugin.Manifest) *plugin.Manifest {

	plugins, err := fetchPlugins(pluginsDir)
	fmt.Printf("plugins: %v\n", plugins[0])
	if err != nil {
		return nil
	}

	var targetManifest *plugin.Manifest
	for _, manifest := range plugins {
		if m.ID == manifest.ID {
			targetManifest = manifest
		}
	}

	if len(targetManifest.DownloadURLs) == 0 {
		return nil
	}

	downUrl := targetManifest.DownloadURLs[0]
	name := tools.ExtractFileNameFromUrl(downUrl)
	name = tools.SanitizeFileName(name)

	tmpDir := filepath.Join(app.location, "tmp")

	err = os.MkdirAll(tmpDir, os.ModePerm)
	if err != nil {
		return nil
	}

	zipPath := filepath.Join(tmpDir, name)
	targetDir := filepath.Join(pluginsDir, targetManifest.ID)
	err = os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		return nil
	}

	downloader, err := downloader.New(
		context.Background(),
		downUrl,
		zipPath,
		true)
	if err != nil {
		return nil
	}

	go downloader.Download()

	tk := time.NewTicker(time.Second)
	defer tk.Stop()

	// 下载
	for range tk.C {
		targetManifest.Status = downloader.Status
		if downloader.State == 1 {
			// 下载中
			runtime.EventsEmit(app.ctx, "plugin.downloading", targetManifest)
		} else {
			// 下载失败
			runtime.EventsEmit(app.ctx, "plugin.downloading", targetManifest)
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
		return nil
	}

	return targetManifest
}

// 运行插件, 并建立连接
func (app *App) RunPlugin(m plugin.Manifest) *plugin.Manifest {
	plugin, ok := app.plugins[m.ID]
	if !ok {
		return nil
	}

	// 运行
	err := plugin.Run(context.Background())
	if err != nil {
		app.logger.Infof("插件开启失败:%s", err)
		return nil
	}

	err = plugin.Init(context.Background())
	if err != nil {
		app.logger.Infof("插件开启失败:%s", err)
		return nil
	}

	return plugin.GetManifest()
}

// 更新插件参数
func (app *App) UpdatePlugin(m plugin.Manifest) *plugin.Manifest {
	plugin, ok := app.plugins[m.ID]
	if !ok {
		return nil
	}

	// TODO 跟新
	err := plugin.Update(context.Background())
	if err != nil {
		app.logger.Infof("插件开启失败:%s", err)
		return nil
	}

	return plugin.GetManifest()
}

func (app *App) StopPlugin(m plugin.Manifest) *plugin.Manifest {
	plugin, ok := app.plugins[m.ID]
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
func (app *App) EnablePlugin(m plugin.Manifest) *plugin.Manifest {
	plugin, ok := app.plugins[m.ID]
	if !ok {
		return nil
	}
	manifest := plugin.GetManifest()
	// 保存配置
	p2 := app.SavePluginConfig(manifest.ID, manifest.PluginConfig)
	if p2 != nil {
		return nil
	}
	return p2
}

// 关闭插件,并禁用插件
func (app *App) DisablePlugin(m plugin.Manifest) *plugin.Manifest {
	plugin, ok := app.plugins[m.ID]
	if !ok {
		app.logger.Infof("没有找到插件:%s", m.ID)
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
	manifest.PluginConfig.Enable = false
	manifest.State = 3

	p2 := app.SavePluginConfig(manifest.ID, manifest.PluginConfig)
	if p2 != nil {
		return nil
	}
	return p2
}

// 保存插件配置
func (app *App) SavePluginConfig(id string, pluginConfig *config.PluginConfig) *plugin.Manifest {
	plugin, ok := app.plugins[id]
	if !ok {
		return nil
	}

	manifest := plugin.GetManifest()

	err := app.SavePluginsConfig(id, pluginConfig)
	if err != nil {
		return nil
	}
	return manifest
}
