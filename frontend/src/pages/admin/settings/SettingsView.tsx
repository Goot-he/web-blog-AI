import { useEffect, useState } from 'react'
import { aiSettingsApi } from '@/shared/api/v1/ai-settings'
import type { AiSettings } from '@/shared/types/ai-settings'
import { logRequestError } from '@/shared/lib/request-error'
import './SettingsView.css'

export default function SettingsView() {
  const [settings, setSettings] = useState<AiSettings | null>(null)
  const [saving, setSaving] = useState(false)
  const [switching, setSwitching] = useState(false)
  const [notice, setNotice] = useState('')

  useEffect(() => { void aiSettingsApi.get().then((result) => setSettings({ ...result, apiKey: result.apiKey || '' })).catch((error: unknown) => { logRequestError('加载助手配置失败', error); setNotice('加载配置失败，请稍后重试') }) }, [])

  const saveSettings = async () => {
    if (!settings) return
    setSaving(true)
    try {
      const saved = await aiSettingsApi.update({ provider: 'zhipu', modelId: settings.modelId.trim(), apiKey: settings.apiKey.trim() })
      setSettings(saved)
      setNotice('模型与 API Key 配置已保存。')
    } catch (error) { logRequestError('保存助手配置失败', error); setNotice('保存配置失败，请稍后重试') } finally { setSaving(false) }
  }

  const toggleAssistant = async () => {
    if (!settings || switching) return
    setSwitching(true)
    try {
      const status = await aiSettingsApi.updateStatus(!settings.assistantEnabled)
      setSettings({ ...settings, assistantEnabled: status.assistantEnabled })
      setNotice(status.assistantEnabled ? '站内助手已启用。' : '站内助手已关闭。')
    } catch (error) { logRequestError('更新助手状态失败', error); setNotice('更新状态失败，请稍后重试') } finally { setSwitching(false) }
  }

  if (!settings) return <section className="admin-card"><p>{notice || '正在加载助手配置...'}</p></section>
  return <section className="admin-card ai-settings"><div className="admin-card-head"><div><span className="eyebrow">AI ASSISTANT</span><h2>站内助手设置</h2><p>配置智谱模型与 API Key，知识库由用户在聊天界面选择。</p></div></div><div className="ai-settings-row"><div><strong>启用站内助手</strong><span>这个开关独立生效，不会修改下方配置。</span></div><button className={`settings-switch ${settings.assistantEnabled ? 'on' : ''}`} type="button" role="switch" aria-label="启用站内助手" aria-checked={settings.assistantEnabled} disabled={switching} onClick={() => void toggleAssistant()}><i /></button></div><div className="ai-settings-grid"><label htmlFor="ai-model">智谱模型<input id="ai-model" value={settings.modelId} onChange={(event) => setSettings({ ...settings, modelId: event.target.value })} placeholder="例如 glm-4-flash" autoComplete="off" /></label><label htmlFor="ai-api-key">智谱 API Key<input id="ai-api-key" value={settings.apiKey} onChange={(event) => setSettings({ ...settings, apiKey: event.target.value })} placeholder="留空表示不配置 API Key" autoComplete="off" /></label></div>{notice && <p className="settings-notice" role="status">{notice}</p>}<div className="settings-actions"><button className="primary-button" type="button" disabled={saving || !settings.modelId.trim()} onClick={() => void saveSettings()}>{saving ? '保存中...' : '保存配置'}</button></div></section>
}
