package downloader

import (
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
	batchSize := int64(math.Sqrt(float64(contentLength) / (1024 * 1024 * 50))) // 1MB  *5 chunks
	batchSize = max(minBatchSize, min(batchSize, maxBatchSize))
	return batchSize
}

// 准备文件
//
// recover: 是否继续原文件下载
func prepareOutputFile(targetPath string, recover bool) (*os.File, error) {
	var f *os.File
	var err error

	if recover {
		f, err = os.OpenFile(targetPath, os.O_RDWR, 0644)
		if os.IsNotExist(err) {
			f, err = os.Create(targetPath)
		}
	} else {
		f, err = os.Create(targetPath)
	}

	if err != nil {
		log.Printf("无法创建或打开文件：%v", err)
		return nil, err
	}

	if recover {
		log.Print("使用现有文件")
	}

	return f, nil
}
