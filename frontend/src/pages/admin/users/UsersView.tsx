import { useEffect, useState } from 'react'
import { userApi } from '@/shared/api/v1/user'
import type { ManagedUser } from '@/shared/types/user'
import { logRequestError } from '@/shared/lib/request-error'
import { resolveStorageUrl } from '@/shared/lib/storage'
import './UsersView.css'

const statusLabel: Record<ManagedUser['status'], string> = { approved: '正常', blocked: '已拉黑' }

export default function UsersView() {
  const [users, setUsers] = useState<ManagedUser[]>([])
  const [keyword, setKeyword] = useState('')
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)
  const [updatingID, setUpdatingID] = useState<number | null>(null)
  const [notice, setNotice] = useState('')

  const loadUsers = async (search = keyword) => {
    setLoading(true)
    try {
      const result = await userApi.list(search)
      if (!Array.isArray(result?.users)) throw new Error('invalid user list response')
      setUsers(result.users)
      setTotal(Number.isFinite(result.total) ? result.total : result.users.length)
      setNotice('')
    } catch (error) {
      logRequestError('加载用户列表失败', error)
      setNotice('用户列表加载失败，请稍后重试。')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { void loadUsers('') }, [])

  const updateBlacklist = async (user: ManagedUser, blacklisted: boolean) => {
    if ((user.status === 'blocked') === blacklisted) return
    setUpdatingID(user.id)
    try {
      await userApi.updateBlacklist(user.id, blacklisted)
      const status = blacklisted ? 'blocked' : 'approved'
      setUsers((current) => current.map((item) => item.id === user.id ? { ...item, status } : item))
      setNotice(blacklisted ? '用户已加入黑名单。' : '用户已移出黑名单。')
    } catch (error) {
      logRequestError('更新用户黑名单失败', error)
      setNotice('状态更新失败，请稍后重试。')
    } finally {
      setUpdatingID(null)
    }
  }

  return <section className="users-view">
    <header className="users-header"><div><span className="eyebrow">USER MANAGEMENT</span><h2>用户管理</h2><p>查看账号状态，并管理成员的访问权限。</p></div><strong>{total} 位用户</strong></header>
    <form className="users-filter" onSubmit={(event) => { event.preventDefault(); void loadUsers() }}><input value={keyword} onChange={(event) => setKeyword(event.target.value)} placeholder="搜索用户名或邮箱" aria-label="搜索用户名或邮箱" /><button className="primary-button" type="submit" disabled={loading}>搜索</button></form>
    {notice && <p className="settings-notice" role="status">{notice}</p>}
    <div className="users-table-wrap"><table className="users-table"><thead><tr><th>用户</th><th>角色</th><th>状态</th><th>手机</th><th>黑名单</th></tr></thead><tbody>{loading ? <tr><td colSpan={5} className="users-state">正在加载用户...</td></tr> : users.length === 0 ? <tr><td colSpan={5} className="users-state">未找到符合条件的用户。</td></tr> : users.map((user) => <tr key={user.id}>
      <td><div className="user-identity"><div className="user-avatar">{user.avatar ? <img src={resolveStorageUrl(user.avatar)} alt="" /> : user.username.slice(0, 1).toUpperCase()}</div><div><strong>{user.username}</strong><small>{user.email || '未填写邮箱'}</small></div></div></td><td>{user.role === 'admin' ? '管理员' : '普通用户'}</td><td><span className={`user-status ${user.status}`}>{statusLabel[user.status]}</span></td><td>{user.phone || '未填写'}</td><td>{user.role === 'admin' ? <span className="user-protected">受保护</span> : <label className="blacklist-switch" title={user.status === 'blocked' ? '移出黑名单' : '加入黑名单'}><input type="checkbox" checked={user.status === 'blocked'} disabled={updatingID === user.id} onChange={(event) => void updateBlacklist(user, event.target.checked)} aria-label={`切换 ${user.username} 的黑名单状态`} /><span aria-hidden="true" /></label>}</td>
    </tr>)}</tbody></table></div>
  </section>
}
