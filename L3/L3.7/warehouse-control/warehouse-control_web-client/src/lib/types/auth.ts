export type UserRole = "Кладовщик" | "Менеджер" | "Аудитор"

export interface User {
  id: number
  username: string
  role: UserRole
}

export interface AuthTokens {
  access_token: string
  refresh_token: string
}
