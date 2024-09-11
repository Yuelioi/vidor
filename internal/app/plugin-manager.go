package app

import "github.com/Yuelioi/vidor/internal/plugin"

type PluginManager struct {
	downloadHandler *DownloadHandler
	extractHandler  *ExtractHandler
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		downloadHandler: &DownloadHandler{},
		extractHandler:  &ExtractHandler{},
	}
}

func (pm *PluginManager) Download(m *plugin.Manifest) {
	pm.downloadHandler.SetNext(pm.extractHandler)
	pm.downloadHandler.Handle(m)
}
