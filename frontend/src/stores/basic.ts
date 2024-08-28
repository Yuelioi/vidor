import { Task } from '@/models/go'
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

  const tasks = reactive<Task[]>([])

  const task1 = new Task()
  task1.title = '11'
  task1.state = 1

  const task2 = new Task()
  task2.title = '22'
  task2.state = 2

  const task3 = new Task()
  task3.title = '33'
  task3.state = 3

  tasks.push(task1)
  tasks.push(task2)
  tasks.push(task3)

  return { configs, tasks, plugins }
})
