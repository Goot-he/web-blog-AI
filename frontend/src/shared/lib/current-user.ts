import { CURRENT_USER_KEY } from '@/shared/config/auth'
import type { CurrentUser } from '@/shared/types/user'

export type StoredCurrentUser = Partial<CurrentUser>

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value)
}

export function readCurrentUser(): StoredCurrentUser {
  const raw = localStorage.getItem(CURRENT_USER_KEY)
  if (!raw) return {}

  try {
    const parsed = JSON.parse(raw) as unknown
    return isRecord(parsed) ? parsed : {}
  } catch {
    localStorage.removeItem(CURRENT_USER_KEY)
    return {}
  }
}

export function saveCurrentUser(user?: StoredCurrentUser | null) {
  if (!user || !isRecord(user)) {
    localStorage.removeItem(CURRENT_USER_KEY)
    return
  }
  localStorage.setItem(CURRENT_USER_KEY, JSON.stringify(user))
  window.dispatchEvent(new Event('renai:user-updated'))
}

export function clearCurrentUser() {
  localStorage.removeItem(CURRENT_USER_KEY)
  window.dispatchEvent(new Event('renai:user-updated'))
}
