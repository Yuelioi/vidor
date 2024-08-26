package globals

import "errors"

var (
	ErrPluginConfigSave = errors.New("插件保存配置失败")
	ErrPluginNotFound   = errors.New("未找到插件")
	ErrPluginRun        = errors.New("插件运行失败")
	ErrPluginInit       = errors.New("插件初始化失败")
	ErrPluginConnect    = errors.New("插件连接失败")
	ErrPluginUpdate     = errors.New("插件更新失败")
	ErrPluginShutdown   = errors.New("插件关闭失败")
)

var (
	ErrFileOrDirCreation = errors.New("无法创建文件或文件夹")
	ErrFileRead          = errors.New("无法读取文件")
	ErrConfigConversion  = errors.New("配置转换失败")
)
