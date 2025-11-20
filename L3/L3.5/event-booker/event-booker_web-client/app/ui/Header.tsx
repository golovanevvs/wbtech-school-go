"use client"

import { Typography } from "@mui/material"

export default function Header() {
  return (
    <Typography
      variant="h1"
      align="center"
      sx={{
        fontSize: { xs: "1.6rem", sm: "2rem" },
        fontWeight: "bold",
        mb: 4,
      }}
    >
      Event Booker
    </Typography>
  )
}
