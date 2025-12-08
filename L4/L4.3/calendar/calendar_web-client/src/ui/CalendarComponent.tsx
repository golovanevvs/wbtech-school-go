"use client"

import { useState, useEffect, useCallback, useMemo } from "react"
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
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [currentDate, setCurrentDate] = useState(new Date())

  const loadEvents = useCallback(async (date: Date) => {
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
  }, [])

  useEffect(() => {
    // Загружаем события для текущего месяца при монтировании
    loadEvents(new Date())
  }, [loadEvents])

  const handleDatesSet = useCallback((arg: { start: Date }) => {
    // Обновляем только если дата действительно изменилась
    const newDate = arg.start
    const currentDateStr = currentDate.toISOString().split('T')[0]
    const newDateStr = newDate.toISOString().split('T')[0]
    
    if (currentDateStr !== newDateStr) {
      setCurrentDate(newDate)
      loadEvents(newDate)
    }
  }, [currentDate, loadEvents])

  const handleEventClick = useCallback((clickInfo: { event: { id: string } }) => {
    const event = events.find(e => e.id === clickInfo.event.id)
    if (event && onEventClick) {
      onEventClick(event)
    }
  }, [events, onEventClick])

  const handleDateClick = useCallback((clickInfo: { date: Date }) => {
    if (onDateClick) {
      onDateClick(clickInfo.date)
    }
  }, [onDateClick])

  const calendarEvents = useMemo(() => {
    if (!events) return []
    return events.map(event => ({
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
    }))
  }, [events])

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