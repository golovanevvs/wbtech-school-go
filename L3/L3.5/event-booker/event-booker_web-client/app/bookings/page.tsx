"use client"

import { useState, useEffect } from "react"
import { Box, Stack, Typography, Alert, Button } from "@mui/material"
import Header from "../ui/Header"
import BookingList from "../ui/bookings/BookingList"
import { getUserBookings, confirmBooking, cancelBooking } from "../api/bookings"
import { Booking } from "../lib/types"

export default function BookingsPage() {
  const [bookings, setBookings] = useState<Booking[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchBookings = async () => {
      try {
        setLoading(true)
        const data = await getUserBookings()
        setBookings(data)
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load bookings")
      } finally {
        setLoading(false)
      }
    }

    fetchBookings()
  }, [])

  const handleConfirm = async (id: number) => {
    try {
      await confirmBooking(id)
      // Обновляем список бронирований
      const updatedBookings = await getUserBookings()
      setBookings(updatedBookings)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to confirm booking")
    }
  }

  const handleCancel = async (id: number) => {
    try {
      await cancelBooking(id)
      // Обновляем список бронирований
      const updatedBookings = await getUserBookings()
      setBookings(updatedBookings)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to cancel booking")
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
          <Header />
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
        <Header />
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
