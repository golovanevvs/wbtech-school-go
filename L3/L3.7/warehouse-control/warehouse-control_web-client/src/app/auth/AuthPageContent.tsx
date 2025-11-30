"use client"

import { useSearchParams } from "next/navigation"
import AuthForm from "@/ui/auth/AuthForm"
import { Box } from "@mui/material"

export default function AuthPageContent() {
  const searchParams = useSearchParams()
  const mode = searchParams.get("mode") === "register" ? "register" : "login"

  console.log("AuthPage render - mode:", mode)

  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: "center",
        px: { xs: 2, sm: 2 },
      }}
    >
      <AuthForm initialMode={mode} />
    </Box>
  )
}
