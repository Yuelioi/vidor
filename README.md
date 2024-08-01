纯go实现的下载器

目前支持

- [X] youtube

  - [X] 单个/列表
  - [ ] 字幕
  - [X] 封面
- [X] vimeo

  - [X] 单个
  - [ ] vimeo有列表么?
  - [ ] 字幕
- [X] bilibili (1080p+需要登录)

  - [X] 单个
  - [X] 多p
  - [X] 系列
  - [ ] 字幕

## TODO

- ffmpeg提示 放lib目录
- 获取列表 优先解析输入的链接?
- home页面允许直接修改文件名
- 字幕
- 翻译
- clean

## 自行构建

后端

`go mod tidy`

`go install github.com/wailsapp/wails/v2/cmd/wails@latest`

`wails dev`

`wails build -ldflags -H=windowsgui` 使用了 https://github.com/energye/systray

前端

`cd frontend`

`pnpm i`

## 鸣谢

ytb下载: github.com/kkdai/youtube/v2

deepl翻译: github.com/michimani/deepl-sdk-go

ffmpeg视频处理: github.com/u2takey/ffmpeg-go
