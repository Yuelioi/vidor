package shared

import (
	"context"
	// "github.com/hashicorp/go-plugin"
)

type Downloader interface {
	PluginMeta() PluginMeta // 获取插件信息

	Show(context.Context, string) (*PlaylistInfo, error)         // 主页搜索展示信息
	Parse(context.Context, *PlaylistInfo) (*PlaylistInfo, error) // 解析
	Do(context.Context, *Part) error                             // 下载封面/图片工作
	Cancel(context.Context, *Part) error                         // 取消下载
}

type InputInfo struct {
	Ctx  context.Context
	Link string
}

type DownloadArgs struct {
	Ctx  context.Context
	Part *Part
}

// type DownloaderRPC struct{ client *rpc.Client }

// ---------------------------------- TODO 插件系统 晚点搞 ------------------------------

// ---------------------------------- 客户端信息 ------------------------------

// func (g *DownloaderRPC) Init(ctx context.Context, config Config) Downloader {
// 	var resp = new(Downloader)
// 	if err := g.client.Call("Plugin.Init", config, &resp); err != nil {
// 		return nil
// 	}
// 	return *resp
// }

// func (g *DownloaderRPC) PluginMeta() PluginMeta {
// 	var resp = new(PluginMeta)
// 	if err := g.client.Call("Plugin.PluginMeta", new(struct{}), &resp); err != nil {
// 		return PluginMeta{}
// 	}
// 	return *resp
// }

// func (g *DownloaderRPC) Show(ctx context.Context, link string) (*PlaylistInfo, error) {
// 	var resp PlaylistInfo
// 	err := g.client.Call("Plugin.ShowInfo", InputInfo{Link: link, Ctx: ctx}, &resp)
// 	return &resp, err
// }

// func (g *DownloaderRPC) Parse(ctx context.Context, part *Part) error {
// 	args := &DownloadArgs{Ctx: ctx, Part: part}
// 	return g.client.Call("Plugin.Parse", args, new(struct{}))
// }

// func (g *DownloaderRPC) Do(ctx context.Context, part *Part) error {
// 	args := &DownloadArgs{Ctx: ctx, Part: part}
// 	return g.client.Call("Plugin.Do", args, new(struct{}))
// }
// func (g *DownloaderRPC) Cancel(ctx context.Context, part *Part) error {
// 	args := &DownloadArgs{Ctx: ctx, Part: part}
// 	return g.client.Call("Plugin.Cancel", args, new(struct{}))
// }

// ---------------------------------- 服务端信息 ------------------------------

type DownloaderRPCServer struct {
	Impl Downloader
}

// func (s *DownloaderRPCServer) PluginMeta(args struct{}, resp *PluginMeta) error {
// 	*resp = s.Impl.PluginMeta()
// 	return nil
// }

// func (s *DownloaderRPCServer) ShowInfo(args InputInfo, resp *PlaylistInfo) error {
// 	info, err := s.Impl.Show(args.Link, args.Callback)
// 	if err != nil {
// 		return err
// 	}
// 	*resp = *info
// 	return nil
// }

// func (s *DownloaderRPCServer) Do(args *DownloadArgs, resp *struct{}) error {
// 	return s.Impl.Do(args.Part, args.Callback)
// }

// func (s *DownloaderRPCServer) Cancel(args *DownloadArgs, resp *struct{}) error {
// 	return s.Impl.Cancel(args.Part, args.Callback)
// }

// type DownloaderRPCPlugin struct {
// 	Impl Downloader
// }

// func (p *DownloaderRPCPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
// 	return &DownloaderRPCServer{Impl: p.Impl}, nil
// }

// func (DownloaderRPCPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
// 	return &DownloaderRPC{client: c}, nil
// }
