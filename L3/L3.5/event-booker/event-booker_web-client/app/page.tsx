"use client"

import { useState, useEffect } from "react"
import { Box, Stack, Typography, Button } from "@mui/material"
import { getEvents } from "./api/events"
import { Event } from "./lib/types"

export default function HomePage() {
  const [events, setEvents] = useState<Event[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchEvents = async () => {
      try {
        const data = await getEvents()
        // Берем только 3 последних события
        setEvents(data.slice(0, 3))
      } catch (err) {
        console.error("Failed to load events:", err)
        setEvents([])
      } finally {
        setLoading(false)
      }
    }

    fetchEvents()
  }, [])

  return (
    <Box
      sx={{
        width: "100%",
        px: { xs: 0, sm: 2 },
        py: 2,
        bgcolor: "background.default",
        maxWidth: "100vw",
        mx: "auto",
      }}
    >
      <Stack spacing={4} alignItems="center">
        <Box sx={{ textAlign: "center", px: 2, maxWidth: 320 }}>
          <Typography variant="h6" color="text.secondary" sx={{ mb: 3 }}>
            Сервис бронирования мероприятий
          </Typography>
          <Stack spacing={2} alignItems="center" width="100%">
            <Button
              variant="contained"
              size="large"
              href="/events"
              sx={{ mt: 2, width: "100%" }}
            >
              Просмотреть мероприятия
            </Button>
            <Button
              variant="outlined"
              size="large"
              href="/auth"
              sx={{ mt: 2, width: "100%" }}
            >
              Войти / Зарегистрироваться
            </Button>
          </Stack>
        </Box>

        <Box sx={{ width: "100%", maxWidth: 1200, px: 2 }}>
          <Typography variant="h4" align="center" sx={{ mb: 2 }}>
            Ближайшие мероприятия
          </Typography>
          {loading ? (
            <Typography variant="h6" align="center">
              Загрузка мероприятий...
            </Typography>
          ) : events.length > 0 ? (
            <Box
              sx={{
                display: "grid",
                gridTemplateColumns: "repeat(auto-fit, minmax(300px, 1fr))",
                gap: 2,
                width: "100%",
              }}
            >
              {events.map((event) => (
                <Box
                  key={event.id}
                  sx={{
                    p: 2,
                    border: "1px solid #ccc",
                    borderRadius: 1,
                    minHeight: 120,
                  }}
                >
                  <Typography variant="h6">{event.title}</Typography>
                  <Typography variant="body2" color="text.secondary">
                    {new Date(event.date).toLocaleString("ru-RU")}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Свободных мест: {event.availablePlaces}
                  </Typography>
                </Box>
              ))}
            </Box>
          ) : (
            <Typography variant="h6" align="center">
              Нет доступных мероприятий
            </Typography>
          )}
        </Box>
      </Stack>
    </Box>
  )
}
