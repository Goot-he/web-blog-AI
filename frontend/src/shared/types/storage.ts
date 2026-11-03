export interface UploadResponse { url: string; key: string; provider: string }
export interface DeleteUploadRequest { key: string }

export type UploadKind = 'avatar' | 'cover' | 'content'
