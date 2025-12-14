import { createRouter, createWebHistory } from 'vue-router'
import FileBrowser from '../views/FileBrowser.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: FileBrowser
    }
  ]
})

export default router