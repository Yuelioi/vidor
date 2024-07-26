import { defineStore } from 'pinia'
import { Config } from '@/models/config'
import { Part } from '@/models/task'

export const useBasicStore = defineStore('basic', () => {
    const config = reactive<Config>({
        Theme: '',
        ProxyURL: '',
        UseProxy: false,
        DownloadDir: '',
        MagicName: '',
        DownloadVideo: false,
        DownloadAudio: false,
        DownloadSubtitle: false,
        DownloadCombine: false,
        SESSDATA: '',
        FFMPEG: '',
        ScaleFactor: 16,
        DownloadLimit: 1
    })

    const tasks = reactive<Part[]>([])

    return { config, tasks }
})
