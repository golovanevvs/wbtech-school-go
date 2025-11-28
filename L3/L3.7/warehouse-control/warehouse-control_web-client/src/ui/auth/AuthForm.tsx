"use client"

import { useState } from "react"
import { useAuth } from "@/lib/contexts/AuthContext"
import { UserRole } from "@/lib/types/auth"
import { useRouter } from "next/navigation"
import { Typography, Box, Alert, Link as MuiLink } from "@mui/material"
import VCard from "../VCard"
import Button from "../Button"
import Input from "../Input"

interface AuthFormProps {
  mode: "login" | "register"
  onAuthSuccess?: () => void
}

export default function AuthForm({ mode, onAuthSuccess }: AuthFormProps) {
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")
  const [name, setName] = useState("")
  const [role, setRole] = useState<UserRole>("Кладовщик")
  const [error, setError] = useState("")
  const [isLoading, setIsLoading] = useState(false)

  const router = useRouter()
  const { login, register, error: authError, clearError } = useAuth()

  // Логирование для отладки
  console.log("AuthForm render - authError:", authError)
  console.log("AuthForm render - mode:", mode)
  console.log("AuthForm render - isLoading:", isLoading)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    clearError()
    setIsLoading(true)

    try {
      let success = false

      if (mode === "login") {
        success = await login(username, password)
      } else {
        if (password !== confirmPassword) {
          setError("Пароли не совпадают")
          setIsLoading(false)
          return
        }
        success = await register(username, password, name, role)
      }

      if (success) {
        if (onAuthSuccess) {
          onAuthSuccess()
        }
        router.push("/")
      } else {
        // Если нет локальной ошибки, но есть ошибка из AuthContext, показываем ее
        if (!authError) {
          setError(
            mode === "login"
              ? "Неверный логин или пароль"
              : "Не удалось зарегистрироваться. Проверьте введенные данные."
          )
        }
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Произошла ошибка")
    } finally {
      setIsLoading(false)
    }
  }

  const handleModeChange = () => {
    const newMode = mode === "login" ? "register" : "login"
    router.push(`/auth?mode=${newMode}`)
  }

  const handleInputChange = () => {
    clearError()
    setError("")
  }

  return (
    <VCard>
      <Typography variant="h4" component="h1" gutterBottom>
        {mode === "login" ? "Вход в систему" : "Регистрация"}
      </Typography>

      <Box component="form" onSubmit={handleSubmit} sx={{ width: "100%" }}>
        <Input
          label="Логин"
          type="text"
          value={username}
          onChange={(e) => {
            setUsername(e.target.value)
            handleInputChange()
          }}
          required
        />

        <Input
          label="Пароль"
          type="password"
          value={password}
          onChange={(e) => {
            setPassword(e.target.value)
            handleInputChange()
          }}
          required
        />

        {mode === "register" && (
          <>
            <Input
              label="Подтверждение пароля"
              type="password"
              value={confirmPassword}
              onChange={(e) => {
                setConfirmPassword(e.target.value)
                handleInputChange()
              }}
              required
              error={
                password !== "" &&
                confirmPassword !== "" &&
                password !== confirmPassword
              }
              helperText={
                password !== "" &&
                confirmPassword !== "" &&
                password !== confirmPassword
                  ? "Пароли не совпадают"
                  : ""
              }
            />

            <Input
              label="Имя"
              type="text"
              value={name}
              onChange={(e) => {
                setName(e.target.value)
                handleInputChange()
              }}
              required
            />

            <Box sx={{ mt: 2, mb: 1 }}>
              <Typography variant="body2" color="text.secondary" gutterBottom>
                Роль
              </Typography>
              <Box sx={{ display: "flex", gap: 1, flexWrap: "wrap" }}>
                {(["Кладовщик", "Менеджер", "Аудитор"] as UserRole[]).map(
                  (r) => (
                    <Button
                      key={r}
                      type="button"
                      variant={role === r ? "contained" : "outlined"}
                      onClick={() => setRole(r)}
                      sx={{
                        flex: 1,
                        minWidth: "100px",
                        fontSize: "0.8rem",
                      }}
                    >
                      {r}
                    </Button>
                  )
                )}
              </Box>
            </Box>
          </>
        )}

        {(error || authError) && (
          <Alert severity="error" sx={{ mt: 2, width: "100%" }}>
            {error || authError}
          </Alert>
        )}

        {/* Логирование условия для Alert */}
        {(() => {
          const showAlert = error || authError
          console.log(
            "Alert condition check - error:",
            error,
            "authError:",
            authError,
            "showAlert:",
            showAlert
          )
          return null
        })()}

        <Button type="submit" disabled={isLoading} sx={{ mt: 3 }}>
          {isLoading
            ? "Загрузка..."
            : mode === "login"
            ? "Войти"
            : "Зарегистрироваться"}
        </Button>
      </Box>

      <Box sx={{ mt: 3, textAlign: "center" }}>
        <Typography variant="body2" color="text.secondary">
          {mode === "login" ? "Нет аккаунта? " : "Уже есть аккаунт? "}
          <MuiLink
            component="button"
            type="button"
            onClick={handleModeChange}
            sx={{
              color: "primary.main",
              textDecoration: "underline",
              cursor: "pointer",
              border: "none",
              background: "none",
              font: "inherit",
              padding: 0,
            }}
          >
            {mode === "login" ? "Зарегистрироваться" : "Войти"}
          </MuiLink>
        </Typography>
      </Box>
    </VCard>
  )
}
