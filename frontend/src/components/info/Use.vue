<template>
    <div class="p-4">
        <div class="flex flex-col space-y-6">
            <div class="">
                <h2 class="text-md font-bold py-2">下载</h2>
                <p class="text-sm">
                    <span class="text-primary">B站</span>
                    下载高清视频,请在
                    <router-link
                        :to="{
                            name: 'setting'
                        }"
                        class="link link-primary">
                        设置页面 -> 三方设置
                    </router-link>
                    配置SESSDATA, 可以选择直接填写, 也可以扫码(后续会支持)
                </p>
            </div>
            <div class="">
                <h2 class="text-md font-bold py-2">合并</h2>
                <p class="text-sm">
                    请安装
                    <span
                        class="link link-secondary"
                        @click="BrowserOpenURL('https://zhuanlan.zhihu.com/p/118362010')">
                        FFmpeg
                    </span>
                    配置环境变量, 或者在
                    <router-link
                        :to="{
                            name: 'setting'
                        }"
                        class="link link-secondary">
                        设置页面
                    </router-link>
                    手动指定安装路径
                </p>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { BrowserOpenURL } from '@wailsjs/runtime/runtime'
</script>
