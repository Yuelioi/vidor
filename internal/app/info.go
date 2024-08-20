package app

var (
	version = "0.0.1" // 软件版本号
	name    = "vidor" // 软件名

)

// 软件基础信息 aa
type AppInfo struct {
	name    string
	version string
}

func NewAppInfo() *AppInfo {
	return &AppInfo{
		name:    name,
		version: version,
	}
}
