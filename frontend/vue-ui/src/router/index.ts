import { createRouter, createWebHistory } from 'vue-router'
import { ref } from 'vue'
import { useSessionStore } from '../stores/session'

export const navigationPending = ref(true)
export const navigationError = ref(false)

export const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', redirect: '/workspace' },
    { path: '/login', name: 'login', component: () => import('../pages/LoginPage.vue'), meta: { title: 'Sign in', guest: true } },
    { path: '/register', name: 'register', component: () => import('../pages/RegisterPage.vue'), meta: { title: 'Create account', guest: true } },
    { path: '/workspace', component: () => import('../pages/WorkspacePage.vue'), meta: { title: 'Workspace', authenticated: true } },
    {
      path: '/admin/users', name: 'users', component: () => import('../pages/admin/UsersPage.vue'), meta: { title: 'Users', admin: true },
      children: [
        { path: 'new', name: 'user-create', component: () => import('../components/users/UserEditorModal.vue'), meta: { title: 'Create user', modal: true } },
        { path: ':id/edit', name: 'user-edit', component: () => import('../components/users/UserEditorModal.vue'), props: { editing: true }, meta: { title: 'Edit user', modal: true } },
        { path: ':id', name: 'user-detail', component: () => import('../components/users/UserDetailModal.vue'), meta: { title: 'User details', modal: true } },
      ],
    },
    { path: '/:pathMatch(.*)*', component: () => import('../pages/NotFoundPage.vue'), meta: { title: 'Page not found' } },
  ],
  scrollBehavior: (to, from, savedPosition) => to.matched[0]?.path === from.matched[0]?.path ? false : savedPosition || { top: 0 },
})

router.beforeEach(async (to) => {
  navigationPending.value = true
  navigationError.value = false
  const session = useSessionStore()
  await session.restore()
  if ((to.meta.authenticated || to.meta.admin) && !session.user) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.meta.admin && !session.isAdmin) return '/workspace'
  if (to.meta.guest && session.user) return session.isAdmin ? '/admin/users' : '/workspace'
})

router.afterEach((to) => {
  navigationPending.value = false
  document.title = `${to.meta.title || 'Workspace'} · Ent Go Vue`
})

router.onError(() => {
  navigationPending.value = false
  navigationError.value = true
})
