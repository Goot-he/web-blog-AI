import { useEffect, useRef, useState, type MouseEvent } from 'react'
import { aiChatApi } from '@/shared/api/v1/ai-chat'
import { aiSettingsApi } from '@/shared/api/v1/ai-settings'
import { vectorApi } from '@/shared/api/v1/vector'
import { navigate } from '@/shared/hooks/usePathname'
import { ACCESS_TOKEN_KEY, CURRENT_USER_KEY } from '@/shared/config/auth'
import { logRequestError } from '@/shared/lib/request-error'
import type { AiChatMessage, AiChatSession, AiChatStreamChunk } from '@/shared/types/ai-chat'
import type { VectorCollection } from '@/shared/types/vector'
import KnowledgePicker from './KnowledgePicker'
import './AssistantPage.css'
import './AssistantDeleteModal.css'
import './AssistantWorkspaceEnhancements.css'
import './AssistantComposerCard.css'
import './AssistantKnowledge.css'

type ChatMessage = AiChatMessage & { localID: string }
const COLLECTION_STORAGE_KEY = 'renai_ai_collection'

const suggestions = ['帮我推荐一篇本站的文章', '这篇文章主要讲了什么？', '从作者知识库中查找相关内容']

function toLocalMessage(message: AiChatMessage, index: number): ChatMessage {
  return { ...message, localID: `${message.created_at}-${index}` }
}

export default function AssistantPage() {
  const [sessions, setSessions] = useState<AiChatSession[]>([])
  const [activeSessionID, setActiveSessionID] = useState('')
  const [messages, setMessages] = useState<ChatMessage[]>([])
  const [draft, setDraft] = useState('')
  const [loading, setLoading] = useState(true)
  const [sending, setSending] = useState(false)
  const [notice, setNotice] = useState('')
  const [pendingDeleteSessionID, setPendingDeleteSessionID] = useState<string | null>(null)
  const [deletingSessionID, setDeletingSessionID] = useState('')
  const [editingSessionID, setEditingSessionID] = useState('')
  const [sessionTitleDraft, setSessionTitleDraft] = useState('')
  const [renamingSessionID, setRenamingSessionID] = useState('')
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false)
  const [assistantEnabled, setAssistantEnabled] = useState(true)
  const [collections, setCollections] = useState<VectorCollection[]>([])
  const [collectionName, setCollectionName] = useState('')
  const [loginRequired, setLoginRequired] = useState(() => !localStorage.getItem(ACCESS_TOKEN_KEY))
  const messageEnd = useRef<HTMLDivElement>(null)
  const loggedIn = !loginRequired
  const isEmpty = !loading && messages.length === 0

  const requireLogin = () => {
    setLoginRequired(true)
    setNotice('')
  }

  const handleRequestError = (error: unknown) => {
    const message = (error as Error).message || ''
    logRequestError('站内助手请求失败', error)
    if (/invalid or expired token|unauthorized|authentication required/i.test(message)) {
      localStorage.removeItem(ACCESS_TOKEN_KEY)
      localStorage.removeItem(CURRENT_USER_KEY)
      setLoginRequired(true)
      setNotice('')
      return
    }
    setNotice('操作未完成，请稍后重试')
  }

  const loadMessages = async (sessionID: string) => {
    const result = await aiChatApi.listMessages(sessionID)
    setMessages(result.map(toLocalMessage))
  }

  const selectSession = async (sessionID: string) => {
    if (sessionID === activeSessionID || sending) return
    setActiveSessionID(sessionID)
    setMessages([])
    try {
      await loadMessages(sessionID)
    } catch (error) {
      handleRequestError(error)
    }
  }

  const createSession = async () => {
    const session = await aiChatApi.createSession()
    setSessions((current) => [session, ...current])
    setActiveSessionID(session.session_id)
    setMessages([])
    return session
  }

  useEffect(() => {
    if (!loggedIn) {
      setLoading(false)
    } else {
      let cancelled = false
      void aiChatApi.listSessions()
        .then(async (result) => {
          if (cancelled) return
          setSessions(result)
          if (result[0]) {
            setActiveSessionID(result[0].session_id)
            const history = await aiChatApi.listMessages(result[0].session_id)
            if (!cancelled) setMessages(history.map(toLocalMessage))
          }
        })
        .catch((error: Error) => !cancelled && handleRequestError(error))
        .finally(() => !cancelled && setLoading(false))
      return () => { cancelled = true }
    }
  }, [loggedIn])

  useEffect(() => {
    void Promise.all([aiSettingsApi.status(), vectorApi.collections()]).then(([settings, list]) => {
      setAssistantEnabled(settings.assistantEnabled)
      setCollections(list)
      const saved = localStorage.getItem(COLLECTION_STORAGE_KEY) || ''
      const selected = list.some((collection) => collection.name === saved) ? saved : list[0]?.name || ''
      setCollectionName(selected)
      if (selected) localStorage.setItem(COLLECTION_STORAGE_KEY, selected)
    }).catch((error: Error) => handleRequestError(error))
  }, [loggedIn])

  useEffect(() => {
    messageEnd.current?.scrollIntoView({ behavior: 'smooth', block: 'end' })
  }, [messages, sending])

  const handleNewSession = async () => {
    if (sending) return
    if (!loggedIn) { requireLogin(); return }
    try {
      setNotice('')
      await createSession()
    } catch (error) {
      handleRequestError(error)
    }
  }

  const deleteSession = async (sessionID: string) => {
    if (sending || deletingSessionID) return
    try {
      setDeletingSessionID(sessionID)
      await aiChatApi.deleteSession(sessionID)
      const remaining = sessions.filter((session) => session.session_id !== sessionID)
      setSessions(remaining)
      if (activeSessionID === sessionID) {
        const next = remaining[0]
        setActiveSessionID(next?.session_id ?? '')
        setMessages([])
        if (next) await loadMessages(next.session_id)
      }
      setPendingDeleteSessionID(null)
    } catch (error) {
      handleRequestError(error)
    } finally {
      setDeletingSessionID('')
    }
  }

  const handleDeleteSession = (event: MouseEvent<HTMLButtonElement>, sessionID: string) => {
    event.stopPropagation()
    if (!sending && !deletingSessionID) setPendingDeleteSessionID(sessionID)
  }

  const startRenameSession = (event: MouseEvent<HTMLButtonElement>, session: AiChatSession) => {
    event.stopPropagation()
    if (sending || deletingSessionID || renamingSessionID) return
    setEditingSessionID(session.session_id)
    setSessionTitleDraft(session.title || '新建会话')
  }

  const cancelRenameSession = () => {
    setEditingSessionID('')
    setSessionTitleDraft('')
  }

  const renameSession = async (sessionID: string) => {
    const title = sessionTitleDraft.trim()
    const currentTitle = sessions.find((session) => session.session_id === sessionID)?.title.trim()
    if (!title || renamingSessionID) return
    if (title === currentTitle) {
      cancelRenameSession()
      return
    }
    try {
      setRenamingSessionID(sessionID)
      setNotice('')
      const session = await aiChatApi.updateSession(sessionID, title)
      setSessions((current) => current.map((item) => item.session_id === sessionID ? session : item))
      cancelRenameSession()
    } catch (error) {
      handleRequestError(error)
    } finally {
      setRenamingSessionID('')
    }
  }

  const send = async (text = draft) => {
    const content = text.trim()
    if (!content || sending || !assistantEnabled || !collectionName) return
    if (!loggedIn) { requireLogin(); return }
    let sessionID = activeSessionID
    try {
      setNotice('')
      setSending(true)
      if (!sessionID) sessionID = (await createSession()).session_id
      const now = new Date().toISOString()
      const placeholderID = `assistant-${Date.now()}`
      setMessages((current) => [...current,
        { localID: `user-${Date.now()}`, session_id: sessionID, role: 'user', content, created_at: now },
        { localID: placeholderID, session_id: sessionID, role: 'assistant', content: '', created_at: now },
      ])
      setDraft('')

      const response = await aiChatApi.chat({ session_id: sessionID, content, collection_name: collectionName })
      const reader = response.body?.getReader()
      if (!reader) throw new Error('浏览器不支持流式响应')
      const decoder = new TextDecoder()
      let buffer = ''
      let finished = false

      const consume = (block: string) => {
        const data = block.split('\n').filter((line) => line.startsWith('data:')).map((line) => line.slice(5).trim()).join('\n')
        if (!data || data === '[DONE]') return
        const chunk = JSON.parse(data) as AiChatStreamChunk
        if (chunk.error) throw new Error(chunk.error)
        if (chunk.content) {
          setMessages((current) => current.map((item) => item.localID === placeholderID ? { ...item, content: item.content + chunk.content } : item))
        }
        finished = chunk.done
      }

      while (!finished) {
        const { done, value } = await reader.read()
        buffer += decoder.decode(value ?? new Uint8Array(), { stream: !done })
        buffer = buffer.replace(/\r\n/g, '\n')
        const blocks = buffer.split('\n\n')
        buffer = blocks.pop() ?? ''
        blocks.forEach(consume)
        if (done) break
      }
    } catch (error) {
      handleRequestError(error)
      setMessages((current) => current.filter((item) => !item.localID.startsWith('assistant-') || item.content))
    } finally {
      setSending(false)
    }
  }

  return <main className="assistant-page">
    <section className={`assistant-workspace${sidebarCollapsed ? ' sidebar-collapsed' : ''}`} aria-label="站内助手对话">
      <aside className="assistant-sidebar">
        <div className="assistant-side-brand"><span>R</span><div><strong>站内助手</strong><small>{collectionName ? `当前：${collectionName}` : '未选择知识库'}</small></div></div>
        <button className="assistant-sidebar-toggle" type="button" aria-label={sidebarCollapsed ? '展开会话栏' : '收起会话栏'} title={sidebarCollapsed ? '展开会话栏' : '收起会话栏'} onClick={() => setSidebarCollapsed((current) => !current)}>{sidebarCollapsed ? '›' : '‹'}</button>
        <button className="assistant-new-chat" type="button" onClick={() => void handleNewSession()} disabled={sending}><span aria-hidden="true">+</span><span className="assistant-new-chat-label">新建对话</span></button>
        <div className="assistant-session-label">最近会话</div>
        <div className="assistant-session-list">
          {loading && <span className="assistant-side-status">正在加载...</span>}
          {!loading && sessions.length === 0 && <span className="assistant-side-status">{loggedIn ? '还没有会话' : '登录后保存会话'}</span>}
          {sessions.map((session) => <div className={`assistant-session${session.session_id === activeSessionID ? ' active' : ''}`} key={session.session_id}>
            {editingSessionID === session.session_id
              ? <input className="assistant-session-title-input" value={sessionTitleDraft} maxLength={256} autoFocus disabled={renamingSessionID === session.session_id} onClick={(event) => event.stopPropagation()} onChange={(event) => setSessionTitleDraft(event.target.value)} onBlur={() => void renameSession(session.session_id)} onKeyDown={(event) => { if (event.key === 'Enter') { event.preventDefault(); event.currentTarget.blur() } if (event.key === 'Escape') cancelRenameSession() }} aria-label="会话名称" />
              : <><button className="assistant-session-select" type="button" onClick={() => void selectSession(session.session_id)} aria-label={`打开会话：${session.title || '新建会话'}`}><span>◌</span></button><button className="assistant-session-title" type="button" title="点击修改会话名称" disabled={Boolean(deletingSessionID) || Boolean(renamingSessionID)} onClick={(event) => startRenameSession(event, session)}><b>{session.title || '新建会话'}</b></button></>}
            <button className="assistant-delete-session" type="button" title="删除会话" aria-label="删除会话" disabled={Boolean(deletingSessionID) || Boolean(renamingSessionID)} onClick={(event) => handleDeleteSession(event, session.session_id)}>×</button>
          </div>)}
        </div>
        <div className="assistant-side-foot"><i /> AI 回复仅供参考</div>
      </aside>

      <section className={`assistant-chat-panel${isEmpty ? ' is-empty' : ''}`}>
        <div className="assistant-message-list">
          {isEmpty && <div className="assistant-welcome"><h2>从一个问题开始</h2><p>查找本站文章、了解作者观点，或查询作者创建的知识库。</p><div>{suggestions.map((item) => <button type="button" key={item} disabled={sending} onClick={() => void send(item)}>{item}</button>)}</div></div>}
          {messages.map((message) => <article className={`assistant-message ${message.role === 'user' ? 'user' : 'bot'}`} key={message.localID}>
            <div className="assistant-message-avatar">{message.role === 'user' ? '我' : 'R'}</div>
            <div><div className="assistant-message-role">{message.role === 'user' ? '你' : '小爱'}</div><p>{message.content || (sending ? <span className="assistant-typing">正在思考</span> : '')}</p></div>
          </article>)}
          <div ref={messageEnd} />
        </div>
        {!loggedIn && <div className="assistant-notice assistant-login-notice" role="status"><span>需要登录后才能使用站内助手</span><button type="button" onClick={() => navigate('/login')}>去登录</button></div>}
        {notice && loggedIn && <div className="assistant-notice" role="status"><span>{notice}</span></div>}
        {!assistantEnabled && <div className="assistant-notice" role="status">站内助手已由管理员关闭，暂时无法发送消息。</div>}
        <form className="assistant-composer" onSubmit={(event) => { event.preventDefault(); void send() }}><textarea value={draft} onChange={(event) => setDraft(event.target.value)} onKeyDown={(event) => { if (event.key === 'Enter' && !event.shiftKey) { event.preventDefault(); void send() } }} placeholder={assistantEnabled ? '给站内助手发送消息' : '站内助手已关闭'} disabled={sending || !assistantEnabled} rows={1} /><div className="assistant-composer-bottom"><div className="assistant-composer-tools"><KnowledgePicker collections={collections} value={collectionName} disabled={sending || collections.length === 0} onChange={(value) => { setCollectionName(value); localStorage.setItem(COLLECTION_STORAGE_KEY, value) }} /><small>当前对话检索：{collectionName || '未选择'}</small></div><button className="assistant-send-button" type="submit" disabled={sending || !assistantEnabled || !collectionName || !draft.trim()} aria-label="发送消息"><svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 19V5M6.5 10.5 12 5l5.5 5.5" /></svg></button></div></form>
      </section>
      {pendingDeleteSessionID && <div className="assistant-delete-modal" role="presentation" onMouseDown={(event) => { if (event.target === event.currentTarget && !deletingSessionID) setPendingDeleteSessionID(null) }}><section className="assistant-delete-card" role="dialog" aria-modal="true" aria-labelledby="assistant-delete-title"><div className="assistant-delete-icon">!</div><h2 id="assistant-delete-title">删除这个会话？</h2><p>会话中的全部聊天记录将被永久删除，无法恢复。</p><div><button type="button" disabled={Boolean(deletingSessionID)} onClick={() => setPendingDeleteSessionID(null)}>取消</button><button className="assistant-delete-confirm" type="button" disabled={Boolean(deletingSessionID)} onClick={() => void deleteSession(pendingDeleteSessionID)}>{deletingSessionID ? '删除中...' : '删除会话'}</button></div></section></div>}
    </section>
  </main>
}
