<script setup lang="ts">
import { ref } from 'vue'
import { PhEye, PhEyeSlash } from '@phosphor-icons/vue'

withDefaults(defineProps<{
  id: string
  name: string
  label?: string
  autocomplete?: 'current-password' | 'new-password'
  help?: string
  validation?: string
}>(), { label: 'Password', autocomplete: 'new-password', validation: 'required|password_bytes' })
const visible = ref(false)
</script>

<template>
  <FormKit :id="id" :name="name" :type="visible ? 'text' : 'password'" :label="label"
    :autocomplete="autocomplete" :help="help" :validation="validation" input-class="pr-24">
    <template #suffix="{ disabled }">
      <button class="password-toggle btn btn-ghost btn-sm gap-1 text-base-content/60" type="button" :disabled="disabled"
        :aria-controls="id" :aria-pressed="visible" :aria-label="`${visible ? 'Hide' : 'Show'} ${label.toLowerCase()}`" @click="visible = !visible">
        <PhEyeSlash v-if="visible" :size="16" aria-hidden="true" />
        <PhEye v-else :size="16" aria-hidden="true" />
        {{ visible ? 'Hide' : 'Show' }}
      </button>
    </template>
  </FormKit>
</template>
