<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { PhArrowClockwise, PhMagnifyingGlass, PhUserPlus } from '@phosphor-icons/vue'
import { useUsersQuery } from '../../queries/users'
import { useSessionStore } from '../../stores/session'
import { userError } from '../../api/errors'
import PageHeader from '../../components/ui/PageHeader.vue'
import LoadingState from '../../components/ui/LoadingState.vue'
import FeedbackAlert from '../../components/ui/FeedbackAlert.vue'
import UserTable from '../../components/users/UserTable.vue'
import CustomSelect from '../../components/forms/CustomSelect.vue'

const session = useSessionStore()
const query = useUsersQuery()
const search = ref('')
const sort = ref('newest')
const sortOptions = [
  { value: 'newest', label: 'Newest first' },
  { value: 'oldest', label: 'Oldest first' },
  { value: 'name', label: 'Name, A–Z' },
]
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
  <PageHeader title="Users">
    <RouterLink to="/admin/users/new" class="btn btn-primary"><PhUserPlus :size="20" aria-hidden="true" />Create user</RouterLink>
  </PageHeader>
  <section class="surface" aria-label="User directory">
    <div class="flex flex-wrap items-center gap-3 border-b border-base-300 p-5">
      <div class="w-full min-w-0 sm:w-auto sm:flex-1">
        <label for="users-search" class="sr-only">Search users</label>
        <div class="input w-full sm:max-w-sm">
          <PhMagnifyingGlass :size="20" class="shrink-0 text-base-content/55" aria-hidden="true" />
          <input id="users-search" v-model="search" class="min-w-0 grow" type="search" placeholder="Search users…" />
        </div>
      </div>
      <div class="min-w-0 flex-1 sm:w-44 sm:flex-none">
        <CustomSelect id="users-sort" v-model="sort" label="Sort by" :options="sortOptions" />
      </div>
      <button class="btn btn-ghost text-base-content/65" :disabled="query.isFetching.value" @click="query.refetch()">
        <span v-if="query.isFetching.value" class="loading loading-spinner loading-xs" aria-hidden="true"></span>
        <PhArrowClockwise v-else :size="20" aria-hidden="true" />{{ query.isPending.value ? 'Loading…' : query.isFetching.value ? 'Refreshing…' : 'Refresh' }}
      </button>
    </div>
    <LoadingState v-if="query.isPending.value" variant="table" label="Loading users…" />
    <div v-else-if="query.isError.value && !query.data.value" class="space-y-4 p-6">
      <FeedbackAlert :message="userError(query.error.value, 'Could not load users. Please try again.')" />
      <button class="btn btn-outline" :disabled="query.isFetching.value" @click="query.refetch()">
        <span v-if="query.isFetching.value" class="loading loading-spinner loading-xs" aria-hidden="true"></span>
        {{ query.isFetching.value ? 'Retrying…' : 'Try again' }}
      </button>
    </div>
    <template v-else>
      <FeedbackAlert v-if="query.isError.value" class="mx-5 mt-4" message="Could not refresh users. Showing the last loaded results. Please try again." />
      <p v-if="query.isFetching.value" class="sr-only" role="status">Updating users…</p>
      <div :aria-busy="query.isFetching.value">
        <UserTable v-if="visible.length" :users="visible" :current-user-id="session.user?.id" />
        <div v-else class="px-6 py-16 text-center">
          <h2 class="text-lg font-semibold">{{ users.length ? 'No matching users' : 'Add your first user' }}</h2>
          <p class="mt-2 text-sm text-base-content/65">{{ users.length ? 'Try another name or email.' : 'Invite someone into your workspace.' }}</p>
          <button v-if="users.length" class="btn btn-ghost btn-sm mt-5" @click="search = ''">Clear search</button>
          <RouterLink v-else to="/admin/users/new" class="btn btn-primary btn-sm mt-5"><PhUserPlus :size="16" aria-hidden="true" />Create user</RouterLink>
        </div>
      </div>
      <div class="flex flex-wrap items-center justify-between gap-4 border-t border-base-300 px-5 py-4">
        <p class="text-xs text-base-content/60" role="status">{{ filtered.length }} {{ filtered.length === 1 ? 'user' : 'users' }}{{ search.trim() ? ' found' : '' }}</p>
        <nav v-if="pageCount > 1" class="flex flex-wrap items-center gap-4" aria-label="Users pagination">
          <p class="text-xs text-base-content/60" role="status">Page {{ page }} of {{ pageCount }}</p>
          <div class="join">
            <button class="btn btn-ghost btn-sm join-item" :disabled="page === 1" @click="page--">Previous</button>
            <button class="btn btn-ghost btn-sm join-item" :disabled="page === pageCount" @click="page++">Next</button>
          </div>
        </nav>
      </div>
    </template>
  </section>
</template>
