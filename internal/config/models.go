package config

// 此乃系统配置
type Config struct {
	BaseDir          string // 配置所在文件夹
	Theme            string `json:"theme"`
	ScaleFactor      int    `json:"scale_factor"`
	ProxyURL         string `json:"proxy_url"`
	UseProxy         bool   `json:"use_proxy"`
	MagicName        string `json:"magic_name"`
	DownloadDir      string `json:"download_dir"`
	DownloadVideo    bool   `json:"download_video"`
	DownloadAudio    bool   `json:"download_audio"`
	DownloadSubtitle bool   `json:"download_subtitle"`
	DownloadCombine  bool   `json:"download_combine"`
	DownloadLimit    int    `json:"download_limit"`
}
