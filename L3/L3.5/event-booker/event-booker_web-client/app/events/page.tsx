"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import { Box, Typography, Stack, Alert, Button } from "@mui/material"
import EventList from "../ui/events/EventList"
import { getEvents } from "../api/events"
import { useAuth } from "../context/AuthContext"
import { Event } from "../lib/types"

export default function EventsPage() {
  const { user } = useAuth()
  const [events, setEvents] = useState<Event[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
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

  const handleCreateEvent = () => {
    if (!user) {
      router.push("/auth?mode=login")
      return
    }
    router.push("/events/create")
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
          
          <EventList events={events} />
        </Box>
      </Stack>
    </Box>
  )
}