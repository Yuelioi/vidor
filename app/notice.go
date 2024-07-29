package app

import (
	"github.com/Yuelioi/vidor/shared"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// 基础插件回调
var callback = func(data shared.NoticeData) {
	runtime.EventsEmit(Application.ctx, data.EventName, data.Message.(*shared.Part))
}

// todo 插件通知
type appNotice struct {
	app *App
}

// todo 外部插件更新
func (notice *appNotice) ProgressUpdate(part shared.Part) {
	runtime.EventsEmit(notice.app.ctx, "updateInfo", part)
}
