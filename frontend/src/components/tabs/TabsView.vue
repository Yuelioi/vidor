<template>
    <div class="flex h-full">
        <div class="left w-36 border-r-4 border-zinc-200 dark:border-base-100">
            <TabSideBar :tabs="tabs" v-model:tabId="tabId"></TabSideBar>
        </div>
        <div
            class="right flex-1 h-full w-full flex border-l-4 border-neutral-100 dark:border-base-300">
            <div class="flex flex-col space-y-4 p-2 pt-4 flex-1 text-sm select-none">
                <template v-for="tab in tabs" :key="tab.id">
                    <component :is="tab.component" v-if="tabId == tab.id" :tab="tab"></component>
                </template>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
const tabId = defineModel<number>('tabId')
defineProps<{ tabs: Tab[] }>()
</script>
