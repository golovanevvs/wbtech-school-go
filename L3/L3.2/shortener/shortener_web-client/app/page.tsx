"use client"

import { Container, Stack } from "@mui/material"
import URLShortenerForm from "./ui/URLShortenerForm"
import AnalyticsSection from "./ui/AnaliticsSection"
import Header from "./ui/Header"

export default function Home() {
  return (
    <Container
      disableGutters
      sx={{
        width: "100%",
        maxWidth: 600,
        mx: "auto",
        px: { xs: 0, sm: 2 },
        py: 2,
        bgcolor: "background.default",
      }}
    >
      <Stack spacing={4}>
        <Header />
        <URLShortenerForm />
        <AnalyticsSection />
      </Stack>
    </Container>
  )
}
