import { Box } from "@mui/material"
import BookingCard from "./BookingCard"
import { Booking } from "../../lib/types"

interface BookingListProps {
  bookings: Booking[]
  bookingErrors?: Record<number, string>
  onConfirm?: (id: number) => void
  onCancel?: (id: number) => void
}

export default function BookingList({
  bookings,
  bookingErrors = {},
  onConfirm,
  onCancel,
}: BookingListProps) {
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
      {bookings.map((booking) => (
        <BookingCard
          key={booking.id}
          booking={booking}
          error={bookingErrors[booking.id] || null}
          onConfirm={onConfirm}
          onCancel={onCancel}
        />
      ))}
    </Box>
  )
}
