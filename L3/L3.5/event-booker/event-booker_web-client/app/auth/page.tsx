"use client"

import { Suspense } from "react"
import { useSearchParams } from "next/navigation"
import { Box } from "@mui/material"
import AuthForm from "../ui/auth/AuthForm"

function AuthFormWrapper() {
  const searchParams = useSearchParams()

  const mode = searchParams.get("mode") === "register" ? "register" : "login"

  const handleAuthSuccess = () => {
    console.log("Auth successful!")
  }

  return <AuthForm mode={mode} onAuthSuccess={handleAuthSuccess} />
}

export default function AuthPage() {
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
      <Suspense fallback={<div>Загрузка...</div>}>
        <AuthFormWrapper />
      </Suspense>
    </Box>
  )
}
