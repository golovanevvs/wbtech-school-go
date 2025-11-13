"use client"

import { Container, Stack } from "@mui/material"
import URLShortenerForm from "./ui/URLShortenerForm"
import AnalyticsSection from "./ui/AnaliticsSection"

export default function Home() {
  return (
    <Container
      disableGutters
      sx={{
        width: "100%",
        maxWidth: 800,
        mx: "auto", // центрирование
        px: { xs: 2, sm: 3 },
        py: 2,
        bgcolor: "background.default",
      }}
    >
      <Stack spacing={4}>
        <URLShortenerForm />
        <AnalyticsSection />
      </Stack>
    </Container>
  )
}
