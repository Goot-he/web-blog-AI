import { request } from '@/shared/api/http'
import { requestStream } from '@/shared/api/sse'
import type {
  AiChatChatRequest,
  AiChatMessage,
  AiChatSession,
} from '@/shared/types/ai-chat'

export const aiChatApi = {
  createSession() {
    return request<AiChatSession>('/v1/ai-chat/session', { method: 'POST' })
  },

  listSessions() {
    return request<AiChatSession[]>('/v1/ai-chat/sessions', { method: 'GET' })
  },

  updateSession(id: string, title: string) {
    return request<AiChatSession>(`/v1/ai-chat/session/${encodeURIComponent(id)}`, { method: 'PUT', body: { title } })
  },

	deleteSession(id: string) {
		return request<void>(`/v1/ai-chat/session/${encodeURIComponent(id)}`, { method: 'DELETE' })
	},

  listMessages(id: string) {
    return request<AiChatMessage[]>(`/v1/ai-chat/session/${encodeURIComponent(id)}/messages`, { method: 'GET' })
  },

  chat(payload: AiChatChatRequest) {
    return requestStream('/v1/ai-chat/chat', payload)
  },
}
