import {
  Card,
  CardContent,
  CardActions,
  Typography,
  Button,
  Chip,
  Box,
  IconButton,
} from "@mui/material"
import EditIcon from "@mui/icons-material/Edit"
import DeleteIcon from "@mui/icons-material/Delete"
import { Event } from "../../lib/types"

interface EventCardProps {
  event: Event
  onBook?: (eventId: number) => void
  onEdit?: (eventId: number) => void
  onDelete?: (eventId: number) => void
  currentUserId?: number
}

export default function EventCard({ 
  event, 
  onBook, 
  onEdit, 
  onDelete, 
  currentUserId 
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

  const handleBookClick = () => {
    if (onBook) {
      onBook(event.id)
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
      </CardContent>
      <CardActions>
        <Button
          size="small"
          disabled={!isAvailable}
          onClick={handleBookClick}
        >
          Забронировать
        </Button>
      </CardActions>
    </Card>
  )
}
