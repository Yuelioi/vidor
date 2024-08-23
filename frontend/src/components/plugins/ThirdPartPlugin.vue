<template>
  <div class="flex flex-col">
    <TabCard :tab="tab">
      <label class="flex items-center input-bordered">晚点发布</label>
    </TabCard>
    <TabCard :tab="tab">
      <label class="flex items-center input-bordered">晚点发布</label>
    </TabCard>
  </div>
</template>

<script setup lang="ts">
defineProps<{ tab: Tab }>()
</script>
