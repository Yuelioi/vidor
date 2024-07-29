<template>
    <TabCard :tab="tab">
        <span>当前版本 v0.0.1</span>
    </TabCard>
</template>

<script setup lang="ts">
defineProps<{ tab: Tab }>()
</script>
