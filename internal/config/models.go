package config

import (
	"github.com/Yuelioi/vidor/internal/models"
)

type Config struct {
	baseDir       string
	SystemConfig  *models.SystemConfig     `json:"system"`
	PluginConfigs map[string]*PluginConfig `json:"plugin"`
}

type PluginConfig struct {
	ID       string            `json:"id"`
	Enable   bool              `json:"enable"` // 建立连接 (Run)
	Settings map[string]string `json:"settings"`
}
