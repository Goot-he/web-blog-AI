import { request } from '@/shared/api/http'
import type { AiChatStatus, AiSettings, UpdateAiSettingsRequest } from '@/shared/types/ai-settings'
export const aiSettingsApi = { status: () => request<AiChatStatus>('/v1/ai-chat/status'), updateStatus: (assistantEnabled: boolean) => request<AiChatStatus>('/v1/ai-chat/status', { method: 'PUT', body: { assistantEnabled } }), get: () => request<AiSettings>('/v1/ai-chat/settings'), update: (payload: UpdateAiSettingsRequest) => request<AiSettings>('/v1/ai-chat/settings', { method: 'PUT', body: payload }) }
