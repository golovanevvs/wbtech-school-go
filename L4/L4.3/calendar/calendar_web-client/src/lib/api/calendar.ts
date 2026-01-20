import { apiClient } from "./client"
import { CalendarEvent, CreateEventRequest, MonthEventsResponse, CreateEventResponse } from "@/lib/types/calendar"

export const calendarApi = {
  async getMonthEvents(year: number, month: number): Promise<CalendarEvent[]> {
    const response = await apiClient.request<MonthEventsResponse>(
      `/events/month?year=${year}&month=${month}`
    )
    return response.events
  },

  async createEvent(eventData: CreateEventRequest): Promise<CalendarEvent> {
    const response = await apiClient.request<CreateEventResponse>(
      "/event/create",
      {
        method: "POST",
        body: JSON.stringify(eventData),
      }
    )
    return response.event
  },

  async getEvent(id: string): Promise<CalendarEvent> {
    return await apiClient.request<CalendarEvent>(`/event/${id}`)
  },

  async updateEvent(id: string, eventData: Partial<CreateEventRequest>): Promise<CalendarEvent> {
    const response = await apiClient.request<CreateEventResponse>(
      `/event/${id}`,
      {
        method: "PUT",
        body: JSON.stringify(eventData),
      }
    )
    return response.event
  },

  async deleteEvent(id: string): Promise<void> {
    await apiClient.request(`/event/${id}`, {
      method: "DELETE",
    })
  },

  async getDayEvents(date: string): Promise<CalendarEvent[]> {
    const response = await apiClient.request<{ events: CalendarEvent[] }>(
      `/events/day?date=${date}`
    )
    return response.events
  },
}