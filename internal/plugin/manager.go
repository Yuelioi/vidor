package plugin

import (
	"context"
	"errors"
	"regexp"
)

type PluginManager struct {
	ctx     context.Context
	plugins map[string]Plugin
}

func NewPluginManager(ctx context.Context) *PluginManager {
	return &PluginManager{
		ctx: ctx,
	}
}

// 获取下载器实例
func (pm *PluginManager) Select(url string) (*DownloadPlugin, error) {
	for _, p := range pm.plugins {

		base := p.GetManifest()
		if base.Type == "downloader" {
			downloadPlugin, ok := p.(*DownloadPlugin)
			if !ok {
				return nil, nil
			}

			for _, match := range downloadPlugin.Manifest.Matches {
				reg, err := regexp.Compile(match)
				if err != nil {
					return nil, errors.New("插件正则表达式编译失败: " + err.Error())
				}
				if reg.MatchString(url) {
					return downloadPlugin, nil
				}
			}
		}

	}
	return nil, errors.New("未找到")
}

// 检查插件是否存在

func (pm *PluginManager) Check(id string) (Plugin, bool) {
	p, ok := pm.plugins[id]
	if !ok {
		return nil, false
	}
	return p, true
}

// 获取配置信息
func (pm *PluginManager) Manifests() map[string]Manifest {
	ms := make(map[string]Manifest, 0)

	for _, plugin := range pm.plugins {
		mf := plugin.GetManifest()
		ms[mf.ID] = *mf
	}
	return ms
}

func (pm *PluginManager) Register(m *Manifest) error {
	// 注册
	var p Plugin

	switch m.Type {
	case "downloader":
		pd := NewDownloader(m)
		p = pd

	default:
		return errors.New("未知的插件类型")
	}

	// 注册配置以及插件
	pm.plugins[m.ID] = p
	return nil
}

// ------------------------------------ Handlers ------------------------------------

func (pm *PluginManager) createHandlerChain(handlers ...PluginHandler) PluginHandler {
	if len(handlers) == 0 {
		return nil
	}

	for i := 0; i < len(handlers)-1; i++ {
		handlers[i].SetNext(handlers[i+1])
	}

	return handlers[0]
}

// 下载插件
//
// 1.下载
// 2.解压
// 3.注册到主机
// 4.运行
func (pm *PluginManager) Download(m *Manifest) error {
	handlerChain := pm.createHandlerChain(
		&DownloadHandler{},
		&ExtractHandler{},
		&RunHandler{},
	)
	return handlerChain.Handle(pm.ctx, m)
}

// 更新插件
//
// 1.禁用当前插件并删除
// 2.下载并解压
// 3.注册到主机
// 4.运行
func (pm *PluginManager) UpdatePlugin(m *Manifest) error {
	handlerChain := pm.createHandlerChain(
	// &RegisterHandler{},
	// &DownloadHandler{},
	// &ExtractHandler{},
	// &RegisterHandler{},
	)

	return handlerChain.Handle(pm.ctx, m)
}

// 删除插件
//
// 1.禁用当前插件并删除
// 2.注销插件
func (pm *PluginManager) RemovePlugin(m *Manifest) error {
	handlerChain := pm.createHandlerChain(
	// &RegisterHandler{},
	// &DownloadHandler{},
	// &ExtractHandler{},
	// &RegisterHandler{},
	)
	return handlerChain.Handle(pm.ctx, m)
}

// 运行插件
//
// 1.禁用当前插件并删除
// 2.注销插件
func (pm *PluginManager) RunPlugin(m *Manifest) error {
	handlerChain := pm.createHandlerChain(
	// &RegisterHandler{},
	// &DownloadHandler{},
	// &ExtractHandler{},
	// &RegisterHandler{},
	)
	return handlerChain.Handle(pm.ctx, m)
}
