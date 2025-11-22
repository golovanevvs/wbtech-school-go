import {
  Box,
  TextField,
  Button,
  Typography,
  Alert,
  FormControlLabel,
  Switch,
  Stack,
  Divider,
} from "@mui/material"
import { useState, useEffect } from "react"
import { User, UpdateUserRequest } from "../../lib/types"

interface ProfileFormProps {
  user: User
  onUpdate: (userData: UpdateUserRequest) => void
  onSubscribeToTelegram: () => void
  onDeleteProfile: () => void
  isLoading?: boolean
  error?: string
}

export default function ProfileForm({
  user,
  onUpdate,
  onSubscribeToTelegram,
  onDeleteProfile,
  isLoading,
  error,
}: ProfileFormProps) {
  const [name, setName] = useState(user.name)
  const [telegramUsername, setTelegramUsername] = useState(
    user.telegramUsername || ""
  )
  const [telegramNotifications, setTelegramNotifications] = useState(
    user.telegramNotifications || false
  )
  const [emailNotifications, setEmailNotifications] = useState(
    user.emailNotifications || false
  )

  // Переключатель Telegram активен только если есть аккаунт
  const telegramEnabled = !!user.telegramChatID
  // Кнопка подписки активна только если включены telegram уведомления и есть аккаунт
  const canSubscribeToTelegram = telegramEnabled && telegramNotifications

  // Синхронизируем локальное состояние с данными пользователя
  useEffect(() => {
    setName(user.name)
    setTelegramUsername(user.telegramUsername || "")
    setTelegramNotifications(user.telegramNotifications || false)
    setEmailNotifications(user.emailNotifications || false)
  }, [user])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onUpdate({
      name,
      telegramUsername: telegramUsername || null,
      telegramNotifications,
      emailNotifications,
    })
  }

  const handleDeleteProfile = () => {
    if (
      window.confirm(
        "Вы уверены, что хотите удалить профиль? Это действие нельзя отменить. Все ваши мероприятия будут также удалены."
      )
    ) {
      onDeleteProfile()
    }
  }

  return (
    <Box
      component="form"
      onSubmit={handleSubmit}
      sx={{ width: "100%", maxWidth: 500, mx: "auto", p: 2 }}
    >
      <Typography variant="h5" component="h2" gutterBottom align="center">
        Профиль пользователя
      </Typography>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <TextField
        label="Имя"
        fullWidth
        margin="normal"
        value={name}
        onChange={(e) => setName(e.target.value)}
        required
      />

      <TextField
        label="Telegram username"
        fullWidth
        margin="normal"
        value={telegramUsername}
        onChange={(e) => setTelegramUsername(e.target.value)}
        placeholder="без @, например: ivan_ivanov"
        helperText="Введите ваш Telegram username (без @)"
      />

      <Divider sx={{ my: 3 }} />

      <Typography variant="h6" gutterBottom>
        Настройки уведомлений
      </Typography>

      <Stack spacing={2} sx={{ mt: 2 }}>
        <FormControlLabel
          control={
            <Switch
              checked={emailNotifications}
              onChange={(e) => setEmailNotifications(e.target.checked)}
            />
          }
          label="Уведомления по электронной почте"
        />

        <FormControlLabel
          control={
            <Switch
              checked={telegramNotifications}
              onChange={(e) => setTelegramNotifications(e.target.checked)}
              disabled={!telegramEnabled}
            />
          }
          label={
            telegramEnabled 
              ? "Уведомления через Telegram" 
              : "Уведомления через Telegram (недоступно - подключите аккаунт)"
          }
        />
      </Stack>

      <Box
        sx={{
          display: "flex",
          gap: 2,
          mt: 3,
          flexDirection: { xs: "column", sm: "row" },
        }}
      >
        <Button
          type="submit"
          variant="contained"
          fullWidth
          disabled={isLoading}
        >
          {isLoading ? "Сохранение..." : "Сохранить"}
        </Button>

        <Button
          variant="outlined"
          fullWidth
          onClick={onSubscribeToTelegram}
          disabled={!canSubscribeToTelegram}
        >
          {user.telegramChatID ? "Подписан на Telegram-бота" : "Подписаться на Telegram-бота"}
        </Button>
      </Box>

      <Divider sx={{ my: 3 }} />

      <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
        <strong>Как это работает:</strong>
        <br />
        1. Если вы находитесь на устройстве, где есть Telegram, введите Telegram username и сохраните профиль.
        <br />
        2. Включите переключатель "Уведомления через Telegram".
        <br />
        3. Нажмите "Подписаться на Telegram-бота".
        <br />
        4. Откроется Telegram-бот в Telegram. Нажмите /Start. В Telegram должно появиться сообщение, что Telegram-бот успешно привязан.
      </Typography>

      <Button
        variant="outlined"
        color="error"
        fullWidth
        onClick={handleDeleteProfile}
        disabled={isLoading}
        sx={{ mt: 2 }}
      >
        Удалить профиль
      </Button>
    </Box>
  )
}
