"use client"

import { useAuth } from "@/lib/contexts/AuthContext"
import { Box, Typography, Button, Stack } from "@mui/material"

export default function Home() {
  const { user, hasRole, isLoading } = useAuth()

  // Показываем загрузку, пока проверяется авторизация
  if (isLoading) {
    return (
      <Box 
        sx={{ 
          display: "flex", 
          justifyContent: "center", 
          alignItems: "center", 
          minHeight: "50vh" 
        }}
      >
        <Typography>Загрузка...</Typography>
      </Box>
    )
  }

  // Если не авторизован, ничего не показываем (useAuthGuard перенаправит)
  if (!user) {
    return null
  }

  // Если авторизован, показываем главную страницу
  return (
    <Box sx={{ maxWidth: 800, mx: "auto", p: 3 }}>
      <Typography variant="h3" component="h1" gutterBottom sx={{ textAlign: "center", mb: 4 }}>
        Добро пожаловать в Warehouse Control!
      </Typography>

      <Typography variant="h6" sx={{ textAlign: "center", mb: 4 }}>
        Здравствуйте, {user.name}! Ваша роль: {user.user_role}
      </Typography>

      <Stack spacing={3} sx={{ alignItems: "center" }}>
        <Typography variant="h5" gutterBottom>
          Выберите действие:
        </Typography>

        {hasRole(["Кладовщик", "Менеджер"]) && (
          <Button 
            variant="contained" 
            size="large"
            onClick={() => window.location.href = "/items"}
            sx={{ minWidth: 200 }}
          >
            Список товаров
          </Button>
        )}

        <Button 
          variant="contained" 
          size="large"
          onClick={() => window.location.href = "/history"}
          sx={{ minWidth: 200 }}
        >
          История действий
        </Button>
      </Stack>
    </Box>
  )
}
