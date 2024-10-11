package app

import (
	"context"

	"github.com/Yuelioi/vidor/internal/notify"
	"github.com/Yuelioi/vidor/internal/plugin"
)

// 返回主机注册的插件
func (a *App) GetPlugins() map[string]plugin.Manifest {
	return a.manager.Manifests()
}

// 返回主机注册的插件
func (a *App) GetMarketPlugins() []*plugin.Manifest {

	ms, err := a.manager.NetManifests()
	if err != nil {
		return nil
	}
	return ms
}

// 下载插件
func (a *App) DownloadPlugin(id string) bool {
	m := &plugin.Manifest{
		ID: id,
	}

	ctx := a.injectMetadata()

	if err := a.manager.Download(m, ctx); err != nil {
		a.notification.Send(notify.Notice{
			EventName:  "system.notice",
			Content:    "下载插件失败: " + err.Error(),
			NoticeType: "info",
			Provider:   "system",
		})
		return false
	}
	return true
}

// 更新插件本体
//
// 1.禁用当前插件并删除
// 2.下载并解压
// 3.注册到主机/运行
func (a *App) UpdatePlugin(id string) bool {
	p, ok := a.manager.Check(id)
	if !ok {
		return false
	}

	ctx := a.injectMetadata()

	if err := a.manager.UpdatePlugin(p.GetManifest(), ctx); err != nil {
		a.notification.Send(notify.Notice{
			EventName:  "system.notice",
			Content:    "更新插件失败: " + err.Error(),
			NoticeType: "info",
			Provider:   "system",
		})
		return false
	}
	return true
}

// 删除插件
//
// 1.禁用当前插件并删除
func (a *App) RemovePlugin(id string) bool {
	p, ok := a.manager.Check(id)
	if !ok {
		return false
	}

	if err := a.manager.RemovePlugin(p.GetManifest()); err != nil {
		a.notification.Send(notify.Notice{
			EventName:  "system.notice",
			Content:    "移除插件失败: " + err.Error(),
			NoticeType: "info",
			Provider:   "system",
		})
		return false
	}
	return true
}

// 运行插件, 并建立连接
func (a *App) RunPlugin(id string) bool {
	p, ok := a.manager.Check(id)
	if !ok {
		return false
	}

	ctx := a.injectMetadata()
	m := p.GetManifest()
	if err := a.manager.RunPlugin(p.GetManifest(), ctx); err != nil {
		a.notification.Send(notify.Notice{
			EventName:  "system.notice",
			Content:    "运行插件失败: " + err.Error(),
			NoticeType: "info",
			Provider:   "system",
		})
		return false
	}

	return a.savePluginConfig(m)
}

func (a *App) UpdatePluginPrams(id string, settings map[string]string) bool {
	p, ok := a.manager.Check(id)
	if !ok {
		return false
	}
	manifest := p.GetManifest()
	manifest.Settings = settings

	if manifest.State == plugin.Working {
		if err := a.manager.UpdatePluginParams(p.GetManifest()); err != nil {
			a.notification.Send(notify.Notice{
				EventName:  "system.notice",
				Content:    "更新插件参数失败: " + err.Error(),
				NoticeType: "info",
				Provider:   "system",
			})
			return false
		}
	}
	return a.savePluginConfig(manifest)
}

// 停止插件
func (a *App) StopPlugin(id string) bool {
	p, ok := a.manager.Check(id)
	if !ok {
		return false
	}
	manifest := p.GetManifest()
	// 停止
	if manifest.State == plugin.Working {
		if err := a.manager.StopPlugin(manifest); err != nil {
			a.notification.Send(notify.Notice{
				EventName:  "system.notice",
				Content:    "插件关闭失败: " + err.Error(),
				NoticeType: "info",
				Provider:   "system",
			})

			manifest.State = plugin.NotWork

			pn := notify.NewPluginNotification(a.ctx)
			pn.Send(*manifest)

			return false
		}
	}

	return true
}

// 启用插件, 但是不会运行
func (a *App) EnablePlugin(id string) bool {
	p, ok := a.manager.Check(id)
	if !ok {
		return false
	}
	if !ok {
		a.notification.Send(notify.Notice{
			EventName:  "system.notice",
			Content:    "未找到当前插件, 请联系作者",
			NoticeType: "info",
			Provider:   "system",
		})
		return false
	}
	manifest := p.GetManifest()
	manifest.Enable = true

	// 更新插件设置
	return a.savePluginConfig(manifest)
}

// 关闭插件,并禁用插件
func (a *App) DisablePlugin(id string) bool {
	p, ok := a.manager.Check(id)
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

	return a.savePluginConfig(manifest)
}

// 保存插件配置
func (a *App) savePluginConfig(m *plugin.Manifest) bool {
	p, ok := a.manager.Check(m.ID)
	if !ok {
		return false
	}

	manifest := p.GetManifest()
	manifest.Settings = m.Settings

	return manifest.Save() == nil
}

// 保存插件配置并更新参数
func (a *App) SavePluginConfig(m *plugin.Manifest) bool {
	p, ok := a.manager.Check(m.ID)
	if !ok {
		return false
	}

	manifest := p.GetManifest()
	manifest.Settings = m.Settings

	err := manifest.Save()
	if err != nil {
		return false
	}

	return a.manager.UpdatePluginParams(manifest) == nil
}
