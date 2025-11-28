"use client"

import { useAuthGuard } from "@/lib/hooks/useAuthGuard"
import { useAuth } from "@/lib/contexts/AuthContext"
import { Box, Typography, Card, CardContent } from "@mui/material"

export default function ProfilePage() {
  const { user, logout } = useAuth()
  // Проверяем авторизацию для всех ролей
  useAuthGuard()

  if (!user) return null

  return (
    <Box sx={{ maxWidth: 600, mx: "auto", mt: 4 }}>
      <Typography variant="h4" gutterBottom>
        Профиль пользователя
      </Typography>

      <Card>
        <CardContent>
          <Typography variant="h6">
            Имя пользователя: {user.username}
          </Typography>
          <Typography variant="h6">Роль: {user.user_role}</Typography>
          <Typography variant="h6">ID: {user.id}</Typography>

          <button onClick={logout} style={{ marginTop: 16 }}>
            Выйти из системы
          </button>
        </CardContent>
      </Card>
    </Box>
  )
}
