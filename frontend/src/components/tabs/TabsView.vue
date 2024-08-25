<template>
  <div class="flex h-full flex-col w-full">
    <div class="left w-full border-zinc-200 border-b-2 py-4 dark:border-base-100">
      <TabSideBar :tabs="tabs" v-model:tabId="tabId"></TabSideBar>
    </div>
    <div
      class="right flex-1 w-full overflow-y-hidden flex border-l-4 border-neutral-100 dark:border-base-300">
      <div class="flex flex-col h-full space-y-4 p-2 flex-1 text-sm select-none overflow-y-auto">
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
