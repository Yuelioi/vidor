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
                        :src="videoInfo.Thumbnail"
                        class="h-full object-contain rounded-md"
                        alt="Video Thumbnail" />
                </div>

                <div class="flex-1 pl-4 p-2">
                    <h2 class="card-title line-clamp-1">{{ videoInfo.WorkDirName }}</h2>
                    <p class="mt-2">作者: {{ videoInfo.Author }}</p>
                    <p class="line-clamp-1 opacity-40">{{ videoInfo.PubDate }}</p>
                    <p class="line-clamp-1 opacity-40">{{ videoInfo.Description }}</p>
                </div>
            </div>

            <div class="px-4">
                <h2 class="p-2 py-3 font-bold">下载选项</h2>
                <div class="flex flex-wrap gap-y-3 items-center select-none">
                    <div class="px-2 basis-1/4">
                        <div class="tooltip tooltip-top flex items-center" data-tip="视频">
                            <label for="downloadVideo" class="flex items-center cursor-pointer">
                                <span class="icon-[lucide--file-video-2] size-6"></span>
                            </label>
                            <input
                                type="checkbox"
                                id="downloadVideo"
                                class="ml-2 checkbox checkbox-xs"
                                v-model="config.DownloadVideo" />
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
                                v-model="config.DownloadAudio" />
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
                                v-model="config.DownloadSubtitle" />
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
                                v-model="config.DownloadCombine" />
                        </div>
                    </div>
                </div>
            </div>

            <div class="overflow-x-auto px-4">
                <h2 class="p-2 py-3 font-bold">分P选择</h2>
                <table class="table">
                    <thead>
                        <tr>
                            <th>
                                <label>
                                    <input
                                        type="checkbox"
                                        class="checkbox"
                                        @change="handleSelectedAll" />
                                </label>
                            </th>
                            <th>序号</th>
                            <th>标题</th>
                            <th>时长</th>
                            <th>画质</th>
                            <th>音质</th>
                            <th>编码</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr
                            v-for="(part, index) in videoInfo.Parts"
                            :key="index"
                            class="overflow-hidden px-2 mb-2">
                            <th>
                                <label>
                                    <input type="checkbox" class="checkbox" />
                                </label>
                            </th>
                            <td>
                                {{ index + 1 }}
                            </td>
                            <td>
                                <div class="font-bold">{{ part.Title }}</div>
                            </td>
                            <td>{{ part.Duration }}</td>
                            <th>
                                <button class="btn btn-ghost btn-xs">画质</button>
                            </th>
                            <th>
                                <button class="btn btn-ghost btn-xs">音质</button>
                            </th>
                            <th>
                                <button class="btn btn-ghost btn-xs">编码</button>
                            </th>
                        </tr>
                    </tbody>
                    <!-- foot -->
                    <tfoot>
                        <tr>
                            <th></th>
                            <th>Index</th>
                            <th>Title</th>
                            <th>Duration</th>
                            <th>Video</th>
                            <th>Audio</th>
                            <th>Code</th>
                        </tr>
                    </tfoot>
                </table>
            </div>

            <div class="px-4 pb-6">
                <div class="flex flex-wrap text-sm py-4">
                    <label class="flex items-center cursor-pointer">
                        <input
                            :checked="selectAll"
                            type="checkbox"
                            class="hidden peer"
                            @change="handleSelectedAll" />
                        <span
                            class="border-primary select-none border text-sm peer-checked:bg-primary peer-checked:text-white ml-2 px-3 py-1 rounded">
                            全选
                        </span>
                    </label>
                    <button class="btn btn-sm ml-auto mx-4" @click="showPlaylistInfo = false">
                        取消
                    </button>
                    <button
                        class="btn btn-primary btn-sm mr-2"
                        @click="addTasks"
                        :disabled="
                            isDownloadBtnDisabled ||
                            selectedVideos.length == 0 ||
                            selectedFormat == ''
                        ">
                        下载
                    </button>
                </div>
            </div>
        </div>
    </VDialog>
</template>
<script lang="ts" setup>
import { VDialog } from '@/plugins/dialog/index.js'
import { Part, PlaylistInfo } from '@/models/go'
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
            videoInfo.Codecs = vi.Codecs
            videoInfo.Parts = vi.Parts
            videoInfo.Description = vi.Description
            videoInfo.PubDate = vi.PubDate

            showPlaylistInfo.value = true

            selectedFormat.value = videoInfo.Qualities[videoInfo.Qualities.length - 1]
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
