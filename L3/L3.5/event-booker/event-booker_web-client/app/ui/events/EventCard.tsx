import {
  Card,
  CardContent,
  CardActions,
  Typography,
  Button,
  Chip,
  Box,
  IconButton,
  Alert,
} from "@mui/material"
import EditIcon from "@mui/icons-material/Edit"
import DeleteIcon from "@mui/icons-material/Delete"
import { Event } from "../../lib/types"

interface EventCardProps {
  event: Event
  onBook?: (eventId: number) => void
  onConfirmBooking?: (eventId: number) => void
  onEdit?: (eventId: number) => void
  onDelete?: (eventId: number) => void
  currentUserId?: number
  // Новое: информация о статусе брони для текущего пользователя
  bookingStatus?: "pending" | "confirmed" | null
  // Новое: время до истечения брони (в миллисекундах)
  bookingExpiresAt?: number | null
  // Новое: текущее время для расчета истечения брони
  currentTime?: number
}

export default function EventCard({ 
  event, 
  onBook, 
  onConfirmBooking,
  onEdit, 
  onDelete, 
  currentUserId,
  bookingStatus,
  bookingExpiresAt,
  currentTime
}: EventCardProps) {
  // Отладочное логирование
  console.log("EventCard received event data:", event)
  console.log("Event fields:", {
    id: event?.id,
    title: event?.title,
    totalPlaces: event?.totalPlaces,
    availablePlaces: event?.availablePlaces,
    bookingDeadline: event?.bookingDeadline,
    date: event?.date,
    description: event?.description
  })
  console.log("Booking status:", bookingStatus, "Expires at:", bookingExpiresAt, "Current time:", currentTime)

  const handleBookClick = () => {
    if (onBook) {
      onBook(event.id)
    }
  }

  const handleConfirmClick = () => {
    if (onConfirmBooking) {
      onConfirmBooking(event.id)
    }
  }

  const handleEditClick = () => {
    if (onEdit) {
      onEdit(event.id)
    }
  }

  const handleDeleteClick = () => {
    if (onDelete) {
      onDelete(event.id)
    }
  }

  // Проверяем, является ли текущий пользователь владельцем мероприятия
  const isOwner = currentUserId && event.ownerId === currentUserId

  // Проверяем доступность мест
  const isAvailable = event.availablePlaces > 0
  console.log("Available places check:", event.availablePlaces, "isAvailable:", isAvailable)

  // Проверяем, есть ли у пользователя активная бронь на это мероприятие
  const hasActiveBooking = bookingStatus === "pending" || bookingStatus === "confirmed"

  // Только показываем таймер если currentTime передан (т.е. мы в контексте с таймером)
  const showTimer = currentTime !== undefined && bookingStatus === "pending" && bookingExpiresAt

  // Проверяем, истекла ли бронь (только если currentTime доступен)
  const isBookingExpired = showTimer && currentTime > bookingExpiresAt!

  // Вычисляем оставшееся время для отображения (только если currentTime доступен)
  const timeLeft = showTimer 
    ? Math.max(0, Math.ceil((bookingExpiresAt! - currentTime) / 1000))
    : 0

  return (
    <Card sx={{ minWidth: 300, maxWidth: 400, margin: "10px" }}>
      <CardContent>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
          <Typography variant="h5" component="div" sx={{ flex: 1 }}>
            {event.title || 'Без названия'}
          </Typography>
          {isOwner && (
            <Box>
              <IconButton 
                size="small" 
                onClick={handleEditClick}
                sx={{ color: 'primary.main' }}
              >
                <EditIcon fontSize="small" />
              </IconButton>
              <IconButton 
                size="small" 
                onClick={handleDeleteClick}
                sx={{ color: 'error.main' }}
              >
                <DeleteIcon fontSize="small" />
              </IconButton>
            </Box>
          )}
        </Box>
        
        <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
          {event.description || 'Без описания'}
        </Typography>
        
        <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
          Дата: {event.date ? new Date(event.date).toLocaleString("ru-RU") : 'Дата не указана'}
        </Typography>
        
        <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
          Всего мест: {event.totalPlaces ?? 'Не указано'} | Свободных: {event.availablePlaces ?? 'Не указано'}
        </Typography>
        
        <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
          Срок бронирования: {event.bookingDeadline ?? 'Не указано'} мин
        </Typography>
        
        {isOwner && (
          <Typography variant="caption" color="primary" sx={{ mt: 1, display: 'block' }}>
            Вы создали это мероприятие
          </Typography>
        )}
        
        <Chip
          label={isAvailable ? "Доступно" : "Заполнено"}
          color={isAvailable ? "success" : "error"}
          sx={{ mt: 1 }}
        />

        {/* Показываем предупреждение о истекающей брони */}
        {bookingStatus === "pending" && bookingExpiresAt && (
          <Alert severity="warning" sx={{ mt: 1 }}>
            У вас есть бронирование. Подтвердите его до истечения времени.
            <br />
            Осталось: {timeLeft} сек
          </Alert>
        )}

        {bookingStatus === "confirmed" && (
          <Alert severity="success" sx={{ mt: 1 }}>
            Ваше бронирование подтверждено
          </Alert>
        )}

        {isBookingExpired && (
          <Alert severity="error" sx={{ mt: 1 }}>
            Ваше бронирование истекло
          </Alert>
        )}
      </CardContent>
      <CardActions>
        {hasActiveBooking ? (
          <Button
            size="small"
            variant="contained"
            color="warning"
            onClick={handleConfirmClick}
            disabled={!!isBookingExpired}
          >
            Подтвердить бронь
          </Button>
        ) : (
          <Button
            size="small"
            disabled={!isAvailable || !currentUserId}
            onClick={handleBookClick}
          >
            {!currentUserId ? "Войдите для бронирования" : "Забронировать"}
          </Button>
        )}
      </CardActions>
    </Card>
  )
}
