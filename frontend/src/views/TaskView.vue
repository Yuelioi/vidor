<template>
  <TabsView :tabs="tabs" v-model:tabId="tabId"></TabsView>
</template>

<script setup lang="ts">
import TaskTab from '@/components/tasks/TaskTab.vue'

const tabId = ref(1)
const tabs: Tab[] = [
  {
    id: 1,
    name: '下载中',
    icon: 'icon-[ic--round-downloading]',
    component: TaskTab,
    color: 'text-primary'
  },
  {
    id: 2,
    name: '队列中',
    icon: 'icon-[lucide--square-stack]',
    component: TaskTab,
    color: 'text-secondary'
  },
  {
    id: 3,
    name: '已完成',
    icon: 'icon-[ic--outline-expand-circle-down]',
    component: TaskTab,
    color: 'text-accent'
  }
]
</script>
