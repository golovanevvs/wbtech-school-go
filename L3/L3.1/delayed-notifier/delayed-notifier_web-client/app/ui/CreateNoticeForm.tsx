"use client"

import { Button, Paper, Stack, TextField, Typography } from "@mui/material"
import { useState } from "react"
import { useForm, Controller } from "react-hook-form"
import SendIcon from "@mui/icons-material/Send"
import { FieldRow } from "@/app/ui/FieldRow"

interface NotifyFormValues {
  message: string
  telegram: string
  email: string
  sentAt: string
}

export default function CreateNoticeForm() {
  const [status, setStatus] = useState<string>("Ожидание запроса")

  const form = useForm<NotifyFormValues>({
    defaultValues: {
      message: "",
      telegram: "",
      email: "",
      sentAt: "",
    },
  })

  const statusColor = (s: string) => {
    if (s.includes("Ошибка")) return "error.main"
    if (s.includes("Успешно")) return "success.main"
    return "primary.main"
  }

  const onSubmit = async (data: NotifyFormValues) => {
    setStatus("Отправка запроса...")

    const apiBase = process.env.NEXT_PUBLIC_API_URL

    try {
      const payload = {
        user_id: 1,
        message: data.message,
        channels: [
          { type: "email", value: data.email },
          { type: "telegram", value: data.telegram },
        ],
        sent_at: new Date(data.sentAt).toISOString(),
      }

      const response = await fetch(`${apiBase}/notify`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      })

      const resData = await response.json()

      if (!response.ok) {
        throw new Error(resData.error || `HTTP error: ${response.status}`)
      }

      setStatus(`Данные успешно отправлены. ID: ${resData.id ?? "неизвестен"}`)
    } catch (err) {
      const msg = err instanceof Error ? err.message : "Неизвестная ошибка"
      setStatus(`Ошибка: ${msg}`)
    }
  }
  return (
    <Paper sx={{ p: 3 }}>
      <Typography
        variant="h2"
        sx={{ textAlign: "center", color: "primary.dark", mb: 3 }}
      >
        Добавление уведомления
      </Typography>
      <form onSubmit={form.handleSubmit(onSubmit)}>
        <Stack spacing={2}>
          <Controller
            name="message"
            control={form.control}
            render={({ field }) => (
              <TextField {...field} label="Сообщение уведомления" fullWidth />
            )}
          />
          <Controller
            name="sentAt"
            control={form.control}
            render={({ field }) => (
              <TextField
                {...field}
                type="datetime-local"
                label="Дата и время уведомления"
                InputLabelProps={{ shrink: true }}
                fullWidth
              />
            )}
          />
          <Controller
            name="telegram"
            control={form.control}
            render={({ field }) => (
              <TextField {...field} label="Аккаунт Telegram" fullWidth />
            )}
          />
          <Controller
            name="email"
            control={form.control}
            render={({ field }) => (
              <TextField {...field} label="E-mail" fullWidth />
            )}
          />
          <Stack direction={{ xs: "column", sm: "row" }} spacing={2}>
            <Button type="submit" variant="contained" startIcon={<SendIcon />}>
              Добавить уведомление
            </Button>

            <Button
              variant="outlined"
              onClick={() =>
                window.open("https://t.me/v_delayed_notifier_bot", "_blank")
              }
            >
              Подписаться на Telegram-бота
            </Button>
          </Stack>
          <FieldRow
            label="Результат"
            value={status}
            statusColor={statusColor(status)}
          />
        </Stack>
      </form>
    </Paper>
  )
}
