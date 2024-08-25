<template>
  <TabsView :tabs="tabs" v-model:tabId="tabId"></TabsView>
</template>

<script setup lang="ts">
import PluginMarket from '@/components/plugins/PluginMarket.vue'
import PluginManager from '@/components/plugins/PluginManager.vue'

const tabId = ref(1)

const tabs: Tab[] = [
  {
    id: 1,
    name: '插件市场',
    icon: 'icon-[lucide--plug-2]',
    component: PluginMarket,
    color: 'text-primary'
  },
  {
    id: 2,
    name: '插件管理',
    component: PluginManager,
    icon: 'icon-[lucide--sailboat]',
    color: 'text-secondary'
  }
]
</script>
