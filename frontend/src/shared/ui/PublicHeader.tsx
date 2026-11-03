import { useEffect, useState, type FormEvent, type ReactNode } from 'react'
import { authApi } from '@/shared/api/v1/auth'
import { ACCESS_TOKEN_KEY, AUTH_TOKEN_EVENT } from '@/shared/config/auth'
import { clearCurrentUser, readCurrentUser } from '@/shared/lib/current-user'
import { navigate } from '@/shared/hooks/usePathname'
import { resolveStorageUrl } from '@/shared/lib/storage'
import AppLink from '@/shared/ui/AppLink'
import './PublicHeader.css'

type CurrentUser = { username?: string; avatar?: string; role?: string }
type IconName = 'home' | 'diary' | 'blog' | 'assistant' | 'about' | 'message' | 'favorite' | 'creator' | 'search'

function PublicIcon({ name }: { name: IconName }) {
  const paths: Record<IconName, ReactNode> = {
    home: <><path d="m3 10 9-7 9 7v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" /><path d="M9 21v-6h6v6" /></>,
    diary: <><rect x="5" y="3" width="14" height="18" rx="2" /><path d="M8 7h8M8 11h8M8 15h5" /></>,
    blog: <><path d="M4 19.5V5a2 2 0 0 1 2-2h12v18H6a2 2 0 0 1-2-1.5Z" /><path d="M8 7h7M8 11h7" /></>,
    assistant: <><path d="M12 3a7 7 0 0 0-7 7v3a3 3 0 0 0 3 3h1v-5H6" /><path d="M12 3a7 7 0 0 1 7 7v3a3 3 0 0 1-3 3h-1v-5h3" /><path d="M10 20h4" /></>,
    about: <><circle cx="12" cy="12" r="9" /><path d="M12 10v6M12 7h.01" /></>,
    message: <><rect x="3" y="5" width="18" height="14" rx="2" /><path d="m3 7 9 6 9-6" /></>,
    favorite: <path d="m12 3 2.7 5.6 6.1.9-4.4 4.3 1 6.1L12 17l-5.4 2.9 1-6.1-4.4-4.3 6.1-.9L12 3Z" />,
    creator: <path d="M12 5v14M5 12h14" />,
    search: <><circle cx="11" cy="11" r="7" /><path d="m20 20-3.5-3.5" /></>,
  }
  return <svg className="public-nav-icon" viewBox="0 0 24 24" aria-hidden="true" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">{paths[name]}</svg>
}

function AccountChevron() {
  return <svg className="nav-account-chevron" viewBox="0 0 24 24" aria-hidden="true" fill="none" stroke="currentColor" strokeWidth="2.8" strokeLinecap="round" strokeLinejoin="round"><path d="m9 18 6-6-6-6" /></svg>
}

const navItems = [
  { to: '/', label: '首页', icon: 'home' as const },
  { to: '/blog', label: '知识专栏', icon: 'blog' as const },
  { to: '/diary', label: '生活小记', icon: 'diary' as const },
  { to: '/assistant', label: '站内助手', icon: 'assistant' as const },
]

export default function PublicHeader({ pathname, overlay = false, scrolled = false }: { pathname: string; overlay?: boolean; scrolled?: boolean }) {
  const [keyword, setKeyword] = useState('')
  const [tip, setTip] = useState('')
  const [loggedIn, setLoggedIn] = useState(() => Boolean(localStorage.getItem(ACCESS_TOKEN_KEY)))
  const currentUser = readCurrentUser() as CurrentUser
  const avatar = resolveStorageUrl(currentUser.avatar ?? '')
  const initial = (currentUser.username || 'U').slice(0, 1).toUpperCase()
  const transparent = overlay && !scrolled
  const isActive = (path: string) => pathname === path || (path !== '/' && pathname.startsWith(`${path}/`))

  useEffect(() => {
    if (!tip) return
    const timer = window.setTimeout(() => setTip(''), 1600)
    return () => window.clearTimeout(timer)
  }, [tip])

  useEffect(() => {
    const syncAuth = () => setLoggedIn(Boolean(localStorage.getItem(ACCESS_TOKEN_KEY)))
    window.addEventListener(AUTH_TOKEN_EVENT, syncAuth)
    return () => window.removeEventListener(AUTH_TOKEN_EVENT, syncAuth)
  }, [])

  const submitSearch = (event: FormEvent) => {
    event.preventDefault()
    const q = keyword.trim()
    navigate(q ? `/blog?q=${encodeURIComponent(q)}` : '/blog')
    window.dispatchEvent(new CustomEvent('renai:blog-search', { detail: q }))
  }

  const logout = async () => {
    try {
      await authApi.logout()
    } catch {
      // authApi.logout clears the local token even if the server already rejected it.
    }
    clearCurrentUser()
    navigate('/login')
  }

  return <header className={`public-header${transparent ? ' transparent' : ''}${scrolled ? ' is-scrolled' : ''}`}>
    <AppLink className="public-brand" to="/"><img src="/ccebd.png" alt="CCEBD" /></AppLink>
    <nav className="public-nav" aria-label="主导航">{navItems.map((item) => <AppLink key={item.to} className={isActive(item.to) ? 'active' : ''} to={item.to}><PublicIcon name={item.icon} /><span>{item.label}</span></AppLink>)}</nav>
    <form className="public-search" role="search" onSubmit={submitSearch}><PublicIcon name="search" /><input value={keyword} onChange={(event) => setKeyword(event.target.value)} placeholder="搜索技术文章、笔记、源码" aria-label="搜索博客" /><button type="submit">搜索</button></form>
    <div className="public-actions">
      {loggedIn ? <div className="nav-account public-account"><AppLink className="nav-avatar" to="/admin/profile" aria-label="个人中心">{avatar ? <img src={avatar} alt="" /> : initial}</AppLink><div className="nav-account-popover" role="menu"><div className="nav-account-identity"><strong>{currentUser.username || 'Renai 用户'}</strong><span>个人账户</span></div><div className="nav-account-actions"><AppLink role="menuitem" to="/admin/profile">个人中心 <AccountChevron /></AppLink><AppLink role="menuitem" to="/admin/my-articles">内容管理 <AccountChevron /></AppLink>{currentUser.role === 'admin' && <AppLink role="menuitem" to="/admin/blog">管理后台 <AccountChevron /></AppLink>}</div><button role="menuitem" type="button" onClick={logout}>退出登录 <AccountChevron /></button></div></div> : <AppLink className="public-login" to="/login">登录</AppLink>}
      <AppLink className={`public-action${isActive('/about') ? ' active' : ''}`} to="/about" aria-label="关于"><PublicIcon name="about" /><span>关于</span></AppLink>
      <button type="button" className="public-action" onClick={() => setTip('消息暂未开发')} aria-label="消息"><PublicIcon name="message" /><span>消息</span></button>
      <button type="button" className="public-action" onClick={() => loggedIn ? navigate('/admin/profile', { profileTab: 'favorites' }) : navigate('/login')} aria-label="收藏"><PublicIcon name="favorite" /><span>收藏</span></button>
      <button type="button" className="public-action public-creator" onClick={() => loggedIn ? navigate('/admin/creator') : navigate('/login')}><span>创作中心</span></button>
    </div>
    {tip && <div className="public-tip" role="status">{tip}</div>}
  </header>
}
