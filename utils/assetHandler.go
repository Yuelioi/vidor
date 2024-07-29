package utils

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func saveImage(client *http.Client, url string, path string) (string, error) {

	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", fmt.Errorf("资产处理 创建文件夹失败%s", err)
	}

	body, err := doReqBody(client, url)
	if err != nil {
		return "", fmt.Errorf("资产处理 请求网络资源失败 %s", err)
	}
	err = os.WriteFile(path, body, 0644)
	if err != nil {
		return "", fmt.Errorf("资产处理 写入图片失败 %s", err)
	}
	return "/files/" + path, nil
}

// 生成缩略图本地路径
func GetThumbnail(client *http.Client, url, filePath string) (string, error) {
	return saveImage(client, url, filePath)
}

// 直接下载封面
func GetCover(client *http.Client, link, path string) (string, error) {
	return saveImage(client, link, path)
}
