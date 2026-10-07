<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useQueryClient } from '@tanstack/vue-query'
import { useSessionStore } from '../../stores/session'
import { useFeedbackStore } from '../../stores/feedback'
import BrandMark from './BrandMark.vue'
import FeedbackAlert from '../ui/FeedbackAlert.vue'

const session = useSessionStore()
const feedback = useFeedbackStore()
const router = useRouter()
const queryClient = useQueryClient()

function signOut() {
  session.signOut()
  queryClient.clear()
  feedback.clear()
  void router.replace('/login')
}
</script>

<template>
  <a href="#main-content" class="btn btn-primary fixed left-4 top-4 z-50 -translate-y-24 focus:translate-y-0">Skip to content</a>
  <header class="border-b border-base-300">
    <div class="navbar mx-auto max-w-6xl flex-wrap gap-4 px-6 py-5">
      <BrandMark />
      <div class="ml-auto flex min-w-0 items-center gap-3">
        <span class="hidden max-w-xs truncate text-sm text-base-content/70 sm:block">{{ session.user?.name }}</span>
        <span v-if="session.isAdmin" class="badge badge-primary badge-outline badge-sm">Root</span>
        <button class="btn btn-ghost btn-sm" @click="signOut">Sign out</button>
      </div>
      <nav class="w-full" aria-label="Main navigation">
        <ul class="menu menu-horizontal gap-1 p-0">
          <li><RouterLink to="/workspace" active-class="menu-active">Workspace</RouterLink></li>
          <li v-if="session.isAdmin"><RouterLink to="/admin/users" active-class="menu-active">Users</RouterLink></li>
        </ul>
      </nav>
    </div>
  </header>
  <main id="main-content" class="mx-auto max-w-6xl px-6 py-10 sm:py-14">
    <FeedbackAlert v-if="feedback.message" class="mb-6" kind="success" :message="feedback.message" dismissible @dismiss="feedback.clear" />
    <slot />
  </main>
</template>
