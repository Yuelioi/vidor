package plugin

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/Yuelioi/vidor/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	Initializing = iota + 1 // 插件正在初始化过程中
	Working                 // 插件正在处理任务
	NotWork                 // 插件不工作（可能因为故障或其他原因）
	Stopping                // 插件正在停止过程中
	Paused                  // 插件已暂停
	Failed                  // 插件遇到错误或故障
	Disabled                // 插件被禁用
)

type Plugin interface {
	GetManifest() *Manifest // 获取插件基础信息

	Run(ctx context.Context) error

	Init(ctx context.Context) error
	Check(ctx context.Context) error
	Update(ctx context.Context) error
	Shutdown(ctx context.Context) error
	Talk(ctx context.Context) error
}

type Manifest struct {
	PluginConfig    *config.PluginConfig // 配置(app生成,放在config目录)
	BaseDir         string               // 插件所在的文件夹路径(app生成)。
	ID              string               `json:"id"`
	ManifestVersion int                  `json:"manifest_version"` // 插件清单的版本号。
	Type            string               `json:"type"`             // 插件类型
	Name            string               `json:"name"`             // 插件的名称。
	Description     string               `json:"description"`      // 插件的描述信息。
	Author          string               `json:"author"`           // 插件的作者。
	Version         string               `json:"version"`          // 插件的版本号。
	HomePage        string               `json:"homepage"`         // 插件的主页地址。
	Color           string               `json:"color"`            // 可能表示插件在用户界面中的颜色标识。
	DocsURL         string               `json:"docs_url"`         // 插件文档的URL。
	Addr            string               `json:"addr"`             // 可能表示插件运行时的服务地址或端口。
	DownloadURLs    []string             `json:"download_urls"`    // 插件的下载链接列表。
	Matches         []string             `json:"matches"`          // 可能表示插件适用的内容匹配规则或模式。
	Categories      []string             `json:"categories"`       // 插件所属的类别。
	Tags            []string             `json:"tags"`             // 插件的标签，用于分类或搜索。
	Executable      string               `json:"executable"`       // 软件执行文件的全名，即启动程序的名字。
	State           int                  `json:"state"`            // 插件的状态码，可能用于内部状态管理。
	Status          string               `json:"status"`           // 插件的状态内容，可能是状态的文字描述。
}

func NewManifest(baseDir string) *Manifest {
	return &Manifest{
		BaseDir:      baseDir,
		PluginConfig: &config.PluginConfig{},
		DownloadURLs: []string{},
		Matches:      []string{},
		Categories:   []string{},
		Tags:         []string{},
	}
}

func (m *Manifest) Run(ctx context.Context) error {
	if m.Addr == "" {
		// 手动生成本地地址
		pluginPath := filepath.Join(m.BaseDir, m.Executable)
		addr, err := getLocalAddr(pluginPath)
		if err != nil {
			return err
		}
		m.Addr = addr
	}

	conn, err := connect(m.Addr)
	if err != nil {
		return err
	}
	conn.Connect()
	return nil
}

func getLocalAddr(pluginPath string) (string, error) {
	// 获取可用端口
	port, err := getAvailablePort()
	if err != nil {
		return "", errors.New("获取可用端口失败: " + err.Error())
	}

	// 获取命令
	cmd := exec.Command(pluginPath, "--port", strconv.Itoa(port))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	// 启动进程
	err = cmd.Start()
	if err != nil {
		return "", errors.New("启动进程失败: " + err.Error())
	}

	return "localhost:" + strconv.Itoa(port), nil
}

// 启动插件
func connect(addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.New("连接失败: " + err.Error())
	}
	return conn, nil
}

// 获取可用端口
func getAvailablePort() (int, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, fmt.Errorf("failed to find an available port: %v", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	return port, nil
}
