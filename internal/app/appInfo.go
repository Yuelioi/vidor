package app

import "github.com/Yuelioi/vidor/internal/globals"

// 软件基础信息 aa
type AppInfo struct {
	name    string
	version string

	LogDir     string
	ConfigDir  string
	AssetsDir  string
	PluginsDir string
}

func NewAppInfo() AppInfo {
	return AppInfo{
		name:    globals.Name,
		version: globals.Version,
	}
}
