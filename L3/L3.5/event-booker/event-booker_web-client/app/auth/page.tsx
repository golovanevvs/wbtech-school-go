// app/auth/page.tsx
"use client"

import { useMemo } from "react"
import { useSearchParams } from "next/navigation"
import { Box } from "@mui/material"
import AuthForm from "../ui/auth/AuthForm"

export default function AuthPage() {
  const searchParams = useSearchParams()
  
  const mode = useMemo(() => {
    const modeParam = searchParams.get("mode")
    return modeParam === "register" ? "register" : "login"
  }, [searchParams])

  // 👇 Колбэк для успешной авторизации
  const handleAuthSuccess = () => {
    // Header автоматически обновится через контекст AuthContext
    console.log("Auth successful!")
  }

  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: "center",
        px: { xs: 2, sm: 2 },
        py: 2,
        bgcolor: "background.default",
        maxWidth: 500,
      }}
    >
      <AuthForm 
        mode={mode} 
        onAuthSuccess={handleAuthSuccess} // 👈 Теперь без ошибки
      />
    </Box>
  )
}