<template>
    <div class="flex w-12 flex-col items-center h-full gap-6 justify-between">
        <div class="pt-1 w-full text-center">
            <router-link :to="{ name: 'home' }">
                <div class="avatar pt-2">
                    <IconLogoLight class="text-primary size-8"></IconLogoLight>
                </div>
            </router-link>
        </div>

        <div class="space-y-4 flex flex-col text-center w-full items-center flex-1">
            <router-link
                class="shadow-md w-full rounded-e-md cursor-pointer"
                :to="{ name: 'home' }"
                @click="setActiveLink('home')">
                <span
                    :class="activeLink === 'home' ? 'icon-[ic--home]' : 'icon-[ic--outline-home]'"
                    class="size-[1.75rem]"></span>
            </router-link>
            <router-link
                class="shadow-md w-full rounded-e-md cursor-pointer"
                :to="{ name: 'task' }"
                @click="setActiveLink('task')">
                <span
                    :class="
                        activeLink === 'task'
                            ? 'icon-[ic--round-download]'
                            : 'icon-[ic--outline-download]'
                    "
                    class="size-[1.75rem]"></span>
            </router-link>
            <router-link
                class="shadow-md w-full rounded-e-md cursor-pointer"
                :to="{ name: 'setting' }"
                @click="setActiveLink('setting')">
                <span
                    :class="
                        activeLink === 'setting'
                            ? 'icon-[ic--round-settings]'
                            : 'icon-[ic--outline-settings]'
                    "
                    class="size-6"></span>
            </router-link>
            <router-link
                class="shadow-md w-full rounded-e-md cursor-pointer"
                :to="{ name: 'plugins' }"
                @click="setActiveLink('plugins')">
                <span
                    :class="
                        activeLink === 'plugins'
                            ? 'icon-[ic--round-extension]'
                            : 'icon-[ic--outline-extension]'
                    "
                    class="size-6"></span>
            </router-link>
        </div>

        <div class="flex mt-auto flex-col space-y-2 pb-4">
            <router-link class=" " :to="{ name: 'info' }" @click="setActiveLink('about')">
                <span
                    :class="
                        activeLink === 'about' ? 'icon-[ic--round-info]' : 'icon-[lucide--info]'
                    "
                    class="size-5"></span>
            </router-link>
            <span class="icon-[lucide--power] size-5" @click="showClose()"></span>
        </div>
    </div>
</template>

<script setup lang="ts">
const activeLink = ref<string>('home')

function setActiveLink(link: string) {
    activeLink.value = link
}
function showClose() {
    VAlert({ alert: '确认要退出么' }).then((confirmed) => {
        if (confirmed) {
            Quit()
        }
    })
}
</script>
