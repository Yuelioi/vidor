package app

import (
	"github.com/Yuelioi/vidor/internal/config"
	"github.com/Yuelioi/vidor/internal/models"
)

// 获取主机配置信息
func (app *App) GetConfig() *config.Config {
	return app.config
}

// 保存配置文件到本地
func (app *App) SaveConfig(config *config.Config) bool {
	// 保存配置文件
	err := app.config.Save()
	if err != nil {
		app.logger.Warnf("保存设置失败%s", err)
	} else {
		app.logger.Info("保存设置成功")
	}
	return err == nil

}

// 修改系统配置(不会保存)
func (app *App) UpdateSystemConfig(systemConfig *models.SystemConfig) *App {
	app.config.SystemConfig = systemConfig
	return app
}

// 修改插件配置(不会保存)
func (app *App) UpdatePluginsConfig(id string, pluginConfig *config.PluginConfig) *App {
	plugin, ok := app.plugins[id]
	if ok {
		plugin.PluginConfig = pluginConfig
	}
	return app
}
