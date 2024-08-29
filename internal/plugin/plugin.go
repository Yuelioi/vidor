package plugin

import (
	"github.com/Yuelioi/vidor/internal/config"
	pb "github.com/Yuelioi/vidor/internal/proto"
)

type Plugin struct {
	*config.PluginConfig
	baseDir         string                   // 插件所在文件夹
	ManifestVersion int                      `json:"manifest_version"`
	Name            string                   `json:"name"`
	Description     string                   `json:"description"`
	Author          string                   `json:"author"`
	Version         string                   `json:"version"`
	HomePage        string                   `json:"homepage"`
	Color           string                   `json:"color"`
	DocsURL         string                   `json:"docs_url"`
	DownloadURLs    []string                 `json:"download_urls"`
	Matches         []string                 `json:"matches"`
	Categories      []string                 `json:"categories"`
	Tags            []string                 `json:"tags"`
	Executable      string                   `json:"executable"` // 软件执行文件全名
	State           int                      `json:"state"`      // 1.运行中 2.运行中 尚未检测通信结果 3.未启动
	Status          string                   `json:"status"`
	Port            int                      `json:"port"`
	PID             int                      `json:"pid"`
	Service         pb.DownloadServiceClient `json:"-"`
}

func New(baseDir string) *Plugin {
	return &Plugin{
		PluginConfig: &config.PluginConfig{
			Settings: make(map[string]string),
		},
		baseDir:      baseDir,
		DownloadURLs: make([]string, 0),
		Categories:   make([]string, 0),
		Tags:         make([]string, 0),
		Matches:      make([]string, 0),
	}
}
