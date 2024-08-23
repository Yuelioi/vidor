<template>
  <div class="flex flex-col">
    <div v-for="(plugin, index) in plugins" :key="plugin.id">
      <div class="h-full">
        <div
          class="card w-full h-full p-4 pr-6 space-y-6 overflow-x-hidden overflow-y-auto select-none"
          :class="{ hovered: hoveredIndex === index }"
          @mouseover="hoveredIndex = index"
          @mouseout="hoveredIndex = null">
          <div
            class="p-4 bg-base-100 rounded-xl overflow-hidden flex flex-col hover:shadow-base-content/20 hover:shadow-xl border border-base-300">
            <div class="flex items-center pb-4 pl-2">
              <span class="size-6" :class="tab.icon"></span>
              <span
                class="ml-2 font-bold align-middle text-[var(--hover-color)]"
                :style="{ '--hover-color': hoveredIndex == index ? plugin.color : '' }">
                {{ plugin.name }}
              </span>
              <span class="ml-auto">
                <template v-if="plugin.id"></template>
                <span><span class="size-6 text-warning icon-[lucide--unplug]"></span></span>
                <span class="size-6 text-success icon-[ic--outline-check-circle-outline]"></span>
              </span>
            </div>
            <div class="m-2 mb-4 space-y-2">{{ plugin.description }}</div>
            <div class="divider my-0"></div>
            <div class="flex items-center">
              <span class="size-6 icon-[ic--round-play-arrow] hover:text-success"></span>
              <span class="size-6 icon-[ic--baseline-stop] hover:text-error"></span>
              <span class="ml-auto space-x-2">
                <span class="size-6 icon-[ic--round-home]"></span>
                <span class="size-6 icon-[iconoir--book-solid]"></span>
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const { plugins } = storeToRefs(useBasicStore())
const hoveredIndex = ref(0)
defineProps<{ tab: Tab }>()
</script>
