package plugins

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/Yuelioi/vidor/shared"
	"github.com/Yuelioi/vidor/utils"

	"github.com/kkdai/youtube/v2"
)

type YouTubeDownloader struct {
	Client  *youtube.Client
	TaskDir string
}

func NewYTBDownloader() *YouTubeDownloader {
	return &YouTubeDownloader{}
}

func (yd *YouTubeDownloader) GetClient(proxyURL string, useProxy bool) (*youtube.Client, error) {
	client, err := utils.GetClient(proxyURL, useProxy)
	if err != nil {
		return nil, err
	}

	return &youtube.Client{
		HTTPClient:  client,
		MaxRoutines: 12,
	}, nil

}

func (yd *YouTubeDownloader) ShowInfo(link string, config shared.Config) (*shared.PlaylistInfo, error) {

	ytbClient, err := yd.GetClient(config.ProxyURL, config.UseProxy)

	if err != nil {
		return nil, err
	}

	yd.Client = ytbClient
	var pi shared.PlaylistInfo

	pi.Url = link
	pi.StreamInfos = make([]shared.StreamInfo, 0)

	cacheThumbnail := ""

	if isPlaylist(link) {
		playlist, err := yd.Client.GetPlaylist(link)
		if err != nil {
			return nil, err
		}

		if playlist == nil {
			return nil, errors.New("无法获取播放列表")
		}

		pi.WorkDirName = playlist.Title
		yd.TaskDir = playlist.Title
		pi.Author = playlist.Author
		cacheThumbnail = getBestThumbnail(playlist.Videos[0].Thumbnails).URL

		var wg sync.WaitGroup

		parts := []struct {
			Index     int
			Url       string
			Title     string
			MaxHeight int
		}{}

		for index, v := range playlist.Videos {
			wg.Add(1)
			go func(index int, v *youtube.PlaylistEntry) {
				defer wg.Done()

				video, err := yd.Client.VideoFromPlaylistEntry(v)

				bestHeight := getBestHighFormatHeight(video.Formats)

				if err != nil {
					fmt.Printf("Error fetching video from playlist entry: %v", err)
					return
				}

				parts = append(parts, struct {
					Index     int
					Url       string
					Title     string
					MaxHeight int
				}{
					Index:     index,
					Url:       fmt.Sprintf("https://www.youtube.com/watch?v=%s", v.ID),
					Title:     v.Title,
					MaxHeight: bestHeight,
				})

			}(index, v)
		}
		wg.Wait()

		maxHeight := 0
		for _, part := range parts {
			if part.MaxHeight > maxHeight {
				maxHeight = part.MaxHeight
			}
		}

		sort.Slice(parts, func(i, j int) bool {
			return parts[i].Index < parts[j].Index
		})

		for _, part := range parts {
			pi.StreamInfos = append(pi.StreamInfos, shared.StreamInfo{
				TaskID: part.Title,
				// Url:   part.Url,
				// Title: part.Title,
			})
		}
		//todo
		// pi.Qualities = utils.GetQualities(maxHeight)

	} else {
		video, err := yd.getVideoDataByID(extractID(link))

		for _, format := range video.Formats {
			fmt.Printf("ItagNo: %d\n", format.ItagNo)
			fmt.Printf("format.Quality: %v\n", format.Quality)
			fmt.Printf("format.QualityLabel: %v\n", format.QualityLabel)
		}

		if err != nil {
			return nil, err
		}
		pi.Author = video.Author
		pi.WorkDirName = video.Title
		pi.StreamInfos = append(pi.StreamInfos, shared.StreamInfo{})

		yd.TaskDir = utils.SanitizeFileName(video.Title)
		cacheThumbnail = getBestThumbnail(video.Thumbnails).URL
	}

	img, err := utils.GetThumbnail(yd.Client.HTTPClient, cacheThumbnail, "")
	if err != nil {
		return nil, err
	}
	pi.Cover = img

	return &pi, nil
}

func (yd *YouTubeDownloader) GetMeta(ctx context.Context, part *shared.Part, callback shared.Callback) error {

	video, err := yd.getVideoDataByID(extractID(part.Url))
	if err != nil {
		return err
	}
	// 下载视频应该关闭超时

	part.Author = video.Author
	part.Title = video.Title
	part.Thumbnail = getBestThumbnail(video.Thumbnails).URL

	part.DownloadDir = filepath.Join(part.DownloadDir, yd.TaskDir)
	if err := utils.CreateDirs([]string{part.DownloadDir}); err != nil {
		return err
	}

	taskImagePath := filepath.Join(part.DownloadDir, utils.SanitizeFileName(video.Title)+".jpg")

	if _, err = os.Stat(taskImagePath); os.IsNotExist(err) {
		_, err := utils.GetCover(yd.Client.HTTPClient, part.Thumbnail, taskImagePath)
		if err != nil {
			return fmt.Errorf("缓存缩略图失败: %v", err.Error())
		}
	}

	part.Status = "获取封面"

	targetHeight, _ := utils.GetQualityID(part.Quality, []shared.StreamQuality{})
	part.Quality, _ = utils.GetQualityLabel(targetHeight, []shared.StreamQuality{})

	part.Status = "获取元数据"

	format_v := getTargetYtbVideo(targetHeight, video.Formats)
	part.Quality = format_v.QualityLabel
	// filePureName := utils.SanitizeFileName(part.Title)

	// input_v := filepath.Join(part.DownloadDir, fmt.Sprintf("%s_temp.mp4", filePureName))
	// input_a := filepath.Join(part.DownloadDir, fmt.Sprintf("%s_temp.mp3", filePureName))

	// 下载视频
	req_v, err := http.NewRequestWithContext(ctx, "GET", format_v.URL, nil)
	if err != nil {
		return err
	}
	req_v.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0")

	// if err = utils.ReqWriter(yd.Client.HTTPClient, req_v, part, input_v, make(chan struct{}), callback); err != nil {
	// 	part.Status = "下载视频出错"
	// 	return err
	// }

	// format_a := getBestHighFormat(filterFormats(video.Formats, "audio"))
	// req, err := http.NewRequestWithContext(ctx, "GET", format_a.URL, nil)
	// if err != nil {
	// 	return err
	// }
	// if err = utils.ReqWriter(yd.Client.HTTPClient, req, part, input_a, make(chan struct{}), callback); err != nil {
	// 	part.Status = "下载音频出错"
	// 	return err
	// }

	part.Status = "正在合并"

	// todo? 报错提醒
	// if err := utils.ClearDirs([]string{path_v, path_a}); err != nil {
	// 	part.Status = "删除缓存文件失败"
	// }
	return nil
}

func (yd *YouTubeDownloader) DownloadVideo(part *shared.Part) error { return nil }
func (yd *YouTubeDownloader) DownloadAudio(part *shared.Part) error { return nil }
func (yd *YouTubeDownloader) DownloadThumbnail(ctx context.Context, part *shared.Part) error {
	return nil
}
func (yd *YouTubeDownloader) DownloadSubtitle(ctx context.Context, part *shared.Part) error {
	return nil
}
func (yd *YouTubeDownloader) Clear() error { return nil }

func getBestThumbnail(thumbnails youtube.Thumbnails) *youtube.Thumbnail {
	var bestThumbnail *youtube.Thumbnail

	var height = uint(0)

	for _, tn := range thumbnails {

		if tn.Height > height {
			height = tn.Height
			bestThumbnail = &tn
		}
	}

	return bestThumbnail
}
func (yd *YouTubeDownloader) getVideoDataByID(videoID string) (*youtube.Video, error) {

	video, err := yd.Client.GetVideo(videoID)
	if err != nil {
		return nil, err
	}

	if video == nil {
		return nil, errors.New("啥也木有")
	}

	return video, nil
}

func formats2Qualities(formats youtube.FormatList) []string {
	maxHeight := 0
	for _, format := range formats {
		if format.Height > maxHeight {
			maxHeight = format.Height
		}
	}

	return nil
	// todo
	// return utils.GetQualities(maxHeight)
}

func (yd *YouTubeDownloader) downloadSubtitle(video *youtube.Video) {
	yd.Client.GetTranscriptCtx(context.Background(), video, "en")
	transcript, _ := yd.Client.GetTranscript(video, "en")
	srt := utils.TranscriptToSrt(transcript)
	utils.WriteSrt(".srt", srt)
}

func isPlaylist(url string) bool {
	re := regexp.MustCompile(`list=([^&]+)`)
	match := re.FindStringSubmatch(url)

	return len(match) > 1
}

func extractID(url string) string {
	re := regexp.MustCompile(`(v|list)=([^&]+)`)
	match := re.FindStringSubmatch(url)

	if len(match) > 2 {
		id := match[2]
		return id
	} else {
		return url
	}
}

func getBestHighFormatHeight(formats youtube.FormatList) int {
	var MaxHeight = 0
	for _, format := range formats {
		if format.Height > MaxHeight {
			MaxHeight = format.Height
		}
	}
	return MaxHeight
}

// 如果有 直接用, 没有就找比它高1级的 kind: video/audio
func filterFormats(formats youtube.FormatList, kind string) []youtube.Format {
	var filteredFormats []youtube.Format
	for _, format := range formats {
		if strings.Contains(format.MimeType, kind) {
			filteredFormats = append(filteredFormats, format)
		}
	}
	return filteredFormats
}

func getBestHighFormat(formats []youtube.Format) youtube.Format {
	var bestFormat youtube.Format
	for _, format := range formats {
		if format.Bitrate > bestFormat.Bitrate {
			bestFormat = format
		}
	}
	return bestFormat
}

func getTargetYtbVideo(targetHeight int, videos youtube.FormatList) *youtube.Format {

	sort.Slice(videos, func(i, j int) bool {
		return videos[i].Height < videos[j].Height
	})

	var closestVideo *youtube.Format
	for _, video := range videos {
		if video.Height >= targetHeight {
			return &video
		}
		closestVideo = &video
	}
	return closestVideo
}
