package main

const (
	apiURL = "https://api.bilibili.com"
)

type userStatus int

const (
	NoLogin userStatus = iota
	Login
	Vip
)



// https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/docs/video/videostream_url.md
// {ID: 6, Label: "240P"},      // 仅 MP4 格式支持, 仅 platform=html5 时有效
// {ID: 16, Label: "360P"},     // 流畅
// {ID: 32, Label: "480P"},     // 清晰
// {ID: 64, Label: "720P"},     // 高清, WEB 端默认值, B站前端需要登录才能选择，但是直接发送请求可以不登录就拿到 720P 的取流地址, 无 720P 时则为 720P60
// {ID: 74, Label: "720P60"},   // 高帧率, 登录认证
// {ID: 80, Label: "1080P"},    // 高清, TV 端与 APP 端默认值, 登录认证
// {ID: 112, Label: "1080P+"},  // 高码率, 大会员认证
// {ID: 116, Label: "1080P60"}, // 高帧率, 大会员认证
// {ID: 120, Label: "4K"},      // 需要 fnval&128=128 且 fourk=1, 大会员认证
// {ID: 125, Label: "HDR 真彩色"}, // 仅支持 DASH 格式, 需要 fnval&64=64, 大会员认证
// {ID: 126, Label: "杜比视界"},    // 仅支持 DASH 格式, 需要 fnval&512=512, 大会员认证
// {ID: 127, Label: "8K"},      // 仅支持 DASH 格式, 需要 fnval&1024=1024, 大会员认证
