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
  console.log(`Making request to: ${API_BASE_URL}${endpoint}`)
  
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
    credentials: 'include', // Включаем cookies в запросы
    ...options,
  })

  console.log(`Response status: ${response.status}`)
  
  if (!response.ok) {
    const errorData = await response.text()
    console.error(`API Error: ${response.status}`, errorData)
    throw new ApiError(errorData, response.status)
  }

  return response.json()
}

export const getUserBookings = async (): Promise<Booking[]> => {
  console.log("Fetching user bookings")
  return apiRequest<Booking[]>("/bookings")
}

export const getBookingById = async (id: number): Promise<Booking> => {
  console.log(`Fetching booking ${id}`)
  return apiRequest<Booking>(`/bookings/${id}`)
}

export const bookEvent = async (
  bookingData: Omit<
    Booking,
    "id" | "userId" | "status" | "createdAt" | "expiresAt"
  >
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
