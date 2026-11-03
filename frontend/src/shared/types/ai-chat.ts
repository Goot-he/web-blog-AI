export interface AiChatSession {
  session_id: string
  title: string
}

export interface AiChatMessage {
  session_id: string
  role: string
  content: string
  created_at: string
}

export interface AiChatChatRequest {
  session_id: string
  content: string
  collection_name: string
}

export interface AiChatStreamChunk {
  session_id: string
  content: string
  done: boolean
  knowledge?: string[]
  error?: string
}
