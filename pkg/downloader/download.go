package downloader

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"

	"github.com/go-resty/resty/v2"
)

// 下载子进程(内部)
type pair struct {
	start atomic.Int64
	end   atomic.Int64
}

// 下载子进程
type Pair struct {
	start int64
	end   int64
}

type Downloader struct {
	ctx             context.Context
	client          *resty.Client
	bufferSize      int
	chunkSize       int64
	batchSize       int64
	out             *os.File
	totalBytesRead  atomic.Int64
	contentLength   int64
	timeInterval    int64
	state           int64 // 0尚未下载 1下载中 2下载完成 3下载出错
	url             string
	targetPath      string
	supportsChunked bool
	segments        []*pair
	cancel          context.CancelFunc
}

// 创建新的下载器
// isBatch: 是否分块下载(github不支持)
func New(ctx context.Context, url, targetPath string, isBatch bool) (*Downloader, error) {
	newCtx, cancel := context.WithCancel(ctx)

	d := &Downloader{
		ctx:          newCtx,
		client:       resty.New(),
		bufferSize:   1024 * 256,
		chunkSize:    5 * 1024 * 1024,
		batchSize:    1,
		state:        0,
		timeInterval: 2333,
		url:          url,
		targetPath:   targetPath,
		segments:     make([]*pair, 0),
		cancel:       cancel,
	}

	resp, err := d.client.R().
		SetHeader("Accept-Ranges", "bytes").
		SetHeader("Range", "bytes=0-").
		SetDoNotParseResponse(true).
		Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == http.StatusPartialContent && resp.Header().Get("Accept-Ranges") == "bytes" {
		d.supportsChunked = true
		//
	} else {
		// 服务器不支持分块下载
		d.supportsChunked = false
		fmt.Print("不支持分块下载")
		return d, nil
	}

	// 获取目标长度
	contentLengthStr := resp.Header().Get("Content-Length")
	if contentLengthStr != "" {
		d.contentLength, err = strconv.ParseInt(contentLengthStr, 10, 64)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Content-Length header is missing")
	}

	if isBatch {
		d.batchSize = autoSetBatchSize(d.contentLength, 2, 5)
	}

	fmt.Printf("Content-Length %+v\n", d.contentLength)

	return d, nil
}
