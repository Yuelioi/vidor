<template>
  <div class="flex h-full w-full text-sm select-none space-x-4 px-4">
    <div
      class="p-2 flex items-center space-x-2 cursor-pointer rounded-lg"
      v-for="(tab, index) in tabs"
      :class="tabStyle(tab, index)"
      @click="tabChange(index)"
      :key="tab.id">
      <span :class="tab.icon" class="size-5"></span>
      <span>{{ tab.name }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Tab } from '@/models/ui'

const tabId = defineModel<number>('tabId')

defineProps<{ tabs: Tab[] }>()

function tabStyle(tab: Tab, index: number) {
  let styleText = `hover:` + tab.color
  if (index + 1 == tabId.value) {
    styleText += ' shadow-md shadow-base-content/20 '
  }
  return styleText
}

function tabChange(index: number) {
  tabId.value = index + 1
}
</script>
