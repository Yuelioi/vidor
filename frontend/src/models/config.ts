export interface Config {
    Theme: string
    ProxyURL: string
    UseProxy: boolean
    DownloadDir: string
    FFMPEG: string
    ScaleFactor: number
    DownloadVideo: boolean
    DownloadAudio: boolean
    DownloadSubtitle: boolean
    DownloadCombine: boolean

    SESSDATA: string

    DownloadLimit: number
}
