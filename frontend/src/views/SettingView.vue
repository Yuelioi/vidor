<template>
  <TabsView :tabs="tabs" v-model:tabId="tabId"></TabsView>
</template>

<script setup lang="ts">
const tabId = ref(1)

import BaseSetting from '@/components/settings/BaseSetting.vue'
import DownloadSetting from '@/components/settings/DownloadSetting.vue'
import TranslateSetting from '@/components/settings/TranslateSetting.vue'
const tabs: Tab[] = [
  {
    id: 1,
    name: '基础设置',
    icon: 'icon-[ic--outline-color-lens] ',
    component: BaseSetting,
    color: 'text-primary'
  },
  {
    id: 2,
    name: '下载设置',
    icon: 'icon-[lucide--download]',
    component: DownloadSetting,
    color: 'text-secondary'
  },
  {
    id: 3,
    name: '翻译设置',
    icon: 'icon-[ic--twotone-translate]',
    component: TranslateSetting,
    color: 'text-accent'
  }
]
</script>
