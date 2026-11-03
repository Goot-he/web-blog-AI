import { useEffect, useMemo, useRef, useState, type FocusEvent, type ReactNode } from 'react'
import { articleApi } from '@/shared/api/v1/article'
import { categoryApi } from '@/shared/api/v1/category'
import type { Category, CategoryType } from '@/shared/types/category'
import type { Article, ArticleImageUploadResponse } from '@/shared/types/article'
import { logRequestError } from '@/shared/lib/request-error'
import './CreateArticleView.css'

type Point = { x: number; y: number }
type ImageSize = { width: number; height: number }
type EditorBlockType = 'paragraph' | 'heading1' | 'heading2' | 'heading3' | 'heading4' | 'heading5' | 'quote' | 'unordered' | 'ordered' | 'divider' | 'image'
type EditorBlock = { id: string; type: EditorBlockType; text: string; url?: string; alt?: string; bold?: boolean; italic?: boolean; code?: boolean }
type FormatOption = { type: EditorBlockType; label: string; hint: string }

const COVER_RATIO = 16 / 9
const readEditingArticle = () => {
  const state = window.history.state as { article?: Article } | null
  return state?.article ?? null
}
const createBlockID = () => `block-${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 8)}`
const emptyBlock = (): EditorBlock => ({ id: createBlockID(), type: 'paragraph', text: '' })
const imagePattern = /^!\[([^\]]*)]\(([^)]+)\)\s*$/
const headingPattern = /^(#{1,5})\s+(.+)$/
const formatOptions: FormatOption[] = [
  { type: 'paragraph', label: '正文', hint: '普通段落' },
  { type: 'heading1', label: '标题一', hint: '# 标题' },
  { type: 'heading2', label: '标题二', hint: '## 标题' },
  { type: 'heading3', label: '标题三', hint: '### 标题' },
  { type: 'heading4', label: '标题四', hint: '#### 标题' },
  { type: 'heading5', label: '标题五', hint: '##### 标题' },
]

const applyInlineMarks = (block: EditorBlock, text: string) => {
  if (!text) return ''
  if (block.code) return `\`${text}\``
  if (block.bold && block.italic) return `***${text}***`
  if (block.bold) return `**${text}**`
  if (block.italic) return `*${text}*`
  return text
}

const extractInlineMarks = (text: string): Pick<EditorBlock, 'text' | 'bold' | 'italic' | 'code'> => {
  const trimmed = text.trim()
  if (/^`[^`]+`$/.test(trimmed)) return { text: trimmed.slice(1, -1), code: true }
  if (/^\*\*\*[\s\S]+\*\*\*$/.test(trimmed)) return { text: trimmed.slice(3, -3), bold: true, italic: true }
  if (/^\*\*[\s\S]+\*\*$/.test(trimmed)) return { text: trimmed.slice(2, -2), bold: true }
  if (/^\*[^*][\s\S]*\*$/.test(trimmed)) return { text: trimmed.slice(1, -1), italic: true }
  return { text }
}

const blockToMarkdown = (block: EditorBlock) => {
  const text = applyInlineMarks(block, block.text.trimEnd())
  switch (block.type) {
    case 'heading1': return text ? `# ${text}` : ''
    case 'heading2': return text ? `## ${text}` : ''
    case 'heading3': return text ? `### ${text}` : ''
    case 'heading4': return text ? `#### ${text}` : ''
    case 'heading5': return text ? `##### ${text}` : ''
    case 'quote': return text ? text.split('\n').map((line) => `> ${line}`).join('\n') : ''
    case 'unordered': return text ? text.split('\n').map((line) => `- ${line}`).join('\n') : ''
    case 'ordered': return text ? text.split('\n').map((line, index) => `${index + 1}. ${line}`).join('\n') : ''
    case 'divider': return '---'
    case 'image': return block.url ? `![${block.alt || block.text || 'image'}](${block.url})` : ''
    default: return text
  }
}

const blocksToMarkdown = (blocks: EditorBlock[]) => blocks.map(blockToMarkdown).filter(Boolean).join('\n\n')
const blockRows = (text: string) => text.split('\n').reduce((total, line) => total + Math.max(1, Math.ceil(line.length / 62)), 0)

const markdownToBlocks = (markdown: string): EditorBlock[] => {
  const blocks: EditorBlock[] = []
  const lines = markdown.replace(/\r\n/g, '\n').split('\n')
  let paragraph: string[] = []
  const flushParagraph = () => {
    if (!paragraph.length) return
    blocks.push({ id: createBlockID(), type: 'paragraph', ...extractInlineMarks(paragraph.join('\n').trimEnd()) })
    paragraph = []
  }

  for (const line of lines) {
    const trimmed = line.trim()
    if (!trimmed) { flushParagraph(); continue }
    const image = line.match(imagePattern)
    const heading = line.match(headingPattern)
    const quote = line.match(/^>\s?(.*)$/)
    const unordered = line.match(/^[-*]\s+(.+)$/)
    const ordered = line.match(/^\d+\.\s+(.+)$/)
    const divider = /^-{3,}$/.test(trimmed)
    if (image) {
      flushParagraph()
      blocks.push({ id: createBlockID(), type: 'image', text: image[1] || 'image', alt: image[1] || 'image', url: image[2] })
    } else if (divider) {
      flushParagraph()
      blocks.push({ id: createBlockID(), type: 'divider', text: '' })
    } else if (heading) {
      flushParagraph()
      blocks.push({ id: createBlockID(), type: `heading${heading[1].length}` as EditorBlockType, ...extractInlineMarks(heading[2]) })
    } else if (quote) {
      flushParagraph()
      blocks.push({ id: createBlockID(), type: 'quote', ...extractInlineMarks(quote[1]) })
    } else if (unordered) {
      flushParagraph()
      blocks.push({ id: createBlockID(), type: 'unordered', ...extractInlineMarks(unordered[1]) })
    } else if (ordered) {
      flushParagraph()
      blocks.push({ id: createBlockID(), type: 'ordered', ...extractInlineMarks(ordered[1]) })
    } else {
      paragraph.push(line)
    }
  }
  flushParagraph()
  return blocks.length ? blocks : [emptyBlock()]
}

function MarkdownRichEditor({
  value,
  busy,
  onChange,
  onUploadImage,
  onRemoveImage,
  onNotice,
  onSaveDraft,
  onPublish,
  draftLabel,
  publishLabel,
  stats,
  header,
}: {
  value: string
  busy: string
  onChange: (value: string) => void
  onUploadImage: (file?: File) => Promise<ArticleImageUploadResponse | null>
  onRemoveImage: (url: string) => void
  onNotice: (value: string) => void
  onSaveDraft: () => void
  onPublish: () => void
  draftLabel: string
  publishLabel: string
  stats: { words: number; headings: number; images: number }
  header: ReactNode
}) {
  const [blocks, setBlocks] = useState<EditorBlock[]>(() => markdownToBlocks(value))
  const [activeID, setActiveID] = useState('')
  const [formatOpen, setFormatOpen] = useState(false)
  const [sourceMode, setSourceMode] = useState(false)
  const inputRefs = useRef(new Map<string, HTMLTextAreaElement>())
  const lastMarkdownRef = useRef(blocksToMarkdown(blocks))

  useEffect(() => {
    if (value === lastMarkdownRef.current) return
    const next = markdownToBlocks(value)
    setBlocks(next)
    lastMarkdownRef.current = blocksToMarkdown(next)
  }, [value])

  const commit = (next: EditorBlock[]) => {
    setBlocks(next)
    const markdown = blocksToMarkdown(next)
    lastMarkdownRef.current = markdown
    onChange(markdown)
  }

  const focusBlock = (id: string) => {
    requestAnimationFrame(() => {
      const input = inputRefs.current.get(id)
      input?.focus()
      input?.setSelectionRange(input.value.length, input.value.length)
    })
  }

  const updateBlock = (id: string, patch: Partial<EditorBlock>) => {
    commit(blocks.map((block) => block.id === id ? { ...block, ...patch } : block))
  }

  const insertBlockAfter = (id: string, block: EditorBlock) => {
    const index = blocks.findIndex((item) => item.id === id)
    const next = [...blocks]
    next.splice(index >= 0 ? index + 1 : next.length, 0, block)
    commit(next)
    focusBlock(block.id)
  }

  const removeBlock = (id: string) => {
    const removed = blocks.find((block) => block.id === id)
    const next = blocks.filter((block) => block.id !== id)
    if (removed?.url) onRemoveImage(removed.url)
    commit(next.length ? next : [emptyBlock()])
  }

  const activeBlock = blocks.find((block) => block.id === activeID && block.type !== 'image' && block.type !== 'divider') ?? blocks.find((block) => block.type !== 'image' && block.type !== 'divider')
  const changeActiveType = (type: EditorBlockType) => {
    if (!activeBlock) return
    updateBlock(activeBlock.id, { type })
    focusBlock(activeBlock.id)
  }

  const selectFormat = (type: EditorBlockType) => {
    changeActiveType(type)
    setFormatOpen(false)
  }

  const closeFormatMenu = (event: FocusEvent<HTMLDivElement>) => {
    if (!event.currentTarget.contains(event.relatedTarget as Node | null)) setFormatOpen(false)
  }

  const wrapSelection = (prefix: string, suffix = prefix) => {
    if (!activeBlock) return
    const input = inputRefs.current.get(activeBlock.id)
    const start = input?.selectionStart ?? activeBlock.text.length
    const end = input?.selectionEnd ?? activeBlock.text.length
    if (start === end) return false
    const selected = activeBlock.text.slice(start, end)
    const nextText = `${activeBlock.text.slice(0, start)}${prefix}${selected}${suffix}${activeBlock.text.slice(end)}`
    updateBlock(activeBlock.id, { text: nextText })
    requestAnimationFrame(() => inputRefs.current.get(activeBlock.id)?.setSelectionRange(start + prefix.length, start + prefix.length + selected.length))
    return true
  }

  const toggleInlineMark = (mark: 'bold' | 'italic' | 'code') => {
    if (!activeBlock) return
    const wrapped = mark === 'bold' ? wrapSelection('**') : mark === 'italic' ? wrapSelection('*') : wrapSelection('`')
    if (wrapped) return
    const patch: Partial<EditorBlock> = { [mark]: !activeBlock[mark] }
    if (mark === 'code' && !activeBlock.code) Object.assign(patch, { bold: false, italic: false })
    if (mark !== 'code' && activeBlock.code) patch.code = false
    updateBlock(activeBlock.id, patch)
    focusBlock(activeBlock.id)
  }

  const handleKeyDown = (event: React.KeyboardEvent<HTMLTextAreaElement>, block: EditorBlock) => {
    if (event.key === 'Enter' && !event.shiftKey) {
      event.preventDefault()
      insertBlockAfter(block.id, emptyBlock())
    }
    if (event.key === 'Backspace' && !block.text && blocks.length > 1) {
      event.preventDefault()
      removeBlock(block.id)
    }
  }

  const insertImage = async (file?: File) => {
    const result = await onUploadImage(file)
    if (!result) return
    const block = { id: createBlockID(), type: 'image' as const, text: file?.name ?? 'image', alt: file?.name ?? 'image', url: result.url }
    const targetID = activeID || blocks.at(-1)?.id || ''
    insertBlockAfter(targetID, block)
  }

  const importMarkdown = async (file?: File) => {
    if (!file) return
    if (!/(\.md$|\.markdown$|text\/)/i.test(`${file.name} ${file.type}`)) {
      onNotice('请选择 Markdown 文件')
      return
    }
    const text = await file.text()
    const next = markdownToBlocks(text)
    commit(next)
    onNotice('Markdown 文件已导入')
  }

  const insertDivider = () => {
    const targetID = activeBlock?.id || activeID || blocks.at(-1)?.id || ''
    insertBlockAfter(targetID, { id: createBlockID(), type: 'divider', text: '' })
  }

  let cursor = 0
  return <div className="rich-editor" ref={(node) => { if (node) node.dataset.ready = 'true' }}>
    <div className="rich-toolbar" aria-label="文章编辑工具">
      <div className="rich-toolbar-group rich-format-picker" onBlur={closeFormatMenu}>
        <button type="button" className="rich-format-trigger" aria-haspopup="listbox" aria-expanded={formatOpen} onClick={() => setFormatOpen((open) => !open)}>
          <span><b>H</b><em>格式</em></span><i aria-hidden="true" />
        </button>
        {formatOpen && <div className="rich-format-menu" role="listbox" aria-label="标题格式">
          {formatOptions.map((option) => <button type="button" role="option" aria-selected={activeBlock?.type === option.type} className={activeBlock?.type === option.type ? 'selected' : ''} key={option.type} onMouseDown={(event) => event.preventDefault()} onClick={() => selectFormat(option.type)}>
            <span>{option.label}</span><small>{option.hint}</small>
          </button>)}
        </div>}
      </div>
      <div className="rich-toolbar-group">
        <button type="button" className={`rich-tool-button${activeBlock?.bold ? ' active' : ''}`} title="加粗" onClick={() => toggleInlineMark('bold')}><span>B</span><em>加粗</em></button>
        <button type="button" className={`rich-tool-button${activeBlock?.italic ? ' active' : ''}`} title="斜体" onClick={() => toggleInlineMark('italic')}><span>I</span><em>斜体</em></button>
        <button type="button" className={`rich-tool-button${activeBlock?.code ? ' active' : ''}`} title="行内代码" onClick={() => toggleInlineMark('code')}><span>&lt;/&gt;</span><em>代码</em></button>
        <button type="button" className="rich-tool-button" title="引用" onClick={() => changeActiveType('quote')}><span>“</span><em>引用</em></button>
        <button type="button" className="rich-tool-button" title="无序列表" onClick={() => changeActiveType('unordered')}><span>☷</span><em>列表</em></button>
        <button type="button" className="rich-tool-button" title="水平线" onClick={insertDivider}><span>—</span><em>水平线</em></button>
      </div>
      <div className="rich-toolbar-group">
        <button type="button" className="rich-tool-button" title="新增段落" onClick={() => activeBlock && insertBlockAfter(activeBlock.id, emptyBlock())}><span>＋</span><em>段落</em></button>
        <label className="rich-tool-button" title="插入图片"><input type="file" accept="image/jpeg,image/png,image/gif,image/webp,image/avif" onChange={(event) => { void insertImage(event.target.files?.[0]); event.currentTarget.value = '' }} disabled={busy !== ''} /><span>▧</span><em>{busy === 'content' ? '上传中' : '图像'}</em></label>
        <button type="button" className={`rich-tool-button${sourceMode ? ' active' : ''}`} title="查看 Markdown 格式" onClick={() => setSourceMode((current) => !current)}><span>MD</span><em>{sourceMode ? '富文本' : '查看 MD'}</em></button>
        <label className="rich-tool-button" title="导入 Markdown"><input type="file" accept=".md,.markdown,text/markdown,text/plain" onChange={(event) => { void importMarkdown(event.target.files?.[0]); event.currentTarget.value = '' }} /><span>⇪</span><em>导入 MD</em></label>
      </div>
      <div className="rich-toolbar-spacer" />
      <div className="rich-toolbar-stats"><span>{stats.words} 字</span><span>{stats.headings} 章</span><span>{stats.images} 图</span></div>
      <div className="rich-toolbar-actions">
        <button type="button" className="rich-draft-action" disabled={busy !== ''} onClick={onSaveDraft}><span>▣</span>{busy === 'draft' ? '保存中...' : draftLabel}</button>
        <button type="button" className="rich-publish-action" disabled={busy !== ''} onClick={onPublish}><span>↗</span>{busy === 'publish' ? '发布中...' : publishLabel}</button>
      </div>
    </div>
    {header}
    {sourceMode ? <textarea className="rich-source-editor" value={value} spellCheck={false} onChange={(event) => onChange(event.target.value)} /> : <div className="rich-canvas" data-empty={blocks.length === 1 && !blocks[0].text && blocks[0].type === 'paragraph'}>
      {blocks.map((block) => {
        const source = blockToMarkdown(block)
        const headingIndex = block.type.startsWith('heading') ? cursor : undefined
        cursor += source.length + (source ? 2 : 0)
        return block.type === 'divider'
          ? <div className="rich-divider-block" key={block.id}>
            <hr />
            <button type="button" onClick={() => removeBlock(block.id)}>删除水平线</button>
          </div>
          : block.type === 'image'
          ? <figure className="rich-image-block" key={block.id}>
            <img src={block.url} alt={block.alt || '文章图片'} />
            <figcaption>{block.alt || '文章图片'}</figcaption>
            <button type="button" onClick={() => removeBlock(block.id)}>删除图片</button>
          </figure>
          : <div className={`rich-block rich-${block.type}${block.bold ? ' is-bold' : ''}${block.italic ? ' is-italic' : ''}${block.code ? ' is-code' : ''}`} key={block.id} data-heading-index={headingIndex}>
            <textarea
              ref={(node) => { if (node) inputRefs.current.set(block.id, node); else inputRefs.current.delete(block.id) }}
              value={block.text}
              rows={Math.max(1, blockRows(block.text))}
              placeholder={block.type === 'paragraph' ? '输入正文，或从上方选择标题、引用、列表...' : ''}
              onFocus={() => setActiveID(block.id)}
              onChange={(event) => updateBlock(block.id, { text: event.target.value })}
              onKeyDown={(event) => handleKeyDown(event, block)}
            />
            {block.type === 'quote' && <span aria-hidden="true" />}
          </div>
      })}
    </div>}
  </div>
}

export default function CreateArticleView() {
  const [title, setTitle] = useState('')
  const [summary, setSummary] = useState('')
  const [content, setContent] = useState('')
  const [categoryID, setCategoryID] = useState(0)
  const [categoryGroups, setCategoryGroups] = useState<Array<{ type: CategoryType; categories: Category[] }>>([])
  const [categoryTypeID, setCategoryTypeID] = useState(0)
  const [openCategoryMenu, setOpenCategoryMenu] = useState<'type' | 'category' | null>(null)
  const [collapsedOutline, setCollapsedOutline] = useState<Set<number>>(() => new Set())
  const [cover, setCover] = useState<ArticleImageUploadResponse | null>(null)
  const [contentImages, setContentImages] = useState<ArticleImageUploadResponse[]>([])
  const [busy, setBusy] = useState<'cover' | 'content' | 'publish' | 'draft' | ''>('')
  const [notice, setNotice] = useState('')
  const [editingArticle, setEditingArticle] = useState<Article | null>(() => readEditingArticle())
  const [cropUrl, setCropUrl] = useState<string | null>(null)
  const [cropFile, setCropFile] = useState<File | null>(null)
  const [cropSize, setCropSize] = useState<ImageSize | null>(null)
  const [cropScale, setCropScale] = useState(1)
  const [cropOffset, setCropOffset] = useState<Point>({ x: 0, y: 0 })
  const editorRef = useRef<HTMLDivElement>(null)
  const cropStageRef = useRef<HTMLDivElement>(null)
  const cropPreviewRef = useRef<HTMLDivElement>(null)
  const cropImageRef = useRef<HTMLImageElement>(null)
  const dragStartRef = useRef<{ point: Point; offset: Point } | null>(null)
  const categorySelectsRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    void categoryApi.listTypes().then(async ({ types }) => {
      const groups = await Promise.all(types.map((item) => categoryApi.listCategories({ parentID: item.id })))
      setCategoryGroups(types.map((type, index) => ({ type, categories: groups[index].categories ?? [] })))
    }).catch((error: unknown) => { logRequestError('加载文章分类失败', error); setNotice('加载分类失败，请稍后重试') })
  }, [])

  useEffect(() => {
    if (!editingArticle) return
    setTitle(editingArticle.title ?? '')
    setSummary(editingArticle.summary ?? '')
    setContent(editingArticle.content ?? '')
    setCategoryID(editingArticle.categoryID ?? 0)
    setCover(editingArticle.coverImage ? { url: editingArticle.coverImage } : null)
  }, [editingArticle])

  useEffect(() => {
    if (!editingArticle || !categoryGroups.length) return
    const group = categoryGroups.find((item) => item.categories.some((category) => category.id === editingArticle.categoryID))
    if (group) setCategoryTypeID(group.type.id)
  }, [categoryGroups, editingArticle])

  useEffect(() => () => { if (cropUrl) URL.revokeObjectURL(cropUrl) }, [cropUrl])

  useEffect(() => {
    const closeMenu = (event: PointerEvent) => { if (!categorySelectsRef.current?.contains(event.target as Node)) setOpenCategoryMenu(null) }
    document.addEventListener('pointerdown', closeMenu)
    return () => document.removeEventListener('pointerdown', closeMenu)
  }, [])

  const headings = useMemo(() => Array.from(content.matchAll(/^(#{1,5})\s+(.+)$/gm)).map((match) => ({ level: match[1].length, text: match[2].trim(), index: match.index ?? 0 })), [content])
  const outlineItems = useMemo(() => headings.map((heading, index) => ({ ...heading, hasChildren: (headings[index + 1]?.level ?? 0) > heading.level })), [headings])
  const visibleOutlineItems = useMemo(() => {
    const collapsedLevels: number[] = []
    return outlineItems.filter((heading) => {
      while (collapsedLevels.length && collapsedLevels[collapsedLevels.length - 1] >= heading.level) collapsedLevels.pop()
      const hidden = collapsedLevels.length > 0
      if (!hidden && heading.hasChildren && collapsedOutline.has(heading.index)) collapsedLevels.push(heading.level)
      return !hidden
    })
  }, [collapsedOutline, outlineItems])
  const wordCount = content.replace(/\s/g, '').length

  const closeCrop = () => {
    setCropUrl(null)
    setCropFile(null)
    setCropSize(null)
    setCropScale(1)
    setCropOffset({ x: 0, y: 0 })
  }

  const openCrop = (file?: File) => {
    if (!file) return
    if (!file.type.startsWith('image/')) { setNotice('请选择图片文件'); return }
    setCropFile(file)
    setCropUrl(URL.createObjectURL(file))
    setCropSize(null)
    setCropScale(1)
    setCropOffset({ x: 0, y: 0 })
  }

  const cropMetrics = (scale = cropScale) => {
    const stage = cropStageRef.current
    if (!stage || !cropSize) return null
    const { width, height } = stage.getBoundingClientRect()
    const baseScale = Math.max(width / cropSize.width, height / cropSize.height)
    const imageWidth = cropSize.width * baseScale * scale
    const imageHeight = cropSize.height * baseScale * scale
    return { width, height, imageWidth, imageHeight }
  }

  const clampCropOffset = (next: Point, scale = cropScale) => {
    const metrics = cropMetrics(scale)
    if (!metrics) return next
    const limitX = Math.max(0, (metrics.imageWidth - metrics.width) / 2)
    const limitY = Math.max(0, (metrics.imageHeight - metrics.height) / 2)
    return { x: Math.max(-limitX, Math.min(limitX, next.x)), y: Math.max(-limitY, Math.min(limitY, next.y)) }
  }

  const changeCropScale = (amount: number) => {
    const next = Math.max(0.7, Math.min(2.5, Number((cropScale + amount).toFixed(2))))
    setCropScale(next)
    setCropOffset((current) => clampCropOffset(current, next))
  }

  const startDrag = (event: React.PointerEvent<HTMLDivElement>) => {
    if (!cropSize) return
    event.currentTarget.setPointerCapture(event.pointerId)
    dragStartRef.current = { point: { x: event.clientX, y: event.clientY }, offset: cropOffset }
  }

  const moveDrag = (event: React.PointerEvent<HTMLDivElement>) => {
    const start = dragStartRef.current
    if (!start) return
    setCropOffset(clampCropOffset({ x: start.offset.x + event.clientX - start.point.x, y: start.offset.y + event.clientY - start.point.y }))
  }

  const finishDrag = () => { dragStartRef.current = null }

  const confirmCover = async () => {
    const image = cropImageRef.current
    const metrics = cropMetrics()
    if (!image || !cropFile || !metrics) return
    try {
      setBusy('cover')
      const canvas = document.createElement('canvas')
      canvas.width = 1600
      canvas.height = Math.round(canvas.width / COVER_RATIO)
      const context = canvas.getContext('2d')
      if (!context) throw new Error('无法处理图片，请重试')
      const outputScale = canvas.width / metrics.width
      const drawWidth = metrics.imageWidth * outputScale
      const drawHeight = metrics.imageHeight * outputScale
      const drawX = (canvas.width - drawWidth) / 2 + cropOffset.x * outputScale
      const drawY = (canvas.height - drawHeight) / 2 + cropOffset.y * outputScale
      context.drawImage(image, drawX, drawY, drawWidth, drawHeight)
      const blob = await new Promise<Blob | null>((resolve) => canvas.toBlob(resolve, 'image/jpeg', 0.9))
      if (!blob) throw new Error('图片裁剪失败，请重试')
      const name = `${cropFile.name.replace(/\.[^.]+$/, '') || 'cover'}-cover.jpg`
      setCover(await articleApi.uploadCover(new File([blob], name, { type: 'image/jpeg' })))
      setNotice('封面已上传')
      closeCrop()
    } catch (error) { logRequestError('上传封面失败', error); setNotice('上传封面失败，请稍后重试') } finally { setBusy('') }
  }

  const uploadContentImage = async (file?: File) => {
    if (!file) return null
    try {
      setBusy('content')
      const result = await articleApi.uploadContent(file)
      setContentImages((current) => [...current, result])
      setNotice('图片已插入正文')
      return result
    } catch (error) { logRequestError('上传正文图片失败', error); setNotice('上传图片失败，请稍后重试'); return null } finally { setBusy('') }
  }

  const removeContentImageByURL = (url: string) => {
    setContentImages((current) => current.filter((item) => item.url !== url))
  }

  const resetForm = () => {
    setTitle('')
    setSummary('')
    setContent('')
    setCategoryID(0)
    setCategoryTypeID(0)
    setCover(null)
    setContentImages([])
    setEditingArticle(null)
    window.history.replaceState({}, '', window.location.pathname)
  }

  const buildPayload = (isPublish: boolean) => ({
    title: title.trim(),
    summary: summary.trim(),
    content,
    coverImage: cover?.url,
    ...(categoryID ? { categoryID } : {}),
    isPublish,
  })

  const publish = async () => {
    if (!title.trim() || !content.trim()) { setNotice('请填写标题和正文'); return }
    try {
      setBusy('publish')
      const payload = buildPayload(true)
      const result = editingArticle && !editingArticle.isPublished ? await articleApi.publishDraft({ id: editingArticle.id }) : editingArticle ? await articleApi.edit({ ...payload, id: editingArticle.id }) : await articleApi.create(payload)
      if (!result.success) throw new Error('作品发布失败')
      resetForm()
      setNotice(editingArticle ? '作品已更新并发布' : '作品已发布')
    } catch (error) { logRequestError('发布作品失败', error); setNotice('发布失败，请稍后重试') } finally { setBusy('') }
  }

  const saveDraft = async () => {
    if (!title.trim() || !content.trim()) { setNotice('请填写标题和正文后再保存草稿'); return }
    try {
      setBusy('draft')
      const payload = buildPayload(false)
      const result = editingArticle ? await articleApi.edit({ ...payload, id: editingArticle.id }) : await articleApi.saveDraft(payload)
      if (!result.success) throw new Error('保存草稿失败')
      resetForm()
      setNotice(editingArticle ? '草稿已更新' : '已保存至草稿箱')
    } catch (error) { logRequestError('保存草稿失败', error); setNotice('保存草稿失败，请稍后重试') } finally { setBusy('') }
  }

  const jumpToHeading = (index: number) => {
    editorRef.current?.querySelector<HTMLElement>(`[data-heading-index="${index}"]`)?.scrollIntoView({ behavior: 'smooth', block: 'center' })
  }

  const toggleOutline = (index: number) => {
    setCollapsedOutline((current) => {
      const next = new Set(current)
      if (next.has(index)) next.delete(index)
      else next.add(index)
      return next
    })
  }

  const imageStyle = (() => {
    const metrics = cropMetrics()
    return metrics ? { width: metrics.imageWidth, height: metrics.imageHeight, transform: `translate(calc(-50% + ${cropOffset.x}px), calc(-50% + ${cropOffset.y}px))` } : undefined
  })()
  const previewImageStyle = (() => {
    if (!cropSize || !cropStageRef.current || !cropPreviewRef.current) return undefined
    const stage = cropStageRef.current.getBoundingClientRect()
    const preview = cropPreviewRef.current.getBoundingClientRect()
    const baseScale = Math.max(preview.width / cropSize.width, preview.height / cropSize.height)
    return { width: cropSize.width * baseScale * cropScale, height: cropSize.height * baseScale * cropScale, transform: `translate(calc(-50% + ${cropOffset.x * preview.width / stage.width}px), calc(-50% + ${cropOffset.y * preview.height / stage.height}px))` }
  })()

  return <section className="editor-page">
    {notice && <div className="editor-notice" role="status">{notice}<button type="button" aria-label="关闭提示" onClick={() => setNotice('')}>×</button></div>}
    <div className="editor-layout">
      <main className="editor-canvas">
        <div ref={editorRef}>
          <MarkdownRichEditor
            value={content}
            busy={busy}
            onChange={setContent}
            onUploadImage={uploadContentImage}
            onRemoveImage={removeContentImageByURL}
            onNotice={setNotice}
            onSaveDraft={() => void saveDraft()}
            onPublish={() => void publish()}
            draftLabel={editingArticle ? '更新草稿' : '保存草稿'}
            publishLabel={editingArticle ? '更新发布' : '发布作品'}
            stats={{ words: wordCount, headings: headings.length, images: contentImages.length }}
            header={<label className="rich-title-field"><span>作品标题</span><input value={title} onChange={(event) => setTitle(event.target.value)} placeholder="请输入文章标题（5～100个字）" maxLength={120} /></label>}
          />
        </div>
        <section className="article-details" aria-label="文章发布信息">
          <div className="details-intro"><span>WORK DETAILS</span><h3>补充发布信息</h3><p>完成正文后，再选择分类、补充摘要和封面。</p></div>
          <div className="details-form">
            <div className="details-row media-row"><div className="details-label"><span>添加封面</span><small>16:9，可裁剪</small></div><div className="details-control media-controls"><label className={cover ? 'compact-cover has-cover' : 'compact-cover'}>{cover ? <img src={cover.url} alt="作品封面预览" /> : <div><b>+</b><span>从本地上传</span></div>}<input type="file" accept="image/jpeg,image/png,image/gif,image/webp,image/avif" onChange={(event) => { openCrop(event.target.files?.[0]); event.currentTarget.value = '' }} disabled={busy !== ''} /><em>{cover ? '重新裁剪' : '选择图片'}</em></label></div></div>
            <label className="details-row"><span className="details-label">作品摘要</span><span className="details-control"><textarea className="summary-input" value={summary} onChange={(event) => setSummary(event.target.value)} placeholder="用一两句话概括作品内容（可选）" maxLength={240} /></span></label>
            <div className="details-row"><span className="details-label">作品分类</span><div className="details-control category-selects" ref={categorySelectsRef}><div className="category-select"><span>一级分类</span><button type="button" aria-haspopup="listbox" aria-expanded={openCategoryMenu === 'type'} onClick={() => setOpenCategoryMenu((current) => current === 'type' ? null : 'type')}>{categoryGroups.find((group) => group.type.id === categoryTypeID)?.type.name ?? '其它'}<i aria-hidden="true" /></button>{openCategoryMenu === 'type' && <div className="category-select-menu" role="listbox"><button type="button" role="option" aria-selected={!categoryTypeID} className={!categoryTypeID ? 'selected' : ''} onClick={() => { setCategoryTypeID(0); setCategoryID(0); setOpenCategoryMenu(null) }}>其它</button>{categoryGroups.length ? categoryGroups.map((group) => <button type="button" role="option" aria-selected={group.type.id === categoryTypeID} className={group.type.id === categoryTypeID ? 'selected' : ''} key={group.type.id} onClick={() => { setCategoryTypeID(group.type.id); setCategoryID(0); setOpenCategoryMenu(null) }}>{group.type.name}</button>) : <p>暂无一级分类</p>}</div>}</div><div className="category-select"><span>二级分类</span><button type="button" disabled={!categoryTypeID} aria-haspopup="listbox" aria-expanded={openCategoryMenu === 'category'} onClick={() => setOpenCategoryMenu((current) => current === 'category' ? null : 'category')}>{categoryGroups.find((group) => group.type.id === categoryTypeID)?.categories.find((item) => item.id === categoryID)?.name ?? (categoryTypeID ? '选择二级分类' : '其它')}<i aria-hidden="true" /></button>{openCategoryMenu === 'category' && <div className="category-select-menu" role="listbox">{categoryGroups.find((group) => group.type.id === categoryTypeID)?.categories.length ? categoryGroups.find((group) => group.type.id === categoryTypeID)?.categories.map((item) => <button type="button" role="option" aria-selected={item.id === categoryID} className={item.id === categoryID ? 'selected' : ''} key={item.id} onClick={() => { setCategoryID(item.id); setOpenCategoryMenu(null) }}>{item.name}</button>) : <p>暂无二级分类</p>}</div>}</div></div></div>
          </div>
        </section>
      </main>
      <aside className="editor-aside"><nav className="toc-panel" aria-label="文章目录"><div className="aside-heading"><span>OUTLINE</span><h3>文章目录</h3></div>{headings.length ? <div className="toc-list">{visibleOutlineItems.map((heading) => <div className={`toc-item toc-level-${heading.level}${collapsedOutline.has(heading.index) ? ' collapsed' : ''}`} key={heading.index}>{heading.hasChildren ? <button type="button" className="toc-toggle" aria-label={collapsedOutline.has(heading.index) ? '展开目录' : '折叠目录'} aria-expanded={!collapsedOutline.has(heading.index)} onClick={() => toggleOutline(heading.index)} /> : <span className="toc-toggle-spacer" />}<button type="button" className="toc-link" onClick={() => jumpToHeading(heading.index)}>{heading.text}</button></div>)}</div> : <p>使用 `#`、`##` 或 `###` 标题后，目录会自动生成。</p>}</nav></aside>
    </div>
    {cropUrl && <div className="crop-modal-backdrop" role="presentation" onMouseDown={(event) => { if (event.target === event.currentTarget && busy !== 'cover') closeCrop() }}><section className="crop-modal" role="dialog" aria-modal="true" aria-labelledby="crop-title"><header><div><span>IMAGE EDITOR</span><h3 id="crop-title">图片编辑</h3></div><button type="button" aria-label="关闭图片编辑" disabled={busy === 'cover'} onClick={closeCrop}>×</button></header><div className="crop-workspace"><div className="crop-stage" ref={cropStageRef} onPointerDown={startDrag} onPointerMove={moveDrag} onPointerUp={finishDrag} onPointerCancel={finishDrag}><img ref={cropImageRef} src={cropUrl} alt="待裁剪封面" draggable={false} style={imageStyle} onLoad={(event) => { const image = event.currentTarget; setCropSize({ width: image.naturalWidth, height: image.naturalHeight }) }} /></div><div className="crop-side"><span>封面预览</span><div className="crop-preview" ref={cropPreviewRef}><img src={cropUrl} alt="封面裁剪预览" style={previewImageStyle} /></div><div className="zoom-controls"><span>缩放</span><div><button type="button" aria-label="缩小图片" onClick={() => changeCropScale(-0.1)} disabled={cropScale <= 0.7}>−</button><output>{Math.round(cropScale * 100)}%</output><button type="button" aria-label="放大图片" onClick={() => changeCropScale(0.1)} disabled={cropScale >= 2.5}>+</button></div></div></div></div><footer><p>拖动图片调整位置，裁剪区域固定为 16:9。</p><div><button type="button" className="crop-cancel" disabled={busy === 'cover'} onClick={closeCrop}>取消</button><button type="button" className="crop-confirm" disabled={!cropSize || busy === 'cover'} onClick={() => void confirmCover()}>{busy === 'cover' ? '上传中...' : '确认上传'}</button></div></footer></section></div>}
  </section>
}
