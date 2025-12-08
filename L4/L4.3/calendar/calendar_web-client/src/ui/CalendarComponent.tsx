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

  // Тестовые события для проверки работы
  const testEvents = useMemo(() => [
    {
      id: '1',
      title: 'Тестовое событие',
      start: new Date().toISOString(),
    }
  ], [])

  const loadEvents = useCallback(async (date: Date) => {
    try {
      setLoading(true)
      setError(null)
      
      console.log('Загружаем события для даты:', date)
      const year = date.getFullYear()
      const month = date.getMonth() + 1
      
      const monthEvents = await calendarApi.getMonthEvents(year, month)
      console.log('Получены события:', monthEvents)
      setEvents(monthEvents)
    } catch (err) {
      console.error('Failed to load events:', err)
      setError(err instanceof Error ? err.message : "Failed to load events")
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    console.log('CalendarComponent mounted, loading events...')
    loadEvents(currentDate)
  }, [loadEvents, currentDate])

  // УБРАЛИ datesSet - это была причина бесконечного цикла!

  const handleEventClick = useCallback((clickInfo: { event: { id: string } }) => {
    const event = events.find(e => e.id === clickInfo.event.id)
    if (event && onEventClick) {
      onEventClick(event)
    }
  }, [events, onEventClick])

  const handleDateClick = useCallback((clickInfo: { date: Date }) => {
    console.log('Date clicked:', clickInfo.date)
    if (onDateClick) {
      onDateClick(clickInfo.date)
    }
  }, [onDateClick])

  // Функция для обновления даты (вызывается вручную)
  const handleDateChange = useCallback((newDate: Date) => {
    console.log('Меняем дату на:', newDate)
    setCurrentDate(newDate)
    loadEvents(newDate)
  }, [loadEvents])

  const calendarEvents = useMemo(() => {
    console.log('Формируем calendarEvents, events:', events)
    if (!events || events.length === 0) {
      return testEvents // Используем тестовые события если нет реальных
    }
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
  }, [events, testEvents])

  console.log('CalendarComponent render, events:', events, 'calendarEvents:', calendarEvents)

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
        initialDate={currentDate}
        editable={true}
        selectable={true}
        selectMirror={true}
        dayMaxEvents={true}
        weekends={true}
        events={calendarEvents}
        // УБРАЛИ datesSet - это была причина бесконечного цикла!
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
        navLinks={true}
        // Добавим ручную навигацию через кнопки
        datesSet={(arg) => {
          // Только логирование, без обновления состояния
          console.log('FullCalendar dates set:', arg.start)
        }}
      />
      
      {/* Ручные кнопки для навигации */}
      <Box sx={{ mt: 2, display: 'flex', gap: 1, justifyContent: 'center' }}>
        <button onClick={() => {
          const newDate = new Date(currentDate)
          newDate.setMonth(newDate.getMonth() - 1)
          handleDateChange(newDate)
        }}>
          Предыдущий месяц
        </button>
        <button onClick={() => handleDateChange(new Date())}>
          Сегодня
        </button>
        <button onClick={() => {
          const newDate = new Date(currentDate)
          newDate.setMonth(newDate.getMonth() + 1)
          handleDateChange(newDate)
        }}>
          Следующий месяц
        </button>
      </Box>
    </Box>
  )
}