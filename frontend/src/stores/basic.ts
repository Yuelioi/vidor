import { defineStore } from 'pinia'
// import { Part } from '@/models/go'
import { Config, SystemConfig } from '@/models/go'

export const useBasicStore = defineStore('basic', () => {
    const config = reactive<Config>({
        system: new SystemConfig(),
        plugins: []
    })

    const tasks = reactive([])

    return { config, tasks }
})
