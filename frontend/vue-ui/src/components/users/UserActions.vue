<script setup lang="ts">
import { PhEye, PhPencilSimple } from '@phosphor-icons/vue'
import type { User } from '../../api/users'
defineProps<{ user: User }>()
</script>

<template>
  <div class="flex justify-end gap-1">
    <RouterLink :to="`/admin/users/${user.id}`" class="btn btn-ghost btn-sm" :aria-label="`View ${user.name}`"><PhEye :size="20" aria-hidden="true" />View</RouterLink>
    <RouterLink :to="`/admin/users/${user.id}/edit`" class="btn btn-ghost btn-sm" :aria-label="`Edit ${user.name}`"><PhPencilSimple :size="20" aria-hidden="true" />Edit</RouterLink>
  </div>
</template>
