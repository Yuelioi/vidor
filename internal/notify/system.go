package notify

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// 系统消息分发(发给前端)
type System struct {
	ctx context.Context
}

func NewSystem(ctx context.Context) *System {
	return &System{
		ctx: ctx,
	}
}

func (s *System) Send(ctx context.Context, provider, eventName, message, messageType string) {
	runtime.EventsEmit(s.ctx, eventName, &Notice{
		Message:     message,
		MessageType: messageType,
		Provider:    provider,
	})
}
