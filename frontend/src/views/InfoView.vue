<template>
  <TabsView :tabs="tabs" v-model:tabId="tabId"></TabsView>
</template>

<script setup lang="ts">
import AppUse from '@/components/info/AppUse.vue'
import AppUpdate from '@/components/info/AppUpdate.vue'
import AppAbout from '@/components/info/AppAbout.vue'

const tabId = ref(1)
const tabs: Tab[] = [
  {
    id: 1,
    name: '使用说明',
    icon: 'icon-[lucide--notebook-pen]',
    component: AppUse,
    color: 'text-primary'
  },
  {
    id: 2,
    name: '检查更新',
    icon: 'icon-[lucide--circle-arrow-up] ',
    component: AppUpdate,
    color: 'text-secondary'
  },
  {
    id: 3,
    name: '关于',
    icon: 'icon-[lucide--ghost]',
    component: AppAbout,
    color: 'text-accent'
  }
]
</script>
