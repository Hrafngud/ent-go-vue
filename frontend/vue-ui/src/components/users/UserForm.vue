<script setup lang="ts">
import { reactive, ref } from 'vue'
import type { UserInput } from '../../api/users'
import PasswordField from '../forms/PasswordField.vue'
import FeedbackAlert from '../ui/FeedbackAlert.vue'

const props = withDefaults(defineProps<{
  initial?: UserInput
  editing?: boolean
  busy: boolean
  error: string
  submitLabel: string
  confirmPassword?: boolean
  lockEmail?: boolean
}>(), { editing: false, confirmPassword: false })
const emit = defineEmits<{ submit: [input: UserInput] }>()
const fields = reactive({ name: props.initial?.name || '', email: props.initial?.email || '', password: '' })
const confirmation = ref('')
const validationError = ref('')

function submit() {
  if (props.busy) return
  validationError.value = ''
  if (!fields.name.trim()) { validationError.value = 'Enter a name.'; return }
  const bytes = new TextEncoder().encode(fields.password).length
  if ((!props.editing || fields.password) && (bytes < 6 || bytes > 72)) {
    validationError.value = 'Use a password between 6 and 72 bytes. Some characters count as more than one byte.'
    return
  }
  if (props.confirmPassword && fields.password !== confirmation.value) {
    validationError.value = 'Your passwords do not match.'
    return
  }
  emit('submit', {
    name: fields.name.trim(), email: fields.email.trim(),
    ...(fields.password ? { password: fields.password } : {}),
  })
}
</script>

<template>
  <form class="max-w-xl space-y-6" :aria-busy="busy" @submit.prevent="submit">
    <fieldset class="space-y-6" :disabled="busy">
      <legend class="sr-only">Account information</legend>
      <div class="space-y-2">
        <label for="user-name" class="block text-sm font-medium">Full name</label>
        <input id="user-name" v-model="fields.name" class="input w-full" name="name" autocomplete="name" required maxlength="255" placeholder="Full name" />
      </div>
      <div class="space-y-2">
        <label for="user-email" class="block text-sm font-medium">Email address</label>
        <input id="user-email" v-model="fields.email" class="input w-full" name="email" type="email" autocomplete="email" required maxlength="255"
          placeholder="you@example.com" :readonly="lockEmail" :aria-describedby="lockEmail ? 'root-email-hint' : undefined" />
        <p v-if="lockEmail" id="root-email-hint" class="text-xs text-base-content/60">The root email is managed in the server configuration.</p>
      </div>
      <PasswordField id="user-password" v-model="fields.password" :disabled="busy" :required="!editing" :min-length="6"
        :label="editing ? 'New password (optional)' : 'Password'" :hint="editing ? 'Leave blank to keep the current password.' : 'Use at least 6 characters, up to 72 bytes.'" />
      <PasswordField v-if="confirmPassword" id="user-password-confirm" v-model="confirmation" label="Confirm password" :disabled="busy" required :min-length="6" />
    </fieldset>
    <FeedbackAlert v-if="validationError || error" :message="validationError || error" />
    <div class="flex flex-wrap items-center gap-3">
      <button type="submit" class="btn btn-primary" :class="{ 'w-full': confirmPassword }" :disabled="busy">
        <span v-if="busy" class="loading loading-spinner loading-sm" aria-hidden="true"></span>
        {{ busy ? 'Saving…' : submitLabel }}
      </button>
      <slot name="actions" :busy="busy" />
    </div>
  </form>
</template>
