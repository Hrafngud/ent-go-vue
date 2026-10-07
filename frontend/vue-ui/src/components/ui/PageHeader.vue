<script setup lang="ts">
defineProps<{ title: string; description?: string; eyebrow?: string }>()
</script>

<template>
  <header class="mb-8 flex flex-wrap items-center justify-between gap-5">
    <div class="min-w-0">
      <p v-if="eyebrow" class="mb-3 text-xs font-medium tracking-widest text-primary uppercase">{{ eyebrow }}</p>
      <h1 tabindex="-1" class="page-title break-words">{{ title }}</h1>
      <p v-if="description" class="mt-2 text-sm text-base-content/65">{{ description }}</p>
    </div>
    <slot />
  </header>
</template>
