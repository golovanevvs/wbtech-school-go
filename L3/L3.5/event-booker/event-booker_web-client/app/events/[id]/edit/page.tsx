"use client"

import { useState, useEffect } from "react"
import { useRouter, useParams } from "next/navigation"
import { Box, Stack, Alert, Paper } from "@mui/material"
import EventForm from "../../../ui/events/EventForm"
import { getEventById, updateEvent } from "../../../api/events"
import { useAuth } from "../../../context/AuthContext"
import { Event } from "../../../lib/types"

export default function EditEventPage() {
  const { user, loading } = useAuth()
  const [event, setEvent] = useState<Event | null>(null)
  const [loadingEvent, setLoadingEvent] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const router = useRouter()
  const params = useParams()
  const eventId = parseInt(params.id as string)

  useEffect(() => {
    // Если не загружается контекст аутентификации, ждем
    if (loading) return

    // Если пользователь не аутентифицирован, перенаправляем на /auth
    if (!user) {
      router.push("/auth")
      return
    }

    // Загружаем данные мероприятия
    const fetchEvent = async () => {
      try {
        setLoadingEvent(true)
        const eventData = await getEventById(eventId)
        
        // Проверяем, является ли пользователь владельцем мероприятия
        if (eventData.ownerId !== user.id) {
          setError("Вы не можете редактировать это мероприятие")
          return
        }
        
        setEvent(eventData)
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load event")
      } finally {
        setLoadingEvent(false)
      }
    }

    fetchEvent()
  }, [user, loading, router, eventId])

  const handleSubmit = async (
    eventData: Omit<Event, "id" | "createdAt" | "updatedAt" | "availablePlaces">
  ) => {
    try {
      setError(null)
      await updateEvent(eventId, eventData)
      router.push("/events")
      router.refresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update event")
    }
  }

  const handleCancel = () => {
    router.push("/events")
  }

  if (loading || loadingEvent) {
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
        <Box>Загрузка...</Box>
      </Box>
    )
  }

  if (error) {
    return (
      <Box
        sx={{
          width: "100%",
          px: { xs: 2, sm: 2 },
          py: 1,
          bgcolor: "background.default",
          maxWidth: "100vw",
          mx: "auto",
        }}
      >
        <Stack spacing={2} alignItems="center">
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
          px: { xs: 2, sm: 2 },
          py: 1,
          bgcolor: "background.default",
          maxWidth: "100vw",
          mx: "auto",
        }}
      >
        <Stack spacing={2} alignItems="center">
          <Alert severity="warning">Мероприятие не найдено</Alert>
        </Stack>
      </Box>
    )
  }

  return (
    <Box
      sx={{
        width: "100%",
        px: { xs: 2, sm: 2 },
        py: 1,
        bgcolor: "background.default",
        maxWidth: "100vw",
        mx: "auto",
      }}
    >
      <Stack spacing={2} alignItems="center">
        <Paper
          sx={{
            p: { xs: 0, sm: 2 },
            width: "100%",
            maxWidth: 500,
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
          }}
        >
          <EventForm 
            onSubmit={handleSubmit} 
            onCancel={handleCancel} 
            event={event}
          />
        </Paper>
      </Stack>
    </Box>
  )
}