<template>
    <div class="flex flex-col space-y-4 p-2 pt-4 flex-1 text-sm select-none">
        <template v-for="tab in tabs" :key="tab.id">
            <component :is="tab.component" v-if="tabId == tab.id" :tab="tab"></component>
        </template>
    </div>
</template>

<script setup lang="ts">
import { Tab } from '@/models/ui'

defineProps<{ tabs: Tab[]; tabId: number }>()
</script>
