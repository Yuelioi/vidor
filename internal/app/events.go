package app

// APP 接受前端事件
// 前端 >> 后端

import (
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func registerEvents(a *App) {

	runtime.EventsOn(a.ctx, "ceshi", func(optionalData ...interface{}) {
		fmt.Print(optionalData)
	})

}
