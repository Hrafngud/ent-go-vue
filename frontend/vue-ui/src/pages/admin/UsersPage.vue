<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { PhArrowClockwise, PhUserPlus } from '@phosphor-icons/vue'
import { useUsersQuery } from '../../queries/users'
import { useSessionStore } from '../../stores/session'
import { userError } from '../../api/errors'
import PageHeader from '../../components/ui/PageHeader.vue'
import LoadingState from '../../components/ui/LoadingState.vue'
import FeedbackAlert from '../../components/ui/FeedbackAlert.vue'
import UserTable from '../../components/users/UserTable.vue'

const session = useSessionStore()
const query = useUsersQuery()
const search = ref('')
const sort = ref('newest')
const page = ref(1)
const pageSize = 10
const users = computed(() => query.data.value?.data || [])
const filtered = computed(() => {
  const term = search.value.trim().toLocaleLowerCase()
  return users.value.filter(user => `${user.name} ${user.email}`.toLocaleLowerCase().includes(term)).sort((a, b) => {
    if (sort.value === 'name') return a.name.localeCompare(b.name)
    const delta = new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
    return sort.value === 'oldest' ? -delta : delta
  })
})
const pageCount = computed(() => Math.max(1, Math.ceil(filtered.value.length / pageSize)))
const visible = computed(() => filtered.value.slice((page.value - 1) * pageSize, page.value * pageSize))
watch([search, sort], () => { page.value = 1 })
watch(pageCount, count => { page.value = Math.min(page.value, count) })
</script>

<template>
  <PageHeader title="Users" description="Manage the people in your workspace." eyebrow="Administration">
    <RouterLink to="/admin/users/new" class="btn btn-primary"><PhUserPlus :size="20" aria-hidden="true" />Create user</RouterLink>
  </PageHeader>
  <div class="mb-6 flex flex-wrap items-end gap-4">
    <div class="w-full min-w-0 space-y-2 sm:w-auto sm:flex-1">
      <label for="users-search" class="block text-sm font-medium">Search users</label>
      <input id="users-search" v-model="search" class="input w-full" type="search" placeholder="Search by name or email" />
    </div>
    <div class="space-y-2">
      <label for="users-sort" class="block text-sm font-medium">Sort by</label>
      <select id="users-sort" v-model="sort" class="select w-full">
        <option value="newest">Newest first</option><option value="oldest">Oldest first</option><option value="name">Name, A–Z</option>
      </select>
    </div>
    <button class="btn btn-ghost" :disabled="query.isFetching.value" @click="query.refetch()">
      <span v-if="query.isFetching.value" class="loading loading-spinner loading-xs" aria-hidden="true"></span>
      <PhArrowClockwise v-else :size="20" aria-hidden="true" />
      Refresh
    </button>
  </div>
  <LoadingState v-if="query.isPending.value" label="Loading users…" />
  <div v-else-if="query.isError.value" class="space-y-4">
    <FeedbackAlert :message="userError(query.error.value, 'Could not load users. Please try again.')" />
    <button class="btn btn-outline" :disabled="query.isFetching.value" @click="query.refetch()">Try again</button>
  </div>
  <template v-else>
    <p class="mb-4 text-sm text-base-content/60" role="status">{{ filtered.length }} {{ filtered.length === 1 ? 'user' : 'users' }}{{ search.trim() ? ' matching your search' : '' }}</p>
    <UserTable v-if="visible.length" :users="visible" :current-user-id="session.user?.id" />
    <div v-else class="border-y border-base-300 py-12">
      <h2 class="text-xl font-semibold">{{ users.length ? 'No matching users' : 'Add your first user' }}</h2>
      <p class="mt-2 text-base-content/70">{{ users.length ? 'Try another name or email address.' : 'Create an account for someone joining your workspace.' }}</p>
      <button v-if="users.length" class="btn btn-ghost mt-4" @click="search = ''">Clear search</button>
      <RouterLink v-else to="/admin/users/new" class="btn btn-primary mt-4"><PhUserPlus :size="20" aria-hidden="true" />Create user</RouterLink>
    </div>
    <nav v-if="pageCount > 1" class="mt-6 flex flex-wrap items-center justify-between gap-4" aria-label="Users pagination">
      <p class="text-sm text-base-content/60" role="status">Page {{ page }} of {{ pageCount }}</p>
      <div class="join">
        <button class="btn btn-outline btn-sm join-item" :disabled="page === 1" @click="page--">Previous</button>
        <button class="btn btn-outline btn-sm join-item" :disabled="page === pageCount" @click="page++">Next</button>
      </div>
    </nav>
  </template>
</template>
