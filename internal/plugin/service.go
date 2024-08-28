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
	"time"

	"github.com/Yuelioi/vidor/internal/config"
	"github.com/Yuelioi/vidor/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

// 启动插件
func (p *Plugin) Run(c *config.Config) error {
	// 获取可用端口
	port, err := getAvailablePort()
	if err != nil {
		return errors.New("获取可用端口失败: " + err.Error())
	}
	p.Port = port

	// 获取命令
	pluginPath := filepath.Join(p.baseDir, p.Executable)
	cmd := exec.Command(pluginPath, "--port", strconv.Itoa(port))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	// 启动进程
	err = cmd.Start()
	if err != nil {
		return errors.New("启动进程失败: " + err.Error())
	}

	// 获取 exe 运行的 PID
	pid := cmd.Process.Pid
	p.PID = pid
	p.Enable = true
	p.State = 2

	conn, err := grpc.NewClient("localhost:"+strconv.Itoa(port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return errors.New("连接失败: " + err.Error())
	}

	p.Service = proto.NewDownloadServiceClient(conn)

	return nil
}

// response, _ := p.service.GetInfo(context.Background(), &pb.InfoRequest{
// 	Url: "https://www.bilibili.com/video/BV1bA411R7BN",
// })

// fmt.Printf("response: %v\n", response)

// 初始化/更新插件配置
func (p *Plugin) Init() error {
	timeout := time.Second * 10

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return errors.New("timeout: failed to connect to plugin within the given time")
		case <-ticker.C:
			_, err := p.Service.Init(ctx, nil)
			if err == nil {
				p.State = 1
				return nil
			}
			return fmt.Errorf("连接%s失败:%s", p.Name, err)
		}
	}

}

func (p *Plugin) Update() error {
	// 获取插件配置上下文
	// ctx = p.injectMetadata(c.injectMetadata(ctx))
	return nil
}

// 停止插件
func StopPlugin(p *Plugin) (*Plugin, error) {
	if p == nil {
		return nil, fmt.Errorf("插件不存在")
	}

	_, err := p.Service.Shutdown(context.TODO(), &emptypb.Empty{})
	return nil, err
}

func getAvailablePort() (int, error) {
	// 监听 "localhost:0" 让系统分配一个可用端口
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, fmt.Errorf("failed to find an available port: %v", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	return port, nil
}
