package notify

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// 系统消息分发(发给前端)
type SystemNotification struct {
	ctx context.Context
}

func NewSystem(ctx context.Context) *SystemNotification {
	return &SystemNotification{
		ctx: ctx,
	}
}

func (s *SystemNotification) Send(ctx context.Context, nc Notice) error {
	if nc.EventName == "" {
		return fmt.Errorf("event name is required")
	}
	runtime.EventsEmit(s.ctx, nc.EventName, nc)
	return nil
}
