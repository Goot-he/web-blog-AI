import { useEffect, useState } from 'react'
import { favorApi } from '@/shared/api/v1/favor'
import type { Article } from '@/shared/types/article'
import { logRequestError } from '@/shared/lib/request-error'
import { resolveStorageUrl } from '@/shared/lib/storage'

export default function FavoritesView() {
  const [articles, setArticles] = useState<Article[]>([])
  const [loading, setLoading] = useState(true)
  const [message, setMessage] = useState('')
  const loadFavorites = async () => {
    try { setLoading(true); setArticles((await favorApi.list({ page: 1, pageSize: 100 })).articles ?? []) }
    catch (error) { logRequestError('加载收藏列表失败', error); setMessage('加载失败，请稍后重试') }
    finally { setLoading(false) }
  }
  useEffect(() => { void loadFavorites() }, [])
  return <section className="admin-card"><div className="admin-card-head"><div><span className="eyebrow">FAVORITES</span><h2>我的收藏</h2><p>收藏的文章会保存在默认收藏夹中。</p></div><button className="primary-button" type="button" onClick={() => void loadFavorites()} disabled={loading}>刷新</button></div>{message && <p className="articles-message" role="alert">{message}</p>}{loading ? <div className="articles-empty">正在加载收藏夹...</div> : articles.length === 0 ? <div className="articles-empty"><strong>收藏夹还是空的</strong><span>收藏文章后，它们会出现在这里。</span></div> : <div className="article-list">{articles.map((article) => <article className="article-row" key={article.id} role="link" tabIndex={0} onClick={() => window.open(`/blog/${article.id}`, '_self')} onKeyDown={(event) => { if (event.key === 'Enter' || event.key === ' ') { event.preventDefault(); window.open(`/blog/${article.id}`, '_self') } }}><div className="article-cover">{article.coverImage ? <img src={resolveStorageUrl(article.coverImage)} alt="文章封面" /> : <span>NO COVER</span>}</div><div className="article-info"><h3>{article.title || '未命名文章'}</h3><p>{article.summary || '暂无摘要'}</p><div className="article-meta"><span>{new Date(article.createdAt * 1000).toLocaleDateString('zh-CN')}</span><span>收藏 {article.favorCount}</span><span>点赞 {article.likeCount}</span></div></div></article>)}</div>}</section>
}
