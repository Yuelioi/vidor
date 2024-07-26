<template>
    <div class="flex h-full">
        <div class="left w-36 border-r-4 border-zinc-200 dark:border-base-100">
            <!-- 任务左侧边栏 -->
            <TabLeftSideBar :tabs="tabs" v-model:tabId="tabId"></TabLeftSideBar>
        </div>
        <div
            class="right flex-1 h-full w-full flex border-l-4 border-neutral-100 dark:border-base-300">
            <TabContent :tabs="tabs" :tabId="tabId"></TabContent>
        </div>
    </div>
</template>

<script setup lang="ts">
import { Tab } from '@/models/ui'
import AppUse from '@/components/info/AppUse.vue'
import AppUpdate from '@/components/info/AppUpdate.vue'
import AppAbout from '@/components/info/AppAbout.vue'

const tabId = ref(1)
const tabs: Tab[] = [
    {
        id: 1,
        name: '使用说明',
        icon: 'icon-[lucide--notebook-pen]',
        component: AppUse
    },
    {
        id: 2,
        name: '检查更新',
        icon: 'icon-[lucide--circle-arrow-up] ',
        component: AppUpdate
    },
    {
        id: 3,
        name: '关于',
        icon: 'icon-[lucide--ghost]',
        component: AppAbout
    }
]
</script>
