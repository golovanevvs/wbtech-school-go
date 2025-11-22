// app/auth/page.tsx
"use client"

import { useState, useEffect } from "react"
import { Box } from "@mui/material"
import AuthForm from "../ui/auth/AuthForm"

export default function AuthPage() {
  const [mode, setMode] = useState<"login" | "register">("login")

  useEffect(() => {
    const handlePopState = () => {
      const urlParams = new URLSearchParams(window.location.search)
      const modeParam = urlParams.get("mode")
      setMode(modeParam === "register" ? "register" : "login")
    }

    handlePopState()

    window.addEventListener("popstate", handlePopState)

    return () => {
      window.removeEventListener("popstate", handlePopState)
    }
  }, [])

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