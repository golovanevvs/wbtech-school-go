"use client"

import { useState, useEffect } from "react"
import { useAuthGuard } from "@/lib/hooks/useAuthGuard"
import { useAuth } from "@/lib/contexts/AuthContext"
import { authAPI } from "@/lib/api/auth"
import {
  Box,
  Typography,
  Card,
  CardContent,
  Button,
  Alert,
  CircularProgress,
} from "@mui/material"

export default function ProfilePage() {
  const { user, deleteUser } = useAuth()
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [userData, setUserData] = useState(user)

  // Проверка авторизации для всех ролей
  useAuthGuard()

  // Загрузка данных пользователя с сервера
  const loadUserData = async () => {
    try {
      setLoading(true)
      setError(null)
      const currentUser = await authAPI.getCurrentUser()
      setUserData(currentUser)
    } catch (err) {
      console.error("Failed to load user data:", err)
      setError(
        err instanceof Error
          ? err.message
          : "Не удалось загрузить данные пользователя"
      )
    } finally {
      setLoading(false)
    }
  }

  // Загрузка данных при монтировании компонента
  useEffect(() => {
    loadUserData()
  }, [])

  // Обработчик удаления профиля
  const handleDeleteProfile = async () => {
    if (
      !window.confirm(
        "Вы уверены, что хотите удалить свой профиль? Это действие нельзя отменить."
      )
    ) {
      return
    }

    try {
      setLoading(true)
      setError(null)

      const success = await deleteUser()
      if (success) {
      }
    } catch (err) {
      console.error("Failed to delete profile:", err)
      setError(
        err instanceof Error ? err.message : "Не удалось удалить профиль"
      )
    } finally {
      setLoading(false)
    }
  }

  if (!userData) return null

  return (
    <Box sx={{ maxWidth: 600, mx: "auto", mt: 4 }}>
      <Typography variant="h4" gutterBottom>
        Профиль пользователя
      </Typography>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <Card>
        <CardContent>
          <Box sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
            <Typography variant="h6">
              Имя пользователя: {userData.username}
            </Typography>
            <Typography variant="h6">
              Отображаемое имя: {userData.name}
            </Typography>
            <Typography variant="h6">Роль: {userData.user_role}</Typography>
            <Typography variant="h6">ID: {userData.id}</Typography>

            <Box sx={{ display: "flex", gap: 2, mt: 3 }}>
              <Button
                variant="outlined"
                onClick={loadUserData}
                disabled={loading}
                startIcon={loading ? <CircularProgress size={20} /> : null}
              >
                Обновить данные
              </Button>

              <Button
                variant="contained"
                color="error"
                onClick={handleDeleteProfile}
                disabled={loading}
              >
                Удалить профиль
              </Button>
            </Box>
          </Box>
        </CardContent>
      </Card>
    </Box>
  )
}
