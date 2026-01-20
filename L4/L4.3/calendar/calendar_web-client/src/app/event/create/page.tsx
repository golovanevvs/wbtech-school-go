"use client"

import { Box } from "@mui/material"
import { useRouter } from "next/navigation"
import EventForm from "@/ui/EventForm"

export default function CreateEventPage() {
  const router = useRouter()

  const handleSuccess = () => {
    router.push("/")
  }

  const handleCancel = () => {
    router.push("/")
  }

  return (
    <Box
      sx={{
        width: "100%",
        maxWidth: "800px",
        mx: "auto",
        px: { xs: 1, sm: 2 },
        py: 3,
      }}
    >
      <EventForm onSuccess={handleSuccess} onCancel={handleCancel} />
    </Box>
  )
}