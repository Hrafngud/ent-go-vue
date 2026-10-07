<script setup lang="ts">
import { computed, ref } from 'vue'
import { onBeforeRouteLeave, useRoute, useRouter } from 'vue-router'
import { PhPencilSimple } from '@phosphor-icons/vue'
import { formatDate } from '../../api/users'
import AppModal from '../ui/AppModal.vue'
import UserRecord from './UserRecord.vue'
import UserAvatar from './UserAvatar.vue'
import DeleteUserPanel from './DeleteUserPanel.vue'

const route = useRoute()
const router = useRouter()
const id = computed(() => String(route.params.id))
const busy = ref(false)
function close() { if (!busy.value) void router.push('/admin/users') }
onBeforeRouteLeave(() => !busy.value)
</script>

<template>
  <AppModal title="User details" @close="close" :busy="busy">
    <UserRecord :id="id" v-slot="{ user }">
      <section aria-label="Account details">
        <div class="flex items-center gap-4 border-b border-base-300 pb-6">
          <UserAvatar :name="user.name" size="lg" />
          <div class="min-w-0">
            <h3 class="break-words text-xl font-semibold tracking-tight">{{ user.name }}</h3>
            <span class="badge badge-sm mt-2" :class="user.is_admin ? 'badge-primary badge-soft' : 'badge-ghost'">{{ user.is_admin ? 'Root administrator' : 'Member' }}</span>
          </div>
        </div>
        <dl class="grid gap-6 py-6 sm:grid-cols-2">
          <div><dt class="text-xs text-base-content/60">Email address</dt><dd class="mt-2 break-all text-sm">{{ user.email }}</dd></div>
          <div><dt class="text-xs text-base-content/60">Joined</dt><dd class="mt-2 text-sm">{{ formatDate(user.created_at) }}</dd></div>
          <div class="sm:col-span-2"><dt class="text-xs text-base-content/60">User ID</dt><dd class="mt-2 break-all text-sm text-base-content/70">{{ user.id }}</dd></div>
        </dl>
        <RouterLink :to="`/admin/users/${user.id}/edit`" class="btn btn-primary"><PhPencilSimple :size="20" aria-hidden="true" />Edit user</RouterLink>
      </section>
      <DeleteUserPanel :user="user" @busy="busy = $event" />
    </UserRecord>
  </AppModal>
</template>
