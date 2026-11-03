import { request, upload } from '@/shared/api/http'
import type { DeleteUploadRequest, UploadResponse } from '@/shared/types/storage'

const auth = (token?: string) => (token ? { token } : {})

export const storageApi = {
  uploadAvatar(file: File) {
    return upload<UploadResponse>('/v1/storage/avatar', file)
  },

  uploadCover(file: File) {
    return upload<UploadResponse>('/v1/storage/cover', file)
  },

  uploadContent(file: File) {
    return upload<UploadResponse>('/v1/storage/content', file)
  },

  delete(payload: DeleteUploadRequest, token?: string) {
    return request<null>('/v1/storage/delete', { method: 'POST', body: payload, ...auth(token) })
  },
}
