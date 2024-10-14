package config

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"google.golang.org/grpc/metadata"
)

// 注入系统配置元数据 用于与插件系统通信
func (c *Config) InjectMetadata(ctx context.Context) context.Context {
	v := reflect.ValueOf(c).Elem()
	t := v.Type()

	// Iterate over all fields in the struct
	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name
		fieldValue := v.Field(i).Interface()

		// Convert field value to string (consider different types)
		var valueStr string
		switch val := fieldValue.(type) {
		case string:
			valueStr = val
		case int:
			valueStr = fmt.Sprintf("%d", val)
		case bool:
			valueStr = fmt.Sprintf("%t", val)
		default:
			valueStr = fmt.Sprintf("%v", val)
		}

		ctx = metadata.AppendToOutgoingContext(ctx, "system."+strings.ToLower(fieldName), valueStr)
	}

	return ctx
}
