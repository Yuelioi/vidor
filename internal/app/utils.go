package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	utils "github.com/Yuelioi/vidor/internal/tools"
	"github.com/go-resty/resty/v2"
)

// 获取当前exe所在目录
func ExePath() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Dir(exePath)
}

// 检查是否存在FFmpeg
func CheckFFmpeg(target string) bool {
	if err := utils.SetFFmpegPath(target); err != nil {
		return false
	}
	return true
}

const (
	bufferSize = 1024 * 256      // 500kb buffer size
	chunkSize  = 5 * 1024 * 1024 // 5MB chunk size
)

func autoSetBatchSize(contentLength int64) int64 {
	minBatchSize := int64(2)
	maxBatchSize := int64(5)

	batchSize := int64(math.Sqrt(float64(contentLength) / (1024 * 1024))) // 1MB chunks
	batchSize = int64(math.Max(float64(minBatchSize), float64(math.Min(float64(batchSize), float64(maxBatchSize)))))
	return batchSize
}

func Down(url string) error {
	client := &resty.Client{}
	req := client.R().SetHeader("Accept-Ranges", "bytes").
		SetHeader("Range", "bytes=0-").SetDoNotParseResponse(true)
	resp, err := req.Get(url)
	if err != nil {
		return nil
	}
	contentLengthStr := resp.Header().Get("Content-Length")
	var contentLength int64
	if contentLengthStr != "" {
		contentLength, err = strconv.ParseInt(contentLengthStr, 10, 64)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Content-Length header is missing")
	}

	fmt.Printf("contentLength: %v\n", contentLength)

	return nil
}

func (app *App) download(ctx context.Context, id string, url string, contentLength int64, tempPath string) error {

	batchSize := autoSetBatchSize(contentLength)
	chunkSize := contentLength / batchSize
	if chunkSize*batchSize < contentLength {
		chunkSize += 1
	}

	out, err := os.Create(tempPath)
	if err != nil {
		log.Printf("无法创建文件：%v", err)
		return err
	}
	defer out.Close()

	var wg sync.WaitGroup
	var totalBytesRead atomic.Int64

	var finished bool

	timeInterval := 333

	ticker := time.NewTicker(time.Duration(timeInterval) * time.Millisecond)

	go func() {
		var previousBytesRead int64
		defer ticker.Stop()

		for range ticker.C {
			if !finished {
				currentBytesRead := totalBytesRead.Load()
				bytesRead := currentBytesRead - previousBytesRead
				previousBytesRead = currentBytesRead

				speedByte := float64(bytesRead) // Speed in B/s

				speed := fmt.Sprintf("%.2f MB/s", speedByte*1000/(1024*1024*float64(timeInterval)))
				runtime.EventsEmit(app.ctx, "plugin-download", speed)
			} else {
				ticker.Stop()
				return
			}

		}

	}()

	for i := int64(0); i < batchSize; i++ {
		start := i * chunkSize
		end := start + chunkSize - 1
		if i == batchSize-1 {
			end = contentLength - 1
		}

		wg.Add(1)
		go func(chunkStart, chunkEnd int64) {
			defer wg.Done()
			downloadChunk(ctx, id, url, chunkStart, chunkEnd, out, &totalBytesRead)
		}(start, end)
	}
	wg.Wait()
	finished = true

	return nil
}

func downloadChunk(ctx context.Context, id string, url string, chunkStart, chunkEnd int64, out *os.File, totalBytesRead *atomic.Int64) error {
	var c resty.Client

	req := c.R().
		SetHeader("Accept-Ranges", "bytes").
		SetHeader("Range", fmt.Sprintf("bytes=%d-%d", chunkStart, chunkEnd)).
		SetDoNotParseResponse(true)

	resp, err := req.Get(url)
	if err != nil {
		log.Println("请求失败:", err)
		return err
	}
	defer resp.RawBody().Close()

	buffer := make([]byte, bufferSize)

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
