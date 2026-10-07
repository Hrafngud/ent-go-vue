<script setup lang="ts">
import type { User } from '../../api/users'
import UserAvatar from './UserAvatar.vue'
defineProps<{ user: User; currentUserId?: string }>()
</script>

<template>
  <div class="flex min-w-0 items-center gap-3">
    <UserAvatar :name="user.name" />
    <div class="min-w-0">
      <div class="flex flex-wrap items-center gap-2">
        <RouterLink :to="`/admin/users/${user.id}`" class="link link-hover max-w-xs break-words font-medium">{{ user.name }}</RouterLink>
        <span v-if="user.is_admin" class="badge badge-primary badge-soft badge-sm">Root</span>
        <span v-if="user.id === currentUserId" class="text-xs text-base-content/60">You</span>
      </div>
      <p class="mt-1 max-w-xs break-all text-xs text-base-content/60">{{ user.email }}</p>
    </div>
  </div>
</template>
