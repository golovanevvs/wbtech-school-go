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

    // Установка начального значения
    handlePopState()

    // Подписка на изменение URL
    window.addEventListener("popstate", handlePopState)

    return () => {
      window.removeEventListener("popstate", handlePopState)
    }
  }, [])

  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: "center",
        px: { xs: 0, sm: 2 },
        py: 2,
        bgcolor: "background.default",
        maxWidth:330,
      }}
    >
      <AuthForm mode={mode} />
    </Box>
  )
}
