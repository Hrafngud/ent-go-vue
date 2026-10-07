<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import type { UserInput } from '../../api/users'
import PasswordField from '../forms/PasswordField.vue'
import FeedbackAlert from '../ui/FeedbackAlert.vue'
import { accountError, normalizeAccount } from '../../validation/account'
import { useCooldown } from '../../composables/useCooldown'
import { ApiError } from '../../api/client'

const props = withDefaults(defineProps<{
  initial?: UserInput
  editing?: boolean
  busy: boolean
  error: string
  submitLabel: string
  busyLabel?: string
  confirmPassword?: boolean
  lockEmail?: boolean
  retryAfter?: number
}>(), { editing: false, confirmPassword: false, busyLabel: 'Saving…' })
const emit = defineEmits<{ submit: [input: UserInput] }>()
const fields = reactive({ name: props.initial?.name || '', email: props.initial?.email || '', password: '' })
const confirmation = ref('')
const validationError = ref('')
const cooldown = useCooldown()
const cooldownSeconds = cooldown.seconds
watch(() => props.retryAfter, (delay) => { if (delay) cooldown.start(new ApiError(429, delay)) })

function submit() {
  if (props.busy || cooldownSeconds.value) return
  validationError.value = accountError(fields, props.editing)
  if (validationError.value) return
  if (props.confirmPassword && fields.password !== confirmation.value) {
    validationError.value = 'Your passwords do not match.'
    return
  }
  emit('submit', normalizeAccount(fields))
}
</script>

<template>
  <form class="w-full space-y-6" :class="{ 'surface max-w-3xl p-6 sm:p-8': !confirmPassword }" :aria-busy="busy" @submit.prevent="submit">
    <fieldset class="grid gap-6" :class="{ 'sm:grid-cols-2': !confirmPassword }" :disabled="busy">
      <legend :class="confirmPassword ? 'sr-only' : 'mb-6 text-base font-semibold'">Account details</legend>
      <div class="space-y-2">
        <label for="user-name" class="block text-sm font-medium">Full name</label>
        <input id="user-name" v-model="fields.name" class="input w-full" name="name" autocomplete="name" required maxlength="255" placeholder="Full name" />
      </div>
      <div class="space-y-2">
        <label for="user-email" class="block text-sm font-medium">Email address</label>
        <input id="user-email" v-model="fields.email" class="input w-full" name="email" type="email" autocomplete="email" required maxlength="254"
          placeholder="you@example.com" :readonly="lockEmail" :aria-describedby="lockEmail ? 'root-email-hint' : undefined" />
        <p v-if="lockEmail" id="root-email-hint" class="text-xs text-base-content/60">The root email is managed in the server configuration.</p>
      </div>
      <PasswordField id="user-password" v-model="fields.password" :disabled="busy" :required="!editing"
        :class="{ 'sm:col-span-2 sm:max-w-md': !confirmPassword }"
        :label="editing ? 'New password (optional)' : 'Password'" :hint="editing ? 'Leave blank to keep the current password. Otherwise use 8–72 bytes.' : 'Use 8–72 bytes. Some characters count as more than one byte.'" />
      <PasswordField v-if="confirmPassword" id="user-password-confirm" v-model="confirmation" label="Confirm password" :disabled="busy" required />
    </fieldset>
    <FeedbackAlert v-if="validationError || error" :message="validationError || error" />
    <div class="flex flex-wrap items-center gap-3" :class="{ 'justify-end border-t border-base-300 pt-6': !confirmPassword }">
      <slot name="actions" :busy="busy" />
      <button type="submit" class="btn btn-primary" :class="{ 'w-full': confirmPassword }" :disabled="busy || cooldownSeconds > 0">
        <span v-if="busy" class="loading loading-spinner loading-sm" aria-hidden="true"></span>
        {{ busy ? busyLabel : cooldownSeconds ? `Try again in ${cooldownSeconds}s` : submitLabel }}
      </button>
      <span v-if="busy" class="sr-only" role="status">{{ busyLabel }}</span>
    </div>
  </form>
</template>
