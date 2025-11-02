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
        bgcolor: "background.default",
        p: { xs: 2, sm: 3 },
      }}
    >
      <NotifyForm />
    </Container>
  )
}
