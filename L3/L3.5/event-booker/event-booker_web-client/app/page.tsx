"use client"

import { useState, useEffect } from "react"
import { Box, Stack, Typography } from "@mui/material"
import { getEvents } from "./api/events"
import { Event } from "./lib/types"
import EventCard from "./ui/events/EventCard"

export default function HomePage() {
  const [events, setEvents] = useState<Event[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchEvents = async () => {
      try {
        const data = await getEvents()
        setEvents(data.slice(0, 6))
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
        <Box sx={{ textAlign: "center", px: 2 }}>
          <Typography variant="h4" sx={{ mb: 1 }}>
            Event Booker
          </Typography>
          <Typography variant="h6" color="text.secondary" sx={{ mb: 3 }}>
            Сервис бронирования мероприятий
          </Typography>
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
                <EventCard key={event.id} event={event} />
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