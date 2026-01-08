import {
  Card,
  CardContent,
  CardActions,
  Typography,
  Button,
  Chip,
  Alert,
} from "@mui/material"
import { Booking } from "../../lib/types"

interface BookingCardProps {
  booking: Booking
  onConfirm?: (id: number) => void
  onCancel?: (id: number) => void
  error?: string | null
}

export default function BookingCard({
  booking,
  onConfirm,
  onCancel,
  error,
}: BookingCardProps) {
  const handleConfirm = () => {
    if (onConfirm) {
      onConfirm(booking.id)
    }
  }

  const handleCancel = () => {
    if (onCancel) {
      onCancel(booking.id)
    }
  }

  const getStatusColor = () => {
    switch (booking.status) {
      case "confirmed":
        return "success"
      case "cancelled":
        return "error"
      case "pending":
        return "warning"
      default:
        return "default"
    }
  }

  return (
    <Card sx={{ minWidth: 300, maxWidth: 600, margin: "10px" }}>
      <CardContent>
        <Typography variant="h6" component="div">
          Бронирование #{booking.id}
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Событие ID: {booking.eventId}
        </Typography>
        <Chip
          label={booking.status}
          color={
            getStatusColor() as
              | "default"
              | "primary"
              | "secondary"
              | "error"
              | "info"
              | "success"
              | "warning"
          }
          sx={{ mt: 1 }}
        />
        <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
          Создано: {new Date(booking.createdAt).toLocaleString("ru-RU")}
        </Typography>
        {booking.expiresAt && (
          <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
            Истекает: {new Date(booking.expiresAt).toLocaleString("ru-RU")}
          </Typography>
        )}
        {booking.confirmedAt && (
          <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
            Подтверждено:{" "}
            {new Date(booking.confirmedAt).toLocaleString("ru-RU")}
          </Typography>
        )}
        {error && (
          <Alert severity="error" sx={{ mt: 2 }}>
            {error}
          </Alert>
        )}
      </CardContent>
      <CardActions>
        {booking.status === "pending" && (
          <>
            <Button
              size="small"
              onClick={handleConfirm}
              disabled={booking.status !== "pending"}
            >
              Подтвердить
            </Button>
            <Button size="small" color="error" onClick={handleCancel}>
              Отменить
            </Button>
          </>
        )}
      </CardActions>
    </Card>
  )
}
