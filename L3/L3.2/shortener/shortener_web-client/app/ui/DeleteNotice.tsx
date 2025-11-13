"use client"

import { useState } from "react"
import {
  Paper,
  TextField,
  Stack,
  InputAdornment,
  IconButton,
  Typography,
} from "@mui/material"
import DeleteIcon from "@mui/icons-material/Delete"
import { FieldRow } from "@/app/ui/FieldRow"

export default function DeleteNotification() {
  const [deleteId, setDeleteId] = useState("")
  const [status, setStatus] = useState("Ожидание запроса")

  const statusColor = (s: string) => {
    if (s.includes("Ошибка")) return "error.main"
    if (s.includes("успеш")) return "success.main"
    return "primary.main"
  }

  const onDelete = async () => {
    if (!deleteId) {
      setStatus("Ошибка: Укажите ID уведомления")
      return
    }

    setStatus("Удаление уведомления...")

    const apiBase = process.env.NEXT_PUBLIC_API_URL

    try {
      const response = await fetch(`${apiBase}/notify/${deleteId}`, {
        method: "DELETE",
      })
      const data = await response.json()

      if (!response.ok)
        throw new Error(data.error || `HTTP error: ${response.status}`)

      setStatus(`Уведомление успешно удалено`)
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
        Удаление уведомления
      </Typography>
      <Stack spacing={2}>
        <TextField
          label="ID уведомления"
          value={deleteId}
          onChange={(e) => setDeleteId(e.target.value)}
          fullWidth
          InputProps={{
            endAdornment: (
              <InputAdornment position="end">
                <IconButton onClick={onDelete}>
                  <DeleteIcon />
                </IconButton>
              </InputAdornment>
            ),
          }}
        />
        <FieldRow
          label="Результат"
          value={status}
          statusColor={statusColor(status)}
        />
      </Stack>
    </Paper>
  )
}
