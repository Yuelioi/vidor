<template>
  <TabCard :tab="tab">
    <p class="flex items-center gap-2">尽情期待</p>
  </TabCard>
</template>

<script setup lang="ts">
defineProps<{ tab: Tab }>()

const { configs } = storeToRefs(useBasicStore())

watch(configs.value, async () => {
  SaveConfig(configs.value).then((result) => {
    console.log(result)
  })
})
</script>
