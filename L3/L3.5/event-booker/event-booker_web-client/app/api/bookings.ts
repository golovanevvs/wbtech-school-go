import { Booking } from "../lib/types"

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

export const getUserBookings = async (): Promise<Booking[]> => {
  return apiRequest<Booking[]>("/bookings")
}

export const getBookingById = async (id: number): Promise<Booking> => {
  return apiRequest<Booking>(`/bookings/${id}`)
}

export const bookEvent = async (
  bookingData: Omit<
    Booking,
    "id" | "userId" | "status" | "createdAt" | "expiresAt"
  >
): Promise<Booking> => {
  return apiRequest<Booking>("/bookings", {
    method: "POST",
    body: JSON.stringify(bookingData),
  })
}

export const confirmBooking = async (id: number): Promise<Booking> => {
  return apiRequest<Booking>(`/bookings/${id}/confirm`, {
    method: "POST",
  })
}

export const cancelBooking = async (id: number): Promise<Booking> => {
  return apiRequest<Booking>(`/bookings/${id}/cancel`, {
    method: "POST",
  })
}

export const deleteBooking = async (id: number): Promise<void> => {
  await apiRequest(`/bookings/${id}`, {
    method: "DELETE",
  })
}
