<template>
  <TabCard :tab="tab">
    <label class="flex items-center gap-2 input input-bordered">
      名称
      <input
        type="text"
        class="ml-2 grow"
        v-model.lazy="configs.magic_name"
        placeholder="下载文件魔法名称" />
    </label>
    <label class="flex items-center gap-2 input input-bordered">
      <span class="pr-2 label-text text-nowrap">并行</span>
      <input
        type="range"
        min="1"
        max="7"
        v-model.number.lazy="configs.download_limit"
        value="5"
        class="range range-xs bg-primary [--range-shdw:#788091]"
        step="1" />
      <span class="pl-2 text-nowrap">{{ configs.download_limit }}</span>
    </label>
    <label class="flex items-center gap-2 input input-bordered">
      代理
      <input
        type="text"
        class="ml-2 grow"
        v-model.lazy="configs.proxy_url"
        placeholder="请输入代理链接" />

      <input type="checkbox" v-model="configs.use_proxy" class="checkbox" />
    </label>
    <label class="flex items-center input input-bordered">
      路径
      <input
        type="text"
        class="ml-2 truncate grow"
        v-model.lazy="configs.download_dir"
        placeholder="设置下载文件夹" />
      <button class="btn btn-square btn-sm" @click="openDownloadDir">
        <span class="icon-[lucide--folder-search]"></span>
      </button>
    </label>
  </TabCard>
</template>

<script setup lang="ts">
import { SetDownloadDir } from '@wailsjs/go/app/App'
defineProps<{ tab: Tab }>()
const { configs } = storeToRefs(useBasicStore())

function openDownloadDir() {
  SetDownloadDir('请选择文件夹').then((result) => {
    if (result != '') {
      configs.value.download_dir = result
    } else {
      Message({ message: '用户取消', type: 'warn' })
    }
  })
}

watch(configs.value, async () => {
  SaveConfig(configs.value).then((result) => {
    console.log(result)
  })
})
</script>
