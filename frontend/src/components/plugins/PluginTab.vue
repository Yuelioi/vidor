<template>
    <TabCard :tab="tab">
        <label class="flex items-center input-bordered">
            <p>插件信息</p>
            <p>插件版本</p>
            <p>插件用法</p>
            <p>插件文档</p>
        </label>
    </TabCard>
</template>

<script setup lang="ts">
defineProps<{ tab: Tab }>()
</script>
