import { defineStore } from 'pinia'
// import { Part } from '@/models/go'

// import { Config, SystemConfig } from '@/models/go'

import { Plugin } from '@/models/go'

import { config, models } from '@wailsjs/go/models'

export const useBasicStore = defineStore('basic', () => {
  const configInstance = new config.Config({
    system: new models.SystemConfig(),
    plugins: []
  })

  const plugins: Record<string, Plugin> = reactive({})

  const configs = reactive(configInstance)

  const tasks = reactive([])

  return { configs, tasks, plugins }
})
