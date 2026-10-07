<script setup lang="ts">
import { ref } from 'vue'
import { PhEye, PhEyeSlash } from '@phosphor-icons/vue'
const model = defineModel<string>({ default: '' })
withDefaults(defineProps<{
  id: string
  label?: string
  disabled?: boolean
  required?: boolean
  autocomplete?: 'current-password' | 'new-password'
  hint?: string
  minLength?: number
}>(), { label: 'Password', autocomplete: 'new-password' })
const visible = ref(false)
</script>

<template>
  <div class="space-y-2">
    <label :for="id" class="block text-sm font-medium">{{ label }}</label>
    <div class="password-input">
      <input :id="id" v-model="model" :name="id" :type="visible ? 'text' : 'password'" class="input w-full"
        :autocomplete="autocomplete" :disabled="disabled" :required="required" :minlength="minLength" maxlength="72"
        :aria-describedby="hint ? `${id}-hint` : undefined" />
      <button class="password-toggle btn btn-ghost btn-sm gap-1 text-base-content/60" type="button" :disabled="disabled" :aria-controls="id" :aria-pressed="visible"
        :aria-label="`${visible ? 'Hide' : 'Show'} ${label.toLowerCase()}`" @click="visible = !visible">
        <PhEyeSlash v-if="visible" :size="16" aria-hidden="true" />
        <PhEye v-else :size="16" aria-hidden="true" />
        {{ visible ? 'Hide' : 'Show' }}
      </button>
    </div>
    <p v-if="hint" :id="`${id}-hint`" class="text-xs text-base-content/60">{{ hint }}</p>
  </div>
</template>
