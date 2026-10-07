<script setup lang="ts">
import { nextTick, watch } from 'vue'
import { useRoute } from 'vue-router'
import AppShell from './components/layout/AppShell.vue'

const route = useRoute()
watch(() => route.fullPath, async () => {
  await nextTick()
  document.querySelector<HTMLElement>('h1')?.focus({ preventScroll: true })
})
</script>

<template>
  <div data-theme="dark" class="min-h-screen bg-base-100 text-base-content">
    <AppShell v-if="route.meta.authenticated || route.meta.admin">
      <RouterView :key="route.path" />
    </AppShell>
    <RouterView v-else :key="route.path" />
  </div>
</template>
