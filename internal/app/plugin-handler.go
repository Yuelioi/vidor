package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Yuelioi/vidor/internal/notify"
	"github.com/Yuelioi/vidor/internal/plugin"
	"github.com/Yuelioi/vidor/internal/tools"
	"github.com/Yuelioi/vidor/pkg/downloader"
	"golift.io/xtractr"
)

type PluginHandler interface {
	Handle(ctx context.Context, a *App, m *plugin.Manifest) error
	SetNext(next PluginHandler) PluginHandler
}

type BaseHandler struct {
	next PluginHandler
}

func (bh *BaseHandler) Handle(ctx context.Context, a *App, m *plugin.Manifest) error {
	if bh.next != nil {
		return bh.Handle(ctx, a, m)
	}
	return nil
}

func (bh *BaseHandler) SetNext(next PluginHandler) PluginHandler {
	bh.next = next
	return next
}

// 下载
type DownloadHandler struct {
	BaseHandler
}

func (d *DownloadHandler) Handle(ctx context.Context, a *App, m *plugin.Manifest) error {
	// 获取插件列表
	plugins, err := fetchPlugins(a.appDirs.Plugins)
	if err != nil {
		a.notification.Send(a.ctx, notify.Notice{
			EventName:  "system.notice",
			Content:    "获取插件信息失败:" + err.Error(),
			NoticeType: "info",
			Provider:   "system",
		})
		return err
	}

	var targetManifest *plugin.Manifest
	for _, manifest := range plugins {
		if m.ID == manifest.ID {
			targetManifest = manifest
		}
	}

	if len(targetManifest.DownloadURLs) == 0 {
		a.notification.Send(a.ctx, notify.Notice{
			EventName:  "system.notice",
			Content:    "获取插件链接失败",
			NoticeType: "info",
			Provider:   "system",
		})
		return errors.New("没有找到插件下载链接")
	}

	downUrl := targetManifest.DownloadURLs[0]
	name := tools.ExtractFileNameFromUrl(downUrl)
	name = tools.SanitizeFileName(name)

	zipPath := filepath.Join(a.appDirs.Temps, name)
	targetDir := filepath.Join(a.appDirs.Plugins, targetManifest.ID)
	err = os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		return err
	}

	dl, err := downloader.New(
		context.Background(),
		downUrl,
		zipPath,
		true)
	if err != nil {
		a.notification.Send(a.ctx, notify.Notice{
			EventName:  "system.notice",
			Content:    "下载插件失败" + err.Error(),
			NoticeType: "info",
			Provider:   "system",
		})
		return err
	}
	errCh := make(chan error)
	go func() {
		err := dl.Download()
		errCh <- err
	}()

	err = <-errCh
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, "zipPath", zipPath)
	return d.BaseHandler.Handle(ctx, a, m)
}

// 解压
type ExtractHandler struct {
	BaseHandler
}

func (e *ExtractHandler) Handle(ctx context.Context, a *App, m *plugin.Manifest) error {
	zipPathFromCtx := ctx.Value("zipPath")
	zipPath, ok := zipPathFromCtx.(string)
	if !ok {
		fmt.Println("zipPath from context:", zipPath)
	} else {
		return errors.New("未找到解压路径")
	}

	targetDir := filepath.Join(a.appDirs.Plugins, m.ID)

	x := &xtractr.XFile{
		FilePath:  zipPath,
		OutputDir: targetDir,
	}

	_, files, _, err := xtractr.ExtractFile(x)
	if err != nil || files == nil {
		a.notification.Send(a.ctx, notify.Notice{
			EventName:  "system.notice",
			Content:    "解压插件失败" + err.Error(),
			NoticeType: "info",
			Provider:   "system",
		})
	}

	return e.BaseHandler.Handle(ctx, a, m)
}

//

// 禁用并删除
type RemoveHandler struct {
	BaseHandler
}

func (r *RemoveHandler) Handle(ctx context.Context, a *App, m *plugin.Manifest) error {

	p, ok := a.plugins[m.ID]
	if !ok {
		a.notification.Send(a.ctx, notify.Notice{
			EventName:  "system.notice",
			Content:    "未找到当前插件, 请联系作者",
			NoticeType: "info",
			Provider:   "system",
		})
		return errors.New("未找到当前插件")
	}
	if err := p.Shutdown(context.Background()); err != nil {
		a.notification.Send(a.ctx, notify.Notice{
			EventName:  "system.notice",
			Content:    "禁用插件失败插件失败" + err.Error(),
			NoticeType: "info",
			Provider:   "system",
		})
		return err
	}

	// 删除插件文件夹
	if err := os.RemoveAll(p.GetManifest().BaseDir); err != nil {
		return err
	}
	return r.BaseHandler.next.Handle(ctx, a, m)
}

// 注册
type RegisterHandler struct {
	BaseHandler
}

func (r *RegisterHandler) Handle(ctx context.Context, a *App, m *plugin.Manifest) error {
	var p plugin.Plugin
	switch m.Type {
	case "downloader":
		pd := plugin.NewDownloader(m)
		p = pd

	default:
		return errors.New("未知的插件类型")
	}

	// 注册配置以及插件
	a.plugins[m.ID] = p
	return r.BaseHandler.next.Handle(ctx, a, m)
}

// 注销
type DeRegisterHandler struct {
	BaseHandler
}

func (r *DeRegisterHandler) Handle(ctx context.Context, a *App, m *plugin.Manifest) error {
	for id := range a.plugins {
		if id == m.ID {
			delete(a.plugins, id)
		}
	}
	return r.BaseHandler.next.Handle(ctx, a, m)
}

// 保存配置
type SaveHandler struct {
	BaseHandler
}

func (r *SaveHandler) Handle(ctx context.Context, a *App, m *plugin.Manifest) error {
	if err := m.Save(); err != nil {
		return err
	}
	return r.BaseHandler.next.Handle(ctx, a, m)
}
