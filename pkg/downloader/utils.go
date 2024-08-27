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
func autoSetBatchSize(contentLength int64, minBatchSize, maxBatchSize int64) int64 {
	batchSize := int64(math.Sqrt(float64(contentLength) / (1024 * 1024))) // 1MB chunks
	batchSize = max(minBatchSize, min(batchSize, maxBatchSize))
	return batchSize
}

// 准备文件
func prepareOutputFile(targetPath string) (*os.File, error) {
	_, err := os.Stat(targetPath)
	if err == nil {
		out, err := os.OpenFile(targetPath, os.O_RDWR, 0666)
		if err != nil {
			log.Printf("无法打开文件：%v", err)
			return out, err
		}
		log.Print("使用现有文件")
		return out, nil
	} else if errors.Is(err, os.ErrNotExist) {
		out, err := os.Create(targetPath)
		if err != nil {
			log.Printf("无法创建文件：%v", err)
			return out, err
		}
		return out, nil
	} else {
		// An unexpected error occurred
		log.Printf("无法检查文件状态：%v", err)
		return nil, err
	}
}
