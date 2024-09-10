<template>
  <div class="w-full h-full relative">
    <div class="flex flex-col">
      <!-- 刷新 -->
      <button class="btn absolute right-4 bottom-4">
        <span class="icon-[lucide--refresh-ccw] size-5" @click="refreshPlugins"></span>
      </button>
      <div v-for="(plugin, pluginKey) in plugins" :key="pluginKey" class="group">
        <div class="h-full">
          <div
            class="card w-full h-full p-4 pr-6 space-y-6 overflow-x-hidden overflow-y-auto select-none"
            :class="{ hovered: hoveredIndex === pluginKey }"
            @mouseover="hoveredIndex = pluginKey"
            @mouseout="hoveredIndex = pluginKey">
            <div
              class="p-4 bg-base-100 rounded-xl overflow-hidden flex flex-col hover:shadow-base-content/20 hover:shadow-xl border border-base-300">
              <!-- 状态栏 -->
              <div class="flex items-center pb-4 pl-2">
                <span
                  class="ml-2 text-lg align-middle font-bold group-hover:text-[var(--hover-color)]"
                  :style="{ '--hover-color': plugin.color }">
                  {{ plugin.name }}
                </span>

                <span
                  class="ml-auto"
                  :class="{
                    'opacity-80 filter grayscale': !plugin.enable,
                    'opacity-100': plugin.enable
                  }">
                  <template v-if="plugin.state == 1">
                    <span
                      class="size-6 text-success icon-[ic--outline-check-circle-outline]"></span>
                  </template>
                  <template v-else-if="plugin.state == 2">
                    <span class="size-6 text-warning icon-[lucide--plug-zap]"></span>
                  </template>
                  <template v-else>
                    <span><span class="size-6 text-error icon-[lucide--unplug]"></span></span>
                  </template>
                </span>
              </div>

              <!-- 插件描述 -->
              <div class="m-2 mb-4 space-y-2">{{ plugin.description }}</div>

              <!-- 设置 -->
              <div class="my-2" v-for="(value, key) in plugin.settings" :key="key">
                <label class="flex items-center gap-2 input input-bordered">
                  {{ key }}
                  <input
                    type="text"
                    class="ml-2 grow"
                    v-model.lazy="plugin.settings[key]"
                    :disabled="plugin.lock || !plugin.enable"
                    @change="savePlugin(plugin)" />
                </label>
              </div>

              <div class="divider my-0"></div>

              <!-- 底部命令工具 -->
              <div class="flex items-center">
                <template v-if="plugin.enable">
                  <template v-if="plugin.state !== 1">
                    <span
                      class="size-6 icon-[ic--round-play-arrow] hover:text-success"
                      :disabled="plugin.lock"
                      @click="runPlugin(plugin)"></span>
                  </template>
                  <template v-if="plugin.state === 1">
                    <span
                      class="size-6 icon-[ic--baseline-stop] hover:text-error"
                      :disabled="plugin.lock"
                      @click="stopPlugin(plugin)"></span>
                  </template>
                </template>

                <template v-if="plugin.enable == true">
                  <span
                    class="ml-2 icon-[fluent--presence-blocked-12-regular] hover:warning"
                    :disabled="plugin.lock"
                    @click="disenablePlugin(plugin)"></span>
                </template>
                <template v-else>
                  <span
                    class="ml-1 size-5 icon-[lucide--plug-2] hover:text-success"
                    :disabled="plugin.lock"
                    @click="enablePlugin(plugin)"></span>
                </template>
                <span class="ml-auto space-x-2">
                  <span
                    class="size-6 icon-[ic--round-home]"
                    @click="BrowserOpenURL(plugin.homepage)"></span>
                  <span
                    class="size-6 icon-[iconoir--book-solid]"
                    @click="BrowserOpenURL(plugin.docs_url)"></span>
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { BrowserOpenURL } from '@wailsjs/runtime/runtime'

import { Plugin } from '@/models/go'

import { GetPlugins } from '@wailsjs/go/app/App'
const { plugins } = storeToRefs(useBasicStore())

const hoveredIndex = ref('')
defineProps<{ tab: Tab }>()

EventsOn('updateInfo', (plugin: Plugin) => {
  for (let key in plugins.value) {
    if (plugins.value[key].id == plugin.id) {
      Object.assign(plugins.value[key], plugin)
    }
  }
})

async function savePlugin(plugin: Plugin) {
  const fetchedPlugin = await SavePluginConfig(plugin.id, plugin)

  console.log(fetchedPlugin)

  if (fetchedPlugin) {
    Object.assign(plugin, fetchedPlugin)
  }
}

async function runPlugin(plugin: Plugin) {
  await updatePluginStatus(plugin, RunPlugin)
}
async function stopPlugin(plugin: Plugin): Promise<void> {
  await updatePluginStatus(plugin, StopPlugin)
}
async function enablePlugin(plugin: Plugin): Promise<void> {
  await updatePluginStatus(plugin, EnablePlugin)
}
async function disenablePlugin(plugin: Plugin): Promise<void> {
  await updatePluginStatus(plugin, DisablePlugin)
}

async function updatePluginStatus<T>(
  plugin: Plugin,
  action: (plugin: Plugin) => Promise<T>
): Promise<void> {
  plugin.lock = true
  try {
    const fetchedPlugin = await action(plugin)
    if (fetchedPlugin) {
      Object.assign(plugin, fetchedPlugin)
    }
  } finally {
    plugin.lock = false
  }
}

async function refreshPlugins() {
  const fetchedPlugins = await GetPlugins()
  if (fetchedPlugins) {
    console.log(plugins.value)
    Object.assign(plugins.value, fetchedPlugins)
    console.log('加载插件')
    console.log(plugins.value)
  }
}

onMounted(async () => {
  // 加载插件
  await refreshPlugins()
})
</script>
