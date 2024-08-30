package plugin

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Yuelioi/vidor/internal/config"
	pb "github.com/Yuelioi/vidor/internal/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DownloadPlugin struct {
	BasePlugin
	Service pb.DownloadServiceClient `json:"-"`
}

func NewDownloader(baseDir string) *DownloadPlugin {
	bp := BasePlugin{
		PluginConfig:    &config.PluginConfig{Settings: make(map[string]string)},
		BaseDir:         baseDir,
		ManifestVersion: 0,
		Name:            "",
		Description:     "",
		Author:          "",
		Version:         "",
		HomePage:        "",
		Color:           "",
		DocsURL:         "",
		Addr:            "",
		DownloadURLs:    []string{},
		Matches:         []string{},
		Categories:      []string{},
		Tags:            []string{},
		Executable:      "",
		State:           0,
		Status:          "",
	}

	return &DownloadPlugin{
		BasePlugin: bp,
	}
}

func (p DownloadPlugin) Init() error {
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

// 关闭插件 停止进程
func (p DownloadPlugin) Kill() error {
	_, err := p.Service.Shutdown(context.Background(), nil)
	return err
}

func (p DownloadPlugin) Update(ctx context.Context) error {
	_, err := p.Service.Update(ctx, &emptypb.Empty{})
	return err
}

// 停止插件
func (p DownloadPlugin) Stop(ctx context.Context) (*Plugin, error) {
	_, err := p.Service.Shutdown(ctx, &emptypb.Empty{})
	return nil, err
}
