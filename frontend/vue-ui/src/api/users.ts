import { get, post, put, remove } from './client'

export interface User {
  id: string
  name: string
  email: string
  created_at: string
  is_admin: boolean
}

export interface UserInput {
  name: string
  email: string
  password?: string
}

interface UserResponse { data: User }
const userPath = (id: string) => `/admin/users/${encodeURIComponent(id)}`

export const usersApi = {
  list: (token: string, signal?: AbortSignal) => get<{ data: User[] }>('/admin/users', signal, token),
  detail: (id: string, token: string, signal?: AbortSignal) => get<UserResponse>(userPath(id), signal, token),
  create: (input: UserInput, token: string) => post<UserResponse>('/admin/users', input, token),
  update: (id: string, input: UserInput, token: string) => put<UserResponse>(userPath(id), input, token),
  delete: (id: string, token: string) => remove(userPath(id), token),
}

export function formatDate(value: string): string {
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? 'Unknown' : new Intl.DateTimeFormat(undefined, { dateStyle: 'medium' }).format(date)
}
