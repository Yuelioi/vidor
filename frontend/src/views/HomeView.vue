<template>
    <div class="container mx-auto h-full overflow-hidden items-center">
        <div class="h-full w-full flex flex-col items-center justify-stretch">
            <div
                class="font-[Racer] basis-1/2 text-[2.75rem] select-none flex items-center justify-center">
                <div class="overflow-hidden">
                    <div class="text-primary relative">
                        <h1>Vidor</h1>
                    </div>
                </div>
            </div>
            <div class="join flex w-full min-w-[350px] items-start justify-center">
                <input
                    class="input basis-3/5 w-full h-full input-bordered join-item"
                    v-model="link"
                    placeholder="Link" />
                <button class="btn btn-primary join-item" @click="extractPlaylistInfo">Go</button>
            </div>
        </div>
    </div>

    <VDialog v-model:show="showPlaylistInfo">
        <template #header>视频信息</template>

        <div class="flex flex-col rounded-md space-y-4">
            <div class="w-full flex cursor-pointer h-32 p-2 join-item">
                <div class="h-full">
                    <img
                        :src="videoInfo.Thumbnail"
                        class="h-full object-contain rounded-md"
                        alt="Video Thumbnail" />
                </div>

                <div class="flex-1 pl-4 p-2">
                    <div
                        class="tooltip tooltip-top flex items-center"
                        :data-tip="videoInfo.WorkDirName">
                        <h2 class="card-title line-clamp-2">{{ videoInfo.WorkDirName }}</h2>
                    </div>
                    <p class="mt-2">作者: {{ videoInfo.Author }}</p>
                </div>
            </div>

            <div class="flex flex-wrap px-4 py-4 gap-y-3 items-center justify-between select-none">
                <div class="px-2">
                    <div class="tooltip tooltip-top flex items-center" data-tip="视频">
                        <label for="downloadVideo" class="flex items-center cursor-pointer">
                            <span class="icon-[lucide--file-video-2]"></span>
                        </label>
                        <input
                            type="checkbox"
                            id="downloadVideo"
                            class="ml-2 checkbox checkbox-xs"
                            v-model="config.DownloadVideo" />
                    </div>
                </div>
                <div class="px-2">
                    <div class="tooltip tooltip-top flex items-center" data-tip="音频">
                        <label for="downloadAudio" class="flex items-center cursor-pointer">
                            <span class="icon-[lucide--file-audio]"></span>
                        </label>
                        <input
                            type="checkbox"
                            id="downloadAudio"
                            class="ml-2 checkbox checkbox-xs"
                            v-model="config.DownloadAudio" />
                    </div>
                </div>
                <div class="px-2">
                    <div class="tooltip tooltip-top flex items-center" data-tip="字幕">
                        <label for="downloadSubtitle" class="flex items-center cursor-pointer">
                            <span class="icon-[lucide--subtitles]"></span>
                        </label>
                        <input
                            type="checkbox"
                            id="downloadSubtitle"
                            class="ml-2 checkbox checkbox-xs"
                            v-model="config.DownloadSubtitle" />
                    </div>
                </div>
                <div class="px-2">
                    <div class="tooltip tooltip-top flex items-center" data-tip="合并">
                        <label for="downloadCombine" class="flex items-center cursor-pointer">
                            <span class="icon-[lucide--combine]"></span>
                        </label>
                        <input
                            type="checkbox"
                            id="downloadCombine"
                            class="ml-2 checkbox checkbox-xs"
                            v-model="config.DownloadCombine" />
                    </div>
                </div>
            </div>

            <div class="flex flex-wrap px-4 py-4 gap-y-3">
                <div
                    v-for="format in videoInfo.Qualities"
                    :key="format"
                    class="px-2 form-control sm:basis-1/4 md:basis--1/5 lg:basis--1/6 xl:basis--1/6">
                    <label class="cursor-pointer flex items-center">
                        <input
                            type="radio"
                            :value="format"
                            v-model="selectedFormat"
                            name="format"
                            class="radio radio-primary radio-sm" />
                        <span class="label-text ml-1">{{ format }}</span>
                    </label>
                </div>
            </div>

            <div class="flex flex-wrap text-sm px-2 py-4">
                <div class="flex flex-wrap w-full">
                    <div
                        v-for="(video, index) in videoInfo.Parts"
                        :key="index"
                        class="w-1/2 sm:w-1/3 md:w-1/4 lg:w-1/5 xl:w-1/6 px-2 mb-2">
                        <label class="flex items-center cursor-pointer space-x-2">
                            <div
                                class="tooltip tooltip-top w-full text-left"
                                :data-tip="video.Title">
                                <input
                                    type="checkbox"
                                    class="hidden peer"
                                    :value="video"
                                    v-model="selectedVideos" />

                                <div
                                    class="peer-checked:bg-primary border peer-checked:text-white truncate w-full px-2 py-1 rounded">
                                    {{ video.Title }}
                                </div>
                            </div>
                        </label>
                    </div>
                </div>
            </div>

            <div class="flex flex-wrap text-sm px-2 py-4">
                <label class="flex items-center cursor-pointer">
                    <input
                        :checked="selectAll"
                        type="checkbox"
                        class="hidden peer"
                        @change="handleSelectedAll" />
                    <span
                        class="border-secondary border text-sm peer-checked:bg-secondary peer-checked:text-white ml-2 px-3 py-1 rounded">
                        全选
                    </span>
                </label>
                <button class="btn btn-sm ml-auto mx-4" @click="showPlaylistInfo = false">
                    取消
                </button>
                <button
                    class="btn btn-primary btn-sm"
                    @click="addTasks"
                    :disabled="
                        isDownloadBtnDisabled || selectedVideos.length == 0 || selectedFormat == ''
                    ">
                    下载
                </button>
            </div>
        </div>
    </VDialog>
</template>
<script lang="ts" setup>
import { VDialog } from '@/plugins/dialog/index.js'
import { Part, PlaylistInfo } from '@/models/task'
import { ShowDownloadInfo, AddDownloadTasks } from '@wailsjs/go/app/App'

const { config, tasks } = storeToRefs(useBasicStore())
const isDownloadBtnDisabled = ref(false)
const link = ref('')
const selectedFormat = ref('')

const showPlaylistInfo = ref(false)
const videoInfo = reactive<PlaylistInfo>(new PlaylistInfo())

const selectedVideos = ref([])
const selectAll = ref(false)

const router = useRouter()

function handleSelectedAll() {
    if (selectedVideos.value.length === videoInfo.Parts.length) {
        selectedVideos.value = []
    } else {
        selectedVideos.value = videoInfo.Parts.map((video) => video)
    }
}

function extractPlaylistInfo() {
    Message({ message: '获取视频信息中...请稍后', duration: 300 })
    ShowDownloadInfo(link.value).then((vi: PlaylistInfo) => {
        if (vi.WorkDirName == '') {
            Message({ message: '获取视频信息失败, 请检查设置, 以及日志文件', type: 'warn' })
        } else {
            videoInfo.Url = vi.Url
            videoInfo.Thumbnail = vi.Thumbnail
            videoInfo.WorkDirName = vi.WorkDirName
            videoInfo.Author = vi.Author
            videoInfo.Qualities = vi.Qualities
            videoInfo.Parts = vi.Parts
            showPlaylistInfo.value = true

            selectedFormat.value = videoInfo.Qualities[0]
        }
    })
}

function addTasks() {
    isDownloadBtnDisabled.value = true
    const parts: Part[] = []
    for (let i = 0; i < selectedVideos.value.length; i++) {
        const video = selectedVideos.value[i]

        const part = new Part(
            video['Url'],
            video['Title'],
            video['Thumbnail'],
            selectedFormat.value
        )
        parts.push(part)
    }

    setTimeout(() => {
        isDownloadBtnDisabled.value = false
    }, 1000)

    AddDownloadTasks(parts, videoInfo.WorkDirName).then((parts: Part[]) => {
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
watch(config.value, async () => {
    SaveConfig(config.value).then(() => {
        console.log('保存配置成功')
    })
})
</script>

../../wailsjs/go/main/App.js
