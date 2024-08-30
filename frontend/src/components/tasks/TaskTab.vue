<template>
  <!-- <div class="p-4 h-full flex flex-col">
        <div class="h-full" v-if="filteredTasks.length > 0">
            <div class="py-2 flex space-x-4 items-center">
                <div class="badge ">
                    任务:{{ filteredTasks.length }}
                </div>

                <div class="" v-if="tab.id == 1" @click="removeAll">
                    <span class="icon-[ic--outline-stop-circle] size-6"></span>
                </div>
                <div class="" v-else-if="tab.id == 2" @click="removeAll">
                    <span class="icon-[icon-park-outline--clear-format] size-6"></span>
                </div>
                <div class="" v-else @click="removeAll">
                    <span class="icon-[ant-design--clear-outlined] size-6"></span>
                </div>
            </div>
            <div class="overflow-y-auto space-y-3 h-full text">
                <div
                    v-for="(task, index) in filteredTasks"
                    :key="index"
                    class="bg-base-100 rounded-md">
                    <div class="flex p-2 h-20 group relative">
                        <div
                            class="absolute group-hover:opacity-0 right-2 top-3 badge badge-sm opacity-45 bg-base-200">
                            {{ task.Quality }}
                        </div>
                        <div v-if="task.Thumbnail" class="relative">
                            <img
                                class="object-contain h-full"
                                :src="task.Thumbnail"
                                :alt="task.Title" />
                            <div
                                @click="OpenFileWithSystemPlayer(task.Path)"
                                :class="tab.color"
                                class="transition-opacity duration-300 opacity-0 group-hover:opacity-100 absolute icon-[lucide--circle-play] size-8 top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2"></div>
                        </div>
                        <div v-else>
                            <div class="h-16 w-24 skeleton shrink-0"></div>
                        </div>

                        <div class="flex-1 pl-4 flex flex-col">
                            <h2 class="flex-1 font-bold line-clamp-1">
                                <span
                                    :class="'group-hover:' + tab.color"
                                    class="group-hover:link font-bold"
                                    @click="BrowserOpenURL(task.URL)">
                                    {{ task.Title ? task.Title : '标题正在加载中...' }}
                                </span>
                            </h2>

                            <div v-if="tab.id == 1">
                                <div class="flex justify-between truncate">
                                    <span class="text-sm text-gray-500">
                                        {{ task.Status }}
                                    </span>
                                    <span class="text-sm text-gray-500">
                                        {{ task.DownloadSpeed }}
                                    </span>
                                </div>
                                <progress
                                    class="progress progress-success"
                                    :value="task.DownloadPercent"
                                    max="100"></progress>
                            </div>

                            <div v-if="tab.id == 3" class="text-xs text-base-content/40">
                                <div>2024年7月24日</div>
                                <div>52.3M</div>
                            </div>
                        </div>
                        <div
                            class="ml-3 h-full flex space-x-2 items-center justify-between transition-opacity duration-300 opacity-0 group-hover:opacity-100">
                            <span
                                class="icon-[lucide--trash-2] size-8"
                                @click="removeTask(task.id)"></span>
                            <span
                                class="icon-[ic--baseline-folder-open] size-8"
                                @click="OpenExplorer(task.DownloadDir)"></span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="pt-2 font-bold" v-else>
            <span>还木有任务捏~</span>
        </div>
    </div> -->
</template>

<script setup lang="ts">
import { OpenExplorer, RemoveTask, OpenFileWithSystemPlayer } from '@wailsjs/go/app/App'
import { BrowserOpenURL } from '@wailsjs/runtime/runtime'
import { Tab } from '@/models/ui'
const props = defineProps<{ tab: Tab }>()

const { tasks } = storeToRefs(useBasicStore())

const filteredTasks = computed(() => {
  if (props.tab.id == 1) {
    return tasks.value.filter((task) => task.state === 2)
  } else if (props.tab.id == 2) {
    return tasks.value.filter((task) => task.state === 1)
  } else {
    return tasks.value.filter((task) => task.state === 3)
  }
})

const removeTask = (uid: string) => {
  RemoveTask(uid).then((ok) => {
    if (ok) {
      Message({ message: '删除任务成功', type: 'success' })
      const index = tasks.value.findIndex((task) => task.id === uid)
      if (index !== -1) {
        tasks.value.splice(index, 1)
      }
    } else {
      Message({ message: '删除任务失败', type: 'error' })
    }
  })
}
</script>
