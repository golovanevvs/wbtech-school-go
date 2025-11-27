"use client"

import { useEffect } from "react"
import { useRouter } from "next/navigation"
import { useAuth } from "@/lib/contexts/AuthContext"
import { Box, Typography, Button, Stack } from "@mui/material"

export default function Home() {
  const { isAuthenticated, user, hasRole } = useAuth()
  const router = useRouter()

  useEffect(() => {
    // Если пользователь не авторизован, перенаправляем на страницу входа
    if (!isAuthenticated) {
      router.push("/auth")
    }
  }, [isAuthenticated, router])

  // Показываем загрузку, пока проверяется авторизация
  if (!isAuthenticated) {
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

  // Если авторизован, показываем главную страницу
  return (
    <Box sx={{ maxWidth: 800, mx: "auto", p: 3 }}>
      <Typography variant="h3" component="h1" gutterBottom sx={{ textAlign: "center", mb: 4 }}>
        Добро пожаловать в Warehouse Control!
      </Typography>

      {user && (
        <Typography variant="h6" sx={{ textAlign: "center", mb: 4 }}>
          Здравствуйте, {user.name}! Ваша роль: {user.user_role}
        </Typography>
      )}

      <Stack spacing={3} sx={{ alignItems: "center" }}>
        <Typography variant="h5" gutterBottom>
          Выберите действие:
        </Typography>

        {hasRole(["Кладовщик", "Менеджер"]) && (
          <Button 
            variant="contained" 
            size="large"
            onClick={() => router.push("/items")}
            sx={{ minWidth: 200 }}
          >
            Список товаров
          </Button>
        )}

        <Button 
          variant="contained" 
          size="large"
          onClick={() => router.push("/history")}
          sx={{ minWidth: 200 }}
        >
          История действий
        </Button>
      </Stack>
    </Box>
  )
}
