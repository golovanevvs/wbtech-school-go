import {
  Card,
  CardContent,
  CardActions,
  Typography,
  Button,
  Chip,
} from "@mui/material"
import { Event } from "../../lib/types"

interface EventCardProps {
  event: Event
  onBook?: (eventId: number) => void
}

export default function EventCard({ event, onBook }: EventCardProps) {
  const handleBookClick = () => {
    if (onBook) {
      onBook(event.id)
    }
  }

  return (
    <Card sx={{ minWidth: 300, maxWidth: 400, margin: "10px" }}>
      <CardContent>
        <Typography variant="h5" component="div">
          {event.title}
        </Typography>
        <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
          {event.description}
        </Typography>
        <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
          Дата: {new Date(event.date).toLocaleString("ru-RU")}
        </Typography>
        <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
          Всего мест: {event.totalPlaces} | Свободных: {event.availablePlaces}
        </Typography>
        <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
          Срок бронирования: {event.bookingDeadline} мин
        </Typography>
        <Chip
          label={event.availablePlaces > 0 ? "Доступно" : "Заполнено"}
          color={event.availablePlaces > 0 ? "success" : "error"}
          sx={{ mt: 1 }}
        />
      </CardContent>
      <CardActions>
        <Button
          size="small"
          disabled={event.availablePlaces === 0}
          onClick={handleBookClick}
        >
          Забронировать
        </Button>
      </CardActions>
    </Card>
  )
}
