export type AdminTab = 'profile' | 'creatorDashboard' | 'create' | 'drafts' | 'myArticles' | 'notes' | 'likes' | 'favorites' | 'dashboard' | 'articles' | 'categories' | 'users' | 'knowledge' | 'settings'

export const TAB_PATHS: Record<AdminTab, string> = { profile: '/admin/profile', creatorDashboard: '/admin/creator', create: '/admin/create', drafts: '/admin/drafts', myArticles: '/admin/my-articles', notes: '/admin/notes', likes: '/admin/likes', favorites: '/admin/favorites', dashboard: '/admin', articles: '/admin/blog', categories: '/admin/categories', users: '/admin/users', knowledge: '/admin/knowledge', settings: '/admin/settings' }
export const TAB_TITLES: Record<AdminTab, string> = { profile: '个人中心', creatorDashboard: '数据看板', create: '投稿', drafts: '草稿箱', myArticles: '我的文章', notes: '小记管理', likes: '点赞列表', favorites: '我的收藏', dashboard: '工作台', articles: '文章管理', categories: '分类管理', users: '用户管理', knowledge: '知识库管理', settings: '系统设置' }

export interface AdminMenuItem { id: AdminTab; label: string; disabled?: boolean }
export interface AdminMenuGroup { label: string; items: AdminMenuItem[]; children?: { label: string; items: AdminMenuItem[] }[]; adminOnly?: boolean; collapsible?: boolean }
export const ADMIN_MENU: AdminMenuGroup[] = [
  { label: '个人中心', items: [{ id: 'profile', label: '个人中心' }] },
  { label: '工作台', items: [{ id: 'dashboard', label: '工作台' }] },
  {
    label: '创作中心',
    collapsible: true,
    items: [{ id: 'creatorDashboard', label: '数据看板' }],
    children: [
      { label: '内容管理', items: [{ id: 'create', label: '投稿' }, { id: 'drafts', label: '草稿箱' }, { id: 'myArticles', label: '我的文章' }] },
      { label: '小记管理', items: [{ id: 'notes', label: '暂未开发', disabled: true }] },
    ],
  },
  { label: '互动反馈', items: [{ id: 'likes', label: '点赞列表' }, { id: 'favorites', label: '我的收藏' }] },
  { label: '系统管理', adminOnly: true, collapsible: true, items: [{ id: 'articles', label: '文章管理' }, { id: 'categories', label: '分类管理' }, { id: 'users', label: '用户管理' }, { id: 'knowledge', label: '知识库管理' }, { id: 'settings', label: '系统设置' }] },
]
