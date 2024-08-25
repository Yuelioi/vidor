package app

import "errors"

var (
	pluginConfigSaveError = errors.New("插件保存配置失败")
	pluginNotFound        = errors.New("未找到插件")
	pluginRunFailed       = errors.New("插件运行失败")
	pluginInitFailed      = errors.New("插件初始化失败")
	pluginConnectFailed   = errors.New("插件连接失败")
	pluginUpdateFailed    = errors.New("插件更新失败")
	pluginShutdownFailed  = errors.New("插件关闭失败")
)

var (
	fileOrDirCreationFailed = errors.New("无法创建文件或文件夹")
	fileReadFailed          = errors.New("无法读取文件")
	configConversionFailed  = errors.New("配置转换失败")
)
