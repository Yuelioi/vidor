package plugin

import (
	"github.com/Yuelioi/vidor/internal/config"
	pb "github.com/Yuelioi/vidor/internal/proto"
)

type Plugin struct {
	*config.PluginConfig
	baseDir         string   // 插件所在文件夹
	ManifestVersion int      `json:"manifest_version"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Author          string   `json:"author"`
	Version         string   `json:"version"`
	HomePage        string   `json:"homepage"`
	Color           string   `json:"color"`
	DocsURL         string   `json:"docs_url"`
	DownloadURL     string   `json:"download_url"`
	Matches         []string `json:"matches"`
	Type            string   `json:"type"`     // downloader/other
	Location        string   `json:"location"` // 软件执行文件全名
	State           int      `json:"state"`    // 1.运行中 2.运行中 尚未检测通信结果 3.未启动
	Port            int      `json:"port"`
	PID             int      `json:"pid"`
	Service         pb.DownloadServiceClient
}

func New(baseDir string) *Plugin {
	return &Plugin{
		PluginConfig: &config.PluginConfig{
			Settings: make(map[string]string),
		},
		baseDir: baseDir,
		Matches: make([]string, 0),
	}
}
