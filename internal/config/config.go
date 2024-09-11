package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Yuelioi/vidor/internal/tools"
)

var defaultSystemConfig = &Config{
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

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return &Config{
		BaseDir:          baseDir,
		Theme:            defaultSystemConfig.Theme,
		ScaleFactor:      defaultSystemConfig.ScaleFactor,
		MagicName:        defaultSystemConfig.MagicName,
		DownloadVideo:    defaultSystemConfig.DownloadVideo,
		DownloadAudio:    defaultSystemConfig.DownloadAudio,
		DownloadSubtitle: defaultSystemConfig.DownloadSubtitle,
		DownloadCombine:  defaultSystemConfig.DownloadCombine,
		DownloadLimit:    defaultSystemConfig.DownloadLimit,
		DownloadDir:      filepath.Join(home, "downloads"),
	}, nil
}

// 加载配置
func (c *Config) Load() error {
	if err := tools.MkDirs(c.BaseDir); err != nil {
		return err
	}

	configFile := filepath.Join(c.BaseDir, "config.json")

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

	tempConfig := &Config{}
	if err := json.Unmarshal(configData, tempConfig); err != nil {
		return err
	}

	// 初始化下载文件夹
	if _, err := os.Stat(c.DownloadDir); os.IsNotExist(err) {

		if err := c.Save(); err != nil {
			return err
		}
	}
	return nil
}

// 保存配置
func (c *Config) Save() error {

	configData, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	configFile := filepath.Join(c.BaseDir, "config.json")

	err = os.WriteFile(configFile, configData, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
