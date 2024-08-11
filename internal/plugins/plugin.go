package plugins

import (
	plugin_bilibili "github.com/Yuelioi/vidor/plugins/bilibili"
	"github.com/Yuelioi/vidor/shared"
)

func SystemPlugins(notice shared.Notice) []shared.PluginMeta {
	plugins := make([]shared.PluginMeta, 0)

	bd := plugin_bilibili.Downloader{}
	plugin := shared.PluginMeta{
		Name:   bd.PluginMeta().Name,
		Regexs: bd.PluginMeta().Regexs,
		Type:   "System",
		New:    plugin_bilibili.New,
	}

	plugins = append(plugins, plugin)

	return plugins
}
