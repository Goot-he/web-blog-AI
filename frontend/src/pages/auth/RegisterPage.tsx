import { useState, type FormEvent } from 'react'
import AppLink from '@/shared/ui/AppLink'
import { authApi } from '@/shared/api/v1/auth'
import { navigate } from '@/shared/hooks/usePathname'
import { logRequestError } from '@/shared/lib/request-error'
import './auth-form.css'

export default function RegisterPage() {
  const [message, setMessage] = useState('')
  const [sending, setSending] = useState(false)
  const sendCode = async (email: string) => {
    if (!email.trim()) { setMessage('请输入邮箱'); return }
    setSending(true)
    try { await authApi.sendEmailCode({ email, purpose: 'register' }); setMessage('验证码已发送，请查收邮箱') }
    catch (error) { logRequestError('发送注册验证码失败', error); setMessage('验证码发送失败，请稍后重试') }
    finally { setSending(false) }
  }
  const submit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    const data = new FormData(event.currentTarget)
    if (![data.get('username'), data.get('email'), data.get('code'), data.get('password')].every((value) => String(value ?? '').trim())) { setMessage('请填写完整信息'); return }
    try {
      await authApi.register({ username: String(data.get('username')), email: String(data.get('email')), password: String(data.get('password')), code: String(data.get('code')) })
      setMessage('注册成功，正在跳转登录'); window.setTimeout(() => navigate('/login'), 500)
    } catch (error) { logRequestError('注册失败', error); setMessage('注册失败，请稍后重试') }
  }
  return <main className="auth-page"><section className="auth-card login-card"><h1>创建账号</h1><form className="auth-form" noValidate onSubmit={submit}><input name="username" required minLength={2} placeholder="请输入用户名" /><input name="email" required type="email" placeholder="请输入邮箱" /><div className="code-input"><input name="code" required maxLength={6} placeholder="请输入验证码" /><button type="button" disabled={sending} onClick={(event) => sendCode(String(new FormData(event.currentTarget.form ?? undefined).get('email') ?? ''))}>{sending ? '发送中' : '获取验证码'}</button></div><input name="password" required minLength={6} type="password" placeholder="请输入密码" /><button className="login-submit" type="submit">注册</button>{message && <p className="form-message" role="status">{message}</p>}</form><p className="auth-switch">已有账号？ <AppLink to="/login">返回登录</AppLink></p></section></main>
}
