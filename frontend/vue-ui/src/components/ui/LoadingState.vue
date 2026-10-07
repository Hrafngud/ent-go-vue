<script setup lang="ts">
withDefaults(defineProps<{
  label?: string
  variant?: 'default' | 'table' | 'detail' | 'form'
}>(), { label: 'Loading…', variant: 'default' })
</script>

<template>
  <div :class="{ 'surface max-w-3xl p-6 sm:p-8': variant === 'detail' || variant === 'form' }">
    <p class="flex items-center gap-3 text-sm text-base-content/65" :class="variant === 'default' ? 'py-12' : variant === 'table' ? 'px-5 py-4' : 'mb-6'" role="status">
      <span class="loading loading-spinner loading-sm" aria-hidden="true"></span>{{ label }}
    </p>
    <div v-if="variant === 'table'" aria-hidden="true">
      <div class="hidden items-center justify-between border-y border-base-300 bg-base-100/40 px-6 py-4 md:flex">
        <div class="skeleton h-3 w-20"></div><div class="skeleton h-3 w-16"></div><div class="skeleton h-3 w-16"></div>
      </div>
      <div class="divide-y divide-base-300">
        <div v-for="row in 5" :key="row" class="flex flex-wrap items-center gap-4 px-5 py-5 md:px-6">
          <div class="skeleton size-10 shrink-0 rounded-full"></div>
          <div class="min-w-0 flex-1 space-y-2">
            <div class="skeleton h-4 w-32 max-w-full"></div><div class="skeleton h-3 w-48 max-w-full"></div>
          </div>
          <div class="skeleton hidden h-4 w-24 md:block"></div>
          <div class="flex w-full items-center justify-between gap-3 md:w-auto md:justify-end">
            <div class="skeleton h-3 w-24 md:hidden"></div>
            <div class="skeleton h-8 w-24"></div>
          </div>
        </div>
      </div>
      <div class="border-t border-base-300 px-5 py-4"><div class="skeleton h-3 w-20"></div></div>
    </div>
    <div v-else-if="variant === 'detail'" aria-hidden="true">
      <div class="flex items-center gap-4 border-b border-base-300 pb-6">
        <div class="skeleton size-14 rounded-full"></div>
        <div class="min-w-0 flex-1 space-y-3"><div class="skeleton h-6 w-48 max-w-full"></div><div class="skeleton h-4 w-20"></div></div>
      </div>
      <div class="grid gap-6 pt-6 sm:grid-cols-2">
        <div v-for="field in 3" :key="field" class="space-y-3"><div class="skeleton h-3 w-20"></div><div class="skeleton h-4 w-48 max-w-full"></div></div>
      </div>
    </div>
    <div v-else-if="variant === 'form'" class="space-y-6" aria-hidden="true">
      <div class="skeleton h-5 w-32"></div>
      <div class="grid gap-6 sm:grid-cols-2">
        <div v-for="field in 3" :key="field" class="space-y-2" :class="{ 'sm:col-span-2 sm:max-w-md': field === 3 }">
          <div class="skeleton h-4 w-24"></div><div class="skeleton h-10 w-full"></div>
        </div>
      </div>
      <div class="flex justify-end border-t border-base-300 pt-6"><div class="skeleton h-10 w-32"></div></div>
    </div>
  </div>
</template>
