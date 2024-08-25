package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"google.golang.org/grpc/metadata"
)

var defaultSystemConfig = &SystemConfig{
	Theme:            "dark",
	ScaleFactor:      16,
	MagicName:        "{{Index}}-{{Title}}",
	DownloadVideo:    true,
	DownloadAudio:    true,
	DownloadSubtitle: true,
	DownloadCombine:  true,
	DownloadLimit:    3,
}

type Config struct {
	logDir        string
	configDir     string
	assetsDir     string
	pluginsDir    string
	SystemConfig  *SystemConfig            `json:"system"`
	PluginConfigs map[string]*PluginConfig `json:"plugins"`
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
	ID       string            `json:"id"`
	Enable   bool              `json:"enable"` // 建立连接 (Run)
	Settings map[string]string `json:"settings"`
}

func NewConfig() (*Config, error) {
	appDir := ExePath()
	fmt.Printf("appDir: %v\n", appDir)

	c := &Config{
		configDir:     filepath.Join(appDir, "configs"),
		pluginsDir:    filepath.Join(appDir, "plugins"),
		assetsDir:     filepath.Join(appDir, "assets"),
		logDir:        filepath.Join(appDir, "logs"),
		SystemConfig:  defaultSystemConfig,
		PluginConfigs: make(map[string]*PluginConfig),
	}

	err := mkDirs(c.logDir, c.configDir, c.assetsDir, c.pluginsDir)
	if err != nil {
		log.Fatal("无法创建文件夹")
	}

	err = c.load()
	return c, err
}

func mkDirs(dirs ...string) error {
	for _, dir := range dirs {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
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

	configFile := filepath.Join(c.configDir, "config.json")

	err = os.WriteFile(configFile, configData, 0644)
	if err != nil {
		return err
	}
	return nil
}

// 加载/创建/初始化配置
func (c *Config) load() error {
	configFile := filepath.Join(c.configDir, "config.json")

	// 检查配置文件是否存在，如果不存在则创建一个空的配置文件
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := c.Save(); err != nil {
			return fileOrDirCreationFailed
		}
	}

	// 读取配置文件
	configData, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	config := &Config{}
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

// 注入系统配置元数据(插件元数据请在插件下注入)
func (c *Config) injectMetadata(ctx context.Context) context.Context {
	v := reflect.ValueOf(c.SystemConfig).Elem()
	t := v.Type()

	// Iterate over all fields in the struct
	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name
		fieldValue := v.Field(i).Interface()

		// Convert field value to string (consider different types)
		var valueStr string
		switch val := fieldValue.(type) {
		case string:
			valueStr = val
		case int:
			valueStr = fmt.Sprintf("%d", val)
		case bool:
			valueStr = fmt.Sprintf("%t", val)
		default:
			valueStr = fmt.Sprintf("%v", val)
		}

		ctx = metadata.AppendToOutgoingContext(ctx, "app."+strings.ToLower(fieldName), valueStr)
	}

	return ctx
}
