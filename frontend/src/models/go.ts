import { plugin } from '@wailsjs/go/models'
import { proto } from '@wailsjs/go/models'

export class Plugin extends plugin.Plugin {
  lock?: boolean = true
}

export class Playlist {
  title?: string
  cover?: string
  author?: string
  tasks?: Task[] = []
}

export class Task extends proto.Task {
  magicName: string = ''
  state: number = 0
  status: string = ''
  speed: string = ''
  size: string = ''
  duration: string = ''
  percent: string = ''
  selected?: boolean = false
  segments?: Segment[]
}

export class Segment extends proto.Segment {
  selected?: boolean = false
  formats?: Format[]
}

export class Format extends proto.Format {
  selected?: boolean = false
}
