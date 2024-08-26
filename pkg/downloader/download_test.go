package downloader

import (
	"context"
	"testing"
)

func TestDown(t *testing.T) {
	url := "https://cdn.yuelili.com/market/vidor/plugins/video-plugin-bilibili/bilibili.exe"

	targetPath := "./test.exe"
	ctx := context.WithoutCancel(context.Background())

	dl, _ := New(ctx, url, targetPath)

	// go dl.Download()

	// time.Sleep(time.Second * 5)
	// dl.Parse()

	segments := []*Pair{
		{3407872, 5242879},
		{6029312, 10485759},
		{11272192, 16090111},
	}
	dl.Recover(segments)

}
