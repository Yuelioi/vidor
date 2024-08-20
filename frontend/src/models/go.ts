export class Config {
    system: SystemConfig = new SystemConfig()
    plugins: PluginConfig[] = []
}

export class SystemConfig {
    theme: string = ''
    scale_factor: number = 0
    proxy_url: string = ''
    use_proxy: boolean = false
    magic_name: string = ''
    download_dir: string = ''
    download_video: boolean = false
    download_audio: boolean = false
    download_subtitle: boolean = false
    download_combine: boolean = false
    download_limit: number = 0
}

export class PluginConfig {
    manifest_version: number = 0
    name: string = ''
    description: string = ''
    author: string = ''
    version: string = ''
    url: string = ''
    docs_url: string = ''
    download_url: string = ''
    matches: string[] = []
    settings: string[] = []
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
    }
}
