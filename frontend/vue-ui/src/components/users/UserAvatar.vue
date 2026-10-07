<script setup lang="ts">
import { computed } from 'vue'
const props = withDefaults(defineProps<{ name: string; size?: 'sm' | 'lg' }>(), { size: 'sm' })
const initials = computed(() => props.name.trim().split(/\s+/).slice(0, 2).map(part => Array.from(part)[0] || '').join('').toLocaleUpperCase() || '?')
</script>

<template>
  <div class="avatar avatar-placeholder shrink-0" aria-hidden="true">
    <div class="rounded-full bg-primary/10 font-medium text-primary" :class="size === 'lg' ? 'size-16 text-xl' : 'size-10 text-sm'"><span>{{ initials }}</span></div>
  </div>
</template>
