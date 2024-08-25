import { defineStore } from 'pinia'
// import { Part } from '@/models/go'

// import { Config, SystemConfig } from '@/models/go'

import { Plugin } from '@/models/go'

import { app } from '@wailsjs/go/models'

export const useBasicStore = defineStore('basic', () => {
  const configInstance = new app.Config({
    system: new app.SystemConfig(),
    plugins: []
  })

  const plugins: Record<string, Plugin> = reactive({})

  const config = reactive(configInstance)

  const tasks = reactive([])

  return { config, tasks, plugins }
})
