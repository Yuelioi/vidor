package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

// 获取合集信息
func processSeasonData(data biliPlayListData, aid int, bvid_input string) biliBaseParams {
	episodes := data.UgcSeason.Sections[0].Episodes
	index := 0
	for epid, epi := range episodes {
		if epi.Aid == aid || epi.Bvid == bvid_input {
			index = epid
		}
	}
	workDirname := SanitizeFileName(data.UgcSeason.Title)
	coverName := fmt.Sprintf("%02d_%s.jpg", index+1, SanitizeFileName(workDirname))

	return biliBaseParams{
		index:    index,
		cid:      episodes[index].CID,
		bvid:     data.UgcSeason.Sections[0].Episodes[index].Bvid,
		title:    episodes[index].Title,
		author:   data.Owner.Name,
		pubDate:  episodes[index].Arc.Pubdate,
		duration: episodes[index].Arc.Duration,

		coverURL:     data.UgcSeason.Cover,
		coverName:    coverName,
		thumbnailURL: episodes[index].Arc.Pic,
	}
}

// 获取普通分p信息
func processPagesData(bpi biliPlayListData, index int) biliBaseParams {
	var title, thumbnailURL string

	if len(bpi.Pages) == 1 {
		title = bpi.Title
	} else {
		title = bpi.Pages[index].Title
	}

	if bpi.Pages[index].Thumbnail != "" {
		thumbnailURL = bpi.Pages[index].Thumbnail
	} else {
		thumbnailURL = bpi.Pic
	}

	workDirname := bpi.Title

	return biliBaseParams{
		index:        index,
		cid:          bpi.Pages[index].CID,
		bvid:         bpi.BVID,
		title:        SanitizeFileName(title),
		author:       bpi.Owner.Name,
		coverURL:     bpi.Pic,
		coverName:    SanitizeFileName(workDirname) + ".jpg",
		thumbnailURL: thumbnailURL,
	}
}

// 音频取最大的, 反正不大
func getTargetAudio(audios []Audio) *Audio {

	if len(audios) == 0 {
		return nil
	}
	sort.Slice(audios, func(i, j int) bool {
		return audios[i].Bandwidth > audios[j].Bandwidth
	})

	return &audios[0]
}

// 如果有 直接用, 没有就找比它高1级的
// TODO 视频编码选择
func getTargetVideo(targeID int, videos []Video) *Video {
	// 倒序, 向上取
	sort.Slice(videos, func(i, j int) bool {
		return videos[i].ID > videos[j].ID
	})

	var closestVideo *Video
	for _, video := range videos {
		if video.ID <= targeID {
			return &video
		}
		closestVideo = &video
	}

	return closestVideo
}

// 获取分p, 减1 以拿到索引
func extractIndex(url string) (int, error) {
	re := regexp.MustCompile(`p=(\d+)`)
	match := re.FindStringSubmatch(url)
	if len(match) > 1 {
		i, err := strconv.Atoi(match[1])
		return i - 1, err

	}
	return 0, fmt.Errorf("no index found in URL")
}

// 获取aid或者bvid
func extractAidBvid(link string) (aid int, bvid string) {
	aidRegex := regexp.MustCompile(`av(\d+)`)
	bvidRegex := regexp.MustCompile(`BV\w+`)

	aidMatches := aidRegex.FindStringSubmatch(link)
	if len(aidMatches) > 1 {
		aid, _ = strconv.Atoi(aidMatches[1])
	}

	bvidMatches := bvidRegex.FindStringSubmatch(link)
	if len(bvidMatches) > 0 {
		bvid = bvidMatches[0]
	}
	return
}

// 合并音频与视频
func CombineAV(ctx context.Context, ffmpegPath string, input_v, input_a, output_v, logFile string) (err error) {

	input := []*ffmpeg_go.Stream{ffmpeg_go.Input(input_v), ffmpeg_go.Input(input_a)}
	out := ffmpeg_go.OutputContext(ctx, input, output_v, ffmpeg_go.KwArgs{"c:v": "copy", "c:a": "aac"})

	_, err = os.Stat(ffmpegPath)

	if err == nil {
		out = out.SetFfmpegPath(ffmpegPath)
	}

	logDir := filepath.Dir(logFile)

	// 检查目录是否存在
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		// 创建目录，使用 0755 权限
		if err := os.MkdirAll(logDir, 0755); err != nil {
			fmt.Println("Error creating directory:", err)
			return err
		}
	}

	err = out.OverWriteOutput().WithOutput().Run()
	// logAdapter, err := createLogAdapter(logFile)
	// if err != nil {
	// 	return err
	// }
	// defer logAdapter.file.Close()

	// cmd := out.OverWriteOutput().WithOutput(logAdapter, logAdapter).Compile()

	// // TODO关闭cmd弹窗
	// // cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	// err = cmd.Run()
	return err
}

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
