package shared

type Format struct {
	IDtag         int      // ★ 标签ID
	Quality       string   // ★ 质量标签
	Selected      bool     // ★ 是否选中当前格式
	URL           string   // 链接
	MimeType      string   // video/mp4...
	ContentLength int64    // 内容长度
	DurationMs    int      // 时长
	Codecs        []string // 编码类型

	// 图片+视频

	Width  int
	Height int

	// 仅视频

	FPS            int // FPS
	Bitrate        int // 码率
	AverageBitrate int // 平均码率

	// 仅音频

	AudioSampleRate string // 音频码率
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
