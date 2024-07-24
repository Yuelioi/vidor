package shared

import "github.com/hashicorp/go-plugin"

// 握手规则
var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "VIDOR_PLUGIN",
	MagicCookieValue: "DOWNLOADER",
}
