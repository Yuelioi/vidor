package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/energye/systray"

	"github.com/Yuelioi/vidor/plugins"
	"github.com/Yuelioi/vidor/shared"
	"github.com/Yuelioi/vidor/utils"
	"github.com/hashicorp/go-plugin"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/sirupsen/logrus"
)

var Application *App

var (
	version      = "0.0.1"
	name         = "vidor"
	appConfigDir = ""
	appTempDir   = ""
)

var (
	logger *logrus.Logger
)

// 软件基础信息
type appInfo struct {
	name      string
	version   string
	configDir string
}

// 应用实例
type App struct {
	appInfo
	ctx         context.Context
	config      *shared.Config
	downloaders []shared.PluginMeta
	taskQueue   *TaskQueue
	tasks       []Task // 所有任务 包括不在下载范围内的
	Notice      shared.Notice
	Logger      *logrus.Logger
	Callback    shared.Callback
}

var callback = func(data shared.NoticeData) {
	runtime.EventsEmit(Application.ctx, data.EventName, data.Message.(*shared.Part))
}

type appNotice struct {
	app *App
}

func (notice *appNotice) ProgressUpdate(part shared.Part) {
	fmt.Printf("notice.app.tasks: %v\n", notice.app.tasks)
	fmt.Printf("ProgressUpdate: %v\n", part.DownloadPercent)
	runtime.EventsEmit(notice.app.ctx, "updateInfo", part)
}

func init() {
	Application = NewApp()

	// 创建必要文件夹
	createAppDirs()

	// 创建日志
	appLogger, err := utils.CreateLogger(appTempDir)
	if err != nil {
		log.Fatal(err)
	}
	logger = appLogger
}

func NewApp() *App {
	return &App{}
}

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

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	logger.Info("App 正在初始化...")
	// 加载软件基础信息 startup创建
	a.appInfo = appInfo{
		name:      name,
		version:   version,
		configDir: appConfigDir,
	}
	a.Logger = logger

	// 测试托盘
	go func() {
		systray.Run(systemTray, func() {})
	}()

	// 加载通知
	a.Notice = &appNotice{app: a}

	a.Callback = callback

	// 加载任务列表
	logger.Infof("App 加载任务列表 %s...", a.appInfo.configDir)
	a.loadTasks()

	// 加载配置信息
	logger.Infof("App 加载配置文件 %s...", a.appInfo.configDir)
	a.loadConfig()

	// 注册事件
	registerEvents(a)

	// TODO 加载本地插件
	logger.Infof("App 加载插件 %s...", a.appInfo.configDir)
	a.loadDownloaders()

	// 加载FFmpeg
	err := utils.SetFFmpegPath(a.config.FFMPEG)
	if err != nil {
		logger.Infof("FFmpeg 加载失败:%s", err)
	} else {
		logger.Info("FFmpeg 加载成功")
	}

	a.Logger.Info("APP 加载完毕")
}

func (a *App) Shutdown(ctx context.Context) {
}

func createAppDirs() {
	tempDir := os.TempDir()
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	// 存储缩略图位置
	imgCacheDir := filepath.Join(os.Getenv("LOCALAPPDATA"), name)

	appConfigDir = filepath.Join(userConfigDir, name)
	appTempDir = filepath.Join(tempDir, name)

	if err := utils.CreateDirs([]string{imgCacheDir, appConfigDir, appTempDir}); err != nil {
		log.Fatal(err)
	}
}

// 保存单个任务
func saveTask(task *Task, tasks []Task, configDir string) error {
	parts := make([]shared.Part, 0)

	for _, tk := range tasks {
		if task.part.UID == tk.part.UID {
			parts = append(parts, *task.part)
		} else {
			parts = append(parts, *tk.part)
		}
	}

	tasksData, err := json.MarshalIndent(parts, "", "  ")
	if err != nil {
		logger.Error(err)
	}

	err = os.WriteFile(filepath.Join(configDir, "tasks.json"), tasksData, 0644)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

// 启动时加载App任务列表
func (a *App) loadTasks() {
	tasks := make([]Task, 0)

	configFile := filepath.Join(a.appInfo.configDir, "tasks.json")
	configData, err := os.ReadFile(configFile)
	if err != nil {
		logger.Errorf("Cannot read/find task file: %v", err)
		a.tasks = tasks
		return
	}

	parts := make([]shared.Part, 0)
	err = json.Unmarshal(configData, &parts)
	if err != nil {
		logger.Errorf("Cannot convert task data: %v", err)
		a.tasks = tasks
		return
	}

	for _, part := range parts {
		// 过滤掉不存在的任务
		if _, err = os.Stat(part.DownloadDir); err == nil {
			newPart := part
			tasks = append(tasks, Task{
				part: &newPart,
			})
		}
	}
	a.tasks = tasks
	saveTasks(tasks, a.configDir)
}

func (a *App) newDownloader(link string) (*shared.Downloader, error) {
	for _, downloader := range a.downloaders {
		for _, regex := range downloader.Regexs {
			if regex.MatchString(link) {
				return &downloader.Impl, nil
			}
		}
	}
	return nil, errors.New("没有对应的下载器")
}

func (a *App) loadDownloaders() {
	_plugins := make([]shared.PluginMeta, 0)

	// System Plugins
	_plugins = append(_plugins, downloader2plugin(plugins.NewBiliDownloader(a.Notice)))

	// Local Plugins
	pluginsDir := "./plugins"
	glob := "*.exe"
	pluginFiles, err := plugin.Discover(glob, pluginsDir)
	if err != nil {
		logger.Infof("Failed to discover plugins: %v", err)
	}

	for _, pluginPath := range pluginFiles {
		_downloader, err := loadLocalPlugin(pluginPath)
		if err != nil {
			logger.Error(err)
			continue
		}
		_plugins = append(_plugins, downloader2plugin(_downloader))
	}

	a.downloaders = _plugins
}

func loadLocalPlugin(pluginPath string) (shared.Downloader, error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"downloader": &shared.DownloaderRPCPlugin{},
		},
		Cmd: exec.Command(pluginPath),
	})

	rpcClient, err := client.Client()
	if err != nil {
		return nil, fmt.Errorf("error creating client for plugin %s: %v", pluginPath, err)
	}

	raw, err := rpcClient.Dispense("downloader")
	if err != nil {
		return nil, fmt.Errorf("error dispensing plugin %s: %v", pluginPath, err)
	}

	downloader, ok := raw.(shared.Downloader)
	if !ok {
		return nil, fmt.Errorf("plugin %s does not implement the expected interface", pluginPath)
	}

	// Now you can use the downloader plugin
	return downloader, nil
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
		if err := saveConfig(appConfigDir, config); err != nil {
			log.Fatal(err)
		}
	}
	a.config = &config
	return nil
}

func downloader2plugin(downloader shared.Downloader) shared.PluginMeta {
	return shared.PluginMeta{
		Name:   downloader.PluginMeta().Name,
		Regexs: downloader.PluginMeta().Regexs,
		Impl:   downloader,
	}
}

func saveConfig(configDir string, config shared.Config) error {
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(configDir, "config.json"), configData, 0644)
	if err != nil {
		return err
	}
	return nil
}

// task更新时 保存
func saveTasks(tasks []Task, configDir string) error {
	parts := make([]shared.Part, 0)

	for _, task := range tasks {
		part := *task.part
		parts = append(parts, part)
	}

	tasksData, err := json.MarshalIndent(parts, "", "  ")
	if err != nil {
		logger.Error(err)
		return err
	}

	err = os.WriteFile(filepath.Join(configDir, "tasks.json"), tasksData, 0644)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
