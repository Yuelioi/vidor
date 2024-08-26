package tools

import (
	"context"
	"testing"
)

func TestDown(t *testing.T) {
	url := "https://cdn.yuelili.com/market/vidor/plugins/video-plugin-bilibili/bilibili.exe"

	targetPath := "./test.exe"
	dl, _ := NewDownloader(url, targetPath)
	dl.Download(context.Background())

}
