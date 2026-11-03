import { useEffect, useMemo, useState } from 'react'
import { articleApi } from '@/shared/api/v1/article'
import { userApi } from '@/shared/api/v1/user'
import { navigate } from '@/shared/hooks/usePathname'
import { resolveStorageUrl } from '@/shared/lib/storage'
import type { Article } from '@/shared/types/article'
import type { CurrentUser } from '@/shared/types/user'
import { logRequestError } from '@/shared/lib/request-error'
import { readCurrentUser } from '@/shared/lib/current-user'
import './ArticlesView.css'

const formatCount = (value: number) => new Intl.NumberFormat('zh-CN').format(value ?? 0)
const formatDate = (value: number) => value ? new Date(value * 1000).toLocaleDateString('zh-CN') : '未发布'
const DAY_MS = 24 * 60 * 60 * 1000

interface TrendPoint {
  label: string
  total: number
  fans: number
}

function sumBy(items: Article[], pick: (item: Article) => number | undefined) {
  return items.reduce((sum, item) => sum + (pick(item) ?? 0), 0)
}

function getDayKey(date: Date) {
  return `${date.getFullYear()}-${date.getMonth()}-${date.getDate()}`
}

function formatAxisDate(date: Date) {
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${date.getFullYear()}/${month}/${day}`
}

function buildSevenDayTrend(articles: Article[]): TrendPoint[] {
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  const points = Array.from({ length: 7 }, (_, index) => {
    const date = new Date(today.getTime() - (6 - index) * DAY_MS)
    return { label: formatAxisDate(date), key: getDayKey(date), total: 0, fans: 0 }
  })
  const byKey = new Map(points.map((point) => [point.key, point]))

  articles.forEach((article) => {
    const time = article.publishedAt || article.createdAt
    if (!time) return
    const date = new Date(time * 1000)
    date.setHours(0, 0, 0, 0)
    const point = byKey.get(getDayKey(date))
    if (point) point.total += article.viewCount ?? 0
  })

  return points.map(({ label, total, fans }) => ({ label, total, fans }))
}

function PlaybackTrend({ points }: { points: TrendPoint[] }) {
  const width = 920
  const height = 300
  const padding = { top: 26, right: 28, bottom: 46, left: 36 }
  const chartWidth = width - padding.left - padding.right
  const chartHeight = height - padding.top - padding.bottom
  const maxData = Math.max(...points.map((point) => Math.max(point.total, point.fans)), 0)
  const maxValue = Math.max(4, Math.ceil(maxData / 4) * 4)
  const ticks = Array.from({ length: 5 }, (_, index) => maxValue - (maxValue / 4) * index)
  const x = (index: number) => padding.left + (points.length === 1 ? 0 : (chartWidth / (points.length - 1)) * index)
  const y = (value: number) => padding.top + chartHeight - (value / maxValue) * chartHeight
  const linePoints = (pick: (point: TrendPoint) => number) => points.map((point, index) => `${x(index)},${y(pick(point))}`).join(' ')

  return <section className="creator-chart-card">
    <header className="creator-chart-head">
      <h3>近7天播放量</h3>
      <div className="creator-chart-legend" aria-label="播放量图例">
        <span><i className="legend-line total" />总播放量</span>
        <span><i className="legend-line fans" />粉丝播放量</span>
      </div>
    </header>
    <svg className="creator-chart-svg" viewBox={`0 0 ${width} ${height}`} role="img" aria-label="近7天播放量趋势">
      {ticks.map((tick) => {
        const tickY = y(tick)
        return <g key={tick}>
          <line className="chart-grid" x1={padding.left} y1={tickY} x2={width - padding.right} y2={tickY} />
          <text className="chart-axis-label y" x={padding.left - 12} y={tickY + 4}>{formatCount(tick)}</text>
        </g>
      })}
      {points.map((point, index) => {
        const tickX = x(index)
        return <g key={point.label}>
          <line className="chart-grid vertical" x1={tickX} y1={padding.top} x2={tickX} y2={height - padding.bottom} />
          <text className="chart-axis-label x" x={tickX} y={height - 16}>{point.label}</text>
        </g>
      })}
      <polyline className="chart-line fans" points={linePoints((point) => point.fans)} />
      <polyline className="chart-line total" points={linePoints((point) => point.total)} />
      {points.map((point, index) => <circle className="chart-dot total" key={`total-${point.label}`} cx={x(index)} cy={y(point.total)} r="3.5" />)}
      {points.map((point, index) => <circle className="chart-dot fans" key={`fans-${point.label}`} cx={x(index)} cy={y(point.fans)} r="3.5" />)}
    </svg>
  </section>
}

function ArticleRow({ article, draft = false, onDelete, onEdit }: { article: Article; draft?: boolean; onDelete: (article: Article) => void; onEdit: (article: Article) => void }) {
  const open = () => {
    if (!draft) navigate(`/blog/${article.id}`)
  }

  return <article className={`article-row${draft ? ' draft' : ''}`} key={article.id} role={draft ? undefined : 'link'} tabIndex={draft ? undefined : 0} onClick={open} onKeyDown={(event) => { if (!draft && event.key === 'Enter') open() }}>
    <div className="article-cover">{article.coverImage ? <img src={resolveStorageUrl(article.coverImage)} alt="作品封面" /> : <span>{draft ? 'DRAFT' : 'NO COVER'}</span>}</div>
    <div className="article-info">
      <h3>{article.title || '未命名作品'}</h3>
      <p>{article.summary || (draft ? '草稿还没有摘要。' : '暂无摘要')}</p>
      <div className="article-meta">
        <span>{draft ? `更新 ${formatDate(article.updatedAt || article.createdAt)}` : formatDate(article.createdAt)}</span>
        <span>浏览 {formatCount(article.viewCount)}</span>
        <span>评论 {formatCount(article.commentCount)}</span>
        <span>点赞 {formatCount(article.likeCount)}</span>
        <span>收藏 {formatCount(article.favorCount)}</span>
      </div>
    </div>
    {draft && <span className="article-status-badge">草稿</span>}
    <button className="article-edit" type="button" onClick={(event) => { event.stopPropagation(); onEdit(article) }}>编辑</button>
    <button className="article-delete" type="button" onClick={(event) => { event.stopPropagation(); onDelete(article) }}>删除</button>
  </article>
}

export default function ArticlesView({ mode = 'dashboard' }: { mode?: 'dashboard' | 'published' | 'drafts' }) {
  const [articles, setArticles] = useState<Article[]>([])
  const [drafts, setDrafts] = useState<Article[]>([])
  const [publishedTotal, setPublishedTotal] = useState(0)
  const [draftTotal, setDraftTotal] = useState(0)
  const [profileStats, setProfileStats] = useState<Partial<CurrentUser>>(() => readCurrentUser())
  const [loading, setLoading] = useState(true)

  const stats = useMemo(() => {
    const receiveLikeCount = profileStats.receive_like_count ?? sumBy(articles, (item) => item.likeCount)
    const receiveFavorCount = profileStats.receive_favor_count ?? sumBy(articles, (item) => item.favorCount)

    const viewCount = profileStats.view_count ?? sumBy(articles, (item) => item.viewCount)
    const commentCount = profileStats.comment_count ?? sumBy(articles, (item) => item.commentCount)

    return [
      { label: '已发布', value: publishedTotal },
      { label: '草稿', value: draftTotal },
      { label: '浏览', value: viewCount },
      { label: '收到点赞', value: receiveLikeCount },
      { label: '收到收藏', value: receiveFavorCount },
      { label: '收到评论', value: commentCount },
      { label: '全部作品', value: publishedTotal + draftTotal },
    ]
  }, [articles, draftTotal, profileStats.comment_count, profileStats.receive_favor_count, profileStats.receive_like_count, profileStats.view_count, publishedTotal])
  const trend = useMemo(() => buildSevenDayTrend(articles), [articles])

  const loadArticles = async () => {
    try {
      setLoading(true)
      const currentUser = readCurrentUser() as { id?: number }
      const [publishedResult, draftResult, profileResult] = await Promise.all([
        articleApi.listByUser({ authorID: currentUser.id, page: 1, pageSize: 50 }),
        articleApi.listDrafts({ page: 1, pageSize: 50 }),
        userApi.getProfile().catch((error) => {
          logRequestError('加载用户统计失败', error)
          return null
        }),
      ])
      setArticles(publishedResult.articles ?? [])
      setDrafts(draftResult.articles ?? [])
      setPublishedTotal(publishedResult.total ?? publishedResult.articles?.length ?? 0)
      setDraftTotal(draftResult.total ?? draftResult.articles?.length ?? 0)
      if (profileResult) setProfileStats(profileResult)
    } catch (error) {
      logRequestError('加载我的文章失败', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { void loadArticles() }, [])

  const removeArticle = async (article: Article, draft = false) => {
    if (!window.confirm(`确定要删除《${article.title || '未命名作品'}》吗？\n\n删除后无法恢复。`)) return
    try {
      if (draft) {
        await articleApi.deleteDraft({ id: article.id })
      } else {
        await articleApi.delete({ id: article.id })
      }
      if (draft) {
        setDrafts((current) => current.filter((item) => item.id !== article.id))
        setDraftTotal((current) => Math.max(0, current - 1))
      } else {
        setArticles((current) => current.filter((item) => item.id !== article.id))
        setPublishedTotal((current) => Math.max(0, current - 1))
      }
    } catch (error) {
      logRequestError('删除我的文章失败', error)
    }
  }

  const editArticle = (article: Article) => {
    navigate('/admin/create', { article })
  }

  const published = <section className="article-section">
    <header><div><h3>我的文章</h3><span>{formatCount(publishedTotal)} 个作品</span></div><button type="button" onClick={() => navigate('/admin/create')}>发布作品</button></header>
    {loading ? <div className="articles-empty">正在加载作品...</div> : articles.length === 0 ? <div className="articles-empty"><strong>还没有发布文章</strong><span>完成一篇投稿后，它会出现在这里。</span></div> : <div className="article-list">{articles.map((article) => <ArticleRow article={article} key={article.id} onEdit={editArticle} onDelete={(item) => void removeArticle(item)} />)}</div>}
  </section>

  const draftList = <section className="article-section">
    <header><div><h3>草稿箱</h3><span>{formatCount(draftTotal)} 个草稿</span></div><button type="button" onClick={() => navigate('/admin/create')}>继续投稿</button></header>
    {loading ? <div className="articles-empty">正在加载草稿...</div> : drafts.length === 0 ? <div className="articles-empty"><strong>草稿箱是空的</strong><span>保存至草稿箱的作品会出现在这里。</span></div> : <div className="article-list">{drafts.map((article) => <ArticleRow article={article} draft key={article.id} onEdit={editArticle} onDelete={(item) => void removeArticle(item, true)} />)}</div>}
  </section>

  return <section className="articles-view">
    {mode === 'dashboard' && <>
      <div className="creator-stats" aria-label="投稿数据统计">
        {stats.map((item) => <article key={item.label}><span>{item.label}</span><strong>{formatCount(item.value)}</strong></article>)}
      </div>
      {loading ? <div className="articles-empty">正在加载数据...</div> : <PlaybackTrend points={trend} />}
    </>}

    {mode === 'drafts' ? draftList : mode === 'published' ? published : null}
  </section>
}
