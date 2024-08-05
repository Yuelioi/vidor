package shared

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type Downloader interface {
	PluginMeta() PluginMeta // 获取插件信息

	Show(link string, callback Callback) (*PlaylistInfo, error) // 主页搜索展示信息
	Parse(*PlaylistInfo) (*PlaylistInfo, error)                 // 解析
	Do(part *Part, callback Callback) error                     // 下载封面/图片工作
	Cancel(part *Part, callback Callback) error                 // 取消下载

}

type InputInfo struct {
	Link     string
	Callback Callback
}

type DownloadArgs struct {
	Part     *Part
	Callback Callback
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

func (g *DownloaderRPC) Show(link string, callback Callback) (*PlaylistInfo, error) {
	var resp PlaylistInfo
	err := g.client.Call("Plugin.ShowInfo", InputInfo{Link: link, Callback: callback}, &resp)
	return &resp, err
}

func (g *DownloaderRPC) Parse(part *Part, callback Callback) error {
	args := &DownloadArgs{Part: part, Callback: callback}
	return g.client.Call("Plugin.Parse", args, new(struct{}))
}

func (g *DownloaderRPC) Do(part *Part, callback Callback) error {
	args := &DownloadArgs{Part: part, Callback: callback}
	return g.client.Call("Plugin.Do", args, new(struct{}))
}
func (g *DownloaderRPC) Cancel(part *Part, callback Callback) error {
	args := &DownloadArgs{Part: part, Callback: callback}
	return g.client.Call("Plugin.Cancel", args, new(struct{}))
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
	info, err := s.Impl.Show(args.Link, args.Callback)
	if err != nil {
		return err
	}
	*resp = *info
	return nil
}

<<<<<<< HEAD
func (s *DownloaderRPCServer) Download(args *DownloadArgs, resp *struct{}) error {
	return s.Impl.Download(args.Part, args.Callback)
}

func (s *DownloaderRPCServer) StopDownload(args *DownloadArgs, resp *struct{}) error {
	return s.Impl.StopDownload(args.Part, args.Callback)
=======
func (s *DownloaderRPCServer) Do(args *DownloadArgs, resp *struct{}) error {
	return s.Impl.Do(args.Part, args.Callback)
}

func (s *DownloaderRPCServer) Cancel(args *DownloadArgs, resp *struct{}) error {
	return s.Impl.Cancel(args.Part, args.Callback)
>>>>>>> 5b0913086fa11fa6c3bc8737a0a8bb8be8541fdd
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
