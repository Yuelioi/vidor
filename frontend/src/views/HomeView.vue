<template>
  <div class="container mx-auto h-full overflow-hidden items-center">
    <div class="h-full w-full flex flex-col items-center justify-center gap-12">
      <div class="join w-full flex min-w-[350px] items-start justify-center">
        <input
          class="input basis-3/5 w-full h-full input-bordered join-item"
          v-model="link"
          placeholder="Link" />
        <button
          class="btn btn-primary join-item text-white/85 dark:text-base-content/85"
          @click="extractPlaylistInfo">
          Go
        </button>
      </div>
    </div>
  </div>

  <VDialog v-model:show="showPlaylistInfo">
    <template #header><span class="text-base-content font-bold">视频信息</span></template>

    <div class="flex flex-col rounded-md space-y-6">
      <div class="w-full flex h-32 p-2 join-item">
        <div class="h-full">
          <img
            :src="videoInfo.cover"
            class="h-full object-contain rounded-md"
            alt="Video Thumbnail" />
        </div>

        <div class="flex-1 pl-4 p-2">
          <h2 class="card-title line-clamp-1">{{ videoInfo.title }}</h2>
          <p class="mt-2">作者: {{ videoInfo.author }}</p>
        </div>
      </div>

      <div class="px-4">
        <h2 class="p-2 py-3 font-bold">下载选项</h2>
        <div class="flex flex-wrap gap-y-3 items-center select-none py-4">
          <div class="px-2 basis-1/4">
            <div class="tooltip tooltip-top flex items-center" data-tip="视频">
              <label for="downloadVideo" class="flex items-center cursor-pointer">
                <span class="icon-[lucide--file-video-2] size-6"></span>
              </label>
              <input
                type="checkbox"
                id="downloadVideo"
                class="ml-2 checkbox checkbox-xs"
                v-model="configs.download_video" />
            </div>
          </div>
          <div class="px-2 basis-1/4">
            <div class="tooltip tooltip-top flex items-center" data-tip="音频">
              <label for="downloadAudio" class="flex items-center cursor-pointer">
                <span class="icon-[lucide--file-audio] size-6"></span>
              </label>
              <input
                type="checkbox"
                id="downloadAudio"
                class="ml-2 checkbox checkbox-xs"
                v-model="configs.download_audio" />
            </div>
          </div>
          <div class="px-2 basis-1/4">
            <div class="tooltip tooltip-top flex items-center" data-tip="字幕">
              <label for="downloadSubtitle" class="flex items-center cursor-pointer">
                <span class="icon-[lucide--subtitles] size-6"></span>
              </label>
              <input
                type="checkbox"
                id="downloadSubtitle"
                class="ml-2 checkbox checkbox-xs"
                v-model="configs.download_subtitle" />
            </div>
          </div>
          <div class="px-2 basis-1/4">
            <div class="tooltip tooltip-top flex items-center" data-tip="合并">
              <label for="downloadCombine" class="flex items-center cursor-pointer">
                <span class="icon-[lucide--combine] size-6"></span>
              </label>
              <input
                type="checkbox"
                id="downloadCombine"
                class="ml-2 checkbox checkbox-xs"
                v-model="configs.download_combine" />
            </div>
          </div>
        </div>

        <div class="py-4 pl-2">
          <label class="flex gap-2 items-center justify-between px-0 input input-bordered">
            <span class="px-4">魔法名称</span>
            <input
              type="text"
              class="grow"
              v-model.lazy="configs.magic_name"
              placeholder="下载文件魔法名称" />
            <button class="btn" @click="applyMagicName">应用</button>
          </label>
        </div>
      </div>

      <div class="px-4">
        <h2 class="p-2 py-3 font-bold">分P选择</h2>
        <div class="flex flex-col">
          <table class="table">
            <thead>
              <tr>
                <th>
                  <label>
                    <input
                      type="checkbox"
                      class="checkbox"
                      :checked="isSelectAll(videoInfo.tasks)"
                      @change="handleSelectedAll(videoInfo.tasks)" />
                  </label>
                </th>
                <th>序号</th>
                <th>
                  <button
                    :class="!showMagicName ? 'font-bold text-sm' : ''"
                    @click="showMagicName = !showMagicName">
                    标题
                  </button>
                  /
                  <button
                    :class="showMagicName ? 'font-bold' : ''"
                    @click="showMagicName = !showMagicName">
                    魔法名称
                  </button>
                </th>
                <template v-if="hasVideo">
                  <th>视频</th>
                </template>
                <template v-if="hasAudio">
                  <th>音频</th>
                </template>
                <template v-if="hasImage">
                  <th>图片</th>
                </template>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(task, index) in videoInfo.tasks" :key="index" class="px-2 mb-2">
                <th>
                  <label>
                    <input
                      type="checkbox"
                      class="checkbox"
                      :checked="task.selected"
                      @change="task.selected = !task.selected" />
                  </label>
                </th>
                <td>
                  {{ index + 1 }}
                </td>
                <td class="w-full">
                  <div class="font-bold" v-if="!showMagicName">
                    {{ task.title }}
                  </div>
                  <div class="font-bold" v-if="showMagicName">
                    <input
                      type="text"
                      name=""
                      class="input w-full"
                      id=""
                      v-model="task.magicName" />
                  </div>
                </td>
                <template v-for="type in mediaTypes" :key="type">
                  <media-cell :task="task" :type="type" :tasks="videoInfo.tasks" />
                </template>
              </tr>
            </tbody>
            <!-- foot -->
            <tfoot>
              <tr>
                <th></th>
                <th>Index</th>
                <th>Title</th>
                <template v-if="hasVideo">
                  <th>Video</th>
                </template>
                <template v-if="hasAudio">
                  <th>Audio</th>
                </template>
                <template v-if="hasImage">
                  <th>Image</th>
                </template>
              </tr>
            </tfoot>
          </table>
        </div>
      </div>
    </div>
    <template #footer>
      <div class="px-4">
        <div class="flex flex-wrap text-sm py-2">
          <label class="flex items-center cursor-pointer">
            <input
              type="checkbox"
              :checked="isSelectAll(videoInfo.tasks)"
              class="hidden peer"
              @change="handleSelectedAll(videoInfo.tasks)" />
            <span
              class="border-primary select-none border text-sm peer-checked:bg-primary peer-checked:text-white ml-2 px-3 py-1 rounded">
              全选
            </span>
          </label>
          <button
            class="btn btn-sm mx-4"
            @click="parsePlaylistInfo"
            :disabled="isSelectAtLessOne(videoInfo.tasks)">
            解析
          </button>
          <button class="btn btn-sm ml-auto mx-4" @click="showPlaylistInfo = false">取消</button>
          <button
            class="btn btn-primary btn-sm mr-2"
            @click="addTasks"
            :disabled="isSelectAtLessOne(videoInfo.tasks)">
            下载
          </button>
        </div>
      </div>
    </template>
  </VDialog>
</template>
<script lang="ts" setup>
import { VDialog } from '@/plugins/dialog/index.js'
import { proto, app } from '@wailsjs/go/models'

const mediaTypes = ['video', 'audio']

import { Playlist, Task } from '@/models/go'
import { ShowDownloadInfo, AddDownloadTasks, ParsePlaylist } from '@wailsjs/go/app/App'
import { MagicName } from '@/utils/util'
// import { MagicName } from '@/utils/util'

const { configs } = storeToRefs(useBasicStore())
const isDownloadBtnDisabled = ref(false)
const link = ref('https://www.bilibili.com/video/BV1k14y117di/')

const showPlaylistInfo = ref(false)
const showMagicName = ref(false)
const videoInfo = reactive<Playlist>(new Playlist(configs.value.download_dir))

const hasVideo = computed(() => {
  return videoInfo.tasks[0].segments.some((segment) => segment.mime_type === 'video')
})

const hasAudio = computed(() => {
  return videoInfo.tasks[0].segments.some((segment) => segment.mime_type === 'audio')
})
const hasImage = computed(() => {
  return videoInfo.tasks[0].segments.some((segment) => segment.mime_type === 'image')
})

// const router = useRouter()

function isSelectAll(tasks: Task[]) {
  return tasks.every((task: Task) => {
    return task.selected
  })
}
function isSelectAtLessOne(tasks: Task[]) {
  return !tasks.some((task: Task) => {
    return task.selected
  })
}

function handleSelectedAll(tasks: Task[]) {
  const status = isSelectAll(tasks)
  tasks.forEach((task) => {
    task.selected = !status
  })
}

// 获取视频信息
async function extractPlaylistInfo() {
  Message({ message: '获取视频信息中...请稍后', duration: 300 })
  const result = await ShowDownloadInfo(link.value)

  console.log(result)

  if (result.title !== '') {
    showPlaylistInfo.value = true
    Object.assign(videoInfo, result)
    selectBest(videoInfo)
  }
}

// 解析视频
function parsePlaylistInfo() {
  const ids: string[] = []

  videoInfo.tasks.forEach((task) => {
    if (task.selected) {
      ids.push(task.id)
    }
  })

  ParsePlaylist(ids).then((vi: proto.TasksResponse) => {
    if (vi.id === '') {
      Message({ message: '获取视频信息失败, 请检查设置, 以及日志文件', type: 'warn' })
    } else {
      Message({ message: '解析成功', type: 'success' })
      Object.assign(videoInfo, vi)
      console.log(videoInfo)

      selectBest(videoInfo)
    }
  })
}

// 选择最高画质
function selectBest(videoInfo: Playlist) {
  videoInfo.tasks.forEach((task) => {
    task.segments.forEach((seg) => {
      if (seg.formats.length > 0) {
        seg.formats[0].selected = true
      }
    })
  })
}

function addTasks() {
  const mps = new app.taskMap()
  mps['1'] = 1

  isDownloadBtnDisabled.value = true
  for (let i = 0; i < videoInfo.tasks.length; i++) {
    const task = videoInfo.tasks[i]

    if (!task.selected) {
      continue
    }
  }

  setTimeout(() => {
    isDownloadBtnDisabled.value = false
  }, 1000)

  AddDownloadTasks(videoInfo.tasks).then((result: boolean) => {
    console.log(result)

    // if (parts.length == 0) {
    //   Message({ message: '添加失败', type: 'warn' })
    // } else {
    //   tasks.value.push(...parts)
    //   Message({ message: '添加成功', type: 'success' })
    //   router.push({
    //     name: 'task'
    //   })
    // }
  })
}

function applyMagicName() {
  videoInfo.tasks.forEach((element, index) => {
    element.magicName = MagicName(
      configs.value.magic_name,
      videoInfo.workDir,
      element.title,
      index + 1
    )
  })
}

watch(configs.value, async () => {
  SaveConfig(configs.value).then(() => {
    console.log('保存配置成功')
  })
})
</script>

../../wailsjs/go/main/App.js
