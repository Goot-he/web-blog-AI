import { useState, type FormEvent } from 'react'
import AppLink from '@/shared/ui/AppLink'
import { authApi } from '@/shared/api/v1/auth'
import { userApi } from '@/shared/api/v1/user'
import type { EmailCodePurpose } from '@/shared/types/auth'
import { navigate } from '@/shared/hooks/usePathname'
import { logRequestError } from '@/shared/lib/request-error'
import { saveCurrentUser } from '@/shared/lib/current-user'
import './auth-form.css'

type Mode = 'password' | 'email'

export default function LoginPage() {
  const [mode, setMode] = useState<Mode>('password')
  const [forgot, setForgot] = useState(false)
  const [message, setMessage] = useState('')
  const [sending, setSending] = useState(false)

  const sendCode = async (email: string, purpose: EmailCodePurpose) => {
    if (!email.trim()) { setMessage('请输入邮箱'); return }
    setSending(true)
    try { await authApi.sendEmailCode({ email, purpose }); setMessage('验证码已发送，请查收邮箱') }
    catch (error) { logRequestError('发送登录验证码失败', error); setMessage('验证码发送失败，请稍后重试') }
    finally { setSending(false) }
  }

  const submit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    const data = new FormData(event.currentTarget)
    const email = String(data.get('email') ?? '').trim()
    const code = String(data.get('code') ?? '').trim()
    const username = String(data.get('username') ?? '').trim()
    const password = String(data.get('password') ?? '')
    if (forgot && (!email || !code || !password)) { setMessage('请填写完整信息'); return }
    if (!forgot && mode === 'password' && (!username || !password)) { setMessage('请输入用户名和密码'); return }
    if (!forgot && mode === 'email' && (!email || !code)) { setMessage('请输入邮箱和验证码'); return }
    try {
      if (forgot) {
        const password = String(data.get('password'))
        await authApi.forgotPassword({ email: String(data.get('email')), code: String(data.get('code')), new_password: password, confirm_password: password })
        setForgot(false); setMessage('密码已重置，请登录'); return
      }
      const result = mode === 'password'
        ? await authApi.login({ username: String(data.get('username')), password: String(data.get('password')) })
        : await authApi.emailLogin({ email: String(data.get('email')), code: String(data.get('code')) })
      saveCurrentUser({ id: result.id, username: result.username, avatar: result.avatar, role: result.role, status: result.status })
      try {
        saveCurrentUser(await userApi.current(result.access_token))
      } catch (profileError) {
        logRequestError('登录后加载用户信息失败', profileError)
      }
      const redirectTo = typeof window.history.state?.redirectTo === 'string' && window.history.state.redirectTo.startsWith('/')
        ? window.history.state.redirectTo
        : '/'
      setMessage('登录成功'); window.setTimeout(() => navigate(redirectTo), 300)
    } catch (error) {
      logRequestError('登录失败', error)
      const detail = error instanceof Error ? error.message : ''
      setMessage(detail === 'USER_BLOCKED' ? '该账号已被管理员拉黑，无法登录。' : '登录失败，请稍后重试')
    }
  }

  return <main className="auth-page"><section className="auth-card login-card">
    <h1>{forgot ? '重置密码' : '登录 Renai'}</h1>
    {!forgot && <div className="login-tabs"><button className={mode === 'password' ? 'active' : ''} onClick={() => setMode('password')} type="button">用户名登录</button><button className={mode === 'email' ? 'active' : ''} onClick={() => setMode('email')} type="button">邮箱登录</button></div>}
    <form className="auth-form" noValidate onSubmit={submit}>
      {mode === 'password' && !forgot && <input name="username" required placeholder="请输入用户名" />}
      {(mode === 'email' || forgot) && <input name="email" required type="email" placeholder="请输入邮箱" />}
      {mode === 'email' && !forgot && <CodeField sending={sending} purpose="login" onSend={sendCode} />}
      {(mode === 'password' || forgot) && <input name="password" required minLength={6} type="password" placeholder={forgot ? '请输入新密码' : '请输入密码'} />}
      {forgot && <CodeField sending={sending} purpose="forgot_password" onSend={sendCode} />}
      <button className="login-submit" type="submit">{forgot ? '重置密码' : '登录'}</button>
      {message && <p className="form-message" role="status">{message}</p>}
    </form>
    <div className="login-links"><button type="button" onClick={() => { setForgot(!forgot); setMessage('') }}>{forgot ? '返回登录' : '忘记密码？'}</button>{!forgot && <AppLink to="/register">注册账号</AppLink>}</div>
    {!forgot && <><div className="social-divider">其他方式登录</div><div className="social-login"><button type="button" onClick={() => setMessage('QQ 登录暂未开发')}>Q</button><button type="button" onClick={() => setMessage('微信登录暂未开发')}>微</button></div></>}
  </section></main>
}

function CodeField({ sending, purpose, onSend }: { sending: boolean; purpose: EmailCodePurpose; onSend: (email: string, purpose: EmailCodePurpose) => void }) {
  return <div className="code-input"><input name="code" required maxLength={6} placeholder="请输入验证码" /><button type="button" disabled={sending} onClick={(event) => onSend(String(new FormData(event.currentTarget.form ?? undefined).get('email') ?? ''), purpose)}>{sending ? '发送中' : '获取验证码'}</button></div>
}
