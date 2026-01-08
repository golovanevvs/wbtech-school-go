import { Booking, CreateBookingRequest, BookingFromServer } from "../lib/types"

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
  console.log(`Making request to: ${API_BASE_URL}${endpoint}`)

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
    credentials: "include",
    ...options,
  })

  console.log(`Response status: ${response.status}`)

  if (!response.ok) {
    const errorText = await response.text()
    console.error(`API Error: ${response.status}`, errorText)
    // Попытка извлечь сообщение об ошибке из JSON
    try {
      const errorJson = JSON.parse(errorText)
      const errorMessage = errorJson.error || errorJson.message || errorText
      throw new ApiError(errorMessage, response.status)
    } catch {
      // Если не JSON, используем текст как есть
      throw new ApiError(errorText, response.status)
    }
  }

  return response.json()
}

export const getUserBookings = async (): Promise<BookingFromServer[]> => {
  console.log("Fetching user bookings")
  return apiRequest<BookingFromServer[]>("/bookings")
}

export const getUserBookingByEventId = async (
  eventId: number
): Promise<Booking | null> => {
  console.log(`Fetching user booking for event ${eventId}`)
  try {
    const booking = await apiRequest<Booking>(`/bookings/user/event/${eventId}`)
    return booking
  } catch (error) {
    if (error instanceof ApiError && error.status === 404) {
      return null
    }
    throw error
  }
}

export const getBookingById = async (id: number): Promise<Booking> => {
  console.log(`Fetching booking ${id}`)
  return apiRequest<Booking>(`/bookings/${id}`)
}

export const bookEvent = async (
  bookingData: CreateBookingRequest
): Promise<Booking> => {
  console.log("Booking event with data:", bookingData)
  return apiRequest<Booking>("/bookings", {
    method: "POST",
    body: JSON.stringify(bookingData),
  })
}

export const confirmBooking = async (id: number): Promise<Booking> => {
  console.log(`Confirming booking ${id}`)
  return apiRequest<Booking>(`/bookings/${id}/confirm`, {
    method: "POST",
  })
}

export const cancelBooking = async (id: number): Promise<Booking> => {
  console.log(`Cancelling booking ${id}`)
  return apiRequest<Booking>(`/bookings/${id}/cancel`, {
    method: "POST",
  })
}

export const deleteBooking = async (id: number): Promise<void> => {
  console.log(`Deleting booking ${id}`)
  await apiRequest(`/bookings/${id}`, {
    method: "DELETE",
  })
}
