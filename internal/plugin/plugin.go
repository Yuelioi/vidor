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
	Running      = iota + 1 // 插件正在运行
	NotWork                 // 插件不工作（可能因为故障或其他原因）
	Working                 // 插件正在处理任务
	Initializing            // 插件正在初始化过程中
	Stopping                // 插件正在停止过程中
	Paused                  // 插件已暂停
	Failed                  // 插件遇到错误或故障
	Disabled                // 插件被禁用
)

type Plugin interface {
	Run(ctx context.Context) error
	Check(ctx context.Context) error
	Update(ctx context.Context) error
	ShutDown(ctx context.Context) error
	Talk(ctx context.Context) error
}

func (p BasePlugin) Check(ctx context.Context) error {
	return nil
}

func (p BasePlugin) ShutDown(ctx context.Context) error {
	return nil
}

func (p BasePlugin) Talk(ctx context.Context) error {
	return nil
}

func (p BasePlugin) Run(ctx context.Context) error {
	if p.Addr == "" {
		// 手动生成本地地址
		pluginPath := filepath.Join(p.BaseDir, p.Executable)
		addr, err := getLocalAddr(pluginPath)
		if err != nil {
			return err
		}
		p.Addr = addr
	}

	conn, err := connect(p.Addr)
	if err != nil {
		return err
	}
	conn.Connect()
	return nil
}

func (p BasePlugin) Update(ctx context.Context) error {
	panic("not implemented") // TODO: Implement
}

type BasePlugin struct {
	*config.PluginConfig
	BaseDir         string   // 插件所在文件夹
	ManifestVersion int      `json:"manifest_version"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Author          string   `json:"author"`
	Version         string   `json:"version"`
	HomePage        string   `json:"homepage"`
	Color           string   `json:"color"`
	DocsURL         string   `json:"docs_url"`
	Addr            string   `json:"addr"`
	DownloadURLs    []string `json:"download_urls"`
	Matches         []string `json:"matches"`
	Categories      []string `json:"categories"`
	Tags            []string `json:"tags"`
	Executable      string   `json:"executable"` // 软件执行文件全名
	State           int      `json:"state"`      // 1.运行中 2.运行中 尚未检测通信结果 3.未启动
	Status          string   `json:"status"`     // 仅前端
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
