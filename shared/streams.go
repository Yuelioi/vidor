package shared

import "time"

type Format struct {
	IDtag         int    // 标签ID
	URL           string // 链接
	MimeType      string // video/mp4...
	ContentLength int64  // 内容长度
	DurationMs    int    // 时长

	// 仅视频
	FPS            int
	Width          int
	Height         int
	Bitrate        int // 码率
	AverageBitrate int // 平均码率

	// 仅音频
	AudioSampleRate string
}

type Stream struct {
	ID              string // youtubeID bilibiliID...
	SessionId       string // biliCID...
	URL             string
	Title           string
	Description     string
	Author          string
	ChannelID       string
	Views           int
	Duration        time.Duration
	PublishDate     time.Time
	Formats         []Format
	DASHManifestURL string // URI of the DASH manifest file
	HLSManifestURL  string // URI of the HLS manifest file
}

type Thumbnail struct {
	URL    string
	Label  string
	Width  uint
	Height uint
}

type CaptionTrack struct {
	BaseURL        string
	Name           string
	LanguageCode   string
	Kind           string
	IsTranslatable bool
}
