// Both production Nginx and the local Vite proxy serve this same namespace.
const API_BASE_URL = '/api'

export class ApiError extends Error {
  readonly status: number
  readonly retryAfter: number

  constructor(status: number, retryAfter = 60) {
    super(`Request failed (${status})`)
    this.name = 'ApiError'
    this.status = status
    this.retryAfter = Number.isFinite(retryAfter) && retryAfter > 0 ? Math.min(3600, Math.ceil(retryAfter)) : 60
  }
}

async function request<T>(path: string, method: string, body?: unknown, token?: string, signal?: AbortSignal): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    method,
    headers: {
      Accept: 'application/json',
      ...(body !== undefined ? { 'Content-Type': 'application/json' } : {}),
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    body: body !== undefined ? JSON.stringify(body) : undefined,
    signal,
  })
  if (!response.ok) {
    if (token && response.status === 401) window.dispatchEvent(new Event('session-expired'))
    throw new ApiError(response.status, Number(response.headers.get('Retry-After')))
  }
  // Registration returns an empty 201 response and deletion returns 204.
  const text = await response.text()
  return (text ? JSON.parse(text) : undefined) as T
}

export function get<T>(path: string, signal?: AbortSignal, token?: string): Promise<T> {
  return request<T>(path, 'GET', undefined, token, signal)
}

export function post<T>(path: string, body: unknown, token?: string): Promise<T> {
  return request<T>(path, 'POST', body, token)
}

export function put<T>(path: string, body: unknown, token: string): Promise<T> {
  return request<T>(path, 'PUT', body, token)
}

export function remove(path: string, token: string): Promise<void> {
  return request<void>(path, 'DELETE', undefined, token)
}
