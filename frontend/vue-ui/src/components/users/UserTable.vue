<script setup lang="ts">
import { formatDate, type User } from '../../api/users'
import UserIdentity from './UserIdentity.vue'
import UserActions from './UserActions.vue'
defineProps<{ users: User[]; currentUserId?: string }>()
</script>

<template>
  <ul class="list divide-y divide-base-300 border-y border-base-300 md:hidden" aria-label="Workspace users">
    <li v-for="user in users" :key="user.id" class="space-y-4 py-5">
      <UserIdentity :user="user" :current-user-id="currentUserId" />
      <div class="flex flex-wrap items-center justify-between gap-3">
        <p class="text-xs text-base-content/60">Created {{ formatDate(user.created_at) }}</p>
        <UserActions :user="user" />
      </div>
    </li>
  </ul>
  <div class="hidden overflow-x-auto border-y border-base-300 md:block">
    <table class="table">
      <caption class="sr-only">Workspace users and account actions</caption>
      <thead><tr><th scope="col">Account</th><th scope="col">Created</th><th scope="col" class="text-right">Actions</th></tr></thead>
      <tbody>
        <tr v-for="user in users" :key="user.id" class="hover:bg-base-200">
          <td>
            <UserIdentity :user="user" :current-user-id="currentUserId" />
          </td>
          <td class="whitespace-nowrap text-sm text-base-content/70">{{ formatDate(user.created_at) }}</td>
          <td>
            <UserActions :user="user" />
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
