<template>
  <TabCard :tab="tab">
    <label class="flex items-center input-bordered input">
      样式
      <div class="ml-auto h-full items-center flex dropdown dropdown-bottom dropdown-end">
        <div
          tabindex="0"
          role="button"
          class="btn btn-sm btn-outline"
          @click="showThemeOption = true">
          {{ configs.theme }}
        </div>
        <ul
          tabindex="0"
          v-if="showThemeOption"
          class="dropdown-content menu bg-base-100 rounded-box z-[1] w-52 p-2 shadow">
          <li
            :value="_theme"
            @click="changeTheme(_theme)"
            v-for="_theme in themes"
            :key="_theme"
            class="">
            <span class="px-2">{{ _theme }}</span>
          </li>
        </ul>
      </div>
    </label>
    <label class="flex items-center gap-2 input input-bordered">
      <span class="pr-2 label-text text-nowrap">缩放</span>
      <input
        type="range"
        min="12"
        max="24"
        v-model.number.lazy="configs.scale_factor"
        value="16"
        class="range range-xs [--range-shdw:#788091]"
        @change="changeScaleFactor"
        step="1" />
      <span class="pl-2 text-nowrap">{{ configs.scale_factor }}</span>
    </label>
  </TabCard>
</template>

<script setup lang="ts">
defineProps<{ tab: Tab }>()
const { configs } = storeToRefs(useBasicStore())

const _themes = ['light', 'dark']
const { themes, switchTheme } = useTheme(_themes)
const showThemeOption = ref(false)

function changeScaleFactor() {
  document.documentElement.style.fontSize = `${configs.value.scale_factor}px`
}

function changeTheme(theme) {
  configs.value.theme = theme
  showThemeOption.value = false
}

watch(configs.value, async () => {
  switchTheme(configs.value.theme)
  SaveConfig(configs.value).then((result) => {
    console.log(result)
  })
})
</script>
