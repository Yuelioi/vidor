package downloader

// func (d *Downloader) monitorDownloadSpeed() {
// 	ticker := time.NewTicker(time.Duration(d.timeInterval) * time.Millisecond)
// 	var previousBytesRead int64
// 	defer ticker.Stop()

// 	for range ticker.C {
// 		select {
// 		case <-d.ctx.Done():
// 			fmt.Println("Context canceled")
// 			return
// 		default:
// 			currentBytesRead := d.totalBytesRead.Load()
// 			bytesRead := currentBytesRead - previousBytesRead
// 			previousBytesRead = currentBytesRead

// 			speedByte := float64(bytesRead)
// 			speed := fmt.Sprintf("%.2f MB/s", speedByte*1000/(1024*1024*float64(d.timeInterval)))
// 			fmt.Println(speed)
// 			d.Status = speed
// 		}
// 	}
// }
