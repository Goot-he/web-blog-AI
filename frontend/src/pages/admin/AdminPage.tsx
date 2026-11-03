import { useEffect, useState } from 'react'
import { authApi } from '@/shared/api/v1/auth'
import { navigate, usePathname } from '@/shared/hooks/usePathname'
import PublicHeader from '@/shared/ui/PublicHeader'
import AdminSidebar from '@/pages/admin/ui/AdminSidebar'
import CategoriesView from '@/pages/admin/categories/CategoriesView'
import CreateArticleView from '@/pages/admin/create/CreateArticleView'
import ArticlesView from '@/pages/admin/articles/ArticlesView'
import ArticleManagementView from '@/pages/admin/articles/ArticleManagementView'
import LikesView from '@/pages/admin/likes/LikesView'
import FavoritesView from '@/pages/admin/favorites/FavoritesView'
import ProfileView from '@/pages/admin/profile/ProfileView'
import SettingsView from '@/pages/admin/settings/SettingsView'
import KnowledgeView from '@/pages/admin/knowledge/KnowledgeView'
import UsersView from '@/pages/admin/users/UsersView'
import Dashboard from '@/pages/admin/Dashboard'
import { TAB_PATHS, TAB_TITLES, type AdminTab } from '@/pages/admin/menu'
import { clearCurrentUser, readCurrentUser } from '@/shared/lib/current-user'
import './AdminPage.css'
import './AdminTone.css'

const tabFromPath = (path: string) => (Object.keys(TAB_PATHS).find((key) => TAB_PATHS[key as AdminTab] === path) ?? 'profile') as AdminTab
const SYSTEM_ADMIN_TABS = new Set<AdminTab>(['articles', 'categories', 'users', 'knowledge', 'settings'])
const CREATOR_TABS = new Set<AdminTab>(['creatorDashboard', 'myArticles', 'drafts', 'create', 'notes'])

function NotesPlaceholder() {
  return <section className="admin-card"><div className="admin-card-head"><div><span className="eyebrow">COMING SOON</span><h2>小记管理</h2><p>生活小记管理模块暂未开发，后续会接入小记发布、草稿和数据能力。</p></div></div></section>
}

function ComingSoon({ tab, onLogout }: { tab: AdminTab; onLogout: () => void }) {
  return <section className="admin-card"><div className="admin-card-head"><div><span className="eyebrow">COMING SOON</span><h2>{TAB_TITLES[tab]}</h2><p>这个工作区正在准备中，后续会接入完整的内容与知识库管理能力。</p></div><button className="primary-button" type="button" onClick={onLogout}>退出当前账号</button></div></section>
}

export default function AdminPage() {
  const path = usePathname()
  const [tab, setTab] = useState<AdminTab>(() => tabFromPath(path))
  const currentUser = readCurrentUser() as { role?: string }

  useEffect(() => setTab(tabFromPath(path)), [path])

  const logout = () => {
    void authApi.logout().catch(() => undefined)
    clearCurrentUser()
    navigate('/login')
  }

  const changeTab = (nextTab: AdminTab) => {
    setTab(nextTab)
    navigate(TAB_PATHS[nextTab])
  }

  if (tab === 'profile') return <ProfileView />
  const sidebarMode = SYSTEM_ADMIN_TABS.has(tab) ? 'system' : CREATOR_TABS.has(tab) ? 'creator' : 'full'
  const content = tab === 'dashboard'
    ? <Dashboard />
    : tab === 'categories'
      ? <CategoriesView />
      : tab === 'articles'
        ? <ArticleManagementView />
        : tab === 'creatorDashboard'
          ? <ArticlesView mode="dashboard" />
          : tab === 'create'
            ? <CreateArticleView />
            : tab === 'drafts'
              ? <ArticlesView mode="drafts" />
              : tab === 'myArticles'
                ? <ArticlesView mode="published" />
                : tab === 'notes'
                  ? <NotesPlaceholder />
                  : tab === 'likes'
                    ? <LikesView />
                    : tab === 'favorites'
                      ? <FavoritesView />
                    : tab === 'settings'
                      ? <SettingsView />
                      : tab === 'knowledge'
                        ? <KnowledgeView />
                        : tab === 'users'
                          ? <UsersView />
                          : <ComingSoon tab={tab} onLogout={logout} />

  return (
    <>
      <PublicHeader pathname={path} />
      <main className={`admin-page admin-page-with-public-nav${tab === 'create' ? ' admin-page-compose' : ''}`}>
        <AdminSidebar tab={tab} role={currentUser.role ?? 'user'} mode={sidebarMode} onChange={changeTab} />
        <section className="admin-main">
          <div className="admin-content">
            {content}
          </div>
        </section>
      </main>
    </>
  )
}
