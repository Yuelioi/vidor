package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "embed"

	"github.com/Yuelioi/vidor/internal/config"
	"github.com/Yuelioi/vidor/internal/globals"
	"github.com/Yuelioi/vidor/internal/logger"
	"github.com/Yuelioi/vidor/internal/notify"
	"github.com/Yuelioi/vidor/internal/plugin"
	"github.com/Yuelioi/vidor/internal/task"
	"github.com/Yuelioi/vidor/internal/tools"
	"github.com/energye/systray"
	"github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed icon.ico
var iconData []byte

var pluginsDir string
var assetsDir string

// 应用实例
type App struct {
	ctx          context.Context
	location     string // 软件路径
	appInfo      AppInfo
	config       *config.Config           // 软件配置信息
	taskQueue    *task.TaskQueue          // 任务队列 用于分发任务
	plugins      map[string]plugin.Plugin // 插件
	cache        *Cache                   // 缓存
	notification *notify.Notification     // 消息分发
	logger       *logrus.Logger
}

func NewApp() *App {
	a := &App{
		plugins: make(map[string]plugin.Plugin),
	}

	appDir, err := tools.ExeDir()
	if err != nil {
		log.Fatal()
	}

	appDir, _ = filepath.Abs(appDir)
	a.location = appDir

	// 初始化文件夹
	loggerDir := filepath.Join(appDir, "logs")
	configDir := filepath.Join(appDir, "configs")
	pluginsDir = filepath.Join(appDir, "plugins")
	assetsDir = filepath.Join(appDir, "assets")
	libDir := filepath.Join(appDir, "lib")
	tools.MkDirs(loggerDir, configDir, libDir, pluginsDir)

	// 加载配置
	fmt.Printf("configDir: %v\n", configDir)
	a.config = config.New(configDir)
	err = a.config.Load()
	if err != nil {
		log.Fatalf("Start: %s", err.Error())
	}

	// 创建日志
	appLogger, err := logger.New(loggerDir)
	if err != nil {
		log.Fatal("init: ", err.Error())
	}
	a.logger = appLogger

	// 初始化软件信息
	a.appInfo = NewAppInfo()
	a.logger.Infof("当前版本%s", a.appInfo.version)

	return a

}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	go func() {
		// 添加托盘
		a.logger.Info("加载系统托盘")
		systray.Run(a.systemTray, func() {})
	}()

	go func() {
		// 任务队列
		a.logger.Info("任务队列加载中")
		a.taskQueue = task.New()
	}()

	go func() {
		// 注册事件
		a.logger.Info("注册事件")
		registerEvents(a)
	}()

	go func() {
		// 加载本地插件
		a.logger.Info("加载插件")
		a.loadPlugins()
	}()

	go func() {
		// 消息注册
		a.notification = notify.New(ctx, a.logger)

	}()

	// 缓存
	a.logger.Info("缓存器加载中")
	a.cache = NewCache()

}

func (a *App) Shutdown(ctx context.Context) {
	// 如果刚运行就关闭 有可能资源泄露

	// 关闭插件
	// for _, plugin := range a.plugins {
	// 	if plugin.State == 1 {
	// 		plugin.Service.Shutdown(context.Background(), nil)
	// 	}
	// }

	// 保存配置

	// 关闭托盘
	systray.Quit()
}

// 系统托盘
func (a *App) systemTray() {
	// iconPath := filepath.Join(assetsDir, "icon.ico")

	// iconData, err := os.ReadFile(iconPath)
	// if err != nil {
	// 	a.logger.Info("加载托盘图标失败")
	// }

	systray.SetIcon(iconData)

	show := systray.AddMenuItem("显示", "Show The Window")
	hide := systray.AddMenuItem("隐藏", "Hide The Window")
	systray.AddSeparator()
	exit := systray.AddMenuItem("退出", "Quit The Program")

	show.Click(func() { runtime.WindowShow(a.ctx) })
	hide.Click(func() { runtime.WindowHide(a.ctx) })
	exit.Click(func() { os.Exit(0) })

	systray.SetOnClick(func(menu systray.IMenu) { runtime.WindowShow(a.ctx) })
	systray.SetOnRClick(func(menu systray.IMenu) { menu.ShowMenu() })

}

// 加载插件
func (a *App) loadPlugins() {
	dirs, err := os.ReadDir(pluginsDir)
	if err != nil {
		log.Fatal(globals.ErrFileRead.Error())
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			// 加载插件主体信息
			pluginManifestPath := filepath.Join(pluginsDir, dir.Name(), "manifest.json")
			manifestData, err := os.ReadFile(pluginManifestPath)
			if err != nil {
				a.logger.Infof(globals.ErrFileRead.Error())
				continue
			}
			pluginDir := filepath.Join(pluginsDir, dir.Name())
			manifest := plugin.NewManifest(pluginDir)

			err = json.Unmarshal(manifestData, manifest)
			if err != nil {
				a.logger.Infof(globals.ErrConfigConversion.Error())
				continue
			}

			// 加载插件配置
			pc := &config.PluginConfig{}
			pluginConfig, ok := a.config.PluginConfigs[manifest.ID]
			if ok {
				pc = pluginConfig
			}
			manifest.PluginConfig = pc

			var p plugin.Plugin

			if manifest.Type == "downloader" {
				dp := plugin.NewDownloader(manifest)

				if pc.Enable {
					dp.Run(context.Background())
				}
				p = dp
			}

			// 加载时, 需要绑定插件配置地址
			a.config.PluginConfigs[manifest.ID] = pc
			a.plugins[manifest.ID] = p
		}
	}
}

func (a *App) registerPlugin(m *plugin.Manifest, pluginConfig *config.PluginConfig) (plugin.Plugin, error) {

	var p plugin.Plugin
	// 先写个下载器的
	if m.Type == "downloader" {
		p := plugin.NewDownloader(m)
		p.Manifest.PluginConfig.Settings = pluginConfig.Settings

		// 运行插件
		if pluginConfig.Enable {
			err := p.Run(context.Background())
			if err != nil {
				a.logger.Warnf(globals.ErrPluginRun.Error())
			}

			// 延迟加载插件
			go func() {
				err = p.Init(context.Background())
				if err != nil {
					a.logger.Warnf(globals.ErrPluginRun.Error())
				}
				runtime.EventsEmit(a.ctx, "plugin-update", p)
			}()

		}

	} else {
		return nil, errors.New("没有匹配的插件类型")
	}

	return p, nil

}
