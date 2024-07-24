package utils

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func saveImage(client *http.Client, url string, path string) (string, error) {
	body, err := doReqBody(client, url)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(path, body, 0644)
	if err != nil {
		fmt.Printf("保存图片失败: %v", err)
		return "", err
	}
	return "/files/" + path, nil
}

// 生成缩略图本地路径
func GetThumbnail(client *http.Client, url string) (string, error) {

	// 获取缓存文件夹路径
	cacheDir := filepath.Join(os.Getenv("LOCALAPPDATA"), "vidor")
	err := os.MkdirAll(cacheDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	// 基于时间戳生成文件名
	timestamp := time.Now().Format("20060102_150405.000")
	fileName := fmt.Sprintf("image_%s.jpg", timestamp)
	filePath := filepath.Join(cacheDir, fileName)

	return saveImage(client, url, filePath)
}

// 直接下载封面
func GetCover(client *http.Client, link string, path string) (string, error) {
	return saveImage(client, link, path)
}
