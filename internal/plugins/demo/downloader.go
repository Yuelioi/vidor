package demo

import (
	"regexp"

	"github.com/Yuelioi/vidor/shared"
)

type Downloader struct {
}

func New(notice shared.Notice) *Downloader {
	return &Downloader{}
}

func (bd *Downloader) PluginMeta() shared.PluginMeta {
	return shared.PluginMeta{
		Name:   "demo",
		Regexs: []*regexp.Regexp{},
	}
}

func (bd *Downloader) ShowInfo(link string, config shared.Config, callback shared.Callback) (*shared.PlaylistInfo, error) {
	var playList shared.PlaylistInfo
	return &playList, nil
}

func (bd *Downloader) GetMeta(config shared.Config, part *shared.Part, callback shared.Callback) error {
	return nil
}

func (bd *Downloader) DownloadThumbnail(part *shared.Part, callback shared.Callback) error {
	return nil
}

func (bd *Downloader) DownloadVideo(part *shared.Part, callback shared.Callback) error {
	return nil
}

func (bd *Downloader) DownloadAudio(part *shared.Part, callback shared.Callback) error {
	return nil
}

func (bd *Downloader) DownloadSubtitle(part *shared.Part, callback shared.Callback) error {
	return nil
}
func (bd *Downloader) Combine(ffmpegPath string, part *shared.Part) error {
	return nil
}
func (bd *Downloader) Clear(part *shared.Part, callback shared.Callback) error {
	return nil
}
func (bd *Downloader) StopDownload(part *shared.Part, callback shared.Callback) error {

	return nil
}
func (bd *Downloader) PauseDownload(part *shared.Part, callback shared.Callback) error {
	return nil
}
