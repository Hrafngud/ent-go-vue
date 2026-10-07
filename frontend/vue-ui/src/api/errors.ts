import { ApiError } from './client'

export function userError(error: unknown, fallback: string): string {
  if (error instanceof ApiError) {
    if (error.status === 409) return 'This email address is already in use. Choose another address.'
    if (error.status === 403) return 'This action is not allowed. User management is reserved for the root account.'
    if (error.status === 404) return 'This user no longer exists. Return to the users list.'
    if (error.status === 422 || error.status === 400) return 'Check your name, email, and password, then try again.'
    if (error.status === 401) return 'Your session has ended. Please sign in again.'
  }
  return fallback
}
