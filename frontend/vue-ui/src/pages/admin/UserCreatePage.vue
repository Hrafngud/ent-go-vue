<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useQueryClient } from '@tanstack/vue-query'
import { usersApi, type UserInput } from '../../api/users'
import { userError } from '../../api/errors'
import { useSessionStore } from '../../stores/session'
import { useFeedbackStore } from '../../stores/feedback'
import UsersBreadcrumbs from '../../components/users/UsersBreadcrumbs.vue'
import UserForm from '../../components/users/UserForm.vue'
import PageHeader from '../../components/ui/PageHeader.vue'

const session = useSessionStore()
const feedback = useFeedbackStore()
const router = useRouter()
const queryClient = useQueryClient()
const busy = ref(false)
const error = ref('')

async function create(input: UserInput) {
  if (busy.value) return
  busy.value = true
  error.value = ''
  try {
    const response = await usersApi.create(input, session.token)
    void queryClient.invalidateQueries({ queryKey: ['users'] })
    feedback.message = `${response.data.name} was created.`
    await router.push(`/admin/users/${response.data.id}`)
  } catch (err) {
    error.value = userError(err, 'Could not create this user. Please try again.')
  } finally { busy.value = false }
}
</script>

<template>
  <UsersBreadcrumbs label="Create user" />
  <PageHeader title="Create user" />
  <UserForm :busy="busy" :error="error" submit-label="Create user" busy-label="Creating user…" @submit="create">
    <template #actions><RouterLink v-if="!busy" to="/admin/users" class="btn btn-ghost">Cancel</RouterLink></template>
  </UserForm>
</template>
