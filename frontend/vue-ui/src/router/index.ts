import { createRouter, createWebHistory } from 'vue-router'
import { useSessionStore } from '../stores/session'

export const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', redirect: '/workspace' },
    { path: '/login', name: 'login', component: () => import('../pages/LoginPage.vue'), meta: { title: 'Sign in', guest: true } },
    { path: '/register', name: 'register', component: () => import('../pages/RegisterPage.vue'), meta: { title: 'Create account', guest: true } },
    { path: '/workspace', component: () => import('../pages/WorkspacePage.vue'), meta: { title: 'Workspace', authenticated: true } },
    { path: '/admin/users', name: 'users', component: () => import('../pages/admin/UsersPage.vue'), meta: { title: 'Users', admin: true } },
    { path: '/admin/users/new', name: 'user-create', component: () => import('../pages/admin/UserCreatePage.vue'), meta: { title: 'Create user', admin: true } },
    { path: '/admin/users/:id/edit', name: 'user-edit', component: () => import('../pages/admin/UserEditPage.vue'), meta: { title: 'Edit user', admin: true } },
    { path: '/admin/users/:id', name: 'user-detail', component: () => import('../pages/admin/UserDetailPage.vue'), meta: { title: 'User details', admin: true } },
    { path: '/:pathMatch(.*)*', component: () => import('../pages/NotFoundPage.vue'), meta: { title: 'Page not found' } },
  ],
  scrollBehavior: () => ({ top: 0 }),
})

router.beforeEach(async (to) => {
  const session = useSessionStore()
  await session.restore()
  if ((to.meta.authenticated || to.meta.admin) && !session.user) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.meta.admin && !session.isAdmin) return '/workspace'
  if (to.meta.guest && session.user) return session.isAdmin ? '/admin/users' : '/workspace'
})

router.afterEach((to) => {
  document.title = `${to.meta.title || 'Workspace'} · Ent Go Vue`
})
