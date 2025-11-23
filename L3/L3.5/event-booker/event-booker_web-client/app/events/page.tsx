"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import { Box, Typography, Stack, Alert, Button } from "@mui/material"
import EventList from "../ui/events/EventList"
import { getEvents, deleteEvent } from "../api/events"
import { bookEvent, confirmBooking, cancelBooking, getUserBookings, getUserBookingByEventId } from "../api/bookings"
import { useAuth } from "../context/AuthContext"
import { Event, transformBookingFromServer, Booking } from "../lib/types"

type BookingInfo = {
  status: "pending" | "confirmed" | null
  expiresAt?: number | null
  bookingId?: number // Добавляем bookingId для быстрого доступа
}

export default function EventsPage() {
  const { user, loading: authLoading } = useAuth()
  const [events, setEvents] = useState<Event[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [bookingsMap, setBookingsMap] = useState<Record<number, BookingInfo>>({})
  const router = useRouter()

  // Загружаем мероприятия
  useEffect(() => {
    const fetchEvents = async () => {
      try {
        setLoading(true)
        const data = await getEvents()
        setEvents(data)
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load events")
      } finally {
        setLoading(false)
      }
    }

    fetchEvents()
  }, [])

  // Загружаем брони пользователя
  useEffect(() => {
    const fetchUserBookings = async () => {
      console.log("=== FETCH USER BOOKINGS DEBUG ===")
      console.log("authLoading:", authLoading)
      console.log("user:", user)
      
      // Если аутентификация еще загружается, не загружаем брони
      if (authLoading) {
        console.log("Auth still loading, skipping bookings fetch")
        return
      }

      if (!user) {
        console.log("No user, setting empty bookings map")
        setBookingsMap({})
        return
      }

      try {
        console.log("🔍 FETCHING USER BOOKINGS FROM SERVER...")
        const rawBookings = await getUserBookings()
        console.log("📦 Raw bookings from server:", rawBookings)
        console.log("📊 Bookings count:", rawBookings?.length || 0)
        
        // Трансформируем данные из snake_case в camelCase
        const bookings: Booking[] = rawBookings.map(transformBookingFromServer)
        console.log("🔄 Transformed bookings:", bookings)
        
        const bookingsMap: Record<number, BookingInfo> = {}
        
        bookings.forEach(booking => {
          console.log("Processing booking:", booking)
          if (booking.status === "pending" || booking.status === "confirmed") {
            const expiresAt = booking.expiresAt ? new Date(booking.expiresAt).getTime() : null
            
            if (booking.eventId) {
              bookingsMap[booking.eventId] = {
                status: booking.status,
                expiresAt: expiresAt,
                bookingId: booking.id // Сохраняем bookingId
              }
              console.log(`Added booking for event ${booking.eventId}:`, bookingsMap[booking.eventId])
            } else {
              console.warn("Booking has no eventId:", booking)
            }
          }
        })
        
        console.log("Final bookingsMap:", bookingsMap)
        setBookingsMap(bookingsMap)
        console.log("=== END FETCH USER BOOKINGS ===")
      } catch (err) {
        console.error("Failed to load user bookings:", err)
      }
    }

    fetchUserBookings()
  }, [user, authLoading])

  const handleCreateEvent = () => {
    if (!user) {
      router.push("/auth?mode=login")
      return
    }
    router.push("/events/create")
  }

  const handleBookEvent = async (eventId: number) => {
    if (!user) {
      router.push("/auth?mode=login")
      return
    }

    // Проверяем, есть ли уже бронь у пользователя на это мероприятие
    const existingBooking = bookingsMap[eventId]
    if (existingBooking && (existingBooking.status === "pending" || existingBooking.status === "confirmed")) {
      setError("У вас уже есть бронь на это мероприятие")
      return
    }

    try {
      const booking = await bookEvent({ event_id: eventId })
      
      // Обновляем карту брони с новой информацией
      setBookingsMap(prev => ({
        ...prev,
        [eventId]: {
          status: "pending",
          expiresAt: new Date(booking.expiresAt).getTime(),
          bookingId: booking.id
        }
      }))

      console.log("Booking created:", booking)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to book event")
    }
  }

  const handleConfirmBooking = async (eventId: number) => {
    if (!user) {
      router.push("/auth?mode=login")
      return
    }

    try {
      // Проверяем, есть ли bookingId в bookingsMap
      const existingBooking = bookingsMap[eventId]
      let bookingId: number

      if (existingBooking?.bookingId) {
        // Если bookingId уже есть в состоянии, используем его
        bookingId = existingBooking.bookingId
      } else {
        // Иначе получаем бронь через API
        const booking = await getUserBookingByEventId(eventId)
        
        if (!booking) {
          setError("Бронь не найдена")
          return
        }
        bookingId = booking.id
      }

      // Подтверждаем бронь по её ID
      const confirmedBooking = await confirmBooking(bookingId)
      
      // Обновляем статус брони на подтвержденную
      setBookingsMap(prev => ({
        ...prev,
        [eventId]: {
          status: "confirmed",
          expiresAt: null,
          bookingId: bookingId // Сохраняем bookingId
        }
      }))

      console.log("Booking confirmed:", confirmedBooking)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to confirm booking")
    }
  }

  const handleCancelBooking = async (eventId: number) => {
    if (!user) {
      router.push("/auth?mode=login")
      return
    }

    try {
      // Проверяем, есть ли bookingId в bookingsMap
      const existingBooking = bookingsMap[eventId]
      let bookingId: number

      if (existingBooking?.bookingId) {
        // Если bookingId уже есть в состоянии, используем его
        bookingId = existingBooking.bookingId
      } else {
        // Иначе получаем бронь через API
        const booking = await getUserBookingByEventId(eventId)
        
        if (!booking) {
          setError("Бронь не найдена")
          return
        }
        bookingId = booking.id
      }

      // Отменяем бронь по её ID
      const cancelledBooking = await cancelBooking(bookingId)
      
      // Удаляем бронь из bookingsMap (так как она отменена)
      setBookingsMap(prev => {
        const newMap = { ...prev }
        delete newMap[eventId]
        return newMap
      })

      console.log("Booking cancelled:", cancelledBooking)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to cancel booking")
    }
  }

  const handleEditEvent = (eventId: number) => {
    if (!user) {
      router.push("/auth?mode=login")
      return
    }
    router.push(`/events/${eventId}/edit`)

  }

  const handleDeleteEvent = async (eventId: number) => {
    if (!user) {
      router.push("/auth?mode=login")
      return
    }

    if (window.confirm("Вы уверены, что хотите удалить это мероприятие?")) {
      try {
        await deleteEvent(eventId)
        // Обновляем список мероприятий
        setEvents(events.filter(event => event.id !== eventId))
        // Удаляем бронь из карты, если она была
        setBookingsMap(prev => {
          const newMap = { ...prev }
          delete newMap[eventId]
          return newMap
        })
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to delete event")
      }
    }
  }

  

  if (loading) {
    return (
      <Box
        sx={{
          width: "100%",
          minHeight: "100vh",
          px: { xs: 0, sm: 2 },
          py: 2,
          bgcolor: "background.default",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <Typography variant="h6">Загрузка мероприятий...</Typography>
      </Box>
    )
  }

  if (error) {
    return (
      <Box
        sx={{
          width: "100%",
          minHeight: "100vh",
          px: { xs: 0, sm: 2 },
          py: 2,
          bgcolor: "background.default",
        }}
      >
        <Stack spacing={4}>
          <Alert severity="error">{error}</Alert>
        </Stack>
      </Box>
    )
  }

  return (
    <Box
      sx={{
        width: "100%",
        minHeight: "100vh",
        px: { xs: 0, sm: 2 },
        py: 2,
        bgcolor: "background.default",
        maxWidth: "100vw",
        mx: "auto",
      }}
    >
      <Stack spacing={4} alignItems="center">
        <Box sx={{ width: "100%", maxWidth: 1200, px: { xs: 0, sm: 2 } }}>
          <Typography variant="h2" align="center" sx={{ mb: 2 }}>
            Мероприятия
          </Typography>
          
          {user && (
            <Box sx={{ textAlign: "center", mb: 2 }}>
              <Button 
                variant="contained" 
                onClick={handleCreateEvent}
                sx={{ mb: 2 }}
              >
                Создать мероприятие
              </Button>
            </Box>
          )}
          
          <EventList 
            events={events}
            onBook={handleBookEvent}
            onConfirmBooking={handleConfirmBooking}
            onCancelBooking={handleCancelBooking}
            onEdit={handleEditEvent}
            onDelete={handleDeleteEvent}
            currentUserId={user?.id}
            bookingsMap={bookingsMap}
          />
        </Box>
      </Stack>
    </Box>
  )
}