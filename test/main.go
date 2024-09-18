package main

import (
	"context"
	"fmt"
)

type keyType struct{}

var KeyApp = keyType{}

func main() {
	ctx := context.Background() // 或者使用其他已存在的 Context
	appValue := "YourAppValue"  // 你的应用程序相关的值
	ctxWithValue := context.WithValue(ctx, KeyApp, appValue)

	// 假设 ctxWithValue 已经被正确地传递到了这里
	if value, ok := ctxWithValue.Value(KeyApp).(string); ok {
		fmt.Println("App Value:", value)
	} else {
		fmt.Println("No App Value found in the context")
	}
}
