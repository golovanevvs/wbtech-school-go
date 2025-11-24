"use client"

import { useState, useEffect } from "react"
import { Box, Stack, Typography, Alert, Button } from "@mui/material"
import Header from "../../ui/Header"
import { getEventById } from "../../api/events"
import { Event } from "../../lib/types"

interface PageProps {
  params: {
    id: string
  }
}

export default function EventDetailPage({ params }: PageProps) {
  const [event, setEvent] = useState<Event | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchEvent = async () => {
      try {
        setLoading(true)
        const data = await getEventById(parseInt(params.id))
        setEvent(data)
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load event")
      } finally {
        setLoading(false)
      }
    }

    fetchEvent()
  }, [params.id])

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
        <Typography variant="h6">Загрузка мероприятия...</Typography>
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

  if (!event) {
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
          <Alert severity="error">Мероприятие не найдено</Alert>
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
      <Stack spacing={4}>
        <Header />
        <Box sx={{ maxWidth: 600, mx: "auto", p: 2 }}>
          <Typography variant="h4" component="h1" gutterBottom>
            {event.title}
          </Typography>
          <Typography variant="body1" paragraph>
            {event.description}
          </Typography>
          <Typography variant="body2" color="text.secondary" paragraph>
            Дата: {new Date(event.date).toLocaleString("ru-RU")}
          </Typography>
          <Typography variant="body2" color="text.secondary" paragraph>
            Всего мест: {event.totalPlaces} | Свободных: {event.availablePlaces}
          </Typography>
          <Typography variant="body2" color="text.secondary" paragraph>
            Срок бронирования: {event.bookingDeadline} минут
          </Typography>
          <Button
            variant="contained"
            disabled={event.availablePlaces === 0}
            onClick={() => {
              console.log("Book event:", event.id)
            }}
          >
            Забронировать место
          </Button>
        </Box>
      </Stack>
    </Box>
  )
}
