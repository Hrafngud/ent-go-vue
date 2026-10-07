<script setup lang="ts">
import { watch } from 'vue'
import { isNode, type FormKitNode } from '@formkit/core'
import type { UserInput } from '../../api/users'
import PasswordField from '../forms/PasswordField.vue'
import { useCooldown } from '../../composables/useCooldown'
import { ApiError } from '../../api/client'

const props = withDefaults(defineProps<{
  initial?: UserInput
  editing?: boolean
  busy: boolean
  error: string
  save: (input: UserInput) => Promise<void>
  submitLabel: string
  busyLabel?: string
  confirmPassword?: boolean
  lockEmail?: boolean
  retryAfter?: number
}>(), { editing: false, confirmPassword: false, busyLabel: 'Saving…' })
// Seed once per mounted account. Query refreshes must not replace an in-progress draft.
const initialValues = { name: props.initial?.name || '', email: props.initial?.email || '', password: '' }
const cooldown = useCooldown()
const cooldownSeconds = cooldown.seconds
watch(() => props.retryAfter, (delay) => { if (delay) cooldown.start(new ApiError(429, delay)) })

async function submit(values: UserInput) {
  if (props.busy || cooldownSeconds.value) return
  await props.save({ name: values.name, email: values.email, ...(values.password ? { password: values.password } : {}) })
}

function focusInvalid(node: FormKitNode) {
  const invalid = node.children.find(child => isNode(child) && child.ledger.value('blocking'))
  if (invalid) document.getElementById(String(invalid.props.id))?.focus()
}
</script>

<template>
  <FormKit v-slot="{ state }" type="form" :value="initialValues" :actions="false" form-class="w-full space-y-6"
    :disabled="busy || cooldownSeconds > 0" :errors="error ? [error] : []" :aria-busy="busy"
    @submit="submit" @submit-invalid="focusInvalid">
    <fieldset class="grid gap-6" :class="{ 'sm:grid-cols-2': !confirmPassword }">
      <legend class="sr-only">Account details</legend>
      <FormKit type="text" name="name" label="Full name" autocomplete="name" placeholder="Full name"
        sanitize="name" validation="required|account_name" />
      <FormKit type="email" name="email" label="Email address" autocomplete="email" placeholder="you@example.com"
        sanitize="email" validation="required|account_email" :readonly="lockEmail"
        :help="lockEmail ? 'The root email is managed in the server configuration.' : undefined" />
      <PasswordField id="user-password" name="password" :class="{ 'sm:col-span-2': !confirmPassword }"
        :validation="editing ? 'password_bytes' : 'required|password_bytes'"
        :label="editing ? 'New password (optional)' : 'Password'"
        :help="editing ? 'Leave blank to keep the current password. Otherwise use 8–72 bytes.' : 'Use 8–72 bytes. Some characters count as more than one byte.'" />
      <PasswordField v-if="confirmPassword" id="user-password-confirm" name="password_confirm" label="Confirm password" validation="required|confirm:password" />
    </fieldset>
    <div class="flex flex-wrap items-center gap-3" :class="{ 'justify-end border-t border-base-300 pt-6': !confirmPassword }">
      <slot name="actions" :busy="busy" :dirty="state.dirty" />
      <button type="submit" class="btn btn-primary" :class="{ 'w-full': confirmPassword }" :disabled="busy || cooldownSeconds > 0">
        <span v-if="busy" class="loading loading-spinner loading-sm" aria-hidden="true"></span>
        {{ busy ? busyLabel : cooldownSeconds ? `Try again in ${cooldownSeconds}s` : submitLabel }}
      </button>
      <span v-if="busy" class="sr-only" role="status">{{ busyLabel }}</span>
    </div>
  </FormKit>
</template>
