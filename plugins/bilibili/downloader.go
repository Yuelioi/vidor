package bilibili

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"

	"github.com/Yuelioi/vidor/shared"
	"github.com/Yuelioi/vidor/utils"
)

type Downloader struct {
	ctx       context.Context
	cancel    context.CancelFunc
	client    *Client
	config    shared.Config
	magicName string
	Notice    shared.Notice
	biliDownloadParams
	biliThumbnailParams
}

func New(ctx context.Context, config shared.Config) shared.Downloader {

	client := NewClient(config.SESSDATA)
	ctx, cancel := context.WithCancel(context.Background())

	return &Downloader{
		ctx:       ctx,
		magicName: config.MagicName,
		cancel:    cancel,
		client:    client,
	}
}

func (bd *Downloader) PluginMeta() shared.PluginMeta {
	return shared.PluginMeta{
		Name: "bilibili",
		Regexs: []*regexp.Regexp{
			regexp.MustCompile(`https://www\.bilibili\.com/video/BV.+`),
			regexp.MustCompile(`https://www\.bilibili\.com/video/av.+`)},
	}
}

func (bd *Downloader) Show(ctx context.Context, link string) (*shared.PlaylistInfo, error) {

	// 获取b站播放列表信息
	aid, bvid := extractAidBvid(link)
	biliPlayList, err := bd.client.GetPlaylistInfo(aid, bvid)
	if err != nil {
		return nil, fmt.Errorf("ShowInfo %s", err)
	}

	// 填充列表信息
	var playList shared.PlaylistInfo
	if biliPlayList.Data.IsSeason {
		playList = biliSeasonToPlaylistInfo(*biliPlayList)
	} else {
		playList = biliPageToPlaylistInfo(biliPlayList.Data.BVID, *biliPlayList)
	}

	thumbnailPath := filepath.Join(os.TempDir(), "vidor", "info_thumbnail.jpg")
	fmt.Printf("thumbnailPath: %v\n", thumbnailPath)

	img, err := utils.GetThumbnail(bd.client.HTTPClient, playList.Cover, thumbnailPath)
	if err != nil {
		return nil, err
	}
	playList.Cover = img

	fmt.Printf("playList: %v\n", playList)

	return &playList, nil
}

func (bd *Downloader) Parse(ctx context.Context, playlist *shared.PlaylistInfo) (*shared.PlaylistInfo, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	errors := make(chan error, len(playlist.StreamInfos))

	for index, streamInfo := range playlist.StreamInfos {
		if !streamInfo.Selected {
			continue
		}

		wg.Add(1)
		go func(streamInfo shared.StreamInfo, index int) {
			defer wg.Done()

			cid, err := strconv.Atoi(streamInfo.SessionId)
			if err != nil {
				errors <- err
				return
			}

			biliDownInfo, err := bd.client.GetVideoDownloadInfo(streamInfo.ID, cid)
			if err != nil {
				errors <- err
				return
			}

			mu.Lock()
			playlist.StreamInfos[index].Videos = make([]shared.Format, 0)
			playlist.StreamInfos[index].Audios = make([]shared.Format, 0)
			mu.Unlock()

			mu.Lock()
			for _, format := range biliDownInfo.Data.SupportFormats {
				playlist.StreamInfos[index].Videos = append(playlist.StreamInfos[index].Videos, shared.Format{
					IDtag:   format.Quality,
					Quality: format.DisplayDesc,
					Codecs:  format.Codecs,
				})
			}
			for _, audio := range biliDownInfo.Data.Dash.Audios {
				playlist.StreamInfos[index].Audios = append(playlist.StreamInfos[index].Audios, shared.Format{
					IDtag:   audio.ID,
					Quality: fmt.Sprint(audio.Bandwidth),
					Codecs:  []string{audio.Codecs},
				})
			}
			mu.Unlock()
		}(streamInfo, index)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			return nil, err
		}
	}

	return playlist, nil
}

func (bd *Downloader) Do(ctx context.Context, part *shared.Part) error {
	return nil
}

func (bd *Downloader) GetMeta(part *shared.Part, callback shared.Callback) error {

	// 提取分P索引 如果有的话, 没有就是0
	index, _ := extractIndex(part.URL)
	aid, bvid := extractAidBvid(part.URL)

	// 获取视频基础信息
	biliPlayList, err := bd.client.GetPlaylistInfo(aid, bvid)
	if err != nil {
		return fmt.Errorf("获取播放列表信息失败: %v", err.Error())
	}

	var biliBase biliBaseParams
	if biliPlayList.Data.IsSeason {
		biliBase = processSeasonData(biliPlayList.Data, aid, bvid)
	} else {
		biliBase = processPagesData(biliPlayList.Data, index)
	}

	// 创建必要文件夹
	if err := utils.CreateDirs([]string{part.DownloadDir}); err != nil {
		return err
	}
	part.Author = biliBase.author
	part.Index = biliBase.index
	part.Title = utils.SanitizeFileName(biliBase.title)
	part.MagicName = utils.MagicName(bd.magicName, part.WorkDirName, part.Title, part.Index+1)
	part.Path = filepath.Join(part.DownloadDir, part.MagicName+".mp4")
	coverPath := filepath.Join(part.DownloadDir, biliBase.coverName)

	// 获取封面参数
	bd.biliThumbnailParams = biliThumbnailParams{
		thumbnailURL: biliBase.thumbnailURL,
		coverPath:    coverPath,
		coverURL:     biliBase.coverURL,
	}

	// 获取视频下载信息
	biliDownInfo, err := bd.client.GetVideoDownloadInfo(biliBase.bvid, biliBase.cid)
	if err != nil {
		return fmt.Errorf("获取视频下载信息失败: %v", err.Error())
	}
	if biliDownInfo.Code != 0 {
		return fmt.Errorf("视频下载信息返回错误: %v", biliDownInfo.Message)
	}

	// 开始下载
	targetID := 0
	// targetID, err := utils.GetQualityID(part.Quality, qualities)
	// if err != nil {
	// 	return err
	// }

	var videoURL, audioURL string
	video := getTargetVideo(targetID, biliDownInfo.Data.Dash.Videos)
	if video == nil {
		videoURL = ""
	} else {
		videoURL = video.BaseURL
	}

	audio := getTargetAudio(biliDownInfo.Data.Dash.Audios)
	if audio == nil {
		audioURL = ""
	} else {
		audioURL = audio.BaseURL
	}

	bd.biliDownloadParams = biliDownloadParams{
		videoURL: videoURL,
		audioURL: audioURL,
		index:    index,
	}

	// part.Quality, err = utils.GetQualityLabel(targetID, qualities)
	// if err != nil {
	// 	return err
	// }

	part.Status = "加载元数据成功"
	callback(shared.NoticeData{EventName: "updateInfo", Message: part})
	return nil
}

func (bd *Downloader) DownloadThumbnail(part *shared.Part, callback shared.Callback) error {
	// 缩略图
	thumbnailPath := filepath.Join(part.DownloadDir, "data", "thumbnail_"+part.MagicName+".jpg")
	thumbnailLocalPath, err := utils.GetThumbnail(bd.client.HTTPClient, bd.biliThumbnailParams.thumbnailURL, thumbnailPath)
	if err != nil {
		return err
	}
	part.Thumbnail = thumbnailLocalPath

	// 封面
	if _, err = os.Stat(bd.coverPath); os.IsNotExist(err) {
		_, err := utils.GetCover(bd.client.HTTPClient, bd.biliThumbnailParams.coverURL, bd.coverPath)
		if err != nil {
			return fmt.Errorf("缓存封面失败: %v", err.Error())
		}
	}
	return nil

}
func (bd *Downloader) DownloadVideo(part *shared.Part, callback shared.Callback) error {
	part.Status = "下载视频"
	callback(shared.NoticeData{EventName: "updateInfo", Message: part})
	return bd.download(part, bd.biliDownloadParams.videoURL, "mp4", callback)
}
func (bd *Downloader) DownloadAudio(part *shared.Part, callback shared.Callback) error {
	part.Status = "下载音频"
	callback(shared.NoticeData{EventName: "updateInfo", Message: part})
	return bd.download(part, bd.biliDownloadParams.audioURL, "mp3", callback)
}

func (bd *Downloader) DownloadSubtitle(part *shared.Part, callback shared.Callback) error {
	return nil
}
func (bd *Downloader) Combine(ffmpegPath string, part *shared.Part) error {

	input_v := filepath.Join(part.DownloadDir, fmt.Sprintf("%s_temp.mp4", part.MagicName))
	input_a := filepath.Join(part.DownloadDir, fmt.Sprintf("%s_temp.mp3", part.MagicName))
	output_v := filepath.Join(part.DownloadDir, fmt.Sprintf("%s.mp4", part.MagicName))
	logFilePath := filepath.Join(part.DownloadDir, "data", fmt.Sprintf("log_%s.txt", part.MagicName))

	utils.CombineAV(bd.ctx, ffmpegPath, input_v, input_a, output_v, logFilePath)
	return nil
}

func (bd *Downloader) Clear(ctx context.Context, part *shared.Part) error {
	return nil
}

func (bd *Downloader) Cancel(ctx context.Context, part *shared.Part) error {
	bd.cancel()

	// callback(shared.NoticeData{
	// 	EventName: "updateInfo",
	// 	Message:   part,
	// })

	return nil
}
func (bd *Downloader) Pause(ctx context.Context, part *shared.Part) error {
	return nil
}

func (bd *Downloader) download(part *shared.Part, link, ext string, callback shared.Callback) error {
	tempPath := filepath.Join(part.DownloadDir, fmt.Sprintf("%s_temp.%s", part.MagicName, ext))

	req, err := bd.client.NewRequest("Get", link, nil)

	if err != nil {

		return err
	}
	mediaURL, err := url.Parse(link)
	if err != nil {

		return err
	}
	req.URL = mediaURL
	utils.ReqWriter(bd.ctx, bd.client.HTTPClient, req, part, tempPath, callback)
	return nil
}
