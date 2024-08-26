package downloader

import (
	"errors"
	"log"
	"math"
	"os"
	"sync/atomic"
)

func newPair(start, end int64) *pair {
	pair := &pair{
		start: atomic.Int64{},
		end:   atomic.Int64{},
	}
	pair.start.Store(start)
	pair.end.Store(end)

	return pair
}

// 自适应batch
func autoSetBatchSize(contentLength int64) int64 {
	minBatchSize := int64(2)
	maxBatchSize := int64(5)

	batchSize := int64(math.Sqrt(float64(contentLength) / (1024 * 1024))) // 1MB chunks
	batchSize = int64(math.Max(float64(minBatchSize), float64(math.Min(float64(batchSize), float64(maxBatchSize)))))
	return batchSize
}

// 准备文件
func (d *Downloader) prepareOutputFile() error {
	_, err := os.Stat(d.targetPath)
	if err == nil {
		out, err := os.OpenFile(d.targetPath, os.O_RDWR, 0666)
		if err != nil {
			log.Printf("无法打开文件：%v", err)
			return err
		}
		log.Print("使用现有文件")
		d.out = out
	} else if errors.Is(err, os.ErrNotExist) {
		out, err := os.Create(d.targetPath)
		if err != nil {
			log.Printf("无法创建文件：%v", err)
			return err
		}
		d.out = out
	} else {
		// An unexpected error occurred
		log.Printf("无法检查文件状态：%v", err)
		return err
	}
	return nil
}
