<template>
    <div class="flex h-full">
        <div class="left w-36 border-r-4 border-zinc-200 dark:border-base-100">
            <!-- 任务左侧边栏 -->
            <TabLeftSideBar :tabs="tabs" v-model:tabId="tabId"></TabLeftSideBar>
        </div>
        <div class="right flex-1 h-full w-full border-l-4 border-neutral-100 dark:border-base-300">
            <!-- 任务内容区域 -->
            <TabContent :tabs="tabs" :tabId="tabId"></TabContent>
        </div>
    </div>
</template>

<script setup lang="ts">
const tabId = ref(1)

import BaseSetting from '@/components/settings/BaseSetting.vue'
import DownloadSetting from '@/components/settings/DownloadSetting.vue'
import ThirdPart from '@/components/settings/ThirdPart.vue'
const tabs = [
    {
        id: 1,
        name: '基础设置',
        icon: 'icon-[ic--outline-color-lens] ',
        component: BaseSetting
    },
    {
        id: 2,
        name: '下载设置',
        icon: 'icon-[lucide--download]',
        component: DownloadSetting
    },
    {
        id: 3,
        name: '三方数据',
        icon: 'icon-[lucide--send] size-6',
        component: ThirdPart
    }
]
</script>
