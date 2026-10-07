<script setup lang="ts">
import { useSessionStore } from '../stores/session'
import PageHeader from '../components/ui/PageHeader.vue'
import ConnectionStatus from '../components/ConnectionStatus.vue'
const session = useSessionStore()
</script>

<template>
  <PageHeader title="Your workspace" :description="`Welcome back, ${session.user?.name || ''}.`" eyebrow="Overview">
    <RouterLink v-if="session.isAdmin" class="btn btn-primary" to="/admin/users">Manage users</RouterLink>
  </PageHeader>
  <ConnectionStatus />
</template>
