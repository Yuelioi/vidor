package tools

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func MkDirs(dirs ...string) error {
	for _, dir := range dirs {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func CleanDir(dirs ...string) error {
	for _, dir := range dirs {
		err := os.RemoveAll(dir)
		if err != nil {
			return err
		}
	}
	return nil
}

// 获取当前exe所在目录
func ExeDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Dir(exePath), nil
}

// SanitizeFileName takes an input string and returns a sanitized file name
func SanitizeFileName(input string) string {
	re := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)
	sanitized := re.ReplaceAllString(input, "_")

	sanitized = strings.TrimSpace(sanitized)
	sanitized = strings.Trim(sanitized, ".")

	if len(sanitized) > 255 {
		sanitized = sanitized[:255]
	}

	return sanitized
}

func ExtractFileNameFromUrl(url string) string {
	// 删除？之后的内容
	if i := strings.LastIndex(url, "?"); i >= 0 {
		url = url[:i]
	}

	// 删除最后一个/
	url = strings.TrimSuffix(url, "/")

	splitUrl := strings.Split(url, "/")
	name := splitUrl[len(splitUrl)-1]
	return name
}

func randomString(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// 生成魔法名称
func MagicName(template, workDirname, title string, index int) string {
	template = strings.ReplaceAll(template, "{{Title}}", title)
	template = strings.ReplaceAll(template, "{{Index}}", fmt.Sprintf("%0.3d", index))
	template = strings.ReplaceAll(template, "{{RndInt}}", fmt.Sprintf("%0.3d", rand.Intn(999)))
	template = strings.ReplaceAll(template, "{{RndChr}}", randomString(3))
	return template
}

// 创建文件夹
func CreateDirs(dirs []string) (err error) {
	for _, dir := range dirs {
		if _, err = os.Stat(dir); os.IsNotExist(err) {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return
			}
		}
	}
	return
}

// 清理文件夹
func ClearDirs(dirs []string) (err error) {
	for _, dir := range dirs {
		err = os.Remove(dir)
		if err != nil {
			return
		}
	}
	return
}
