<template>
    <div class="flex h-full">
        <div class="left w-36 border-r-4 border-zinc-200 dark:border-base-100">
            <!-- 插件左侧边栏 -->
            <TabLeftSideBar :tabs="tabs" v-model:tabId="tabId"></TabLeftSideBar>
        </div>
        <div class="right flex-1 h-full w-full border-l-4 border-neutral-100 dark:border-base-300">
            <PluginTab></PluginTab>
        </div>
    </div>
</template>

<script setup lang="ts">
const tabId = ref(1)
const tabs = [
    {
        id: 1,
        name: '系统插件',
        icon: 'icon-[lucide--plug-2]'
    },
    {
        id: 2,
        name: '三方插件',
        icon: 'icon-[lucide--sailboat]'
    }
]
</script>
