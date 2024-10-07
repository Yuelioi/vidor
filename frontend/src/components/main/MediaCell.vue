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
          @click="selectFormat(formats, index, $event)">
          {{ format.label }}
        </button>
      </template>
    </ul>
  </th>
</template>

<script setup lang="ts">
import { proto } from '@wailsjs/go/models'

const props = defineProps<{
  tasks: proto.Task[]
  task: proto.Task
  type: string
}>()

function filterSegments(segments: proto.Segment[], mimeType: string) {
  const result = segments.find((segment) => segment.mime_type === mimeType) || new proto.Segment()
  return result
}

const formats = computed(() => {
  return filterSegments(props.task.segments, props.type).formats
})
const currentFormat = (formats: proto.Format[]) => {
  const selectedFormat = formats.find((format) => format.selected)
  return selectedFormat ? selectedFormat.label : '选择'
}

function selectFormat(formats: proto.Format[], index: number, event: MouseEvent) {
  if (event.shiftKey) {
    // 批量选择
    const currentFormat = formats[index]
    props.tasks.forEach((task: proto.Task) => {
      task.segments.forEach((seg: proto.Segment) => {
        if (seg.mime_type === props.type) {
          seg.formats.forEach((format: proto.Format) => {
            if (format.label === currentFormat.label) {
              format.selected = true
            } else {
              format.selected = false
            }
          })
        }
      })
    })
  } else {
    // 正常选择单个label
    for (let i = 0; i < formats.length; i++) {
      if (i == index) {
        formats[i].selected = true
      } else {
        formats[i].selected = false
      }
    }
  }
}
</script>
