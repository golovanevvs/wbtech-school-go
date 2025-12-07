"use client"

import { useState, useEffect } from "react"
import FullCalendar from "@fullcalendar/react"
import dayGridPlugin from "@fullcalendar/daygrid"
import timeGridPlugin from "@fullcalendar/timegrid"
import interactionPlugin from "@fullcalendar/interaction"
import { Box, CircularProgress, Alert } from "@mui/material"
import { calendarApi } from "@/lib/api/calendar"
import { CalendarEvent } from "@/lib/types/calendar"

interface CalendarComponentProps {
  onEventClick?: (event: CalendarEvent) => void
  onDateClick?: (date: Date) => void
}

export default function CalendarComponent({ onEventClick, onDateClick }: CalendarComponentProps) {
  const [events, setEvents] = useState<CalendarEvent[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [currentDate, setCurrentDate] = useState(new Date())

  const loadEvents = async (date: Date) => {
    try {
      setLoading(true)
      setError(null)
      
      const year = date.getFullYear()
      const month = date.getMonth() + 1 // FullCalendar uses 0-based months
      
      const monthEvents = await calendarApi.getMonthEvents(year, month)
      setEvents(monthEvents)
    } catch (err) {
      console.error("Failed to load events:", err)
      setError(err instanceof Error ? err.message : "Failed to load events")
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadEvents(currentDate)
  }, [currentDate])

  const handleDatesSet = (arg: { start: Date }) => {
    setCurrentDate(arg.start)
  }

  const handleEventClick = (clickInfo: { event: { id: string } }) => {
    const event = events.find(e => e.id === clickInfo.event.id)
    if (event && onEventClick) {
      onEventClick(event)
    }
  }

  const handleDateClick = (clickInfo: { date: Date }) => {
    if (onDateClick) {
      onDateClick(clickInfo.date)
    }
  }

  const calendarEvents = events?.map(event => ({
    id: event.id,
    title: event.title,
    start: event.start,
    end: event.end,
    allDay: event.allDay || false,
    extendedProps: {
      description: event.description,
      reminder: event.reminder,
      reminderTime: event.reminderTime,
    },
  })) || []

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
      </Box>
    )
  }

  if (error) {
    return (
      <Alert severity="error" sx={{ mb: 2 }}>
        {error}
      </Alert>
    )
  }

  return (
    <Box>
      <FullCalendar
        plugins={[dayGridPlugin, timeGridPlugin, interactionPlugin]}
        headerToolbar={{
          left: "prev,next today",
          center: "title",
          right: "dayGridMonth,timeGridWeek,timeGridDay",
        }}
        initialView="dayGridMonth"
        editable={true}
        selectable={true}
        selectMirror={true}
        dayMaxEvents={true}
        weekends={true}
        events={calendarEvents}
        datesSet={handleDatesSet}
        eventClick={handleEventClick}
        dateClick={handleDateClick}
        height="auto"
        locale="ru"
        buttonText={{
          today: "Сегодня",
          month: "Месяц",
          week: "Неделя",
          day: "День",
          prev: "Назад",
          next: "Вперед",
        }}
        eventDisplay="block"
        eventTimeFormat={{
          hour: "2-digit",
          minute: "2-digit",
          meridiem: false,
        }}
      />
    </Box>
  )
}