import { defineStore } from 'pinia'
// import { Part } from '@/models/go'

// import { Config, SystemConfig } from '@/models/go'

import { app } from '@wailsjs/go/models'

export const useBasicStore = defineStore('basic', () => {
    const configInstance = new app.Config({
        system: new app.SystemConfig(),
        plugins: []
    })

    const config = reactive(configInstance)

    const tasks = reactive([])

    return { config, tasks }
})
