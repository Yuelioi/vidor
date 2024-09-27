package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Yuelioi/vidor/internal/tools"
	"github.com/Yuelioi/vidor/pkg/downloader"
	"github.com/go-resty/resty/v2"
	"golift.io/xtractr"
)

type key int

// 上下文键
const (
	KeyAppPath key = iota
	keyZipPath
)

type PluginHandler interface {
	Handle(ctx context.Context, m *Manifest) error
	SetNext(next PluginHandler) PluginHandler
}

type BaseHandler struct {
	next PluginHandler
}

func (bh *BaseHandler) Handle(ctx context.Context, m *Manifest) error {
	if bh.next != nil {
		return bh.next.Handle(ctx, m)
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

// 获取网络插件列表
func fetchPlugins() ([]*Manifest, error) {
	pluginsUrl := "https://cdn.yuelili.com/market/vidor/plugins.json"

	resp, err := resty.New().R().Get(pluginsUrl)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("链接失败")
	}

	var rawPlugins []map[string]interface{}
	err = json.Unmarshal(resp.Body(), &rawPlugins)
	if err != nil {
		return nil, err
	}

	plugins := make([]*Manifest, 0, len(rawPlugins))
	for _, rawPlugin := range rawPlugins {
		manifest := NewManifest("")

		manifestData, err := json.Marshal(rawPlugin)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(manifestData, manifest)
		if err != nil {
			return nil, err
		}

		plugins = append(plugins, manifest)
	}

	return plugins, nil
}

func (d *DownloadHandler) Handle(ctx context.Context, m *Manifest) error {
	// 获取插件列表
	plugins, err := fetchPlugins()
	if err != nil {
		return err
	}

	var targetManifest *Manifest
	for _, manifest := range plugins {
		if m.ID == manifest.ID {
			targetManifest = manifest
		}
	}

	if len(targetManifest.DownloadURLs) == 0 {
		return errors.New("没有找到插件下载链接")
	}

	downUrl := targetManifest.DownloadURLs[0]
	name := tools.ExtractFileNameFromUrl(downUrl)
	name = tools.SanitizeFileName(name)

	appDir, ok := ctx.Value(KeyAppPath).(string)
	if !ok {
		return errors.New("未找到上下文值")
	}

	zipPath := filepath.Join(appDir, "temps", name)
	targetDir := filepath.Join(appDir, "plugin", targetManifest.ID)
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
	ctx = context.WithValue(ctx, keyZipPath, zipPath)
	return d.BaseHandler.Handle(ctx, m)
}

// 解压
type ExtractHandler struct {
	BaseHandler
}

func (e *ExtractHandler) Handle(ctx context.Context, m *Manifest) error {
	zipPath, ok := ctx.Value(keyZipPath).(string)
	if !ok {
		fmt.Println("zipPath from context:", zipPath)
	} else {
		return errors.New("未找到解压路径")
	}

	appDir, ok := ctx.Value(KeyAppPath).(string)
	if !ok {
		return errors.New("未找到上下文值")
	}

	targetDir := filepath.Join(appDir, "plugins", m.ID)

	x := &xtractr.XFile{
		FilePath:  zipPath,
		OutputDir: targetDir,
	}

	_, files, _, err := xtractr.ExtractFile(x)
	if err != nil || files == nil {
		return errors.New("解压插件失败" + err.Error())

	}

	return e.BaseHandler.Handle(ctx, m)
}

//

// 删除
type RemoveHandler struct {
	BaseHandler
}

func (r *RemoveHandler) Handle(ctx context.Context, m *Manifest) error {
	// 删除插件文件夹
	if err := os.RemoveAll(m.BaseDir); err != nil {
		return err
	}
	return r.BaseHandler.Handle(ctx, m)
}

// 保存配置
type SaveHandler struct {
	BaseHandler
}

func (r *SaveHandler) Handle(ctx context.Context, m *Manifest) error {
	if err := m.Save(); err != nil {
		return err
	}
	return r.BaseHandler.Handle(ctx, m)
}

// -------------------------------------------------------------------------------

// 注册
type RegisterPMHandler struct {
	BaseHandler
	pm *PluginManager
}

func (r *RegisterPMHandler) Handle(ctx context.Context, m *Manifest) error {
	fmt.Printf("Register Handle\n")
	var p Plugin

	switch m.Type {
	case "downloader":
		pd := NewDownloader(m)
		p = pd

	default:
		return errors.New("未知的插件类型")
	}

	r.pm.plugins[m.ID] = p
	return r.BaseHandler.Handle(ctx, m)
}

// 注销
type DeRegisterPMHandler struct {
	BaseHandler
	pm *PluginManager
}

func (r *DeRegisterPMHandler) Handle(ctx context.Context, m *Manifest) error {
	fmt.Printf("DeRegister Handle\n")
	for key := range r.pm.plugins {
		if key == m.ID {
			delete(r.pm.plugins, key)
			return r.BaseHandler.Handle(ctx, m)
		}
	}
	return errors.New("deRegister 未找到插件")
}

type RunnerPMHandler struct {
	BaseHandler
	pm *PluginManager
}

func (r *RunnerPMHandler) Handle(ctx context.Context, m *Manifest) error {
	fmt.Printf("Runner Handle\n")

	for key, p := range r.pm.plugins {
		if key == m.ID {
			err := p.Run(ctx)
			if err == nil {
				m.State = Working
				r.BaseHandler.Handle(ctx, m)
			}
			return err
		}
	}
	return errors.New("runner 未找到插件")
}

type InitPMHandler struct {
	BaseHandler
	pm *PluginManager
}

func (r *InitPMHandler) Handle(ctx context.Context, m *Manifest) error {
	fmt.Printf("Init Handler\n")
	for key, p := range r.pm.plugins {
		if key == m.ID {
			err := p.Init(ctx)
			if err == nil {
				r.BaseHandler.Handle(ctx, m)
			}
			return err
		}
	}
	return errors.New("init 未找到插件")
}

type StopperPMHandler struct {
	BaseHandler
	pm *PluginManager
}

func (r *StopperPMHandler) Handle(ctx context.Context, m *Manifest) error {
	fmt.Printf("Stopper Handle\n")

	for key, p := range r.pm.plugins {
		if key == m.ID {
			err := p.Shutdown(ctx)
			if err == nil {
				m.State = NotWork
				return r.BaseHandler.Handle(ctx, m)
			}
			return err
		}
	}
	return errors.New("stopper 未找到插件")
}

// 注入插件参数
type UpdatePluginParamsPMHandler struct {
	BaseHandler
	pm *PluginManager
}

func (r *UpdatePluginParamsPMHandler) Handle(ctx context.Context, m *Manifest) error {
	fmt.Printf("UpdatePluginParams Handler\n")
	for key, p := range r.pm.plugins {
		if key == m.ID {
			ctx = InjectMetadata(ctx, m.Settings)
			if err := p.Update(ctx); err != nil {
				return err
			}
			return r.BaseHandler.Handle(ctx, m)
		}
	}

	return errors.New("updatePluginParams 未找到插件")
}

// 同时注入插件参数与系统参数
type UpdateParamsPMHandler struct {
	BaseHandler
	pm *PluginManager
}

func (r *UpdateParamsPMHandler) Handle(ctx context.Context, m *Manifest) error {
	fmt.Printf("UpdateParams Handler\n")

	for key, p := range r.pm.plugins {
		if key == m.ID {
			ctx = InjectMetadata(ctx, m.Settings)

			if err := p.Update(ctx); err != nil {
				return err
			}
			return r.BaseHandler.Handle(ctx, m)

		}
	}
	return errors.New("updateSystemParams 未找到插件")
}

// 注入系统参数(请提前传正确的ctx)
type UpdateSystemParamsPMHandler struct {
	BaseHandler
	pm *PluginManager
}

func (r *UpdateSystemParamsPMHandler) Handle(ctx context.Context, m *Manifest) error {
	fmt.Printf("UpdateSystemParams Handler\n")

	for key, p := range r.pm.plugins {
		if key == m.ID {
			if err := p.Update(ctx); err != nil {
				return err
			}
			return r.BaseHandler.Handle(ctx, m)

		}
	}
	return errors.New("updateSystemParams 未找到插件")
}
