<script setup lang="ts">
import { computed } from 'vue'
import { PhArrowClockwise, PhCloud, PhDatabase } from '@phosphor-icons/vue'
import { useHealthQuery, useReadinessQuery } from '../queries/health'

const health = useHealthQuery()
const readiness = useReadinessQuery()
const checking = computed(() => health.isFetching.value || readiness.isFetching.value)
const connections = computed(() => [
  { name: 'Backend', icon: PhCloud, pending: health.isFetching.value || health.isPending.value, error: health.isError.value, connected: health.data.value?.status === 'ok' },
  { name: 'Database', icon: PhDatabase, pending: readiness.isFetching.value || readiness.isPending.value, error: readiness.isError.value, connected: readiness.data.value?.database === 'connected' },
])

function refresh() {
  void health.refetch()
  void readiness.refetch()
}
</script>

<template>
  <section class="surface p-6" aria-labelledby="connection-title">
    <div class="flex flex-wrap items-center justify-between gap-4 border-b border-base-300 pb-5">
      <h2 id="connection-title" class="text-base font-semibold">Connections</h2>
      <button class="btn btn-ghost btn-sm text-base-content/65" :disabled="checking" @click="refresh">
        <span v-if="checking" class="loading loading-spinner loading-xs" aria-hidden="true"></span>
        <PhArrowClockwise v-else :size="16" aria-hidden="true" />
        {{ checking ? 'Checking…' : 'Refresh' }}
      </button>
    </div>
    <ul class="divide-y divide-base-300" aria-live="polite" aria-label="Connection status">
      <li v-for="connection in connections" :key="connection.name" class="flex flex-wrap items-center justify-between gap-3 py-5 last:pb-0">
        <span class="flex items-center gap-3 text-sm"><component :is="connection.icon" :size="20" class="text-base-content/55" aria-hidden="true" />{{ connection.name }}</span>
        <span class="badge badge-soft badge-sm gap-2" :class="connection.pending ? 'badge-ghost' : connection.error || !connection.connected ? 'badge-error' : 'badge-success'">
          <span v-if="connection.pending" class="loading loading-spinner loading-xs" aria-hidden="true"></span>
          <span v-else class="status status-xs" aria-hidden="true"></span>
          {{ connection.pending ? 'Checking…' : connection.error ? 'Unavailable' : connection.connected ? 'Connected' : 'Unexpected response' }}
        </span>
      </li>
    </ul>
    <p v-if="health.isError.value || readiness.isError.value" class="mt-4 text-sm">
      Connection could not be confirmed. Check again in a moment.
    </p>
  </section>
</template>
