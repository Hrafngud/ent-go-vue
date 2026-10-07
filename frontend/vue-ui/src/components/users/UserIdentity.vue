<script setup lang="ts">
import type { User } from '../../api/users'
defineProps<{ user: User; currentUserId?: string }>()
</script>

<template>
  <div>
    <div class="flex flex-wrap items-center gap-2">
      <RouterLink :to="`/admin/users/${user.id}`" class="link link-hover max-w-xs break-words font-semibold">{{ user.name }}</RouterLink>
      <span v-if="user.is_admin" class="badge badge-primary badge-outline badge-sm">Root</span>
      <span v-if="user.id === currentUserId" class="badge badge-ghost badge-sm">You</span>
    </div>
    <p class="mt-1 max-w-xs break-all text-sm text-base-content/60">{{ user.email }}</p>
  </div>
</template>
