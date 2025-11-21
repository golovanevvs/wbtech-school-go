"use client"

import { useState, useEffect } from "react"
import { Box, Typography, Stack, Alert } from "@mui/material"
import ProfileForm from "../ui/profile/ProfileForm"
import { getCurrentUser, updateUser } from "../api/auth"
import { User } from "../lib/types"

export default function ProfilePage() {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const userData = await getCurrentUser()
        setUser(userData)
      } catch (err) {
        setError(
          err instanceof Error ? err.message : "Failed to load user data"
        )
      } finally {
        setLoading(false)
      }
    }

    fetchUser()
  }, [])

  const handleUpdate = async (userData: Partial<User>) => {
    if (!user) return

    setSaving(true)
    setError(null)

    try {
      const updatedUser = await updateUser({ ...user, ...userData })
      setUser(updatedUser)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update user")
    } finally {
      setSaving(false)
    }
  }

  const handleSubscribeToTelegram = () => {
    // Открываем Telegram бота в новом окне
    window.open("https://t.me/your_bot_name", "_blank")
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
        <Typography variant="h6">Загрузка профиля...</Typography>
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

  if (!user) {
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
          <Alert severity="error">Пользователь не найден</Alert>
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
        <ProfileForm
          user={user}
          onUpdate={handleUpdate}
          onSubscribeToTelegram={handleSubscribeToTelegram}
          isLoading={saving}
          error={error || undefined}
        />
      </Stack>
    </Box>
  )
}
