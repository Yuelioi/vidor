package app

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	_ "embed"

	"github.com/Yuelioi/vidor/internal/config"
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

// 应用实例
type App struct {
	ctx          context.Context
	appDirs      *AppDirs                    // 软件路径
	appInfo      AppInfo                     // 软件信息
	config       *config.Config              // 软件配置信息
	taskQueue    *task.TaskQueue             // 任务队列 用于分发任务
	manager      *plugin.PluginManager       // 插件管理系统
	cache        *Cache                      // 缓存
	notification *notify.LoggingNotification // 消息分发
	logger       *logrus.Logger              // 日志系统
}

// 初始化必要文件夹
func initDirs() (*AppDirs, error) {
	root, err := tools.ExeDir()
	if err != nil {
		return nil, err
	}

	appDirs := &AppDirs{
		AppRoot: root,
		Libs:    filepath.Join(root, "libs"),
		Temps:   filepath.Join(root, "temps"),
		Plugins: filepath.Join(root, "plugins"),
		Configs: filepath.Join(root, "configs"),
		Logs:    filepath.Join(root, "logs"),
	}

	err = tools.MkDirs(appDirs.Libs, appDirs.Temps, appDirs.Plugins, appDirs.Configs, appDirs.Logs)
	if err != nil {
		return nil, err
	}

	return appDirs, nil
}

func NewApp() (*App, error) {
	a := &App{}

	// 初始化文件夹
	dirs, err := initDirs()
	if err != nil {
		return nil, err
	}
	a.appDirs = dirs

	// 加载配置
	cf, err := config.New(a.appDirs.Configs)
	if err != nil {
		return nil, err
	}
	a.config = cf
	if err := a.config.Load(); err != nil {
		return nil, err
	}

	// 创建日志
	appLogger, err := logger.New(a.appDirs.Logs)
	if err != nil {
		return nil, err
	}
	a.logger = appLogger

	// 初始化软件信息
	a.appInfo = NewAppInfo()
	a.logger.Infof("加载成功, 当前版本%s", a.appInfo.version)
	return a, nil
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	a.manager = plugin.NewPluginManager(ctx)

	var wg sync.WaitGroup
	wg.Add(5)

	go func() {
		defer wg.Done()
		// 添加托盘
		a.logger.Info("加载系统托盘")
		systray.Run(a.systemTray, func() {})
	}()

	go func() {
		defer wg.Done()
		// 任务队列
		a.logger.Info("任务队列加载中")
		a.taskQueue = task.New()
	}()

	go func() {
		defer wg.Done()
		// 注册事件
		a.logger.Info("注册事件")
		registerEvents(a)
	}()

	go func() {
		defer wg.Done()
		// 加载本地插件
		a.logger.Info("加载插件")
		if err := a.loadPlugins(); err != nil {
			a.logger.Infof("加载插件失败,err:%s", err)
		}
	}()

	go func() {
		defer wg.Done()
		// 消息注册
		systemNotification := notify.NewSystem(a.ctx)
		a.notification = notify.NewLoggingNotification(a.logger, systemNotification)

	}()

	// 缓存
	a.logger.Info("缓存器加载中")
	a.cache = NewCache()

	wg.Wait()
	a.logger.Info("应用启动完成")
}

// TODO 关闭
func (a *App) Shutdown(ctx context.Context) {
	// 如果刚运行就关闭 有可能资源泄露

	// 关闭插件
	// for _, plugin := range a.plugins {
	// 	if plugin.State == 1 {
	// 		plugin.Service.Shutdown(context.Background(), nil)
	// 	}
	// }

	// 清理tmp文件夹
	err := tools.CleanDir(a.appDirs.Temps)
	if err != nil {
		a.logger.Warnf("清理临时文件夹失败: %s", err)
	}

	// 保存配置

	// 关闭托盘
	systray.Quit()
}

// 系统托盘
func (a *App) systemTray() {
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
func (a *App) loadPlugins() error {
	dirs, err := os.ReadDir(a.appDirs.Plugins)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			// 读取插件信息
			pluginDir := filepath.Join(a.appDirs.Plugins, dir.Name())
			pluginManifestPath := filepath.Join(pluginDir, "manifest.json")
			manifestData, err := os.ReadFile(pluginManifestPath)
			if err != nil {
				a.logger.Infof("读取插件文件失败: %s", err)
				continue
			}

			// 解析插件信息
			manifest := plugin.NewManifest(pluginDir)
			err = json.Unmarshal(manifestData, manifest)
			if err != nil {
				a.logger.Infof("插件配置转换失败: %s", err)
				continue
			}

			// 注册插件
			if err := a.manager.Register(manifest); err != nil {
				a.logger.Warnf("注册插件失败: %s", err)
				continue
			}
		}
	}
	return nil
}
