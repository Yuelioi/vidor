package plugin

// 插件注入数据 用于与插件系统通信

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

func InjectMetadata(ctx context.Context, settings map[string]string) context.Context {
	for key, setting := range settings {
		ctx = metadata.AppendToOutgoingContext(ctx, "plugin."+strings.ToLower(key), setting)
	}
	return ctx
}
