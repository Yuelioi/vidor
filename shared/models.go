package shared

import (
	"regexp"
	"time"
)

// 软件配置信息
type Config struct {
	Theme       string
	ScaleFactor int

	ProxyURL         string
	UseProxy         bool
	FFMPEG           string
	MagicName        string
	DownloadDir      string
	DownloadVideo    bool
	DownloadAudio    bool
	DownloadSubtitle bool
	DownloadCombine  bool
	DownloadLimit    int

	SESSDATA string
}

// 任务片段信息
type Part struct {
	TaskID      string // 唯一标识 ;需创建task时初始化
	WorkDirName string // 工作文件夹名 ;需创建task时初始化
	DownloadDir string // 下载文件夹完整路径 ;需创建task时初始化
	MagicName   string // 下载文件名 不带后缀
	Path        string // 下载文件完整路径

	Url         string    // 链接
	Index       int       // 所在父级索引 0开始
	Author      string    // 作者
	Title       string    // 标题
	Description string    // 描述
	Thumbnail   string    // 封面
	Quality     string    // 质量标签
	Width       int       // 宽度
	Height      int       // 高度
	Size        int       // 字节数
	Duration    int       // 持续时间 秒
	CreatedAt   time.Time // 任务创建日期
	PubDate     time.Time // 发布日期

	State           string // 状态
	Status          string // 进度
	DownloadPercent int    // 下载百分比
	DownloadSpeed   string // 下载速度
}

// Home页面搜索展示所需信息
type PlaylistInfo struct {
	Url         string       // 下载链接
	Cover       string       // 封面
	WorkDirName string       // 工作路径名 一般是视频标题/合集标题 后续用来创建下载文件夹
	Author      string       // 作者
	Description string       // 描述
	PubDate     time.Time    // 发布日期
	StreamInfos []StreamInfo // 分段合集
}

type StreamInfo struct {
	TaskID     string
	Thumbnails []Thumbnail
	Videos     []Stream
	Audios     []Stream
	Captions   []CaptionTrack
}

type status struct {
	Queue       string
	Downloading string
	Pause       string
	Finished    string

	Stopped              string
	GettingMetadata      string
	DownloadingVideo     string
	DownloadingAudio     string
	DownloadingSubtitle  string
	DownloadingThumbnail string
	Merging              string
	Failed               string
	Unknown              string
}

// 流媒体质量
type StreamQuality struct {
	ID    int    // 流媒体ID
	Label string // 流媒体标签
}

// 插件信息
type PluginMeta struct {
	Name   string
	Type   string // System/ThirdPart
	Regexs []*regexp.Regexp
	Impl   func(notice Notice) Downloader
}

// 任务状态
var TaskStatus = status{
	// 任务本体状态
	Queue:       "队列中",
	Downloading: "下载中",
	Pause:       "已暂停",
	Finished:    "已完成",

	// 任务下载状态
	Stopped:              "已取消",
	GettingMetadata:      "获取元数据",
	DownloadingThumbnail: "下载封面",
	DownloadingVideo:     "下载视频",
	DownloadingAudio:     "下载音频",
	DownloadingSubtitle:  "下载字幕",
	Merging:              "合并中",
	Failed:               "下载失败",
	Unknown:              "未知状态",
}
