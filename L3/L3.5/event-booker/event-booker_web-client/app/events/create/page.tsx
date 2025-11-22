"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import { Box, Stack, Alert, Paper } from "@mui/material"
import EventForm from "../../ui/events/EventForm"
import { createEvent } from "../../api/events"
import { useAuth } from "../../context/AuthContext"
import { Event } from "../../lib/types"

export default function CreateEventPage() {
  const { user, loading } = useAuth()
  const [error, setError] = useState<string | null>(null)
  const router = useRouter()

  useEffect(() => {
    // Если не загружается контекст аутентификации, ждем
    if (loading) return

    // Если пользователь не аутентифицирован, перенаправляем на /auth
    if (!user) {
      router.push("/auth")
      return
    }
  }, [user, loading, router])

  const handleSubmit = async (
    eventData: Omit<Event, "id" | "createdAt" | "updatedAt" | "availablePlaces">
  ) => {
    try {
      setError(null)
      await createEvent(eventData)
      router.push("/events")
      router.refresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create event")
    }
  }

  const handleCancel = () => {
    router.push("/events")
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
        <Box>Загрузка...</Box>
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
        {error && <Alert severity="error">{error}</Alert>}
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
          <EventForm onSubmit={handleSubmit} onCancel={handleCancel} />
        </Paper>
      </Stack>
    </Box>
  )
}
