package plugin

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
