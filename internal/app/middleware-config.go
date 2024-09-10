package app

import (
	"github.com/Yuelioi/vidor/internal/config"
)

// 获取主机所有配置信息
func (app *App) GetConfig() *config.Config {
	return app.config
}

// 保存所有配置到本地
func (app *App) SaveConfig(config *config.Config) bool {
	app.config = config
	// 保存配置文件
	err := app.config.Save()
	if err != nil {
		app.logger.Warnf("保存设置失败%s", err)
	} else {
		app.logger.Info("保存设置成功")
	}
	return err == nil

}
