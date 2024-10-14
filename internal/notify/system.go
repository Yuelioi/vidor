package notify

// 系统消息分发(发给前端)

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type SystemNotification struct {
	ctx context.Context
}

func NewSystem(ctx context.Context) *SystemNotification {
	return &SystemNotification{
		ctx: ctx,
	}
}

func (s *SystemNotification) Send(nc Notice) error {
	if nc.EventName == "" {
		return fmt.Errorf("event name is required")
	}
	runtime.EventsEmit(s.ctx, nc.EventName, nc)
	return nil
}
