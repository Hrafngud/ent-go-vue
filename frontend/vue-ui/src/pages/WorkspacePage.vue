<script setup lang="ts">
import { useSessionStore } from '../stores/session'
import PageHeader from '../components/ui/PageHeader.vue'
import ConnectionStatus from '../components/ConnectionStatus.vue'
import UserAvatar from '../components/users/UserAvatar.vue'
import { PhArrowUpRight } from '@phosphor-icons/vue'
const session = useSessionStore()
</script>

<template>
  <PageHeader title="Overview" :description="`Welcome back, ${session.user?.name || ''}.`">
    <RouterLink v-if="session.isAdmin" class="btn btn-primary" to="/admin/users">Manage users<PhArrowUpRight :size="20" aria-hidden="true" /></RouterLink>
  </PageHeader>
  <div class="grid items-start gap-6 xl:grid-cols-3">
    <ConnectionStatus class="xl:col-span-2" />
    <section v-if="session.user" class="surface p-6" aria-labelledby="account-title">
      <h2 id="account-title" class="text-base font-semibold">Your account</h2>
      <div class="mt-6 flex items-center gap-3">
        <UserAvatar :name="session.user.name" />
        <div class="min-w-0">
          <p class="break-words text-sm font-medium">{{ session.user.name }}</p>
          <p class="mt-1 text-xs text-base-content/60">{{ session.isAdmin ? 'Root administrator' : 'Member' }}</p>
        </div>
      </div>
      <p class="mt-5 break-all border-t border-base-300 pt-5 text-sm text-base-content/65">{{ session.user.email }}</p>
    </section>
  </div>
</template>
