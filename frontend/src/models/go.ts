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
    Url: string
    Thumbnail: string
    WorkDirName: string
    Author: string
    PubDate: Date
    Description: string
    Qualities: string[]
    Codecs: string[]
    Parts: {
        URL: string
        Title: string
    }[]

    constructor() {
        this.Url = ''
        this.Thumbnail = ''
        this.WorkDirName = ''
        this.Author = ''
        this.Qualities = []
        this.Parts = []
    }
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
