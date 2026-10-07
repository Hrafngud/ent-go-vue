<script setup lang="ts">
import { PhSignOut, PhSquaresFour, PhUsers } from '@phosphor-icons/vue'
import { useRoute, useRouter } from 'vue-router'
import { useQueryClient } from '@tanstack/vue-query'
import { useSessionStore } from '../../stores/session'
import { useFeedbackStore } from '../../stores/feedback'
import BrandMark from './BrandMark.vue'
import FeedbackAlert from '../ui/FeedbackAlert.vue'
import UserAvatar from '../users/UserAvatar.vue'

const session = useSessionStore()
const feedback = useFeedbackStore()
const router = useRouter()
const route = useRoute()
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
  <div class="app-shell">
    <header class="border-b border-base-300 bg-base-200 lg:sticky lg:top-0 lg:h-dvh lg:border-r lg:border-b-0">
      <div class="flex h-full flex-col px-4 py-4 lg:px-5 lg:py-8">
        <div class="flex items-center justify-between gap-3 lg:px-2">
          <BrandMark />
          <button class="btn btn-ghost btn-sm lg:hidden" @click="signOut"><PhSignOut :size="20" aria-hidden="true" />Sign out</button>
        </div>
        <nav class="app-nav mt-4 lg:mt-12" aria-label="Main navigation">
          <p class="mb-3 hidden px-3 text-xs font-medium tracking-widest text-base-content/55 uppercase lg:block">Workspace</p>
          <ul class="menu menu-horizontal w-full gap-1 p-0 text-sm lg:menu-vertical">
            <li class="flex-1 lg:flex-none"><RouterLink to="/workspace" active-class="menu-active"><PhSquaresFour :size="20" aria-hidden="true" />Overview</RouterLink></li>
            <li v-if="session.isAdmin" class="flex-1 lg:flex-none"><RouterLink to="/admin/users" :class="{ 'menu-active': route.path.startsWith('/admin/users') }" :aria-current="route.path.startsWith('/admin/users') ? 'location' : undefined"><PhUsers :size="20" aria-hidden="true" />Users</RouterLink></li>
          </ul>
        </nav>
        <div class="mt-auto hidden border-t border-base-300 pt-5 lg:block">
          <div class="flex min-w-0 items-center gap-3 px-2">
            <UserAvatar :name="session.user?.name || ''" />
            <div class="min-w-0">
              <p class="truncate text-sm font-medium">{{ session.user?.name }}</p>
              <p class="mt-0.5 text-xs text-base-content/60">{{ session.isAdmin ? 'Root administrator' : 'Member' }}</p>
            </div>
          </div>
          <button class="btn btn-ghost btn-sm mt-4 w-full justify-start text-base-content/70" @click="signOut"><PhSignOut :size="20" aria-hidden="true" />Sign out</button>
        </div>
      </div>
    </header>
    <main id="main-content" class="mx-auto w-full min-w-0 max-w-6xl px-5 py-8 sm:px-8 lg:px-12 lg:py-12">
      <FeedbackAlert v-if="feedback.message" class="mb-6" kind="success" :message="feedback.message" dismissible @dismiss="feedback.clear" />
      <div :key="$route.path" class="page-enter"><slot /></div>
    </main>
  </div>
</template>
