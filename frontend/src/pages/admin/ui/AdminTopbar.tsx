import { useEffect, useState } from 'react'
import { TAB_TITLES, type AdminTab } from '@/pages/admin/menu'
import { resolveStorageUrl } from '@/shared/lib/storage'
import { readCurrentUser } from '@/shared/lib/current-user'

type TopbarUser = { username?: string; role?: string; avatar?: string }

export default function AdminTopbar({ tab, section = '管理后台' }: { tab: AdminTab; section?: string }) {
  const [currentUser, setCurrentUser] = useState<TopbarUser>(() => readCurrentUser() as TopbarUser)
  const avatar = resolveStorageUrl(currentUser.avatar ?? '')

  useEffect(() => {
    const updateUser = () => setCurrentUser(readCurrentUser() as TopbarUser)
    window.addEventListener('renai:user-updated', updateUser)
    return () => window.removeEventListener('renai:user-updated', updateUser)
  }, [])

  return <header className="admin-top"><div className="admin-top-title"><span className="crumb">{section} <b>/</b> {TAB_TITLES[tab]}</span><h1>{TAB_TITLES[tab]}</h1></div><div className="admin-user"><div className="admin-avatar">{avatar ? <img src={avatar} alt="用户头像" /> : (currentUser.username || 'U').slice(0, 1).toUpperCase()}</div><div className="admin-user-name"><b>{currentUser.username || '管理员'}</b><small>{currentUser.role === 'admin' ? '超级管理员' : '社团成员'}</small></div></div></header>
}
