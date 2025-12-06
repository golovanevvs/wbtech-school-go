export interface CalendarEvent {
  id: string
  title: string
  description?: string
  start: string
  end?: string
  allDay?: boolean
  reminder?: boolean
  reminderTime?: string
  createdAt: string
  updatedAt: string
}

export interface CreateEventRequest {
  title: string
  description?: string
  start: string
  end?: string
  allDay?: boolean
  reminder?: boolean
  reminderTime?: string
}

export interface MonthEventsResponse {
  events: CalendarEvent[]
}

export interface CreateEventResponse {
  event: CalendarEvent
}