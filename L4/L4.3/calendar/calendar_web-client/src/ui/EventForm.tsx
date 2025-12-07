"use client"

import { useState } from "react"
import {
  Box,
  TextField,
  Button,
  FormControlLabel,
  Switch,
  Typography,
  Paper,
} from "@mui/material"
import { CreateEventRequest } from "@/lib/types/calendar"
import { calendarApi } from "@/lib/api/calendar"

interface EventFormProps {
  onSuccess?: () => void
  onCancel?: () => void
}

export default function EventForm({ onSuccess, onCancel }: EventFormProps) {
  const [formData, setFormData] = useState<CreateEventRequest>({
    title: "",
    description: "",
    start: "",
    end: "",
    allDay: false,
    reminder: false,
    reminderTime: "",
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleChange = (field: keyof CreateEventRequest) => (
    event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const value = event.target.type === "checkbox" 
      ? (event.target as HTMLInputElement).checked
      : event.target.value
    setFormData(prev => ({ ...prev, [field]: value }))
  }

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault()
    
    if (!formData.title.trim()) {
      setError("Название события обязательно")
      return
    }

    if (!formData.start) {
      setError("Дата и время начала обязательны")
      return
    }

    try {
      setLoading(true)
      setError(null)
      
      // Добавляем timezone к времени для совместимости с Go
      const eventData = {
        ...formData,
        start: formData.start.includes('T') ? formData.start + ':00Z' : formData.start,
        end: formData.end ? (formData.end.includes('T') ? formData.end + ':00Z' : formData.end) : undefined,
        reminderTime: formData.reminderTime ? (formData.reminderTime.includes('T') ? formData.reminderTime + ':00Z' : formData.reminderTime) : undefined,
      }
      
      await calendarApi.createEvent(eventData)
      
      if (onSuccess) {
        onSuccess()
      }
    } catch (err) {
      console.error("Failed to create event:", err)
      setError(err instanceof Error ? err.message : "Ошибка при создании события")
    } finally {
      setLoading(false)
    }
  }

  const handleCancel = () => {
    if (onCancel) {
      onCancel()
    }
  }

  // Установка значений по умолчанию
  const setDefaultDateTime = () => {
    const now = new Date()
    const defaultStart = new Date(now.getTime() + 60 * 60 * 1000) // +1 час
    const defaultEnd = new Date(now.getTime() + 2 * 60 * 60 * 1000) // +2 часа

    setFormData(prev => ({
      ...prev,
      start: defaultStart.toISOString().slice(0, 16),
      end: defaultEnd.toISOString().slice(0, 16),
    }))
  }

  return (
    <Paper
      component="form"
      onSubmit={handleSubmit}
      sx={{
        p: 3,
        maxWidth: "600px",
        mx: "auto",
      }}
    >
      <Typography variant="h5" component="h2" gutterBottom>
        Создание события
      </Typography>

      {error && (
        <Typography color="error" sx={{ mb: 2 }}>
          {error}
        </Typography>
      )}

      <Box sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
        <TextField
          required
          label="Название события"
          value={formData.title}
          onChange={handleChange("title")}
          fullWidth
        />

        <TextField
          label="Описание"
          value={formData.description}
          onChange={handleChange("description")}
          fullWidth
          multiline
          rows={3}
        />

        <TextField
          required
          label="Дата и время начала"
          type="datetime-local"
          value={formData.start}
          onChange={handleChange("start")}
          fullWidth
          InputLabelProps={{ shrink: true }}
        />

        <TextField
          label="Дата и время окончания"
          type="datetime-local"
          value={formData.end}
          onChange={handleChange("end")}
          fullWidth
          InputLabelProps={{ shrink: true }}
        />

        <FormControlLabel
          control={
            <Switch
              checked={formData.allDay}
              onChange={handleChange("allDay")}
            />
          }
          label="Событие на весь день"
        />

        <FormControlLabel
          control={
            <Switch
              checked={formData.reminder}
              onChange={handleChange("reminder")}
            />
          }
          label="Напоминание"
        />

        {formData.reminder && (
          <TextField
            label="Время напоминания (до начала события)"
            type="datetime-local"
            value={formData.reminderTime}
            onChange={handleChange("reminderTime")}
            fullWidth
            InputLabelProps={{ shrink: true }}
          />
        )}

        <Box sx={{ display: "flex", gap: 2, mt: 2 }}>
          <Button
            type="submit"
            variant="contained"
            disabled={loading}
            fullWidth
          >
            {loading ? "Создание..." : "Создать"}
          </Button>
          <Button
            type="button"
            variant="outlined"
            onClick={handleCancel}
            disabled={loading}
            fullWidth
          >
            Отмена
          </Button>
        </Box>

        <Button
          type="button"
          variant="text"
          onClick={setDefaultDateTime}
          disabled={loading}
          sx={{ mt: 1 }}
        >
          Установить время по умолчанию
        </Button>
      </Box>
    </Paper>
  )
}