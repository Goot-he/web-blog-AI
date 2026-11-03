import { env } from '@/shared/config/env'
import { ACCESS_TOKEN_KEY, AUTH_TOKEN_EVENT } from '@/shared/config/auth'
import type { ApiResponse } from '@/shared/types/api'

export interface RequestOptions {
  method?: string
  body?: unknown
  auth?: boolean
  token?: string | null
}

export interface UploadOptions {
  method?: string
}

export function getAccessToken() {
  return localStorage.getItem(ACCESS_TOKEN_KEY)
}

function notifyAuthTokenUpdated() {
  window.dispatchEvent(new Event(AUTH_TOKEN_EVENT))
}

export function setAccessToken(token?: string | null) {
  const normalized = token?.trim()
  if (!normalized) {
    clearAccessToken()
    return
  }
  localStorage.setItem(ACCESS_TOKEN_KEY, normalized)
  notifyAuthTokenUpdated()
}

export function clearAccessToken() {
  localStorage.removeItem(ACCESS_TOKEN_KEY)
  notifyAuthTokenUpdated()
}

function headers(auth: boolean, tokenOverride?: string | null) {
  const result = new Headers({ Accept: 'application/json', 'Content-Type': 'application/json' })
  const token = auth ? (tokenOverride === undefined ? getAccessToken() : tokenOverride)?.trim() : null
  if (token) result.set('Authorization', `Bearer ${token}`)
  return result
}

function unwrap<T>(payload: unknown): T {
  if (typeof payload === 'object' && payload !== null && 'code' in payload && 'msg' in payload) {
    const response = payload as ApiResponse<T>
    if (response.code !== 200) throw new Error(response.msg || '请求失败')
    return response.data
  }
  return payload as T
}

async function parseResponse(response: Response): Promise<unknown> {
  const content = await response.text()
  if (!content.trim()) return undefined
  try {
    return JSON.parse(content) as unknown
  } catch {
    throw new Error(content.trim() || `请求失败 (${response.status})`)
  }
}

export async function request<T>(path: string, options: RequestOptions = {}) {
  const response = await fetch(`${env.apiBaseUrl}${path}`, {
    method: options.method ?? 'GET',
    headers: headers(options.auth !== false, options.token),
    body: options.body === undefined ? undefined : JSON.stringify(options.body),
    // refresh_token 存在 HttpOnly Cookie 中，所有 API 请求都必须携带它，
    // 否则 access token 过期时后端无法执行自动刷新。
    credentials: 'include',
  })
  const refreshedAuthorization = response.headers.get('Authorization')
  if (refreshedAuthorization?.startsWith('Bearer ')) {
    setAccessToken(refreshedAuthorization.slice('Bearer '.length))
  }
  if (response.status === 204) return undefined as T
  const payload = await parseResponse(response)
  if (!response.ok) throw new Error((payload as { msg?: string }).msg || `请求失败 (${response.status})`)
  return unwrap<T>(payload)
}

export async function upload<T>(path: string, file: File, fields: Record<string, string> = {}, options: UploadOptions = {}) {
  const form = new FormData()
  form.set('file', file)
  Object.entries(fields).forEach(([key, value]) => form.set(key, value))
  const token = getAccessToken()
  const response = await fetch(`${env.apiBaseUrl}${path}`, {
    method: options.method ?? 'POST',
    headers: token ? { Authorization: `Bearer ${token}` } : undefined,
    body: form,
    credentials: 'include',
  })
  const payload = await parseResponse(response)
  if (!response.ok) throw new Error((payload as { msg?: string }).msg || `请求失败 (${response.status})`)
  return unwrap<T>(payload)
}
