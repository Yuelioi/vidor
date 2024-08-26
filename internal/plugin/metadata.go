package plugin

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

func (p *Plugin) InjectMetadata(ctx context.Context) context.Context {
	for key, setting := range p.Settings {
		ctx = metadata.AppendToOutgoingContext(ctx, "plugin."+strings.ToLower(key), setting)
	}
	return ctx
}
