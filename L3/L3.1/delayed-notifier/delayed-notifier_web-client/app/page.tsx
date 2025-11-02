"use client"

import { Container } from "@mui/material"
import NotifyForm from "./ui/NotifyForm"

export default function Home() {
  return (
    <Container
      disableGutters
      sx={{
        width: "100%",
        maxWidth: 600,
        mx: "auto", // центрирование
        px: { xs: 0, sm: 2 },
        py: 2,
        bgcolor: "background.default",
      }}
    >
      <NotifyForm />
    </Container>
  )
}
