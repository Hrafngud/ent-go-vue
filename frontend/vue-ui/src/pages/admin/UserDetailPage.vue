<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useUserQuery } from '../../queries/users'
import { userError } from '../../api/errors'
import { formatDate } from '../../api/users'
import UsersBreadcrumbs from '../../components/users/UsersBreadcrumbs.vue'
import DeleteUserPanel from '../../components/users/DeleteUserPanel.vue'
import PageHeader from '../../components/ui/PageHeader.vue'
import LoadingState from '../../components/ui/LoadingState.vue'
import FeedbackAlert from '../../components/ui/FeedbackAlert.vue'

const route = useRoute()
const query = useUserQuery(() => String(route.params.id))
const user = computed(() => query.data.value?.data)
</script>

<template>
  <UsersBreadcrumbs label="User details" />
  <PageHeader title="User details" description="Account information and management.">
    <RouterLink v-if="user" :to="`/admin/users/${user.id}/edit`" class="btn btn-primary">Edit user</RouterLink>
  </PageHeader>
  <LoadingState v-if="query.isPending.value" label="Loading user…" />
  <div v-else-if="query.isError.value" class="space-y-4">
    <FeedbackAlert :message="userError(query.error.value, 'Could not load this user. Please try again.')" />
    <button class="btn btn-outline" :disabled="query.isFetching.value" @click="query.refetch()">Try again</button>
  </div>
  <template v-else-if="user">
    <dl class="grid max-w-3xl gap-6 border-y border-base-300 py-6 sm:grid-cols-2">
      <div><dt class="text-sm text-base-content/60">Full name</dt><dd class="mt-2 break-words font-medium">{{ user.name }}</dd></div>
      <div><dt class="text-sm text-base-content/60">Email address</dt><dd class="mt-2 break-all">{{ user.email }}</dd></div>
      <div><dt class="text-sm text-base-content/60">Created</dt><dd class="mt-2">{{ formatDate(user.created_at) }}</dd></div>
      <div><dt class="text-sm text-base-content/60">Access</dt><dd class="mt-2"><span class="badge" :class="user.is_admin ? 'badge-primary badge-outline' : 'badge-ghost'">{{ user.is_admin ? 'Root administrator' : 'Member' }}</span></dd></div>
      <div class="sm:col-span-2"><dt class="text-sm text-base-content/60">User ID</dt><dd class="mt-2 break-all text-sm">{{ user.id }}</dd></div>
    </dl>
    <DeleteUserPanel :user="user" />
  </template>
</template>
