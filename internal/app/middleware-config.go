package app

import (
	"github.com/Yuelioi/vidor/internal/config"
	"github.com/Yuelioi/vidor/internal/models"
)

// 获取主机所有配置信息
func (app *App) GetConfig() *config.Config {
	return app.config
}

// 保存配置文件到本地
func (app *App) SaveConfig(config *config.Config) bool {
	app.config.PluginConfigs = config.PluginConfigs
	app.config.SystemConfig = config.SystemConfig
	// 保存配置文件
	err := app.config.Save()
	if err != nil {
		app.logger.Warnf("保存设置失败%s", err)
	} else {
		app.logger.Info("保存设置成功")
	}
	return err == nil

}

// 保存系统配置
func (app *App) SaveSystemConfig(systemConfig *models.SystemConfig) error {
	app.config.SystemConfig = systemConfig

	return app.config.Save()
}

// 保存插件配置
func (app *App) SavePluginsConfig(id string, pluginConfig *config.PluginConfig) error {
	plugin, ok := app.plugins[id]
	if ok {
		plugin.PluginConfig = pluginConfig
	}
	return app.config.Save()
}
