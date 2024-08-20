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
                v-model="config.system.download_video" />
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
                v-model="config.system.download_audio" />
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
                v-model="config.system.download_subtitle" />
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
                v-model="config.system.download_combine" />
            </div>
          </div>
        </div>

        <div class="py-4 pl-2">
          <label class="flex gap-2 items-center justify-between px-0 input input-bordered">
            <span class="px-4">魔法名称</span>
            <input
              type="text"
              class="grow"
              v-model.lazy="config.system.magic_name"
              placeholder="下载文件魔法名称" />
            <button class="btn" @click="applyMagicName">应用</button>
          </label>
        </div>
      </div>

      <div class="overflow-y-hidden px-4">
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
                      :checked="isSelectAll(videoInfo.stream_infos)"
                      @change="handleselectedAll(videoInfo.stream_infos)" />
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
                <th>视频</th>
                <th>音频</th>
                <th>图片</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="(streamInfo, index) in videoInfo.stream_infos"
                :key="index"
                class="px-2 mb-2">
                <th>
                  <label>
                    <input
                      type="checkbox"
                      class="checkbox"
                      :checked="streamInfo.selected"
                      @change="streamInfo.selected = !streamInfo.selected" />
                  </label>
                </th>
                <td>
                  {{ index + 1 }}
                </td>
                <!-- 标题 -->
                <td class="w-full">
                  <div class="font-bold" v-if="!showMagicName">
                    {{ streamInfo.title }}
                  </div>
                  <div class="font-bold" v-if="showMagicName">
                    <input
                      type="text"
                      name=""
                      class="input w-full"
                      id=""
                      v-model="streamInfo.magicName" />
                  </div>
                </td>
                <!-- 视频 -->
                <th class="relative group">
                  <div tabindex="0" role="button" class="min-w-max btn btn-sm btn-outline">
                    {{ currentFormat(videoSegments(streamInfo.streams, 'video').formats) }}
                  </div>

                  <ul
                    tabindex="0"
                    class="dropdown-content duration-500 transition-opacity absolute rounded-lg opacity-0 invisible group-hover:visible group-hover:opacity-100 group-hover:block top-[80%] w-full menu bg-base-300 z-[1]">
                    <template
                      v-for="(format, index) in videoSegments(streamInfo.streams, 'video').formats"
                      :key="index">
                      <button
                        :class="format.selected ? '' : ''"
                        class="py-1 btn btn-xs"
                        @click="selectFormat(streamInfo.streams, index)">
                        {{ format.label }}
                      </button>
                    </template>
                  </ul>
                </th>

                <th>
                  <button class="btn btn-ghost btn-xs min-w-max">暂不支持</button>
                </th>
              </tr>
            </tbody>
            <!-- foot -->
            <tfoot>
              <tr>
                <th></th>
                <th>Index</th>
                <th>Title</th>
                <th>Video</th>
                <th>Audio</th>
                <th>Code</th>
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
              :checked="isSelectAll(videoInfo.stream_infos)"
              class="hidden peer"
              @change="handleselectedAll(videoInfo.stream_infos)" />
            <span
              class="border-primary select-none border text-sm peer-checked:bg-primary peer-checked:text-white ml-2 px-3 py-1 rounded">
              全选
            </span>
          </label>
          <button
            class="btn btn-sm mx-4"
            @click="parsePlaylistInfo"
            :disabled="isSelectAtLessOne(videoInfo.stream_infos)">
            解析
          </button>
          <button class="btn btn-sm ml-auto mx-4" @click="showPlaylistInfo = false">取消</button>
          <button
            class="btn btn-primary btn-sm mr-2"
            @click="addTasks"
            :disabled="isSelectAtLessOne(videoInfo.stream_infos)">
            下载
          </button>
        </div>
      </div>
    </template>
  </VDialog>
</template>
<script lang="ts" setup>
import { VDialog } from '@/plugins/dialog/index.js'
import { proto } from '@wailsjs/go/models'

import { Part, Playlist, Task, Segment, Format } from '@/models/go'
import { ShowDownloadInfo, AddDownloadTasks, ParsePlaylist } from '@wailsjs/go/app/App'
// import { MagicName } from '@/utils/util'

const { config, tasks } = storeToRefs(useBasicStore())
const isDownloadBtnDisabled = ref(false)
const link = ref('')

const showPlaylistInfo = ref(false)
const showMagicName = ref(false)
const videoInfo = reactive<Playlist>(new Playlist())

const router = useRouter()

function isSelectAll(streamInfos: Task[]) {
  return streamInfos.every((streamInfo: Task) => {
    return streamInfo.selected
  })
}
function isSelectAtLessOne(streamInfos: Task[]) {
  return !streamInfos.some((streamInfo: Task) => {
    return streamInfo.selected
  })
}

function handleselectedAll(streamInfos: Task[]) {
  const status = isSelectAll(streamInfos)

  streamInfos.forEach((streamInfo) => {
    streamInfo.selected = !status
  })
}

function extractPlaylistInfo() {
  Message({ message: '获取视频信息中...请稍后', duration: 300 })
  ShowDownloadInfo(link.value).then((vi: proto.VideoInfoResponse) => {
    if (vi.title == '') {
      Message({ message: '获取视频信息失败, 请检查设置, 以及日志文件', type: 'warn' })
    } else {
      showPlaylistInfo.value = true
      console.log(vi)
      Object.assign(videoInfo, vi)
    }
  })
}

function parsePlaylistInfo() {
  ParsePlaylist([]).then((vi: proto.ParseResponse) => {
    if (vi.id == '') {
      Message({ message: '获取视频信息失败, 请检查设置, 以及日志文件', type: 'warn' })
    } else {
      Message({ message: '解析成功', type: 'success' })

      Object.assign(videoInfo, vi)

      selectBest(videoInfo)
    }
  })
}

function videoSegments(streams: Segment[], mimeType: string) {
  for (const stream of streams) {
    if ((stream.mime_type = mimeType)) {
      return stream
    }
  }
  return new Segment()
}

// 选择最高画质
function selectBest(videoInfo: Playlist) {
  videoInfo.stream_infos.forEach((element) => {
    element.streams
    // element.Videos[0].selected = true
    // element.Audios[0].selected = true
  })
}

function addTasks() {
  isDownloadBtnDisabled.value = true
  // const parts: Part[] = []
  for (let i = 0; i < videoInfo.stream_infos.length; i++) {
    const streamInfo = videoInfo.stream_infos[i]

    if (!streamInfo.selected) {
      continue
    }

    // const part = new Part(
    //     video['Url'],
    //     video['Title'],
    //     video['Thumbnail'],
    //     selectedFormat.value
    // )
    // parts.push(part)
  }

  setTimeout(() => {
    isDownloadBtnDisabled.value = false
  }, 1000)

  AddDownloadTasks(videoInfo.stream_infos).then((parts: Part[]) => {
    console.log(parts)

    if (parts.length == 0) {
      Message({ message: '添加失败', type: 'warn' })
    } else {
      tasks.value.push(...parts)
      Message({ message: '添加成功', type: 'success' })
      router.push({
        name: 'task'
      })
    }
  })
}

function applyMagicName() {
  // videoInfo.stream_infos.forEach((element, index) => {
  //   // element.magicNamee = MagicName(
  //   //   config.value.system.magic_name,
  //   //   videoInfo.WorkDirName,
  //   //   element.Name,
  //   //   index + 1
  //   // )
  // })
}

const currentFormat = (formats: Format[]) => {
  const selectedFormat = formats.find((format) => format.selected)
  return selectedFormat ? selectedFormat.label : '选择'
}

function selectFormat(formats: Format[], index: number) {
  for (let i = 0; i < formats.length; i++) {
    if (i == index) {
      formats[i].selected = true
    } else {
      formats[i].selected = false
    }
  }
  console.log(formats)
}

watch(config.value, async () => {
  SaveConfig(config.value).then(() => {
    console.log('保存配置成功')
  })
})
</script>

../../wailsjs/go/main/App.js
