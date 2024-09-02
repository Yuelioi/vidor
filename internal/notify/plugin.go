package notify

import (
	"context"
)

// 插件通信消息分发(发给插件)
type Plugin struct {
	ctx context.Context
}

func NewPlugin(ctx context.Context) *Plugin {
	return &Plugin{
		ctx: ctx,
	}
}
