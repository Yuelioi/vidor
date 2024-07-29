<template>
    <div class="flex items-center space-x-4">
        <div class="h-full flex-1" @dblclick="WindowToggleMaximise" style="widows: 1"></div>

        <div class="flex ml-auto items-center pr-2 space-x-2">
            <span class="icon-[ic--baseline-minus] size-6" @click="WindowMinimise"></span>
            <span class="icon-[lucide--square]" @click="WindowToggleMaximise"></span>
            <span class="icon-[ic--round-close] size-6" @click="WindowHide"></span>
        </div>
    </div>
</template>

<script setup lang="ts">
import { WindowToggleMaximise, WindowHide, WindowMinimise } from '@wailsjs/runtime'
</script>
