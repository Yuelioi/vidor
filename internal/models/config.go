package models

type SystemConfig struct {
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

type PluginConfig struct {
	ID       string            `json:"id"`
	Enable   bool              `json:"enable"` // 建立连接 (Run)
	Settings map[string]string `json:"settings"`
}
