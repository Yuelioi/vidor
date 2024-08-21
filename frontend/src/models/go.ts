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
  tasks?: Task[] = []
}

export class Task extends proto.Task {
  selected: boolean = false
  magicName = ''
  segments: Segment[] = []
}

export class Segment extends proto.Segment {
  selected: boolean = false
  formats?: Format[]
}

export class Format extends proto.Format {
  selected: boolean = false
}
