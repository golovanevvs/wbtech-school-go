"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import { Box, Typography, Stack, Alert, Paper } from "@mui/material"
import ProfileForm from "../ui/profile/ProfileForm"
import { getCurrentUser, updateUser, deleteUser } from "../api/auth"
import { useAuth } from "../context/AuthContext"
import { User, UpdateUserRequest } from "../lib/types"

export default function ProfilePage() {
  const { user: authUser, loading: authLoading } = useAuth()
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const router = useRouter()

  useEffect(() => {
    if (authLoading) return

    if (!authUser) {
      router.push("/auth")
      return
    }

    const fetchUser = async () => {
      try {
        const userData = authUser || (await getCurrentUser())
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
  }, [authUser, authLoading, router])

  const handleUpdate = async (userData: UpdateUserRequest, shouldLaunchTelegram?: boolean) => {
    if (!user) return

    setSaving(true)
    setError(null)

    try {
      const updatedUser = await updateUser(userData)
      setUser(updatedUser)
      
      // Если нужно запустить Telegram и у пользователя есть username
      if (shouldLaunchTelegram && updatedUser.telegramUsername) {
        // Небольшая задержка, чтобы дать время пользователю увидеть обновленные данные
        setTimeout(() => {
          window.open("tg://resolve?domain=v_delayed_notifier_bot", "_blank")
        }, 500)
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update user")
    } finally {
      setSaving(false)
    }
  }

  const handleDeleteProfile = async () => {
    if (!user) return

    setSaving(true)
    setError(null)

    try {
      // Вызываем API для удаления профиля
      await deleteUser()
      
      // Перенаправляем на страницу авторизации
      router.push("/auth")
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete profile")
    } finally {
      setSaving(false)
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
          <ProfileForm
            key={user.id}
            user={user}
            onUpdate={handleUpdate}
            onDeleteProfile={handleDeleteProfile}
            isLoading={saving}
            error={error || undefined}
          />
        </Paper>
      </Stack>
    </Box>
  )
}
