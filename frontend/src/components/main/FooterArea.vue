<template><div></div></template>

<script setup lang="ts">
const _themes = ['light', 'dark']
const { switchTheme } = useTheme(_themes)
const { configs, tasks } = storeToRefs(useBasicStore())
import { Task } from '@/models/go'

EventsOn('updateInfo', (optionalData?: Task) => {
  const index = tasks.value.findIndex((task) => task.id === optionalData?.id)
  if (index !== -1 && optionalData) {
    tasks.value.splice(index, 1, optionalData)
  }
})

EventsOn('system.message', (messageData: MessageData) => {
  Message({ message: messageData['message'], type: messageData['messageType'] })
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
    Object.assign(configs.value, fetchedConfig)
    console.log('加载配置')
    console.log(configs.value)
  }

  // 加载任务
  const fetchedTasks = (await GetTaskParts()) as Task[]
  tasks.value.splice(0, tasks.value.length, ...fetchedTasks)

  // 切换主题
  switchTheme(configs.value.system.theme)

  // 设置字体大小
  const scale = Math.min(Math.max(configs.value.system.scale_factor, 12), 32)
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
