<template>
  <TabsView :tabs="tabs" v-model:tabId="tabId"></TabsView>
</template>

<script setup lang="ts">
import SystemPlugin from '@/components/plugins/SystemPlugin.vue'
import ThirdPartPlugin from '@/components/plugins/ThirdPartPlugin.vue'

const tabId = ref(1)

const tabs: Tab[] = [
  {
    id: 1,
    name: '系统插件',
    icon: 'icon-[lucide--plug-2]',
    component: SystemPlugin,
    color: 'text-primary'
  },
  {
    id: 2,
    name: '三方插件',
    component: ThirdPartPlugin,
    icon: 'icon-[lucide--sailboat]',
    color: 'text-secondary'
  }
]
</script>
