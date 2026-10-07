<script setup lang="ts">
import { computed } from 'vue'
import { useUserQuery } from '../../queries/users'
import { userError } from '../../api/errors'
import LoadingState from '../ui/LoadingState.vue'
import FeedbackAlert from '../ui/FeedbackAlert.vue'

const props = defineProps<{ id: string; variant?: 'form' | 'detail' }>()
const query = useUserQuery(() => props.id)
const user = computed(() => query.data.value?.data)
</script>

<template>
  <LoadingState v-if="query.isPending.value" :variant="variant || 'detail'" label="Loading user…" />
  <div v-else-if="query.isError.value && !user" class="space-y-4">
    <FeedbackAlert :message="userError(query.error.value, 'Could not load this user. Please try again.')" />
    <button class="btn btn-outline" :disabled="query.isFetching.value" @click="query.refetch()">
      <span v-if="query.isFetching.value" class="loading loading-spinner loading-xs" aria-hidden="true"></span>
      {{ query.isFetching.value ? 'Retrying…' : 'Try again' }}
    </button>
  </div>
  <template v-else-if="user">
    <FeedbackAlert v-if="query.isError.value" class="mb-4" message="Could not refresh this user. Showing the last loaded details." />
    <p v-if="query.isFetching.value" class="mb-4 flex items-center gap-2 text-sm text-base-content/65" role="status">
      <span class="loading loading-spinner loading-xs" aria-hidden="true"></span>Refreshing account…
    </p>
    <slot :user="user" />
  </template>
</template>
