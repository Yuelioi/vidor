package bilibili

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/Yuelioi/vidor/shared"
	"github.com/Yuelioi/vidor/utils"
)

// 获取合集信息
func processSeasonData(data biliPlayListData, aid int, bvid_input string) biliBaseParams {
	episodes := data.UgcSeason.Sections[0].Episodes
	index := 0
	for epid, epi := range episodes {
		if epi.Aid == aid || epi.Bvid == bvid_input {
			index = epid
		}
	}
	workDirname := utils.SanitizeFileName(data.UgcSeason.Title)
	coverName := fmt.Sprintf("%02d_%s.jpg", index+1, utils.SanitizeFileName(workDirname))

	return biliBaseParams{
		index:    index,
		cid:      episodes[index].CID,
		bvid:     data.UgcSeason.Sections[0].Episodes[index].Bvid,
		title:    episodes[index].Title,
		author:   data.Owner.Name,
		pubDate:  episodes[index].Arc.Pubdate,
		duration: episodes[index].Arc.Duration,

		coverURL:     data.UgcSeason.Cover,
		coverName:    coverName,
		thumbnailURL: episodes[index].Arc.Pic,
	}
}

// 获取普通分p信息
func processPagesData(bpi biliPlayListData, index int) biliBaseParams {
	var title, thumbnailURL string

	if len(bpi.Pages) == 1 {
		title = bpi.Title
	} else {
		title = bpi.Pages[index].Title
	}

	if bpi.Pages[index].Thumbnail != "" {
		thumbnailURL = bpi.Pages[index].Thumbnail
	} else {
		thumbnailURL = bpi.Pic
	}

	workDirname := bpi.Title

	return biliBaseParams{
		index:        index,
		cid:          bpi.Pages[index].CID,
		bvid:         bpi.BVID,
		title:        utils.SanitizeFileName(title),
		author:       bpi.Owner.Name,
		coverURL:     bpi.Pic,
		coverName:    utils.SanitizeFileName(workDirname) + ".jpg",
		thumbnailURL: thumbnailURL,
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

// B站视频列表信息转为通用的列表信息
func biliPlaylistInfoToPlaylistInfo(biliInfo biliPlaylistInfo) shared.PlaylistInfo {
	var videoInfo shared.PlaylistInfo

	videoInfo.URL = fmt.Sprintf("https://www.bilibili.com/video/%s", biliInfo.Data.BVID)
	videoInfo.WorkDirName = utils.SanitizeFileName(biliInfo.Data.Title)
	videoInfo.Author = biliInfo.Data.Owner.Name
	videoInfo.Description = biliInfo.Data.Desc
	videoInfo.Cover = biliInfo.Data.Pic

	videoInfo.StreamInfos = make([]shared.StreamInfo, 0)
	return videoInfo
}

// B站分P视频列表信息转为通用的列表信息
func biliPageToPlaylistInfo(bvid string, biliInfo biliPlaylistInfo) shared.PlaylistInfo {
	videoInfo := biliPlaylistInfoToPlaylistInfo(biliInfo)

	videoInfo.Author = biliInfo.Data.Owner.Name
	videoInfo.WorkDirName = utils.SanitizeFileName(biliInfo.Data.Title)
	videoInfo.PubDate = time.Unix(int64(biliInfo.Data.PubDate), 0)
	videoInfo.StreamInfos = make([]shared.StreamInfo, 0)

	for _, page := range biliInfo.Data.Pages {

		videoInfo.StreamInfos = append(videoInfo.StreamInfos, shared.StreamInfo{
			ID:        bvid,
			SessionId: fmt.Sprint(page.CID),
			Name:      page.Title,
			Videos: []shared.Format{
				{
					IDtag:   9999,
					Quality: "尚未解析", Selected: true,
				},
			},
			Audios: []shared.Format{
				{
					IDtag:   9999,
					Quality: "尚未解析", Selected: true,
				},
			},
			Captions:   []shared.CaptionTrack{{Name: "需要解析"}},
			Thumbnails: []shared.Thumbnail{{URL: page.Thumbnail}},
		})

	}

	return videoInfo
}

// B站合集视频列表信息转为通用的列表信息
func biliSeasonToPlaylistInfo(biliInfo biliPlaylistInfo) shared.PlaylistInfo {
	videoInfo := biliPlaylistInfoToPlaylistInfo(biliInfo)

	videoInfo.Author = biliInfo.Data.Owner.Name
	videoInfo.WorkDirName = utils.SanitizeFileName(biliInfo.Data.Title)
	videoInfo.PubDate = time.Unix(int64(biliInfo.Data.PubDate), 0)
	videoInfo.StreamInfos = make([]shared.StreamInfo, 0)

	for _, episode := range biliInfo.Data.UgcSeason.Sections[0].Episodes {

		videoInfo.StreamInfos = append(videoInfo.StreamInfos, shared.StreamInfo{
			ID:        episode.Bvid,
			SessionId: fmt.Sprint(episode.CID),
			Name:      episode.Title,
			Videos: []shared.Format{
				{
					IDtag:   9999,
					Quality: "尚未解析", Selected: true,
				},
			},
			Audios: []shared.Format{
				{
					IDtag:   9999,
					Quality: "尚未解析", Selected: true,
				},
			},
			Captions:   []shared.CaptionTrack{{Name: "需要解析"}},
			Thumbnails: []shared.Thumbnail{{URL: episode.Arc.Pic}},
		})

	}

	return videoInfo
}
