package app

import (
	"context"

	"github.com/Yuelioi/vidor/internal/plugin"
)

type PluginManager struct {
	a   *App
	ctx context.Context

	downloadHandler *DownloadHandler
	extractHandler  *ExtractHandler
	registerHandler *RegisterHandler
	remover         *RegisterHandler
}

func NewPluginManager(a *App) *PluginManager {
	return &PluginManager{
		a:               a,
		ctx:             context.Background(),
		downloadHandler: &DownloadHandler{},
		extractHandler:  &ExtractHandler{},
		registerHandler: &RegisterHandler{},
		remover:         &RegisterHandler{},
	}
}

// 下载插件
//
// 1.下载
// 2.解压
// 3.注册到主机
// 4.运行
func (pm *PluginManager) Download(m *plugin.Manifest) {
	pm.downloadHandler.SetNext(pm.extractHandler).SetNext(pm.registerHandler)
	pm.downloadHandler.Handle(pm.ctx, pm.a, m)
}

// 更新插件
//
// 1.禁用当前插件并删除
// 2.下载并解压
// 3.注册到主机
// 4.运行
func (pm *PluginManager) UpdatePlugin(m *plugin.Manifest) {
	pm.remover.SetNext(pm.downloadHandler).SetNext(pm.registerHandler)
	pm.remover.Handle(pm.ctx, pm.a, m)
}

// 删除插件
//
// 1.禁用当前插件并删除
// 2.注销插件
func (pm *PluginManager) RemovePlugin(m *plugin.Manifest) {
	pm.remover.SetNext(pm.downloadHandler).SetNext(pm.registerHandler)

}
