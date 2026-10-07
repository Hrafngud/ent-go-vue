<script setup lang="ts">
withDefaults(defineProps<{ message: string; kind?: 'error' | 'success' | 'info'; dismissible?: boolean }>(), { kind: 'error' })
defineEmits<{ dismiss: [] }>()
</script>

<template>
  <div class="alert alert-soft" :class="{ 'alert-error': kind === 'error', 'alert-success': kind === 'success', 'alert-info': kind === 'info' }" :role="kind === 'error' ? 'alert' : 'status'">
    <span class="min-w-0 break-words">{{ message }}</span>
    <button v-if="dismissible" type="button" class="btn btn-ghost btn-xs" aria-label="Dismiss message" @click="$emit('dismiss')">Dismiss</button>
  </div>
</template>
