<script setup lang="ts">
import { nextTick, watch } from 'vue'
import { useRoute } from 'vue-router'
import { navigationError, navigationPending } from './router'
import AppShell from './components/layout/AppShell.vue'
import BrandMark from './components/layout/BrandMark.vue'
import LoadingState from './components/ui/LoadingState.vue'
import FeedbackAlert from './components/ui/FeedbackAlert.vue'

const route = useRoute()
function reload() { window.location.reload() }
watch(() => [route.fullPath, route.meta.modal] as const, async ([, modal], [, previousModal]) => {
  await nextTick()
  if (modal) return
  if (previousModal && document.activeElement instanceof HTMLElement && document.activeElement !== document.body) return
  document.querySelector<HTMLElement>('h1')?.focus({ preventScroll: true })
})
</script>

<template>
  <div data-theme="dark" class="min-h-screen bg-base-100 text-base-content">
    <div v-if="!route.matched.length" class="flex min-h-dvh flex-col items-center justify-center gap-8 px-6">
      <BrandMark />
      <LoadingState v-if="!navigationError" label="Opening your workspace…" />
      <template v-else>
        <FeedbackAlert message="Could not open this page. Please reload to try again." />
        <button class="btn btn-primary" @click="reload">Reload page</button>
      </template>
    </div>
    <template v-else>
      <div v-if="navigationPending" class="fixed inset-x-0 top-0 z-50" role="status">
        <progress class="progress progress-primary block h-1 w-full" aria-label="Loading page"></progress>
        <span class="sr-only">Loading page…</span>
      </div>
      <div v-if="navigationError" class="fixed inset-x-4 top-4 z-50 mx-auto max-w-xl">
        <FeedbackAlert message="Could not open this page. Please reload to try again." />
        <button class="btn btn-primary btn-sm mt-3" @click="reload">Reload page</button>
      </div>
      <AppShell v-if="route.meta.authenticated || route.meta.admin">
        <RouterView :key="route.matched[0]?.path" />
      </AppShell>
      <RouterView v-else :key="route.path" />
    </template>
  </div>
</template>
