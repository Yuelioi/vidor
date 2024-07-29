<template>
    <TabCard :tab="tab">
        <div class="flex-1">Vidor 是一款由Golang编写的高性能下载器</div>
        <div class="divider mt-auto"></div>
        <div class="flex flex-col space-y-2 h-full">
            <div class="">
                <div
                    class="link flex items-center space-x-2"
                    @click="BrowserOpenURL('https://space.bilibili.com/4279370')">
                    <span class="icon-[tabler--brand-bilibili] size-5 text-primary"></span>
                    <span>https://space.bilibili.com/4279370</span>
                </div>
            </div>

            <div class="">
                <div
                    class="link flex items-center space-x-2"
                    @click="BrowserOpenURL('https://space.bilibili.com/4279370')">
                    <span class="icon-[lucide--github] size-5"></span>
                    <span>https://github.com/Yuelioi/vidor</span>
                </div>
            </div>
        </div>
    </TabCard>
</template>

<script setup lang="ts">
import { BrowserOpenURL } from '@wailsjs/runtime/runtime'
defineProps<{ tab: Tab }>()
</script>
