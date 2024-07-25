package plugins

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"

	"github.com/Yuelioi/vidor/shared"
	"github.com/Yuelioi/vidor/utils"
)

const (
	bilibiliApiURL = "https://api.bilibili.com"
)

type page struct {
	CID       int    `json:"cid"`
	Page      int    `json:"page"`
	Title     string `json:"part"`
	Duration  int    `json:"duration"`
	Thumbnail string `json:"first_frame"`
	Dimension struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"dimension"`
}

type biliUserInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		MID       int `json:"mid"`
		VIPType   int `json:"vip_type"`
		VIPStatus int `json:"vip_status"`
		Label     struct {
			Text string `json:"text"`
		} `json:"label"`
	} `json:"data"`
}

type biliPlayListData struct {
	TName   string `json:"tname"`
	AID     int    `json:"aid"`
	BVID    string `json:"bvid"`
	CID     int    `json:"cid"`
	Pic     string `json:"pic"`
	PubDate int    `json:"pubdate"`
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Owner   struct {
		Name string `json:"name"`
	} `json:"owner"`
	Pages     []page     `json:"pages"`
	Subtitle  *struct{}  `json:"subtitle"`
	UgcSeason *ugcSeason `json:"ugc_season"`
	IsSeason  bool       `json:"is_season_display"`
}

type biliPlaylistInfo struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    biliPlayListData `json:"data"`
}

type biliVideo struct {
	ID        int      `json:"id"`
	BaseUrl   string   `json:"baseUrl"`
	BaseURL   string   `json:"base_url"`
	BackupUrl []string `json:"backupUrl"`
	BackupURL []string `json:"backup_url"`
	Codecs    string   `json:"codecs"`
	Width     int      `json:"width"`
	Height    int      `json:"height"`
	FrameRate string   `json:"frameRate"`
}

type biliAudio struct {
	BaseUrl   string   `json:"baseUrl"`
	BaseURL   string   `json:"base_url"`
	BackupUrl []string `json:"backupUrl"`
	BackupURL []string `json:"backup_url"`
	Bandwidth int      `json:"bandwidth"`
	Codecs    string   `json:"codecs"`
}

type episode struct {
	Aid   int    `json:"aid"`
	CID   int    `json:"cid"`
	Bvid  string `json:"bvid"`
	Title string `json:"title"`
	Arc   struct {
		Pic      string `json:"pic"`
		Title    string `json:"title"`
		Pubdate  int    `json:"pubdate"`
		Ctime    int    `json:"ctime"`
		Desc     string `json:"desc"`
		Duration int    `json:"duration"`
	} `json:"arc"`
	Page page `json:"page"`
}

type ugcSeason struct {
	Title    string `json:"title"`
	Cover    string `json:"cover"`
	Sections []struct {
		Title    string    `json:"title"`
		Episodes []episode `json:"episodes"`
	} `json:"sections"`
}

type biliDownloadInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Timelength int `json:"timelength"`
		Dash       struct {
			Videos []biliVideo `json:"video"`
			Audios []biliAudio `json:"audio"`
		} `json:"dash"`
		SupportFormats []struct {
			Quality     int      `json:"quality"`
			Format      string   `json:"format"`
			NewLabel    string   `json:"new_description"`
			DisplayDesc string   `json:"display_desc"`
			Superscript string   `json:"superscript"`
			Codecs      []string `json:"codecs"`
		} `json:"support_formats"`
	} `json:"data"`
}

// 视频基础参数
type biliBaseParams struct {
	cid  int
	bvid string

	title  string
	author string

	pubDate  int
	duration int

	coverUrl     string
	coverName    string
	thumbnailUrl string
}

// 视频下载参数
type biliDownloadParams struct {
	videoUrl string
	audioUrl string
	index    int
	sessdata string
}

// 视频封面参数
type biliThumbnailParams struct {
	thumbnailUrl string
	coverPath    string
	coverUrl     string
}
type userStatus int

const (
	NoLogin userStatus = iota
	Login
	Vip
)

type BilibiliDownloader struct {
	Client      *http.Client
	Notice      shared.Notice
	stopChannel chan struct{} // 在GetMate时 初始化chan
	userState   userStatus    // 0未登录 1已登录 2Vip
	biliBaseParams
	biliDownloadParams
	biliThumbnailParams
}

// https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/docs/video/videostream_url.md
var qualities = []shared.VideoQuality{
	{ID: 6, Label: "240P"},      // 仅 MP4 格式支持, 仅 platform=html5 时有效
	{ID: 16, Label: "360P"},     // 流畅
	{ID: 32, Label: "480P"},     // 清晰
	{ID: 64, Label: "720P"},     // 高清, WEB 端默认值, B站前端需要登录才能选择，但是直接发送请求可以不登录就拿到 720P 的取流地址, 无 720P 时则为 720P60
	{ID: 74, Label: "720P60"},   // 高帧率, 登录认证
	{ID: 80, Label: "1080P"},    // 高清, TV 端与 APP 端默认值, 登录认证
	{ID: 112, Label: "1080P+"},  // 高码率, 大会员认证
	{ID: 116, Label: "1080P60"}, // 高帧率, 大会员认证
	{ID: 120, Label: "4K 超清"},   // 需要 fnval&128=128 且 fourk=1, 大会员认证
	{ID: 125, Label: "HDR 真彩色"}, // 仅支持 DASH 格式, 需要 fnval&64=64, 大会员认证
	{ID: 126, Label: "杜比视界"},    // 仅支持 DASH 格式, 需要 fnval&512=512, 大会员认证
	{ID: 127, Label: "8K 超高清"},  // 仅支持 DASH 格式, 需要 fnval&1024=1024, 大会员认证
	{ID: 999, Label: "最高画质(💗)"}, // 仅支持 DASH 格式, 需要 fnval&1024=1024, 大会员认证
}

func NewBiliDownloader(notice shared.Notice) *BilibiliDownloader {
	return &BilibiliDownloader{
		Notice: notice,
	}
}

func (bd *BilibiliDownloader) PluginMeta() shared.PluginMeta {
	return shared.PluginMeta{
		Name: "bilibili",
		Regexs: []*regexp.Regexp{
			regexp.MustCompile(`https://www\.bilibili\.com/video/BV.+`),
			regexp.MustCompile(`https://www\.bilibili\.com/video/av.+`)},
	}
}

func (bd *BilibiliDownloader) getUserStates(sessdata string) {
	apiUrl := bilibiliApiURL + "/x/vip/web/user/info"

	data, err := doBiliReq(*bd.Client, apiUrl, sessdata)
	var biliUserInfo biliUserInfo
	err = json.Unmarshal(data, &biliUserInfo)
	if err != nil {
		bd.userState = NoLogin
		return
	}

	if biliUserInfo.Code == -101 {
		bd.userState = NoLogin
		return
	} else if biliUserInfo.Data.VIPStatus == 1 {
		bd.userState = Vip
		return
	}
	bd.userState = Login
}

func (bd *BilibiliDownloader) ShowInfo(link string, config shared.Config, callback shared.Callback) (*shared.PlaylistInfo, error) {
	client, err := utils.GetClient(config.ProxyURL, config.UseProxy)
	if err != nil {
		return nil, err
	}

	// 初始化Client
	bd.Client = client

	// 获取登录信息
	bd.getUserStates(config.SESSDATA)

	aid, bvid := extractAidBvid(link)
	bpi, err := getPlaylistInfo(*client, aid, bvid, config.SESSDATA)
	if err != nil {
		return nil, err
	}

	var pi shared.PlaylistInfo

	if bpi.Data.IsSeason {
		pi = biliSeasonToPlaylistInfo(*bpi)
	} else {
		pi = biliPageToPlaylistInfo(*bpi)
	}

	pi.Description = bpi.Data.Desc

	copyQualities := make([]shared.VideoQuality, len(qualities))
	copy(copyQualities, qualities)
	if bd.userState == NoLogin {
		copyQualities = copyQualities[1:4]
	} else if bd.userState == Login {
		copyQualities = copyQualities[2:6]
	} else {
		copyQualities = copyQualities[2:]
	}

	for _, qu := range copyQualities {
		pi.Qualities = append(pi.Qualities, qu.Label)
	}

	img, err := utils.GetThumbnail(client, pi.Thumbnail)
	if err != nil {
		return nil, err
	}
	pi.Thumbnail = img

	return &pi, nil
}

func (bd *BilibiliDownloader) GetMeta(config shared.Config, part *shared.Part, callback shared.Callback) error {

	bd.stopChannel = make(chan struct{})
	// 提取分P索引 如果有的话, 没有就是0
	index, _ := extractIndex(part.Url)
	aid, bvid := extractAidBvid(part.Url)

	// 获取视频基础信息
	bpi, err := getPlaylistInfo(*bd.Client, aid, bvid, config.SESSDATA)
	if err != nil {
		return fmt.Errorf("获取播放列表信息失败: %v", err.Error())
	}

	var biliBase biliBaseParams

	if bpi.Data.IsSeason {
		biliBase = processSeasonData(bpi.Data, aid, bvid)
	} else {
		biliBase = processPagesData(bpi.Data, index)
	}

	bd.biliBaseParams = biliBase
	part.Author = biliBase.author

	// 创建必要文件夹
	if err := utils.CreateDirs([]string{part.DownloadDir}); err != nil {
		return err
	}

	coverPath := filepath.Join(part.DownloadDir, biliBase.coverName)
	part.Title = utils.SanitizeFileName(biliBase.title)

	part.Path = filepath.Join(part.DownloadDir, part.Title+".mp4")

	// 获取封面参数
	bd.biliThumbnailParams = biliThumbnailParams{
		thumbnailUrl: biliBase.thumbnailUrl,
		coverPath:    coverPath,
		coverUrl:     biliBase.coverUrl,
	}

	// 获取视频下载信息
	bdi, err := getVideoDownloadInfo(*bd.Client, biliBase.bvid, biliBase.cid, config.SESSDATA)
	if err != nil {
		return fmt.Errorf("获取视频下载信息失败: %v", err.Error())
	}
	if bdi.Code != 0 {
		return fmt.Errorf("视频下载信息返回错误: %v", bdi.Message)
	}

	// 开始下载
	targetID, err := utils.GetQualityID(part.Quality, qualities)
	if err != nil {
		return err
	}

	var videoUrl, audioUrl string
	video := getTargetVideo(targetID, bdi.Data.Dash.Videos)
	if video == nil {
		videoUrl = ""
	} else {
		videoUrl = video.BaseURL
	}

	audio := getTargetAudio(bdi.Data.Dash.Audios)
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

func (bd *BilibiliDownloader) DownloadThumbnail(part *shared.Part, callback shared.Callback) error {
	// 缩略图
	thumbnailLocal, err := utils.GetThumbnail(bd.Client, bd.biliThumbnailParams.thumbnailUrl)
	if err != nil {
		return err
	}
	part.Thumbnail = thumbnailLocal

	// 封面
	if _, err = os.Stat(bd.coverPath); os.IsNotExist(err) {
		_, err := utils.GetCover(bd.Client, bd.biliThumbnailParams.coverUrl, bd.coverPath)
		if err != nil {
			return fmt.Errorf("缓存封面失败: %v", err.Error())
		}
	}
	return nil

}
func (bd *BilibiliDownloader) DownloadVideo(part *shared.Part, callback shared.Callback) error {
	part.Status = "下载视频"
	callback(shared.NoticeData{EventName: "updateInfo", Message: part})
	return bd.download(part, bd.biliDownloadParams.videoUrl, "mp4", callback)
}
func (bd *BilibiliDownloader) DownloadAudio(part *shared.Part, callback shared.Callback) error {
	part.Status = "下载音频"
	callback(shared.NoticeData{EventName: "updateInfo", Message: part})
	return bd.download(part, bd.biliDownloadParams.audioUrl, "mp3", callback)
}

func (bd *BilibiliDownloader) DownloadSubtitle(part *shared.Part, callback shared.Callback) error {
	return nil
}
func (bd *BilibiliDownloader) Combine(ffmpegPath string, part *shared.Part) error {
	utils.CombineAV(ffmpegPath, part, bd.stopChannel)
	return nil
}
func (bd *BilibiliDownloader) Clear(part *shared.Part, callback shared.Callback) error {
	return nil
}
func (bd *BilibiliDownloader) StopDownload(part *shared.Part, callback shared.Callback) error {
	if bd.stopChannel != nil {
		println("关闭通道")
		close(bd.stopChannel)
		bd.stopChannel = nil
	}
	part.State = shared.TaskStatus.Finished
	callback(shared.NoticeData{
		EventName: "updateInfo",
		Message:   part,
	})

	return nil
}
func (bd *BilibiliDownloader) PauseDownload(part *shared.Part, callback shared.Callback) error {
	return nil
}

func (bd *BilibiliDownloader) download(part *shared.Part, link, ext string, callback shared.Callback) error {
	path := filepath.Join(part.DownloadDir, fmt.Sprintf("%s_temp.%s", part.Title, ext))

	req, err := downloadReq(bd.sessdata)
	if err != nil {

		return err
	}
	mediaUrl, err := url.Parse(link)
	if err != nil {

		return err
	}
	req.URL = mediaUrl
	utils.ReqWriter(bd.Client, req, part, path, bd.stopChannel, callback)
	return nil
}

// 获取合集信息
func processSeasonData(data biliPlayListData, aid int, bvid_input string) biliBaseParams {
	episodes := data.UgcSeason.Sections[0].Episodes
	index := 0
	for epid, epi := range episodes {
		if epi.Aid == aid || epi.Bvid == bvid_input {
			index = epid
		}
	}
	workDirname := data.UgcSeason.Title
	coverName := fmt.Sprintf("%02d_%s.jpg", index+1, utils.SanitizeFileName(workDirname))

	return biliBaseParams{
		cid:      episodes[index].CID,
		bvid:     data.UgcSeason.Sections[0].Episodes[index].Bvid,
		title:    fmt.Sprintf("%02d_%s", index+1, utils.SanitizeFileName(episodes[index].Title)),
		author:   data.Owner.Name,
		pubDate:  episodes[index].Arc.Pubdate,
		duration: episodes[index].Arc.Duration,

		coverUrl:     data.UgcSeason.Cover,
		coverName:    coverName,
		thumbnailUrl: episodes[index].Arc.Pic,
	}
}

// 获取普通分p信息
func processPagesData(bpi biliPlayListData, index int) biliBaseParams {
	var title, thumbnailUrl string

	if len(bpi.Pages) == 1 {
		title = bpi.Title
	} else {
		title = bpi.Pages[index].Title
	}

	if bpi.Pages[index].Thumbnail != "" {
		thumbnailUrl = bpi.Pages[index].Thumbnail
	} else {
		thumbnailUrl = bpi.Pic
	}

	workDirname := bpi.Title

	return biliBaseParams{
		cid:          bpi.Pages[index].CID,
		bvid:         bpi.BVID,
		title:        fmt.Sprintf("%02d_%s", index+1, utils.SanitizeFileName(title)),
		author:       bpi.Owner.Name,
		coverUrl:     bpi.Pic,
		coverName:    utils.SanitizeFileName(workDirname) + ".jpg",
		thumbnailUrl: thumbnailUrl,
	}
}

func getTargetAudio(audios []biliAudio) *biliAudio {

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
func getTargetVideo(targeID int, videos []biliVideo) *biliVideo {
	// 从小到大排序, 应对竖版视频
	sort.Slice(videos, func(i, j int) bool {
		return videos[i].Height < videos[j].Height
	})

	var closestVideo *biliVideo
	for _, video := range videos {
		if video.ID >= targeID {
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

func doBiliReq(client http.Client, link, sessdata string) (body []byte, err error) {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0")
	req.Header.Set("Referer", "https://www.bilibili.com/vedio")

	if len(sessdata) > 0 {
		req.Header.Set("Cookie", "SESSDATA="+sessdata)
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	return
}

func getPlaylistInfo(client http.Client, aid int, bvid, sessdata string) (*biliPlaylistInfo, error) {
	var bv biliPlaylistInfo

	body, err := doBiliReq(client, fmt.Sprintf(`%s/x/web-interface/view?aid=%d&bvid=%s`, bilibiliApiURL, aid, bvid), sessdata)
	if err != nil {
		fmt.Println("Error cannot fetch Aid:", err)
		return nil, err
	}
	err = json.Unmarshal(body, &bv)
	if err != nil {
		fmt.Println("Error  JSON:", err)
		return nil, err
	}

	if bv.Code != 0 {
		return nil, errors.New(bv.Message)
	}

	return &bv, nil

}

func getVideoDownloadInfo(client http.Client, bvid string, cid int, sessdata string) (*biliDownloadInfo, error) {
	body, err := doBiliReq(client, fmt.Sprintf("%s/x/player/wbi/playurl?bvid=%s&cid=%d&fnval=4048&fourk=1", bilibiliApiURL, bvid, cid), sessdata)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	bi := biliDownloadInfo{}
	if err = json.Unmarshal(body, &bi); err != nil {
		return nil, err
	}
	return &bi, nil
}

// 将B站视频列表信息转为通用的列表信息
func biliPageToPlaylistInfo(biliInfo biliPlaylistInfo) shared.PlaylistInfo {
	var videoInfo shared.PlaylistInfo

	videoInfo.Url = fmt.Sprintf("https://www.bilibili.com/video/%s", biliInfo.Data.BVID)
	videoInfo.Thumbnail = biliInfo.Data.Pic

	videoInfo.WorkDirName = utils.SanitizeFileName(biliInfo.Data.Title)
	videoInfo.Author = biliInfo.Data.Owner.Name

	videoInfo.Parts = make([]shared.Part, 0)

	var MaxHeight = 0

	// TODO 是否在队列就显示封面?
	for _, page := range biliInfo.Data.Pages {
		// var thumbnailUrl string

		// if biliInfo.Data.Pages[index].Thumbnail != "" {
		// 	thumbnailUrl = biliInfo.Data.Pages[index].Thumbnail
		// } else {
		// 	thumbnailUrl = biliInfo.Data.Pic
		// }

		videoInfo.Parts = append(videoInfo.Parts, shared.Part{
			Url:       fmt.Sprintf("https://www.bilibili.com/video/%s?p=%d", biliInfo.Data.BVID, page.Page),
			Title:     page.Title,
			Thumbnail: "",
		})

		min := func(a, b int) int {
			if a < b {
				return a
			} else {
				return b
			}
		}(page.Dimension.Width, page.Dimension.Height)

		if MaxHeight < min {
			MaxHeight = min
		}

	}

	// 如果只有一个 就用主标题
	if len(biliInfo.Data.Pages) == 1 {
		videoInfo.Parts[0].Title = videoInfo.WorkDirName
	}
	// TODO
	// videoInfo.Qualities = utils.GetQualities(MaxHeight)

	return videoInfo
}

func biliSeasonToPlaylistInfo(biliInfo biliPlaylistInfo) shared.PlaylistInfo {
	var videoInfo shared.PlaylistInfo

	videoInfo.Url = fmt.Sprintf("https://www.bilibili.com/video/%s", biliInfo.Data.BVID)
	videoInfo.Thumbnail = biliInfo.Data.Pic
	videoInfo.WorkDirName = utils.SanitizeFileName(biliInfo.Data.UgcSeason.Title)
	videoInfo.Author = biliInfo.Data.Owner.Name

	videoInfo.Parts = make([]shared.Part, 0)
	var Height = 0

	for _, episode := range biliInfo.Data.UgcSeason.Sections[0].Episodes {
		videoInfo.Parts = append(videoInfo.Parts, shared.Part{
			Url:   fmt.Sprintf("https://www.bilibili.com/video/%s", episode.Bvid),
			Title: episode.Title,
			// Thumbnail: biliInfo.Data.UgcSeason.Sections[0].Episodes[index].Arc.Pic,
		})

		min := func(a, b int) int {
			if a < b {
				return a
			} else {
				return b
			}
		}(episode.Page.Dimension.Width, episode.Page.Dimension.Height)
		if Height < int(min) {
			Height = int(min)
		}
	}

	// 如果只有一个 就用主标题
	if len(biliInfo.Data.UgcSeason.Sections[0].Episodes) == 1 {
		videoInfo.Parts[0].Title = videoInfo.WorkDirName
	}
	// TODO
	// videoInfo.Qualities = utils.GetQualities(Height)

	return videoInfo

}

func downloadReq(SESSDATA string) (*http.Request, error) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0")
	req.Header.Set("Referer", "https://www.bilibili.com")
	if SESSDATA != "" {
		req.Header.Set("Cookie", "SESSDATA="+SESSDATA)
	}

	return req, nil
}
