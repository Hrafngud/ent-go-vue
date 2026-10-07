<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useQueryClient } from '@tanstack/vue-query'
import { ApiError, get, post } from './api/client'
import ConnectionStatus from './components/ConnectionStatus.vue'
import LoginPage from './components/LoginPage.vue'

type User = { id: string; name: string; email: string }
const sessionKey = 'ent-go-vue.session'
const user = ref<User | null>(null)
const restoring = ref(true)
const busy = ref(false)
const error = ref('')
const queryClient = useQueryClient()

function saveToken(token: string | null) {
  try {
    if (token) sessionStorage.setItem(sessionKey, token)
    else sessionStorage.removeItem(sessionKey)
  } catch {
    // Sign-in still works in memory when browser storage is unavailable.
  }
}

async function loadProfile(token: string) {
  const profile = await get<{ data: User }>('/users/me', undefined, token)
  user.value = profile.data
}

async function signIn(email: string, password: string) {
  if (busy.value) return
  busy.value = true
  error.value = ''
  try {
    const { token } = await post<{ token: string }>('/auth/login', { email, password })
    await loadProfile(token)
    saveToken(token)
  } catch (err) {
    error.value = err instanceof ApiError && (err.status === 401 || err.status === 422)
      ? 'The email or password is incorrect. Please try again.'
      : 'Could not sign in. Please try again in a moment.'
  } finally {
    busy.value = false
  }
}

function signOut() {
  saveToken(null)
  user.value = null
  error.value = ''
  queryClient.clear()
}

onMounted(async () => {
  let token: string | null = null
  try {
    token = sessionStorage.getItem(sessionKey)
  } catch {
    // Storage may be disabled by the browser.
  }
  try {
    if (token) await loadProfile(token)
  } catch (err) {
    saveToken(null)
    error.value = err instanceof ApiError && (err.status === 401 || err.status === 404)
      ? 'Your session has ended. Please sign in again.'
      : 'Could not restore your session. Please sign in again.'
  } finally {
    restoring.value = false
  }
})
</script>

<template>
  <main class="min-h-screen bg-base-100 text-base-content">
    <div v-if="restoring" class="flex min-h-screen items-center justify-center gap-3" role="status">
      <span class="loading loading-spinner loading-sm" aria-hidden="true"></span>
      Restoring your session…
    </div>
    <template v-else-if="user">
      <header class="mx-auto flex max-w-3xl flex-wrap items-center justify-between gap-4 px-6 pt-8">
        <p class="min-w-0 break-words text-sm">Signed in as <strong>{{ user.name }}</strong></p>
        <button class="btn btn-ghost btn-sm" @click="signOut">Sign out</button>
      </header>
      <ConnectionStatus />
    </template>
    <LoginPage v-else :busy="busy" :error="error" @submit="signIn" />
  </main>
</template>
