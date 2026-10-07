<script setup lang="ts">
import { nextTick, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useQueryClient } from '@tanstack/vue-query'
import { usersApi, type User } from '../../api/users'
import { userError } from '../../api/errors'
import { useSessionStore } from '../../stores/session'
import { useFeedbackStore } from '../../stores/feedback'
import FeedbackAlert from '../ui/FeedbackAlert.vue'

const props = defineProps<{ user: User }>()
const open = ref(false)
const busy = ref(false)
const error = ref('')
const confirmButton = ref<HTMLButtonElement>()
const deleteButton = ref<HTMLButtonElement>()
const session = useSessionStore()
const feedback = useFeedbackStore()
const router = useRouter()
const queryClient = useQueryClient()

async function toggle(value: boolean) {
  open.value = value
  error.value = ''
  await nextTick()
  const button = value ? confirmButton.value : deleteButton.value
  button?.focus()
}

async function deleteUser() {
  if (busy.value) return
  busy.value = true
  error.value = ''
  try {
    await usersApi.delete(props.user.id, session.token)
    queryClient.removeQueries({ queryKey: ['users', props.user.id] })
    void queryClient.invalidateQueries({ queryKey: ['users'] })
    feedback.message = `${props.user.name} was deleted.`
    await router.push('/admin/users')
  } catch (err) {
    error.value = userError(err, 'Could not delete this user. Please try again.')
  } finally { busy.value = false }
}
</script>

<template>
  <section class="mt-12 border-t border-base-300 pt-6" aria-labelledby="delete-title">
    <h2 id="delete-title" class="text-lg font-semibold">Delete account</h2>
    <p v-if="user.is_admin" class="mt-2 text-sm text-base-content/70">The configured root account is protected from deletion.</p>
    <template v-else>
      <p class="mt-2 text-sm text-base-content/70">Deleting this user permanently removes their account.</p>
      <button v-if="!open" ref="deleteButton" class="btn btn-error btn-outline btn-sm mt-4" @click="toggle(true)">Delete user</button>
      <div v-else class="mt-4 max-w-xl space-y-4" :aria-busy="busy">
        <p class="break-words text-sm">Delete <strong>{{ user.name }}</strong> ({{ user.email }})? This cannot be undone.</p>
        <FeedbackAlert v-if="error" :message="error" />
        <div class="flex flex-wrap gap-3">
          <button ref="confirmButton" class="btn btn-error btn-sm" :disabled="busy" @click="deleteUser">
            <span v-if="busy" class="loading loading-spinner loading-xs" aria-hidden="true"></span>
            {{ busy ? 'Deleting…' : 'Yes, delete user' }}
          </button>
          <button class="btn btn-ghost btn-sm" :disabled="busy" @click="toggle(false)">Cancel</button>
        </div>
      </div>
    </template>
  </section>
</template>
