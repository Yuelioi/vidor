package app

import (
	"github.com/Yuelioi/vidor/internal/plugin"
)

type Cache struct {
	downloader *plugin.DownloadPlugin
}

func NewCache() *Cache {
	return &Cache{}
}

// 下载器缓存
func (c *Cache) Downloader() *plugin.DownloadPlugin {
	return c.downloader
}
func (c *Cache) SetDownloader(p *plugin.DownloadPlugin) {
	c.downloader = p
}

// 插件列表缓存
