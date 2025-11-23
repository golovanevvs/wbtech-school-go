import { useState } from "react"
import { Box, Button, Typography, Alert } from "@mui/material"
import { CreateBookingRequest } from "../../lib/types"

interface BookingFormProps {
  eventId: number
  onBook: (bookingData: CreateBookingRequest) => void
  onCancel?: () => void
  eventTitle: string
}

export default function BookingForm({
  eventId,
  onBook,
  onCancel,
  eventTitle,
}: BookingFormProps) {
  const [error, setError] = useState("")

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()

    // Передаем данные в формате, который ожидает сервер (snake_case)
    onBook({
      event_id: eventId,
    })
  }

  return (
    <Box
      component="form"
      onSubmit={handleSubmit}
      sx={{ width: "100%", maxWidth: 400, mx: "auto", p: 2 }}
    >
      <Typography variant="h5" component="h2" gutterBottom align="center">
        Бронирование места
      </Typography>

      <Typography variant="body1" paragraph align="center">
        Мероприятие: <strong>{eventTitle}</strong>
      </Typography>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <Typography variant="body2" color="text.secondary" paragraph>
        Ваша бронь будет действительна в течение установленного времени. Если вы
        не подтвердите бронь в течение этого времени, она будет автоматически
        отменена.
      </Typography>

      <Box sx={{ display: "flex", gap: 2, mt: 2 }}>
        <Button type="submit" variant="contained" fullWidth>
          Забронировать
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
