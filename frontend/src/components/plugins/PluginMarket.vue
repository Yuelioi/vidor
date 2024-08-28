<template>
  <div class="w-full h-full">
    <div class="flex flex-col items-center">
      <div class="font-bold text-2xl py-4">插件商店</div>

      <div class="py-3">当前共有 {{ marketPlugins.length }} 个插件</div>
      <label class="input input-bordered w-4/5 flex items-center gap-2">
        <input type="text" class="grow" v-model="search" placeholder="搜索" />
        <span class="icon-[lucide--search]"></span>
      </label>

      <!-- 插件列表 -->
      <div class="w-full flex flex-col items-center">
        <div class="my-4 w-4/5 h-32 group" v-for="plugin in marketPlugins" :key="plugin.id">
          <div class="card w-full h-full card-side bg-base-100 shadow-xl">
            <figure class="basis-3/12 relative">
              <img
                src="https://img.daisyui.com/images/stock/photo-1635805737707-575885ab0820.webp"
                alt="Movie" />
            </figure>

            <div class="card-body p-6 basis-9/12 relative">
              <div class="flex">
                <h2 class="text-center text-lg">{{ plugin.name }}</h2>
                <div class="flex items-center ml-3">
                  <span class="size-4 mb-1 icon-[ic--outline-cloud-download]"></span>
                  <span class="mx-1">100</span>
                </div>

                <span class="ml-auto space-x-2 opacity-0 group-hover:opacity-100">
                  <span
                    class="size-6 icon-[ic--round-home]"
                    @click="BrowserOpenURL(plugin.homepage)"></span>
                  <span
                    class="size-6 icon-[iconoir--book-solid]"
                    @click="BrowserOpenURL(plugin.docs_url)"></span>
                </span>
              </div>

              <p class="line-clamp-2 w-5/6">
                {{
                  plugin.description +
                  plugin.description +
                  plugin.description +
                  plugin.description +
                  plugin.description +
                  plugin.description +
                  plugin.description +
                  plugin.description
                }}
              </p>

              <div class="absolute right-4 bottom-4">
                <button class="btn btn-sm btn-primary" @click="download(plugin)">下载</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Plugin } from '@/models/go'

import { BrowserOpenURL } from '@wailsjs/runtime/runtime'

const search = ref('')
const marketPlugins = reactive<Plugin[]>([])

EventsOn('plugin.downloading', (plugin?: Plugin) => {
  console.log(plugin)
})

onMounted(async () => {
  const resp = await fetch('https://cdn.yuelili.com/market/vidor/plugins.json')
  const data = await resp.json()
  Object.assign(marketPlugins, data)
  console.log(data)
  console.log(marketPlugins)
})

async function download(plugin: Plugin) {
  console.log(plugin)
  const p = await DownloadPlugin(plugin)
  console.log(p)
}
</script>
