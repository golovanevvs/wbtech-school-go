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
        minWidth: 250,
        mx: "auto", // центрирование
        bgcolor: "background.default",
        p: 2,
      }}
      >
        <NotifyForm/>
    </Container>
  )
}
