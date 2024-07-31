export interface MessageData {
    Message: string
    MessageType: string
}

export interface Config {
    Theme: string
    ProxyURL: string
    UseProxy: boolean
    DownloadDir: string
    FFMPEG: string
    ScaleFactor: number
    MagicName: string
    DownloadVideo: boolean
    DownloadAudio: boolean
    DownloadSubtitle: boolean
    DownloadCombine: boolean

    SESSDATA: string

    DownloadLimit: number
}

export class Task {
    path: string // 下载后的路径
    url: string // 原始Url
    author: string
    title: string
    thumbnail: string
    createdAt: Date // 任务创建时间
    size: number
    resolution: number
    status: string // 已完成/下载中50%/队列中
    constructor() {
        this.path = ''
        this.url = ''
        this.author = ''
        this.title = ''
        this.thumbnail = ''
        this.createdAt = new Date(0)
        this.size = 0
        this.resolution = 0
        this.status = ''
    }
}
// 直接获取播放列表信息
export class PlaylistInfo {
    Url: string = ''
    Cover: string = ''
    WorkDirName: string = ''
    Author: string = ''
    PubDate: Date = new Date()
    Description: string = ''
    StreamInfos: StreamInfo[] = []
}

export class StreamInfo {
    Name = ''
    TaskID = ''
    Selected: boolean = false
    Thumbnails: Thumbnail[] = []
    Videos: Stream = new Stream()
    Audios: Stream = new Stream()
    Captions = []
}

export class Stream {
    ID = '' // youtubeID bilibiliID...
    SessionId = '' // biliCID...
    URL = ''
    Title = ''
    Description = ''
    Author = ''
    ChannelID = ''
    Views = 0
    Duration = 0 // You can convert this to a proper duration format if needed
    PublishDate = new Date()
    Formats: Format[] = []
    DASHManifestURL = '' // URI of the DASH manifest file
    HLSManifestURL = '' // URI of the HLS manifest file
}

export class Format {
    IDtag = 0 // 标签ID
    URL = '' // 链接
    MimeType = '' // video/mp4...
    Quality = '' // 质量标签
    ContentLength = 0 // 内容长度
    DurationMs = 0 // 时长
    Selected = false

    // 图片+视频
    Width = 0
    Height = 0

    // 仅视频
    FPS = 0 // FPS
    Bitrate = 0 // 码率
    AverageBitrate = 0 // 平均码率

    // 仅音频
    AudioSampleRate = '' // 音频码率
}

export class Thumbnail {
    URL = ''
    Label = ''
    Width = 0
    Height = 0
}

export class CaptionTrack {
    BaseURL = ''
    Name = ''
    LanguageCode = ''
    Kind = ''
    IsTranslatable = false
}

export class Part {
    TaskID: string
    DownloadDir: string

    Author: string
    Title: string
    Url: string
    Quality: string
    Resolution: number
    State: string
    Status: string
    Description: string

    Size: number
    Path: string

    Thumbnail: string
    CreatedAt: Date
    PubDate: Date
    DownloadPercent: number
    DownloadSpeed: number

    constructor(url: string, title: string, Thumbnail: string, quality: string) {
        this.Url = url
        this.Title = title
        this.Quality = quality
        this.Thumbnail = Thumbnail
    }
}
