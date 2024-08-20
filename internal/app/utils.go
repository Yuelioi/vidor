package app

import (
	"log"
	"os"
	"path/filepath"

	utils "github.com/Yuelioi/vidor/internal/tools"
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
		logger.Error(err)
		return false
	}
	return true
}
