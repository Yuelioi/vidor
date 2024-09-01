package plugin

import (
	"context"
	"errors"
	"fmt"
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
	return p.Manifest.Run(ctx)
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
