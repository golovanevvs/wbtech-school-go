"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import { Box, Stack, Typography, Alert, Button } from "@mui/material"
import BookingList from "../ui/bookings/BookingList"
import { getUserBookings, confirmBooking, cancelBooking } from "../api/bookings"
import { useAuth } from "../context/AuthContext"
import { Booking, transformBookingFromServer } from "../lib/types"

export default function BookingsPage() {
  const { user, loading: authLoading } = useAuth()
  const [bookings, setBookings] = useState<Booking[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const router = useRouter()

  useEffect(() => {
    if (authLoading) return

    if (!user) {
      router.push("/auth")
      return
    }

    const fetchBookings = async () => {
      try {
        setLoading(true)
        const data = await getUserBookings()
        setBookings(data.map(transformBookingFromServer))
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load bookings")
      } finally {
        setLoading(false)
      }
    }

    fetchBookings()
  }, [user, authLoading, router])

  const handleConfirm = async (id: number) => {
    try {
      await confirmBooking(id)
      const updatedBookings = await getUserBookings()
      setBookings(updatedBookings.map(transformBookingFromServer))
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to confirm booking"
      // Если бронь не в статусе pending (истекла или уже подтверждена), обновляем список
      if (errorMessage.includes("not in pending status")) {
        const updatedBookings = await getUserBookings()
        setBookings(updatedBookings.map(transformBookingFromServer))
        setError("Срок брони истёк. Бронь была отменена.")
      } else {
        setError(errorMessage)
      }
    }
  }

  const handleCancel = async (id: number) => {
    try {
      await cancelBooking(id)
      const updatedBookings = await getUserBookings()
      setBookings(updatedBookings.map(transformBookingFromServer))
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to cancel booking")
    }
  }

  if (authLoading || loading) {
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
        <Typography variant="h6">Загрузка бронирований...</Typography>
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
          <Alert severity="error" sx={{ mx: 2 }}>
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
        <Typography variant="h2" align="center">
          Мои бронирования
        </Typography>
        {bookings.length > 0 ? (
          <BookingList
            bookings={bookings}
            onConfirm={handleConfirm}
            onCancel={handleCancel}
          />
        ) : (
          <Typography variant="body1" align="center">
            У вас нет активных бронирований
          </Typography>
        )}
      </Stack>
    </Box>
  )
}
