package bilibili

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

type userInfo struct {
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

	coverUrl     string
	coverName    string
	thumbnailUrl string
}

// 视频下载参数
type biliDownloadParams struct {
	videoUrl string
	audioUrl string
	index    int
}

// 视频封面参数
type biliThumbnailParams struct {
	thumbnailUrl string
	coverPath    string
	coverUrl     string
}

type Video struct {
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

type Audio struct {
	BaseUrl   string   `json:"baseUrl"`
	BaseURL   string   `json:"base_url"`
	BackupUrl []string `json:"backupUrl"`
	BackupURL []string `json:"backup_url"`
	Bandwidth int      `json:"bandwidth"`
	Codecs    string   `json:"codecs"`
}
