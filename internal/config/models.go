package config

import (
	"github.com/Yuelioi/vidor/internal/models"
)

type Config struct {
	baseDir       string
	SystemConfig  *models.SystemConfig           `json:"system"`
	PluginConfigs map[string]models.PluginConfig `json:"plugins"`
	Test          *models.SystemConfig           `json:"test"`
}
