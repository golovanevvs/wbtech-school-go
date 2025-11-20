import { Event } from "../lib/types"

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"

// Интерфейс для обработки ошибок API
class ApiError extends Error {
  constructor(message: string, public status: number) {
    super(message)
    this.name = "ApiError"
  }
}

// Функция для выполнения API запросов
const apiRequest = async <T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> => {
  const token = getToken()

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...options.headers,
    },
    ...options,
  })

  if (!response.ok) {
    const errorData = await response.text()
    throw new ApiError(errorData, response.status)
  }

  return response.json()
}

// Функция для получения токена из localStorage
const getToken = (): string | null => {
  if (typeof window !== "undefined") {
    return localStorage.getItem("token")
  }
  return null
}

// Функция для получения всех событий
export const getEvents = async (): Promise<Event[]> => {
  return apiRequest<Event[]>("/events")
}

// Функция для получения события по ID
export const getEventById = async (id: number): Promise<Event> => {
  return apiRequest<Event>(`/events/${id}`)
}

// Функция для создания события
export const createEvent = async (
  eventData: Omit<Event, "id" | "createdAt" | "updatedAt" | "availablePlaces">
): Promise<Event> => {
  return apiRequest<Event>("/events", {
    method: "POST",
    body: JSON.stringify(eventData),
  })
}

// Функция для обновления события
export const updateEvent = async (
  id: number,
  eventData: Partial<Event>
): Promise<Event> => {
  return apiRequest<Event>(`/events/${id}`, {
    method: "PUT",
    body: JSON.stringify(eventData),
  })
}

// Функция для удаления события
export const deleteEvent = async (id: number): Promise<void> => {
  await apiRequest(`/events/${id}`, {
    method: "DELETE",
  })
}
