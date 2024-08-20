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

  constructor(url: string, title: string, quality: string) {
    this.URL = url
    this.Title = title
    this.Quality = quality
  }
}

import { proto } from '@wailsjs/go/models'
export class Playlist {
  title?: string
  cover?: string
  author?: string
  stream_infos?: StreamInfo[] = []
}

export class StreamInfo extends proto.StreamInfo {
  selected: boolean = false
  magicName = ''
  streams: Stream[] = []
}

export class Stream extends proto.Stream {
  selected: boolean = false
  formats?: Format[]
}

export class Format extends proto.Format {
  selected: boolean = false
}
