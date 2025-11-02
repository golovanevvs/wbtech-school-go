"use client"

import { Box, Link } from "@mui/material"

export default function Footer() {
  return (
    <Box
      component="footer"
      sx={{
        textAlign: "center",
        fontSize: { xs: "0.75rem", sm: "0.875rem" },
        color: "text.secondary",
        mt: 4,
        mb: 2,
      }}
    >
      Designed by{" "}
      <Link
        href="https://t.me/golovanevvs"
        target="_blank"
        rel="noopener noreferrer"
      >
        @golovanevvs
      </Link>
    </Box>
  )
}
