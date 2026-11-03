import type { VisitorProfile } from './user'
export interface FollowListResponse { users: VisitorProfile[]; total: number }
export interface FollowStatusResponse { isFollowing: boolean }
