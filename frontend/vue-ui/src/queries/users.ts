import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'
import { ApiError } from '../api/client'
import { usersApi } from '../api/users'
import { useSessionStore } from '../stores/session'

const retry = (count: number, error: Error) => !(error instanceof ApiError && error.status < 500) && count < 1

export function useUsersQuery() {
  const session = useSessionStore()
  return useQuery({
    queryKey: ['users'],
    queryFn: ({ signal }) => usersApi.list(session.token, signal),
    enabled: () => session.isAdmin,
    retry,
  })
}

export function useUserQuery(id: () => string) {
  const session = useSessionStore()
  return useQuery({
    queryKey: computed(() => ['users', id()]),
    queryFn: ({ signal }) => usersApi.detail(id(), session.token, signal),
    enabled: () => session.isAdmin,
    retry,
  })
}
