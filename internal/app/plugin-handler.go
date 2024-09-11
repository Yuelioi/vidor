package app

import (
	"github.com/Yuelioi/vidor/internal/plugin"
)

type PluginHandler interface {
	Handle(m *plugin.Manifest) error
	SetNext(next PluginHandler) PluginHandler
}

type BaseHandler struct {
	next PluginHandler
}

func (bh *BaseHandler) Handle(m *plugin.Manifest) error {
	if bh.next != nil {
		return bh.Handle(m)
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

func (d *DownloadHandler) Handle(m *plugin.Manifest) error {
	return d.BaseHandler.Handle(m)
}

// 解压
type ExtractHandler struct {
	BaseHandler
}

func (e *ExtractHandler) Handle(m *plugin.Manifest) error {
	return e.BaseHandler.Handle(m)
}

// 删除
type RemoveHandler struct {
	BaseHandler
}

type Register struct {
	BaseHandler
}
