package app

// APP 缓存

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
