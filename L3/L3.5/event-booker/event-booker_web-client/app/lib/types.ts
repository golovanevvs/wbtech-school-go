export interface User {
  id: number
  email: string
  name: string
  telegramUsername?: string | null
  telegramChatID?: number | null
  telegramNotifications?: boolean
  emailNotifications?: boolean
  created_at: string
  updated_at: string
}

export interface Event {
  id: number
  title: string
  date: string
  description: string
  totalPlaces: number
  availablePlaces: number
  bookingDeadline: number
  ownerId: number
  telegramNotifications?: boolean
  emailNotifications?: boolean
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

// Интерфейс для бронирования с сервера (snake_case)
export interface BookingFromServer {
  id: number
  user_id: number
  event_id: number
  status: "pending" | "confirmed" | "cancelled"
  created_at: string
  expires_at: string
  confirmed_at?: string
  cancelled_at?: string
}

// Функция трансформации snake_case → camelCase
export function transformBookingFromServer(booking: BookingFromServer): Booking {
  return {
    id: booking.id,
    userId: booking.user_id,
    eventId: booking.event_id,
    status: booking.status,
    createdAt: booking.created_at,
    expiresAt: booking.expires_at,
    confirmedAt: booking.confirmed_at,
    cancelledAt: booking.cancelled_at,
  }
}

// Новый интерфейс для запроса бронирования (соответствует серверу)
export interface CreateBookingRequest {
  event_id: number
  booking_deadline_minutes?: number
}

export interface AuthResponse {
  user: User
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
  message: string
}

export interface RegisterResponse {
  message: string
  user: User
}

export interface UpdateUserRequest {
  name?: string
  telegramUsername?: string | null
  telegramNotifications?: boolean
  emailNotifications?: boolean
  resetTelegramChatID?: boolean // специальный флаг для сброса chatID на сервере
}