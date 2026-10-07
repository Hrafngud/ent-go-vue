import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { ApiError, get, post } from '../api/client'
import type { User } from '../api/users'

const sessionKey = 'ent-go-vue.session'

export const useSessionStore = defineStore('session', () => {
  const user = ref<User | null>(null)
  const token = ref('')
  const notice = ref('')
  const isAdmin = computed(() => user.value?.is_admin === true)
  let restoration: Promise<void> | undefined

  function saveToken(value: string) {
    token.value = value
    try {
      if (value) sessionStorage.setItem(sessionKey, value)
      else sessionStorage.removeItem(sessionKey)
    } catch { /* A session can still work in memory when storage is disabled. */ }
  }

  async function refreshProfile() {
    const profile = await get<{ data: User }>('/users/me', undefined, token.value)
    user.value = profile.data
  }

  function signOut(message = '') {
    saveToken('')
    user.value = null
    notice.value = message
  }

  function restore() {
    restoration ??= (async () => {
      try {
        token.value = sessionStorage.getItem(sessionKey) || ''
      } catch { /* Storage may be disabled. */ }
      if (!token.value) return
      try {
        await refreshProfile()
      } catch (error) {
        signOut(error instanceof ApiError && (error.status === 401 || error.status === 404)
          ? 'Your session has ended. Please sign in again.'
          : 'Could not restore your session. Please sign in again.')
      }
    })()
    return restoration
  }

  async function signIn(email: string, password: string) {
    const response = await post<{ token: string }>('/auth/login', { email, password })
    saveToken(response.token)
    try {
      await refreshProfile()
      notice.value = ''
    } catch (error) {
      signOut()
      throw error
    }
  }

  return { user, token, notice, isAdmin, restore, refreshProfile, signIn, signOut }
})
