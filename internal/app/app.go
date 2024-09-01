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
	ctx       context.Context
	location  string // 软件路径
	appInfo   AppInfo
	config    *config.Config           // 软件配置信息
	taskQueue *task.TaskQueue          // 任务队列 用于分发任务
	plugins   map[string]plugin.Plugin // 插件
	cache     *Cache                   // 缓存
	logger    *logrus.Logger
}

func NewApp() *App {
	app := &App{
		plugins: make(map[string]plugin.Plugin),
	}

	appDir, err := tools.ExeDir()
	if err != nil {
		log.Fatal()
	}

	appDir, _ = filepath.Abs(appDir)
	app.location = appDir

	// 初始化文件夹
	loggerDir := filepath.Join(appDir, "logs")
	configDir := filepath.Join(appDir, "configs")
	pluginsDir = filepath.Join(appDir, "plugins")
	assetsDir = filepath.Join(appDir, "assets")
	libDir := filepath.Join(appDir, "lib")
	tools.MkDirs(loggerDir, configDir, libDir, pluginsDir)

	// 加载配置
	fmt.Printf("configDir: %v\n", configDir)
	app.config = config.New(configDir)
	err = app.config.Load()
	if err != nil {
		log.Fatalf("Start: %s", err.Error())
	}

	// 创建日志
	appLogger, err := logger.New(loggerDir)
	if err != nil {
		log.Fatal("init: ", err.Error())
	}

	app.logger = appLogger

	// 初始化软件信息
	app.appInfo = NewAppInfo()
	app.logger.Infof("当前版本%s", app.appInfo.version)

	return app

}

func (app *App) Startup(ctx context.Context) {
	app.ctx = ctx

	go func() {
		// 添加托盘
		app.logger.Info("加载系统托盘")
		systray.Run(app.systemTray, func() {})
	}()

	go func() {
		// 任务队列
		app.logger.Info("任务队列加载中")
		app.taskQueue = task.New()
	}()

	go func() {
		// 注册事件
		app.logger.Info("注册事件")
		registerEvents(app)
	}()

	go func() {
		// 加载本地插件
		app.logger.Info("加载插件")
		app.loadPlugins()
	}()

	// 缓存
	app.logger.Info("缓存器加载中")
	app.cache = NewCache()

}

func (app *App) Shutdown(ctx context.Context) {
	// 如果刚运行就关闭 有可能资源泄露

	// 关闭插件
	// for _, plugin := range app.plugins {
	// 	if plugin.State == 1 {
	// 		plugin.Service.Shutdown(context.Background(), nil)
	// 	}
	// }

	// 保存配置

	// 关闭托盘
	systray.Quit()
}

// 系统托盘
func (app *App) systemTray() {
	// iconPath := filepath.Join(assetsDir, "icon.ico")

	// iconData, err := os.ReadFile(iconPath)
	// if err != nil {
	// 	app.logger.Info("加载托盘图标失败")
	// }

	systray.SetIcon(iconData)

	show := systray.AddMenuItem("显示", "Show The Window")
	hide := systray.AddMenuItem("隐藏", "Hide The Window")
	systray.AddSeparator()
	exit := systray.AddMenuItem("退出", "Quit The Program")

	show.Click(func() { runtime.WindowShow(app.ctx) })
	hide.Click(func() { runtime.WindowHide(app.ctx) })
	exit.Click(func() { os.Exit(0) })

	systray.SetOnClick(func(menu systray.IMenu) { runtime.WindowShow(app.ctx) })
	systray.SetOnRClick(func(menu systray.IMenu) { menu.ShowMenu() })

}

// 加载插件
func (app *App) loadPlugins() {
	dirs, err := os.ReadDir(pluginsDir)
	if err != nil {
		log.Fatal(globals.ErrFileRead.Error())
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			pluginManifestPath := filepath.Join(pluginsDir, dir.Name(), "manifest.json")
			manifestData, err := os.ReadFile(pluginManifestPath)
			if err != nil {
				app.logger.Infof(globals.ErrFileRead.Error())
				continue
			}
			pluginDir := filepath.Join(pluginsDir, dir.Name())
			manifest := plugin.NewManifest(pluginDir)

			// p := plugin.NewDownloader(pluginDir)
			err = json.Unmarshal(manifestData, manifest)
			if err != nil {
				app.logger.Infof(globals.ErrConfigConversion.Error())
				continue
			}

			// 加载插件配置
			var pc *config.PluginConfig
			pluginConfig, ok := app.config.PluginConfigs[manifest.ID]
			if ok {
				pc = pluginConfig
			} else {
				pc, err := config.New(app.appInfo.ConfigDir)

			}

			// 加载时, 需要绑定插件配置地址
			app.config.PluginConfigs[p.PluginConfig.ID] = p.PluginConfig

			app.logger.Infof("成功加载插件:%v\n\n", p.PluginConfig)
			app.plugins[p.PluginConfig.ID] = p
		}
	}
}

func (app *App) registerPlugin(m *plugin.Manifest, pluginConfig *config.PluginConfig) (plugin.Plugin, error) {

	var p plugin.Plugin
	// 先写个下载器的
	if m.Type == "downloader" {
		p := plugin.NewDownloader(m)
		p.Manifest.PluginConfig.Settings = pluginConfig.Settings

		// 运行插件
		if pluginConfig.Enable {
			err := p.Run(context.Background())
			if err != nil {
				app.logger.Warnf(globals.ErrPluginRun.Error())
			}

			// 延迟加载插件
			go func() {
				err = p.Init(context.Background())
				if err != nil {
					app.logger.Warnf(globals.ErrPluginRun.Error())
				}
				runtime.EventsEmit(app.ctx, "plugin-update", p)
			}()

		}

	} else {
		return nil, errors.New("没有匹配的插件类型")
	}

	return p, nil

}
