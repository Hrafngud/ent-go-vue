import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { QueryClient, VueQueryPlugin } from '@tanstack/vue-query'
import { router } from './router'
import { useSessionStore } from './stores/session'
import { useFeedbackStore } from './stores/feedback'

import App from './App.vue'
import './style.css'

const app = createApp(App)

app.use(createPinia())
const queryClient = new QueryClient()
app.use(VueQueryPlugin, { queryClient })
app.use(router)

window.addEventListener('session-expired', () => {
  const session = useSessionStore()
  if (!session.token) return
  session.signOut('Your session has ended. Please sign in again.')
  queryClient.clear()
  useFeedbackStore().clear()
  void router.replace({ name: 'login', query: { redirect: router.currentRoute.value.fullPath } })
})

app.mount('#app')
