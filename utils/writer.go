package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Yuelioi/vidor/shared"
)

const (
	bufferSize = 1024 * 1024     // 1MB buffer size
	chunkSize  = 5 * 1024 * 1024 // 5MB chunk size
)

type chunkClip struct {
	start int64
	end   int64
}

func autoSetBatchSize(contentLength int64) int64 {
	minBatchSize := int64(2)
	maxBatchSize := int64(20)

	batchSize := int64(math.Sqrt(float64(contentLength) / (1024 * 1024))) // 1MB chunks
	batchSize = int64(math.Max(float64(minBatchSize), float64(math.Min(float64(batchSize), float64(maxBatchSize)))))
	return batchSize
}

func ReqWriter(ctx context.Context, client *http.Client, req *http.Request, part *shared.Part, path string, callback shared.Callback) error {

	req.Header.Set("Accept-Ranges", "bytes")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 创建文件
	out, err := os.Create(path)
	if err != nil {
		log.Println("创建文件失败：", err)
		return err
	}
	defer func() {
		if err := out.Close(); err != nil {
			log.Println("Close out文件失败：", err)
		}
	}()

	contentLength := resp.ContentLength
	batchSize := autoSetBatchSize(contentLength)
	fmt.Printf("batchSize: %v\n", batchSize)

	var totalBytesRead atomic.Int64
	var lastBytesRead int64 = 0

	startTime := time.Now()
	lastTime := time.Now()

	ticker := time.NewTicker(time.Millisecond * 300)

	var downloading = true

	go func() {
		for range ticker.C {
			if downloading {
				now := time.Now()
				elapsed := now.Sub(lastTime)
				lastTime = now

				currentBytesRead := totalBytesRead.Load()
				if lastBytesRead == 0 {
					lastBytesRead = currentBytesRead
				} else {
					currentDownloadSpeed := float64(currentBytesRead-lastBytesRead) / elapsed.Seconds()
					lastBytesRead = currentBytesRead
					part.DownloadSpeed = fmt.Sprintf("%.2f MB/s", currentDownloadSpeed/(1024*1024))
					part.DownloadPercent = int(float64(currentBytesRead) / float64(contentLength) * 100)
					fmt.Printf("part.DownloadSpeed: %v\n", part.DownloadSpeed)
					callback(shared.NoticeData{
						EventName: "updateInfo",
						Message:   part,
					})

				}
			} else {
				finishTime := time.Now()
				elapsed := finishTime.Sub(startTime)
				currentDownloadSpeed := float64(contentLength) / elapsed.Seconds()
				part.DownloadSpeed = fmt.Sprintf("%.2f MB/s", currentDownloadSpeed/(1024*1024))
				ticker.Stop()
				return
			}
		}

	}()

	chunkSize := contentLength / batchSize
	if chunkSize*batchSize < contentLength {
		chunkSize += 1
	}

	var chunks []chunkClip
	for i := int64(0); i <= batchSize-1; i++ {
		start := i * chunkSize
		end := start + chunkSize - 1
		if i == batchSize-1 {
			end = contentLength - 1
		}
		chunks = append(chunks, chunkClip{start, end})
	}

	var wg sync.WaitGroup
	errChan := make(chan error, batchSize)

	for index, chunk := range chunks {
		wg.Add(1)
		go func(idx int, chunk chunkClip) {
			defer wg.Done()
			err := downloadChunk(ctx, client, req, chunk.start, chunk.end, out, &totalBytesRead)
			if err != nil {
				errChan <- err
			}
		}(index, chunk)
	}

	go func() {
		wg.Wait()
		close(errChan)
		downloading = false
	}()

	select {
	case <-ctx.Done():
		// 上下文被取消
		return ctx.Err()
	case err := <-errChan:
		// 下载过程中出错
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadChunk(ctx context.Context, client *http.Client, req *http.Request, start, end int64, out *os.File, totalBytesRead *atomic.Int64) error {

	reqCopy := req.Clone(req.Context())
	reqCopy.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	resp, err := client.Do(reqCopy)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	buffer := make([]byte, bufferSize)

	for {
		select {
		case <-ctx.Done():
			// 上下文被取消
			fmt.Println("ctx 退出")
			return ctx.Err()
		default:
			n, err := resp.Body.Read(buffer)
			if n > 0 {
				_, writeErr := out.WriteAt(buffer[:n], start)
				if writeErr != nil {
					log.Printf("写入文件失败：%v", writeErr)
					return writeErr
				}
				start += int64(n)
				totalBytesRead.Add(int64(n))
			}

			if err != nil {
				if err == io.EOF {
					return nil // 读取完毕，正常退出
				}
				log.Printf("读取响应体失败：%v", err)
				return err // 读取过程中出错，返回错误
			}
		}
	}
}
