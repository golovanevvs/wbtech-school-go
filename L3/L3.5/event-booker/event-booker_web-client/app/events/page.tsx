"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import { Box, Typography, Stack, Alert, Button } from "@mui/material"
import EventList from "../ui/events/EventList"
import { getEvents, deleteEvent } from "../api/events"
import { bookEvent, confirmBooking, getUserBookings } from "../api/bookings"
import { useAuth } from "../context/AuthContext"
import { Event } from "../lib/types"

type BookingInfo = {
  status: "pending" | "confirmed" | null
  expiresAt?: number | null
}

export default function EventsPage() {
  const { user } = useAuth()
  const [events, setEvents] = useState<Event[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [bookingsMap, setBookingsMap] = useState<Record<number, BookingInfo>>({})
  const [currentTime, setCurrentTime] = useState(Date.now())
  const router = useRouter()

  // Обновляем текущее время каждую секунду для отслеживания истечения брони
  useEffect(() => {
    const interval = setInterval(() => {
      setCurrentTime(Date.now())
    }, 1000)

    return () => clearInterval(interval)
  }, [])

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
      if (!user) {
        setBookingsMap({})
        return
      }

      try {
        const bookings = await getUserBookings()
        const bookingsMap: Record<number, BookingInfo> = {}
        
        bookings.forEach(booking => {
          if (booking.status === "pending" || booking.status === "confirmed") {
            const expiresAt = new Date(booking.expiresAt).getTime()
            bookingsMap[booking.eventId] = {
              status: booking.status,
              expiresAt: expiresAt
            }
          }
        })
        
        setBookingsMap(bookingsMap)
      } catch (err) {
        console.error("Failed to load user bookings:", err)
      }
    }

    fetchUserBookings()
  }, [user])

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

    try {
      const booking = await bookEvent({ eventId })
      
      // Обновляем карту брони с новой информацией
      setBookingsMap(prev => ({
        ...prev,
        [eventId]: {
          status: "pending",
          expiresAt: new Date(booking.expiresAt).getTime()
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
      const booking = await confirmBooking(eventId)
      
      // Обновляем статус брони на подтвержденную
      setBookingsMap(prev => ({
        ...prev,
        [eventId]: {
          status: "confirmed",
          expiresAt: null
        }
      }))

      console.log("Booking confirmed:", booking)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to confirm booking")
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

  // Фильтруем истекшие брони
  useEffect(() => {
    setBookingsMap(prev => {
      const filtered = { ...prev }
      Object.keys(filtered).forEach(eventId => {
        const bookingInfo = filtered[parseInt(eventId)]
        if (bookingInfo.status === "pending" && 
            bookingInfo.expiresAt && 
            currentTime > bookingInfo.expiresAt) {
          // Бронь истекла, удаляем её из карты
          delete filtered[parseInt(eventId)]
        }
      })
      return filtered
    })
  }, [currentTime])

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
        <Box sx={{ width: "100%", maxWidth: 1200, px: 2 }}>
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
            onEdit={handleEditEvent}
            onDelete={handleDeleteEvent}
            currentUserId={user?.id}
            bookingsMap={bookingsMap}
            currentTime={currentTime}
          />
        </Box>
      </Stack>
    </Box>
  )
}