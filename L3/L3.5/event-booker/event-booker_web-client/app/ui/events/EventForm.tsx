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

const formatDateForInput = (isoDate: string): string => {
  if (!isoDate) return ""

  try {
    const date = new Date(isoDate)
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, "0")
    const day = String(date.getDate()).padStart(2, "0")
    const hours = String(date.getHours()).padStart(2, "0")
    const minutes = String(date.getMinutes()).padStart(2, "0")

    return `${year}-${month}-${day}T${hours}:${minutes}`
  } catch (error) {
    console.error("Error formatting date:", error)
    return ""
  }
}

const formatDateFromInput = (inputDate: string): string => {
  if (!inputDate) return ""

  try {
    const date = new Date(inputDate)
    return date.toISOString()
  } catch (error) {
    console.error("Error parsing date:", error)
    return ""
  }
}

export default function EventForm({
  onSubmit,
  onCancel,
  event,
}: EventFormProps) {
  const [title, setTitle] = useState(event?.title || "")
  const [description, setDescription] = useState(event?.description || "")
  const [date, setDate] = useState(
    event?.date ? formatDateForInput(event.date) : ""
  )
  const [totalPlaces, setTotalPlaces] = useState(
    event?.totalPlaces?.toString() || ""
  )
  const [bookingDeadline, setBookingDeadline] = useState(
    event?.bookingDeadline?.toString() || ""
  )
  const [error, setError] = useState("")

  console.log("EventForm received event:", event)
  console.log("Event date:", event?.date)
  console.log(
    "Formatted date for input:",
    event?.date ? formatDateForInput(event.date) : ""
  )

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

    const isoDate = formatDateFromInput(date)
    console.log("Submitting date:", date, "-> ISO:", isoDate)

    const submitData = {
      title,
      description,
      date: isoDate,
      totalPlaces: totalPlacesNum,
      bookingDeadline: bookingDeadlineNum,
      ownerId: event?.ownerId || 0,
      telegramNotifications: event?.telegramNotifications,
      emailNotifications: event?.emailNotifications,
    }

    console.log("EventForm submit data:", submitData)
    console.log(
      "Submitting availablePlaces: undefined (as expected - should not be sent)"
    )

    onSubmit(submitData)
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
