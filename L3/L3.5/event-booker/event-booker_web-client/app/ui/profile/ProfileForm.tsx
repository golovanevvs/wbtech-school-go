import {
  Box,
  TextField,
  Button,
  Typography,
  Alert,
  FormControlLabel,
  Switch,
  Stack,
} from "@mui/material"
import { useState } from "react"
import { User, UpdateUserRequest } from "../../lib/types"

interface ProfileFormProps {
  user: User
  onUpdate: (userData: UpdateUserRequest) => void
  onSubscribeToTelegram: () => void
  isLoading?: boolean
  error?: string
}

export default function ProfileForm({
  user,
  onUpdate,
  onSubscribeToTelegram,
  isLoading,
  error,
}: ProfileFormProps) {
  const [name, setName] = useState(user.name)
  const [telegramUsername, setTelegramUsername] = useState(
    user.telegramUsername || ""
  )
  const [telegramNotifications, setTelegramNotifications] = useState(false)
  const [emailNotifications, setEmailNotifications] = useState(true)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onUpdate({
      name,
      telegramUsername: telegramUsername || null,
      telegramNotifications,
      emailNotifications,
    })
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

      <Stack spacing={2} sx={{ mt: 2 }}>
       <FormControlLabel
          control={
            <Switch
              checked={telegramNotifications}
              onChange={(e) => setTelegramNotifications(e.target.checked)}
            />
          }
          label="Уведомления через Telegram"
        />

        <FormControlLabel
          control={
            <Switch
              checked={emailNotifications}
              onChange={(e) => setEmailNotifications(e.target.checked)}
            />
          }
          label="Уведомления на Email"
        />
      </Stack>

      <Box
        sx={{
          display: "flex",
          gap: 2,
          mt: 2,
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
          disabled={!!user.telegramChatID}
        >
          {user.telegramChatID ? "Подписан" : "Подписаться на Telegram-бота"}
        </Button>
      </Box>

      <Typography variant="body2" color="text.secondary" sx={{ mt: 2 }}>
        <strong>Как это работает:</strong>
        <br />
        1. Если вы находитесь на устройстве, где есть Telegram, введите Telegram username и сохраните профиль.
        <br />
        2. Нажмите Подписаться на Telegram-бота.
        <br />
        3. Откроется Telegram-бот в Telegram. Нажмите /Start. В Telegram должно появиться сообщение, что Telegram-бот укспешно привязан.
      </Typography>
    </Box>
  )
}
