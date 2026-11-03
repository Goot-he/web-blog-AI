export interface AiChatStatus { assistantEnabled: boolean }
export interface AiSettings { assistantEnabled: boolean; provider: string; modelId: string; apiKey: string; apiKeyConfigured: boolean }
export interface UpdateAiSettingsRequest { provider: string; modelId: string; apiKey: string }
