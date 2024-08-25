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
    Object.assign(config.value, fetchedConfig)
    console.log('加载配置')
    console.log(config.value)
  }

  // 加载任务
  const fetchedTasks = (await GetTaskParts()) as Part[]
  tasks.value.splice(0, tasks.value.length, ...fetchedTasks)

  // 切换主题
  switchTheme(config.value.system.theme)

  // 设置字体大小
  const scale = Math.min(Math.max(config.value.system.scale_factor, 12), 32)
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
