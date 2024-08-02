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
    url: string // 原始URL
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
    URL: string = ''
    Cover: string = ''
    WorkDirName: string = ''
    Author: string = ''
    PubDate: Date = new Date()
    Description: string = ''
    StreamInfos: StreamInfo[] = []
}

export class StreamInfo {
    ID: string = '' // youtubeID bilibiliID...
    SessionId: string = '' // biliCID...
    URL = ''
    Name = ''
    MagicName: string = ''
    Selected: boolean = false

    Thumbnails: Thumbnail[] = []
    Videos: Format[] = []
    Audios: Format[] = []
    Captions = []
}

export class StreamQuality {
    IDtag = 0
    Label = ''
}

export class Format {
    IDtag = 0 // 标签ID
    Quality = '' // 质量标签
    Selected = false
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
    URL: string
    Quality: string
    Resolution: number
    State: string
    Status: string
    Description: string

    Thumbnail: string
    Video: StreamQuality = new StreamQuality()
    AQuality: StreamQuality = new StreamQuality()
    Subtitle: string = '' // todo

    Size: number
    Path: string

    CreatedAt: Date
    PubDate: Date
    DownloadPercent: number
    DownloadSpeed: number

    constructor(url: string, title: string, Thumbnail: string, quality: string) {
        this.URL = url
        this.Title = title
        this.Quality = quality
        this.Thumbnail = Thumbnail
    }
}
