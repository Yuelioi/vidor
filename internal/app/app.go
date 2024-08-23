package app

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/energye/systray"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/sirupsen/logrus"
)

// 应用实例
type App struct {
	appInfo   AppInfo
	ctx       context.Context
	config    *Config   // 软件配置信息
	taskQueue TaskQueue // 任务队列 用于分发任务
	plugins   []*Plugin // 插件
	cache     *Cache    // 缓存
	Logger    *logrus.Logger
}

func init() {

	Application = NewApp()
	Application.appInfo = NewAppInfo()
	Application.config = NewConfig()

	// 创建日志
	appLogger, err := createLogger(Application.config.logDir)
	if err != nil {
		log.Fatal("init: ", err.Error())
	}
	logger = appLogger

	Application.taskQueue = NewTaskQueue()
	Application.cache = NewCache()
	Application.Logger = appLogger

	logger.Info("初始化完毕")

	// p, err = RunPlugin(p)
	// if err != nil {
	// 	logger.Errorf("运行插件失败%v", err)
	// }

}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	logger.Info("启动中")

	a.ctx = ctx
	a.Logger = logger

	// 添加托盘
	logger.Info("加载系统托盘")
	go func() {
		systray.Run(systemTray, func() {})
	}()

	// 加载任务列表

	// 注册事件
	logger.Info("注册事件")
	registerEvents(a)

	// 加载本地插件
	logger.Info("加载插件")
	a.loadPlugins()

	// 加载FFmpeg
	// err = utils.SetFFmpegPath(a.config.FFMPEG)
	// if err != nil {
	// 	logger.Infof("FFmpeg 加载失败:%s", err)
	// } else {
	// 	logger.Info("FFmpeg 加载成功")
	// }

	a.Logger.Info("APP 加载完毕")
}

func (a *App) Shutdown(ctx context.Context) {
	a.taskQueue = nil
	systray.Quit()
}

// 系统托盘
func systemTray() {
	iconPath := filepath.Join(Application.config.assetsDir, "icon.ico")
	iconData, err := os.ReadFile(iconPath)
	if err != nil {
		logger.Info("加载托盘图标失败")
	}

	systray.SetIcon(iconData)

	show := systray.AddMenuItem("显示", "Show The Window")
	hide := systray.AddMenuItem("隐藏", "Hide The Window")
	systray.AddSeparator()
	exit := systray.AddMenuItem("退出", "Quit The Program")

	show.Click(func() { runtime.WindowShow(Application.ctx) })
	hide.Click(func() { runtime.WindowHide(Application.ctx) })
	exit.Click(func() { os.Exit(0) })

	systray.SetOnClick(func(menu systray.IMenu) { runtime.WindowShow(Application.ctx) })
	systray.SetOnRClick(func(menu systray.IMenu) { menu.ShowMenu() })

}

func (a *App) selectPlugin(url string) (*Plugin, error) {
	// implement plugin selection logic based on URL scheme or content type
	return a.plugins[0], nil
}

// Load Plugins
func (a *App) loadPlugins() {

	dirs, err := os.ReadDir(a.config.pluginsDir)
	if err != nil {
		logger.Infof("无法加载插件文件夹:%s", err)
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			pluginConfigPath := filepath.Join(a.config.pluginsDir, dir.Name(), "manifest.json")
			data, err := os.ReadFile(pluginConfigPath)
			if err != nil {
				logger.Infof("无法读取插件配置:%s", err)
			}

			plugin := &Plugin{}
			err = json.Unmarshal(data, plugin)
			if err != nil {
				logger.Infof("插件配置转换失败:%s", err)
			}
			logger.Infof("成功加载插件:%s", plugin.Name)

			a.plugins = append(a.plugins, plugin)

		}
	}
}
