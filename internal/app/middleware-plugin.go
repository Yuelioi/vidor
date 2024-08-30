package app

import (
	"context"
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
func (app *App) GetPlugins() map[string]*plugin.Plugin {
	return app.plugins
}

// 获取网络插件列表
func fetchPlugins() ([]*plugin.Plugin, error) {
	pluginsUrl := "https://cdn.yuelili.com/market/vidor/plugins.json"
	plugins := make([]*plugin.Plugin, 0)
	resp, err := resty.New().R().SetResult(&plugins).Get(pluginsUrl)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("链接失败")
	}

	return plugins, nil
}

// 下载插件
//
// 1.下载
// 2.解压
// 3.注册到主机
func (app *App) DownloadPlugin(p *plugin.Plugin) *plugin.Plugin {

	// plugins, err := fetchPlugins()
	// fmt.Printf("plugins: %v\n", plugins[0])
	// if err != nil {
	// 	return nil
	// }

	targetPlugin := plugin.NewDownloader("")

	// for _, plugin := range plugins {
	// 	if p.ID == plugin.ID {
	// 		targetPlugin = plugin
	// 	}
	// }

	p2, ok := (*p).(plugin.BasePlugin)
	if !ok {

	}
	if len(p2.DownloadURLs) == 0 {
		return nil
	}

	downUrl := p2.DownloadURLs[0]
	name := tools.ExtractFileNameFromUrl(downUrl)
	name = tools.SanitizeFileName(name)

	tmpDir := filepath.Join(app.location, "tmp")

	err := os.MkdirAll(tmpDir, os.ModePerm)
	if err != nil {
		return nil
	}

	zipPath := filepath.Join(tmpDir, name)
	targetDir := filepath.Join(pluginsDir, p2.ID)
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
		targetPlugin.Status = downloader.Status
		if downloader.State == 1 {
			// 下载中
			runtime.EventsEmit(app.ctx, "plugin.downloading", targetPlugin)
		} else {
			// 下载失败
			runtime.EventsEmit(app.ctx, "plugin.downloading", targetPlugin)
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

	return targetPlugin
}

// 运行插件, 并建立连接
func (app *App) RunPlugin(p *plugin.Plugin) *plugin.Plugin {
	plugin, ok := app.plugins[p.ID]
	if !ok {
		return nil
	}
	// 运行
	err := plugin.Run()
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
func (app *App) UpdatePlugin(p *plugin.Plugin) *plugin.Plugin {
	plugin, ok := app.plugins[p.ID]
	if !ok {
		return nil
	}
	// TODO 跟新
	err := plugin.Run()
	if err != nil {
		app.logger.Infof("插件开启失败:%s", err)
		return nil
	}

	return plugin
}

func (app *App) StopPlugin(p *plugin.Plugin) *plugin.Plugin {
	plugin, ok := app.plugins[p.ID]
	if !ok {
		return nil
	}
	// 停止
	_, err := plugin.Service.Shutdown(context.Background(), nil)
	if err != nil {
		return nil
	}
	p.State = 3
	return p
}

// 启用插件, 但是不会运行
func (app *App) EnablePlugin(p *plugin.Plugin) (*plugin.Plugin, string) {
	plugin, ok := app.plugins[p.ID]
	if !ok {
		app.logger.Infof("没有找到插件:%s", p.ID)
		return nil, fmt.Sprintf("没有找到插件:%s", p.ID)
	}
	plugin.Enable = true
	// 保存配置
	p2 := app.SavePluginConfig(plugin.ID, plugin.PluginConfig)
	if p2 != nil {
		return nil, fmt.Sprintf("保存插件配置失败:%s", p.ID)
	}
	return p2, fmt.Sprintf("保存插件配置失败:%s", p.ID)
}

// 关闭插件,并禁用插件
func (app *App) DisablePlugin(p *plugin.Plugin) *plugin.Plugin {
	plugin, ok := app.plugins[p.ID]
	if !ok {
		app.logger.Infof("没有找到插件:%s", p.ID)
		return nil
	}

	// 关闭插件
	if plugin.State == 1 {
		_, err := plugin.Service.Shutdown(context.Background(), nil)
		if err != nil {
			return nil
		}
	}

	// 禁用并保存配置
	plugin.Enable = false
	plugin.State = 3

	p2 := app.SavePluginConfig(plugin.ID, plugin.PluginConfig)
	if p2 != nil {
		return nil
	}
	return p2
}

// 保存插件配置
func (app *App) SavePluginConfig(id string, pluginConfig *config.PluginConfig) *plugin.Plugin {
	plugin, ok := app.plugins[id]
	if !ok {
		// globals.ErrPluginConfigSave
		return nil
	}

	err := app.SavePluginsConfig(id, pluginConfig)
	if err != nil {
		return nil
	}
	return plugin
}
