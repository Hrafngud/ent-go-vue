<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ApiError } from '../api/client'
import { useSessionStore } from '../stores/session'
import AuthLayout from '../components/layout/AuthLayout.vue'
import PasswordField from '../components/forms/PasswordField.vue'
import FeedbackAlert from '../components/ui/FeedbackAlert.vue'
import { emailError, passwordError, normalizeLogin } from '../validation/account'
import { userError } from '../api/errors'
import { useCooldown } from '../composables/useCooldown'

const session = useSessionStore()
const route = useRoute()
const router = useRouter()
const email = ref(typeof route.query.email === 'string' ? route.query.email : '')
const password = ref('')
const busy = ref(false)
const error = ref('')
const cooldown = useCooldown()
const cooldownSeconds = cooldown.seconds

async function signIn() {
  if (busy.value || cooldownSeconds.value) return
  error.value = emailError(email.value) || passwordError(password.value, 1)
  if (error.value) return
  busy.value = true
  error.value = ''
  try {
    await session.signIn(normalizeLogin(email.value), password.value)
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : ''
    const safeRedirect = redirect.startsWith('/') && !redirect.startsWith('//') && !['/login', '/register'].includes(redirect.split('?')[0] || '')
    await router.replace(safeRedirect ? redirect : session.isAdmin ? '/admin/users' : '/workspace')
  } catch (err) {
    cooldown.start(err)
    error.value = err instanceof ApiError && (err.status === 401 || err.status === 422)
      ? 'The email or password is incorrect. Please try again.'
      : userError(err, 'Could not sign in. Please try again in a moment.')
  } finally { busy.value = false }
}
</script>

<template>
  <AuthLayout title="Welcome back." description="Your workspace is ready." eyebrow="Your account" split>
    <FeedbackAlert v-if="route.query.registered === '1'" class="mb-6" kind="success" message="Your account is ready. Sign in to get started." />
    <FeedbackAlert v-if="session.notice" class="mb-6" kind="info" :message="session.notice" />
    <form class="space-y-6" :aria-busy="busy" @submit.prevent="signIn">
      <div class="space-y-2">
        <label for="login-email" class="block text-sm font-medium">Email address</label>
        <input id="login-email" v-model="email" class="input w-full" name="email" type="email" autocomplete="username"
          placeholder="you@example.com" required maxlength="254" :disabled="busy" />
      </div>
      <PasswordField id="login-password" v-model="password" autocomplete="current-password" required :disabled="busy" />
      <FeedbackAlert v-if="error" :message="error" />
      <button type="submit" class="btn btn-primary w-full" :disabled="busy || cooldownSeconds > 0">
        <span v-if="busy" class="loading loading-spinner loading-sm" aria-hidden="true"></span>
        {{ busy ? 'Signing in…' : cooldownSeconds ? `Try again in ${cooldownSeconds}s` : 'Sign in' }}
      </button>
      <span v-if="busy" class="sr-only" role="status">Signing in…</span>
    </form>
    <p class="mt-8 border-t border-base-300 pt-6 text-center text-sm text-base-content/65">New here? <RouterLink class="link link-primary link-hover" to="/register">Create an account</RouterLink></p>
  </AuthLayout>
</template>
