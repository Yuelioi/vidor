package main

import "time"

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
			Videos []Video `json:"video"`
			Audios []Audio `json:"audio"`
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
	index int // 0起始
	cid   int
	bvid  string

	title  string
	author string

	pubDate  int
	duration int

	coverURL     string
	coverName    string
	thumbnailURL string
}

// 视频下载参数
type biliDownloadParams struct {
	videoURL string
	audioURL string
	index    int
}

// 视频封面参数
type biliThumbnailParams struct {
	thumbnailURL string
	coverPath    string
	coverURL     string
}

type Video struct {
	ID        int      `json:"id"`
	BaseURL   string   `json:"baseURL"`
	BaseUrl   string   `json:"base_url"`
	BackupURL []string `json:"backupURL"`
	BackupUrl []string `json:"backup_url"`
	Codecs    string   `json:"codecs"`
	Width     int      `json:"width"`
	Height    int      `json:"height"`
	FrameRate string   `json:"frameRate"`
}

type Audio struct {
	ID        int      `json:"id"`
	BaseURL   string   `json:"baseURL"`
	BaseUrl   string   `json:"base_url"`
	BackupURL []string `json:"backupURL"`
	BackupUrl []string `json:"backup_url"`
	Bandwidth int      `json:"bandwidth"`
	Codecs    string   `json:"codecs"`
}

// *------------------------------

type Part struct {
	URL         string // 链接
	TaskID      string // 唯一标识 ;需创建task时初始化
	WorkDirName string // 工作文件夹名 ;需创建task时初始化
	DownloadDir string // 下载文件夹完整路径 ;需创建task时初始化
	MagicName   string // 下载文件名 不带后缀
	Path        string // 下载文件完整路径

	Index       int       // 所在父级索引 0开始
	Author      string    // 作者
	Title       string    // 标题
	Description string    // 描述
	Width       int       // 宽度
	Height      int       // 高度
	Size        int       // 字节数
	Duration    int       // 持续时间 秒
	PubDate     time.Time // 发布日期

	Thumbnail string        // 封面
	Video     StreamQuality // 质量标签
	Audio     StreamQuality // 质量标签
	Subtitle  string        // todo

	State           string    // 状态
	Status          string    // 进度
	CreatedAt       time.Time // 任务创建日期
	DownloadPercent int       // 下载百分比
	DownloadSpeed   string    // 下载速度
}

type StreamQuality struct {
	IDtag int
	Label string
}
