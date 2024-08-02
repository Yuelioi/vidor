package shared

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type Downloader interface {
	PluginMeta() PluginMeta // 获取插件信息

	// 主页搜索展示信息
	//
	// 必须获取 PlaylistInfo 所有信息
	// 其中 StreamInfo ★字段需要获取,★★字段需要初始化
	ShowInfo(link string, callback Callback) (*PlaylistInfo, error)
	ParsePlaylist(*PlaylistInfo) (*PlaylistInfo, error)

	Download(part *Part, callback Callback) error

	PauseDownload(part *Part, callback Callback) error // 暂停下载
	StopDownload(part *Part, callback Callback) error  // 停止下载
}

type InputInfo struct {
	Link     string
	Callback Callback
}

type DownloadArgs struct {
	Part     *Part
	Callback Callback
}
type CombineArgs struct {
	ffmpegPath string
	Part       *Part
}

type DownloaderRPC struct{ client *rpc.Client }

// ---------------------------------- 客户端信息 ------------------------------

func (g *DownloaderRPC) PluginMeta() PluginMeta {
	var resp = new(PluginMeta)
	if err := g.client.Call("Plugin.PluginMeta", new(struct{}), &resp); err != nil {
		return PluginMeta{}
	}
	return *resp
}

func (g *DownloaderRPC) ShowInfo(link string, callback Callback) (*PlaylistInfo, error) {
	var resp PlaylistInfo
	err := g.client.Call("Plugin.ShowInfo", InputInfo{Link: link, Callback: callback}, &resp)
	return &resp, err
}

func (g *DownloaderRPC) GetMeta(part *Part, callback Callback) (*Part, error) {
	var resp Part
	err := g.client.Call("Plugin.GetMeta", DownloadArgs{
		Part:     part,
		Callback: callback,
	}, &resp)
	return &resp, err
}

func (g *DownloaderRPC) DownloadThumbnail(part *Part, callback Callback) error {
	args := &DownloadArgs{Part: part, Callback: callback}
	return g.client.Call("Plugin.DownloadThumbnail", args, new(struct{}))
}
func (g *DownloaderRPC) DownloadVideo(part *Part, callback Callback) error {
	args := &DownloadArgs{Part: part, Callback: callback}
	return g.client.Call("Plugin.DownloadVideo", args, new(struct{}))
}
func (g *DownloaderRPC) DownloadAudio(part *Part, callback Callback) error {
	args := &DownloadArgs{Part: part, Callback: callback}
	return g.client.Call("Plugin.DownloadAudio", args, new(struct{}))
}
func (g *DownloaderRPC) DownloadSubtitle(part *Part, callback Callback) error {
	args := &DownloadArgs{Part: part, Callback: callback}
	return g.client.Call("Plugin.DownloadSubtitle", args, new(struct{}))
}

func (g *DownloaderRPC) StopDownload(part *Part, callback Callback) error {
	args := &DownloadArgs{Part: part, Callback: callback}
	return g.client.Call("Plugin.StopDownload", args, new(struct{}))
}
func (g *DownloaderRPC) Combine(ffmpegPath string, part *Part) error {
	args := &CombineArgs{ffmpegPath: ffmpegPath, Part: part}
	return g.client.Call("Plugin.Combine", args, new(struct{}))
}
func (g *DownloaderRPC) Clear(part *Part, callback Callback) error {
	args := &DownloadArgs{Part: part, Callback: callback}
	return g.client.Call("Plugin.Clear", args, new(struct{}))
}

// ---------------------------------- 服务端信息 ------------------------------

type DownloaderRPCServer struct {
	Impl Downloader
}

func (s *DownloaderRPCServer) PluginMeta(args struct{}, resp *PluginMeta) error {
	*resp = s.Impl.PluginMeta()
	return nil
}

func (s *DownloaderRPCServer) ShowInfo(args InputInfo, resp *PlaylistInfo) error {
	info, err := s.Impl.ShowInfo(args.Link, args.Callback)
	if err != nil {
		return err
	}
	*resp = *info
	return nil
}

func (s *DownloaderRPCServer) Download(args *DownloadArgs, resp *struct{}) error {
	return s.Impl.Download(args.Part, args.Callback)
}

func (s *DownloaderRPCServer) StopDownload(args *DownloadArgs, resp *struct{}) error {
	return s.Impl.StopDownload(args.Part, args.Callback)
}

type DownloaderRPCPlugin struct {
	Impl Downloader
}

func (p *DownloaderRPCPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &DownloaderRPCServer{Impl: p.Impl}, nil
}

func (DownloaderRPCPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &DownloaderRPC{client: c}, nil
}
