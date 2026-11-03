import type { AnchorHTMLAttributes, MouseEvent } from 'react'
import { navigate } from '@/shared/hooks/usePathname'

interface AppLinkProps extends AnchorHTMLAttributes<HTMLAnchorElement> {
  to: string
}

export default function AppLink({ to, onClick, ...props }: AppLinkProps) {
  const handleClick = (event: MouseEvent<HTMLAnchorElement>) => {
    onClick?.(event)
    if (event.defaultPrevented || event.button !== 0 || event.metaKey || event.ctrlKey || event.shiftKey) return
    if (/^\/blog\/\d+$/.test(to) || /^\/users\/\d+$/.test(to)) {
      event.preventDefault()
      window.open(to, '_blank', 'noopener')
      return
    }
    event.preventDefault()
    navigate(to)
  }

  return <a href={to} onClick={handleClick} {...props} />
}
