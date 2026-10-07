<script setup lang="ts">
import BrandMark from './BrandMark.vue'
defineProps<{ title: string; description: string }>()
</script>

<template>
  <div class="flex min-h-screen flex-col">
    <header class="px-6 py-8 sm:px-12"><BrandMark /></header>
    <main id="main-content" class="flex flex-1 items-center px-6 py-10 sm:px-12">
      <section class="mx-auto w-full max-w-sm" aria-labelledby="auth-title">
        <p class="mb-3 text-xs font-medium uppercase tracking-widest text-base-content/60">Your workspace</p>
        <h1 id="auth-title" tabindex="-1" class="text-4xl font-semibold tracking-tight">{{ title }}</h1>
        <p class="mt-3 text-base-content/70">{{ description }}</p>
        <div class="mt-8"><slot /></div>
      </section>
    </main>
    <footer class="px-6 py-6 text-xs text-base-content/60 sm:px-12">Ent Go Vue</footer>
  </div>
</template>
