"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import { Box, Typography, Stack, Alert, Button } from "@mui/material"
import EventList from "../ui/events/EventList"
import { getEvents, deleteEvent } from "../api/events"
import {
  bookEvent,
  confirmBooking,
  cancelBooking,
  getUserBookings,
  getUserBookingByEventId,
} from "../api/bookings"
import { useAuth } from "../context/AuthContext"
import { Event, transformBookingFromServer, Booking } from "../lib/types"

type BookingInfo = {
  status: "pending" | "confirmed" | null
  expiresAt?: number | null
  bookingId?: number
}

export default function EventsPage() {
  const { user, loading: authLoading } = useAuth()
  const [events, setEvents] = useState<Event[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [bookingErrors, setBookingErrors] = useState<Record<number, string>>({})
  const [bookingsMap, setBookingsMap] = useState<Record<number, BookingInfo>>(
    {}
  )
  const router = useRouter()

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

  useEffect(() => {
    const fetchUserBookings = async () => {
      console.log("=== FETCH USER BOOKINGS DEBUG ===")
      console.log("authLoading:", authLoading)
      console.log("user:", user)

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

        const bookings: Booking[] = rawBookings.map(transformBookingFromServer)
        console.log("🔄 Transformed bookings:", bookings)

        const bookingsMap: Record<number, BookingInfo> = {}

        bookings.forEach((booking) => {
          console.log("Processing booking:", booking)
          if (booking.status === "pending" || booking.status === "confirmed") {
            const expiresAt = booking.expiresAt
              ? new Date(booking.expiresAt).getTime()
              : null

            if (booking.eventId) {
              bookingsMap[booking.eventId] = {
                status: booking.status,
                expiresAt: expiresAt,
                bookingId: booking.id,
              }
              console.log(
                `Added booking for event ${booking.eventId}:`,
                bookingsMap[booking.eventId]
              )
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

    setBookingErrors(prev => ({ ...prev, [eventId]: "" }))

    const existingBooking = bookingsMap[eventId]
    if (
      existingBooking &&
      (existingBooking.status === "pending" ||
        existingBooking.status === "confirmed")
    ) {
      setBookingErrors(prev => ({
        ...prev,
        [eventId]: "У вас уже есть бронь на это мероприятие",
      }))
      return
    }

    try {
      const booking = await bookEvent({ event_id: eventId })

      setBookingsMap((prev) => ({
        ...prev,
        [eventId]: {
          status: "pending",
          expiresAt: new Date(booking.expiresAt).getTime(),
          bookingId: booking.id,
        },
      }))

      const updatedEvents = await getEvents()
      setEvents(updatedEvents)

      console.log("Booking created:", booking)
    } catch (err) {
      setBookingErrors(prev => ({
        ...prev,
        [eventId]: err instanceof Error ? err.message : "Failed to book event",
      }))
    }
  }

  const handleConfirmBooking = async (eventId: number) => {
    if (!user) {
      router.push("/auth?mode=login")
      return
    }

    setBookingErrors(prev => ({ ...prev, [eventId]: "" }))

    try {
      const existingBooking = bookingsMap[eventId]
      let bookingId: number

      if (existingBooking?.bookingId) {
        bookingId = existingBooking.bookingId
      } else {
        const booking = await getUserBookingByEventId(eventId)

        if (!booking) {
          setBookingErrors(prev => ({
            ...prev,
            [eventId]: "Бронь не найдена",
          }))
          return
        }
        bookingId = booking.id
      }

      const confirmedBooking = await confirmBooking(bookingId)

      setBookingsMap((prev) => ({
        ...prev,
        [eventId]: {
          status: "confirmed",
          expiresAt: null,
          bookingId: bookingId,
        },
      }))

      console.log("Booking confirmed:", confirmedBooking)
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to confirm booking"
      if (errorMessage.includes("not in pending status")) {
        setBookingErrors(prev => ({
          ...prev,
          [eventId]: "Срок брони истёк. Бронь была отменена.",
        }))
        try {
          await getUserBookings()
          
          setBookingsMap((prev) => {
            const newMap = { ...prev }
            delete newMap[eventId]
            return newMap
          })
        } catch {
        }
      } else {
        setBookingErrors(prev => ({
          ...prev,
          [eventId]: errorMessage,
        }))
      }
    }
  }

  const handleCancelBooking = async (eventId: number) => {
    if (!user) {
      router.push("/auth?mode=login")
      return
    }

    setBookingErrors(prev => ({ ...prev, [eventId]: "" }))

    try {
      const existingBooking = bookingsMap[eventId]
      let bookingId: number

      if (existingBooking?.bookingId) {
        bookingId = existingBooking.bookingId
      } else {
        const booking = await getUserBookingByEventId(eventId)

        if (!booking) {
          setBookingErrors(prev => ({
            ...prev,
            [eventId]: "Бронь не найдена",
          }))
          return
        }
        bookingId = booking.id
      }

      const cancelledBooking = await cancelBooking(bookingId)

      setBookingsMap((prev) => {
        const newMap = { ...prev }
        delete newMap[eventId]
        return newMap
      })

      const updatedEvents = await getEvents()
      setEvents(updatedEvents)

      console.log("Booking cancelled:", cancelledBooking)
    } catch (err) {
      setBookingErrors(prev => ({
        ...prev,
        [eventId]: err instanceof Error ? err.message : "Failed to cancel booking",
      }))
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
        setEvents(events.filter((event) => event.id !== eventId))
        setBookingsMap((prev) => {
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
          <Alert severity="error">
            {error}
            <Button onClick={() => setError(null)} sx={{ mt: 1 }} size="small">
              Закрыть
            </Button>
          </Alert>
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
            bookingErrors={bookingErrors}
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
