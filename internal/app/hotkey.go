package app

// APP 全局热键注册

import (
	"log"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"golang.design/x/hotkey"
)

func registerHotkey(a *App) {

	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.Key1)
	err := hk.Register()
	if err != nil {
		log.Printf("hotkey: failed to register hotkey: %v", err)
		return
	}

	log.Printf("hotkey: %v is registered\n", hk)

	for {
		<-hk.Keydown()
		runtime.WindowShow(a.ctx)
	}

}
