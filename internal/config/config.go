package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Yuelioi/vidor/internal/models"
	"github.com/Yuelioi/vidor/internal/tools"
)

var defaultSystemConfig = &models.SystemConfig{
	Theme:            "dark",
	ScaleFactor:      16,
	MagicName:        "{{Index}}-{{Title}}",
	DownloadVideo:    true,
	DownloadAudio:    true,
	DownloadSubtitle: true,
	DownloadCombine:  true,
	DownloadLimit:    3,
}

func New(baseDir string) *Config {
	return &Config{
		baseDir:       baseDir,
		SystemConfig:  defaultSystemConfig,
		PluginConfigs: make(map[string]*PluginConfig),
	}
}

// 加载配置
func (c *Config) Load() error {
	if err := tools.MkDirs(c.baseDir); err != nil {
		return err
	}

	configFile := filepath.Join(c.baseDir, "config.json")

	// 检查配置文件是否存在，如果不存在则创建一个空的配置文件
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := c.Save(); err != nil {
			return err
		}
	}

	// 读取配置文件
	configData, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	tempConfig := &Config{
		SystemConfig:  defaultSystemConfig,
		PluginConfigs: make(map[string]*PluginConfig),
	}
	if err := json.Unmarshal(configData, tempConfig); err != nil {
		return err
	}

	c.SystemConfig = tempConfig.SystemConfig
	c.PluginConfigs = tempConfig.PluginConfigs

	// 初始化下载文件夹
	if _, err := os.Stat(c.SystemConfig.DownloadDir); os.IsNotExist(err) {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		c.SystemConfig.DownloadDir = filepath.Join(home, "downloads")
		if err := c.Save(); err != nil {
			return err
		}
	}
	return nil
}

// 保存配置
func (c *Config) Save() error {

	config := map[string]interface{}{
		"system":  c.SystemConfig,
		"plugins": c.PluginConfigs,
	}

	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	configFile := filepath.Join(c.baseDir, "config.json")

	err = os.WriteFile(configFile, configData, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
