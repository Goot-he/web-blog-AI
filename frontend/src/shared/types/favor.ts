import type { Article } from './article'

export interface FavorRequest { objectId: number; folderId?: number; userId?: number }
export interface UserFavorListRequest { page?: number; pageSize?: number; folderId?: number; userId?: number }
export interface FavorBoolResponse { success: boolean }
export interface UserFavorListResponse { articles: Article[]; total: number }
export interface FavoriteFolder { id: number; name: string; isDefault: boolean; articleCount: number }
export interface FavoriteFolderListResponse { folders: FavoriteFolder[] }
