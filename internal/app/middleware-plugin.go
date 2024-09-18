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

// 下载插件
func (a *App) DownloadPlugin(id string) bool {
	m := &plugin.Manifest{
		ID: id,
	}

	if err := a.manager.Download(m); err != nil {
		a.notification.Send(a.ctx, notify.Notice{
			EventName:  "system.notice",
			Content:    "下载插件失败" + err.Error(),
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

	if err := a.manager.UpdatePlugin(p.GetManifest()); err != nil {
		a.notification.Send(a.ctx, notify.Notice{
			EventName:  "system.notice",
			Content:    "下载插件失败" + err.Error(),
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
		a.notification.Send(a.ctx, notify.Notice{
			EventName:  "system.notice",
			Content:    "下载插件失败" + err.Error(),
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

	if err := a.manager.RunPlugin(p.GetManifest()); err != nil {
		a.notification.Send(a.ctx, notify.Notice{
			EventName:  "system.notice",
			Content:    "下载插件失败" + err.Error(),
			NoticeType: "info",
			Provider:   "system",
		})
		return false
	}
	return true
}

func (a *App) UpdatePluginPrams(id string, settings map[string]string) bool {
	p, ok := a.manager.Check(id)
	if !ok {
		return false
	}

	p.GetManifest().Settings = settings

	if err := a.manager.UpdatePluginParams(p.GetManifest()); err != nil {
		a.notification.Send(a.ctx, notify.Notice{
			EventName:  "system.notice",
			Content:    "下载插件失败" + err.Error(),
			NoticeType: "info",
			Provider:   "system",
		})
		return false
	}
	return true
}

// 停止插件
func (a *App) StopPlugin(id string) bool {
	p, ok := a.manager.Check(id)
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
	p, ok := a.manager.Check(id)
	if !ok {
		return false
	}
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

	// 更新插件设置
	ok = a.SavePluginConfig(id, manifest)
	return ok
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

	ok = a.SavePluginConfig(id, manifest)
	return ok
}

// 保存插件配置
func (a *App) SavePluginConfig(id string, m *plugin.Manifest) bool {
	p, ok := a.manager.Check(id)
	if !ok {
		return false
	}

	manifest := p.GetManifest()

	err := manifest.Save()
	return err == nil
}
