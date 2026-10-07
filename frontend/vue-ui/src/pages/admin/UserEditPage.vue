<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQueryClient } from '@tanstack/vue-query'
import { useUserQuery } from '../../queries/users'
import { usersApi, type UserInput } from '../../api/users'
import { userError } from '../../api/errors'
import { useSessionStore } from '../../stores/session'
import { useFeedbackStore } from '../../stores/feedback'
import UsersBreadcrumbs from '../../components/users/UsersBreadcrumbs.vue'
import UserForm from '../../components/users/UserForm.vue'
import PageHeader from '../../components/ui/PageHeader.vue'
import LoadingState from '../../components/ui/LoadingState.vue'
import FeedbackAlert from '../../components/ui/FeedbackAlert.vue'

const route = useRoute()
const router = useRouter()
const session = useSessionStore()
const feedback = useFeedbackStore()
const queryClient = useQueryClient()
const query = useUserQuery(() => String(route.params.id))
const user = computed(() => query.data.value?.data)
const busy = ref(false)
const error = ref('')

async function update(input: UserInput) {
  if (busy.value || !user.value) return
  busy.value = true
  error.value = ''
  try {
    const response = await usersApi.update(user.value.id, input, session.token)
    queryClient.setQueryData(['users', response.data.id], response)
    void queryClient.invalidateQueries({ queryKey: ['users'] })
    if (response.data.id === session.user?.id) session.user = response.data
    feedback.message = `${response.data.name} was updated.`
    await router.push(`/admin/users/${response.data.id}`)
  } catch (err) {
    error.value = userError(err, 'Could not save changes. Please try again.')
  } finally { busy.value = false }
}
</script>

<template>
  <UsersBreadcrumbs label="Edit user" />
  <PageHeader title="Edit user" :description="user?.name" />
  <LoadingState v-if="query.isPending.value" variant="form" label="Loading user…" />
  <div v-else-if="query.isError.value && !user" class="space-y-4">
    <FeedbackAlert :message="userError(query.error.value, 'Could not load this user. Please try again.')" />
    <button class="btn btn-outline" :disabled="query.isFetching.value" @click="query.refetch()">
      <span v-if="query.isFetching.value" class="loading loading-spinner loading-xs" aria-hidden="true"></span>
      {{ query.isFetching.value ? 'Retrying…' : 'Try again' }}
    </button>
  </div>
  <template v-else-if="user">
    <FeedbackAlert v-if="query.isError.value" class="mb-4 max-w-3xl" message="Could not refresh this user. Showing the last loaded details." />
    <p v-if="query.isFetching.value" class="mb-4 flex items-center gap-2 text-sm text-base-content/65" role="status">
      <span class="loading loading-spinner loading-xs" aria-hidden="true"></span>Updating user…
    </p>
    <UserForm :initial="user" editing :lock-email="user.is_admin" :busy="busy" :error="error" submit-label="Save changes" @submit="update">
      <template #actions><RouterLink v-if="!busy" :to="`/admin/users/${user.id}`" class="btn btn-ghost">Cancel</RouterLink></template>
    </UserForm>
  </template>
</template>
