import { createWebHashHistory, createRouter } from 'vue-router'

import HomeView from '@/views/HomeView.vue'
import TaskView from '@/views/TaskView.vue'
import PluginView from '@/views/PluginView.vue'
import SettingView from '@/views/SettingView.vue'
import InfoView from '@/views/InfoView.vue'

const routes = [
  { path: '/', component: HomeView, name: 'home' },
  { path: '/plugins', component: PluginView, name: 'plugins' },
  { path: '/tasks/', component: TaskView, name: 'tasks' },
  { path: '/setting', component: SettingView, name: 'setting' },
  { path: '/info', component: InfoView, name: 'info' }
]

export default createRouter({
  history: createWebHashHistory(),
  routes
})
