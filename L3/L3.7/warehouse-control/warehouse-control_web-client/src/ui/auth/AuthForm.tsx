"use client"

import { useState } from "react"
import {
  Box,
  Paper,
  TextField,
  Button,
  Typography,
  ToggleButton,
  ToggleButtonGroup,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Alert,
  Divider,
  Chip,
} from "@mui/material"
import { useAuth } from "@/lib/contexts/AuthContext"
import { useRouter } from "next/navigation"

// Предустановленные аккаунты для быстрого тестирования
const TEST_ACCOUNTS = [
  {
    label: "Кладовщик (добавление, редактирование, удаление)",
    login: "storekeeper",
    password: "password",
    role: "Кладовщик",
  },
  {
    label: "Менеджер (просмотр списка товаров)",
    login: "manager",
    password: "password",
    role: "Менеджер",
  },
  {
    label: "Аудитор (просмотр истории изменений)",
    login: "auditor",
    password: "password",
    role: "Аудитор",
  },
]

interface AuthFormProps {
  initialMode?: "login" | "register"
}

export default function AuthForm({ initialMode = "login" }: AuthFormProps) {
  const router = useRouter()
  const { login, register } = useAuth()

  const [mode, setMode] = useState<"login" | "register">(initialMode)
  const [formData, setFormData] = useState({
    username: "",
    password: "",
    name: "",
    role: "Кладовщик",
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Обработчик выбора предустановленного аккаунта для тестирования
  const handleTestAccountSelect = async (
    account: (typeof TEST_ACCOUNTS)[0]
  ) => {
    try {
      setLoading(true)
      setError(null)

      // Автоматическое заполнение формы и выполнение входа
      setFormData({
        username: account.login,
        password: account.password,
        name: account.login,
        role: account.role,
      })

      await login(account.login, account.password)
      router.push("/")
    } catch (err) {
      console.error("Failed to login with test account:", err)
      setError(err instanceof Error ? err.message : "Ошибка входа")
    } finally {
      setLoading(false)
    }
  }

  // Обработчик изменения полей формы
  const handleInputChange = (field: string, value: string) => {
    setFormData((prev) => ({
      ...prev,
      [field]: value,
    }))
  }

  // Обработчик отправки формы
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    try {
      setLoading(true)
      setError(null)

      if (mode === "login") {
        await login(formData.username, formData.password)
        router.push("/")
      } else {
        await register(
          formData.username,
          formData.password,
          formData.name,
          formData.role
        )
      }
    } catch (err) {
      console.error("Auth error:", err)
      setError(err instanceof Error ? err.message : "Произошла ошибка")
    } finally {
      setLoading(false)
    }
  }

  // Обработчик переключения режима
  const handleModeChange = (
    _: React.MouseEvent<HTMLElement>,
    newMode: "login" | "register"
  ) => {
    if (newMode !== null) {
      setMode(newMode)
      setError(null)
    }
  }

  return (
    <Box sx={{ maxWidth: 500, mx: "auto", p: 1 }}>
      <Paper elevation={3} sx={{ p: 3 }}>
        <Typography variant="h4" sx={{ mb: 3, textAlign: "center" }}>
          Warehouse Control
        </Typography>

        {/* Раздел для быстрого тестирования */}
        <Box sx={{ mb: 1 }}>
          <Typography
            variant="h6"
            sx={{ mb: 2, textAlign: "center", color: "primary.main" }}
          >
            Быстрое тестирование
          </Typography>
          <Typography
            variant="body2"
            sx={{ mb: 2, textAlign: "center", color: "text.secondary" }}
          >
            Выберите предустановленный аккаунт для быстрого входа:
          </Typography>

          <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
            {TEST_ACCOUNTS.map((account, index) => (
              <Button
                key={index}
                variant="outlined"
                onClick={() => handleTestAccountSelect(account)}
                disabled={loading}
                sx={{
                  justifyContent: "flex-start",
                  textAlign: "left",
                  py: 1.5,
                }}
              >
                <Box
                  sx={{
                    display: "flex",
                    alignItems: "center",
                    gap: 1,
                    width: "100%",
                  }}
                >
                  <Typography variant="body1" sx={{ flex: 1 }}>
                    {account.label}
                  </Typography>
                  <Chip
                    label={account.role}
                    size="small"
                    color="primary"
                    variant="outlined"
                  />
                </Box>
              </Button>
            ))}
          </Box>

          <Divider sx={{ my: 3 }}>
            <Typography variant="body2" color="text.secondary">
              или используйте форму ниже
            </Typography>
          </Divider>
        </Box>

        {/* Штатная форма авторизации */}
        <form onSubmit={handleSubmit}>
          <Box sx={{ display: "flex", flexDirection: "column", gap: 3 }}>
            {/* Переключатель режима */}
            <ToggleButtonGroup
              value={mode}
              exclusive
              onChange={handleModeChange}
              sx={{ width: "100%" }}
            >
              <ToggleButton value="login" sx={{ flex: 1 }}>
                Вход
              </ToggleButton>
              <ToggleButton value="register" sx={{ flex: 1 }}>
                Регистрация
              </ToggleButton>
            </ToggleButtonGroup>

            {/* Поле имени пользователя */}
            <TextField
              fullWidth
              label="Имя пользователя"
              value={formData.username}
              onChange={(e) => handleInputChange("username", e.target.value)}
              disabled={loading}
              required
            />

            {/* Поле пароля */}
            <TextField
              fullWidth
              label="Пароль"
              type="password"
              value={formData.password}
              onChange={(e) => handleInputChange("password", e.target.value)}
              disabled={loading}
              required
            />

            {/* Дополнительные поля для регистрации */}
            {mode === "register" && (
              <>
                <TextField
                  fullWidth
                  label="Полное имя"
                  value={formData.name}
                  onChange={(e) => handleInputChange("name", e.target.value)}
                  disabled={loading}
                  required
                />

                <FormControl fullWidth>
                  <InputLabel>Роль</InputLabel>
                  <Select
                    value={formData.role}
                    onChange={(e) => handleInputChange("role", e.target.value)}
                    disabled={loading}
                    label="Роль"
                  >
                    <MenuItem value="Кладовщик">Кладовщик</MenuItem>
                    <MenuItem value="Менеджер">Менеджер</MenuItem>
                    <MenuItem value="Аудитор">Аудитор</MenuItem>
                  </Select>
                </FormControl>
              </>
            )}

            {/* Ошибка */}
            {error && <Alert severity="error">{error}</Alert>}

            {/* Кнопка отправки */}
            <Button
              type="submit"
              variant="contained"
              size="large"
              disabled={loading}
              sx={{ mt: 2 }}
            >
              {loading
                ? "Обработка..."
                : mode === "login"
                ? "Войти"
                : "Зарегистрироваться"}
            </Button>
          </Box>
        </form>

        {/* Информация о тестовых аккаунтах */}
        <Box sx={{ mt: 3, p: 2, bgcolor: "info.light", borderRadius: 1 }}>
          <Typography variant="body2" color="info.contrastText">
            <strong>Тестовые аккаунты:</strong>
            <br />
            • storekeeper / password (Кладовщик)
            <br />
            • manager / password (Менеджер)
            <br />• auditor / password (Аудитор)
          </Typography>
        </Box>
      </Paper>
    </Box>
  )
}
