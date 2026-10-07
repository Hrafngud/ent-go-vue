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
    <div class="flex gap-2">
      <input :id="id" v-model="model" :name="id" :type="visible ? 'text' : 'password'" class="input min-w-0 flex-1"
        :autocomplete="autocomplete" :disabled="disabled" :required="required" :minlength="minLength"
        :aria-describedby="hint ? `${id}-hint` : undefined" />
      <button class="btn btn-ghost" type="button" :disabled="disabled" :aria-controls="id" :aria-pressed="visible"
        :aria-label="`${visible ? 'Hide' : 'Show'} ${label.toLowerCase()}`" @click="visible = !visible">
        <PhEyeSlash v-if="visible" :size="20" aria-hidden="true" />
        <PhEye v-else :size="20" aria-hidden="true" />
        {{ visible ? 'Hide' : 'Show' }}
      </button>
    </div>
    <p v-if="hint" :id="`${id}-hint`" class="text-xs text-base-content/60">{{ hint }}</p>
  </div>
</template>
