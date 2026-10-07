<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import type { User } from '../../api/users'
import { useDeleteUserMutation } from '../../queries/users'
import { userError } from '../../api/errors'
import { useFeedbackStore } from '../../stores/feedback'
import FeedbackAlert from '../ui/FeedbackAlert.vue'

const props = defineProps<{ user: User }>()
const emit = defineEmits<{ busy: [value: boolean] }>()
const open = ref(false)
const mutation = useDeleteUserMutation()
const busy = mutation.isPending
const error = computed(() => mutation.error.value ? userError(mutation.error.value, 'Could not delete this user. Please try again.') : '')
watch(busy, value => emit('busy', value), { flush: 'sync' })
const confirmButton = ref<HTMLButtonElement>()
const deleteButton = ref<HTMLButtonElement>()
const feedback = useFeedbackStore()
const router = useRouter()

async function toggle(value: boolean) {
  open.value = value
  mutation.reset()
  await nextTick()
  const button = value ? confirmButton.value : deleteButton.value
  button?.focus()
}

async function deleteUser() {
  if (busy.value) return
  try {
    await mutation.mutateAsync(props.user.id)
    feedback.message = `${props.user.name} was deleted.`
    await router.push('/admin/users')
  } catch { /* The mutation retains the error for the confirmation panel. */ }
}
</script>

<template>
  <section class="mt-8 border-t border-base-300 pt-6" aria-labelledby="delete-title">
    <div class="flex flex-wrap items-center justify-between gap-4">
      <div>
        <h2 id="delete-title" class="text-sm font-medium">{{ user.is_admin ? 'Account protection' : 'Delete account' }}</h2>
        <p class="mt-1 text-xs leading-relaxed text-base-content/60">{{ user.is_admin ? 'The configured root account is protected from deletion.' : 'Permanently removes this account.' }}</p>
      </div>
      <button v-if="!user.is_admin && !open" ref="deleteButton" class="btn btn-error btn-soft btn-sm" @click="toggle(true)">Delete user</button>
    </div>
    <template v-if="!user.is_admin && open">
      <div class="mt-4 space-y-4 rounded-box border border-error/20 bg-error/5 p-5" :aria-busy="busy">
        <p class="break-words text-sm">Delete <strong>{{ user.name }}</strong> ({{ user.email }})? This cannot be undone.</p>
        <FeedbackAlert v-if="error" :message="error" />
        <div class="flex flex-wrap gap-3">
          <button ref="confirmButton" class="btn btn-error btn-sm" :disabled="busy" @click="deleteUser">
            <span v-if="busy" class="loading loading-spinner loading-xs" aria-hidden="true"></span>
            {{ busy ? 'Deleting…' : 'Yes, delete user' }}
          </button>
          <button class="btn btn-ghost btn-sm" :disabled="busy" @click="toggle(false)">Cancel</button>
          <span v-if="busy" class="sr-only" role="status">Deleting user…</span>
        </div>
      </div>
    </template>
  </section>
</template>
