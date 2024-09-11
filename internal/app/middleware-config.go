package app

import (
	"github.com/Yuelioi/vidor/internal/config"
)

// 获取主机配置
func (a *App) GetConfig() *config.Config {
	return a.config
}

// 保存配置到本地
func (a *App) SaveConfig(config *config.Config) bool {

	if config.BaseDir == "" {
		config.BaseDir = a.config.BaseDir
	}

	// 更新前端传来的配置信息
	a.config = config

	// 保存配置文件
	err := a.config.Save()
	if err != nil {
		a.logger.Warnf("保存设置失败: %s", err)
	} else {
		a.logger.Info("保存设置成功")
	}

	return err == nil
}
