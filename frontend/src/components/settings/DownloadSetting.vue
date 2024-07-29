<template>
    <TabCard :tab="tab">
        <label class="flex items-center gap-2 input input-bordered">
            名称
            <input
                type="text"
                class="ml-2 grow"
                v-model.lazy="config.MagicName"
                placeholder="下载文件魔法名称" />

            <input type="checkbox" v-model="config.UseProxy" class="checkbox" />
        </label>
        <label class="flex items-center gap-2 input input-bordered">
            <span class="pr-2 label-text text-nowrap">并行</span>
            <input
                type="range"
                min="1"
                max="7"
                v-model.number.lazy="config.DownloadLimit"
                value="5"
                class="range range-xs bg-primary [--range-shdw:#788091]"
                step="1" />
            <span class="pl-2 text-nowrap">{{ config.DownloadLimit }}</span>
        </label>
        <label class="flex items-center gap-2 input input-bordered">
            代理
            <input
                type="text"
                class="ml-2 grow"
                v-model.lazy="config.ProxyURL"
                placeholder="请输入代理链接" />

            <input type="checkbox" v-model="config.UseProxy" class="checkbox" />
        </label>
        <label class="flex items-center input input-bordered">
            路径
            <input
                type="text"
                class="ml-2 truncate grow"
                v-model.lazy="config.DownloadDir"
                placeholder="设置下载文件夹" />
            <button class="btn btn-square btn-sm" @click="openDownloadDir">
                <span class="icon-[lucide--folder-search]"></span>
            </button>
        </label>
        <label class="flex items-center input input-bordered">
            FFmpeg
            <input
                type="text"
                class="ml-4 truncate grow"
                v-model.lazy="config.FFMPEG"
                @change="checkFFmpeg"
                placeholder="设置ffmpeg文件路径" />
            <button class="btn btn-sm btn-square" @click="ffmpegChange">
                <span class="icon-[lucide--folder-search]"></span>
            </button>
        </label>
    </TabCard>
</template>

<script setup lang="ts">
defineProps<{ tab: Tab }>()
const { config } = storeToRefs(useBasicStore())

function openDownloadDir() {
    SetDownloadDir('请选择文件夹').then((result) => {
        if (result != '') {
            config.value.DownloadDir = result
        } else {
            Message({ message: '用户取消', type: 'warn' })
        }
    })
}
function ffmpegChange() {
    SetFFmpegPath('请选择FFmpeg文件夹').then((result) => {
        if (result != '') {
            config.value.FFMPEG = result
            Message({ message: '设置成功', type: 'success' })
        } else {
            Message({ message: '用户取消/路径无效', type: 'warn' })
        }
    })
}
function checkFFmpeg() {
    CheckFFmpeg(config.value.FFMPEG).then((result) => {
        if (result) {
            Message({ message: '设置成功', type: 'success' })
        } else {
            Message({ message: '用户取消/路径无效', type: 'warn' })
        }
    })
}

watch(config.value, async () => {
    SaveConfig(config.value).then((result) => {
        console.log(result)
    })
})
</script>
