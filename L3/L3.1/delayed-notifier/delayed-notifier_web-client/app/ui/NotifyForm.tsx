"use client"

import { Button, Container, Paper, Stack, TextField } from "@mui/material"
import { useState } from "react"
import { useForm, Controller } from "react-hook-form"
import SendIcon from "@mui/icons-material/Send"
import { FieldRow } from "@/app/ui/FieldRow"

interface NotifyResponse {
  error?: string
  id?: string
  status?: string
}

interface NotifyFormValues {
  message: string
  telegram: string
  email: string
  sentAt: string
}

export default function NotifyForm() {
  const [createStatus, setCreateStatus] = useState<string>("Ожидание запроса")

  const createForm = useForm<NotifyFormValues>({
    defaultValues: {
      message: "",
      telegram: "",
      email: "",
      sentAt: "",
    },
  })

  const statusColor = (status: string) => {
    if (status.includes("Ошибка")) return "error.main"
    if (status.includes("успеш")) return "success.main"
    return "primary.main"
  }

  //Create

  const onCreateSubmit = async (formData: NotifyFormValues) => {
    setCreateStatus("Отправка запроса...")

    try {
      const payload = {
        user_id: 1,
        message: formData.message,
        channels: [
          { type: "email", value: formData.email },
          { type: "telegram", value: formData.telegram },
        ],
        sent_at: new Date(formData.sentAt).toISOString(),
      }

      const response = await fetch(
        "https://insecurely-fond-shiner.cloudpub.ru/notify",
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(payload),
        }
      )

      const data: NotifyResponse = await response.json()

      if (!response.ok) {
        throw new Error(data.error || `HTTP error: ${response.status}`)
      }

      setCreateStatus(
        `Данные успещно отправлены. ID: ${data.id ?? "неизвестен"}`
      )
    } catch (err) {
      const msg = err instanceof Error ? err.message : "Неизвестная ошибка"
      setCreateStatus(`Ошибка: ${msg}`)
    }
  }
  return (
    <Stack spacing={4}>
      <Container component={Paper} sx={{ p: 3 }}>
        <form onSubmit={createForm.handleSubmit(onCreateSubmit)}>
          <Stack spacing={2}>
            <Controller
              name="message"
              control={createForm.control}
              render={({ field }) => (
                <TextField {...field} label="Сообщение уведомления" fullWidth />
              )}
            />
            <Controller
              name="sentAt"
              control={createForm.control}
              render={({ field }) => (
                <TextField
                  {...field}
                  type="datetime-local"
                  label="Дата и время уведомления"
                  slotProps={{
                    inputLabel: {
                      shrink: true,
                    },
                  }}
                  fullWidth
                />
              )}
            />

            <Controller
              name="telegram"
              control={createForm.control}
              render={({ field }) => (
                <TextField {...field} label="Аккаунт Telegram" fullWidth />
              )}
            />

            <Controller
              name="email"
              control={createForm.control}
              render={({ field }) => (
                <TextField {...field} label="E-mail" fullWidth />
              )}
            />

            <Stack direction="row" spacing={2}>
              <Button
                type="submit"
                variant="contained"
                startIcon={<SendIcon />}
              >
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
              value={createStatus}
              statusColor={statusColor(createStatus)}
            />
          </Stack>
        </form>
      </Container>
    </Stack>
  )
}
