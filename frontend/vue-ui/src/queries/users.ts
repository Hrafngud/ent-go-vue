import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { ApiError } from '../api/client'
import { usersApi, type UserInput } from '../api/users'
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

export function useSaveUserMutation() {
  const session = useSessionStore()
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, input }: { id?: string; input: UserInput }) => id
      ? usersApi.update(id, input, session.token)
      : usersApi.create(input, session.token),
    retry: false,
    onSuccess(response) {
      queryClient.setQueryData(['users', response.data.id], response)
      if (response.data.id === session.user?.id) session.user = response.data
      void queryClient.invalidateQueries({ queryKey: ['users'], exact: true })
    },
  })
}

export function useDeleteUserMutation() {
  const session = useSessionStore()
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => usersApi.delete(id, session.token),
    retry: false,
    onSuccess(_, id) {
      queryClient.removeQueries({ queryKey: ['users', id], exact: true })
      void queryClient.invalidateQueries({ queryKey: ['users'], exact: true })
    },
  })
}

export function useUserQuery(id: () => string) {
  const session = useSessionStore()
  return useQuery({
    queryKey: computed(() => ['users', id()]),
    queryFn: ({ signal }) => usersApi.detail(id(), session.token, signal),
    enabled: () => session.isAdmin && !!id(),
    retry,
  })
}
