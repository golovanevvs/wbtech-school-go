"use client"

import { useSearchParams } from "next/navigation"
import AuthForm from "@/ui/auth/AuthForm"

export default function AuthPage() {
  const searchParams = useSearchParams()
  const mode = searchParams.get("mode") === "register" ? "register" : "login"

  const handleAuthSuccess = () => {
    // Дополнительные действия после успешной авторизации
    // можно добавить здесь, если нужно
  }

  return (
    <div style={{ 
      minHeight: "100vh", 
      display: "flex", 
      alignItems: "center", 
      justifyContent: "center",
      padding: "20px"
    }}>
      <AuthForm mode={mode} onAuthSuccess={handleAuthSuccess} />
    </div>
  )
}
