<template>
  <th class="relative group">
    <div tabindex="0" role="button" class="min-w-max btn btn-sm btn-outline">
      {{ currentFormat(formats) }}
    </div>
    <ul
      tabindex="0"
      class="dropdown-content duration-500 transition-opacity absolute rounded-lg opacity-0 invisible group-hover:visible group-hover:opacity-100 group-hover:block top-[80%] w-full menu bg-base-300 z-[1]">
      <template v-for="(format, index) in formats" :key="index">
        <button
          :class="format.selected ? '' : ''"
          class="py-1 btn btn-xs"
          @click="selectFormat(formats, index)">
          {{ format.label }}
        </button>
      </template>
    </ul>
  </th>
</template>

<script setup lang="ts">
import { Segment, Task, Format } from '@/models/go'

const props = defineProps<{
  task: Task
  type: string
}>()

function filterSegments(segments: Segment[], mimeType: string) {
  const result = segments.find((segment) => segment.mime_type === mimeType) || new Segment()
  console.log(result)

  return result
}

const formats = computed(() => {
  return filterSegments(props.task.segments, props.type).formats
})
const currentFormat = (formats: Format[]) => {
  const selectedFormat = formats.find((format) => format.selected)
  return selectedFormat ? selectedFormat.label : '选择'
}

function selectFormat(formats: Format[], index: number) {
  // 正常选择单个label
  for (let i = 0; i < formats.length; i++) {
    if (i == index) {
      formats[i].selected = true
    } else {
      formats[i].selected = false
    }
  }
}
</script>
