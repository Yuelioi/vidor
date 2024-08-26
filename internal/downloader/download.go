package tools

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-resty/resty/v2"
)

type Downloader struct {
	client        *resty.Client
	bufferSize    int
	chunkSize     int64
	batchSize     int64
	contentLength int64
	timeInterval  int64
	state         int64 // 0尚未下载 1下载中 2下载完成 3下载出错
	url           string
	targetPath    string
}

func NewDownloader(url, targetPath string) (*Downloader, error) {
	d := &Downloader{
		client:       resty.New(),
		bufferSize:   1024 * 256,
		chunkSize:    5 * 1024 * 1024,
		state:        0,
		timeInterval: 2333,
		url:          url,
		targetPath:   targetPath,
	}

	resp, err := d.client.R().
		SetHeader("Accept-Ranges", "bytes").
		SetHeader("Range", "bytes=0-").
		SetDoNotParseResponse(true).
		Get(url)

	if err != nil {
		return nil, err
	}

	contentLengthStr := resp.Header().Get("Content-Length")
	if contentLengthStr != "" {
		d.contentLength, err = strconv.ParseInt(contentLengthStr, 10, 64)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Content-Length header is missing")
	}

	d.batchSize = autoSetBatchSize(d.contentLength)
	return d, nil
}

func (d *Downloader) Download(ctx context.Context) error {
	out, err := os.Create(d.targetPath)
	if err != nil {
		log.Printf("无法创建文件：%v", err)
		return err
	}
	defer out.Close()

	var wg sync.WaitGroup
	var totalBytesRead atomic.Int64

	var finished bool

	ticker := time.NewTicker(time.Duration(d.timeInterval) * time.Millisecond)

	go func() {
		var previousBytesRead int64
		defer ticker.Stop()

		for range ticker.C {
			if !finished {
				currentBytesRead := totalBytesRead.Load()
				bytesRead := currentBytesRead - previousBytesRead
				previousBytesRead = currentBytesRead

				speedByte := float64(bytesRead)
				speed := fmt.Sprintf("%.2f MB/s", speedByte*1000/(1024*1024*float64(d.timeInterval)))
				fmt.Println(speed)
			} else {
				ticker.Stop()
				return
			}
		}
	}()

	for i := int64(0); i < d.batchSize; i++ {
		start := i * d.chunkSize
		end := start + d.chunkSize - 1
		if i == d.batchSize-1 {
			end = d.contentLength - 1
		}

		wg.Add(1)
		go func(chunkStart, chunkEnd int64) {
			defer wg.Done()
			err := d.downloadChunk(ctx, chunkStart, chunkEnd, out, &totalBytesRead)
			if err != nil {
				log.Println("请求失败:", err)
			}
		}(start, end)
	}
	wg.Wait()
	finished = true

	return nil
}

func (d *Downloader) downloadChunk(ctx context.Context, chunkStart, chunkEnd int64, out *os.File, totalBytesRead *atomic.Int64) error {
	req := d.client.R().
		SetHeader("Accept-Ranges", "bytes").
		SetHeader("Range", fmt.Sprintf("bytes=%d-%d", chunkStart, chunkEnd)).
		SetDoNotParseResponse(true)

	resp, err := req.Get(d.url)
	if err != nil {
		log.Println("请求失败:", err)
		return err
	}
	if resp.RawBody() == nil {
		return errors.New("response body is nil")
	}
	defer resp.RawBody().Close()

	buffer := make([]byte, d.bufferSize)

	for {
		select {

		case <-ctx.Done():
			fmt.Println("Context canceled")
			return ctx.Err()
		default:
			n, err := io.ReadFull(resp.RawBody(), buffer)
			if n > 0 {
				_, writeErr := out.WriteAt(buffer[:n], chunkStart)
				if writeErr != nil {
					log.Printf("写入文件失败：%v", writeErr)
					return writeErr
				}
				chunkStart += int64(n)
				totalBytesRead.Add(int64(n))
			}

			if err != nil {
				if err == io.EOF {
					return nil // 读取完毕，正常退出
				}
				return err // 读取过程中出错，返回错误
			}
		}
	}
}

func autoSetBatchSize(contentLength int64) int64 {
	minBatchSize := int64(2)
	maxBatchSize := int64(5)

	batchSize := int64(math.Sqrt(float64(contentLength) / (1024 * 1024))) // 1MB chunks
	batchSize = int64(math.Max(float64(minBatchSize), float64(math.Min(float64(batchSize), float64(maxBatchSize)))))
	return batchSize
}
