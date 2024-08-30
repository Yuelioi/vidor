package config

import (
	"encoding/json"
	"log"
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

func New(baseDir string) (*Config, error) {

	c := &Config{
		baseDir:       baseDir,
		SystemConfig:  defaultSystemConfig,
		PluginConfigs: make(map[string]*PluginConfig),
	}

	err := tools.MkDirs(baseDir)
	if err != nil {
		log.Fatal("无法创建文件夹")
	}

	err = c.load()
	return c, err
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
	log.Println(string(configData)) // log the JSON data

	err = os.WriteFile(configFile, configData, 0644)
	if err != nil {
		return err
	}
	return nil
}

// 加载/创建/初始化配置
func (c *Config) load() error {

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

	config := &Config{
		SystemConfig:  defaultSystemConfig,
		PluginConfigs: map[string]*PluginConfig{},
	}
	err = json.Unmarshal(configData, config)
	if err != nil {
		return err
	}

	c.SystemConfig = config.SystemConfig
	c.PluginConfigs = config.PluginConfigs

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
