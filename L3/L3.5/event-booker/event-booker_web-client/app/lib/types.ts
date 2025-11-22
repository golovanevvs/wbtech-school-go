export interface User {
  id: number
  email: string
  name: string
  telegramUsername?: string | null
  telegramChatID?: number | null
}

export interface Event {
  id: number
  title: string
  date: string
  description: string
  totalPlaces: number
  availablePlaces: number
  bookingDeadline: number
  createdAt: string
  updatedAt: string
}

export interface Booking {
  id: number
  userId: number
  eventId: number
  status: "pending" | "confirmed" | "cancelled"
  createdAt: string
  expiresAt: string
  confirmedAt?: string
  cancelledAt?: string
}

export interface AuthResponse {
  user: User
  token: string
  refreshToken?: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
  name: string
}

export interface LoginResponse {
  token: string
  refreshToken: string
  user?: User
}

export interface RegisterResponse {
  token: string
  refreshToken: string
  user?: User
}

export interface UpdateUserRequest {
  name?: string
  telegramUsername?: string | null
  telegramNotifications?: boolean
  emailNotifications?: boolean
}
