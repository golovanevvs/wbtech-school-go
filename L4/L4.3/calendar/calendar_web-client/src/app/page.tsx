"use client"

import { Box, Button, Typography } from "@mui/material"
import { useRouter } from "next/navigation"
import CalendarComponent from "@/ui/CalendarComponent"
import AddIcon from "@mui/icons-material/Add"

export default function Home() {
  const router = useRouter()

  const handleAddEvent = () => {
    router.push("/event/create")
  }

  const handleEventClick = (event: unknown) => {
    // TODO: Реализовать просмотр/редактирование события
    console.log("Event clicked:", event)
  }

  const handleDateClick = (date: Date) => {
    // TODO: Реализовать быстрое создание события по клику на дату
    console.log("Date clicked:", date)
  }

  return (
    <Box
      sx={{
        width: "100%",
        maxWidth: "1200px",
        mx: "auto",
        px: { xs: 1, sm: 2 },
      }}
    >
      {/* Заголовок с кнопкой добавления */}
      <Box
        sx={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          mb: 3,
          flexWrap: "wrap",
          gap: 2,
        }}
      >
        <Typography variant="h4" component="h1">
          Календарь событий
        </Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={handleAddEvent}
          size="large"
        >
          Добавить
        </Button>
      </Box>

      {/* Календарь */}
      <CalendarComponent
        onEventClick={handleEventClick}
        onDateClick={handleDateClick}
      />
    </Box>
  )
}
