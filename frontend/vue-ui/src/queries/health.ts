import { useQuery } from '@tanstack/vue-query'
import { get } from '../api/client'

interface Health {
  status: 'ok'
}

interface Readiness extends Health {
  database: 'connected'
}

const healthOptions = {
  refetchInterval: 30_000,
  staleTime: 10_000,
  retry: 1,
}

export function useHealthQuery() {
  return useQuery({
    ...healthOptions,
    queryKey: ['health'],
    queryFn: ({ signal }) => get<Health>('/health', signal),
  })
}

export function useReadinessQuery() {
  return useQuery({
    ...healthOptions,
    queryKey: ['readiness'],
    queryFn: ({ signal }) => get<Readiness>('/ready', signal),
  })
}
