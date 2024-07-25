<template>
    <div class="flex h-full">
        <div class="left w-36 border-r-4 border-zinc-200 dark:border-base-100">
            <!-- 任务左侧边栏 -->
            <TabLeftSideBar :tabs="tabs" v-model:tabId="tabId"></TabLeftSideBar>
        </div>
        <div class="right flex-1 h-full w-full border-l-4 border-neutral-100 dark:border-base-300">
            <!-- 任务内容区域 -->
            <div class="h-full" v-if="tabId == 1"><AppUse></AppUse></div>
            <div class="h-full" v-else-if="tabId == 2"><AppUpdate></AppUpdate></div>
            <div class="h-full" v-else><AppAbout></AppAbout></div>
        </div>
    </div>
</template>

<script setup lang="ts">
const tabId = ref(1)
const tabs = [
    {
        id: 1,
        name: '使用说明',
        icon: 'icon-[lucide--notebook-pen]'
    },
    {
        id: 2,
        name: '检查更新',
        icon: 'icon-[lucide--circle-arrow-up]'
    },
    {
        id: 3,
        name: '关于',
        icon: 'icon-[lucide--ghost]'
    }
]
</script>
