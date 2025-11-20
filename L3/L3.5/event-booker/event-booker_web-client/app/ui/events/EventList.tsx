import { Box } from "@mui/material"
import EventCard from "./EventCard"
import { Event } from "../../lib/types"

interface EventListProps {
  events: Event[]
}

export default function EventList({ events }: EventListProps) {
  return (
    <Box
      sx={{
        width: "100%",
        display: "grid",
        gridTemplateColumns: "repeat(auto-fit, minmax(300px, 1fr))",
        gap: 2,
        justifyContent: "center",
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
