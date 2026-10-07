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
  <PageHeader title="Edit user" description="Update account information or set a new password." />
  <LoadingState v-if="query.isPending.value" label="Loading user…" />
  <div v-else-if="query.isError.value" class="space-y-4">
    <FeedbackAlert :message="userError(query.error.value, 'Could not load this user. Please try again.')" />
    <button class="btn btn-outline" :disabled="query.isFetching.value" @click="query.refetch()">Try again</button>
  </div>
  <UserForm v-else-if="user" :initial="user" editing :lock-email="user.is_admin" :busy="busy" :error="error" submit-label="Save changes" @submit="update">
    <template #actions><RouterLink v-if="!busy" :to="`/admin/users/${user.id}`" class="btn btn-ghost">Cancel</RouterLink></template>
  </UserForm>
</template>
