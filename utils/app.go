package utils

import "github.com/Yuelioi/vidor/shared"

func Downloader2plugin(downloader shared.Downloader, plugin_type string) shared.PluginMeta {
	return shared.PluginMeta{
		Name:   downloader.PluginMeta().Name,
		Regexs: downloader.PluginMeta().Regexs,
		Type:   plugin_type,
		Impl:   downloader,
	}
}
