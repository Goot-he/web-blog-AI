import { useEffect, useState } from 'react'
import { likeApi } from '@/shared/api/v1/like'
import type { Article } from '@/shared/types/article'
import { logRequestError } from '@/shared/lib/request-error'

const formatCount = (value: number) => new Intl.NumberFormat('zh-CN').format(value ?? 0)

export default function LikesView() {
  const [articles, setArticles] = useState<Article[]>([])
  const [loading, setLoading] = useState(true)
  const [message, setMessage] = useState('')

  const loadLikes = async () => {
    try {
      setLoading(true)
      const result = await likeApi.list({ objectType: 'article', page: 1, pageSize: 100 })
      setArticles(result.articles ?? [])
    } catch (error) {
      logRequestError('加载点赞列表失败', error)
      setMessage('加载失败，请稍后重试')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { void loadLikes() }, [])

  return (
    <section className="admin-card">
      <div className="admin-card-head">
        <div>
          <span className="eyebrow">LIKES</span>
          <h2>点赞列表</h2>
          <p>查看当前账号点过赞的文章。</p>
        </div>
        <button className="primary-button" type="button" onClick={() => void loadLikes()} disabled={loading}>刷新</button>
      </div>
      {message && <p className="articles-message" role="alert">{message}</p>}
      {loading ? <div className="articles-empty">正在加载点赞列表...</div> : articles.length === 0 ? <div className="articles-empty"><strong>还没有点赞任何文章</strong><span>点过赞的文章会出现在这里。</span></div> : <div className="article-list">{articles.map((article) => <article className="article-row" key={article.id} role="link" tabIndex={0} onClick={() => window.open(`/blog/${article.id}`, '_self')} onKeyDown={(event) => { if (event.key === 'Enter' || event.key === ' ') { event.preventDefault(); window.open(`/blog/${article.id}`, '_self') } }}><div className="article-cover">{article.coverImage ? <img src={article.coverImage} alt="文章封面" /> : <span>NO COVER</span>}</div><div className="article-info"><h3>{article.title || '未命名文章'}</h3><p>{article.summary || '暂无摘要'}</p><div className="article-meta"><span>{new Date(article.createdAt * 1000).toLocaleDateString('zh-CN')}</span><span className="article-views">浏览 {formatCount(article.viewCount)}</span><span>评论 {formatCount(article.commentCount)}</span><span>点赞 {formatCount(article.likeCount)}</span></div></div></article>)}</div>}
    </section>
  )
}
