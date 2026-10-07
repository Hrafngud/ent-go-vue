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

export async function get<T>(path: string, signal?: AbortSignal, token?: string): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    headers: { Accept: 'application/json', ...(token ? { Authorization: `Bearer ${token}` } : {}) },
    signal,
  })
  if (!response.ok) throw new ApiError(response.status)
  return response.json() as Promise<T>
}

export async function post<T>(path: string, body: unknown): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    method: 'POST',
    headers: { Accept: 'application/json', 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  if (!response.ok) throw new ApiError(response.status)
  return response.json() as Promise<T>
}
