import { Box } from "@mui/material"
import EventCard from "./EventCard"
import { Event } from "../../lib/types"

interface EventListProps {
  events: Event[] | null
  onBook?: (eventId: number) => void
  onConfirmBooking?: (eventId: number) => void
  onCancelBooking?: (eventId: number) => void
  onEdit?: (eventId: number) => void
  onDelete?: (eventId: number) => void
  currentUserId?: number
  bookingsMap?: Record<
    number,
    { status: "pending" | "confirmed" | null; expiresAt?: number | null }
  >
}

export default function EventList({
  events,
  onBook,
  onConfirmBooking,
  onCancelBooking,
  onEdit,
  onDelete,
  currentUserId,
  bookingsMap = {},
}: EventListProps) {
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
      {events.map((event) => {
        const bookingInfo = bookingsMap[event.id] || {
          status: null,
          expiresAt: null,
        }
        return (
          <EventCard
            key={event.id}
            event={event}
            onBook={onBook}
            onConfirmBooking={onConfirmBooking}
            onCancelBooking={onCancelBooking}
            onEdit={onEdit}
            onDelete={onDelete}
            currentUserId={currentUserId}
            bookingStatus={bookingInfo.status}
            bookingExpiresAt={bookingInfo.expiresAt}
          />
        )
      })}
    </Box>
  )
}
