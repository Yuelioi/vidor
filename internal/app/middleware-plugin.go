package app

import (
	"context"
	"fmt"

	"github.com/Yuelioi/vidor/internal/config"
	"github.com/Yuelioi/vidor/internal/globals"
	"github.com/Yuelioi/vidor/internal/plugin"
	"github.com/go-resty/resty/v2"
)

func (app *App) GetPlugins() map[string]*plugin.Plugin {
	return app.plugins
}

func (app *App) DownloadPlugin(p *plugin.Plugin) *plugin.Plugin {
	pluginDir := fmt.Sprintf("https://cdn.yuelili.com/market/vidor/plugins/%s", p.ID)

	client := &resty.Client{}
	resp, err := client.R().Get(pluginDir + "/" + p.Location)

	if err != nil {
		return nil
	}
	fmt.Println(resp)

	return nil
}

// 运行插件, 并建立连接
func (app *App) RunPlugin(p *plugin.Plugin) *plugin.Plugin {
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
func (app *App) UpdatePlugin(p *plugin.Plugin) *plugin.Plugin {
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
	p2, err := app.SavePluginConfig(plugin.ID, plugin.PluginConfig)
	if err != nil {
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

	p2, err := app.SavePluginConfig(plugin.ID, plugin.PluginConfig)
	if err != nil {
		return nil
	}
	return p2
}

// 保存插件配置
func (app *App) SavePluginConfig(id string, pluginConfig *config.PluginConfig) (*plugin.Plugin, error) {
	plugin, ok := app.plugins[id]
	if !ok {
		return nil, globals.ErrPluginConfigSave
	}

	err := app.UpdatePluginsConfig(id, pluginConfig).config.Save()
	if err != nil {
		return nil, err
	}
	return plugin, nil
}
