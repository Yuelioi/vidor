package plugin

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	pb "github.com/Yuelioi/vidor/internal/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DownloadPlugin struct {
	Manifest *Manifest
	Service  pb.DownloadServiceClient `json:"-"`
}

func NewDownloader(m *Manifest) *DownloadPlugin {
	return &DownloadPlugin{
		Manifest: m,
	}
}

func (p *DownloadPlugin) GetManifest() *Manifest {
	return p.Manifest
}

func (p *DownloadPlugin) Run(ctx context.Context) error {

	// 获取命令
	pluginPath := filepath.Join(p.Manifest.BaseDir, p.Manifest.Executable)

	if p.Manifest.Addr == "" {
		// 自动生成本地地址
		pluginPath := filepath.Join(p.Manifest.BaseDir, p.Manifest.ID, p.Manifest.Executable)
		addr, err := getLocalAddr(pluginPath)
		if err != nil {
			return err
		}
		p.Manifest.Addr = addr
	}

	// 调试模式 debug:port
	if strings.HasPrefix(p.Manifest.Addr, "debug") {
		addrs := strings.Split(p.Manifest.Addr, ":")
		if len(addrs) == 2 {
			cmd := exec.Command(pluginPath, "--port", addrs[1])
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

			// 启动进程
			err := cmd.Start()
			if err != nil {
				return errors.New("启动进程失败: " + err.Error())
			}
		}

	}

	// 本地启动 localhost[:port]
	if strings.HasPrefix(p.Manifest.Addr, "localhost") {
		addrs := strings.Split(p.Manifest.Addr, ":")

		// 使用端口启动(可用于调试)
		if len(addrs) == 2 {
			cmd := exec.Command(pluginPath, "--port", addrs[1])
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

			// 启动进程
			err := cmd.Start()
			if err != nil {
				return errors.New("启动进程失败: " + err.Error())
			}
		} else {

		}

	}
	// 远程
	if strings.HasPrefix(p.Manifest.Addr, "remote") {
		// ...
	}

	// 获取 exe 运行的 PID
	p.Manifest.Enable = true
	p.Manifest.State = Working

	conn, err := connect(p.Manifest.Addr)
	if err != nil {
		return err
	}
	p.Service = pb.NewDownloadServiceClient(conn)
	conn.Connect()
	return nil
}

func (p *DownloadPlugin) Check(ctx context.Context) error {
	_, err := p.Service.Check(context.Background(), nil)
	return err
}

func (p *DownloadPlugin) Shutdown(ctx context.Context) error {
	_, err := p.Service.Shutdown(context.Background(), nil)
	return err
}

func (p DownloadPlugin) Init(ctx context.Context) error {
	timeout := time.Second * 10

	ctx, cancel := context.WithTimeout(ctx, timeout)
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
				p.Manifest.State = 1
				return nil
			}
			return fmt.Errorf("连接%s失败:%s", p.Manifest.Name, err)
		}
	}
}
func (p *DownloadPlugin) Update(ctx context.Context) error {
	_, err := p.Service.Update(ctx, &emptypb.Empty{})
	return err
}

func (p *DownloadPlugin) Talk(ctx context.Context) error {
	return nil
}
