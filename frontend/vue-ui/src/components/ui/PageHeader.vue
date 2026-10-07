<script setup lang="ts">
defineProps<{ title: string; description?: string; eyebrow?: string }>()
</script>

<template>
  <header class="mb-8 flex flex-wrap items-end justify-between gap-5">
    <div class="min-w-0">
      <p v-if="eyebrow" class="mb-3 text-xs font-medium uppercase tracking-widest text-base-content/60">{{ eyebrow }}</p>
      <h1 tabindex="-1" class="break-words text-4xl font-semibold tracking-tight">{{ title }}</h1>
      <p v-if="description" class="mt-3 text-base-content/70">{{ description }}</p>
    </div>
    <slot />
  </header>
</template>
