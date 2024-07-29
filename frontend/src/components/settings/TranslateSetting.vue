<template>
    <TabCard :tab="tab">
        <p class="flex items-center gap-2">尽情期待</p>
    </TabCard>
</template>

<script setup lang="ts">
defineProps<{ tab: Tab }>()

const { config } = storeToRefs(useBasicStore())

watch(config.value, async () => {
    SaveConfig(config.value).then((result) => {
        console.log(result)
    })
})
</script>
