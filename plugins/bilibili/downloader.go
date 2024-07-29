package bilibili

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/Yuelioi/vidor/shared"
	"github.com/Yuelioi/vidor/utils"
)

type Downloader struct {
	ctx       context.Context
	cancel    context.CancelFunc
	Client    *http.Client
	Notice    shared.Notice
	userState userStatus // 0未登录 1已登录 2Vip
	biliDownloadParams
	biliThumbnailParams
}

func New(notice shared.Notice) shared.Downloader {
	ctx, cancel := context.WithCancel(context.Background())
	return &Downloader{
		Notice: notice,
		ctx:    ctx,
		cancel: cancel,
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

func (bd *Downloader) ShowInfo(link string, config shared.Config, callback shared.Callback) (*shared.PlaylistInfo, error) {
	client, err := utils.GetClient(config.ProxyURL, config.UseProxy)
	if err != nil {
		return nil, err
	}

	// 初始化Client
	bd.Client = client

	// 获取登录信息
	bd.getUserStates(config.SESSDATA)

	// 获取b站播放列表信息
	var playList shared.PlaylistInfo
	aid, bvid := extractAidBvid(link)
	biliPlayList, err := getPlaylistInfo(*client, aid, bvid, config.SESSDATA)
	if err != nil {
		return nil, err
	}
	if biliPlayList.Data.IsSeason {
		playList = biliSeasonToPlaylistInfo(*biliPlayList)
	} else {
		playList = biliPageToPlaylistInfo(*biliPlayList)
	}

	copyQualities := make([]shared.StreamQuality, len(qualities))
	copy(copyQualities, qualities)
	if bd.userState == NoLogin {
		copyQualities = copyQualities[1:4]
	} else if bd.userState == Login {
		copyQualities = copyQualities[2:6]
	} else {
		copyQualities = copyQualities[2:]
	}

	for _, qu := range copyQualities {
		playList.Qualities = append(playList.Qualities, qu.Label)
	}

	thumbnailPath := filepath.Join(os.TempDir(), "vidor", "info_thumbnail.jpg")
	fmt.Printf("thumbnailPath: %v\n", thumbnailPath)

	img, err := utils.GetThumbnail(client, playList.Cover, thumbnailPath)
	if err != nil {
		return nil, err
	}
	playList.Cover = img

	return &playList, nil
}

func (bd *Downloader) GetMeta(config shared.Config, part *shared.Part, callback shared.Callback) error {

	client, err := utils.GetClient(config.ProxyURL, config.UseProxy)
	if err != nil {
		return err
	}

	// 初始化Client
	bd.Client = client

	// 获取登录信息
	bd.getUserStates(config.SESSDATA)

	// 提取分P索引 如果有的话, 没有就是0
	index, _ := extractIndex(part.Url)
	aid, bvid := extractAidBvid(part.Url)

	// 获取视频基础信息
	biliPlayList, err := getPlaylistInfo(*bd.Client, aid, bvid, config.SESSDATA)
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
	part.MagicName = utils.MagicName(config.MagicName, part.WorkDirName, part.Title, part.Index+1)
	part.Path = filepath.Join(part.DownloadDir, part.MagicName+".mp4")
	coverPath := filepath.Join(part.DownloadDir, biliBase.coverName)

	// 获取封面参数
	bd.biliThumbnailParams = biliThumbnailParams{
		thumbnailUrl: biliBase.thumbnailUrl,
		coverPath:    coverPath,
		coverUrl:     biliBase.coverUrl,
	}

	// 获取视频下载信息
	biliDownInfo, err := getVideoDownloadInfo(*bd.Client, biliBase.bvid, biliBase.cid, config.SESSDATA)
	if err != nil {
		return fmt.Errorf("获取视频下载信息失败: %v", err.Error())
	}
	if biliDownInfo.Code != 0 {
		return fmt.Errorf("视频下载信息返回错误: %v", biliDownInfo.Message)
	}

	// 开始下载
	targetID, err := utils.GetQualityID(part.Quality, qualities)
	if err != nil {
		return err
	}

	var videoUrl, audioUrl string
	video := getTargetVideo(targetID, biliDownInfo.Data.Dash.Videos)
	if video == nil {
		videoUrl = ""
	} else {
		videoUrl = video.BaseURL
	}

	audio := getTargetAudio(biliDownInfo.Data.Dash.Audios)
	if audio == nil {
		audioUrl = ""
	} else {
		audioUrl = audio.BaseURL
	}

	bd.biliDownloadParams = biliDownloadParams{
		videoUrl: videoUrl,
		audioUrl: audioUrl,
		sessdata: config.SESSDATA,
		index:    index,
	}

	part.Quality, err = utils.GetQualityLabel(targetID, qualities)
	if err != nil {
		return err
	}

	part.Status = "加载元数据成功"
	callback(shared.NoticeData{EventName: "updateInfo", Message: part})
	return nil
}

func (bd *Downloader) DownloadThumbnail(part *shared.Part, callback shared.Callback) error {
	// 缩略图
	thumbnailPath := filepath.Join(part.DownloadDir, "data", "thumbnail_"+part.MagicName+".jpg")
	thumbnailLocalPath, err := utils.GetThumbnail(bd.Client, bd.biliThumbnailParams.thumbnailUrl, thumbnailPath)
	if err != nil {
		return err
	}
	part.Thumbnail = thumbnailLocalPath

	// 封面
	if _, err = os.Stat(bd.coverPath); os.IsNotExist(err) {
		_, err := utils.GetCover(bd.Client, bd.biliThumbnailParams.coverUrl, bd.coverPath)
		if err != nil {
			return fmt.Errorf("缓存封面失败: %v", err.Error())
		}
	}
	return nil

}
func (bd *Downloader) DownloadVideo(part *shared.Part, callback shared.Callback) error {
	part.Status = "下载视频"
	callback(shared.NoticeData{EventName: "updateInfo", Message: part})
	return bd.download(part, bd.biliDownloadParams.videoUrl, "mp4", callback)
}
func (bd *Downloader) DownloadAudio(part *shared.Part, callback shared.Callback) error {
	part.Status = "下载音频"
	callback(shared.NoticeData{EventName: "updateInfo", Message: part})
	return bd.download(part, bd.biliDownloadParams.audioUrl, "mp3", callback)
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
func (bd *Downloader) Clear(part *shared.Part, callback shared.Callback) error {
	return nil
}
func (bd *Downloader) StopDownload(part *shared.Part, callback shared.Callback) error {
	bd.cancel()

	callback(shared.NoticeData{
		EventName: "updateInfo",
		Message:   part,
	})

	return nil
}
func (bd *Downloader) PauseDownload(part *shared.Part, callback shared.Callback) error {
	return nil
}
