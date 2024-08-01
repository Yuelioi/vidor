package utils

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
)

func MaxInt(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func MinInt(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func ClampInt(src, min, max int) int {
	return MinInt(MaxInt(src, min), max)
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
	template = strings.ReplaceAll(template, "{{WorkDir}}", workDirname)
	template = strings.ReplaceAll(template, "{{Title}}", title)
	template = strings.ReplaceAll(template, "{{Index}}", fmt.Sprintf("%0.3d", index))
	template = strings.ReplaceAll(template, "{{RndInt}}", fmt.Sprintf("%0.3d", rand.Intn(999)))
	template = strings.ReplaceAll(template, "{{RndChr}}", randomString(3))
	return template
}

// // 向上取画质代码
// //
// //	sortOrder?
// func GetQualityID(label string, qualities []shared.StreamQuality) (int, error) {

// 	if len(qualities) == 0 {
// 		return 0, errors.New("视频质量列表为空")
// 	}

// 	qualitiesCopy := make([]shared.StreamQuality, len(qualities))
// 	copy(qualitiesCopy, qualities)

// 	sort.Slice(qualitiesCopy, func(i, j int) bool {
// 		return qualitiesCopy[i].ID > qualitiesCopy[j].ID
// 	})

// 	for _, q := range qualitiesCopy {
// 		if q.Label == label {
// 			return q.ID, nil
// 		}
// 	}
// 	return qualitiesCopy[0].ID, nil
// }

// func GetQualityLabel(id int, qualities []shared.StreamQuality) (string, error) {

// 	if len(qualities) == 0 {
// 		return "", errors.New("视频质量列表为空")
// 	}

// 	// 创建 qualities 的副本
// 	qualitiesCopy := make([]shared.StreamQuality, len(qualities))
// 	copy(qualitiesCopy, qualities)

// 	sort.Slice(qualitiesCopy, func(i, j int) bool {
// 		return qualitiesCopy[i].ID > qualitiesCopy[j].ID
// 	})

// 	for _, q := range qualitiesCopy {
// 		if q.ID == id {
// 			return q.Label, nil
// 		}
// 	}
// 	return qualitiesCopy[0].Label, nil
// }

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
