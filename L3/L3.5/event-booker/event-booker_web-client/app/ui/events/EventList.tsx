import { Box } from "@mui/material"
import EventCard from "./EventCard"
import { Event } from "../../lib/types"

interface EventListProps {
  events: Event[] | null
}

export default function EventList({ events }: EventListProps) {
  if (!events) {
    return (
      <Box sx={{ width: "100%", textAlign: "center", py: 2 }}>
        Ошибка загрузки мероприятий
      </Box>
    )
  }

  return (
    <Box
      sx={{
        display: "grid",
        gridTemplateColumns: "repeat(auto-fit, minmax(300px, 1fr))",
        gap: 2,
        width: "100%",
        maxWidth: "1200px",
        mx: "auto",
        px: 2,
      }}
    >
      {events.map((event) => (
        <EventCard
          key={event.id}
          event={event}
          onBook={(eventId) => {
            // Здесь будет логика бронирования
            console.log("Book event:", eventId)
          }}
        />
      ))}
    </Box>
  )
}
