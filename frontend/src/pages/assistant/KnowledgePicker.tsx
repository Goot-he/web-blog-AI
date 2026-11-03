import { useEffect, useMemo, useRef, useState } from 'react'
import type { VectorCollection } from '@/shared/types/vector'

interface KnowledgePickerProps {
  collections: VectorCollection[]
  value: string
  disabled?: boolean
  onChange: (value: string) => void
}

export default function KnowledgePicker({ collections, value, disabled, onChange }: KnowledgePickerProps) {
  const [open, setOpen] = useState(false)
  const [query, setQuery] = useState('')
  const rootRef = useRef<HTMLDivElement>(null)
  const selected = collections.find((collection) => collection.name === value)
  const filtered = useMemo(() => collections.filter((collection) => collection.name.toLowerCase().includes(query.trim().toLowerCase())), [collections, query])

  useEffect(() => {
    const close = (event: PointerEvent) => {
      if (!rootRef.current?.contains(event.target as Node)) setOpen(false)
    }
    document.addEventListener('pointerdown', close)
    return () => document.removeEventListener('pointerdown', close)
  }, [])

  return <div className={`assistant-knowledge-picker${open ? ' open' : ''}`} ref={rootRef}>
    <button type="button" className="assistant-knowledge-trigger" aria-haspopup="listbox" aria-expanded={open} disabled={disabled} onClick={() => { setOpen((current) => !current); setQuery('') }}><span><small>知识库</small><b>{selected?.name || '请选择知识库'}</b></span><i aria-hidden="true" /></button>
    {open && <div className="assistant-knowledge-popover"><label><span className="assistant-search-icon" aria-hidden="true" /><input autoFocus value={query} onChange={(event) => setQuery(event.target.value)} placeholder="搜索知识库" /></label><div className="assistant-knowledge-options" role="listbox" aria-label="知识库列表">{filtered.length === 0 ? <p>没有匹配的知识库</p> : filtered.map((collection) => <button type="button" role="option" aria-selected={collection.name === value} className={collection.name === value ? 'selected' : ''} key={collection.name} onClick={() => { onChange(collection.name); setOpen(false); setQuery('') }}><span><b>{collection.name}</b><small>{collection.count} 个切块</small></span>{collection.name === value && <i aria-hidden="true">✓</i>}</button>)}</div></div>}
  </div>
}
