import { useEffect, useState, type ReactNode } from 'react'
import AppLink from '@/shared/ui/AppLink'
import PublicHeader from '@/shared/ui/PublicHeader'

export default function SiteLayout({ children, pathname }: { children: ReactNode; pathname: string }) {
  const isHome = pathname === '/'
  const isUserProfile = pathname.startsWith('/users/')
  const hideFooter = pathname === '/assistant' || pathname === '/diary'
  const [hasScrolled, setHasScrolled] = useState(false)

  useEffect(() => {
    if (!isHome) {
      setHasScrolled(false)
      return
    }
    const updateHeader = () => setHasScrolled(window.scrollY > 32)
    updateHeader()
    window.addEventListener('scroll', updateHeader, { passive: true })
    return () => window.removeEventListener('scroll', updateHeader)
  }, [isHome])

  return <div className={`site-shell${isHome ? ' site-shell-home' : ''}${pathname === '/diary' ? ' site-shell-diary' : ''}`}>
    {!isUserProfile && <PublicHeader pathname={pathname} overlay={isHome} scrolled={hasScrolled} />}
    {children}
    {!isUserProfile && !hideFooter && <footer className="site-footer">
      <div className="site-footer-brand"><b>CCEBD</b><p>一个记录社团生活、创作和灵感的交流空间。<br />在这里，分享值得被认真看见。</p><a className="site-footer-record" href="https://beian.miit.gov.cn/" target="_blank" rel="noreferrer">浙ICP备2025156531号-1</a></div>
      <nav className="site-footer-links" aria-label="快速导航"><strong>快速导航</strong><AppLink to="/">首页</AppLink><AppLink to="/blog">知识专栏</AppLink><AppLink to="/diary">生活小记</AppLink><AppLink to="/about">关于我们</AppLink></nav>
      <div className="site-footer-support"><strong>支持</strong><AppLink to="/assistant">站内助手</AppLink></div>
    </footer>}
  </div>
}
