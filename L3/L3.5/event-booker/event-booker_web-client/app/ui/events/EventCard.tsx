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
  onCancelBooking?: (eventId: number) => void
  onEdit?: (eventId: number) => void
  onDelete?: (eventId: number) => void
  currentUserId?: number
  bookingStatus?: "pending" | "confirmed" | null
  bookingExpiresAt?: number | null
}

export default function EventCard({ 
  event, 
  onBook, 
  onConfirmBooking,
  onCancelBooking,
  onEdit, 
  onDelete, 
  currentUserId,
  bookingStatus,
  bookingExpiresAt
}: EventCardProps) {
  

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

  const handleCancelClick = () => {
    if (onCancelBooking) {
      onCancelBooking(event.id)
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
        
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1, mt: 1, alignItems: 'flex-start' }}>
          {isOwner && (
            <Chip
              label="Вы создали это мероприятие"
              color="primary"
            />
          )}
          
          <Chip
            label={isAvailable ? "Доступно" : "Заполнено"}
            color={isAvailable ? "success" : "error"}
          />
        </Box>

        {bookingStatus === "confirmed" && (
          <Alert severity="success" sx={{ mt: 1 }}>
            Ваше бронирование подтверждено
          </Alert>
        )}
      </CardContent>
      <CardActions>
        {bookingStatus === "pending" && (
          <>
            <Button
              size="small"
              color="warning"
              onClick={handleConfirmClick}
            >
              Подтвердить бронь
            </Button>
            <Button
              size="small"
              color="error"
              onClick={handleCancelClick}
            >
              Отменить бронь
            </Button>
          </>
        )}
        
        {!hasActiveBooking && (
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
