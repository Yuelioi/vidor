package app

import (
	"context"
	"log"
	"os"

	"github.com/energye/systray"

	utils "github.com/Yuelioi/vidor/internal/tools"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/sirupsen/logrus"
)

// 应用实例
type App struct {
	appInfo   AppInfo
	ctx       context.Context
	config    *Config   // 软件配置信息
	taskQueue TaskQueue // 任务队列 用于分发任务
	plugins   []*Plugin
	Logger    *logrus.Logger
}

func init() {

	Application = NewApp()
	Application.appInfo = *NewAppInfo()
	Application.config = NewConfig()

	// 创建日志
	appLogger, err := utils.CreateLogger(Application.config.logDir)
	if err != nil {
		log.Fatal("init: ", err.Error())
	}
	logger = appLogger

	Application.taskQueue = NewTaskQueue()
	Application.Logger = logrus.New()

	p, _ := LoadPlugin("bilibili", "location string", "_type string")
	RunPlugin(p)

	Application.plugins = append(Application.plugins, p)
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {

	a.ctx = ctx
	a.Logger = logger

	// 添加托盘
	go func() {
		systray.Run(systemTray, func() {})
	}()

	// 加载任务列表

	// 注册事件
	registerEvents(a)

	// TODO 加载本地插件

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
	// a.SaveConfig(a.config)
	systray.Quit()
}

// 系统托盘
func systemTray() {

	iconData, _ := os.ReadFile("./build/windows/icon.ico")

	systray.SetIcon(iconData) // read the icon from a file

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
