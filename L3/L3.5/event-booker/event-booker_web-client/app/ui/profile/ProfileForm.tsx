import { useState } from "react"
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
import { User } from "../../lib/types"

interface ProfileFormProps {
  user: User
  onUpdate: (userData: Partial<User>) => void
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
  const [telegramEnabled, setTelegramEnabled] = useState(!!user.telegramChatID)
  const [emailNotifications, setEmailNotifications] = useState(true)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onUpdate({
      name,
      telegramChatID: user.telegramChatID ? user.telegramChatID : null,
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

      <Stack spacing={2} sx={{ mt: 2 }}>
        <FormControlLabel
          control={
            <Switch
              checked={telegramEnabled}
              onChange={(e) => setTelegramEnabled(e.target.checked)}
              disabled={!user.telegramChatID}
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
    </Box>
  )
}
