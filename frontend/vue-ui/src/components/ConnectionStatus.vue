<script setup lang="ts">
import { computed } from 'vue'
import { useHealthQuery, useReadinessQuery } from '../queries/health'

const health = useHealthQuery()
const readiness = useReadinessQuery()
const checking = computed(() => health.isFetching.value || readiness.isFetching.value)

function refresh() {
  void health.refetch()
  void readiness.refetch()
}
</script>

<template>
  <section class="mx-auto max-w-3xl px-6 py-12" aria-labelledby="connection-title">
    <div class="flex flex-wrap items-start justify-between gap-4">
      <div>
        <p class="text-sm font-medium tracking-wide text-base-content/70">ENT GO VUE</p>
        <h1 id="connection-title" class="mt-2 text-3xl font-semibold">Connection status</h1>
      </div>
      <button class="btn btn-outline btn-sm" :disabled="checking" @click="refresh">
        {{ checking ? 'Checking…' : 'Check again' }}
      </button>
    </div>

    <div class="mt-8 divide-y divide-base-300 border-y border-base-300" role="status" aria-live="polite">
      <p class="flex items-center gap-3 py-5">
        <span v-if="health.isPending.value" class="loading loading-spinner loading-xs" aria-hidden="true"></span>
        <span v-if="health.isPending.value">Backend: connecting…</span>
        <span v-else-if="health.isError.value" class="text-error">Backend: unavailable</span>
        <span v-else-if="health.data.value?.status === 'ok'" class="flex items-center gap-3">
          <span class="status status-success" aria-hidden="true"></span>
          Backend: connected
        </span>
        <span v-else class="text-error">Backend: unexpected response</span>
      </p>
      <p class="py-5">
        <span v-if="readiness.isPending.value">Database: checking…</span>
        <span v-else-if="readiness.isError.value" class="text-error">Database: unavailable</span>
        <span v-else-if="readiness.data.value?.database === 'connected'" class="flex items-center gap-3">
          <span class="status status-success" aria-hidden="true"></span>
          Database: connected
        </span>
        <span v-else class="text-error">Database: unexpected response</span>
      </p>
    </div>
    <p v-if="health.isError.value || readiness.isError.value" class="mt-4 text-sm">
      Connection could not be confirmed. Check again in a moment.
    </p>
    <p class="mt-4 text-sm text-base-content/70">Status refreshes every 30 seconds.</p>
  </section>
</template>
