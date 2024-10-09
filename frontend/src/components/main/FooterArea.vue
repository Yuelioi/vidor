<template><div></div></template>

<script setup lang="ts">
const _themes = ['light', 'dark']
const { switchTheme } = useTheme(_themes)
const { configs, tasks } = storeToRefs(useBasicStore())
import { Notice } from '@/models/go'
import { proto } from '@wailsjs/go/models'

import { WindowMinimise } from '@wailsjs/runtime'

EventsOn('system.task', (fetchData?: proto.Task) => {
  const index = tasks.value.findIndex((task) => task.id === fetchData?.id)
  if (index !== -1 && fetchData) {
    tasks.value.splice(index, 1, fetchData)
  }
})
EventsOn('system.tasks', (fetchData?: proto.Task[]) => {
  Object.assign(tasks.value, fetchData)
  Message({ message: '任务更新' })

  console.log(tasks.value)
})

EventsOn('system.notice', (messageData: Notice) => {
  Message({ message: messageData.content, type: messageData.noticeType })
})

function blockWindowScale(event: KeyboardEvent) {
  if (event.ctrlKey === true || event.metaKey) {
    event.preventDefault()
  }
}

function handleEscape(event) {
  if (event.key === 'Escape') {
    WindowMinimise()
  }
}

onMounted(async () => {
  // 加载配置
  const fetchedConfig = (await GetConfig()) as Config
  if (fetchedConfig) {
    Object.assign(configs.value, fetchedConfig)
    console.log('加载配置成功')
  }

  // 加载任务
  const fetchedTasks = (await GetTasks()) as proto.Task[]
  tasks.value.splice(0, tasks.value.length, ...fetchedTasks)
})

onMounted(() => {
  // 切换主题
  switchTheme(configs.value.theme)

  // 设置字体大小
  const scale = Math.min(Math.max(configs.value.scale_factor, 12), 32)
  document.documentElement.style.fontSize = `${scale}px`

  // 禁用页面缩放
  document.addEventListener('wheel', blockWindowScale, {
    capture: false,
    passive: false
  })

  // 加载快捷方式
  document.addEventListener('keydown', handleEscape)
})

onUnmounted(() => {
  // 注销页面缩放
  document.removeEventListener('wheel', blockWindowScale)

  // 注销快捷键
  document.removeEventListener('keydown', handleEscape)
})
</script>
