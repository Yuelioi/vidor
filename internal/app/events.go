package app

import (
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func registerEvents(a *App) {
	runtime.EventsOn(a.ctx, "ceshi", func(optionalData ...interface{}) {
		fmt.Print(optionalData)
	})
}
