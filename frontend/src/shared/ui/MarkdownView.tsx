import { useEffect, useMemo, useRef } from 'react'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { resolveMarkdownImages } from '@/shared/lib/storage'

interface Props {
  content: string
  className?: string
}

marked.setOptions({ breaks: true, gfm: true })

const keywordPattern = /\b(?:abstract|and|as|async|await|break|case|catch|class|const|continue|def|default|delete|do|else|enum|export|extends|false|finally|for|from|func|function|if|implements|import|in|interface|let|map|new|nil|none|not|null|of|or|package|private|protected|public|range|return|select|static|string|struct|switch|this|throw|true|try|type|undefined|var|while|with|yield)\b/gi
const tokenPattern = /(\/\/[^\n]*|#[^\n]*|\/\*[\s\S]*?\*\/|<!--[\s\S]*?-->|`[^`]*`|'(?:\\.|[^'\\])*'|"(?:\\.|[^"\\])*"|\b\d+(?:\.\d+)?\b|\b(?:abstract|and|as|async|await|break|case|catch|class|const|continue|def|default|delete|do|else|enum|export|extends|false|finally|for|from|func|function|if|implements|import|in|interface|let|map|new|nil|none|not|null|of|or|package|private|protected|public|range|return|select|static|string|struct|switch|this|throw|true|try|type|undefined|var|while|with|yield)\b|\b[A-Za-z_][\w]*(?=\())/gi

function appendHighlightedCode(target: HTMLElement, value: string) {
  let cursor = 0
  for (const match of value.matchAll(tokenPattern)) {
    const index = match.index ?? 0
    if (index > cursor) target.append(document.createTextNode(value.slice(cursor, index)))
    const token = match[0]
    const className = token.startsWith('//') || token.startsWith('#') || token.startsWith('/*') || token.startsWith('<!--')
      ? 'token-comment'
      : token.startsWith('"') || token.startsWith("'") || token.startsWith('`')
        ? 'token-string'
        : /^\d/.test(token)
          ? 'token-number'
          : keywordPattern.test(token)
            ? 'token-keyword'
            : 'token-function'
    keywordPattern.lastIndex = 0
    const span = document.createElement('span')
    span.className = className
    span.textContent = token
    target.append(span)
    cursor = index + token.length
  }
  if (cursor < value.length) target.append(document.createTextNode(value.slice(cursor)))
}

function renderCodeBlocks(html: string) {
  const template = document.createElement('template')
  template.innerHTML = html

  template.content.querySelectorAll('pre').forEach((pre) => {
    const code = pre.querySelector('code')
    if (!code) return

    pre.dataset.language = (code.className.match(/language-([\w-]+)/)?.[1] ?? 'code').toUpperCase()
    const source = code.textContent?.replace(/\n$/, '') ?? ''
    code.textContent = ''
    code.classList.add('code-with-lines')

    source.split('\n').forEach((line, index) => {
      const row = document.createElement('span')
      row.className = 'code-line'
      const lineNumber = document.createElement('span')
      lineNumber.className = 'code-line-number'
      lineNumber.textContent = String(index + 1)
      const lineContent = document.createElement('span')
      lineContent.className = 'code-line-content'
      appendHighlightedCode(lineContent, line)
      row.append(lineNumber, lineContent)
      code.append(row)
    })

    const button = document.createElement('button')
    button.type = 'button'
    button.className = 'code-copy'
    button.textContent = 'Copy'
    button.setAttribute('aria-label', 'Copy code')
    pre.append(button)
  })

  return template.innerHTML
}

export default function MarkdownView({ content, className }: Props) {
  const rootRef = useRef<HTMLDivElement>(null)
  const html = useMemo(() => {
    if (!content) return ''
    const raw = marked.parse(resolveMarkdownImages(content)) as string
    return renderCodeBlocks(DOMPurify.sanitize(raw))
  }, [content])

  useEffect(() => {
    const root = rootRef.current
    if (!root) return

    const handleCopy = async (event: MouseEvent) => {
      const target = event.target as Element | null
      const button = target?.closest<HTMLButtonElement>('.code-copy')
      if (!button || !root.contains(button)) return

      const code = button.parentElement?.querySelector('code')
      const codeText = Array.from(code?.querySelectorAll('.code-line-content') ?? [])
        .map((line) => line.textContent ?? '')
        .join('\n')
      try {
        await navigator.clipboard.writeText(codeText)
        button.textContent = 'Copied'
      } catch {
        button.textContent = 'Failed'
      }
      window.setTimeout(() => { button.textContent = 'Copy' }, 1400)
    }

    root.addEventListener('click', handleCopy)
    return () => root.removeEventListener('click', handleCopy)
  }, [])

  return <div ref={rootRef} className={className ?? 'markdown-body'} dangerouslySetInnerHTML={{ __html: html }} />
}
