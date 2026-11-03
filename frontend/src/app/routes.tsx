import type { ReactNode } from 'react'
import AboutPage from '@/pages/about/AboutPage'
import BlogPage from '@/pages/blog/BlogPage'
import BlogDetailPage from '@/pages/blog/BlogDetailPage'
import DiaryPage from '@/pages/diary/DiaryPage'
import HomePage from '@/pages/home/HomePage'
import LoginPage from '@/pages/auth/LoginPage'
import RegisterPage from '@/pages/auth/RegisterPage'
import AssistantPage from '@/pages/assistant/AssistantPage'
import AdminPage from '@/pages/admin/AdminPage'
import UserProfilePage from '@/pages/user/UserProfilePage'

export interface AppRoute {
  path: string
  label: string
  element: ReactNode
  showInNav?: boolean
}

// 加入所有的路由，showInNav表示是否在导航栏显示

export const routes: AppRoute[] = [
  { path: '/', label: '首页', element: <HomePage />, showInNav: true },
  { path: '/diary', label: '生活小记', element: <DiaryPage />, showInNav: true },
  { path: '/blog/:id', label: '博客详情', element: <BlogDetailPage /> },
  { path: '/blog', label: '知识专栏', element: <BlogPage />, showInNav: true },
  { path: '/users/:id', label: '用户主页', element: <UserProfilePage /> },
  { path: '/assistant', label: '站内助手', element: <AssistantPage />, showInNav: true },
  { path: '/about', label: '关于我们', element: <AboutPage />, showInNav: true },
  { path: '/login', label: '登录', element: <LoginPage /> },
  { path: '/register', label: '注册', element: <RegisterPage /> },
  { path: '/admin/create', label: '后台创作内容', element: <AdminPage /> },
  { path: '/admin/creator', label: '后台创作数据看板', element: <AdminPage /> },
  { path: '/admin/drafts', label: '后台草稿箱', element: <AdminPage /> },
  { path: '/admin/profile', label: '后台个人中心', element: <AdminPage /> },
  { path: '/admin/my-articles', label: '后台我的文章', element: <AdminPage /> },
  { path: '/admin/notes', label: '后台小记管理', element: <AdminPage /> },
  { path: '/admin/blog', label: '后台文章管理', element: <AdminPage /> },
  { path: '/admin/categories', label: '后台分类管理', element: <AdminPage /> },
  { path: '/admin/users', label: '后台用户详情', element: <AdminPage /> },
  { path: '/admin/knowledge', label: '后台知识库管理', element: <AdminPage /> },
  { path: '/admin/likes', label: '后台点赞列表', element: <AdminPage /> },
  { path: '/admin/favorites', label: '后台我的收藏', element: <AdminPage /> },
  { path: '/admin/settings', label: '后台系统设置', element: <AdminPage /> },
  { path: '/admin', label: '后台管理', element: <AdminPage /> },
]
