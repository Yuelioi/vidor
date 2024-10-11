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
		ctx:     ctx,
		plugins: make(map[string]Plugin, 0),
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
				if reg.MatchString(url) && downloadPlugin.Manifest.State == Working {
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

func (pm *PluginManager) NetManifests() ([]*Manifest, error) {
	return fetchPlugins()
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
// 5.注入上下文
// 6.保存插件配置
//
// 需要传入正确的系统参数上下文
func (pm *PluginManager) Download(m *Manifest, ctx context.Context) error {
	handlerChain := pm.createHandlerChain(
		&DownloadHandler{},
		&ExtractHandler{},
		&RegisterPMHandler{pm: pm},
		&RunnerPMHandler{pm: pm},
		&InitPMHandler{pm: pm},
		&UpdateParamsPMHandler{pm: pm},
		&SaveHandler{},
	)
	return handlerChain.Handle(pm.ctx, m)
}

// 更新插件本体
//
// 1.禁用,注销插件并删除
// 2.下载并解压
// 3.运行
// 4.注入上下文, 保存配置
//
// 需要传入正确的系统参数上下文
func (pm *PluginManager) UpdatePlugin(m *Manifest, ctx context.Context) error {
	handlerChain := pm.createHandlerChain(
		&StopperPMHandler{pm: pm},
		&DeRegisterPMHandler{pm: pm},
		&RemoveHandler{},
		&DownloadHandler{},
		&RegisterPMHandler{pm: pm},
		&RunnerPMHandler{pm: pm},
		&InitPMHandler{pm: pm},
		&UpdateParamsPMHandler{pm: pm},
		&SaveHandler{},
	)

	return handlerChain.Handle(pm.ctx, m)
}

// 删除插件
//
// 1.禁用当前插件并删除
// 2.注销插件
func (pm *PluginManager) RemovePlugin(m *Manifest) error {
	handlerChain := pm.createHandlerChain(
		&StopperPMHandler{pm: pm},
		&DeRegisterPMHandler{pm: pm},
		&RemoveHandler{},
	)
	return handlerChain.Handle(pm.ctx, m)
}

// 运行插件
//
// 1.运行
// 2.注入上下文
//
// 需要传入正确的系统参数上下文
func (pm *PluginManager) RunPlugin(m *Manifest, ctx context.Context) error {
	handlerChain := pm.createHandlerChain(
		&RunnerPMHandler{pm: pm},
		&InitPMHandler{pm: pm},
		&UpdateParamsPMHandler{pm: pm},
	)
	return handlerChain.Handle(ctx, m)
}

// 停止插件
func (pm *PluginManager) StopPlugin(m *Manifest) error {
	handlerChain := pm.createHandlerChain(
		&StopperPMHandler{pm: pm},
	)
	return handlerChain.Handle(pm.ctx, m)
}

// 注册插件
func (pm *PluginManager) Register(m *Manifest) error {
	handlerChain := pm.createHandlerChain(
		&RegisterPMHandler{pm: pm},
	)

	return handlerChain.Handle(pm.ctx, m)
}

// 更新插件参数
func (pm *PluginManager) UpdatePluginParams(m *Manifest) error {
	handlerChain := pm.createHandlerChain(
		&UpdatePluginParamsPMHandler{pm: pm},
	)

	return handlerChain.Handle(pm.ctx, m)
}

// 更新系统参数
func (pm *PluginManager) UpdateSystemParams(ctx context.Context) error {

	handlerChain := pm.createHandlerChain(
		&UpdateSystemParamsPMHandler{pm: pm},
	)

	for _, p := range pm.plugins {
		if p.GetManifest().State == Working {
			handlerChain.Handle(ctx, p.GetManifest())
		}
	}
	return nil
}
