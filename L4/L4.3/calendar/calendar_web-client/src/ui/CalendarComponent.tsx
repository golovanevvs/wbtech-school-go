"use client"

import { useState, useEffect, useCallback, useMemo } from "react"
import FullCalendar from "@fullcalendar/react"
import dayGridPlugin from "@fullcalendar/daygrid"
import timeGridPlugin from "@fullcalendar/timegrid"
import interactionPlugin from "@fullcalendar/interaction"
import { Box, CircularProgress, Alert } from "@mui/material"
import { calendarApi } from "@/lib/api/calendar"
import { CalendarEvent, CreateEventRequest } from "@/lib/types/calendar"
import { EventDropArg } from "@fullcalendar/core"

interface CalendarComponentProps {
  onEventClick?: (event: CalendarEvent) => void
  onDateClick?: (date: Date) => void
}

// Simple logger function to avoid TypeScript issues with console
interface Logger {
  error: (...args: unknown[]) => void
  log: (...args: unknown[]) => void
}

const logger: Logger = {
  error: (...args: unknown[]) => {
    if (typeof window !== 'undefined' && window.console) {
      window.console.error(...args)
    }
  },
  log: (...args: unknown[]) => {
    if (typeof window !== 'undefined' && window.console) {
      window.console.log(...args)
    }
  }
}

export default function CalendarComponent({ onEventClick, onDateClick }: CalendarComponentProps) {
  const [events, setEvents] = useState<CalendarEvent[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [currentDate, setCurrentDate] = useState(() => new Date())

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

  // Обработчик перетаскивания события (eventDrop)
  const handleEventDrop = useCallback(async (info: EventDropArg) => {
    try {
      const event = info.event
      const eventId = event.id
      
      console.log('Событие перетащено:', eventId, 'новая дата:', event.start)
      
      // Подготавливаем данные для обновления
      const updatedEventData: Partial<CreateEventRequest> = {
        title: event.title,
        start: event.start ? event.start.toISOString() : '',
        end: event.end ? event.end.toISOString() : undefined,
        allDay: event.allDay,
      }
      
      // Обновляем событие на сервере
      const updatedEvent = await calendarApi.updateEvent(eventId, updatedEventData)
      
      // Обновляем локальное состояние
      setEvents(prevEvents => 
        prevEvents.map(evt => 
          evt.id === eventId 
            ? { ...evt, start: updatedEvent.start, end: updatedEvent.end }
            : evt
        )
      )
      
      logger.log('Событие успешно обновлено:', updatedEvent)
    } catch (err) {
      // Откатываем изменения в UI
      (info as any).revert()
      
      setError(err instanceof Error ? err.message : "Failed to update event")
    }
  }, [])

  // Обработчик изменения размера события (eventResize)
  const handleEventResize = useCallback(async (info: unknown) => {
    try {
      const event = (info as any).event
      const eventId = event.id
      
      console.log('Размер события изменен:', eventId, 'новое время окончания:', event.end)
      
      // Подготавливаем данные для обновления
      const updatedEventData: Partial<CreateEventRequest> = {
        title: event.title,
        start: event.start ? event.start.toISOString() : '',
        end: event.end ? event.end.toISOString() : undefined,
        allDay: event.allDay,
      }
      
      // Обновляем событие на сервере
      const updatedEvent = await calendarApi.updateEvent(eventId, updatedEventData)
      
      // Обновляем локальное состояние
      setEvents(prevEvents => 
        prevEvents.map(evt => 
          evt.id === eventId 
            ? { ...evt, start: updatedEvent.start, end: updatedEvent.end }
            : evt
        )
      )
      
      logger.log('Событие успешно обновлено после изменения размера:', updatedEvent)
    } catch (err) {
      // Откатываем изменения в UI
      (info as any).revert()
      
      setError(err instanceof Error ? err.message : "Failed to update event")
    }
  }, [])

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

  // Функция для обновления даты (вызывается вручной)
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
        editable={true} // Включаем редактирование
        selectable={true}
        selectMirror={true}
        dayMaxEvents={true}
        weekends={true}
        events={calendarEvents}
        // УБРАЛИ datesSet - это была причина бесконечного цикла!
        eventDrop={handleEventDrop} // Обработчик перетаскивания
        eventResize={handleEventResize} // Обработчик изменения размера
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