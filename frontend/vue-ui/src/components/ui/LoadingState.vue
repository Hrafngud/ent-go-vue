<script setup lang="ts">
withDefaults(defineProps<{ label?: string }>(), { label: 'Loading…' })
</script>

<template>
  <div class="flex items-center gap-3 py-12" role="status">
    <span class="loading loading-spinner loading-sm" aria-hidden="true"></span>{{ label }}
  </div>
</template>
