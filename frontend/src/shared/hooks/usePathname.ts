import { useEffect, useState } from 'react'

export function navigate(path: string, state: Record<string, unknown> = {}) {
  if (window.location.pathname === path) {
    if (Object.keys(state).length) {
      window.history.replaceState(state, '', path)
      window.dispatchEvent(new PopStateEvent('popstate'))
    }
    return
  }
  window.history.pushState(state, '', path)
  window.dispatchEvent(new PopStateEvent('popstate'))
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

export function usePathname() {
  const [pathname, setPathname] = useState(window.location.pathname)

  useEffect(() => {
    const update = () => setPathname(window.location.pathname)
    window.addEventListener('popstate', update)
    return () => window.removeEventListener('popstate', update)
  }, [])

  return pathname
}