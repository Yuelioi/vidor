package app

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/energye/systray"
	"github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// 应用实例
type App struct {
	ctx       context.Context
	appInfo   AppInfo
	config    *Config            // 软件配置信息
	taskQueue TaskQueue          // 任务队列 用于分发任务
	plugins   map[string]*Plugin // 插件
	cache     *Cache             // 缓存
	logger    *logrus.Logger
}

func NewApp() *App {
	return &App{
		plugins: make(map[string]*Plugin),
	}
}

func (app *App) Startup(ctx context.Context) {
	app.ctx = ctx

	// 加载配置
	c, err := NewConfig()
	if err != nil {
		log.Fatalf("Start: %s", err.Error())
	}
	app.config = c

	// 创建日志
	appLogger, err := createLogger(app.config.logDir)
	if err != nil {
		log.Fatal("init: ", err.Error())
	}

	app.logger = appLogger

	// 初始化软件信息
	app.appInfo = NewAppInfo()
	app.logger.Infof("当前版本%s", app.appInfo.version)

	go func() {
		// 添加托盘
		app.logger.Info("加载系统托盘")
		systray.Run(app.systemTray, func() {})
	}()

	go func() {
		// 任务队列
		app.logger.Info("任务队列加载中")
		app.taskQueue = NewTaskQueue()
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

	// 加载FFmpeg
	// err = utils.SetFFmpegPath(app.config.FFMPEG)
	// if err != nil {
	// 	logger.Infof("FFmpeg 加载失败:%s", err)
	// } else {
	// 	logger.Info("FFmpeg 加载成功")
	// }

	app.logger.Info("APP 加载完毕")
}

func (app *App) Shutdown(ctx context.Context) {
	// 如果刚运行就关闭 有可能资源泄露

	// 关闭插件
	for _, plugin := range app.plugins {
		if plugin.State == 1 {
			plugin.service.Shutdown(context.Background(), nil)
		}
	}

	// 保存配置

	// 关闭托盘
	systray.Quit()
}

// 系统托盘
func (app *App) systemTray() {
	iconPath := filepath.Join(app.config.assetsDir, "icon.ico")
	iconData, err := os.ReadFile(iconPath)
	if err != nil {
		app.logger.Info("加载托盘图标失败")
	}

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

// 基于链接获取下载器
func (app *App) selectPlugin(url string) (*Plugin, error) {
	for _, plugin := range app.plugins {
		for _, match := range plugin.Matches {
			reg, err := regexp.Compile(match)
			if err != nil {
				return nil, errors.New("插件正则表达式编译失败: " + err.Error())
			}
			if reg.MatchString(url) {
				return plugin, nil
			}
		}
	}
	return nil, pluginNotFound
}

// 加载插件
func (app *App) loadPlugins() {
	dirs, err := os.ReadDir(app.config.pluginsDir)
	if err != nil {
		log.Fatalf(fileReadFailed.Error())
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			pluginManifestPath := filepath.Join(app.config.pluginsDir, dir.Name(), "manifest.json")
			manifestData, err := os.ReadFile(pluginManifestPath)
			if err != nil {
				app.logger.Infof(fileReadFailed.Error())
				continue
			}

			plugin := NewPlugin()
			err = json.Unmarshal(manifestData, plugin)
			if err != nil {
				app.logger.Infof(configConversionFailed.Error())
				continue
			}

			// 加载插件配置
			pluginConfig, ok := app.config.PluginConfigs[plugin.ID]
			if ok {
				plugin.Settings = pluginConfig.Settings

				// 运行插件
				if pluginConfig.Enable {
					err = plugin.Run(app.config)
					if err != nil {
						app.logger.Warnf(pluginRunFailed.Error())
						continue
					}

					// 延迟加载插件
					go func() {
						err = plugin.Init()
						if err != nil {
							app.logger.Warnf(pluginInitFailed.Error())
						}
						runtime.EventsEmit(app.ctx, "plugin-update", plugin)
					}()

				}
			}
			// 加载时, 需要绑定插件配置地址
			app.config.PluginConfigs[plugin.ID] = plugin.PluginConfig

			app.logger.Infof("成功加载插件:%v\n\n", plugin.PluginConfig)
			app.plugins[plugin.ID] = plugin
		}
	}
}
