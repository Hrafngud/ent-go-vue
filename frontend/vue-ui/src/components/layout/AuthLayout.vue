<script setup lang="ts">
import BrandMark from './BrandMark.vue'
defineProps<{ title: string; description?: string; eyebrow?: string }>()
</script>

<template>
  <div class="flex min-h-dvh flex-col">
    <header class="px-5 py-6 sm:px-10 sm:py-8"><BrandMark /></header>
    <main id="main-content" class="flex flex-1 items-center justify-center px-5 py-6 sm:py-10">
      <section class="surface page-enter w-full max-w-md p-6 sm:p-10" aria-labelledby="auth-title">
        <p v-if="eyebrow" class="mb-5 text-xs font-medium tracking-widest text-primary uppercase">{{ eyebrow }}</p>
        <h1 id="auth-title" tabindex="-1" class="auth-title">{{ title }}</h1>
        <p v-if="description" class="mt-3 text-sm leading-relaxed text-base-content/65">{{ description }}</p>
        <div class="mt-8"><slot /></div>
      </section>
    </main>
  </div>
</template>
