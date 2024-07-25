<template>
    <div class="p-4 h-full">
        <div class="w-full h-full p-4 space-y-6 overflow-x-hidden overflow-y-auto select-none">
            <div class="p-4 bg-base-100 rounded-xl hover:shadow-base-content/20 hover:shadow-xl">
                <div class="flex items-center pb-2 pl-2 text-primary h-full">
                    <span class="icon-[ic--outline-color-lens] size-6"></span>
                    <span class="ml-2 font-bold align-middle">主题</span>
                </div>
                <div class="mt-2 mb-4 space-y-6">
                    <label class="flex items-center input-bordered input">
                        样式
                        <div
                            class="ml-auto h-full items-center flex dropdown dropdown-bottom dropdown-end">
                            <div
                                tabindex="0"
                                role="button"
                                class="btn btn-sm btn-outline"
                                @click="showThemeOption = true">
                                {{ config.Theme }}
                            </div>
                            <ul
                                tabindex="0"
                                v-if="showThemeOption"
                                class="dropdown-content menu bg-base-100 rounded-box z-[1] w-52 p-2 shadow">
                                <li
                                    :value="_theme"
                                    @click="changeTheme(_theme)"
                                    v-for="_theme in themes"
                                    :key="_theme"
                                    class="">
                                    <span class="px-2">{{ _theme }}</span>
                                </li>
                            </ul>
                        </div>
                    </label>
                    <label class="flex items-center gap-2 input input-bordered">
                        <span class="pr-2 label-text text-nowrap">缩放</span>
                        <input
                            type="range"
                            min="12"
                            max="24"
                            v-model.number.lazy="config.ScaleFactor"
                            value="16"
                            class="range range-xs [--range-shdw:#788091]"
                            @change="changeScaleFactor"
                            step="1" />
                        <span class="pl-2 text-nowrap">{{ config.ScaleFactor }}</span>
                    </label>
                </div>
            </div>
            <div class="p-4 bg-base-100 rounded-xl hover:shadow-base-content/20 hover:shadow-xl">
                <div class="flex items-center pb-2 pl-2 text-base text-secondary">
                    <span class="icon-[lucide--download] size-6"></span>
                    <span class="ml-2 font-bold">下载</span>
                </div>

                <div class="mt-2 mb-4 space-y-6">
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
                </div>
            </div>

            <div class="p-4 bg-base-100 rounded-xl hover:shadow-base-content/20 hover:shadow-xl">
                <div class="flex pb-2 pl-2 items-center text-accent">
                    <span class="icon-[lucide--send] size-6"></span>
                    <h2 class="pl-2 text-base font-bold">三方数据</h2>
                </div>

                <div class="mt-2 mb-4 space-y-6">
                    <label class="flex items-center gap-2 input input-bordered">
                        <span class="icon-[tabler--brand-bilibili]"></span>
                        SESSDATA
                        <input
                            type="text"
                            class="mx-2 flex-1"
                            v-model.lazy="config.SESSDATA"
                            placeholder="B站登录信息" />
                        <span class="ml-auto icon-[lucide--scan-line] size-5 text-secondary"></span>
                        <span
                            class="cursor-pointer link-secondary"
                            @click="BrowserOpenURL('https://www.bilibili.com/read/cv25451423/')">
                            文档
                        </span>
                    </label>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { BrowserOpenURL } from '@wailsjs/runtime/runtime'

const { config } = storeToRefs(useBasicStore())

const _themes = ['light', 'dark']
const { themes, switchTheme } = useTheme(_themes)
const showThemeOption = ref(false)

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

function changeScaleFactor() {
    document.documentElement.style.fontSize = `${config.value.ScaleFactor}px`
}

function changeTheme(theme) {
    config.value.Theme = theme
    showThemeOption.value = false
}

watch(config.value, async () => {
    switchTheme(config.value.Theme)
    SaveConfig(config.value).then((result) => {
        console.log(result)
    })
})
</script>
