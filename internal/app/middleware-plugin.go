package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Yuelioi/vidor/internal/notify"
	"github.com/Yuelioi/vidor/internal/plugin"
	"github.com/Yuelioi/vidor/internal/tools"
	"github.com/Yuelioi/vidor/pkg/downloader"
	"github.com/go-resty/resty/v2"

	"golift.io/xtractr"
)

// 返回主机注册的插件
func (a *App) GetPlugins() map[string]plugin.Manifest {
	ms := make(map[string]plugin.Manifest, 0)

	for _, plugin := range a.plugins {
		mf := plugin.GetManifest()
		ms[mf.ID] = *mf
	}
	return ms
}

// 获取网络插件列表
func fetchPlugins(pluginsDir string) ([]*plugin.Manifest, error) {
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

	plugins := make([]*plugin.Manifest, 0, len(rawPlugins))
	for _, rawPlugin := range rawPlugins {
		manifest := plugin.NewManifest(pluginsDir)

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

func (a *App) checkPlugin(id string) (plugin.Plugin, bool) {
	p, ok := a.plugins[id]
	if !ok {
		a.notification.Send(a.ctx, notify.Notice{
			EventName:  "system.notice",
			Content:    "未找到当前插件, 请联系作者",
			NoticeType: "info",
			Provider:   "system",
		})
		return nil, false
	}
	return p, true
}

func (a *App) downloadPlugin(id string) (*plugin.Manifest, error) {
	// 获取插件列表
	plugins, err := fetchPlugins(a.appDirs.Plugins)
	if err != nil {
		a.notification.Send(a.ctx, notify.Notice{
			EventName:  "system.notice",
			Content:    "获取插件信息失败:" + err.Error(),
			NoticeType: "info",
			Provider:   "system",
		})
		return nil, err
	}

	var targetManifest *plugin.Manifest
	for _, manifest := range plugins {
		if id == manifest.ID {
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
		return nil, errors.New("没有找到插件下载链接")
	}

	downUrl := targetManifest.DownloadURLs[0]
	name := tools.ExtractFileNameFromUrl(downUrl)
	name = tools.SanitizeFileName(name)

	zipPath := filepath.Join(a.appDirs.Temps, name)
	targetDir := filepath.Join(a.appDirs.Plugins, targetManifest.ID)
	err = os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	errCh := make(chan error)
	go func() {
		err := dl.Download()
		errCh <- err
	}()

	err = <-errCh
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Download succeeded")
	}

	// 解压
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
		return nil, err
	}
	return targetManifest, nil
}

// 下载插件
//
// 1.下载
// 2.解压
// 3.注册到主机/运行
func (a *App) DownloadPlugin(id string) *plugin.Manifest {
	targetManifest, err := a.downloadPlugin(id)
	if err != nil {
		return nil
	}

	// 注册
	a.registerPlugin(targetManifest)

	return targetManifest
}

// 禁用并移除插件
func (a *App) disableAndDeletePlugin(p plugin.Plugin) error {
	// 禁用插件
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

	return nil
}

// 更新插件
//
// 1.禁用当前插件并删除
// 2.下载并解压
// 3.注册到主机/运行
func (a *App) UpdatePlugin(id string) bool {
	p, ok := a.checkPlugin(id)
	if !ok {
		return false
	}

	// 禁用当前插件
	if err := a.disableAndDeletePlugin(p); err != nil {
		a.notification.Send(a.ctx, notify.Notice{
			EventName:  "system.notice",
			Content:    "禁用当前插件失败" + err.Error(),
			NoticeType: "info",
			Provider:   "system",
		})
		return false
	}

	// 重新下载
	targetManifest, err := a.downloadPlugin(id)
	if err != nil {
		return false
	}

	// 注册
	a.registerPlugin(targetManifest)
	return true

}

// 删除插件
//
// 1.禁用当前插件并删除
func (a *App) RemovePlugin(id string) bool {

	p, ok := a.checkPlugin(id)
	if !ok {
		return false
	}

	err := a.disableAndDeletePlugin(p)
	return err == nil
}

// 运行插件, 并建立连接
func (a *App) RunPlugin(id string) bool {
	plugin, ok := a.checkPlugin(id)
	if !ok {
		return false
	}

	// 运行
	err := plugin.Run(context.Background())
	if err != nil {
		a.logger.Infof("插件开启失败:%s", err)
		return false
	}

	err = plugin.Init(context.Background())
	if err != nil {
		a.logger.Infof("插件开启失败:%s", err)
		return false
	}

	return true
}

// 更新插件参数
func (a *App) UpdatePluginPrams(id string, settings map[string]string) bool {
	p, ok := a.checkPlugin(id)
	if !ok {
		return false
	}

	p.GetManifest().Settings = settings
	err := p.Update(context.Background())
	if err != nil {
		a.logger.Infof("插件开启失败:%s", err)
		return false
	}

	return true
}

func (a *App) StopPlugin(id string) bool {
	p, ok := a.checkPlugin(id)
	if !ok {
		return false
	}
	// 停止
	err := p.Shutdown(context.Background())
	if err != nil {
		a.notification.Send(a.ctx, notify.Notice{
			EventName:  "system.notice",
			Content:    "插件保存失败",
			NoticeType: "info",
			Provider:   "system",
		})
		return false
	}
	p.GetManifest().State = plugin.NotWork
	return true
}

// 启用插件, 但是不会运行
func (a *App) EnablePlugin(id string) bool {
	p, ok := a.plugins[id]
	if !ok {
		a.notification.Send(a.ctx, notify.Notice{
			EventName:  "system.notice",
			Content:    "未找到当前插件, 请联系作者",
			NoticeType: "info",
			Provider:   "system",
		})
		return false
	}
	manifest := p.GetManifest()
	manifest.Enable = true

	return true
}

// 关闭插件,并禁用插件
func (a *App) DisablePlugin(id string) bool {
	p, ok := a.checkPlugin(id)
	if !ok {
		return false
	}

	manifest := p.GetManifest()

	// 关闭插件
	if manifest.State == plugin.Working {
		err := p.Shutdown(context.Background())
		if err != nil {
			return false
		}
	}

	// 禁用并保存配置
	manifest.Enable = false

	ok = a.SavePluginConfig(id, manifest)
	return ok
}

// 保存插件配置
func (a *App) SavePluginConfig(id string, m *plugin.Manifest) bool {
	p, ok := a.checkPlugin(id)
	if !ok {
		return false
	}

	manifest := p.GetManifest()

	err := manifest.Save()
	return err == nil
}
