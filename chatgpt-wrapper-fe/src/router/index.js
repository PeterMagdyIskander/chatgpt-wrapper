import ChatView from '@/views/ChatView.vue'
import { createRouter, createWebHistory  } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Chat',
    component: ChatView
  },
  {
    path: '/chat/:id?',
    name: 'ChatWithId',
    component: ChatView,
    props: true
  }
]

const router = createRouter({
  history: createWebHistory (),
  routes
})

export default router
