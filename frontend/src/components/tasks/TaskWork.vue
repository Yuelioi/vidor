<template>
  <div class="p-4 h-full flex flex-col">
    <div class="h-full flex flex-col" v-if="filteredTasks.length > 0">
      <!-- 任务顶部功能组 -->
      <div class="pb-4 pt-2 flex space-x-4 items-center">
        <div class="badge badge-lg mr-auto">任务:{{ filteredTasks.length }}</div>
        <div class=""><span class="icon-[lucide--play] size-6"></span></div>
        <div class=""><span class="icon-[lucide--pause] size-6"></span></div>

        <div class="" @click="removeAll">
          <span class="icon-[ant-design--clear-outlined] size-6"></span>
        </div>
      </div>

      <!-- 任务列表-->
      <div class="overflow-y-auto space-y-3 h-full text flex-1">
        <div
          v-for="(task, index) in filteredTasks"
          :key="index"
          class="h-24 group hover:shadow-2xl">
          <div class="card w-full h-full card-side bg-base-100 shadow-xl">
            <!-- 封面 -->
            <figure class="basis-2/12 relative" v-if="task.cover !== '1'">
              <img :src="task.cover" :alt="task.title" />
            </figure>
            <div v-else class="basis-2/12">
              <div class="h-24 skeleton shrink-0 rounded-r-none"></div>
            </div>

            <div class="card-body py-4 px-6 basis-10/12 relative">
              <!-- 第一行 -->
              <div class="flex items-center">
                <span
                  class="text-center line-clamp-1 font-bold group-hover:link"
                  @click="BrowserOpenURL(task.url)">
                  {{ task.title ? task.title : '标题正在加载中...' }}
                </span>

                <div class="text-slate-300 space-x-2 ml-auto">
                  <span class="icon-[lucide--play] size-5"></span>
                  <span class="icon-[lucide--pause] size-5"></span>
                  <div
                    @click="OpenFileWithSystemPlayer(task.work_dir)"
                    class="icon-[lucide--circle-play] size-5"></div>
                  <span
                    class="icon-[lucide--trash-2] size-5 cursor-pointer"
                    @click="removeTask(task.id)"></span>
                  <span
                    class="icon-[ic--baseline-folder-open] size-5 cursor-pointer"
                    @click="OpenExplorer(task.work_dir)"></span>
                </div>
              </div>

              <!-- 第4行 -->
              <div class="flex items-center gap-2 text-xs text-base-content/40">
                <span>{{ task.size }}</span>
                <span class="line-clamp-1">
                  {{ task.speed }}
                </span>
                <span class="ml-auto line-clamp-1">
                  {{ task.status }}
                </span>
              </div>

              <!-- 第3行 -->
              <progress
                class="progress progress-success"
                :value="task.percent"
                max="100"></progress>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="pt-2 font-bold" v-else>
      <span>还木有任务捏~</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { OpenExplorer, RemoveTask, OpenFileWithSystemPlayer } from '@wailsjs/go/app/App'
import { BrowserOpenURL } from '@wailsjs/runtime/runtime'
import { Tab } from '@/models/ui'
import { proto } from '@wailsjs/go/models'
defineProps<{ tab: Tab }>()

const { tasks } = storeToRefs(useBasicStore())

const filteredTasks = computed(() => {
  return tasks.value.filter((task) => task.state === 1)
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

const removeAll = () => {
  VAlert({ alert: '确定要清理所有任务吗(不会删除文件)' }).then((ok) => {
    if (ok) {
      console.log(filteredTasks.value)

      RemoveAllTask([]).then((ok) => {
        if (ok) {
          Message({ message: '删除任务成功', type: 'success' })
          console.log(tasks.value, 1)
          tasks.value.splice(
            0,
            tasks.value.length,
            ...subtractTaskLists(tasks.value, filteredTasks.value)
          )
          console.log(tasks.value, 2)
        } else {
          Message({ message: '删除任务失败', type: 'error' })
        }
      })
    }
  })
}

function subtractTaskLists(tasks: proto.Task[], filteredTasks: proto.Task[]): proto.Task[] {
  const filteredTaskids = new Set(filteredTasks.map((task) => task.id))
  return tasks.filter((task) => !filteredTaskids.has(task.id))
}
</script>
