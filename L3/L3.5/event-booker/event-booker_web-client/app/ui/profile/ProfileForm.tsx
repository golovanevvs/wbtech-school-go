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
import { useState } from "react"
import { User, UpdateUserRequest } from "../../lib/types"

interface ProfileFormProps {
  user: User
  onUpdate: (userData: UpdateUserRequest, shouldLaunchTelegram?: boolean) => void
  onDeleteProfile: () => void
  isLoading?: boolean
  error?: string
}

export default function ProfileForm({
  user,
  onUpdate,
  onDeleteProfile,
  isLoading = false,
  error,
}: ProfileFormProps) {
  // Инициализируем состояние один раз из user prop
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

  // Переключатель активен только если есть chatID
  const telegramEnabled = !!user.telegramChatID

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    onUpdate({
      name: name.trim(),
      telegramUsername: telegramUsername.trim() || null,
      telegramNotifications,
      emailNotifications,
    })
  }

  const handleSubscribeToTelegram = async () => {
    // Проверяем, что username не пустой
    const username = telegramUsername.trim()
    if (!username) {
      alert("Пожалуйста, введите Telegram username перед подпиской на бота")
      return
    }

    // Обновляем данные пользователя с текущим Telegram username
    // и сбрасываем chatID (передаем 0 для сброса на сервере)
    onUpdate({
      name: name.trim(),
      telegramUsername: username,
      telegramNotifications,
      emailNotifications,
      resetTelegramChatID: true, // специальный флаг для сброса chatID
    }, true) // передаем shouldLaunchTelegram = true
  }

  const handleCheckSubscription = async () => {
    // Просто обновляем данные пользователя без изменений
    // чтобы получить актуальное состояние от сервера
    onUpdate({
      name: name.trim(),
      telegramUsername: telegramUsername.trim() || null,
      telegramNotifications,
      emailNotifications,
    }, false) // передаем shouldLaunchTelegram = false
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

      {/* 1. Имя */}
      <TextField
        label="Имя"
        fullWidth
        margin="normal"
        value={name}
        onChange={(e) => setName(e.target.value)}
        required
      />

      <Divider sx={{ my: 3 }} />

      {/* 2. Настройки уведомлений */}
      <Typography variant="h6" gutterBottom>
        Настройки уведомлений
      </Typography>

      <Stack spacing={3} sx={{ mt: 2 }}>
        {/* 2.1 Для e-mail */}
        <Box>
          <Typography variant="subtitle1" gutterBottom>
            Электронная почта
          </Typography>
          <FormControlLabel
            control={
              <Switch
                checked={emailNotifications}
                onChange={(e) => setEmailNotifications(e.target.checked)}
              />
            }
            label="Уведомления по электронной почте"
          />
        </Box>

        {/* 2.2 Для Telegram */}
        <Box>
          <Typography variant="subtitle1" gutterBottom>
            Telegram
          </Typography>
          
          {/* 1. Переключатель (не активен пока нет chatID) */}
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
                : "Уведомления через Telegram (требуется подключение)"
            }
          />

          {/* 2. Инструкция */}
          <Typography variant="body2" color="text.secondary" sx={{ mt: 1, mb: 2 }}>
            <strong>Как подключить уведомления через Telegram:</strong>
            <br />
            1. Введите ваш Telegram username (без @) в поле ниже
            <br />
            2. Нажмите кнопку Подписаться на Telegram-бота
            <br />
            3. Откроется Telegram-бот. Нажмите /Start для привязки аккаунта
            <br />
            4. Проверьте состояние подписки кнопкой ниже
          </Typography>

          {/* 3. Поле для ввода Telegram username */}
          <TextField
            label="Telegram username"
            fullWidth
            margin="normal"
            value={telegramUsername}
            onChange={(e) => setTelegramUsername(e.target.value)}
            placeholder="без @, например: ivan_ivanov"
            helperText="Введите ваш Telegram username (без @)"
          />

          {/* 4. Кнопки */}
          <Stack direction="row" spacing={2} sx={{ mt: 2 }}>
            <Button
              variant="outlined"
              fullWidth
              onClick={handleSubscribeToTelegram}
              disabled={isLoading}
            >
              {isLoading ? "Сохранение..." : "Подписаться на Telegram-бота"}
            </Button>
            
            <Button
              variant="outlined"
              fullWidth
              onClick={handleCheckSubscription}
              disabled={isLoading}
            >
              {isLoading ? "Проверка..." : "Проверить состояние подписки"}
            </Button>
          </Stack>
        </Box>
      </Stack>

      <Divider sx={{ my: 3 }} />

      {/* 5. Кнопки управления профилем */}
      <Stack spacing={2}>
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
          color="error"
          fullWidth
          onClick={handleDeleteProfile}
          disabled={isLoading}
        >
          Удалить профиль
        </Button>
      </Stack>
    </Box>
  )
}
