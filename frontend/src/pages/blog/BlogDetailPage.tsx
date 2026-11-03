import { useEffect, useMemo, useRef, useState, type FormEvent } from 'react'
import { articleApi } from '@/shared/api/v1/article'
import { commentApi } from '@/shared/api/v1/comment'
import { likeApi } from '@/shared/api/v1/like'
import { favorApi } from '@/shared/api/v1/favor'
import type { FavoriteFolder } from '@/shared/types/favor'
import AppLink from '@/shared/ui/AppLink'
import MarkdownView from '@/shared/ui/MarkdownView'
import { navigate } from '@/shared/hooks/usePathname'
import type { Article } from '@/shared/types/article'
import type { CommentItem } from '@/shared/types/comment'
import { resolveStorageUrl } from '@/shared/lib/storage'
import { logRequestError } from '@/shared/lib/request-error'
import { readCurrentUser } from '@/shared/lib/current-user'
import './BlogPage.css'

const formatDateTime = (value: string | number) => {
  const date = typeof value === 'number' ? new Date(value > 1_000_000_000_000 ? value : value * 1000) : new Date(value)
  return Number.isNaN(date.getTime()) ? String(value) : date.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
}

const Avatar = ({ src, name, small = false, userID }: { src?: string; name: string; small?: boolean; userID?: number }) => {
  const content = <span className={small ? 'front-comment-avatar small' : 'front-comment-avatar'}>
    {src ? <img src={resolveStorageUrl(src)} alt="" /> : name.slice(0, 1).toUpperCase()}
  </span>

  return userID ? <AppLink className="front-comment-avatar-link" to={`/users/${userID}`} aria-label={`查看${name}的主页`} title={`查看${name}的主页`}>{content}</AppLink> : content
}

function ActionIcon({ name }: { name: 'view' | 'comment' | 'like' | 'favor' | 'top' }) {
  const paths = name === 'view'
    ? <><path d="M2.5 12s3.4-5 9.5-5 9.5 5 9.5 5-3.4 5-9.5 5-9.5-5-9.5-5Z" /><circle cx="12" cy="12" r="2.5" /></>
    : name === 'comment'
      ? <><path d="M4 5.5h16v11H9l-5 3v-14Z" /><path d="M8 10h8M8 13h5" /></>
      : name === 'like'
        ? <path d="M12 20.5 4.8 13.7A4.6 4.6 0 0 1 11.3 7L12 7.8l.7-.8a4.6 4.6 0 0 1 6.5 6.7L12 20.5Z" />
        : name === 'favor'
          ? <path d="m12 3 2.78 5.63 6.22.9-4.5 4.39 1.06 6.2L12 17.2l-5.56 2.92 1.06-6.2L3 9.53l6.22-.9L12 3Z" />
        : <><path d="m6 10 6-6 6 6" /><path d="M12 4v16" /></>
  return <svg viewBox="0 0 24 24" aria-hidden="true" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">{paths}</svg>
}

export default function BlogDetailPage() {
  const articleID = Number(window.location.pathname.split('/').filter(Boolean).at(-1))
  const [article, setArticle] = useState<Article | null>(null)
  const [comments, setComments] = useState<CommentItem[]>([])
  const [replies, setReplies] = useState<Record<number, CommentItem[]>>({})
  const [expandedReplies, setExpandedReplies] = useState<Record<number, boolean>>({})
  const [commentText, setCommentText] = useState('')
  const [replyTarget, setReplyTarget] = useState<CommentItem | null>(null)
  const [replyText, setReplyText] = useState('')
  const [loading, setLoading] = useState(true)
  const [commentsLoading, setCommentsLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [message, setMessage] = useState('')
  const [liked, setLiked] = useState(false)
  const [liking, setLiking] = useState(false)
  const [favored, setFavored] = useState(false)
  const [favoring, setFavoring] = useState(false)
  const [favorPickerOpen, setFavorPickerOpen] = useState(false)
  const [folders, setFolders] = useState<FavoriteFolder[]>([])
  const [newFolderName, setNewFolderName] = useState('')
  const [likingCommentIDs, setLikingCommentIDs] = useState<Record<number, boolean>>({})
  const loadedArticleID = useRef<number | null>(null)
  const loggedIn = Boolean(localStorage.getItem('renai_access_token'))
  const currentUser = readCurrentUser() as { id?: number; username?: string; avatar?: string }

  const loadComments = async () => {
    if (!articleID) return
    setCommentsLoading(true)
    try {
      const result = await commentApi.list({ articleId: articleID, page: 1, size: 50 })
      setComments(result.comments ?? [])
    } catch (reason) {
      logRequestError('加载文章评论失败', reason)
      setMessage('加载评论失败，请稍后重试')
    } finally {
      setCommentsLoading(false)
    }
  }

  const loadReplies = async (rootID: number) => {
    const result = await commentApi.replies({ parentId: rootID, page: 1, size: 50 })
    const nextReplies = result.replies ?? []
    setReplies((current) => ({ ...current, [rootID]: nextReplies }))
    return nextReplies
  }

  useEffect(() => {
    if (!articleID) {
      setLoading(false)
      return
    }
    if (loadedArticleID.current !== articleID) {
      loadedArticleID.current = articleID
      void articleApi.detail({ id: articleID }).then((result) => {
        setArticle(result.article)
        setLiked(result.isLiked)
        setFavored(result.isFavored)
      }).catch((reason: unknown) => { logRequestError('加载文章详情失败', reason); setMessage('加载文章失败，请稍后重试') }).finally(() => setLoading(false))
    }
    void loadComments()
  }, [articleID])

  const headings = useMemo(() => article ? Array.from(article.content.matchAll(/^(#{1,3})\s+(.+)$/gm)).map((match) => ({ level: match[1].length, text: match[2].trim() })) : [], [article])
  const readingMinutes = useMemo(() => Math.max(1, Math.ceil((article?.content.replace(/[#*`_>\s]/g, '').length ?? 0) / 450)), [article])
  const jumpToHeading = (index: number) => document.querySelectorAll('.front-article-body h1,.front-article-body h2,.front-article-body h3').item(index)?.scrollIntoView({ behavior: 'smooth', block: 'start' })

  const submitComment = async (event: FormEvent) => {
    event.preventDefault()
    if (!loggedIn) return navigate('/login')
    if (!commentText.trim()) return
    try {
      setSubmitting(true)
      const result = await commentApi.create({ articleId: articleID, content: commentText.trim() })
      if (!result.success) throw new Error('评论发布失败')
      setCommentText('')
      setMessage('评论已发布')
      await loadComments()
    } catch (reason) {
      logRequestError('发布评论失败', reason)
      setMessage('评论发布失败，请稍后重试')
    } finally {
      setSubmitting(false)
    }
  }

  const toggleReplies = async (comment: CommentItem) => {
    if (expandedReplies[comment.id]) {
      setExpandedReplies((current) => ({ ...current, [comment.id]: false }))
      return
    }
    setExpandedReplies((current) => ({ ...current, [comment.id]: true }))
    if (replies[comment.id]) return
    try {
      await loadReplies(comment.id)
    } catch (reason) {
      logRequestError('加载评论回复失败', reason)
      setMessage('加载回复失败，请稍后重试')
    }
  }

  const submitReply = async (event: FormEvent) => {
    event.preventDefault()
    if (!replyTarget || !replyText.trim()) return
    if (!loggedIn) return navigate('/login')
    const rootID = replyTarget.rootId || replyTarget.id
    try {
      setSubmitting(true)
      const result = await commentApi.reply({ articleId: articleID, parentId: rootID, replyToId: replyTarget.id, content: replyText.trim() })
      if (!result.success) throw new Error('回复发布失败')
      setReplyText('')
      setReplyTarget(null)
      setExpandedReplies((current) => ({ ...current, [rootID]: true }))
      await Promise.all([loadComments(), loadReplies(rootID)])
      setMessage('回复已发布')
    } catch (reason) {
      logRequestError('发布回复失败', reason)
      setMessage('回复发布失败，请稍后重试')
    } finally {
      setSubmitting(false)
    }
  }

  const toggleArticleLike = async () => {
    if (!loggedIn) return navigate('/login')
    if (!article || liking) return
    setLiking(true)
    try {
      const result = liked
        ? await likeApi.cancelThumbUp({ objectType: 'article', objectId: article.id })
        : await likeApi.thumbUp({ objectType: 'article', objectId: article.id })
      if (!result.success) throw new Error('点赞操作失败')
      setArticle({ ...article, likeCount: Math.max(0, article.likeCount + (liked ? -1 : 1)) })
      setLiked(!liked)
    } catch (reason) {
      logRequestError('更新文章点赞失败', reason)
      setMessage('操作失败，请稍后重试')
    } finally {
      setLiking(false)
    }
  }

  const toggleArticleFavor = async (folderId?: number) => {
    if (!loggedIn) return navigate('/login')
    if (!article || favoring) return
    setFavoring(true)
    try {
      const result = favored
        ? await favorApi.cancel({ objectId: article.id })
        : await favorApi.favor({ objectId: article.id, folderId })
      if (!result.success) throw new Error('收藏操作失败')
      setArticle({ ...article, favorCount: Math.max(0, article.favorCount + (favored ? -1 : 1)) })
      setFavored(!favored)
      setFavorPickerOpen(false)
      setMessage(favored ? '已取消收藏' : '已加入我的收藏')
    } catch (reason) {
      logRequestError('更新文章收藏失败', reason)
      setMessage('收藏操作失败，请稍后重试')
    } finally {
      setFavoring(false)
    }
  }

  const openFavorPicker = async () => {
    if (!loggedIn) return navigate('/login')
    if (favored) return void toggleArticleFavor()
    try { setFolders((await favorApi.folders()).folders ?? []); setFavorPickerOpen(true) } catch (reason) { logRequestError('加载收藏夹失败', reason); setMessage('加载收藏夹失败，请稍后重试') }
  }

  const createFolder = async () => {
    const name = newFolderName.trim()
    if (!name) return
    try { const folder = await favorApi.createFolder(name); setFolders((current) => [...current, folder]); setNewFolderName(''); setMessage(`已创建收藏夹“${folder.name}”，请选择收藏夹完成收藏`) } catch (reason) { logRequestError('创建收藏夹失败', reason); setMessage('创建收藏夹失败，请稍后重试') }
  }

  const toggleCommentLike = async (comment: CommentItem) => {
    if (!loggedIn) return navigate('/login')
    if (likingCommentIDs[comment.id]) return

    setLikingCommentIDs((current) => ({ ...current, [comment.id]: true }))
    try {
      const result = comment.isLiked
        ? await likeApi.cancelThumbUp({ objectType: 'comment', objectId: comment.id })
        : await likeApi.thumbUp({ objectType: 'comment', objectId: comment.id })
      if (!result.success) throw new Error('点赞操作失败')

      const update = (item: CommentItem) => item.id === comment.id
        ? { ...item, isLiked: !item.isLiked, likeCount: Math.max(0, item.likeCount + (item.isLiked ? -1 : 1)) }
        : item
      setComments((current) => current.map(update))
      setReplies((current) => Object.fromEntries(Object.entries(current).map(([parentID, items]) => [parentID, items.map(update)])))
    } catch (reason) {
      logRequestError('更新评论点赞失败', reason)
      setMessage('操作失败，请稍后重试')
    } finally {
      setLikingCommentIDs((current) => ({ ...current, [comment.id]: false }))
    }
  }

  if (loading) return <main className="front-detail-state">正在打开文章...</main>
  if (!article) return <main className="front-detail-state"><strong>没有找到这篇文章</strong><AppLink to="/blog">返回博客列表</AppLink>{message && <span>{message}</span>}</main>

  return <main className="front-detail-page">
    <article className="front-reading-column">
      <header className="front-detail-header">
        <AppLink to="/blog" className="front-back-link">返回博客</AppLink>
        <h1>{article.title}</h1>
        <div className="front-detail-author"><Avatar src={article.authorAvatar} name={article.authorName || 'R'} userID={article.authorID} /><div><strong>{article.authorName || 'Renai 成员'}</strong><span>{formatDateTime(article.createdAt)} · 约 {readingMinutes} 分钟阅读</span></div></div>
      </header>
      {article.coverImage && <img className="front-detail-cover" src={resolveStorageUrl(article.coverImage)} alt={`${article.title}封面`} />}
      <section className="front-article-body"><MarkdownView content={article.content} /></section>
      <section className="front-comments" id="comments">
        <header><div><span>一起交流</span><h2>评论</h2></div><strong>{comments.length} 条</strong></header>
        {loggedIn ? <form className="front-comment-form" onSubmit={submitComment}><Avatar src={currentUser.avatar} name={currentUser.username || '我'} userID={currentUser.id} /><div><label htmlFor="new-comment">写下你的想法</label><textarea id="new-comment" value={commentText} maxLength={1000} onChange={(event) => setCommentText(event.target.value)} placeholder="友善交流，让讨论更有价值。" /><footer><span>{commentText.length}/1000</span><button type="submit" disabled={submitting || !commentText.trim()}>{submitting ? '发布中...' : '发布评论'}</button></footer></div></form> : <div className="front-comment-login"><span>登录后参与评论，与社团成员继续讨论。</span><button type="button" onClick={() => navigate('/login')}>去登录</button></div>}
        {message && <p className="front-comment-message" role="status">{message}</p>}
        {commentsLoading ? <div className="front-comments-state">正在加载评论...</div> : comments.length === 0 ? <div className="front-comments-state"><strong>还没有评论</strong><span>成为第一个参与讨论的人。</span></div> : <div className="front-comment-list">{comments.map((comment) => <article className="front-comment" key={comment.id}><Avatar src={comment.userAvatar} name={comment.userName || 'R'} userID={comment.userId} /><div><header><strong>{comment.userName || 'Renai 成员'}</strong><time>{formatDateTime(comment.createdAt)}</time></header><p>{comment.content}</p><div className="front-comment-actions"><button type="button" className={`front-comment-like${comment.isLiked ? ' liked' : ''}`} disabled={likingCommentIDs[comment.id]} onClick={() => void toggleCommentLike(comment)} aria-label={comment.isLiked ? '取消点赞评论' : '点赞评论'} title={comment.isLiked ? '取消点赞' : '点赞'}><ActionIcon name="like" /><span>{comment.likeCount}</span></button><button type="button" onClick={() => loggedIn ? setReplyTarget(comment) : navigate('/login')}>回复</button>{comment.childCount > 0 && <button type="button" onClick={() => void toggleReplies(comment)}>{expandedReplies[comment.id] ? '收起回复' : `查看 ${comment.childCount} 条回复`}</button>}</div>{expandedReplies[comment.id] && <div className="front-replies">{replies[comment.id] ? replies[comment.id].map((reply) => <article key={reply.id}><Avatar src={reply.userAvatar} name={reply.userName || 'R'} userID={reply.userId} small /><div><header><strong>{reply.userName || 'Renai 成员'}</strong><time>{formatDateTime(reply.createdAt)}</time></header><p>{reply.replyToUserName && <span className="front-reply-to">回复 @{reply.replyToUserName}：</span>}{reply.content}</p><div className="front-comment-actions front-reply-actions"><button type="button" className={`front-comment-like${reply.isLiked ? ' liked' : ''}`} disabled={likingCommentIDs[reply.id]} onClick={() => void toggleCommentLike(reply)} aria-label={reply.isLiked ? '取消点赞评论' : '点赞评论'} title={reply.isLiked ? '取消点赞' : '点赞'}><ActionIcon name="like" /><span>{reply.likeCount}</span></button><button type="button" onClick={() => loggedIn ? setReplyTarget(reply) : navigate('/login')}>回复</button></div></div></article>) : <span>正在加载回复...</span>}</div>}</div></article>)}</div>}
      </section>
    </article>
    <aside className="front-detail-aside">
      <nav aria-label="文章目录"><span>文章导航</span><h2>目录</h2>{headings.length ? headings.map((heading, index) => <button type="button" className={`level-${heading.level}`} key={`${heading.text}-${index}`} onClick={() => jumpToHeading(index)}>{heading.text}</button>) : <p>文章暂无章节标题</p>}</nav>
      <div className="front-article-actions" aria-label="文章互动"><span className="front-action-item" title={`${article.viewCount} 次浏览`}><ActionIcon name="view" /><b>{article.viewCount}</b></span><a className="front-action-item" href="#comments" aria-label={`${article.commentCount} 条评论`} title={`${article.commentCount} 条评论`}><ActionIcon name="comment" /><b>{article.commentCount}</b></a><button type="button" className={liked ? 'liked' : ''} disabled={liking} onClick={() => void toggleArticleLike()} aria-label={`${article.likeCount} 个点赞`} title={liked ? '取消点赞' : '点赞'}><ActionIcon name="like" /><b>{article.likeCount}</b></button><button type="button" className={favored ? 'favored' : ''} disabled={favoring} onClick={() => void openFavorPicker()} aria-label={favored ? '取消收藏文章' : '收藏文章'} title={favored ? '取消收藏' : '收藏到收藏夹'}><ActionIcon name="favor" /><b>{article.favorCount}</b></button><button type="button" className="front-action-top" onClick={() => window.scrollTo({ top: 0, behavior: 'smooth' })} aria-label="返回顶部" title="返回顶部"><ActionIcon name="top" /></button></div>
    </aside>
    {replyTarget && <div className="front-reply-overlay" role="presentation" onMouseDown={(event) => { if (event.target === event.currentTarget) setReplyTarget(null) }}><form className="front-reply-dialog" role="dialog" aria-modal="true" aria-labelledby="reply-title" onSubmit={submitReply}><header><div><span>回复</span><strong id="reply-title">{replyTarget.userName}</strong></div><button type="button" aria-label="关闭回复" onClick={() => setReplyTarget(null)}>×</button></header><blockquote>{replyTarget.content}</blockquote><label htmlFor="reply-content">回复内容</label><textarea id="reply-content" autoFocus value={replyText} maxLength={1000} onChange={(event) => setReplyText(event.target.value)} /><footer><span>{replyText.length}/1000</span><button type="submit" disabled={submitting || !replyText.trim()}>{submitting ? '发布中...' : '发布回复'}</button></footer></form></div>}
    {favorPickerOpen && <div className="front-reply-overlay" role="presentation" onMouseDown={(event) => { if (event.target === event.currentTarget) setFavorPickerOpen(false) }}><section className="front-reply-dialog favor-picker" role="dialog" aria-modal="true" aria-label="选择收藏夹"><header><div><span>收藏文章</span><strong>选择收藏夹</strong></div><button type="button" aria-label="关闭收藏夹" onClick={() => setFavorPickerOpen(false)}>×</button></header><div className="favor-folder-list">{folders.map((folder) => <button type="button" key={folder.id} onClick={() => void toggleArticleFavor(folder.id)}><span>{folder.name}{folder.isDefault ? '（默认）' : ''}</span><small>{folder.articleCount} 篇</small></button>)}</div><footer className="favor-folder-create"><input value={newFolderName} maxLength={50} placeholder="新建收藏夹名称" onChange={(event) => setNewFolderName(event.target.value)} /><button type="button" onClick={() => void createFolder()} disabled={!newFolderName.trim() || favoring}>新建收藏夹</button></footer></section></div>}
  </main>
}
