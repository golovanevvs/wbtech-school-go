import { Event } from "../lib/types"

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL

class ApiError extends Error {
  constructor(message: string, public status: number) {
    super(message)
    this.name = "ApiError"
  }
}

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

const getToken = (): string | null => {
  if (typeof window !== "undefined") {
    return localStorage.getItem("token")
  }
  return null
}

export const getEvents = async (): Promise<Event[]> => {
  try {
    const response = await apiRequest<Event[]>("/events")
    return response || []
  } catch (error) {
    console.error("Failed to load events:", error)
    return []
  }
}

export const getEventById = async (id: number): Promise<Event> => {
  return apiRequest<Event>(`/events/${id}`)
}

export const createEvent = async (
  eventData: Omit<Event, "id" | "createdAt" | "updatedAt" | "availablePlaces">
): Promise<Event> => {
  return apiRequest<Event>("/events", {
    method: "POST",
    body: JSON.stringify(eventData),
  })
}

export const updateEvent = async (
  id: number,
  eventData: Partial<Event>
): Promise<Event> => {
  return apiRequest<Event>(`/events/${id}`, {
    method: "PUT",
    body: JSON.stringify(eventData),
  })
}

export const deleteEvent = async (id: number): Promise<void> => {
  await apiRequest(`/events/${id}`, {
    method: "DELETE",
  })
}
