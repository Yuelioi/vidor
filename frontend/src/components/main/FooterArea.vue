<template><div></div></template>

<script setup lang="ts">
import { Part } from '@/models/go'

const _themes = ['light', 'dark']
const { switchTheme } = useTheme(_themes)
const { config, tasks } = storeToRefs(useBasicStore())

EventsOn('updateInfo', (optionalData?: Part) => {
    const index = tasks.value.findIndex((task) => task.TaskID === optionalData?.TaskID)
    if (index !== -1 && optionalData) {
        tasks.value.splice(index, 1, optionalData)
    }
})

EventsOn('message', (messageData: MessageData) => {
    Message({ message: messageData['Message'], type: messageData['MessageType'] })
})

function blockWindowScale(event: KeyboardEvent) {
    if (event.ctrlKey === true || event.metaKey) {
        event.preventDefault()
    }
}

onMounted(async () => {
    // 加载配置
    const fetchedConfig = (await GetConfig()) as Config
    if (fetchedConfig) {
        config.value.Theme = fetchedConfig.Theme
        config.value.ProxyURL = fetchedConfig.ProxyURL
        config.value.UseProxy = fetchedConfig.UseProxy
        config.value.MagicName = fetchedConfig.MagicName
        config.value.DownloadDir = fetchedConfig.DownloadDir
        config.value.DownloadVideo = fetchedConfig.DownloadVideo
        config.value.DownloadAudio = fetchedConfig.DownloadAudio
        config.value.DownloadSubtitle = fetchedConfig.DownloadSubtitle
        config.value.DownloadCombine = fetchedConfig.DownloadCombine

        config.value.SESSDATA = fetchedConfig.SESSDATA

        config.value.ScaleFactor = fetchedConfig.ScaleFactor
        config.value.FFMPEG = fetchedConfig.FFMPEG
        config.value.DownloadLimit = fetchedConfig.DownloadLimit
    }

    // 加载任务
    const fetchedTasks = (await GetTaskParts()) as Part[]
    tasks.value.splice(0, tasks.value.length, ...fetchedTasks)

    // 切换主题
    switchTheme(config.value.Theme)

    // 设置字体大小
    const scale = Math.min(Math.max(config.value.ScaleFactor, 12), 32)
    document.documentElement.style.fontSize = `${scale}px`

    // 禁用页面缩放

    document.addEventListener('mousewheel', blockWindowScale, {
        capture: false,
        passive: false
    })
})

onUnmounted(() => {
    document.removeEventListener('mousewheel', blockWindowScale)
})
</script>
