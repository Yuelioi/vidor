package plugins

// import (
// 	"context"
// 	"encoding/base64"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"net/url"
// 	"os"
// 	"path/filepath"
// 	"sort"
// 	"time"

// 	"github.com/Yuelioi/vidor/utils"

// 	"github.com/Yuelioi/vidor/shared"

// 	"github.com/chromedp/chromedp"
// )

// type Config struct {
// 	Request struct {
// 		Files struct {
// 			Dash struct {
// 				CDNs struct {
// 					AkfireInterconnectQuic struct {
// 						AVCURL string `json:"avc_url"`
// 						URL    string `json:"url"`
// 					} `json:"akfire_interconnect_quic"`
// 				} `json:"cdns"`
// 			} `json:"dash"`
// 			HLS struct {
// 				CDNs struct {
// 					AkfireInterconnectQuic struct {
// 						AVCURL string `json:"avc_url"`
// 						URL    string `json:"url"`
// 					} `json:"akfire_interconnect_quic"`
// 				} `json:"cdns"`
// 			} `json:"hls"`
// 		} `json:"files"`
// 	} `json:"request"`
// 	Video struct {
// 		ID       int    `json:"id"`
// 		Width    int    `json:"width"`
// 		Height   int    `json:"height"`
// 		URL      string `json:"url"`
// 		Title    string `json:"title"`
// 		Duration int    `json:"duration"`
// 		Thumbs   struct {
// 			Size640  string `json:"640"`
// 			Size960  string `json:"960"`
// 			Size1280 string `json:"1280"`
// 			Base     string `json:"base"`
// 		} `json:"thumbs"`
// 		Owner struct {
// 			ID    int    `json:"id"`
// 			Name  string `json:"name"`
// 			Img   string `json:"img"`
// 			Img2x string `json:"img_2x"`
// 			URL   string `json:"url"`
// 		} `json:"owner"`
// 	} `json:"video"`
// }

// type Segment struct {
// 	Start float64 `json:"start"`
// 	End   float64 `json:"end"`
// 	URL   string  `json:"url"`
// 	Size  int     `json:"size"`
// }

// // Video represents the video information.
// type vmVideo struct {
// 	ID                 string    `json:"id"`
// 	AvgID              string    `json:"avg_id"`
// 	BaseURL            string    `json:"base_url"`
// 	Format             string    `json:"format"`
// 	MimeType           string    `json:"mime_type"`
// 	Codecs             string    `json:"codecs"`
// 	Bitrate            int       `json:"bitrate"`
// 	AvgBitrate         int       `json:"avg_bitrate"`
// 	Duration           float64   `json:"duration"`
// 	Framerate          int       `json:"framerate"`
// 	Width              int       `json:"width"`
// 	Height             int       `json:"height"`
// 	MaxSegmentDuration int       `json:"max_segment_duration"`
// 	InitSegment        string    `json:"init_segment"`
// 	IndexSegment       string    `json:"index_segment"`
// 	Segments           []Segment `json:"segments"`
// 	AudioProvenance    int       `json:"AudioProvenance"`
// }

// type Audio struct {
// 	ID                 string    `json:"id"`
// 	AvgID              string    `json:"avg_id"`
// 	BaseURL            string    `json:"base_url"`
// 	Format             string    `json:"format"`
// 	MimeType           string    `json:"mime_type"`
// 	Codecs             string    `json:"codecs"`
// 	Bitrate            int       `json:"bitrate"`
// 	AvgBitrate         int       `json:"avg_bitrate"`
// 	Duration           float64   `json:"duration"`
// 	Channels           int       `json:"channels"`
// 	SampleRate         int       `json:"sample_rate"`
// 	MaxSegmentDuration int       `json:"max_segment_duration"`
// 	InitSegment        string    `json:"init_segment"`
// 	Segments           []Segment `json:"segments"`
// }

// type Clip struct {
// 	ClipID  string    `json:"clip_id"`
// 	BaseURL string    `json:"base_url"`
// 	Videos  []vmVideo `json:"video"`
// 	Audios  []Audio   `json:"audio"`
// }

// type VimeoDownloader struct {
// 	Client      *http.Client
// 	videoConfig *Config
// }

// func NewVimeo() *VimeoDownloader {
// 	return &VimeoDownloader{}
// }

// func (vm *VimeoDownloader) GetClient(proxyURL string, useProxy bool) interface{} {
// 	client, err := utils.GetClient(proxyURL, useProxy)
// 	if err != nil {
// 		return nil
// 	}
// 	return client
// }

// func (vd *VimeoDownloader) ShowInfo(link string, config shared.Config) (*shared.PlaylistInfo, error) {

// 	vd.Client = vd.GetClient(config.ProxyURL, config.UseProxy).(*http.Client)
// 	pli, videoConfig, err := start(*vd.Client, link, config)
// 	vd.videoConfig = videoConfig

// 	return pli, err
// }

// func (vd *VimeoDownloader) GetMeta(ctx context.Context, part *shared.Part) error {

// 	clipMainURL := vd.videoConfig.Request.Files.Dash.CDNs.AkfireInterconnectQuic.URL
// 	clip, err := fetchClip(*vd.Client, clipMainURL)
// 	if err != nil {
// 		return err
// 	}

// 	baseURL, err := resolveURL(clipMainURL, clip.BaseURL)
// 	if err != nil {
// 		return err
// 	}
// 	clipBaseURL := baseURL

// 	part.Thumbnail = vd.videoConfig.Video.Thumbs.Size1280
// 	part.Title = vd.videoConfig.Video.Title
// 	part.Author = vd.videoConfig.Video.Owner.Name
// 	part.URL = vd.videoConfig.Video.URL

// 	taskDir := filepath.Join(part.DownloadDir, utils.SanitizeFileName(part.Title))
// 	filePureName := utils.SanitizeFileName(part.Title)
// 	part.DownloadDir = taskDir
// 	if err := utils.CreateDirs([]string{taskDir}); err != nil {
// 		return err
// 	}

// 	taskImagePath := filepath.Join(taskDir, filePureName+".jpg")
// 	if _, err = os.Stat(taskImagePath); os.IsNotExist(err) {
// 		_, err := utils.GetCover(vd.Client, part.Thumbnail, taskImagePath)
// 		if err != nil {
// 			return fmt.Errorf("缓存缩略图失败: %v", err.Error())
// 		}
// 	}

// 	bestVideo := getTargetVimeoVideo(clip.Videos, 0)
// 	bestAudio := getBestVimeoAudio(clip.Audios)

// 	// 先不要带后缀
// 	path_pure := filepath.Join(part.DownloadDir, filePureName)

// 	input_v := filepath.Join(part.DownloadDir, fmt.Sprintf("%s_temp.mp4", filePureName))
// 	input_a := filepath.Join(part.DownloadDir, fmt.Sprintf("%s_temp.mp3", filePureName))

// 	videosTemp, err := vd.downloadSegments(part, bestVideo.InitSegment, bestVideo.Segments, clipBaseURL, bestVideo.BaseURL, path_pure, "m4s")
// 	if err != nil {
// 		return err
// 	}

// 	utils.CombineSegments(videosTemp, input_v, map[string]interface{}{"v": 1, "a": 0})

// 	part.Status = "音频片段下载中"
// 	audiosTemp, err := vd.downloadSegments(part, bestAudio.InitSegment, bestAudio.Segments, clipBaseURL, bestAudio.BaseURL, path_pure, "mp3")
// 	if err != nil {
// 		return err
// 	}
// 	part.Status = "音频片段合并中"
// 	utils.CombineSegments(audiosTemp, input_a, map[string]interface{}{"v": 0, "a": 1})

// 	if err = utils.ClearDirs(append(videosTemp, audiosTemp...)); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (vd *VimeoDownloader) DownloadVideo(ctx context.Context, part *shared.Part) error { return nil }
// func (vd *VimeoDownloader) DownloadAudio(ctx context.Context, part *shared.Part) error { return nil }
// func (vd *VimeoDownloader) DownloadThumbnail(ctx context.Context, part *shared.Part) error {
// 	return nil
// }
// func (vd *VimeoDownloader) DownloadSubtitle(ctx context.Context, part *shared.Part) error { return nil }
// func (vd *VimeoDownloader) Clear() error                                                  { return nil }

// func start(client http.Client, link string, config shared.Config) (*shared.PlaylistInfo, *Config, error) {

// 	pli := shared.PlaylistInfo{}

// 	masterURL, err := getVimeoConfigURL(link, config)
// 	if err != nil {
// 		log.Printf("获取配置URL时出错: %v", err)
// 		return nil, nil, err
// 	}

// 	// 读取配置
// 	videoConfig, err := masterConfig(client, masterURL)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	pli.Cover = videoConfig.Video.Thumbs.Size1280
// 	pli.WorkDirName = videoConfig.Video.Title
// 	pli.Author = videoConfig.Video.Owner.Name
// 	pli.URL = videoConfig.Video.URL
// 	// todo
// 	// pli.Qualities = utils.GetQualities(videoConfig.Video.Height)

// 	pli.StreamInfos = make([]shared.StreamInfo, 0)

// 	pli.StreamInfos = []shared.StreamInfo{{}}

// 	return &pli, videoConfig, nil
// }
// func getTargetVimeoVideo(videos []vmVideo, targetHeight int) vmVideo {
// 	// 从小到大排序, 应对竖版视频
// 	sort.Slice(videos, func(i, j int) bool {
// 		return videos[i].Height < videos[j].Height
// 	})

// 	var closestVideo vmVideo
// 	for _, video := range videos {
// 		if video.Height >= targetHeight {
// 			return video
// 		}
// 		closestVideo = video
// 	}
// 	return closestVideo
// }

// func getBestVimeoAudio(audios []Audio) Audio {
// 	var bestAudio Audio
// 	var currentRate = 0
// 	for _, audio := range audios {
// 		if audio.SampleRate > currentRate {
// 			bestAudio = audio
// 			currentRate = audio.SampleRate
// 		}
// 	}
// 	return bestAudio
// }

// // 获取所有信息
// func getVimeoConfigURL(videoURL string, config shared.Config) (string, error) {

// 	opts := append(chromedp.DefaultExecAllocatorOptions[:],
// 		chromedp.Flag("headless", true),
// 		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
// 		chromedp.Flag("proxy-bypass-list", "<-loopback>"),
// 	)

// 	if config.UseProxy {
// 		opts = append(opts,
// 			chromedp.ProxyServer(config.ProxyURL),
// 		)
// 	}

// 	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
// 	defer cancel()

// 	ctx, cancel = chromedp.NewContext(ctx)
// 	defer cancel()
// 	// 设置超时时间
// 	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
// 	defer cancel()

// 	var configURL string

// 	if err := chromedp.Run(ctx, chromedp.Navigate(videoURL),
// 		chromedp.WaitReady(`.player`, chromedp.ByQuery),
// 		chromedp.AttributeValue(`.js-player`, "data-config-url", &configURL, nil),
// 	); err != nil {
// 		return "", fmt.Errorf("未找到配置URL")
// 	}

// 	if configURL == "" {
// 		return "", fmt.Errorf("未找到配置URL")
// 	}

// 	return configURL, nil
// }

// func masterConfig(client http.Client, masterURL string) (*Config, error) {

// 	req, err := http.NewRequest("GET", masterURL, nil)
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 		return nil, err
// 	}

// 	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

// 	resp, err := client.Do(req)
// 	if err != nil {

// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	// 读取响应体
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Printf("Failed to read response body: %v", err)
// 		return nil, err
// 	}

// 	// 解析 JSON 数据
// 	var config *Config
// 	err = json.Unmarshal(body, &config)
// 	if err != nil {
// 		log.Printf("Failed to parse master.json: %v", err)
// 	}

// 	return config, nil
// }

// func fetchClip(client http.Client, link string) (*Clip, error) {
// 	req, err := http.NewRequest("GET", link, nil)
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 		return nil, err
// 	}

// 	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Printf("Failed to fetch master JSON")
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	// 读取响应体
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Printf("Failed to read response body: %v", err)
// 		return nil, err
// 	}

// 	// 解析 JSON 数据
// 	var clip *Clip
// 	err = json.Unmarshal(body, &clip)
// 	if err != nil {
// 		log.Printf("Failed to parse master.json: %v", err)
// 	}

// 	return clip, nil
// }

// func (vd *VimeoDownloader) downloadSeg(link string, initSegmentData []byte, filename string) {

// 	// 创建请求
// 	req, err := http.NewRequest("GET", link, nil)
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 		return
// 	}

// 	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

// 	res, err := vd.Client.Do(req)
// 	if err != nil {
// 		fmt.Println("Error sending request:", err)
// 		return
// 	}
// 	defer res.Body.Close()

// 	file, err := os.Create(filename)
// 	if err != nil {
// 		fmt.Println("Error creating file:", err)
// 		return
// 	}
// 	defer file.Close()

// 	// 写入初始化段内容
// 	_, err = file.Write(initSegmentData)
// 	if err != nil {
// 		fmt.Println("Error writing init segment to file:", err)
// 		return
// 	}

// 	_, err = io.Copy(file, res.Body)
// 	if err != nil {
// 		fmt.Println("Error writing to file:", err)
// 		return
// 	}

// }

// func (vd *VimeoDownloader) downloadSegments(part *shared.Part, initSegment string, segs []Segment, clipBaseURL, baseURL, path, ext string) ([]string, error) {
// 	// 下载初始化段

// 	var temps []string
// 	// 读取初始化段内容
// 	initSegData, err := base64.StdEncoding.DecodeString(initSegment)
// 	if err != nil {
// 		fmt.Println("Error decoding init segment:", err)
// 		return temps, err
// 	}

// 	// 下载视频片段并将初始化段写入每个片段文件
// 	for index, seg := range segs {
// 		absURL, _ := resolveURL(clipBaseURL, baseURL)
// 		absURL, _ = resolveURL(absURL, seg.URL)
// 		realLink, _ := url.QueryUnescape(absURL)
// 		filepath := fmt.Sprintf("%s_%d.%s", path, index, ext)
// 		temps = append(temps, filepath)
// 		vd.downloadSeg(realLink, initSegData, filepath)
// 		part.DownloadPercent = int(float64(index)/float64(len(segs))) * 100
// 		fmt.Printf("%+v", part.DownloadPercent)
// 	}

// 	return temps, nil
// }

// func resolveURL(baseURL, relative string) (string, error) {
// 	// Parse the m3u8 URL
// 	baseUrl, err := url.Parse(baseURL)
// 	if err != nil {
// 		return "", fmt.Errorf("error parsing URL: %v", err)
// 	}

// 	// Resolve the relative segment URL against the m3u8 base URL
// 	segmentAbsURL := baseUrl.ResolveReference(&url.URL{Path: relative})
// 	return segmentAbsURL.String(), nil
// }
