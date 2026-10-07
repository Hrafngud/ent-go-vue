<script setup lang="ts">
import { computed } from 'vue'
import { onBeforeRouteLeave, useRoute, useRouter } from 'vue-router'
import type { UserInput } from '../../api/users'
import { ApiError } from '../../api/client'
import { userError } from '../../api/errors'
import { useSaveUserMutation } from '../../queries/users'
import { useFeedbackStore } from '../../stores/feedback'
import AppModal from '../ui/AppModal.vue'
import UserForm from './UserForm.vue'
import UserRecord from './UserRecord.vue'

const props = defineProps<{ editing?: boolean }>()
const route = useRoute()
const router = useRouter()
const feedback = useFeedbackStore()
const mutation = useSaveUserMutation()
const busy = mutation.isPending
const id = computed(() => String(route.params.id || ''))
const error = computed(() => mutation.error.value ? userError(mutation.error.value, 'Could not save this user. Please try again.') : '')
const retryAfter = computed(() => mutation.error.value instanceof ApiError && mutation.error.value.status === 429 ? mutation.error.value.retryAfter : 0)

function close() { if (!busy.value) void router.push('/admin/users') }
onBeforeRouteLeave(() => !busy.value)

async function save(input: UserInput) {
  if (busy.value) return
  try {
    const response = await mutation.mutateAsync({ id: props.editing ? id.value : undefined, input })
    feedback.message = `${response.data.name} was ${props.editing ? 'updated' : 'created'}.`
    await router.replace(`/admin/users/${response.data.id}`)
  } catch { /* The mutation exposes the API error to FormKit and preserves the draft. */ }
}
</script>

<template>
  <AppModal :title="editing ? 'Edit user' : 'Create user'" :busy="busy" @close="close"
    :description="editing ? 'Update account details or set a new password.' : 'Give someone a place in your workspace.'">
    <UserRecord v-if="editing" :id="id" v-slot="{ user }" variant="form">
      <UserForm :key="user.id" :initial="user" editing :lock-email="user.is_admin" :busy="busy" :error="error" :retry-after="retryAfter" :save="save" submit-label="Save changes">
        <template #actions><button type="button" class="btn btn-ghost" :disabled="busy" @click="close">Cancel</button></template>
      </UserForm>
    </UserRecord>
    <UserForm v-else :busy="busy" :error="error" :retry-after="retryAfter" :save="save" submit-label="Create user" busy-label="Creating user…">
      <template #actions><button type="button" class="btn btn-ghost" :disabled="busy" @click="close">Cancel</button></template>
    </UserForm>
  </AppModal>
</template>
