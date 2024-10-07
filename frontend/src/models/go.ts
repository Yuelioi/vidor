import { plugin } from '@wailsjs/go/models'

export class Plugin extends plugin.Manifest {
  lock?: boolean = true
  downloaded?: boolean = false // 仅插件市场
}

export class Notice {
  eventName?: string
  content?: string
  noticeType?: string
  provider?: string
}
