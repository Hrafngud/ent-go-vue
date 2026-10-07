<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { post } from '../api/client'
import { userError } from '../api/errors'
import type { UserInput } from '../api/users'
import AuthLayout from '../components/layout/AuthLayout.vue'
import UserForm from '../components/users/UserForm.vue'

const router = useRouter()
const busy = ref(false)
const error = ref('')

async function register(input: UserInput) {
  if (busy.value) return
  busy.value = true
  error.value = ''
  try {
    await post<void>('/auth/register', input)
    await router.replace({ name: 'login', query: { registered: '1', email: input.email } })
  } catch (err) {
    error.value = userError(err, 'Could not create your account. Please try again in a moment.')
  } finally { busy.value = false }
}
</script>

<template>
  <AuthLayout title="Make yourself at home." description="Create your account to join the workspace.">
    <UserForm :busy="busy" :error="error" submit-label="Create account" confirm-password @submit="register" />
    <p class="mt-6 text-center text-sm text-base-content/70">Already have an account? <RouterLink class="link link-primary" to="/login">Sign in</RouterLink></p>
  </AuthLayout>
</template>
