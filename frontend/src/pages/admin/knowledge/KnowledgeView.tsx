import { useCallback, useEffect, useState } from 'react'
import { vectorApi } from '@/shared/api/v1/vector'
import type { VectorCollection, VectorDocument } from '@/shared/types/vector'
import { logRequestError } from '@/shared/lib/request-error'
import './KnowledgeView.css'

const DOCUMENT_PAGE_SIZE = 5

function KnowledgeIcon() { return <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M5 5.5A2.5 2.5 0 0 1 7.5 3H19v15.5a2.5 2.5 0 0 0-2.5-2.5H5V5.5Z" /><path d="M5 16h11.5a2.5 2.5 0 0 1 2.5 2.5V21H7.5A2.5 2.5 0 0 1 5 18.5V16Z" /><path d="M9 7h6M9 10h5" /></svg> }
function RefreshIcon() { return <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M20 11a8 8 0 0 0-14.8-4.2L3 9" /><path d="M3 4v5h5M4 13a8 8 0 0 0 14.8 4.2L21 15" /><path d="M21 20v-5h-5" /></svg> }

export default function KnowledgeView() {
  const [collections, setCollections] = useState<VectorCollection[]>([])
  const [selected, setSelected] = useState('')
  const [documents, setDocuments] = useState<VectorDocument[]>([])
  const [documentTotal, setDocumentTotal] = useState(0)
  const [documentPage, setDocumentPage] = useState(1)
  const [newName, setNewName] = useState('')
  const [loading, setLoading] = useState(true)
  const [detailLoading, setDetailLoading] = useState(false)
  const [error, setError] = useState('')

  const loadCollections = useCallback(async () => { setLoading(true); setError(''); try { setCollections(await vectorApi.collections()) } catch (requestError) { logRequestError('加载知识库列表失败', requestError); setError('加载知识库失败，请稍后重试') } finally { setLoading(false) } }, [])
  useEffect(() => { void loadCollections() }, [loadCollections])

  const loadDocuments = async (name: string, page: number) => {
    setDetailLoading(true); setError('')
    try { const result = await vectorApi.documents(name, page, DOCUMENT_PAGE_SIZE); setDocuments(result.items); setDocumentTotal(result.total); setDocumentPage(result.page) } catch (requestError) { logRequestError('加载知识库文档失败', requestError); setError('加载文档失败，请稍后重试') } finally { setDetailLoading(false) }
  }
  const openCollection = async (name: string) => {
    setSelected(name)
    await loadDocuments(name, 1)
  }
  const createCollection = async () => {
    const name = newName.trim(); if (!name) return
    try { await vectorApi.createCollection(name); setNewName(''); await loadCollections(); await openCollection(name) } catch (requestError) { logRequestError('创建知识库失败', requestError); setError('创建知识库失败，请稍后重试') }
  }
  const deleteCollection = async (name: string) => {
    if (!window.confirm(`确定删除知识库“${name}”及其中的全部文档吗？`)) return
    try { await vectorApi.deleteCollection(name); if (selected === name) { setSelected(''); setDocuments([]); setDocumentTotal(0) } await loadCollections() } catch (requestError) { logRequestError('删除知识库失败', requestError); setError('删除知识库失败，请稍后重试') }
  }

  const totalPages = Math.max(1, Math.ceil(documentTotal / DOCUMENT_PAGE_SIZE))

  return <section className="knowledge-view"><header className="knowledge-header"><div><span className="eyebrow">KNOWLEDGE BASE</span><h2>知识库管理</h2><p>创建、删除知识库，并查看指定知识库中的文档基础数据。</p></div><button className="knowledge-refresh" type="button" onClick={() => void loadCollections()} disabled={loading} title="刷新知识库列表"><RefreshIcon /><span>刷新</span></button></header><div className="knowledge-create"><input value={newName} onChange={(event) => setNewName(event.target.value)} placeholder="输入新知识库名称" /><button className="primary-button" type="button" disabled={!newName.trim()} onClick={() => void createCollection()}>新增知识库</button></div>{error && <p className="settings-notice" role="status">{error}</p>}{loading ? <div className="knowledge-state"><p>正在加载知识库...</p></div> : collections.length === 0 ? <div className="knowledge-state"><div className="knowledge-icon"><KnowledgeIcon /></div><strong>暂未创建知识库</strong><p>输入名称即可创建一个空知识库。</p></div> : <div className="knowledge-layout"><div className="knowledge-grid">{collections.map((collection) => <article className={`knowledge-card${selected === collection.name ? ' active' : ''}`} key={collection.name} onClick={() => void openCollection(collection.name)}><button className="knowledge-delete" type="button" title="删除知识库" onClick={(event) => { event.stopPropagation(); void deleteCollection(collection.name) }}>×</button><div className="knowledge-icon"><KnowledgeIcon /></div><div className="knowledge-card-body"><h3 title={collection.name}>{collection.name}</h3><p>{selected === collection.name ? '当前查看知识库' : '点击查看文档'}</p></div><div className="knowledge-count"><strong>{collection.count.toLocaleString()}</strong><span>个切块</span></div></article>)}</div>{selected && <section className="knowledge-documents"><header><div><span className="eyebrow">DOCUMENTS</span><h3>{selected}</h3></div><strong>{documentTotal} 个切块</strong></header>{detailLoading ? <p>正在加载文档...</p> : documents.length === 0 ? <div className="knowledge-state"><p>这个知识库还没有文档。</p></div> : <><div className="knowledge-table"><div className="knowledge-table-row head"><span>文档</span><span>类型 / 来源</span><span>内容摘要</span></div>{documents.map((document) => <div className="knowledge-table-row" key={document.chunk_id}><span><b>{document.title || document.doc_id || document.chunk_id}</b><small>切块 #{document.chunk_no}</small></span><span><b>{document.knowledge_type || '-'}</b><small>{document.source || '-'}</small></span><p>{document.content}</p></div>)}</div><footer className="knowledge-pagination"><span>第 {documentPage} / {totalPages} 页，每页 5 条</span><div><button type="button" disabled={documentPage <= 1} onClick={() => void loadDocuments(selected, documentPage - 1)}>上一页</button><button type="button" disabled={documentPage >= totalPages} onClick={() => void loadDocuments(selected, documentPage + 1)}>下一页</button></div></footer></>}</section>}</div>}</section>
}
