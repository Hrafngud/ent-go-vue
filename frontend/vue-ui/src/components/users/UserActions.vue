<script setup lang="ts">
import type { User } from '../../api/users'
defineProps<{ user: User }>()
</script>

<template>
  <div class="flex justify-end gap-1">
    <RouterLink :to="`/admin/users/${user.id}`" class="btn btn-ghost btn-sm" :aria-label="`View ${user.name}`">View</RouterLink>
    <RouterLink :to="`/admin/users/${user.id}/edit`" class="btn btn-ghost btn-sm" :aria-label="`Edit ${user.name}`">Edit</RouterLink>
  </div>
</template>
