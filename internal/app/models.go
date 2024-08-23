package app

import (
	"time"
)

type StreamQuality struct {
	IDtag int
	Label string
}

// 任务片段信息
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
