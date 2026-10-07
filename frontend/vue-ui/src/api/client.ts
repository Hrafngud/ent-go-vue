// Both production Nginx and the local Vite proxy serve this same namespace.
const API_BASE_URL = '/api'

export class ApiError extends Error {
  readonly status: number

  constructor(status: number) {
    super(`Request failed (${status})`)
    this.name = 'ApiError'
    this.status = status
  }
}

export async function get<T>(path: string, signal?: AbortSignal): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    headers: { Accept: 'application/json' },
    signal,
  })
  if (!response.ok) throw new ApiError(response.status)
  return response.json() as Promise<T>
}
