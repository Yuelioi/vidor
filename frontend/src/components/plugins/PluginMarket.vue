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
        <div class="my-4 w-4/5 h-32 group hover:shadow-2xl" v-for="plugin in filteredMarketPlugins" :key="plugin.id">
          <div class="card w-full h-full card-side bg-base-100 shadow-xl">
            <figure class="basis-3/12 relative">
              <img
                src="https://img.daisyui.com/images/stock/photo-1635805737707-575885ab0820.webp"
                alt="Movie" />
            </figure>

            <div class="card-body py-4 px-6 basis-9/12 relative">
              <!-- 内容第一行 -->
              <div class="flex items-center">
                <span class="text-center text-lg font-bold]">
                  {{ plugin.name }}
                </span>
                <span class="flex  ml-3">
                  <span class="size-4 mb-1 icon-[ic--outline-cloud-download]"></span>
                  <span class="mx-1">100</span>
                </span>

                <span class="ml-2 space-x-2 opacity-0 group-hover:opacity-100">
                  <span
                    class="size-6 icon-[ic--round-home]"
                    @click="BrowserOpenURL(plugin.homepage)"></span>
                  <span
                    class="size-6 icon-[iconoir--book-solid]"
                    @click="BrowserOpenURL(plugin.docs_url)"></span>
                </span>
              </div>

              <!-- 内容第一行 -->
              <p class="w-5/6">
                <div class="text-ellipsis line-clamp-1">
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
                </div>
              </p>

              <div class="flex items-center">
                <div class="">{{ plugin.author }}</div>

                <div v-for="tag,index in plugin.tags" :key="tag">
                  <span class="badge badge-neutral">tag</span>
                </div>
              </div>

              <!-- 右上 版本号 -->
              <div
                class="absolute right-4 top-4 badge badge-neutral badge-sm mr-auto text-neutral-content">
                <span class="text-slate-400" >{{ plugin.version }}</span>
              </div>

              <!-- 右下 下载按钮 -->
              <div class="absolute right-4 bottom-4">
                <button
                  :disabled="calculateLock(plugin)"
                  class="btn btn-sm btn-primary"
                  @click="download(plugin)">
                  <span>{{ calculatePluginState(plugin) }}</span>
                </button>
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

import { GetPlugins } from '@wailsjs/go/app/App'
const { plugins } = storeToRefs(useBasicStore())

const search = ref('')
const marketPlugins = reactive<Plugin[]>([])
const filteredMarketPlugins = computed(() => {
  const tmpPlugins: Plugin[] = []

  marketPlugins.forEach((plugin) => {
    if (
      plugin.name.toLowerCase().includes(search.value.toLowerCase()) ||
      plugin.description.toLowerCase().includes(search.value.toLowerCase()) ||
      plugin.id.toLowerCase().includes(search.value.toLowerCase())
    ) {
      tmpPlugins.push(plugin)
    }
  })
  return tmpPlugins
})

EventsOn('plugin.downloading', (plugin?: Plugin) => {
  console.log(plugin)
})

function calculateLock(plugin:Plugin){
  // 下载中
  if (plugin.lock){
    return true
  }

  if (plugin.id in plugins.value) {
    if (plugin.version < plugins[plugin.id]) {
      return false
    } else {
      return true
    }
  }
  return false
  
}

function compareVersions(v1, v2) {
  const v1Parts = v1.split('.').map(Number);
  const v2Parts = v2.split('.').map(Number);

  for (let i = 0; i < Math.max(v1Parts.length, v2Parts.length   
); i++) {
    if (v1Parts[i] !== v2Parts[i]) {
      return v1Parts[i] - v2Parts[i];
    }
  }

  return 0;
}

// 插件按钮所属状态
function calculatePluginState(plugin: Plugin) {
  // 本地有插件
  if (plugin.id in plugins.value) {
    if (plugin.version < plugins[plugin.id]) {
      return '更新'
    } else {
      return '已下载'
    }
  }
  // 下载中
  if (plugin.lock) {
    return '下载中'
  }
  return '下载'
}

onMounted(async () => {
  const resp = await fetch('https://cdn.yuelili.com/market/vidor/plugins.json')
  const data = await resp.json()
  Object.assign(marketPlugins, data)
  console.log(marketPlugins)
})

// 下载插件
async function download(plugin: Plugin) {
  console.log(plugin)
  plugin.lock = true
  const p = await DownloadPlugin(plugin)

  if (p) {
    plugins.value[plugin.id] = p
  }

  plugin.lock = false
  console.log(plugin)
}

function searchPlugin() {}

onMounted(async () => {
  // 加载插件
  const fetchedPlugins = await GetPlugins()
  if (fetchedPlugins) {
    console.log(plugins.value)
    Object.assign(plugins.value, fetchedPlugins)
    console.log('加载插件')
    console.log(plugins.value)
  }
})
</script>
