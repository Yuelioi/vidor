package downloader

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"
	"sync"
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
	contentLength   int64
	url             string // 链接
	targetPath      string // 目标路径
	supportsChunked bool   // 是否支持分段
	out             *os.File
	totalBytesRead  atomic.Int64
	mu              sync.RWMutex
	segments        []*pair
	cancel          context.CancelFunc

	State  int64  // 0尚未下载 1下载中 2下载暂停 3下载完成 4下载出错
	Status string // 状态信息
}

// 创建新的下载器
// isBatch: 是否分块下载(github不支持)
func New(ctx context.Context, url, targetPath string, isBatch bool) (*Downloader, error) {
	newCtx, cancel := context.WithCancel(ctx)

	d := &Downloader{
		url:            url,
		targetPath:     targetPath,
		ctx:            newCtx,
		cancel:         cancel,
		client:         resty.New(),
		bufferSize:     1024 * 256,
		chunkSize:      5 * 1024 * 1024,
		batchSize:      1,
		totalBytesRead: atomic.Int64{},
		mu:             sync.RWMutex{},
		segments:       make([]*pair, 0),
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
		batchSize := autoSetBatchSize(d.contentLength, 1, 5)
		if batchSize == 1 {
			// 只有一段 直接下载
			d.batchSize = 1
			d.supportsChunked = false
		}
	}

	return d, nil
}
