export type UserRole = "Кладовщик" | "Менеджер" | "Аудитор"

export interface User {
  id: number
  username: string
  name: string
  user_role: UserRole
}

export interface AuthTokens {
  access_token: string
  refresh_token: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  password: string
  name: string
  role: UserRole
}

export interface LoginResponse {
  message: string
}

export interface RegisterResponse {
  message: string
  user: User
}