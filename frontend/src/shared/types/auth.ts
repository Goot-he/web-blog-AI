export type EmailCodePurpose = 'register' | 'login' | 'forgot_password'

export interface SendEmailCodeRequest {
  email: string
  purpose: EmailCodePurpose
}

export interface LoginRequest {
  username: string
  password: string
}

export interface EmailLoginRequest {
  email: string
  code: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
  code: string
}

export interface ForgotPasswordRequest {
  email: string
  code: string
  new_password: string
  confirm_password: string
}

export interface AuthUser {
  id: number
  username: string
  email?: string
  phone?: string
  avatar: string
  sex?: string
  birthday?: string
  role: string
  status: string
}

export interface AuthResponse {
  access_token: string
  access_expires: number
  id: number
  username: string
  avatar: string
  role: string
  status: string
}
