package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	logDir        string
	configDir     string
	assetsDir     string
	pluginsDir    string
	SystemConfig  SystemConfig   `json:"system"`
	PluginConfigs []PluginConfig `json:"plugins"`
}

type SystemConfig struct {
	Theme            string `json:"theme"`
	ScaleFactor      int    `json:"scale_factor"`
	ProxyURL         string `json:"proxy_url"`
	UseProxy         bool   `json:"use_proxy"`
	MagicName        string `json:"magic_name"`
	DownloadDir      string `json:"download_dir"`
	DownloadVideo    bool   `json:"download_video"`
	DownloadAudio    bool   `json:"download_audio"`
	DownloadSubtitle bool   `json:"download_subtitle"`
	DownloadCombine  bool   `json:"download_combine"`
	DownloadLimit    int    `json:"download_limit"`
}

type PluginConfig struct {
	ManifestVersion int      `json:"manifest_version"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Author          string   `json:"author"`
	Version         string   `json:"version"`
	URL             string   `json:"url"`
	DocsURL         string   `json:"docs_url"`
	DownloadURL     string   `json:"download_url"`
	Matches         []string `json:"matches"`
	Settings        []string `json:"settings"`
}

func NewConfig() *Config {
	c := &Config{}
	appDir := ExePath()

	c.configDir = filepath.Join(appDir, "configs")
	c.pluginsDir = filepath.Join(appDir, "plugins")
	c.assetsDir = filepath.Join(appDir, "assets")
	c.logDir = filepath.Join(appDir, "logs")

	mkDirs(c.logDir, c.configDir, c.assetsDir, c.pluginsDir)

	c.loadConfig()

	// 加载插件配置

	c.PluginConfigs = []PluginConfig{}
	return c
}

func mkDirs(dirs ...string) {
	for _, dir := range dirs {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Fatal("无法创建文件夹")
		}
	}
}

// 保存配置
func (c *Config) SaveConfig() error {
	// 防止还没有init 前端就监控配置变化？
	if c.configDir == "" {
		return errors.New("还没有初始化")
	}

	config := map[string]interface{}{
		"system":  c.SystemConfig,
		"plugins": c.PluginConfigs,
	}

	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	configFile := filepath.Join(c.configDir, "config.json")

	err = os.WriteFile(configFile, configData, 0644)
	if err != nil {
		return err
	}
	return nil
}

// 加载/创建/初始化配置
func (c *Config) loadConfig() error {
	configFile := filepath.Join(c.configDir, "config.json")

	// 检查配置文件是否存在，如果不存在则创建一个空的配置文件
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Printf("配置文件 '%s' 不存在，将创建一个空的配置文件", configFile)
		c.SystemConfig = SystemConfig{
			Theme:            "dark",
			ScaleFactor:      16,
			MagicName:        "{{Index}}-{{Title}}",
			DownloadVideo:    true,
			DownloadAudio:    true,
			DownloadSubtitle: true,
			DownloadCombine:  true,
			DownloadLimit:    5,
		}
		c.PluginConfigs = []PluginConfig{} // 初始化为空的插件配置

		if err := c.SaveConfig(); err != nil {
			return fmt.Errorf("无法创建配置文件: %v", err)
		}
	}

	configData, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	var config map[string]interface{}
	err = json.Unmarshal(configData, &config)
	if err != nil {
		return err
	}

	if systemConfig, ok := config["system"].(map[string]interface{}); ok {
		if err := mapToStruct(systemConfig, &c.SystemConfig); err != nil {
			return err
		}
	}

	if pluginsConfig, ok := config["plugins"].([]interface{}); ok {
		for _, plugin := range pluginsConfig {
			var pluginConfig PluginConfig
			pluginMap := plugin.(map[string]interface{})
			if err := mapToStruct(pluginMap, &pluginConfig); err != nil {
				return err
			}
			c.PluginConfigs = append(c.PluginConfigs, pluginConfig)
		}
	}

	// 初始化下载文件夹
	if _, err := os.Stat(c.SystemConfig.DownloadDir); os.IsNotExist(err) {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		c.SystemConfig.DownloadDir = filepath.Join(home, "downloads")
		if err := c.SaveConfig(); err != nil {
			return err
		}
	}
	return nil
}

// 辅助函数：将 map 转换为结构体
func mapToStruct(data map[string]interface{}, result interface{}) error {
	configData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(configData, result)
}
