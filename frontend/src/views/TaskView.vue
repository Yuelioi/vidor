<template>
    <div class="flex h-full">
        <div class="left w-36 border-r-4 border-zinc-200 dark:border-base-100">
            <!-- 任务左侧边栏 -->
            <TabLeftSideBar :tabs="tabs" v-model:tabId="tabId"></TabLeftSideBar>
        </div>
        <div class="right flex-1 h-full w-full border-l-4 border-neutral-100 dark:border-base-300">
            <!-- 任务内容区域 -->
            <TaskTab v-if="tabId == 1" :filteredTasks="downloadingTasks" :tabId="tabId"></TaskTab>
            <TaskTab v-else-if="tabId == 2" :filteredTasks="queueTasks" :tabId="tabId"></TaskTab>
            <TaskTab v-else :filteredTasks="finishedTasks" :tabId="tabId"></TaskTab>
        </div>
    </div>
</template>

<script setup lang="ts">
const { tasks } = storeToRefs(useBasicStore())

const downloadingTasks = computed(() => {
    return tasks.value.filter((task) => task.State === '下载中')
})
const queueTasks = computed(() => {
    return tasks.value.filter((task) => task.State === '队列中')
})
const finishedTasks = computed(() => {
    return tasks.value.filter((task) => task.State === '已完成')
})

const tabId = ref(1)
const tabs = [
    {
        id: 1,
        name: '下载中',
        icon: 'icon-[ic--round-downloading]'
    },
    {
        id: 2,
        name: '队列中',
        icon: 'icon-[lucide--square-stack]'
    },
    {
        id: 3,
        name: '已完成',
        icon: 'icon-[ic--outline-expand-circle-down]'
    }
]
</script>
