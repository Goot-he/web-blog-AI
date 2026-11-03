import { clearAccessToken, request, setAccessToken } from '@/shared/api/http'
import type {
  AuthResponse,
  EmailLoginRequest,
  ForgotPasswordRequest,
  LoginRequest,
  RegisterRequest,
  SendEmailCodeRequest,
} from '@/shared/types/auth'

export const authApi = {
  sendEmailCode(payload: SendEmailCodeRequest) {
    return request<null>('/v1/auth/send-email-code', { method: 'POST', body: payload, auth: false })
  },

  register(payload: RegisterRequest) {
    return request<null>('/v1/auth/register', { method: 'POST', body: payload, auth: false })
  },

  login(payload: LoginRequest) {
    return request<AuthResponse>('/v1/auth/login', { method: 'POST', body: payload, auth: false }).then((data) => {
      setAccessToken(data.access_token)
      return data
    })
  },

  emailLogin(payload: EmailLoginRequest) {
    return request<AuthResponse>('/v1/auth/email_login', { method: 'POST', body: payload, auth: false }).then((data) => {
      setAccessToken(data.access_token)
      return data
    })
  },

  forgotPassword(payload: ForgotPasswordRequest) {
    return request<null>('/v1/auth/forgot-password', { method: 'POST', body: payload, auth: false })
  },

  logout(token?: string) {
    return request<null>('/v1/auth/logout', { method: 'GET', ...(token ? { token } : {}) }).finally(() => {
      clearAccessToken()
    })
  },
}
