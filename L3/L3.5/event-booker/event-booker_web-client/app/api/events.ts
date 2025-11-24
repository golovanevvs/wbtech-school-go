import { Event } from "../lib/types"

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL

interface EventFromServer {
  id: number
  title: string
  date: string
  description: string
  total_places: number
  available_places: number
  booking_deadline: number
  owner_id: number
  telegram_notifications?: boolean
  email_notifications?: boolean
  created_at: string
  updated_at: string
}

interface EventUpdateData {
  title?: string
  date?: string
  description?: string
  total_places?: number
  booking_deadline?: number
  telegram_notifications?: boolean
  email_notifications?: boolean
}

const transformEventFromServer = (serverEvent: EventFromServer): Event => {
  console.log("Transforming event from server:", serverEvent)

  const transformed: Event = {
    id: serverEvent.id,
    title: serverEvent.title,
    date: serverEvent.date,
    description: serverEvent.description,
    totalPlaces: serverEvent.total_places,
    availablePlaces: serverEvent.available_places,
    bookingDeadline: serverEvent.booking_deadline,
    ownerId: serverEvent.owner_id,
    telegramNotifications: serverEvent.telegram_notifications,
    emailNotifications: serverEvent.email_notifications,
    createdAt: serverEvent.created_at,
    updatedAt: serverEvent.updated_at,
  }

  console.log("Transformed event:", transformed)
  return transformed
}

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
    const errorData = await response.text()
    console.error(`API Error: ${response.status}`, errorData)
    throw new ApiError(errorData, response.status)
  }

  const contentType = response.headers.get("content-type")
  if (contentType && contentType.includes("application/json")) {
    const result = await response.json()
    console.log(`API Response data:`, result)
    return result
  } else {
    console.log("Response has no JSON content, returning empty object")
    return {} as T
  }
}

const apiRequestWithoutJson = async (
  endpoint: string,
  options: RequestInit = {}
): Promise<void> => {
  console.log(
    `Making request without JSON parsing to: ${API_BASE_URL}${endpoint}`
  )

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
    const errorData = await response.text()
    console.error(`API Error: ${response.status}`, errorData)
    throw new ApiError(errorData, response.status)
  }

  console.log("Request completed successfully")
}

export const getEvents = async (): Promise<Event[]> => {
  try {
    console.log("Fetching events from API...")
    const response = await apiRequest<EventFromServer[]>("/events")
    console.log("Events API response:", response)
    console.log("Events array length:", response?.length)

    if (response && response.length > 0) {
      console.log("First event from server (raw):", response[0])
    }

    const transformedEvents = response.map(transformEventFromServer)
    console.log("Transformed events:", transformedEvents)

    return transformedEvents
  } catch (error) {
    console.error("Failed to load events:", error)
    return []
  }
}

export const getEventById = async (id: number): Promise<Event> => {
  const serverEvent = await apiRequest<EventFromServer>(`/events/${id}`)
  return transformEventFromServer(serverEvent)
}

export const createEvent = async (
  eventData: Omit<Event, "id" | "createdAt" | "updatedAt" | "availablePlaces">
): Promise<Event> => {
  console.log("Creating event with data:", eventData)

  const serverEventData: Partial<EventFromServer> = {
    title: eventData.title,
    date: eventData.date,
    description: eventData.description,
    total_places: eventData.totalPlaces,
    booking_deadline: eventData.bookingDeadline,
    owner_id: eventData.ownerId,
    telegram_notifications: eventData.telegramNotifications,
    email_notifications: eventData.emailNotifications,
  }

  console.log("Sending to server (snake_case):", serverEventData)

  const createdEvent = await apiRequest<EventFromServer>("/events", {
    method: "POST",
    body: JSON.stringify(serverEventData),
  })

  return transformEventFromServer(createdEvent)
}

export const updateEvent = async (
  id: number,
  eventData: Partial<Event>
): Promise<Event> => {
  console.log(`Updating event ${id} with data:`, eventData)

  const serverEventData: EventUpdateData = {}
  if (eventData.title !== undefined) serverEventData.title = eventData.title
  if (eventData.date !== undefined) serverEventData.date = eventData.date
  if (eventData.description !== undefined)
    serverEventData.description = eventData.description
  if (eventData.totalPlaces !== undefined)
    serverEventData.total_places = eventData.totalPlaces
  if (eventData.bookingDeadline !== undefined)
    serverEventData.booking_deadline = eventData.bookingDeadline
  if (eventData.telegramNotifications !== undefined)
    serverEventData.telegram_notifications = eventData.telegramNotifications
  if (eventData.emailNotifications !== undefined)
    serverEventData.email_notifications = eventData.emailNotifications

  console.log("Sending to server (snake_case):", serverEventData)
  console.log("Server event data keys:", Object.keys(serverEventData))

  const updatedEvent = await apiRequest<EventFromServer>(`/events/${id}`, {
    method: "PUT",
    body: JSON.stringify(serverEventData),
  })

  console.log("Server response after update:", updatedEvent)
  console.log("Available places from server:", updatedEvent.available_places)

  const transformedEvent = transformEventFromServer(updatedEvent)
  console.log("Transformed event after update:", transformedEvent)
  console.log(
    "Available places after transform:",
    transformedEvent.availablePlaces
  )

  return transformedEvent
}

export const deleteEvent = async (id: number): Promise<void> => {
  console.log(`Deleting event ${id}`)

  await apiRequestWithoutJson(`/events/${id}`, {
    method: "DELETE",
  })
}
