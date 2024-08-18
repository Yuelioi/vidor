package app

import (
	"log"
	"os"
	"path/filepath"
)

func ExePath() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	// 获取可执行文件所在目录
	exeDir := filepath.Dir(exePath)

	return exeDir
}
