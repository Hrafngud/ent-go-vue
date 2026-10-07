<script setup lang="ts">
import { ref } from 'vue'

defineProps<{ busy: boolean; error: string }>()
const emit = defineEmits<{ submit: [email: string, password: string] }>()
const email = ref('')
const password = ref('')
const showPassword = ref(false)

function submit() {
  emit('submit', email.value.trim(), password.value)
}
</script>

<template>
  <div class="flex min-h-screen flex-col">
    <header class="px-6 py-8 sm:px-12">
      <p class="text-sm font-semibold tracking-widest">ENT GO VUE<span class="text-primary">.</span></p>
    </header>

    <section class="flex flex-1 items-center px-6 pb-16 sm:px-12" aria-labelledby="login-title">
      <div class="mx-auto w-full max-w-sm">
        <p class="mb-3 text-xs font-medium uppercase tracking-widest text-base-content/60">Your workspace</p>
        <h1 id="login-title" class="text-4xl font-semibold tracking-tight">Welcome back.</h1>
        <p class="mt-3 text-base-content/70">Sign in to continue.</p>

        <form class="mt-10 space-y-6" :aria-busy="busy" @submit.prevent="submit">
          <div class="space-y-2">
            <label for="login-email" class="block text-sm font-medium">Email address</label>
            <input
              id="login-email"
              v-model="email"
              class="input w-full"
              name="email"
              type="email"
              autocomplete="username"
              placeholder="you@example.com"
              required
              :disabled="busy"
              :aria-describedby="error ? 'login-error' : undefined"
            />
          </div>
          <div class="space-y-2">
            <label for="login-password" class="block text-sm font-medium">Password</label>
            <div class="flex gap-2">
              <input
                id="login-password"
                v-model="password"
                class="input min-w-0 flex-1"
                name="password"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="current-password"
                required
                :disabled="busy"
                :aria-describedby="error ? 'login-error' : undefined"
              />
              <button
                type="button"
                class="btn btn-ghost"
                aria-controls="login-password"
                :aria-pressed="showPassword"
                :disabled="busy"
                @click="showPassword = !showPassword"
              >{{ showPassword ? 'Hide' : 'Show' }}</button>
            </div>
          </div>

          <p v-if="error" id="login-error" class="text-sm text-error" role="alert">{{ error }}</p>
          <button type="submit" class="btn btn-primary w-full" :disabled="busy">
            <span v-if="busy" class="loading loading-spinner loading-sm" aria-hidden="true"></span>
            {{ busy ? 'Signing in…' : 'Sign in' }}
          </button>
        </form>
      </div>
    </section>
    <footer class="px-6 py-6 text-xs text-base-content/60 sm:px-12">Ent Go Vue</footer>
  </div>
</template>
