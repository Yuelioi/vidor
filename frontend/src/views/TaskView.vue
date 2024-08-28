<template>
  <TabsView :tabs="tabs" v-model:tabId="tabId"></TabsView>
</template>

<script setup lang="ts">
import TaskTab from '@/components/tasks/TaskTab.vue'

import TaskWork from '@/components/tasks/TaskWork.vue'
import TaskQueue from '@/components/tasks/TaskQueue.vue'
import TaskFinished from '@/components/tasks/TaskFinished.vue'

const tabId = ref(1)
const tabs: Tab[] = [
  {
    id: 1,
    name: '下载中',
    icon: 'icon-[ic--round-downloading]',
    component: TaskWork,
    color: 'text-primary'
  },
  {
    id: 2,
    name: '队列中',
    icon: 'icon-[lucide--square-stack]',
    component: TaskQueue,
    color: 'text-secondary'
  },
  {
    id: 3,
    name: '已完成',
    icon: 'icon-[ic--outline-expand-circle-down]',
    component: TaskFinished,
    color: 'text-accent'
  }
]
</script>
