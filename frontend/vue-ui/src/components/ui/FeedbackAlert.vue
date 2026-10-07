<script setup lang="ts">
import { PhCheckCircle, PhInfo, PhWarningCircle, PhX } from '@phosphor-icons/vue'
withDefaults(defineProps<{ message: string; kind?: 'error' | 'success' | 'info'; dismissible?: boolean }>(), { kind: 'error' })
defineEmits<{ dismiss: [] }>()
</script>

<template>
  <div class="alert alert-soft" :class="{ 'alert-error': kind === 'error', 'alert-success': kind === 'success', 'alert-info': kind === 'info' }" :role="kind === 'error' ? 'alert' : 'status'">
    <PhWarningCircle v-if="kind === 'error'" :size="20" aria-hidden="true" />
    <PhCheckCircle v-else-if="kind === 'success'" :size="20" aria-hidden="true" />
    <PhInfo v-else :size="20" aria-hidden="true" />
    <span class="min-w-0 break-words">{{ message }}</span>
    <button v-if="dismissible" type="button" class="btn btn-ghost btn-xs" aria-label="Dismiss message" @click="$emit('dismiss')"><PhX :size="16" aria-hidden="true" />Dismiss</button>
  </div>
</template>
