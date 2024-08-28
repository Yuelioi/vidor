package downloader

import (
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
)

func (d *Downloader) Download() error {
	file, err := prepareOutputFile(d.targetPath, false)
	if err != nil {
		d.State = 4
		return err
	}

	d.out = file
	defer d.out.Close()
	defer d.cancel()

	d.allocateSegments()

	d.State = 1
	return d.start()
}

// 恢复下载
func (d *Downloader) Recover(segments []*Pair) error {
	file, err := prepareOutputFile(d.targetPath, true)
	if err != nil {
		d.State = 4
		return err
	}
	d.out = file
	defer d.out.Close()
	defer d.cancel()

	d.loadSegments(segments)
	d.State = 1
	return d.start()
}

// 暂停
func (d *Downloader) Parse() []Pair {
	fmt.Print("暂停")
	pairs := d.storeWork()
	d.State = 2
	d.cancel()
	return pairs
}

func (d *Downloader) loadSegments(segments []*Pair) {
	for _, seg := range segments {
		pair := newPair(seg.start, seg.end)
		d.segments = append(d.segments, pair)
	}
}

func (d *Downloader) allocateSegments() {
	// 分配下载区间
	for i := int64(0); i < d.batchSize; i++ {
		start := i * d.chunkSize
		end := start + d.chunkSize - 1
		if i == d.batchSize-1 {
			end = d.contentLength - 1
		}

		fmt.Print(start, end, "\n")
		pair := newPair(start, end)
		d.segments = append(d.segments, pair)
	}
}

// 开始
// TODO 错误处理
func (d *Downloader) start() error {
	d.State = 1
	var wg sync.WaitGroup

	if d.supportsChunked {
		for _, seg := range d.segments {
			wg.Add(1)
			go func(pair *pair) {
				defer wg.Done()
				fmt.Printf("seg: %v\n", seg)
				err := d.downloadChunk(seg)
				if err != nil {
					d.mu.Lock()
					defer d.mu.Unlock()
					d.State = 4
					d.Status = err.Error()
				}

			}(seg)
		}
	} else {
		err := d.download()
		if err != nil {
			// 单线程
			d.State = 4
			d.Status = err.Error()
		}
	}

	wg.Wait()
	d.State = 3
	d.Status = "下载完成"
	return nil
}

// 储存工作区
func (d *Downloader) storeWork() []Pair {
	pairs := make([]Pair, 0)
	for _, pair := range d.segments {
		pairs = append(pairs, Pair{
			start: pair.start.Load(),
			end:   pair.end.Load(),
		})
	}
	return pairs
}

// 下载分块
func (d *Downloader) downloadChunk(seg *pair) error {
	chunkStart := seg.start.Load()
	chunkEnd := seg.end.Load()

	req := d.client.R().
		SetHeader("Accept-Ranges", "bytes").
		SetHeader("Range", fmt.Sprintf("bytes=%d-%d", chunkStart, chunkEnd)).
		SetDoNotParseResponse(true)

	resp, err := req.Get(d.url)
	if err != nil {
		log.Println("请求失败:", err)
		return err
	}
	defer resp.RawBody().Close()

	if resp.RawBody() == nil {

		return errors.New("response body is nil")
	}

	buffer := make([]byte, d.bufferSize)

	for {
		select {

		case <-d.ctx.Done():
			fmt.Println("Context canceled")
			return d.ctx.Err()
		default:
			n, err := resp.RawBody().Read(buffer)
			if n > 0 {
				_, writeErr := d.out.WriteAt(buffer[:n], chunkStart)
				if writeErr != nil {

					log.Printf("写入文件失败：%v", writeErr)
					return writeErr
				}
				chunkStart += int64(n)
				seg.start.Add(int64(n))
				d.totalBytesRead.Add(int64(n))
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

// 下载整块
func (d *Downloader) download() error {
	fmt.Println("直接下载")

	resp, err := d.client.R().
		SetDoNotParseResponse(true).Get(d.url)
	if err != nil {
		log.Println("请求失败:", err)
		return err
	}
	defer resp.RawBody().Close()

	if resp.RawBody() == nil {
		return errors.New("response body is nil")
	}

	buffer := make([]byte, d.bufferSize)
	start := int64(0)

	for {
		select {
		case <-d.ctx.Done():
			fmt.Println("Context canceled")
			return d.ctx.Err()
		default:
			n, err := resp.RawBody().Read(buffer)
			if n > 0 {
				_, writeErr := d.out.WriteAt(buffer[:n], start)
				if writeErr != nil {
					log.Printf("写入文件失败：%v", writeErr)
					return writeErr
				}
				start += int64(n)
				d.totalBytesRead.Add(int64(n))
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
