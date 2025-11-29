"use client"

import { Suspense } from "react"
import AuthPageContent from "./AuthPageContent"
import { Box, CircularProgress } from "@mui/material"

function AuthPageSkeleton() {
  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        minHeight: "50vh",
        px: { xs: 2, sm: 2 },
      }}
    >
      <CircularProgress />
    </Box>
  )
}

export default function AuthPage() {
  return (
    <Suspense fallback={<AuthPageSkeleton />}>
      <AuthPageContent />
    </Suspense>
  )
}