import ProfileView from '@/pages/admin/profile/ProfileView'
import AppLink from '@/shared/ui/AppLink'

export default function UserProfilePage() {
  const userID = Number(window.location.pathname.split('/').filter(Boolean).at(-1))

  if (!Number.isFinite(userID) || userID <= 0) {
    return <main className="front-detail-state">
      <strong>没有找到这个用户</strong>
      <AppLink to="/blog">返回博客</AppLink>
    </main>
  }

  return <ProfileView userID={userID} publicView />
}
