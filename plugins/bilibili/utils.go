package bilibili

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"

	"github.com/Yuelioi/vidor/shared"
	"github.com/Yuelioi/vidor/utils"
)

func (bd *Downloader) getUserStates(sessdata string) {
	apiUrl := apiURL + "/x/vip/web/user/info"

	data, err := doBiliReq(*bd.Client, apiUrl, sessdata)
	if err != nil {
		return
	}

	var userInfo userInfo
	err = json.Unmarshal(data, &userInfo)
	if err != nil {
		return
	}

	if userInfo.Data.VIPStatus == 1 {
		bd.userState = Vip
		return
	}
	bd.userState = Login
}

func (bd *Downloader) download(part *shared.Part, link, ext string, callback shared.Callback) error {
	path := filepath.Join(part.DownloadDir, fmt.Sprintf("%s_temp.%s", part.MagicName, ext))

	req, err := downloadReq(bd.sessdata)
	if err != nil {

		return err
	}
	mediaUrl, err := url.Parse(link)
	if err != nil {

		return err
	}
	req.URL = mediaUrl
	utils.ReqWriter(bd.ctx, bd.Client, req, part, path, callback)
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
		index:    index,
		cid:      episodes[index].CID,
		bvid:     data.UgcSeason.Sections[0].Episodes[index].Bvid,
		title:    episodes[index].Title,
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
		index:        index,
		cid:          bpi.Pages[index].CID,
		bvid:         bpi.BVID,
		title:        utils.SanitizeFileName(title),
		author:       bpi.Owner.Name,
		coverUrl:     bpi.Pic,
		coverName:    utils.SanitizeFileName(workDirname) + ".jpg",
		thumbnailUrl: thumbnailUrl,
	}
}

// 音频取最大的, 反正不大
func getTargetAudio(audios []Audio) *Audio {

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
func getTargetVideo(targeID int, videos []Video) *Video {
	// 倒序, 向上取
	sort.Slice(videos, func(i, j int) bool {
		return videos[i].ID > videos[j].ID
	})

	var closestVideo *Video
	for _, video := range videos {
		if video.ID <= targeID {
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

	body, err := doBiliReq(client, fmt.Sprintf(`%s/x/web-interface/view?aid=%d&bvid=%s`, apiURL, aid, bvid), sessdata)
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
	body, err := doBiliReq(client, fmt.Sprintf("%s/x/player/wbi/playurl?bvid=%s&cid=%d&fnval=4048&fourk=1", apiURL, bvid, cid), sessdata)
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

// B站视频列表信息转为通用的列表信息
func biliPlaylistInfoToPlaylistInfo(biliInfo biliPlaylistInfo) shared.PlaylistInfo {
	var videoInfo shared.PlaylistInfo

	videoInfo.Url = fmt.Sprintf("https://www.bilibili.com/video/%s", biliInfo.Data.BVID)
	videoInfo.WorkDirName = utils.SanitizeFileName(biliInfo.Data.Title)
	videoInfo.Author = biliInfo.Data.Owner.Name
	videoInfo.Description = biliInfo.Data.Desc
	videoInfo.Cover = biliInfo.Data.Pic

	videoInfo.StreamInfos = make([]shared.StreamInfo, 0)
	return videoInfo
}

// B站分P视频列表信息转为通用的列表信息
func biliPageToPlaylistInfo(biliInfo biliPlaylistInfo) shared.PlaylistInfo {
	videoInfo := biliPlaylistInfoToPlaylistInfo(biliInfo)

	for _, page := range biliInfo.Data.Pages {

		videoInfo.StreamInfos = append(videoInfo.StreamInfos, shared.StreamInfo{
			TaskID: page.Title,
		})

	}
	// 如果只有一个 就用主标题
	if len(biliInfo.Data.Pages) == 1 {
		videoInfo.StreamInfos[0].TaskID = videoInfo.WorkDirName
	}
	return videoInfo
}

// B站合集视频列表信息转为通用的列表信息
func biliSeasonToPlaylistInfo(biliInfo biliPlaylistInfo) shared.PlaylistInfo {
	videoInfo := biliPlaylistInfoToPlaylistInfo(biliInfo)

	for _, episode := range biliInfo.Data.UgcSeason.Sections[0].Episodes {
		videoInfo.StreamInfos = append(videoInfo.StreamInfos, shared.StreamInfo{
			TaskID: episode.Arc.Pic,

			// Url:      fmt.Sprintf("https://www.bilibili.com/video/%s", episode.Bvid),
			// Title:    episode.Title,
			// Duration: episode.Arc.Duration,

			// Thumbnail: biliInfo.Data.UgcSeason.Sections[0].Episodes[index].Arc.Pic,
		})

	}

	// 如果只有一个 就用主标题
	// if len(biliInfo.Data.UgcSeason.Sections[0].Episodes) == 1 {
	// 	videoInfo.StreamInfos[0].Title = videoInfo.WorkDirName
	// }

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
