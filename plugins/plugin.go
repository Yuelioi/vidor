package plugins

import (
	"github.com/Yuelioi/vidor/plugins/bilibili"
	"github.com/Yuelioi/vidor/shared"
	"github.com/Yuelioi/vidor/utils"
)

func SystemPlugins(notice shared.Notice) []shared.PluginMeta {
	plugins := make([]shared.PluginMeta, 0)

	plugins = append(plugins, utils.Downloader2plugin(&bilibili.Downloader{}, "System"))

	return plugins
}
