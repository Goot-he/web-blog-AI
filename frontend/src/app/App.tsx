import { useEffect, useState } from 'react'
import SiteLayout from '@/app/layouts/SiteLayout'
import LoginPage from '@/pages/auth/LoginPage'
import { navigate, usePathname } from '@/shared/hooks/usePathname'
import { AUTH_TOKEN_EVENT, ACCESS_TOKEN_KEY } from '@/shared/config/auth'
import { routes } from '@/app/routes'

function App() {
  const pathname = usePathname()
  const [isAuthenticated, setIsAuthenticated] = useState(() => Boolean(localStorage.getItem(ACCESS_TOKEN_KEY)))
  const route = routes.find((item) => item.path === pathname || (item.path.includes('/:') && pathname.startsWith(item.path.split('/:')[0] + '/'))) ?? routes[0]

  useEffect(() => {
    if (pathname.startsWith('/admin') && !localStorage.getItem(ACCESS_TOKEN_KEY)) {
      navigate('/login', { redirectTo: pathname })
      return
    }
    setIsAuthenticated(Boolean(localStorage.getItem(ACCESS_TOKEN_KEY)))
  }, [pathname])

  useEffect(() => {
    const syncAuth = () => setIsAuthenticated(Boolean(localStorage.getItem(ACCESS_TOKEN_KEY)))
    window.addEventListener(AUTH_TOKEN_EVENT, syncAuth)
    return () => window.removeEventListener(AUTH_TOKEN_EVENT, syncAuth)
  }, [])

  if (pathname.startsWith('/admin') && !isAuthenticated) return <LoginPage />
  if (pathname.startsWith('/admin')) return route.element
  return <SiteLayout pathname={pathname}>{route.element}</SiteLayout>
}

export default App
