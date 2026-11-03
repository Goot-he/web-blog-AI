import { request } from '@/shared/api/http'
import type { FollowListResponse, FollowStatusResponse } from '@/shared/types/follow'
export const followApi = {
 follow(targetUserId:number){return request<{success:boolean}>('/v1/follow',{method:'POST',body:{targetUserId}})},
 unfollow(targetUserId:number){return request<{success:boolean}>('/v1/follow/cancel',{method:'POST',body:{targetUserId}})},
 status(targetUserId:number){return request<FollowStatusResponse>(`/v1/follow/status/${targetUserId}`)},
 following(userId:number,page=1,pageSize=10){return request<FollowListResponse>(`/v1/follow/following/${userId}?page=${page}&pageSize=${pageSize}`)},
 followers(userId:number,page=1,pageSize=10){return request<FollowListResponse>(`/v1/follow/followers/${userId}?page=${page}&pageSize=${pageSize}`)},
}
