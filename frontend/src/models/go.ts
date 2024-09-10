import { plugin } from '@wailsjs/go/models'
import { proto } from '@wailsjs/go/models'

export class Plugin extends plugin.Manifest {
  lock?: boolean = true
  downloaded?: boolean = false // 仅插件市场
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
  percent: number = 0
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
