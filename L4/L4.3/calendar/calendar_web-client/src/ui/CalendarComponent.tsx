"use client"

import { useState, useEffect, useCallback, useMemo, useRef } from "react"
import FullCalendar from "@fullcalendar/react"
import dayGridPlugin from "@fullcalendar/daygrid"
import timeGridPlugin from "@fullcalendar/timegrid"
import interactionPlugin from "@fullcalendar/interaction"
import { Box, CircularProgress, Alert } from "@mui/material"
import { calendarApi } from "@/lib/api/calendar"
import { CalendarEvent, CreateEventRequest } from "@/lib/types/calendar"
import { 
  FullCalendarEventDropInfo, 
  FullCalendarEventResizeInfo, 
  FullCalendarEventResizeStartInfo, 
  FullCalendarEventResizeStopInfo,
  DateUtils
} from "@/lib/types/fullcalendar-types"

interface CalendarComponentProps {
  onEventClick?: (event: CalendarEvent) => void
  onDateClick?: (date: Date) => void
}

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
  const calendarRef = useRef<FullCalendar>(null)

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

  const handleEventDrop = useCallback(async (info: FullCalendarEventDropInfo) => {
    const event = info.event
    const eventId = event.id || ''
    
    console.log('Событие перетащено:', eventId, 'новая дата:', event.start)
    
    const oldEventData = events.find(evt => evt.id === eventId)
    
    const newStart = DateUtils.toISOString(event.start)
    const newEnd = DateUtils.toISOString(event.end)
    
    setEvents(prevEvents => 
      prevEvents.map(evt => 
        evt.id === eventId 
          ? { ...evt, start: newStart || evt.start, end: newEnd || evt.end }
          : evt
      )
    )
    
    setTimeout(() => {
      if (calendarRef.current) {
        const calendarApi = calendarRef.current.getApi()
        calendarApi.refetchEvents()
      }
    }, 0)
    
    try {
      const updatedEventData: Partial<CreateEventRequest> = {
        title: event.title,
        start: newStart || undefined,
        end: newEnd || undefined,
        allDay: event.allDay || false,
      }
      
      const updatedEvent = await calendarApi.updateEvent(eventId, updatedEventData)
      
      setEvents(prevEvents => 
        prevEvents.map(evt => 
          evt.id === eventId ? updatedEvent : evt
        )
      )
      
      logger.log('Событие успешно обновлено:', updatedEvent)
    } catch (err) {
      if (oldEventData) {
        setEvents(prevEvents => 
          prevEvents.map(evt => 
            evt.id === eventId ? oldEventData : evt
          )
        )
        
        setTimeout(() => {
          if (calendarRef.current) {
            const calendarApi = calendarRef.current.getApi()
            calendarApi.refetchEvents()
          }
        }, 0)
      }
      
      setError(err instanceof Error ? err.message : "Failed to update event")
      logger.error('Ошибка обновления события:', err)
    }
  }, [events])

  const handleEventResize = useCallback(async (info: FullCalendarEventResizeInfo) => {
    const event = info.event
    const eventId = event.id || ''
    
    console.log('Размер события изменен:', eventId, 'новое время окончания:', event.end)
    
    const oldEventData = events.find(evt => evt.id === eventId)
    
    const newStart = DateUtils.toISOString(event.start)
    const newEnd = DateUtils.toISOString(event.end)
    
    setEvents(prevEvents => 
      prevEvents.map(evt => 
        evt.id === eventId 
          ? { ...evt, start: newStart || evt.start, end: newEnd || evt.end }
          : evt
      )
    )
    
    setTimeout(() => {
      if (calendarRef.current) {
        const calendarApi = calendarRef.current.getApi()
        calendarApi.refetchEvents()
      }
    }, 0)
    
    try {
      const updatedEventData: Partial<CreateEventRequest> = {
        title: event.title,
        start: newStart || undefined,
        end: newEnd || undefined,
        allDay: event.allDay || false,
      }
      
      const updatedEvent = await calendarApi.updateEvent(eventId, updatedEventData)
      
      setEvents(prevEvents => 
        prevEvents.map(evt => 
          evt.id === eventId ? updatedEvent : evt
        )
      )
      
      logger.log('Событие успешно обновлено после изменения размера:', updatedEvent)
    } catch (err) {
      if (oldEventData) {
        setEvents(prevEvents => 
          prevEvents.map(evt => 
            evt.id === eventId ? oldEventData : evt
          )
        )
        
        setTimeout(() => {
          if (calendarRef.current) {
            const calendarApi = calendarRef.current.getApi()
            calendarApi.refetchEvents()
          }
        }, 0)
      }
      
      setError(err instanceof Error ? err.message : "Failed to update event")
      logger.error('Ошибка изменения размера события:', err)
    }
  }, [events])

  const handleEventResizeStart = useCallback((info: FullCalendarEventResizeStartInfo) => {
    console.log('Начало изменения размера события:', info.event.title)
  }, [])

  const handleEventResizeStop = useCallback((info: FullCalendarEventResizeStopInfo) => {
    console.log('Окончание изменения размера события:', info.event.title)
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

  const handleDateChange = useCallback((newDate: Date) => {
    console.log('Меняем дату на:', newDate)
    setCurrentDate(newDate)
    loadEvents(newDate)
  }, [loadEvents])

  const calendarEvents = useMemo(() => {
    console.log('Формируем calendarEvents, events:', events)
    if (!events || events.length === 0) {
      return testEvents 
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
        ref={calendarRef}
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
        eventDrop={handleEventDrop} 
        eventResize={handleEventResize} 
        eventResizeStart={handleEventResizeStart} 
        eventResizeStop={handleEventResizeStop} 
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
        datesSet={(arg) => {
          console.log('FullCalendar dates set:', arg.start)
        }}
      />
      
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