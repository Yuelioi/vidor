<template>
    <Teleport to="body">
        <transition name="fade">
            <div
                v-if="show"
                class="fixed top-0 left-0 z-40 w-screen h-screen bg-base-300 opacity-50 overlay base-300"></div>
        </transition>
        <transition name="bounce">
            <div v-if="show" class="fixed inset-0 flex items-center justify-center z-50">
                <div
                    ref="dialogRef"
                    class="w-[80%] overflow-hidden min-w-[375px] min-h-[500px] max-h-[75vh] flex flex-col rounded-lg bg-base-200 overflow-y-auto">
                    <header
                        class="flex px-4 py-2 items-center text-sm text-base-100 justify-between bg-base-300 font-bold">
                        <slot name="header">
                            <span>{{ title }}</span>
                        </slot>

                        <span
                            class="icon-[ic--round-close] size-6 mx-1"
                            @click="show = false"></span>
                    </header>
                    <main
                        class="text-base-content bg-base-100 shadow-lg rounded-lg overflow-x-hidden overflow-y-scroll">
                        <slot></slot>
                    </main>
                    <footer>
                        <slot name="footer"></slot>
                    </footer>
                </div>
            </div>
        </transition>
    </Teleport>
</template>

<script lang="ts" setup>
defineProps({
    title: {
        type: String,
        required: false
    }
})

const show = defineModel('show', { default: false, required: true })
const dialogRef = ref<HTMLElement | null>(null)

const closeDialog = () => {
    show.value = false
}

const handleClickOutside = (e: MouseEvent) => {
    if (dialogRef.value && !dialogRef.value.contains(e.target as Node)) {
        closeDialog()
    }
}

onMounted(() => {
    document.body.classList.add('dialog')
    document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
    document.body.classList.remove('dialog')
    document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.5s;
}

.fade-enter,
.fade-leave-to {
    opacity: 0;
}

.bounce-enter-active {
    animation: bounce-in 0.5s;
}

.bounce-leave-active {
    animation: bounce-in 0.1s reverse;
}

@keyframes bounce-in {
    0% {
        transform: scale(0);
    }
    50% {
        transform: scale(1.1);
    }
    100% {
        transform: scale(1);
    }
}
</style>
