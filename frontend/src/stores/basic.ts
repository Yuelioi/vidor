import { defineStore } from 'pinia'
// import { Part } from '@/models/go'

// import { Config, SystemConfig } from '@/models/go'

import { Plugin } from '@/models/go'

import { config, proto } from '@wailsjs/go/models'

export const useBasicStore = defineStore('basic', () => {
  const configInstance = new config.Config()

  const plugins: Record<string, Plugin> = reactive({})

  const configs = reactive(configInstance)

  const tasks = reactive<proto.Task[]>([])

  const task1 = new proto.Task()
  task1.state = 1
  task1.status = '下载出错下载出错下载出错下载出错下载出错下载出错下载出错'
  task1.speed = '20M/S'
  task1.size = 500
  task1.duration = 300
  task1.percent = 50
  task1.id = '123'
  task1.url = 'http://example.com/task1'
  task1.session_id = 'session-123'
  task1.title = 'Task One Task One Task One Task One Task One'
  task1.cover = 'https://cdn.yuelili.com/docs/web/assert/anime-girl.jpg'
  task1.work_dir = '/path/to/workdir'
  task1.segments = [] // 假设Segment是一个空数组

  const task4 = new proto.Task()
  task4.state = 1
  task4.status = '下载视频中'
  task4.speed = '20M/S'
  task4.size = 5000
  task4.duration = 600
  task4.percent = 50
  task4.id = '123'
  task4.url = 'http://example.com/task1'
  task4.session_id = 'session-123'
  task4.title = 'Task One'
  task4.cover = '1'
  task4.work_dir = '/path/to/workdir'
  task4.segments = [] // 假设Segment是一个空数组

  const task2 = new proto.Task()
  task2.state = 2
  task2.speed = '20M/S'
  task1.size = 6000
  task2.duration = 600
  task2.percent = 75
  task2.id = '456'
  task2.url = 'http://example.com/task2'
  task2.session_id = 'session-456'
  task2.title = 'Task Two'
  task2.cover = 'https://cdn.yuelili.com/docs/web/assert/anime-girl.jpg'
  task2.work_dir = '/path/to/workdir2'
  task2.segments = [] // 假设Segment是一个空数组

  const task3 = new proto.Task()
  task3.state = 3
  task1.speed = '20M/S'
  task1.size = 500
  task3.duration = 600
  task3.percent = 25
  task3.id = '789'
  task3.url = 'http://example.com/task3'
  task3.session_id = 'session-789'
  task3.title = 'Task Three'
  task3.cover = 'https://cdn.yuelili.com/docs/web/assert/anime-girl.jpg'
  task3.work_dir = '/path/to/workdir3'
  task3.segments = [] // 假设Segment是一个空数组

  tasks.push(task1)
  tasks.push(task4)
  tasks.push(task2)
  tasks.push(task3)

  return { configs, tasks, plugins }
})
