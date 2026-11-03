import { useState } from 'react'
import { TAB_PATHS, ADMIN_MENU, type AdminTab } from '@/pages/admin/menu'
import { navigate } from '@/shared/hooks/usePathname'

const icons: Record<string, string> = { '个人中心': '◎', '工作台': '⌂', '创作中心': '✎', '互动反馈': '♡', '系统管理': '⚙' }

export default function AdminSidebar({ tab, role, mode = 'full', onChange }: { tab: AdminTab; role: string; mode?: 'full' | 'creator' | 'system'; onChange: (tab: AdminTab) => void }) {
  const [expanded, setExpanded] = useState<Record<string, boolean>>((): Record<string, boolean> => (
    window.matchMedia('(min-width: 901px)').matches
      ? { '创作中心': true, '系统管理': true }
      : {}
  ))
  const [childExpanded, setChildExpanded] = useState<Record<string, boolean>>(() => ({ 内容管理: true }))
  const selectTab = (nextTab: AdminTab) => { onChange(nextTab); navigate(TAB_PATHS[nextTab]) }
  const groups = ADMIN_MENU
    .filter((group) => mode === 'system' ? group.label === '系统管理' : mode === 'creator' ? group.label === '创作中心' : true)
    .filter((group) => !group.adminOnly || role === 'admin')
  const renderItem = (item: { id: AdminTab; label: string; disabled?: boolean }) => (
    <button key={item.id} type="button" disabled={item.disabled} className={`admin-menu-link ${tab === item.id ? 'active' : ''}`} onClick={() => !item.disabled && selectTab(item.id)}>{item.label}</button>
  )

  return (
    <aside className="admin-sidebar">
      <nav className="admin-menu" aria-label="后台导航">
        {groups.map((group) => {
          const open = mode === 'system' || mode === 'creator' || (expanded[group.label] ?? false)
          const direct = !group.collapsible
          return <div className="admin-menu-group" key={group.label}>
            {direct ? <button className={`admin-menu-link ${tab === group.items[0].id ? 'active' : ''}`} type="button" onClick={() => selectTab(group.items[0].id)}><span className="menu-icon">{icons[group.label] ?? '•'}</span>{group.label}</button> : <button className="admin-menu-heading" type="button" onClick={() => mode === 'system' || mode === 'creator' ? undefined : setExpanded((current) => ({ ...current, [group.label]: !open }))}><span className="menu-icon">{icons[group.label] ?? '•'}</span><span>{group.label}</span>{mode === 'full' && <span className={open ? 'admin-menu-arrow expanded' : 'admin-menu-arrow'}>›</span>}</button>}
            {group.collapsible && open && <div className="admin-menu-items">{group.items.map(renderItem)}{group.children?.map((child) => {
              const childOpen = childExpanded[child.label] ?? false
              return <div className="admin-submenu" key={child.label}>
                <button className="admin-submenu-heading" type="button" onClick={() => setChildExpanded((current) => ({ ...current, [child.label]: !childOpen }))}><span>{child.label}</span><i className={childOpen ? 'expanded' : ''}>›</i></button>
                {childOpen && child.items.map(renderItem)}
              </div>
            })}</div>}
          </div>
        })}
      </nav>
      <div className="admin-sidebar-footer"><div className="footer-line"><span className="status-dot" />系统运行正常</div><small>CCEBD · 2026</small></div>
    </aside>
  )
}
