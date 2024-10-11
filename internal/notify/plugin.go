package notify

import (
	"context"

	"github.com/Yuelioi/vidor/internal/plugin"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// 插件通信消息分发(发给插件)
type PluginNotification struct {
	ctx context.Context
}

func NewPluginNotification(ctx context.Context) *PluginNotification {
	return &PluginNotification{
		ctx: ctx,
	}
}

func (s *PluginNotification) Send(m plugin.Manifest) error {
	runtime.EventsEmit(s.ctx, "system.plugin", m)
	return nil
}
