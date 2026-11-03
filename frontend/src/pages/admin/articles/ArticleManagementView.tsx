import { useEffect, useRef, useState } from 'react'
import { articleApi } from '@/shared/api/v1/article'
import { categoryApi } from '@/shared/api/v1/category'
import { navigate } from '@/shared/hooks/usePathname'
import type { Article } from '@/shared/types/article'
import type { Category } from '@/shared/types/category'
import { logRequestError } from '@/shared/lib/request-error'
import { OTHER_CATEGORY_NAME } from '@/shared/lib/category'
import './ArticlesView.css'

export default function ArticleManagementView() {
  const [articles, setArticles] = useState<Article[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [keyword, setKeyword] = useState('')
  const [categoryID, setCategoryID] = useState('')
  const [loading, setLoading] = useState(true)
  const [message, setMessage] = useState('')
  const [categoryOpen, setCategoryOpen] = useState(false)
  const pickerRef = useRef<HTMLDivElement>(null)

  const load = async (all = false) => {
    try {
      setLoading(true)
      const params = { page: 1, pageSize: 100 }
      const result = all ? await articleApi.list(params) : keyword.trim() ? await articleApi.search(keyword.trim(), params) : categoryID ? await articleApi.byCategory(Number(categoryID), params) : await articleApi.list(params)
      setArticles(result.articles ?? [])
    } catch (error) { logRequestError('加载文章管理列表失败', error); setMessage('加载失败，请稍后重试') } finally { setLoading(false) }
  }

  useEffect(() => {
    void load()
    void categoryApi.listTypes().then(async ({ types }) => Promise.all((types ?? []).map((type) => categoryApi.listCategories({ parentID: type.id })))).then((items) => setCategories(items.flatMap((item) => item.categories ?? []))).catch((error: unknown) => { logRequestError('加载文章分类失败', error); setMessage('加载分类失败，请稍后重试') })
  }, [])
  useEffect(() => { const close = (event: MouseEvent) => { if (!pickerRef.current?.contains(event.target as Node)) setCategoryOpen(false) }; document.addEventListener('mousedown', close); return () => document.removeEventListener('mousedown', close) }, [])

  const selectCategory = (id: string) => { setCategoryID(id); setKeyword(''); setCategoryOpen(false) }
  const selected = categories.find((item) => item.id === Number(categoryID))
  const selectedName = categoryID === '0' ? OTHER_CATEGORY_NAME : selected?.name || '全部分类'
  const remove = async (article: Article) => { if (!window.confirm(`确定要删除《${article.title}》吗？\n\n删除后无法恢复。`)) return; try { await articleApi.delete({ id: article.id }); setArticles((items) => items.filter((item) => item.id !== article.id)) } catch (error) { logRequestError('删除文章失败', error); setMessage('删除失败，请稍后重试') } }
  const clear = () => { setKeyword(''); setCategoryID(''); void load(true) }

  return <section className="articles-view"><form className="articles-filters" onSubmit={(event) => { event.preventDefault(); void load() }}><input value={keyword} onChange={(event) => setKeyword(event.target.value)} placeholder="搜索标题、摘要或正文" /><div className="category-picker" ref={pickerRef}><button className="category-picker-trigger" type="button" aria-haspopup="listbox" aria-expanded={categoryOpen} onClick={() => setCategoryOpen((value) => !value)}><span>{selectedName}</span><i aria-hidden="true">⌄</i></button>{categoryOpen && <div className="category-picker-menu" role="listbox"><button type="button" className={!categoryID ? 'selected' : ''} onClick={() => selectCategory('')}><span>全部分类</span>{!categoryID && <b>✓</b>}</button><button type="button" className={categoryID === '0' ? 'selected' : ''} onClick={() => selectCategory('0')}><span>{OTHER_CATEGORY_NAME}</span>{categoryID === '0' && <b>✓</b>}</button>{categories.length ? categories.map((item) => <button type="button" key={item.id} className={String(item.id) === categoryID ? 'selected' : ''} onClick={() => selectCategory(String(item.id))}><span>{item.name}</span>{String(item.id) === categoryID && <b>✓</b>}</button>) : <span className="category-picker-empty">暂无可用分类</span>}</div>}</div><button className="primary-button" type="submit">搜索</button><button className="articles-reset" type="button" onClick={clear}>清除筛选</button></form>{message && <p className="articles-message">{message}</p>}{loading ? <div className="articles-empty">正在加载文章...</div> : articles.length === 0 ? <div className="articles-empty"><strong>没有找到文章</strong></div> : <div className="article-list">{articles.map((article) => <article className="article-row" key={article.id} onClick={() => navigate(`/blog/${article.id}`)}><div className="article-cover">{article.coverImage ? <img src={article.coverImage} alt="文章封面" /> : <span>NO COVER</span>}</div><div className="article-info"><h3>{article.title}</h3><p>{article.summary || '暂无摘要'}</p><div className="article-meta"><span>{article.authorName || '未知作者'}</span><span>浏览 {article.viewCount}</span><span>评论 {article.commentCount}</span></div></div><button className="article-delete" type="button" onClick={(event) => { event.stopPropagation(); void remove(article) }}>删除</button></article>)}</div>}</section>
}
