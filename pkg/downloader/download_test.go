package downloader

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestDown(t *testing.T) {
	url := ""

	start := time.Now()
	targetPath := "./test.exe"
	ctx := context.WithoutCancel(context.Background())

	dl, _ := New(ctx, url, targetPath, true)

	err := dl.Download()

	if err != nil {
		fmt.Print(err)
	}
	print(time.Since(start))

	// time.Sleep(time.Second * 5)
	// dl.Parse()

	// segments := []*Pair{
	// 	{3407872, 5242879},
	// 	{6029312, 10485759},
	// 	{11272192, 16090111},
	// }
	// dl.Recover(segments)

}
