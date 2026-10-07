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
import UserAvatar from '../../components/users/UserAvatar.vue'
import { PhPencilSimple } from '@phosphor-icons/vue'

const route = useRoute()
const query = useUserQuery(() => String(route.params.id))
const user = computed(() => query.data.value?.data)
</script>

<template>
  <UsersBreadcrumbs label="User details" />
  <PageHeader title="User details">
    <RouterLink v-if="user" :to="`/admin/users/${user.id}/edit`" class="btn btn-primary"><PhPencilSimple :size="20" aria-hidden="true" />Edit user</RouterLink>
  </PageHeader>
  <LoadingState v-if="query.isPending.value" label="Loading user…" />
  <div v-else-if="query.isError.value" class="space-y-4">
    <FeedbackAlert :message="userError(query.error.value, 'Could not load this user. Please try again.')" />
    <button class="btn btn-outline" :disabled="query.isFetching.value" @click="query.refetch()">Try again</button>
  </div>
  <template v-else-if="user">
    <div class="max-w-3xl">
      <section class="surface p-6 sm:p-8" aria-label="Account details">
        <div class="flex items-center gap-4 border-b border-base-300 pb-6">
          <UserAvatar :name="user.name" size="lg" />
          <div class="min-w-0">
            <h2 class="break-words text-xl font-semibold tracking-tight">{{ user.name }}</h2>
            <span class="badge badge-sm mt-2" :class="user.is_admin ? 'badge-primary badge-soft' : 'badge-ghost'">{{ user.is_admin ? 'Root administrator' : 'Member' }}</span>
          </div>
        </div>
        <dl class="grid gap-6 pt-6 sm:grid-cols-2">
          <div><dt class="text-xs text-base-content/60">Email address</dt><dd class="mt-2 break-all text-sm">{{ user.email }}</dd></div>
          <div><dt class="text-xs text-base-content/60">Joined</dt><dd class="mt-2 text-sm">{{ formatDate(user.created_at) }}</dd></div>
          <div class="sm:col-span-2"><dt class="text-xs text-base-content/60">User ID</dt><dd class="mt-2 break-all text-sm text-base-content/70">{{ user.id }}</dd></div>
        </dl>
      </section>
      <DeleteUserPanel :user="user" />
    </div>
  </template>
</template>
