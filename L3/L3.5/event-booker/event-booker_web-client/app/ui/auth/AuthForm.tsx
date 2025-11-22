"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { Box, Typography, Alert } from "@mui/material"
import Card from "../Card"
import Button from "../Button"
import Input from "../Input"
import { LoginRequest, RegisterRequest } from "../../lib/types"
import { login as loginApi, register as registerApi } from "../../api/auth"
import { useAuth } from "../../context/AuthContext"

interface AuthFormProps {
  mode: "login" | "register"
  onAuthSuccess?: () => void
}

export default function AuthForm({ mode, onAuthSuccess }: AuthFormProps) {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")
  const [name, setName] = useState("")
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)

  const router = useRouter()
  const { login } = useAuth()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    setLoading(true)

    if (mode === "register") {
      if (password !== confirmPassword) {
        setError("Пароли не совпадают")
        setLoading(false)
        return
      }
    }

    try {
      let accessToken: string
      let refreshToken: string

      if (mode === "login") {
        const credentials: LoginRequest = { email, password }
        const response = await loginApi(credentials)
        accessToken = response.token

        if (!response.refreshToken) {
          throw new Error("No refresh token received from server")
        }
        refreshToken = response.refreshToken
      } else {
        const credentials: RegisterRequest = { email, password, name }
        const response = await registerApi(credentials)
        accessToken = response.token

        if (!response.refreshToken) {
          throw new Error("No refresh token received from server")
        }
        refreshToken = response.refreshToken
      }

      await login(accessToken, refreshToken)

      if (onAuthSuccess) {
        onAuthSuccess()
      }

      router.push("/events")
      router.refresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Authentication failed")
    } finally {
      setLoading(false)
    }
  }

  return (
    <Card>
      <Typography variant="h5" component="h1" gutterBottom>
        {mode === "login" ? "Вход" : "Регистрация"}
      </Typography>

      <Box component="form" onSubmit={handleSubmit} sx={{ width: "100%" }}>
        {mode === "register" && (
          <Input
            label="Имя"
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
          />
        )}

        <Input
          label="Email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />

        <Input
          label="Пароль"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />

        {mode === "register" && (
          <Input
            label="Подтверждение пароля"
            type="password"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
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
        )}

        {error && (
          <Alert severity="error" sx={{ mt: 2, width: "100%" }}>
            {error}
          </Alert>
        )}

        <Button type="submit" disabled={loading} sx={{ mt: 2 }}>
          {loading
            ? "Загрузка..."
            : mode === "login"
            ? "Войти"
            : "Зарегистрироваться"}
        </Button>
      </Box>

      <Box sx={{ mt: 2, textAlign: "center" }}>
        <Typography variant="body2" color="text.secondary">
          {mode === "login" ? "Нет аккаунта? " : "Уже есть аккаунт? "}
          <a
            href={mode === "login" ? "/auth?mode=register" : "/auth?mode=login"}
          >
            {mode === "login" ? "Зарегистрироваться" : "Войти"}
          </a>
        </Typography>
      </Box>
    </Card>
  )
}
