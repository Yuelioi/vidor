<template>
    <TabCard :tab="tab">
        <label class="flex items-center input-bordered">晚点发布</label>
    </TabCard>
</template>

<script setup lang="ts">
defineProps<{ tab: Tab }>()
</script>
