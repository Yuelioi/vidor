package shared

import "github.com/hashicorp/go-plugin"

// 握手规则
var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "VIDOR_PLUGIN",
	MagicCookieValue: "DOWNLOADER",
}

// 未解析状态
var DefaultFormat = Format{
	IDtag:    9999,
	Quality:  "尚未解析",
	Selected: true,
}
