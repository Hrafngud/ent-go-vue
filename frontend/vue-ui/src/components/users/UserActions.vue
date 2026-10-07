<script setup lang="ts">
import { PhEye, PhPencilSimple } from '@phosphor-icons/vue'
import type { User } from '../../api/users'
defineProps<{ user: User }>()
</script>

<template>
  <div class="flex justify-end gap-1">
    <RouterLink :to="`/admin/users/${user.id}`" class="btn btn-ghost btn-sm gap-1.5 text-base-content/65" :aria-label="`View ${user.name}`"><PhEye :size="16" aria-hidden="true" />View</RouterLink>
    <RouterLink :to="`/admin/users/${user.id}/edit`" class="btn btn-ghost btn-sm gap-1.5 text-base-content/65" :aria-label="`Edit ${user.name}`"><PhPencilSimple :size="16" aria-hidden="true" />Edit</RouterLink>
  </div>
</template>
