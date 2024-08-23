<template>
  <div class="h-full">
    <div class="w-full h-full p-4 pr-6 space-y-6 overflow-x-hidden overflow-y-auto select-none">
      <div
        class="p-4 bg-base-100 rounded-xl overflow-hidden flex flex-col hover:shadow-base-content/20 hover:shadow-xl border border-base-300">
        <div class="flex items-center pb-4 pl-2" :class="tab.color">
          <span class="size-6" :class="tab.icon"></span>
          <span class="ml-2 font-bold align-middle">{{ tab.name }}</span>
        </div>
        <div class="m-2 mb-4 space-y-6">
          <slot></slot>
        </div>
      </div>

      <div
        v-if="$slots.description"
        class="p-4 rounded-xl overflow-hidden flex flex-col border border-base-300">
        <slot name="description"></slot>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{ tab: Tab }>()
</script>
