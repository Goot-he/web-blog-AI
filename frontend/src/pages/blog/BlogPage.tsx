import { useEffect, useMemo, useRef, useState } from 'react'
import { articleApi } from '@/shared/api/v1/article'
import { categoryApi } from '@/shared/api/v1/category'
import AppLink from '@/shared/ui/AppLink'
import type { Article } from '@/shared/types/article'
import type { Category, CategoryType } from '@/shared/types/category'
import { resolveStorageUrl } from '@/shared/lib/storage'
import { logRequestError } from '@/shared/lib/request-error'
import { getCategoryName, OTHER_CATEGORY_NAME } from '@/shared/lib/category'
import './BlogPage.css'

const PAGE_SIZE = 12
type CategoryGroup = { type: CategoryType; categories: Category[] }
type SelectedCategory = number | null
const formatDate = (value: number) => new Date(value > 1_000_000_000_000 ? value : value * 1000).toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })

export default function BlogPage() {
  const [groups, setGroups] = useState<CategoryGroup[]>([])
  const [selectedType, setSelectedType] = useState(0)
  const [selectedCategory, setSelectedCategory] = useState<SelectedCategory>(null)
  const [articles, setArticles] = useState<Article[]>([])
  const [page, setPage] = useState(1)
  const [hasMore, setHasMore] = useState(false)
  const [query, setQuery] = useState('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const secondaryScrollerRef = useRef<HTMLDivElement>(null)

  const applyQueryFromUrl = () => {
    const nextQuery = new URLSearchParams(window.location.search).get('q')?.trim() ?? ''
    setPage(1)
    setSelectedType(0)
    setSelectedCategory(null)
    setQuery(nextQuery)
  }

  useEffect(() => {
    void categoryApi.listTypes().then(async ({types}) => {
      const children = await Promise.all(types.map((type) => categoryApi.listCategories({parentID: type.id})));
      setGroups(types.map((type, index) => ({type, categories: children[index].categories ?? []})))
    }).catch((reason: unknown) => {
      logRequestError('加载文章分类失败', reason);
      setError('加载失败，请稍后重试')
    })
  }, [])
  useEffect(() => {
    applyQueryFromUrl()
    window.addEventListener('popstate', applyQueryFromUrl)
    window.addEventListener('renai:blog-search', applyQueryFromUrl)
    return () => {
      window.removeEventListener('popstate', applyQueryFromUrl)
      window.removeEventListener('renai:blog-search', applyQueryFromUrl)
    }
  }, [])
  useEffect(() => {
    let active = true;
    setLoading(true);
    setError('');
    const params = {page, pageSize: PAGE_SIZE};
    const selectedChildren = groups.find((group) => group.type.id === selectedType)?.categories ?? [];
    const request = query ? articleApi.search(query, params) : selectedCategory !== null ? articleApi.byCategory(selectedCategory, params) : selectedType ? Promise.all(selectedChildren.map((category) => articleApi.byCategory(category.id, params))).then((responses) => ({articles: responses.flatMap((response) => response.articles ?? [])})) : articleApi.list(params);
    void request.then(({articles: nextArticles}) => {
      if (!active) return;
      const next = nextArticles ?? [];
      setArticles((current) => page === 1 ? next : [...current, ...next]);
      setHasMore(next.length === PAGE_SIZE)
    }).catch((reason: unknown) => {
      logRequestError('加载文章列表失败', reason);
      if (active) setError('加载失败，请稍后重试')
    }).finally(() => {
      if (active) setLoading(false)
    });
    return () => {
      active = false
    }
  }, [page, query, selectedCategory, selectedType, groups])

  const categoryNames = useMemo(() => new Map(groups.flatMap((group) => group.categories.map((item) => [item.id, item.name] as const))), [groups])
  const activeChildren = groups.find((group) => group.type.id === selectedType)?.categories ?? []
  const resetAnd = (action: () => void) => {
    setPage(1);
    action()
  }
  const clearFilters = () => resetAnd(() => {
    setSelectedType(0);
    setSelectedCategory(null);
    setQuery('')
  })
  const scrollSecondary = (direction: -1 | 1) => {
    secondaryScrollerRef.current?.scrollBy({ left: direction * 320, behavior: 'smooth' })
  }

  return <main className="front-blog-page">
    <section className="front-channel-bar" aria-label="文章分类">
      <div className="front-channel-primary">
        <button type="button" className="front-channel-hot" disabled>热门</button>
        <div className="front-channel-list">
          <button type="button" className={!selectedType && selectedCategory === null ? 'active' : ''} onClick={clearFilters}>全部</button>
          <button type="button" className={!selectedType && selectedCategory === 0 ? 'active' : ''} onClick={() => resetAnd(() => {
            setSelectedType(0)
            setSelectedCategory(0)
            setQuery('')
          })}>{OTHER_CATEGORY_NAME}</button>
          {groups.map((group) => <button type="button" className={selectedType === group.type.id ? 'active' : ''} key={group.type.id} onClick={() => resetAnd(() => {
            setSelectedType(group.type.id)
            setSelectedCategory(null)
            setQuery('')
          })}>{group.type.name}</button>)}
        </div>
      </div>
      {selectedType > 0 && <div className="front-channel-secondary">
        <button type="button" className="front-channel-arrow" aria-label="向左查看更多二级分类" onClick={() => scrollSecondary(-1)}><svg viewBox="0 0 24 24" aria-hidden="true"><path d="M15 6 9 12l6 6" /></svg></button>
        <div className="front-channel-scroll" ref={secondaryScrollerRef}>
          <button type="button" className={selectedCategory === null ? 'active' : ''} onClick={() => resetAnd(() => setSelectedCategory(null))}>全部细分</button>
          {activeChildren.map((category) => <button type="button" className={selectedCategory === category.id ? 'active' : ''} key={category.id} onClick={() => resetAnd(() => setSelectedCategory(category.id))}>{category.name}</button>)}
        </div>
        <button type="button" className="front-channel-arrow" aria-label="向右查看更多二级分类" onClick={() => scrollSecondary(1)}><svg viewBox="0 0 24 24" aria-hidden="true"><path d="m9 6 6 6-6 6" /></svg></button>
      </div>}
    </section>
    <section className="front-blog-results" aria-live="polite" aria-busy={loading}>
      {error && <div className="front-blog-state error"><strong>{error}</strong>
      </div>}{!error && loading && page === 1 && articles.length === 0 &&
        <div className="front-blog-state">正在加载文章...</div>}{!error && !loading && articles.length === 0 &&
        <div className="front-blog-state"><strong>暂时没有找到文章</strong><span>换个分类或关键词试试。</span></div>}
      <div
          className={`front-article-grid${loading && articles.length > 0 ? ' is-refreshing' : ''}`}>{articles.map((article) =>
          <AppLink className="front-article-card" to={`/blog/${article.id}`} key={article.id}>
            <div className="front-article-cover">{article.coverImage ?
                <img src={resolveStorageUrl(article.coverImage)} alt={`${article.title}封面`} loading="lazy"/> :
                <span>{article.title.slice(0, 1)}</span>}</div>
            <div className="front-article-copy">
              <h3>{article.title}</h3>
              <footer><span className="front-author front-author-text">作者：{article.authorName || 'Renai 成员'}</span><span className="front-card-meta">分类：{getCategoryName(article.categoryID, categoryNames)}</span><time className="front-card-meta">{formatDate(article.createdAt)}</time>
              </footer>
            </div>
          </AppLink>)}</div>
      {hasMore && !error && <button type="button" className="front-load-more" disabled={loading}
                                    onClick={() => setPage((current) => current + 1)}>{loading ? '加载中...' : '加载更多'}</button>}
    </section>
  </main>
};
