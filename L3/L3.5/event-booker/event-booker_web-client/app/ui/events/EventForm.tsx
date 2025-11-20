import { useState } from "react"
import { Box, TextField, Button, Typography, Alert } from "@mui/material"
import { Event } from "../../lib/types"

interface EventFormProps {
  onSubmit: (
    eventData: Omit<Event, "id" | "createdAt" | "updatedAt" | "availablePlaces">
  ) => void
  onCancel?: () => void
  event?: Event
}

export default function EventForm({
  onSubmit,
  onCancel,
  event,
}: EventFormProps) {
  const [title, setTitle] = useState(event?.title || "")
  const [description, setDescription] = useState(event?.description || "")
  const [date, setDate] = useState(event?.date || "")
  const [totalPlaces, setTotalPlaces] = useState(
    event?.totalPlaces?.toString() || ""
  )
  const [bookingDeadline, setBookingDeadline] = useState(
    event?.bookingDeadline?.toString() || ""
  )
  const [error, setError] = useState("")

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()

    if (!title || !date || !totalPlaces || !bookingDeadline) {
      setError("Все поля обязательны для заполнения")
      return
    }

    const totalPlacesNum = parseInt(totalPlaces)
    const bookingDeadlineNum = parseInt(bookingDeadline)

    if (isNaN(totalPlacesNum) || totalPlacesNum <= 0) {
      setError("Количество мест должно быть положительным числом")
      return
    }

    if (isNaN(bookingDeadlineNum) || bookingDeadlineNum <= 0) {
      setError("Срок бронирования должен быть положительным числом")
      return
    }

    onSubmit({
      title,
      description,
      date,
      totalPlaces: totalPlacesNum,
      bookingDeadline: bookingDeadlineNum,
    })
  }

  return (
    <Box
      component="form"
      onSubmit={handleSubmit}
      sx={{ width: "100%", maxWidth: 420, mx: "auto", p: 2 }}
    >
      <Typography variant="h5" component="h2" gutterBottom align="center">
        {event ? "Редактировать мероприятие" : "Создать мероприятие"}
      </Typography>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <TextField
        label="Название"
        fullWidth
        margin="normal"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        required
      />

      <TextField
        label="Описание"
        fullWidth
        margin="normal"
        multiline
        rows={4}
        value={description}
        onChange={(e) => setDescription(e.target.value)}
      />

      <TextField
        label="Дата и время"
        type="datetime-local"
        fullWidth
        margin="normal"
        value={date}
        onChange={(e) => setDate(e.target.value)}
        InputLabelProps={{
          shrink: true,
        }}
        required
      />

      <TextField
        label="Количество мест"
        type="number"
        fullWidth
        margin="normal"
        value={totalPlaces}
        onChange={(e) => setTotalPlaces(e.target.value)}
        required
      />

      <TextField
        label="Срок бронирования (минуты)"
        type="number"
        fullWidth
        margin="normal"
        value={bookingDeadline}
        onChange={(e) => setBookingDeadline(e.target.value)}
        required
        helperText="Время в минутах, в течение которого бронь должна быть подтверждена"
      />

      <Box sx={{ display: "flex", gap: 2, mt: 2 }}>
        <Button type="submit" variant="contained" fullWidth>
          {event ? "Сохранить" : "Создать"}
        </Button>
        {onCancel && (
          <Button variant="outlined" fullWidth onClick={onCancel}>
            Отмена
          </Button>
        )}
      </Box>
    </Box>
  )
}
