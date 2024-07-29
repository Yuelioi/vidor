<template>
    <TabCard :tab="tab">
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
    </TabCard>
</template>

<script setup lang="ts">
defineProps<{ tab: Tab }>()
import { BrowserOpenURL } from '@wailsjs/runtime/runtime'

const { config } = storeToRefs(useBasicStore())

watch(config.value, async () => {
    SaveConfig(config.value).then((result) => {
        console.log(result)
    })
})
</script>
