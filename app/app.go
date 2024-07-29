package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/energye/systray"

	"github.com/Yuelioi/vidor/plugins"
	"github.com/Yuelioi/vidor/shared"
	"github.com/Yuelioi/vidor/utils"
	"github.com/hashicorp/go-plugin"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/sirupsen/logrus"
)

// 软件基础信息 aa
type appInfo struct {
	name       string
	version    string
	appDir     string
	configDir  string
	pluginsDir string
	tempDir    string
	logDir     string
}

// 应用实例
type App struct {
	appInfo
	ctx         context.Context
	config      *shared.Config      // 软件配置信息
	downloaders []shared.PluginMeta // 下载器 解析链接并下载
	taskQueue   *TaskQueue          // 任务队列 用于分发任务 同一时刻只会出现一个队列
	tasks       []*Task             // 所有任务 包括不在下载范围内的
	Notice      shared.Notice
	Logger      *logrus.Logger
	Callback    shared.Callback
}

func init() {
	Application = NewApp()
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	// 初始化app信息 创建必要文件夹
	if err := a.initAppInfo(); err != nil {
		log.Fatal("init: ", err.Error())
	}

	// 创建日志
	appLogger, err := utils.CreateLogger(a.logDir)
	if err != nil {
		log.Fatal("init: ", err.Error())
	}
	logger = appLogger

	a.ctx = ctx
	a.Logger = logger

	// 添加托盘
	go func() {
		systray.Run(systemTray, func() {})
	}()

	// 加载通知
	a.Notice = &appNotice{app: a}
	a.Callback = callback

	// 加载任务列表
	logger.Info("APP: ", "任务列表正在加载...")
	go a.loadTasks()

	// 加载配置信息
	logger.Info("APP: ", "加载配置文件正在加载...")
	if err := a.loadConfig(); err != nil {
		log.Fatal("init config: ", err.Error())
	}

	// 注册事件
	registerEvents(a)

	// TODO 加载本地插件
	logger.Infof("App 加载插件 %s...", a.appInfo.configDir)
	if err := a.loadDownloaders(); err != nil {
		log.Fatal("init plugin: ", err.Error())
	}

	// 加载FFmpeg
	err = utils.SetFFmpegPath(a.config.FFMPEG)
	if err != nil {
		logger.Infof("FFmpeg 加载失败:%s", err)
	} else {
		logger.Info("FFmpeg 加载成功")
	}

	a.Logger.Info("APP 加载完毕")
}

func (a *App) Shutdown(ctx context.Context) {
	a.taskQueue = nil
	a.SaveConfig(a.config)
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

// 初始化创建一些必要文件夹
func (a *App) initAppInfo() (err error) {

	exePath, err := os.Executable()
	if err != nil {
		log.Fatal("Error:", err)
		return
	}

	appDir := filepath.Dir(exePath)

	a.appInfo = appInfo{
		name:       name,
		version:    version,
		appDir:     appDir,
		configDir:  filepath.Join(appDir, "config"),
		pluginsDir: string(filepath.Separator) + "plugins",

		tempDir: filepath.Join(appDir, "tmp"),
		logDir:  filepath.Join(appDir, "log"),
	}

	if err = utils.CreateDirs([]string{
		a.configDir, a.tempDir,
		a.pluginsDir, a.logDir}); err != nil {
		return
	}
	return
}

// 启动时加载App任务列表
func (a *App) loadTasks() error {
	tasks := make([]*Task, 0)

	configFile := filepath.Join(a.appInfo.configDir, "tasks.json")
	configData, err := os.ReadFile(configFile)
	if err != nil {
		logger.Errorf("Cannot read/find task file: %v", err)
		a.tasks = tasks
		return err
	}

	parts := make([]shared.Part, 0)
	err = json.Unmarshal(configData, &parts)
	if err != nil {
		logger.Errorf("Cannot convert task data: %v", err)
		a.tasks = tasks
		return err
	}

	for _, part := range parts {
		// 过滤掉不存在的任务
		if _, err = os.Stat(part.DownloadDir); err == nil {
			newPart := part
			tasks = append(tasks, &Task{
				part: &newPart,
			})
		}
	}
	a.tasks = tasks
	saveTasks(tasks, a.configDir)
	return nil
}

// todo 插件类型 system local
func (a *App) loadDownloaders() error {
	_plugins := make([]shared.PluginMeta, 0)

	// System Plugins
	system_plugins := plugins.SystemPlugins(a.Notice)
	_plugins = append(_plugins, system_plugins...)

	// Local Plugins
	pluginsDir := "./plugins"
	glob := "*.exe"
	pluginFiles, err := plugin.Discover(glob, pluginsDir)
	if err != nil {
		return err
	}

	for _, pluginPath := range pluginFiles {
		_downloader, err := loadLocalPlugin(pluginPath)
		if err != nil {
			logger.Error(err)
			continue
		}
		_plugins = append(_plugins, utils.Downloader2plugin(_downloader, "ThirdPart"))
	}

	a.downloaders = _plugins
	return nil
}

// 加载/创建/初始化配置
func (a *App) loadConfig() error {
	configFile := filepath.Join(a.configDir, "config.json")
	var config shared.Config

	// 检查配置文件是否存在，如果不存在则创建一个空的配置文件
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		logger.Infof("配置文件 '%s' 不存在，将创建一个空的配置文件", configFile)
		config := shared.Config{
			Theme:            "dark",
			ScaleFactor:      16,
			MagicName:        "{{Index}}-{{Title}}",
			DownloadVideo:    true,
			DownloadAudio:    true,
			DownloadSubtitle: true,
			DownloadCombine:  true,
			DownloadLimit:    5,
		}
		if err := saveConfig(a.configDir, config); err != nil {
			return fmt.Errorf("无法创建配置文件: %v", err)
		}

	}

	configData, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(configData, &config)
	if err != nil {
		return err
	}

	// 初始化下载文件夹
	if _, err := os.Stat(config.DownloadDir); os.IsNotExist(err) {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		config.DownloadDir = filepath.Join(home, "downloads")
		if err := saveConfig(a.configDir, config); err != nil {
			log.Fatal(err)
		}
	}
	a.config = &config
	return nil
}
