package config

import (
	"github.com/Yuelioi/vidor/internal/models"
)

type Config struct {
	baseDir       string                   // 配置所在文件夹
	SystemConfig  *models.SystemConfig     `json:"system"` // 系统配置
	PluginConfigs map[string]*PluginConfig `json:"plugin"` // 插件配置
}

type PluginConfig struct {
	ID       string            `json:"id"`       // 插件ID
	Enable   bool              `json:"enable"`   // 插件是否开机启动
	Settings map[string]string `json:"settings"` // 插件设置
}
