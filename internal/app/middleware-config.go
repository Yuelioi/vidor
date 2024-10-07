package app

import (
	"context"
	"path/filepath"

	"github.com/Yuelioi/vidor/internal/config"
	"google.golang.org/grpc/metadata"
)

// 获取主机配置
func (a *App) GetConfig() *config.Config {
	return a.config
}

// 保存配置到本地
func (a *App) SaveConfig(config *config.Config) bool {

	if config.BaseDir == "" {
		config.BaseDir = a.config.BaseDir
	}

	// 更新前端传来的配置信息
	a.config = config

	// 保存配置文件
	err := a.config.Save()
	if err != nil {
		a.logger.Warnf("保存设置失败: %s", err)
		return false
	} else {
		a.logger.Info("保存设置成功")
	}

	a.manager.UpdateSystemParams(a.injectMetadata())

	return err == nil
}

// 注入系统metadata(默认的+config)
func (a *App) injectMetadata() context.Context {
	configs := map[string]string{
		"system.ffmpeg":     filepath.Join(a.appDirs.Libs, "ffmpeg", "ffmpeg.exe"),
		"system.tmpdirs":    a.appDirs.Temps,
		"system.logdirs":    a.appDirs.Logs,
		"system.plugindirs": a.appDirs.Plugins,
		"system.configdirs": a.appDirs.Configs,
	}

	ctx := context.Background()

	for key, value := range configs {
		ctx = metadata.AppendToOutgoingContext(ctx, key, value)
	}

	return a.config.InjectMetadata(ctx)
}
